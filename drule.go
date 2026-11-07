package zabbix

import "context"

// DRuleService wraps the "drule" API namespace.
type DRuleService struct{ client *Client }

// DRule returns the drule service.
func (c *Client) DRule() *DRuleService { return &DRuleService{c} }

// DRule object (network discovery rule).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/drule/object
type DRule struct {
	DRuleID        string           `json:"druleid,omitempty"`
	Name           string           `json:"name,omitempty"`
	IPRange        string           `json:"iprange,omitempty"`
	Delay          string           `json:"delay,omitempty"`
	Status         string           `json:"status,omitempty"`
	ProxyID        string           `json:"proxyid,omitempty"`
	ConcurrencyMax string           `json:"concurrency_max,omitempty"`
	Error          string           `json:"error,omitempty"`
	DChecks        Relation[DCheck] `json:"dchecks,omitzero"`
	DHosts         Relation[DHost]  `json:"dhosts,omitzero"`
}

// DRuleGetParams are the parameters for drule.get.
// DRuleGetParams queries the list form of drule.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type DRuleGetParams struct {
	Output                 any            `json:"output,omitempty"`
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	SortField              []string       `json:"sortfield,omitzero"`
	SortOrder              []string       `json:"sortorder,omitzero"`
	Limit                  int            `json:"limit,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	DRuleIDs               []string       `json:"druleids,omitzero"`
	DHostIDs               []string       `json:"dhostids,omitzero"`
	DServiceIDs            []string       `json:"dserviceids,omitzero"`
	SelectDChecks          any            `json:"selectDChecks,omitempty"`
	SelectDHosts           any            `json:"selectDHosts,omitempty"`
}

// Get retrieves network discovery rules matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/drule/get
func (s *DRuleService) Get(ctx context.Context, params DRuleGetParams) ([]DRule, error) {
	return getList[DRule](ctx, s.client, "drule.get", params)
}

// DRuleCountParams contains the filters for drule.get counts.
type DRuleCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	DRuleIDs               []string       `json:"druleids,omitzero"`
	DHostIDs               []string       `json:"dhostids,omitzero"`
	DServiceIDs            []string       `json:"dserviceids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *DRuleService) Count(ctx context.Context, params DRuleCountParams) (uint64, error) {
	return getCount(ctx, s.client, "drule.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *DRuleService) GetMap(ctx context.Context, params DRuleGetParams) (map[string]DRule, error) {
	return getMap[DRule](ctx, s.client, "drule.get", params)
}

// Create creates discovery rules and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/drule/create
func (s *DRuleService) Create(ctx context.Context, rules ...DRuleCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "drule.create", "druleids", rules)
}

// Update updates discovery rules and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/drule/update
func (s *DRuleService) Update(ctx context.Context, rules ...DRuleUpdateParams) ([]string, error) {
	return updateIDs(ctx, s.client, "drule.update", "druleids", rules)
}

// Delete deletes discovery rules by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/drule/delete
func (s *DRuleService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "drule.delete", "druleids", ids)
}
