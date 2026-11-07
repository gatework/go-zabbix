package zabbix

import "context"

// DHostService wraps the "dhost" API namespace.
type DHostService struct{ client *Client }

// DHost returns the dhost service.
func (c *Client) DHost() *DHostService { return &DHostService{c} }

// DHost object (discovered host).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dhost/object
type DHost struct {
	DHostID   string             `json:"dhostid,omitempty"`
	DRuleID   string             `json:"druleid,omitempty"`
	Status    string             `json:"status,omitempty"`
	LastUp    string             `json:"lastup,omitempty"`
	LastDown  string             `json:"lastdown,omitempty"`
	DServices Relation[DService] `json:"dservices,omitzero"`
}

// DHostGetParams are the parameters for dhost.get.
// DHostGetParams queries the list form of dhost.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type DHostGetParams struct {
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
	DHostIDs               []string       `json:"dhostids,omitzero"`
	DRuleIDs               []string       `json:"druleids,omitzero"`
	DServiceIDs            []string       `json:"dserviceids,omitzero"`
	SelectDServices        any            `json:"selectDServices,omitempty"`
}

// Get retrieves discovered hosts matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dhost/get
func (s *DHostService) Get(ctx context.Context, params DHostGetParams) ([]DHost, error) {
	return getList[DHost](ctx, s.client, "dhost.get", params)
}

// DHostCountParams contains the filters for dhost.get counts.
type DHostCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	DHostIDs               []string       `json:"dhostids,omitzero"`
	DRuleIDs               []string       `json:"druleids,omitzero"`
	DServiceIDs            []string       `json:"dserviceids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *DHostService) Count(ctx context.Context, params DHostCountParams) (uint64, error) {
	return getCount(ctx, s.client, "dhost.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *DHostService) GetMap(ctx context.Context, params DHostGetParams) (map[string]DHost, error) {
	return getMap[DHost](ctx, s.client, "dhost.get", params)
}
