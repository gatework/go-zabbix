package zabbix

import "context"

// IconMapService wraps the "iconmap" API namespace.
type IconMapService struct{ client *Client }

// IconMap returns the iconmap service.
func (c *Client) IconMap() *IconMapService { return &IconMapService{c} }

// IconMap object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/iconmap/object
type IconMap struct {
	IconMapID     string `json:"iconmapid,omitempty"`
	Name          string `json:"name,omitempty"`
	DefaultIconID string `json:"default_iconid,omitempty"`

	Mappings []IconMapMapping `json:"mappings,omitempty"`
}

// IconMapMapping is a single icon mapping within an icon map.
type IconMapMapping struct {
	IconMappingID string `json:"iconmappingid,omitempty"`
	IconMapID     string `json:"iconmapid,omitempty"`
	IconID        string `json:"iconid,omitempty"`
	InventoryLink string `json:"inventory_link,omitempty"`
	Expression    string `json:"expression,omitempty"`
	SortOrder     string `json:"sortorder,omitempty"`
}

// IconMapGetParams are the parameters for iconmap.get.
// IconMapGetParams queries the list form of iconmap.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type IconMapGetParams struct {
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
	IconMapIDs             []string       `json:"iconmapids,omitzero"`
	SysmapIDs              []string       `json:"sysmapids,omitzero"`
	SelectMappings         any            `json:"selectMappings,omitempty"`
}

// Get retrieves icon maps matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/iconmap/get
func (s *IconMapService) Get(ctx context.Context, params IconMapGetParams) ([]IconMap, error) {
	return getList[IconMap](ctx, s.client, "iconmap.get", params)
}

// IconMapCountParams contains the filters for iconmap.get counts.
type IconMapCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	IconMapIDs             []string       `json:"iconmapids,omitzero"`
	SysmapIDs              []string       `json:"sysmapids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *IconMapService) Count(ctx context.Context, params IconMapCountParams) (uint64, error) {
	return getCount(ctx, s.client, "iconmap.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *IconMapService) GetMap(ctx context.Context, params IconMapGetParams) (map[string]IconMap, error) {
	return getMap[IconMap](ctx, s.client, "iconmap.get", params)
}

// Create creates icon maps and returns the new icon map ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/iconmap/create
func (s *IconMapService) Create(ctx context.Context, iconmaps ...IconMapCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "iconmap.create", "iconmapids", iconmaps)
}

// Update updates icon maps and returns the affected icon map ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/iconmap/update
func (s *IconMapService) Update(ctx context.Context, iconmaps ...IconMapUpdateParams) ([]string, error) {
	return updateIDs(ctx, s.client, "iconmap.update", "iconmapids", iconmaps)
}

// Delete deletes icon maps by id and returns the deleted icon map ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/iconmap/delete
func (s *IconMapService) Delete(ctx context.Context, iconmapIDs ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "iconmap.delete", "iconmapids", iconmapIDs)
}
