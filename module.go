package zabbix

import "context"

// ModuleService wraps the "module" API namespace.
type ModuleService struct{ client *Client }

// Module returns the module service.
func (c *Client) Module() *ModuleService { return &ModuleService{c} }

// Module object (Zabbix 7.x frontend module).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/module/object
type Module struct {
	ModuleID     string `json:"moduleid,omitempty"`
	ID           string `json:"id,omitempty"`
	RelativePath string `json:"relative_path,omitempty"`
	Status       string `json:"status,omitempty"`
	Config       any    `json:"config,omitempty"`
}

// ModuleGetParams are the parameters for module.get.
type ModuleGetParams struct {
	GetParams
	ModuleIDs []string `json:"moduleids,omitempty"`
}

// Get retrieves modules matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/module/get
func (s *ModuleService) Get(ctx context.Context, params ModuleGetParams) ([]Module, error) {
	var out []Module
	err := s.client.Call(ctx, "module.get", params, &out)
	return out, err
}

// Create creates modules and returns the new module ids.
func (s *ModuleService) Create(ctx context.Context, items ...Module) ([]string, error) {
	var res struct {
		IDs []string `json:"moduleids"`
	}
	err := s.client.Call(ctx, "module.create", items, &res)
	return res.IDs, err
}

// Update updates modules and returns the affected module ids.
func (s *ModuleService) Update(ctx context.Context, items ...Module) ([]string, error) {
	var res struct {
		IDs []string `json:"moduleids"`
	}
	err := s.client.Call(ctx, "module.update", items, &res)
	return res.IDs, err
}

// Delete deletes modules by id and returns the deleted module ids.
func (s *ModuleService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"moduleids"`
	}
	err := s.client.Call(ctx, "module.delete", ids, &res)
	return res.IDs, err
}
