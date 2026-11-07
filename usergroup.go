package zabbix

import "context"

// UserGroupService wraps the "usergroup" API namespace.
type UserGroupService struct{ client *Client }

// UserGroup returns the usergroup service.
func (c *Client) UserGroup() *UserGroupService { return &UserGroupService{c} }

// UserGroup object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usergroup/object
type UserGroup struct {
	UsrGrpID        string `json:"usrgrpid,omitempty"`
	Name            string `json:"name,omitempty"`
	GUIAccess       string `json:"gui_access,omitempty"`
	UsersStatus     string `json:"users_status,omitempty"`
	DebugMode       string `json:"debug_mode,omitempty"`
	UserDirectoryID string `json:"userdirectoryid,omitempty"`
	MFAStatus       string `json:"mfa_status,omitempty"`
	MFAID           string `json:"mfaid,omitempty"`

	Rights     []UserGroupRight     `json:"rights,omitempty"`
	TagFilters []UserGroupTagFilter `json:"tag_filters,omitempty"`
	Users      []User               `json:"users,omitempty"`
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
type UserGroupGetParams struct {
	GetParams
	UsrGrpIDs        []string `json:"usrgrpids,omitempty"`
	UserIDs          []string `json:"userids,omitempty"`
	Status           *int     `json:"status,omitempty"`
	SelectUsers      any      `json:"selectUsers,omitempty"`
	SelectRights     any      `json:"selectRights,omitempty"`
	SelectTagFilters any      `json:"selectTagFilters,omitempty"`
}

// Get retrieves user groups matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usergroup/get
func (s *UserGroupService) Get(ctx context.Context, params UserGroupGetParams) ([]UserGroup, error) {
	var out []UserGroup
	err := s.client.Call(ctx, "usergroup.get", params, &out)
	return out, err
}

// Create creates user groups and returns the new ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usergroup/create
func (s *UserGroupService) Create(ctx context.Context, items ...UserGroup) ([]string, error) {
	var res struct {
		IDs []string `json:"usrgrpids"`
	}
	err := s.client.Call(ctx, "usergroup.create", items, &res)
	return res.IDs, err
}

// Update updates user groups and returns the affected ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usergroup/update
func (s *UserGroupService) Update(ctx context.Context, items ...UserGroup) ([]string, error) {
	var res struct {
		IDs []string `json:"usrgrpids"`
	}
	err := s.client.Call(ctx, "usergroup.update", items, &res)
	return res.IDs, err
}

// Delete deletes user groups by id and returns the deleted ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/usergroup/delete
func (s *UserGroupService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"usrgrpids"`
	}
	err := s.client.Call(ctx, "usergroup.delete", ids, &res)
	return res.IDs, err
}
