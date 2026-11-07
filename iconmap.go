package zabbix

import "context"

// IconMapService wraps the "iconmap" API namespace.
type IconMapService struct{ client *Client }

// IconMap returns the iconmap service.
func (c *Client) IconMap() *IconMapService { return &IconMapService{c} }

// IconMap object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/iconmap/object
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
type IconMapGetParams struct {
	GetParams
	IconMapIDs     []string `json:"iconmapids,omitempty"`
	SysmapIDs      []string `json:"sysmapids,omitempty"`
	SelectMappings any      `json:"selectMappings,omitempty"`
}

// Get retrieves icon maps matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/iconmap/get
func (s *IconMapService) Get(ctx context.Context, params IconMapGetParams) ([]IconMap, error) {
	var out []IconMap
	err := s.client.Call(ctx, "iconmap.get", params, &out)
	return out, err
}

// Create creates icon maps and returns the new icon map ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/iconmap/create
func (s *IconMapService) Create(ctx context.Context, iconmaps ...IconMap) ([]string, error) {
	var res struct {
		IDs []string `json:"iconmapids"`
	}
	err := s.client.Call(ctx, "iconmap.create", iconmaps, &res)
	return res.IDs, err
}

// Update updates icon maps and returns the affected icon map ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/iconmap/update
func (s *IconMapService) Update(ctx context.Context, iconmaps ...IconMap) ([]string, error) {
	var res struct {
		IDs []string `json:"iconmapids"`
	}
	err := s.client.Call(ctx, "iconmap.update", iconmaps, &res)
	return res.IDs, err
}

// Delete deletes icon maps by id and returns the deleted icon map ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/iconmap/delete
func (s *IconMapService) Delete(ctx context.Context, iconmapIDs ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"iconmapids"`
	}
	err := s.client.Call(ctx, "iconmap.delete", iconmapIDs, &res)
	return res.IDs, err
}
