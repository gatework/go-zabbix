package zabbix

import (
	"context"
	"strings"
)

// Hosts searches device names and technical host names, returning at most 50 matches.
func (s *MonitoringService) Hosts(ctx context.Context, keyword string) (MonitoredHostList, error) {
	result := MonitoredHostList{Items: []MonitoredHost{}}
	keyword = strings.TrimSpace(keyword)
	if keyword == "" || !validMonitoringText(keyword, 128) {
		return result, ErrMonitoringInput
	}
	rows, err := monitoringQuery(ctx, s, func(ctx context.Context) ([]Host, error) {
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
		if !monitoringIdentifier.MatchString(r.HostID) || !validMonitoringText(r.Name, 512) || !validMonitoringText(r.Host, 512) {
			return result, ErrMonitoringData
		}
		result.Items = append(result.Items, MonitoredHost{ID: r.HostID, Name: r.Name, Host: r.Host})
	}
	return result, nil
}

// Items searches enabled bps/Bps items on one host, returning at most 200 matches.
func (s *MonitoringService) Items(ctx context.Context, hostID, keyword string) (MonitoredItemList, error) {
	result := MonitoredItemList{Items: []MonitoredItem{}}
	keyword = strings.TrimSpace(keyword)
	if !monitoringIdentifier.MatchString(hostID) || !validMonitoringText(keyword, 128) {
		return result, ErrMonitoringInput
	}
	params := ItemGetParams{
		Output:    monitoringItemOutput(),
		HostIDs:   []string{hostID},
		Filter:    map[string]any{"value_type": []string{"0", "3"}, "status": "0", "units": []string{"bps", "Bps"}},
		SortField: []string{"name"},
		Limit:     201,
	}
	if keyword != "" {
		params.Search = map[string]any{"name": keyword, "key_": keyword}
		params.SearchByAny = true
	}
	rows, err := monitoringQuery(ctx, s, func(ctx context.Context) ([]Item, error) {
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
		if i.IsRateFor(hostID) && monitoringIdentifier.MatchString(i.ID) && validMonitoringText(i.Name, 512) && validMonitoringText(i.Key, 2048) {
			result.Items = append(result.Items, i)
		}
	}
	return result, nil
}

// Host resolves exactly one visible host.
func (s *MonitoringService) Host(ctx context.Context, id string) (MonitoredHost, error) {
	var result MonitoredHost
	if !monitoringIdentifier.MatchString(id) {
		return result, ErrMonitoringBinding
	}
	rows, err := monitoringQuery(ctx, s, func(ctx context.Context) ([]Host, error) {
		return s.client.Host().Get(ctx, HostGetParams{
			Output: []string{"hostid", "name", "host"}, HostIDs: []string{id}, Limit: 2,
		})
	})
	if err != nil {
		return result, err
	}
	if len(rows) != 1 || rows[0].HostID != id || !validMonitoringText(rows[0].Name, 512) || !validMonitoringText(rows[0].Host, 512) {
		return result, ErrMonitoringBinding
	}
	return MonitoredHost{ID: id, Name: rows[0].Name, Host: rows[0].Host}, nil
}

// Item resolves exactly one visible item.
func (s *MonitoringService) Item(ctx context.Context, id string) (MonitoredItem, error) {
	var result MonitoredItem
	if !monitoringIdentifier.MatchString(id) {
		return result, ErrMonitoringBinding
	}
	rows, err := monitoringQuery(ctx, s, func(ctx context.Context) ([]Item, error) {
		return s.client.Item().Get(ctx, ItemGetParams{
			Output: monitoringItemOutput(), ItemIDs: []string{id}, Limit: 2,
		})
	})
	if err != nil {
		return result, err
	}
	if len(rows) != 1 || rows[0].ItemID != id || !validMonitoringText(rows[0].Name, 512) {
		return result, ErrMonitoringBinding
	}
	return monitoredItem(rows[0]), nil
}
