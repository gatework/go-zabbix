package zabbix

import (
	"math"
	"slices"
	"sort"
	"strconv"
)

func rateValue(value string, multiplier float64) *float64 {
	v, err := strconv.ParseFloat(value, 64)
	v *= multiplier
	if err != nil || math.IsNaN(v) || math.IsInf(v, 0) || v < 0 {
		return nil
	}
	return &v
}
func recordPeak(peak **float64, value *float64) {
	if value != nil && (*peak == nil || **peak < *value) {
		observed := *value
		*peak = &observed
	}
}
func trafficPoint(clock int64, value *float64) TrafficPoint {
	return TrafficPoint{Clock: clock, BPS: value}
}

// Buckets average only actual samples. Missing normal polling intervals are
// omitted; a gap larger than max(5 minutes, 3 x median observed interval) gets
// a null marker. Each marker consumes an otherwise empty bucket, so even with
// gaps the final chart contains at most 720 points.
func historyPoints(rows []MetricSample, from, till int64, multiplier float64) ([]TrafficPoint, *float64, int) {
	points := []TrafficPoint{}
	var peak *float64
	valid := 0
	filtered := make([]MetricSample, 0, len(rows))
	for _, row := range rows {
		if row.Clock >= from && row.Clock <= till {
			filtered = append(filtered, row)
		}
	}
	if len(filtered) == 0 {
		return points, nil, 0
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].Clock < filtered[j].Clock })
	gap := historyGap(filtered)
	width := max(int64(1), (till-from+719)/720)
	type bucket struct {
		mean        float64
		count       int
		invalid     bool
		first, last int64
		seen        bool
	}
	buckets := make([]bucket, 720)
	for _, row := range filtered {
		index := min(int((row.Clock-from)/width), 719)
		b := &buckets[index]
		if !b.seen {
			b.first = row.Clock
			b.seen = true
		}
		b.last = row.Clock
		value := rateValue(row.Value, multiplier)
		if value == nil {
			b.invalid = true
			continue
		}
		b.count++
		b.mean += (*value - b.mean) / float64(b.count)
		valid++
		recordPeak(&peak, value)
	}
	lastClock, lastIndex := from, -1
	for index, b := range buckets {
		if !b.seen {
			continue
		}
		if b.first-lastClock > gap && index > lastIndex+1 {
			points = append(points, trafficPoint(from+int64(lastIndex+1)*width, nil))
		}
		var value *float64
		if b.count > 0 && !b.invalid {
			v := b.mean
			value = &v
		}
		points = append(points, trafficPoint(b.first, value))
		lastClock = b.last
		lastIndex = index
	}
	if till-lastClock > gap && lastIndex < 719 {
		points = append(points, trafficPoint(till, nil))
	}
	return points, peak, valid
}

func historyGap(rows []MetricSample) int64 {
	intervals := []int64{}
	for i := 1; i < len(rows); i++ {
		if delta := rows[i].Clock - rows[i-1].Clock; delta > 0 {
			intervals = append(intervals, delta)
		}
	}
	if len(intervals) == 0 {
		return 300
	}
	slices.Sort(intervals)
	return max(300, 3*intervals[len(intervals)/2])
}
func trendPoints(rows []MetricTrend, from, till int64, multiplier float64) ([]TrafficPoint, *float64, int) {
	points := []TrafficPoint{}
	var peak *float64
	valid := 0
	values := map[int64]*float64{}
	for _, row := range rows {
		if row.Clock < from || row.Clock > till || row.Count <= 0 {
			continue
		}
		avg := rateValue(row.Avg, multiplier)
		maximum := rateValue(row.Max, multiplier)
		values[row.Clock] = avg
		if avg != nil {
			valid++
		}
		recordPeak(&peak, maximum)
	}
	if len(values) == 0 {
		return points, nil, 0
	}
	first := (from + 3599) / 3600 * 3600
	step := max(int64(1), ((till-first)/3600+1+719)/720) * 3600
	if step > 3600 {
		values = bucketTrends(rows, first, till, step, multiplier)
	}
	for clock := first; clock <= till; clock += step {
		points = append(points, trafficPoint(clock, values[clock]))
	}
	return points, peak, valid
}

// Long periods stay bounded without treating missing hours as observed zero.
// Hourly means are weighted by num; the observed peak is retained separately.
func bucketTrends(rows []MetricTrend, from, till, step int64, multiplier float64) map[int64]*float64 {
	type bucket struct {
		mean, weight float64
		invalid      bool
	}
	buckets := map[int64]bucket{}
	for _, row := range rows {
		if row.Clock < from || row.Clock > till || row.Count <= 0 {
			continue
		}
		clock := from + (row.Clock-from)/step*step
		b := buckets[clock]
		value := rateValue(row.Avg, multiplier)
		if value == nil {
			b.invalid = true
		} else {
			weight := float64(row.Count)
			b.weight += weight
			b.mean += (*value - b.mean) * (weight / b.weight)
		}
		buckets[clock] = b
	}
	values := make(map[int64]*float64, len(buckets))
	for clock, b := range buckets {
		if b.weight > 0 && !b.invalid {
			value := b.mean
			values[clock] = &value
		}
	}
	return values
}
