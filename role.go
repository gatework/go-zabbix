package zabbix

import "context"

// RoleService wraps the "role" API namespace.
type RoleService struct{ client *Client }

// Role returns the role service.
func (c *Client) Role() *RoleService { return &RoleService{c} }

// Role object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/role/object
type Role struct {
	RoleID   string `json:"roleid,omitempty"`
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Readonly string `json:"readonly,omitempty"`

	// Rules is a complex nested object (ui, actions, services, modules, api, ...);
	// it is represented untyped.
	Rules any `json:"rules,omitempty"`
}

// RoleGetParams are the parameters for role.get.
type RoleGetParams struct {
	GetParams
	RoleIDs     []string `json:"roleids,omitempty"`
	SelectRules any      `json:"selectRules,omitempty"`
}

// Get retrieves roles matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/role/get
func (s *RoleService) Get(ctx context.Context, params RoleGetParams) ([]Role, error) {
	var out []Role
	err := s.client.Call(ctx, "role.get", params, &out)
	return out, err
}

// Create creates roles and returns the new ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/role/create
func (s *RoleService) Create(ctx context.Context, items ...Role) ([]string, error) {
	var res struct {
		IDs []string `json:"roleids"`
	}
	err := s.client.Call(ctx, "role.create", items, &res)
	return res.IDs, err
}

// Update updates roles and returns the affected ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/role/update
func (s *RoleService) Update(ctx context.Context, items ...Role) ([]string, error) {
	var res struct {
		IDs []string `json:"roleids"`
	}
	err := s.client.Call(ctx, "role.update", items, &res)
	return res.IDs, err
}

// Delete deletes roles by id and returns the deleted ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/role/delete
func (s *RoleService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"roleids"`
	}
	err := s.client.Call(ctx, "role.delete", ids, &res)
	return res.IDs, err
}
