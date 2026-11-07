package zabbix

import "context"

// ModuleService wraps the "module" API namespace.
type ModuleService struct{ client *Client }

// Module returns the module service.
func (c *Client) Module() *ModuleService { return &ModuleService{c} }

// Module object (Zabbix 7.x frontend module).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/module/object
type Module struct {
	ModuleID     string `json:"moduleid,omitempty"`
	ID           string `json:"id,omitempty"`
	RelativePath string `json:"relative_path,omitempty"`
	Status       string `json:"status,omitempty"`
	Config       any    `json:"config,omitempty"`
}

// ModuleGetParams are the parameters for module.get.
// ModuleGetParams queries the list form of module.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type ModuleGetParams struct {
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
	ModuleIDs              []string       `json:"moduleids,omitzero"`
}

// Get retrieves modules matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/module/get
func (s *ModuleService) Get(ctx context.Context, params ModuleGetParams) ([]Module, error) {
	return getList[Module](ctx, s.client, "module.get", params)
}

// ModuleCountParams contains the filters for module.get counts.
type ModuleCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	ModuleIDs              []string       `json:"moduleids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *ModuleService) Count(ctx context.Context, params ModuleCountParams) (uint64, error) {
	return getCount(ctx, s.client, "module.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *ModuleService) GetMap(ctx context.Context, params ModuleGetParams) (map[string]Module, error) {
	return getMap[Module](ctx, s.client, "module.get", params)
}

// Create creates modules and returns the new module ids.
func (s *ModuleService) Create(ctx context.Context, items ...ModuleCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "module.create", "moduleids", items)
}

// Update updates modules and returns the affected module ids.
func (s *ModuleService) Update(ctx context.Context, items ...ModuleUpdateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "module.update", "moduleids", items)
}

// Delete deletes modules by id and returns the deleted module ids.
func (s *ModuleService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "module.delete", "moduleids", ids)
}
