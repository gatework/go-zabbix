package zabbix

import "context"

// UserGroupService wraps the "usergroup" API namespace.
type UserGroupService struct{ client *Client }

// UserGroup returns the usergroup service.
func (c *Client) UserGroup() *UserGroupService { return &UserGroupService{c} }

// UserGroup object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usergroup/object
type UserGroup struct {
	UsrGrpID        string `json:"usrgrpid,omitempty"`
	Name            string `json:"name,omitempty"`
	GUIAccess       string `json:"gui_access,omitempty"`
	UsersStatus     string `json:"users_status,omitempty"`
	DebugMode       string `json:"debug_mode,omitempty"`
	UserDirectoryID string `json:"userdirectoryid,omitempty"`
	MFAStatus       string `json:"mfa_status,omitempty"`
	MFAID           string `json:"mfaid,omitempty"`

	TagFilters          []UserGroupTagFilter `json:"tag_filters,omitempty"`
	Users               []User               `json:"users,omitempty"`
	HostGroupRights     []UserGroupRight     `json:"hostgroup_rights,omitempty"`
	TemplateGroupRights []UserGroupRight     `json:"template_group_rights,omitempty"`
}

// UserGroupRight is a host-group permission entry for a user group.
type UserGroupRight struct {
	Permission string `json:"permission,omitempty"`
	ID         string `json:"id,omitempty"`
}

// UserGroupTagFilter restricts a user group's visibility by tag on a host group.
type UserGroupTagFilter struct {
	GroupID string `json:"groupid,omitempty"`
	Tag     string `json:"tag,omitempty"`
	Value   string `json:"value,omitempty"`
}

// UserGroupGetParams are the parameters for usergroup.get.
// UserGroupGetParams queries the list form of usergroup.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type UserGroupGetParams struct {
	Output                    any            `json:"output,omitempty"`
	Filter                    map[string]any `json:"filter,omitzero"`
	Search                    map[string]any `json:"search,omitzero"`
	SearchByAny               bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled    bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch               bool           `json:"startSearch,omitempty"`
	ExcludeSearch             bool           `json:"excludeSearch,omitempty"`
	SortField                 []string       `json:"sortfield,omitzero"`
	SortOrder                 []string       `json:"sortorder,omitzero"`
	Limit                     int            `json:"limit,omitempty"`
	Editable                  bool           `json:"editable,omitempty"`
	UsrGrpIDs                 []string       `json:"usrgrpids,omitzero"`
	UserIDs                   []string       `json:"userids,omitzero"`
	Status                    *int           `json:"status,omitempty"`
	SelectUsers               any            `json:"selectUsers,omitempty"`
	SelectTagFilters          any            `json:"selectTagFilters,omitempty"`
	SelectHostGroupRights     any            `json:"selectHostGroupRights,omitempty"`
	SelectTemplateGroupRights any            `json:"selectTemplateGroupRights,omitempty"`
}

// Get retrieves user groups matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usergroup/get
func (s *UserGroupService) Get(ctx context.Context, params UserGroupGetParams) ([]UserGroup, error) {
	return getList[UserGroup](ctx, s.client, "usergroup.get", params)
}

// UserGroupCountParams contains the filters for usergroup.get counts.
type UserGroupCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	UsrGrpIDs              []string       `json:"usrgrpids,omitzero"`
	UserIDs                []string       `json:"userids,omitzero"`
	Status                 *int           `json:"status,omitempty"`
}

// Count returns the number of objects matching params.
func (s *UserGroupService) Count(ctx context.Context, params UserGroupCountParams) (uint64, error) {
	return getCount(ctx, s.client, "usergroup.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *UserGroupService) GetMap(ctx context.Context, params UserGroupGetParams) (map[string]UserGroup, error) {
	return getMap[UserGroup](ctx, s.client, "usergroup.get", params)
}

// Create creates user groups and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usergroup/create
func (s *UserGroupService) Create(ctx context.Context, items ...UserGroupCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "usergroup.create", "usrgrpids", items)
}

// Update updates user groups and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usergroup/update
func (s *UserGroupService) Update(ctx context.Context, items ...UserGroupUpdateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "usergroup.update", "usrgrpids", items)
}

// Delete deletes user groups by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usergroup/delete
func (s *UserGroupService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "usergroup.delete", "usrgrpids", ids)
}
