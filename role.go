package zabbix

import (
	"context"
	"encoding/json"
)

// RoleService wraps the "role" API namespace.
type RoleService struct{ client *Client }

// Role returns the role service.
func (c *Client) Role() *RoleService { return &RoleService{c} }

// Role object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/role/object
type Role struct {
	RoleID   string `json:"roleid,omitempty"`
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	ReadOnly string `json:"readonly,omitempty"`

	Rules SelectedObject[RoleRules] `json:"rules,omitzero"`
}

// RoleRules contains the permissions returned by selectRules.
type RoleRules struct {
	UI                   []RoleRuleStatus   `json:"ui,omitempty"`
	UIDefaultAccess      string             `json:"ui.default_access,omitempty"`
	ServicesReadMode     string             `json:"services.read.mode,omitempty"`
	ServicesReadList     []Service          `json:"services.read.list,omitempty"`
	ServicesReadTag      *ServiceAccessTag  `json:"services.read.tag,omitempty"`
	ServicesWriteMode    string             `json:"services.write.mode,omitempty"`
	ServicesWriteList    []Service          `json:"services.write.list,omitempty"`
	ServicesWriteTag     *ServiceAccessTag  `json:"services.write.tag,omitempty"`
	Modules              []RoleModuleStatus `json:"modules,omitempty"`
	ModulesDefaultAccess string             `json:"modules.default_access,omitempty"`
	APIAccess            string             `json:"api.access,omitempty"`
	APIMode              string             `json:"api.mode,omitempty"`
	API                  []string           `json:"api,omitempty"`
	Actions              []RoleRuleStatus   `json:"actions,omitempty"`
	ActionsDefaultAccess string             `json:"actions.default_access,omitempty"`
}

// ServiceAccessTag is a tag condition granting a role access to services.
// Empty Tag and Value fields mean there is no tag-based access condition.
type ServiceAccessTag struct {
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

// RoleRuleStatus controls one UI element or action.
type RoleRuleStatus struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}

// RoleModuleStatus controls access to one module.
type RoleModuleStatus struct {
	// PHP array keys can produce numeric module IDs even though IDs in stored
	// rules are strings. Number preserves both forms without float conversion.
	ModuleID json.Number `json:"moduleid,omitempty"`
	Status   string      `json:"status,omitempty"`
}

// RoleGetParams are the parameters for role.get.
// RoleGetParams queries the list form of role.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type RoleGetParams struct {
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
	RoleIDs                []string       `json:"roleids,omitzero"`
	SelectRules            any            `json:"selectRules,omitempty"`
}

// Get retrieves roles matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/role/get
func (s *RoleService) Get(ctx context.Context, params RoleGetParams) ([]Role, error) {
	return getList[Role](ctx, s.client, "role.get", params)
}

// RoleCountParams contains the filters for role.get counts.
type RoleCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	RoleIDs                []string       `json:"roleids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *RoleService) Count(ctx context.Context, params RoleCountParams) (uint64, error) {
	return getCount(ctx, s.client, "role.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *RoleService) GetMap(ctx context.Context, params RoleGetParams) (map[string]Role, error) {
	return getMap[Role](ctx, s.client, "role.get", params)
}

// Create creates roles and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/role/create
func (s *RoleService) Create(ctx context.Context, items ...RoleCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "role.create", "roleids", items)
}

// Update updates roles and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/role/update
func (s *RoleService) Update(ctx context.Context, items ...RoleUpdateParams) ([]string, error) {
	return updateIDs(ctx, s.client, "role.update", "roleids", items)
}

// Delete deletes roles by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/role/delete
func (s *RoleService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "role.delete", "roleids", ids)
}
