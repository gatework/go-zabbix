package zabbix

import "context"

// TrendService wraps the "trend" API namespace.
type TrendService struct{ client *Client }

// Trend returns the trend service.
func (c *Client) Trend() *TrendService { return &TrendService{c} }

// Trend is an hourly trend value for an item.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/trend/object
type Trend struct {
	ItemID   string `json:"itemid,omitempty"`
	Clock    string `json:"clock,omitempty"`
	Num      string `json:"num,omitempty"`
	ValueMin string `json:"value_min,omitempty"`
	ValueAvg string `json:"value_avg,omitempty"`
	ValueMax string `json:"value_max,omitempty"`
}

// TrendGetParams are the parameters for trend.get.
// TrendGetParams queries the list form of trend.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type TrendGetParams struct {
	Output   any      `json:"output,omitempty"`
	Limit    int      `json:"limit,omitempty"`
	ItemIDs  []string `json:"itemids,omitzero"`
	TimeFrom *int64   `json:"time_from,omitempty"`
	TimeTill *int64   `json:"time_till,omitempty"`
}

// Get retrieves trend values matching params. The trend API is read-only.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/trend/get
func (s *TrendService) Get(ctx context.Context, params TrendGetParams) ([]Trend, error) {
	return getList[Trend](ctx, s.client, "trend.get", params)
}

// TrendCountParams contains the filters for trend.get counts.
type TrendCountParams struct {
	ItemIDs  []string `json:"itemids,omitzero"`
	TimeFrom *int64   `json:"time_from,omitempty"`
	TimeTill *int64   `json:"time_till,omitempty"`
}

// Count returns the number of objects matching params.
func (s *TrendService) Count(ctx context.Context, params TrendCountParams) (uint64, error) {
	return getCount(ctx, s.client, "trend.get", params)
}
