//go:build integration

package zabbix

import (
	"bytes"
	"context"
	"errors"
	"image/png"
	"net/url"
	"testing"
	"time"
)

// Uses only the disposable host/items created by this suite. Sampling and
// rendering are checked independently so a valid empty plot cannot pass both.
func TestIntegrationFrontendCharts(t *testing.T) {
	c, ctx := loginClient(t)
	hostID, _, _ := integrationHost(t, c, ctx)
	itemIDs, err := c.Item().Create(ctx, ItemCreateParams{
		HostID: hostID, Name: "SDK interface traffic", Key: "sdk.traffic", Type: 2,
		ValueType: 3, Units: integrationValue("bps"),
	})
	itemID := integrationID(t, itemIDs, err)
	graphIDs, err := c.Graph().Create(ctx, GraphCreateParams{
		Name: "SDK interface traffic", Width: 900, Height: 240,
		GItems: []GraphItemInput{{ItemID: itemID, Color: "00AA00"}},
	})
	graphID := integrationID(t, graphIDs, err)
	item, err := c.Monitoring().Item(ctx, itemID)
	integrationCheck(t, err)
	if !item.IsRateFor(hostID) {
		t.Fatal("monitoring item lost its host/rate contract")
	}
	till := time.Now().UTC().Truncate(time.Second)
	integrationTrafficSample(t, c, ctx, hostID, itemID, till)
	endpoint, err := url.Parse(*zbxURL)
	integrationCheck(t, err)
	frontend, err := NewFrontend(endpoint.ResolveReference(&url.URL{Path: "./"}).String(), FrontendOptions{Timezone: time.UTC})
	integrationCheck(t, err)
	t.Cleanup(frontend.CloseIdleConnections)
	integrationCheck(t, frontend.Login(ctx, *zbxUser, *zbxPass))
	loggedOut := false
	t.Cleanup(func() {
		if loggedOut {
			return
		}
		cleanup, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()
		if err := frontend.Logout(cleanup); err != nil {
			t.Errorf("frontend cleanup: %v", err)
		}
	})
	for _, request := range []ChartRequest{
		{ItemIDs: []string{itemID}, From: till.Add(-time.Hour), To: till},
		{GraphID: graphID, From: till.Add(-time.Hour), To: till},
	} {
		data, err := frontend.Chart(ctx, request)
		integrationCheck(t, err)
		config, err := png.DecodeConfig(bytes.NewReader(data))
		integrationCheck(t, err)
		if config.Width < 900 || config.Height < 240 {
			t.Fatalf("unexpected plot dimensions: %+v", config)
		}
	}
	integrationCheck(t, frontend.Logout(ctx))
	loggedOut = true
	_, err = frontend.Chart(ctx, ChartRequest{ItemIDs: []string{itemID}, From: till.Add(-time.Hour), To: till})
	if !errors.Is(err, ErrFrontendAuthentication) {
		t.Fatalf("logged-out chart = %v", err)
	}
}

func integrationTrafficSample(t *testing.T, client *Client, ctx context.Context, hostID, itemID string, till time.Time) {
	t.Helper()
	poll, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	pushed := false
	for {
		if !pushed {
			var response struct {
				Response string `json:"response"`
				Data     []struct {
					Error string `json:"error"`
				} `json:"data"`
			}
			// A newly created trapper must first reach the server configuration cache.
			// The raw call is fixture setup, not a new public write API.
			err := client.Call(poll, "history.push", []map[string]any{{
				"itemid": itemID, "clock": till.Add(-time.Minute).Unix(), "value": 12345,
			}}, &response)
			integrationCheck(t, err)
			if response.Response != "success" || len(response.Data) != 1 {
				t.Fatal("unexpected history.push receipt")
			}
			pushed = response.Data[0].Error == ""
		}
		if pushed {
			series, err := client.Monitoring().Traffic(poll, TrafficQuery{
				HostID: hostID, ItemID: itemID, From: till.Add(-time.Hour).Unix(), Till: till.Unix(),
			})
			integrationCheck(t, err)
			if series.SampleCount == 1 && !series.Partial && series.CurrentBPS != nil &&
				*series.CurrentBPS == 12345 && series.PeakBPS != nil && *series.PeakBPS == 12345 {
				return
			}
		}
		select {
		case <-poll.Done():
			t.Fatalf("traffic sample did not become visible: %v", poll.Err())
		case <-ticker.C:
		}
	}
}
