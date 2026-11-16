package zabbix

import (
	"context"
	"strings"
)

// Hosts searches device names and technical host names, returning at most 50 matches.
func (s *MetricsService) Hosts(ctx context.Context, keyword string) (MonitoredHostList, error) {
	result := MonitoredHostList{Items: []MonitoredHost{}}
	keyword = strings.TrimSpace(keyword)
	if keyword == "" || !validMetricText(keyword, 128) {
		return result, MetricErrorInvalidQuery
	}
	rows, err := queryMetrics(ctx, s, func(ctx context.Context) ([]Host, error) {
		return s.client.Host().Get(ctx, HostGetParams{
			Output:      []string{"hostid", "name", "host"},
			Search:      map[string]any{"name": keyword, "host": keyword},
			SearchByAny: true,
			SortField:   []string{"name"},
			Limit:       51,
		})
	})
	if err != nil {
		return result, err
	}
	result.HasMore = len(rows) > 50
	if len(rows) > 50 {
		rows = rows[:50]
	}
	for _, r := range rows {
		if !metricIdentifier.MatchString(r.HostID) || !validMetricText(r.Name, 512) || !validMetricText(r.Host, 512) {
			return result, MetricErrorInvalidData
		}
		result.Items = append(result.Items, MonitoredHost{ID: r.HostID, Name: r.Name, Host: r.Host})
	}
	return result, nil
}

// Items searches enabled bps/Bps items on one host, returning at most 200 matches.
func (s *MetricsService) Items(ctx context.Context, hostID, keyword string) (MonitoredItemList, error) {
	result := MonitoredItemList{Items: []MonitoredItem{}}
	keyword = strings.TrimSpace(keyword)
	if !metricIdentifier.MatchString(hostID) || !validMetricText(keyword, 128) {
		return result, MetricErrorInvalidQuery
	}
	params := ItemGetParams{
		Output:    metricItemOutput(),
		HostIDs:   []string{hostID},
		Filter:    map[string]any{"value_type": []string{"0", "3"}, "status": "0", "units": []string{"bps", "Bps"}},
		SortField: []string{"name"},
		Limit:     201,
	}
	if keyword != "" {
		params.Search = map[string]any{"name": keyword, "key_": keyword}
		params.SearchByAny = true
	}
	rows, err := queryMetrics(ctx, s, func(ctx context.Context) ([]Item, error) {
		return s.client.Item().Get(ctx, params)
	})
	if err != nil {
		return result, err
	}
	result.HasMore = len(rows) > 200
	if len(rows) > 200 {
		rows = rows[:200]
	}
	for _, r := range rows {
		i := monitoredItem(r)
		if i.IsRateFor(hostID) && metricIdentifier.MatchString(i.ID) && validMetricText(i.Name, 512) && validMetricText(i.Key, 2048) {
			result.Items = append(result.Items, i)
		}
	}
	return result, nil
}

// Host resolves exactly one visible host.
func (s *MetricsService) Host(ctx context.Context, id string) (MonitoredHost, error) {
	var result MonitoredHost
	if !metricIdentifier.MatchString(id) {
		return result, MetricErrorUnavailable
	}
	rows, err := queryMetrics(ctx, s, func(ctx context.Context) ([]Host, error) {
		return s.client.Host().Get(ctx, HostGetParams{
			Output: []string{"hostid", "name", "host"}, HostIDs: []string{id}, Limit: 2,
		})
	})
	if err != nil {
		return result, err
	}
	if len(rows) != 1 || rows[0].HostID != id || !validMetricText(rows[0].Name, 512) || !validMetricText(rows[0].Host, 512) {
		return result, MetricErrorUnavailable
	}
	return MonitoredHost{ID: id, Name: rows[0].Name, Host: rows[0].Host}, nil
}

// Item resolves exactly one visible item.
func (s *MetricsService) Item(ctx context.Context, id string) (MonitoredItem, error) {
	var result MonitoredItem
	if !metricIdentifier.MatchString(id) {
		return result, MetricErrorUnavailable
	}
	rows, err := queryMetrics(ctx, s, func(ctx context.Context) ([]Item, error) {
		return s.client.Item().Get(ctx, ItemGetParams{
			Output: metricItemOutput(), ItemIDs: []string{id}, Limit: 2,
		})
	})
	if err != nil {
		return result, err
	}
	if len(rows) != 1 || rows[0].ItemID != id || !validMetricText(rows[0].Name, 512) {
		return result, MetricErrorUnavailable
	}
	return monitoredItem(rows[0]), nil
}
