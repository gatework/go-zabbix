package zabbix

import "context"

// ValueMapService wraps the "valuemap" API namespace.
type ValueMapService struct{ client *Client }

// ValueMap returns the valuemap service.
func (c *Client) ValueMap() *ValueMapService { return &ValueMapService{c} }

// ValueMap object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/valuemap/object
type ValueMap struct {
	ValueMapID string `json:"valuemapid,omitempty"`
	HostID     string `json:"hostid,omitempty"`
	Name       string `json:"name,omitempty"`
	UUID       string `json:"uuid,omitempty"`

	Mappings []ValueMapMapping `json:"mappings,omitempty"`
}

// ValueMapMapping is a single value mapping within a value map.
type ValueMapMapping struct {
	Type      string `json:"type,omitempty"`
	Value     string `json:"value,omitempty"`
	NewValue  string `json:"newvalue,omitempty"`
	SortOrder string `json:"sortorder,omitempty"`
}

// ValueMapGetParams are the parameters for valuemap.get.
type ValueMapGetParams struct {
	GetParams
	ValueMapIDs    []string `json:"valuemapids,omitempty"`
	HostIDs        []string `json:"hostids,omitempty"`
	SelectMappings any      `json:"selectMappings,omitempty"`
}

// Get retrieves value maps matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/valuemap/get
func (s *ValueMapService) Get(ctx context.Context, params ValueMapGetParams) ([]ValueMap, error) {
	var out []ValueMap
	err := s.client.Call(ctx, "valuemap.get", params, &out)
	return out, err
}

// Create creates value maps and returns the new value map ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/valuemap/create
func (s *ValueMapService) Create(ctx context.Context, valuemaps ...ValueMap) ([]string, error) {
	var res struct {
		IDs []string `json:"valuemapids"`
	}
	err := s.client.Call(ctx, "valuemap.create", valuemaps, &res)
	return res.IDs, err
}

// Update updates value maps and returns the affected value map ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/valuemap/update
func (s *ValueMapService) Update(ctx context.Context, valuemaps ...ValueMap) ([]string, error) {
	var res struct {
		IDs []string `json:"valuemapids"`
	}
	err := s.client.Call(ctx, "valuemap.update", valuemaps, &res)
	return res.IDs, err
}

// Delete deletes value maps by id and returns the deleted value map ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/valuemap/delete
func (s *ValueMapService) Delete(ctx context.Context, valuemapIDs ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"valuemapids"`
	}
	err := s.client.Call(ctx, "valuemap.delete", valuemapIDs, &res)
	return res.IDs, err
}
