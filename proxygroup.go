package zabbix

import "context"

// ProxyGroupService wraps the "proxygroup" API namespace.
type ProxyGroupService struct{ client *Client }

// ProxyGroup returns the proxygroup service.
func (c *Client) ProxyGroup() *ProxyGroupService { return &ProxyGroupService{c} }

// ProxyGroup object (Zabbix 7.0+).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxygroup/object
type ProxyGroup struct {
	ProxyGroupID  string          `json:"proxy_groupid,omitempty"`
	Name          string          `json:"name,omitempty"`
	Description   string          `json:"description,omitempty"`
	FailoverDelay string          `json:"failover_delay,omitempty"`
	MinOnline     string          `json:"min_online,omitempty"`
	State         string          `json:"state,omitempty"`
	Proxies       Relation[Proxy] `json:"proxies,omitzero"`
}

// ProxyGroupProxy is a proxy that is a member of a proxy group.
type ProxyGroupProxy struct {
	ProxyID string `json:"proxyid,omitempty"`
	Name    string `json:"name,omitempty"`
}

// ProxyGroupGetParams are the parameters for proxygroup.get.
// ProxyGroupGetParams queries the list form of proxygroup.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type ProxyGroupGetParams struct {
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
	ProxyGroupIDs          []string       `json:"proxy_groupids,omitzero"`
	SelectProxies          any            `json:"selectProxies,omitempty"`
}

// Get retrieves proxy groups matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxygroup/get
func (s *ProxyGroupService) Get(ctx context.Context, params ProxyGroupGetParams) ([]ProxyGroup, error) {
	return getList[ProxyGroup](ctx, s.client, "proxygroup.get", params)
}

// ProxyGroupCountParams contains the filters for proxygroup.get counts.
type ProxyGroupCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	ProxyGroupIDs          []string       `json:"proxy_groupids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *ProxyGroupService) Count(ctx context.Context, params ProxyGroupCountParams) (uint64, error) {
	return getCount(ctx, s.client, "proxygroup.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *ProxyGroupService) GetMap(ctx context.Context, params ProxyGroupGetParams) (map[string]ProxyGroup, error) {
	return getMap[ProxyGroup](ctx, s.client, "proxygroup.get", params)
}

// Create creates proxy groups and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxygroup/create
func (s *ProxyGroupService) Create(ctx context.Context, groups ...ProxyGroupCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "proxygroup.create", "proxy_groupids", groups)
}

// Update updates proxy groups and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxygroup/update
func (s *ProxyGroupService) Update(ctx context.Context, groups ...ProxyGroupUpdateParams) ([]string, error) {
	return updateIDs(ctx, s.client, "proxygroup.update", "proxy_groupids", groups)
}

// Delete deletes proxy groups by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxygroup/delete
func (s *ProxyGroupService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "proxygroup.delete", "proxy_groupids", ids)
}
