package zabbix

import "context"

// DServiceService wraps the "dservice" API namespace.
type DServiceService struct{ client *Client }

// DService returns the dservice service.
func (c *Client) DService() *DServiceService { return &DServiceService{c} }

// DService object (discovered service).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/dservice/object
type DService struct {
	DServiceID string `json:"dserviceid,omitempty"`
	DHostID    string `json:"dhostid,omitempty"`
	Type       string `json:"type,omitempty"`
	Key        string `json:"key_,omitempty"`
	Value      string `json:"value,omitempty"`
	Port       string `json:"port,omitempty"`
	Status     string `json:"status,omitempty"`
	LastUp     string `json:"lastup,omitempty"`
	LastDown   string `json:"lastdown,omitempty"`
	DCheckID   string `json:"dcheckid,omitempty"`
	IP         string `json:"ip,omitempty"`
	DNS        string `json:"dns,omitempty"`
}

// DServiceGetParams are the parameters for dservice.get.
type DServiceGetParams struct {
	GetParams
	DServiceIDs  []string `json:"dserviceids,omitempty"`
	DHostIDs     []string `json:"dhostids,omitempty"`
	DRuleIDs     []string `json:"druleids,omitempty"`
	DCheckIDs    []string `json:"dcheckids,omitempty"`
	SelectDRules any      `json:"selectDRules,omitempty"`
	SelectDHosts any      `json:"selectDHosts,omitempty"`
	SelectHosts  any      `json:"selectHosts,omitempty"`
}

// Get retrieves discovered services matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/dservice/get
func (s *DServiceService) Get(ctx context.Context, params DServiceGetParams) ([]DService, error) {
	var out []DService
	err := s.client.Call(ctx, "dservice.get", params, &out)
	return out, err
}
