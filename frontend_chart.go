package zabbix

import (
	"bytes"
	"context"
	"image/png"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// ChartRequest selects either item IDs or one existing graph, with an exact
// period. Width and Height are plot dimensions, defaulting to 900 by 240 pixels.
// Zabbix may add room for labels and legends to the resulting PNG.
type ChartRequest struct {
	ItemIDs []string
	GraphID string
	From    time.Time
	To      time.Time
	Width   int
	Height  int
}

// Chart downloads a bounded PNG using the frontend session. It never sends an
// API token. A 200 status alone is not success: HTML, small error images and
// invalid or oversized PNGs are rejected. It does not retry expired sessions.
func (f *Frontend) Chart(ctx context.Context, request ChartRequest) ([]byte, error) {
	hc, _, err := f.snapshot(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := f.sessionID(hc); err != nil {
		return nil, err
	}
	endpoint, err := f.chartURL(&request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, requestError("frontend.chart", ErrorInvalidParams, err)
	}
	req.Header.Set("Accept", "image/png")
	resp, err := hc.Do(req)
	if err != nil {
		return nil, requestError("frontend.chart", ErrorTransport, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden ||
		resp.StatusCode >= 300 && resp.StatusCode < 400 {
		return nil, ErrFrontendAuthentication
	}
	if resp.StatusCode != http.StatusOK {
		return nil, requestError("frontend.chart", ErrorHTTP, &HTTPError{StatusCode: resp.StatusCode})
	}
	mediaType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil || mediaType != "image/png" {
		return nil, ErrFrontendImage
	}
	data, err := readResponseBody("frontend.chart", resp.Body, f.maxBytes)
	if err != nil {
		return nil, err
	}
	if err := validateChartImage(data, request.Width, request.Height); err != nil {
		return nil, err
	}
	return data, nil
}

func (f *Frontend) chartURL(request *ChartRequest) (string, error) {
	from, till := request.From.In(f.location), request.To.In(f.location)
	if (len(request.ItemIDs) == 0) == (request.GraphID == "") || len(request.ItemIDs) > 64 ||
		!validMetricWindow(from.Unix(), till.Unix()) || till.Sub(from) < time.Minute ||
		from.Year() < 1000 || till.Year() > 2999 {
		return "", ErrFrontendInput
	}
	if request.Width == 0 {
		request.Width = 900
	}
	if request.Height == 0 {
		request.Height = 240
	}
	if request.Width < 64 || request.Width > 4096 || request.Height < 64 || request.Height > 2048 {
		return "", ErrFrontendInput
	}
	values := url.Values{
		"from":  {from.Format("2006-01-02 15:04:05")},
		"to":    {till.Format("2006-01-02 15:04:05")},
		"width": {strconv.Itoa(request.Width)}, "height": {strconv.Itoa(request.Height)},
		"legend": {"1"}, "resolve_macros": {"1"},
	}
	name := "chart.php"
	if request.GraphID != "" {
		if !monitoringIdentifier.MatchString(request.GraphID) {
			return "", ErrFrontendInput
		}
		name = "chart2.php"
		values.Set("graphid", request.GraphID)
		values.Set("profileIdx", "web.charts.filter")
	} else {
		for _, id := range request.ItemIDs {
			if !monitoringIdentifier.MatchString(id) {
				return "", ErrFrontendInput
			}
			values.Add("itemids[]", id)
		}
		values.Set("batch", "1")
		values.Set("profileIdx", "web.item.graph.filter")
	}
	// Zabbix 7.4 getTimeSelectorPeriod ignores from/to when profileIdx is absent.
	endpoint := f.endpoint(name)
	endpoint.RawQuery = values.Encode()
	return endpoint.String(), nil
}

func validateChartImage(data []byte, width, height int) error {
	config, err := png.DecodeConfig(bytes.NewReader(data))
	if err != nil || config.Width < width || config.Height < height ||
		config.Width > 8192 || config.Height > 4096 || int64(config.Width)*int64(config.Height) > 16<<20 {
		return ErrFrontendImage
	}
	reader := bytes.NewReader(data)
	if _, err := png.Decode(reader); err != nil || reader.Len() != 0 {
		return ErrFrontendImage
	}
	return nil
}
