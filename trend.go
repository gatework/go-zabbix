package zabbix

import "context"

// TrendService wraps the "trend" API namespace.
type TrendService struct{ client *Client }

// Trend returns the trend service.
func (c *Client) Trend() *TrendService { return &TrendService{c} }

// Trend is an hourly trend value for an item.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/trend/object
type Trend struct {
	ItemID   string `json:"itemid,omitempty"`
	Clock    string `json:"clock,omitempty"`
	Num      string `json:"num,omitempty"`
	ValueMin string `json:"value_min,omitempty"`
	ValueAvg string `json:"value_avg,omitempty"`
	ValueMax string `json:"value_max,omitempty"`
}

// TrendGetParams are the parameters for trend.get.
type TrendGetParams struct {
	GetParams
	ItemIDs  []string `json:"itemids,omitempty"`
	TimeFrom int64    `json:"time_from,omitempty"`
	TimeTill int64    `json:"time_till,omitempty"`
}

// Get retrieves trend values matching params. The trend API is read-only.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/trend/get
func (s *TrendService) Get(ctx context.Context, params TrendGetParams) ([]Trend, error) {
	var out []Trend
	err := s.client.Call(ctx, "trend.get", params, &out)
	return out, err
}
