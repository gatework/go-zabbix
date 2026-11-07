package zabbix

import "context"

// DHostService wraps the "dhost" API namespace.
type DHostService struct{ client *Client }

// DHost returns the dhost service.
func (c *Client) DHost() *DHostService { return &DHostService{c} }

// DHost object (discovered host).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/dhost/object
type DHost struct {
	DHostID   string     `json:"dhostid,omitempty"`
	DRuleID   string     `json:"druleid,omitempty"`
	Status    string     `json:"status,omitempty"`
	LastUp    string     `json:"lastup,omitempty"`
	LastDown  string     `json:"lastdown,omitempty"`
	DServices []DService `json:"dservices,omitempty"`
}

// DHostGetParams are the parameters for dhost.get.
type DHostGetParams struct {
	GetParams
	DHostIDs        []string `json:"dhostids,omitempty"`
	DRuleIDs        []string `json:"druleids,omitempty"`
	DServiceIDs     []string `json:"dserviceids,omitempty"`
	SelectDServices any      `json:"selectDServices,omitempty"`
}

// Get retrieves discovered hosts matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/dhost/get
func (s *DHostService) Get(ctx context.Context, params DHostGetParams) ([]DHost, error) {
	var out []DHost
	err := s.client.Call(ctx, "dhost.get", params, &out)
	return out, err
}
