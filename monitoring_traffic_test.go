package zabbix_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/gatework/go-zabbix"
)

type trafficFixture struct {
	samples   []zabbix.MetricSample
	trends    []zabbix.MetricTrend
	item      zabbix.MonitoredItem
	partial   bool
	itemCalls int
}

func (s *trafficFixture) Host(context.Context, string) (zabbix.MonitoredHost, error) {
	return zabbix.MonitoredHost{ID: "1", Name: "router", Host: "router"}, nil
}

func (s *trafficFixture) Item(context.Context, string) (zabbix.MonitoredItem, error) {
	s.itemCalls++
	return s.item, nil
}

func (s *trafficFixture) History(context.Context, string, string, int64, int64) ([]zabbix.MetricSample, bool, error) {
	return s.samples, s.partial, nil
}

func (s *trafficFixture) Trends(context.Context, string, int64, int64) ([]zabbix.MetricTrend, bool, error) {
	return s.trends, s.partial, nil
}

func TestReadTraffic_GoAdminGolden(t *testing.T) {
	t.Parallel()
	// Captured from go-admin 2d919862 before moving its implementation. The
	// legacy millisecond tuples deliberately differ from the SDK's typed points.
	data, err := os.ReadFile("testdata/traffic_go_admin.json")
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Name       string
		From, Till int64
		Units      string
		Trends     bool
		Samples    []zabbix.MetricSample
		Aggregates []zabbix.MetricTrend
		Points     [][2]*float64
		Peak       *float64
		Count      int
	}
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	if len(cases) != 16 {
		t.Fatalf("golden case count = %d", len(cases))
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			source := &trafficFixture{samples: tc.Samples, trends: tc.Aggregates, item: zabbix.MonitoredItem{
				ID: "10", HostID: "1", Status: "0", ValueType: "0", Units: tc.Units,
			}}
			got, err := zabbix.ReadTraffic(t.Context(), source, zabbix.TrafficQuery{
				HostID: "1", ItemID: "10", From: tc.From, Till: tc.Till, UseTrends: tc.Trends,
			})
			if err != nil {
				t.Fatal(err)
			}
			want := make([]zabbix.TrafficPoint, 0, len(tc.Points))
			for _, point := range tc.Points {
				want = append(want, zabbix.TrafficPoint{Clock: int64(*point[0] / 1000), BPS: point[1]})
			}
			if !reflect.DeepEqual(got.Points, want) || !reflect.DeepEqual(got.PeakBPS, tc.Peak) || got.SampleCount != tc.Count {
				t.Fatalf("traffic changed: points=%d/%d count=%d/%d peak=%v/%v", len(got.Points), len(want), got.SampleCount, tc.Count, got.PeakBPS, tc.Peak)
			}
		})
	}
}

func TestReadTraffic_CurrentOwnershipAndPartial(t *testing.T) {
	t.Parallel()
	source := &trafficFixture{partial: true,
		item:    zabbix.MonitoredItem{ID: "10", HostID: "1", Status: "0", ValueType: "0", Units: "Bps", LastValue: "2", LastClock: 3900},
		samples: []zabbix.MetricSample{{Clock: 3900, Value: "100"}},
	}
	query := zabbix.TrafficQuery{HostID: "1", ItemID: "10", From: 400, Till: 4000}
	got, err := zabbix.ReadTraffic(t.Context(), source, query)
	if err != nil || !got.Partial || got.CurrentBPS == nil || *got.CurrentBPS != 16 || got.PeakBPS == nil || *got.PeakBPS != 800 || got.Stale {
		t.Fatalf("rate result = %+v, %v", got, err)
	}
	source.item.LastClock = 3000
	got, err = zabbix.ReadTraffic(t.Context(), source, query)
	if err != nil || !got.Stale || got.CurrentBPS != nil || got.LastSampleAt == nil || *got.LastSampleAt != 3000 {
		t.Fatalf("stale result = %+v, %v", got, err)
	}
	source.item.HostID = "2"
	if _, err := zabbix.ReadTraffic(t.Context(), source, query); !errors.Is(err, zabbix.ErrMonitoringBinding) {
		t.Fatalf("wrong host error = %v", err)
	}
	query.HostID = "2"
	source.itemCalls = 0
	if _, err := zabbix.ReadTraffic(t.Context(), source, query); !errors.Is(err, zabbix.ErrMonitoringBinding) || source.itemCalls != 0 {
		t.Fatalf("wrong host lookup continued to item query: calls=%d, err=%v", source.itemCalls, err)
	}
}

func TestMonitoring_ZeroValueAndTimeBounds(t *testing.T) {
	t.Parallel()
	var client zabbix.Client
	if _, err := client.Monitoring().Host(t.Context(), "1"); !errors.Is(err, zabbix.ErrMonitoringInput) {
		t.Fatalf("zero client = %v", err)
	}
	for _, bounds := range [][2]int64{{-1, 1}, {2, 1}, {0, 9223372036854776}} {
		t.Run(fmt.Sprint(bounds), func(t *testing.T) {
			t.Parallel()
			if _, _, err := client.Monitoring().History(t.Context(), "1", "0", bounds[0], bounds[1]); !errors.Is(err, zabbix.ErrMonitoringInput) {
				t.Fatalf("invalid bounds = %v", err)
			}
		})
	}
}

func TestReadTraffic_LongTrendsAreBounded(t *testing.T) {
	t.Parallel()
	source := &trafficFixture{item: zabbix.MonitoredItem{ID: "10", HostID: "1", Status: "0", ValueType: "0", Units: "bps"},
		trends: []zabbix.MetricTrend{{Clock: 0, Avg: "10", Max: "100", Count: 1}, {Clock: 3600, Avg: "30", Max: "200", Count: 3}},
	}
	got, err := zabbix.ReadTraffic(t.Context(), source, zabbix.TrafficQuery{HostID: "1", ItemID: "10", From: 0, Till: 9223372036854775, UseTrends: true})
	if err != nil || len(got.Points) > 720 || len(got.Points) == 0 || got.Points[0].BPS == nil || *got.Points[0].BPS != 25 || *got.PeakBPS != 200 {
		t.Fatalf("long trends: points=%d result=%+v err=%v", len(got.Points), got, err)
	}
}

func ExampleReadTraffic() {
	source := &trafficFixture{item: zabbix.MonitoredItem{ID: "10", HostID: "1", Status: "0", ValueType: "3", Units: "Bps"},
		samples: []zabbix.MetricSample{{Clock: 3600, Value: "125"}},
	}
	series, err := zabbix.ReadTraffic(context.Background(), source, zabbix.TrafficQuery{HostID: "1", ItemID: "10", From: 3500, Till: 3700})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*series.PeakBPS, series.Partial)
	// Output: 1000 false
}
