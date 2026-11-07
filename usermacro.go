package zabbix

import "context"

// UserMacroService wraps the "usermacro" API namespace. It manages both host/template
// level user macros (UserMacro) and global macros (GlobalMacro).
type UserMacroService struct{ client *Client }

// UserMacro returns the usermacro service.
func (c *Client) UserMacro() *UserMacroService { return &UserMacroService{c} }

// GlobalMacro is a global (server-wide) user macro.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usermacro/object
type GlobalMacro struct {
	GlobalMacroID string `json:"globalmacroid,omitempty"`
	Macro         string `json:"macro,omitempty"`
	Value         string `json:"value,omitempty"`
	Type          string `json:"type,omitempty"`
	Description   string `json:"description,omitempty"`
}

// UserMacroGetParams are the parameters for usermacro.get.
type UserMacroGetParams struct {
	GetParams
	HostIDs        []string `json:"hostids,omitempty"`
	HostMacroIDs   []string `json:"hostmacroids,omitempty"`
	GlobalMacroIDs []string `json:"globalmacroids,omitempty"`
	Globalmacro    bool     `json:"globalmacro,omitempty"`
}

// Get retrieves host/template user macros matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usermacro/get
func (s *UserMacroService) Get(ctx context.Context, params UserMacroGetParams) ([]UserMacro, error) {
	var out []UserMacro
	err := s.client.Call(ctx, "usermacro.get", params, &out)
	return out, err
}

// Create creates host/template user macros and returns the new ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usermacro/create
func (s *UserMacroService) Create(ctx context.Context, items ...UserMacro) ([]string, error) {
	var res struct {
		IDs []string `json:"hostmacroids"`
	}
	err := s.client.Call(ctx, "usermacro.create", items, &res)
	return res.IDs, err
}

// Update updates host/template user macros and returns the affected ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usermacro/update
func (s *UserMacroService) Update(ctx context.Context, items ...UserMacro) ([]string, error) {
	var res struct {
		IDs []string `json:"hostmacroids"`
	}
	err := s.client.Call(ctx, "usermacro.update", items, &res)
	return res.IDs, err
}

// Delete deletes host/template user macros by id and returns the deleted ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usermacro/delete
func (s *UserMacroService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"hostmacroids"`
	}
	err := s.client.Call(ctx, "usermacro.delete", ids, &res)
	return res.IDs, err
}

// GetGlobal retrieves global macros matching params. It calls usermacro.get with
// globalmacro=true.
func (s *UserMacroService) GetGlobal(ctx context.Context, params UserMacroGetParams) ([]GlobalMacro, error) {
	params.Globalmacro = true
	var out []GlobalMacro
	err := s.client.Call(ctx, "usermacro.get", params, &out)
	return out, err
}

// CreateGlobal creates global macros and returns the new ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usermacro/createglobal
func (s *UserMacroService) CreateGlobal(ctx context.Context, items ...GlobalMacro) ([]string, error) {
	var res struct {
		IDs []string `json:"globalmacroids"`
	}
	err := s.client.Call(ctx, "usermacro.createglobal", items, &res)
	return res.IDs, err
}

// UpdateGlobal updates global macros and returns the affected ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usermacro/updateglobal
func (s *UserMacroService) UpdateGlobal(ctx context.Context, items ...GlobalMacro) ([]string, error) {
	var res struct {
		IDs []string `json:"globalmacroids"`
	}
	err := s.client.Call(ctx, "usermacro.updateglobal", items, &res)
	return res.IDs, err
}

// DeleteGlobal deletes global macros by id and returns the deleted ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usermacro/deleteglobal
func (s *UserMacroService) DeleteGlobal(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"globalmacroids"`
	}
	err := s.client.Call(ctx, "usermacro.deleteglobal", ids, &res)
	return res.IDs, err
}
