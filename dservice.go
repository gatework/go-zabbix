package zabbix

import "context"

// DServiceService wraps the "dservice" API namespace.
type DServiceService struct{ client *Client }

// DService returns the dservice service.
func (c *Client) DService() *DServiceService { return &DServiceService{c} }

// DService object (discovered service).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dservice/object
type DService struct {
	DServiceID string         `json:"dserviceid,omitempty"`
	DHostID    string         `json:"dhostid,omitempty"`
	Type       string         `json:"type,omitempty"`
	Key        string         `json:"key_,omitempty"`
	Value      string         `json:"value,omitempty"`
	Port       string         `json:"port,omitempty"`
	Status     string         `json:"status,omitempty"`
	LastUp     string         `json:"lastup,omitempty"`
	LastDown   string         `json:"lastdown,omitempty"`
	DCheckID   string         `json:"dcheckid,omitempty"`
	IP         string         `json:"ip,omitempty"`
	DNS        string         `json:"dns,omitempty"`
	DRules     []DRule        `json:"drules,omitempty"`
	DHosts     []DHost        `json:"dhosts,omitempty"`
	Hosts      Relation[Host] `json:"hosts,omitzero"`
}

// DServiceGetParams are the parameters for dservice.get.
// DServiceGetParams queries the list form of dservice.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type DServiceGetParams struct {
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
	DServiceIDs            []string       `json:"dserviceids,omitzero"`
	DHostIDs               []string       `json:"dhostids,omitzero"`
	DRuleIDs               []string       `json:"druleids,omitzero"`
	DCheckIDs              []string       `json:"dcheckids,omitzero"`
	SelectDRules           any            `json:"selectDRules,omitempty"`
	SelectDHosts           any            `json:"selectDHosts,omitempty"`
	SelectHosts            any            `json:"selectHosts,omitempty"`
}

// Get retrieves discovered services matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dservice/get
func (s *DServiceService) Get(ctx context.Context, params DServiceGetParams) ([]DService, error) {
	return getList[DService](ctx, s.client, "dservice.get", params)
}

// DServiceCountParams contains the filters for dservice.get counts.
type DServiceCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	DServiceIDs            []string       `json:"dserviceids,omitzero"`
	DHostIDs               []string       `json:"dhostids,omitzero"`
	DRuleIDs               []string       `json:"druleids,omitzero"`
	DCheckIDs              []string       `json:"dcheckids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *DServiceService) Count(ctx context.Context, params DServiceCountParams) (uint64, error) {
	return getCount(ctx, s.client, "dservice.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *DServiceService) GetMap(ctx context.Context, params DServiceGetParams) (map[string]DService, error) {
	return getMap[DService](ctx, s.client, "dservice.get", params)
}
