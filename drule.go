package zabbix

import "context"

// DRuleService wraps the "drule" API namespace.
type DRuleService struct{ client *Client }

// DRule returns the drule service.
func (c *Client) DRule() *DRuleService { return &DRuleService{c} }

// DRule object (network discovery rule).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/drule/object
type DRule struct {
	DRuleID        string   `json:"druleid,omitempty"`
	Name           string   `json:"name,omitempty"`
	IPRange        string   `json:"iprange,omitempty"`
	Delay          string   `json:"delay,omitempty"`
	Status         string   `json:"status,omitempty"`
	ProxyID        string   `json:"proxyid,omitempty"`
	ConcurrencyMax string   `json:"concurrency_max,omitempty"`
	Error          string   `json:"error,omitempty"`
	DChecks        []DCheck `json:"dchecks,omitempty"`
}

// DRuleGetParams are the parameters for drule.get.
type DRuleGetParams struct {
	GetParams
	DRuleIDs      []string `json:"druleids,omitempty"`
	DHostIDs      []string `json:"dhostids,omitempty"`
	DServiceIDs   []string `json:"dserviceids,omitempty"`
	SelectDChecks any      `json:"selectDChecks,omitempty"`
	SelectDHosts  any      `json:"selectDHosts,omitempty"`
}

// Get retrieves network discovery rules matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/drule/get
func (s *DRuleService) Get(ctx context.Context, params DRuleGetParams) ([]DRule, error) {
	var out []DRule
	err := s.client.Call(ctx, "drule.get", params, &out)
	return out, err
}

// Create creates discovery rules and returns the new ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/drule/create
func (s *DRuleService) Create(ctx context.Context, rules ...DRule) ([]string, error) {
	var res struct {
		IDs []string `json:"druleids"`
	}
	err := s.client.Call(ctx, "drule.create", rules, &res)
	return res.IDs, err
}

// Update updates discovery rules and returns the affected ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/drule/update
func (s *DRuleService) Update(ctx context.Context, rules ...DRule) ([]string, error) {
	var res struct {
		IDs []string `json:"druleids"`
	}
	err := s.client.Call(ctx, "drule.update", rules, &res)
	return res.IDs, err
}

// Delete deletes discovery rules by id and returns the deleted ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/drule/delete
func (s *DRuleService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"druleids"`
	}
	err := s.client.Call(ctx, "drule.delete", ids, &res)
	return res.IDs, err
}
