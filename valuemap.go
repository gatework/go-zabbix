package zabbix

import "context"

// ValueMapService wraps the "valuemap" API namespace.
type ValueMapService struct{ client *Client }

// ValueMap returns the valuemap service.
func (c *Client) ValueMap() *ValueMapService { return &ValueMapService{c} }

// ValueMap object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/valuemap/object
type ValueMap struct {
	ValueMapID string                    `json:"valuemapid,omitempty"`
	HostID     string                    `json:"hostid,omitempty"`
	Name       string                    `json:"name,omitempty"`
	UUID       string                    `json:"uuid,omitempty"`
	Mappings   Relation[ValueMapMapping] `json:"mappings,omitzero"`
}

// ValueMapMapping is a single value mapping within a value map.
type ValueMapMapping struct {
	Type      string `json:"type,omitempty"`
	Value     string `json:"value,omitempty"`
	NewValue  string `json:"newvalue,omitempty"`
	SortOrder string `json:"sortorder,omitempty"`
}

// ValueMapGetParams are the parameters for valuemap.get.
// ValueMapGetParams queries the list form of valuemap.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type ValueMapGetParams struct {
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
	ValueMapIDs            []string       `json:"valuemapids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	SelectMappings         any            `json:"selectMappings,omitempty"`
}

// Get retrieves value maps matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/valuemap/get
func (s *ValueMapService) Get(ctx context.Context, params ValueMapGetParams) ([]ValueMap, error) {
	return getList[ValueMap](ctx, s.client, "valuemap.get", params)
}

// ValueMapCountParams contains the filters for valuemap.get counts.
type ValueMapCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	ValueMapIDs            []string       `json:"valuemapids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *ValueMapService) Count(ctx context.Context, params ValueMapCountParams) (uint64, error) {
	return getCount(ctx, s.client, "valuemap.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *ValueMapService) GetMap(ctx context.Context, params ValueMapGetParams) (map[string]ValueMap, error) {
	return getMap[ValueMap](ctx, s.client, "valuemap.get", params)
}

// Create creates value maps and returns the new value map ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/valuemap/create
func (s *ValueMapService) Create(ctx context.Context, valuemaps ...ValueMapCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "valuemap.create", "valuemapids", valuemaps)
}

// Update updates value maps and returns the affected value map ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/valuemap/update
func (s *ValueMapService) Update(ctx context.Context, valuemaps ...ValueMapUpdateParams) ([]string, error) {
	return updateIDs(ctx, s.client, "valuemap.update", "valuemapids", valuemaps)
}

// Delete deletes value maps by id and returns the deleted value map ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/valuemap/delete
func (s *ValueMapService) Delete(ctx context.Context, valuemapIDs ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "valuemap.delete", "valuemapids", valuemapIDs)
}
