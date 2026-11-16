package zabbix

import (
	"context"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// MetricsService provides bounded device and interface-metric queries.
// A zero service is not usable; obtain it with Client.Metrics.
type MetricsService struct{ client *Client }

// Metrics returns reusable device and traffic monitoring operations.
func (c *Client) Metrics() *MetricsService { return &MetricsService{client: c} }

func queryMetrics[T any](ctx context.Context, s *MetricsService, query func(context.Context) ([]T, error)) ([]T, error) {
	if ctx == nil || s == nil || s.client == nil || s.client.httpClient == nil {
		return nil, MetricErrorInvalidQuery
	}
	return query(ctx)
}

// Unix seconds must also be safe to represent as epoch milliseconds in consumers.
func validMetricWindow(from, till int64) bool {
	return from >= 0 && till >= from && till <= 9223372036854775
}

// MonitoredHost is a validated device identity used for monitoring.
type MonitoredHost struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Host string `json:"host"`
}

// MonitoredHostList is a bounded host search result; HasMore reports truncation.
type MonitoredHostList struct {
	Items   []MonitoredHost `json:"items"`
	HasMore bool            `json:"hasMore"`
}

// MonitoredItem is a numeric-rate candidate with its latest observation.
type MonitoredItem struct {
	ID        string `json:"id"`
	HostID    string `json:"hostId"`
	Name      string `json:"name"`
	Key       string `json:"key"`
	Units     string `json:"units"`
	ValueType string `json:"-"`
	Status    string `json:"-"`
	State     string `json:"-"`
	LastValue string `json:"-"`
	LastClock int64  `json:"-"`
}

// IsRateFor reports whether the enabled numeric item measures a rate on hostID.
// Unsupported state is checked separately when fetching traffic.
func (i MonitoredItem) IsRateFor(hostID string) bool {
	return i.HostID == hostID && (i.ValueType == "0" || i.ValueType == "3") && i.Status == "0" && (i.Units == "bps" || i.Units == "Bps")
}

// MonitoredItemList is a bounded rate-item search result.
type MonitoredItemList struct {
	Items   []MonitoredItem `json:"items"`
	HasMore bool            `json:"hasMore"`
}

// MetricSample is an observed history value at a Unix-second timestamp.
type MetricSample struct {
	Clock int64  `json:"clock"`
	Value string `json:"value"`
}

// MetricTrend is an hourly aggregate; Max is distinct from Avg.
type MetricTrend struct {
	Clock int64  `json:"clock"`
	Avg   string `json:"avg"`
	Max   string `json:"max"`
	Count int64  `json:"count"`
}

func monitoredItem(r Item) MonitoredItem {
	clock, _ := strconv.ParseInt(r.LastClock, 10, 64)
	return MonitoredItem{
		ID: r.ItemID, HostID: r.HostID, Name: r.Name, Key: r.Key, Units: r.Units,
		ValueType: r.ValueType, Status: r.Status, State: r.State, LastValue: r.LastValue, LastClock: clock,
	}
}

var metricIdentifier = regexp.MustCompile(`^[1-9][0-9]{0,31}$`)

func validMetricText(value string, maxLength int) bool {
	return utf8.ValidString(value) && utf8.RuneCountInString(value) <= maxLength && !strings.ContainsRune(value, 0)
}

func metricItemOutput() []string {
	return []string{"itemid", "hostid", "name", "key_", "units", "value_type", "status", "state", "lastvalue", "lastclock"}
}
