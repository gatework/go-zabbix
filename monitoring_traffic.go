package zabbix

import "context"

// MonitoringLookup resolves host and item identities before reading traffic.
type MonitoringLookup interface {
	Host(ctx context.Context, id string) (MonitoredHost, error)
	Item(ctx context.Context, id string) (MonitoredItem, error)
}

// MetricReader reads bounded history and trend observations, including truncation.
type MetricReader interface {
	History(ctx context.Context, id, valueType string, from, till int64) ([]MetricSample, bool, error)
	Trends(ctx context.Context, id string, from, till int64) ([]MetricTrend, bool, error)
}

// TrafficSource supplies the read-only capabilities needed to resolve a rate series.
type TrafficSource interface {
	MonitoringLookup
	MetricReader
}

var _ TrafficSource = (*MonitoringService)(nil)

// TrafficQuery selects one item's rate on a particular host over Unix seconds.
// UseTrends requests hourly aggregates; otherwise individual history is used.
// All bounds are explicit so retries and alert attachments use the same period.
type TrafficQuery struct {
	HostID    string
	ItemID    string
	From      int64
	Till      int64
	UseTrends bool
}

// TrafficPoint contains a Unix-second timestamp and a nullable bits/second rate.
// A nil rate is missing or invalid evidence, never an observed zero.
type TrafficPoint struct {
	Clock int64    `json:"clock"`
	BPS   *float64 `json:"bps"`
}

// TrafficSeries retains the identities, observed peak and completeness of a query.
// PeakBPS is a peak among received observations, not proof of the full-period peak
// when Partial is true. LastSampleAt is in Unix seconds and independent of trends.
type TrafficSeries struct {
	Host         MonitoredHost  `json:"host"`
	Item         MonitoredItem  `json:"item"`
	Points       []TrafficPoint `json:"points"`
	CurrentBPS   *float64       `json:"current_bps"`
	PeakBPS      *float64       `json:"peak_bps"`
	LastSampleAt *int64         `json:"last_sample_at"`
	SampleCount  int            `json:"sample_count"`
	Partial      bool           `json:"partial"`
	Stale        bool           `json:"stale"`
}

// Traffic revalidates item ownership and returns a bounded, normalized rate series.
func (s *MonitoringService) Traffic(ctx context.Context, query TrafficQuery) (TrafficSeries, error) {
	return ReadTraffic(ctx, s, query)
}

// ReadTraffic composes identity lookup and metric reading without application
// database or authorization dependencies. It preserves source error chains.
// Callers enforce their own circuit/device permissions before invoking it.
func ReadTraffic(ctx context.Context, source TrafficSource, query TrafficQuery) (TrafficSeries, error) {
	result := TrafficSeries{Points: []TrafficPoint{}}
	if ctx == nil || source == nil || !validMetricWindow(query.From, query.Till) ||
		!monitoringIdentifier.MatchString(query.HostID) || !monitoringIdentifier.MatchString(query.ItemID) {
		return result, ErrMonitoringInput
	}
	if err := ctx.Err(); err != nil {
		return result, err
	}
	host, err := source.Host(ctx, query.HostID)
	if err != nil {
		return result, err
	}
	if host.ID != query.HostID {
		return result, ErrMonitoringBinding
	}
	item, err := source.Item(ctx, query.ItemID)
	if err != nil {
		return result, err
	}
	if item.ID != query.ItemID || !item.IsRateFor(query.HostID) || item.State == "1" {
		return result, ErrMonitoringBinding
	}
	result.Host, result.Item = host, item
	result.setLatest(query.Till)
	if err := result.readSamples(ctx, source, query); err != nil {
		return result, err
	}
	return result, nil
}

func (s *TrafficSeries) setLatest(till int64) {
	clock := s.Item.LastClock
	if clock > 0 && clock <= till {
		s.LastSampleAt = &clock
		if till-clock <= 300 {
			s.CurrentBPS = rateValue(s.Item.LastValue, s.multiplier())
		}
	}
	s.Stale = clock > 0 && till-clock > 300
}

func (s *TrafficSeries) readSamples(ctx context.Context, source MetricReader, query TrafficQuery) error {
	if query.UseTrends {
		rows, partial, err := source.Trends(ctx, s.Item.ID, query.From, query.Till)
		s.Partial = partial
		if err != nil {
			return err
		}
		s.Points, s.PeakBPS, s.SampleCount = trendPoints(rows, query.From, query.Till, s.multiplier())
		return nil
	}
	rows, partial, err := source.History(ctx, s.Item.ID, s.Item.ValueType, query.From, query.Till)
	s.Partial = partial
	if err != nil {
		return err
	}
	s.Points, s.PeakBPS, s.SampleCount = historyPoints(rows, query.From, query.Till, s.multiplier())
	return nil
}

func (s *TrafficSeries) multiplier() float64 {
	if s.Item.Units == "Bps" {
		return 8
	}
	return 1
}
