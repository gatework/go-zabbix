package zabbix

import "context"

// UserMacroService wraps the "usermacro" API namespace. It manages both host/template
// level user macros (UserMacro) and global macros (GlobalMacro).
type UserMacroService struct{ client *Client }

// UserMacro returns the usermacro service.
func (c *Client) UserMacro() *UserMacroService { return &UserMacroService{c} }

// GlobalMacro is a global (server-wide) user macro.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/object
type GlobalMacro struct {
	GlobalMacroID string `json:"globalmacroid,omitempty"`
	Macro         string `json:"macro,omitempty"`
	Value         string `json:"value,omitempty"`
	Type          string `json:"type,omitempty"`
	Description   string `json:"description,omitempty"`
}

// UserMacroGetParams are the parameters for usermacro.get.
// UserMacroGetParams queries the list form of usermacro.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type UserMacroGetParams struct {
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
	HostIDs                []string       `json:"hostids,omitzero"`
	HostMacroIDs           []string       `json:"hostmacroids,omitzero"`
}

// Get retrieves host/template user macros matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/get
func (s *UserMacroService) Get(ctx context.Context, params UserMacroGetParams) ([]UserMacro, error) {
	return getList[UserMacro](ctx, s.client, "usermacro.get", params)
}

// UserMacroCountParams contains the filters for usermacro.get counts.
type UserMacroCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	HostIDs                []string       `json:"hostids,omitzero"`
	HostMacroIDs           []string       `json:"hostmacroids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *UserMacroService) Count(ctx context.Context, params UserMacroCountParams) (uint64, error) {
	return getCount(ctx, s.client, "usermacro.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *UserMacroService) GetMap(ctx context.Context, params UserMacroGetParams) (map[string]UserMacro, error) {
	return getMap[UserMacro](ctx, s.client, "usermacro.get", params)
}

// Create creates host/template user macros and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/create
func (s *UserMacroService) Create(ctx context.Context, items ...HostMacroCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "usermacro.create", "hostmacroids", items)
}

// Update updates host/template user macros and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/update
func (s *UserMacroService) Update(ctx context.Context, items ...HostMacroUpdateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "usermacro.update", "hostmacroids", items)
}

// Delete deletes host/template user macros by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/delete
func (s *UserMacroService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "usermacro.delete", "hostmacroids", ids)
}

// GlobalMacroGetParams queries server-wide macros without host-macro filters.
type GlobalMacroGetParams struct {
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

	GlobalMacroIDs []string `json:"globalmacroids,omitzero"`
}

// GlobalMacroCountParams contains the filters for a global macro count.
type GlobalMacroCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`

	GlobalMacroIDs []string `json:"globalmacroids,omitzero"`
}

// GetGlobal retrieves global macros. The entity switch is fixed by this method.
func (s *UserMacroService) GetGlobal(ctx context.Context, params GlobalMacroGetParams) ([]GlobalMacro, error) {
	return getList[GlobalMacro](ctx, s.client, "usermacro.get", struct {
		GlobalMacroGetParams
		GlobalMacro bool `json:"globalmacro"`
	}{params, true})
}

// CountGlobal counts global macros matching params.
func (s *UserMacroService) CountGlobal(ctx context.Context, params GlobalMacroCountParams) (uint64, error) {
	return getCount(ctx, s.client, "usermacro.get", struct {
		GlobalMacroCountParams
		GlobalMacro bool `json:"globalmacro"`
	}{params, true})
}

// GetGlobalMap retrieves global macros keyed by their global macro IDs.
func (s *UserMacroService) GetGlobalMap(ctx context.Context, params GlobalMacroGetParams) (map[string]GlobalMacro, error) {
	return getMap[GlobalMacro](ctx, s.client, "usermacro.get", struct {
		GlobalMacroGetParams
		GlobalMacro bool `json:"globalmacro"`
	}{params, true})
}

// CreateGlobal creates global macros and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/createglobal
func (s *UserMacroService) CreateGlobal(ctx context.Context, items ...GlobalMacroCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "usermacro.createglobal", "globalmacroids", items)
}

// UpdateGlobal updates global macros and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/updateglobal
func (s *UserMacroService) UpdateGlobal(ctx context.Context, items ...GlobalMacroUpdateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "usermacro.updateglobal", "globalmacroids", items)
}

// DeleteGlobal deletes global macros by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/deleteglobal
func (s *UserMacroService) DeleteGlobal(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "usermacro.deleteglobal", "globalmacroids", ids)
}
