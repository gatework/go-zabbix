package zabbix

import (
	"context"
	"strconv"
)

// History reads up to 12001 observations, retaining the truncation probe.
// Partial is true above 12000; callers must not present the range as complete.
func (s *MonitoringService) History(ctx context.Context, id, valueType string, from, till int64) ([]MetricSample, bool, error) {
	if !validMetricWindow(from, till) {
		return nil, false, ErrMonitoringInput
	}
	if !monitoringIdentifier.MatchString(id) || (valueType != "0" && valueType != "3") {
		return nil, false, ErrMonitoringBinding
	}
	history := 0
	if valueType == "3" {
		history = 3
	}
	rows, err := monitoringQuery(ctx, s, func(ctx context.Context) ([]HistoryValue, error) {
		return s.client.History().Get(ctx, HistoryGetParams{
			Output: []string{"clock", "value"}, ItemIDs: []string{id}, History: &history,
			TimeFrom: &from, TimeTill: &till,
			SortField: []string{"clock"}, SortOrder: []string{"DESC"}, Limit: 12001,
		})
	})
	if err != nil {
		return nil, false, err
	}
	partial := len(rows) > 12000
	// The probe is an observed sample too: keep it when calculating the
	// reported peak, while explicitly marking the range as partial.
	if len(rows) > 12001 {
		rows = rows[:12001]
	}
	result := make([]MetricSample, 0, len(rows))
	for _, r := range rows {
		clock, err := strconv.ParseInt(r.Clock, 10, 64)
		if err != nil || clock < 0 {
			return nil, partial, ErrMonitoringData
		}
		result = append(result, MetricSample{Clock: clock, Value: r.Value})
	}
	return result, partial, nil
}

// Trends reads up to 201 hourly aggregates, retaining the truncation probe.
// Results are not ordered by Zabbix; Partial is true above 200.
func (s *MonitoringService) Trends(ctx context.Context, id string, from, till int64) ([]MetricTrend, bool, error) {
	if !validMetricWindow(from, till) {
		return nil, false, ErrMonitoringInput
	}
	if !monitoringIdentifier.MatchString(id) {
		return nil, false, ErrMonitoringBinding
	}
	// Zabbix 7.4 does not sort trend.get. The circuit view queries one item over
	// seven days; the extra row detects an incomplete response if that bound is
	// exceeded. The service orders the observed points before presentation.
	rows, err := monitoringQuery(ctx, s, func(ctx context.Context) ([]Trend, error) {
		return s.client.Trend().Get(ctx, TrendGetParams{
			Output: []string{"clock", "num", "value_avg", "value_max"}, ItemIDs: []string{id},
			TimeFrom: &from, TimeTill: &till, Limit: 201,
		})
	})
	if err != nil {
		return nil, false, err
	}
	partial := len(rows) > 200
	// Keep the observed probe when calculating peaks, just as history does.
	if len(rows) > 201 {
		rows = rows[:201]
	}
	result := make([]MetricTrend, 0, len(rows))
	for _, r := range rows {
		clock, err := strconv.ParseInt(r.Clock, 10, 64)
		count, e := strconv.ParseInt(r.Num, 10, 64)
		if err != nil || e != nil || clock < 0 || count < 0 {
			return nil, partial, ErrMonitoringData
		}
		result = append(result, MetricTrend{Clock: clock, Avg: r.ValueAvg, Max: r.ValueMax, Count: count})
	}
	return result, partial, nil
}
