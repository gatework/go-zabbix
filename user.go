package zabbix

import (
	"context"
	"encoding/json"
)

// UserService wraps the "user" API namespace.
type UserService struct{ client *Client }

// User returns the user service.
func (c *Client) User() *UserService { return &UserService{c} }

// User object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/user/object
type User struct {
	UserID          string               `json:"userid,omitempty"`
	Username        string               `json:"username,omitempty"`
	Passwd          string               `json:"passwd,omitempty"`
	RoleID          string               `json:"roleid,omitempty"`
	AttemptClock    string               `json:"attempt_clock,omitempty"`
	AttemptFailed   string               `json:"attempt_failed,omitempty"`
	AttemptIP       string               `json:"attempt_ip,omitempty"`
	Autologin       string               `json:"autologin,omitempty"`
	Autologout      string               `json:"autologout,omitempty"`
	Lang            string               `json:"lang,omitempty"`
	Name            string               `json:"name,omitempty"`
	Refresh         string               `json:"refresh,omitempty"`
	RowsPerPage     string               `json:"rows_per_page,omitempty"`
	Surname         string               `json:"surname,omitempty"`
	Theme           string               `json:"theme,omitempty"`
	TsProvisioned   string               `json:"ts_provisioned,omitempty"`
	URL             string               `json:"url,omitempty"`
	UserDirectoryID string               `json:"userdirectoryid,omitempty"`
	Timezone        string               `json:"timezone,omitempty"`
	UserGroups      []UserGroup          `json:"usrgrps,omitempty"`
	Medias          []Media              `json:"medias,omitempty"`
	Role            SelectedObject[Role] `json:"role,omitzero"`
	GUIAccess       json.Number          `json:"gui_access,omitempty"`
	DebugMode       json.Number          `json:"debug_mode,omitempty"`
	UsersStatus     json.Number          `json:"users_status,omitempty"`
}

// UsrGrp is a user-group reference used when creating or updating users.
type UsrGrp struct {
	UsrGrpID string `json:"usrgrpid"`
}

// Media object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/user/object
type Media struct {
	MediaTypeID          string `json:"mediatypeid,omitempty"`
	SendTo               any    `json:"sendto,omitempty"` // string or []string depending on media type
	Active               string `json:"active,omitempty"`
	Severity             string `json:"severity,omitempty"`
	Period               string `json:"period,omitempty"`
	UserDirectoryMediaID string `json:"userdirectory_mediaid,omitempty"`
	MediaID              string `json:"mediaid,omitempty"`
}

// UserGetParams are the parameters for user.get.
// UserGetParams queries the list form of user.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type UserGetParams struct {
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
	UserIDs                []string       `json:"userids,omitzero"`
	MediaIDs               []string       `json:"mediaids,omitzero"`
	MediaTypeIDs           []string       `json:"mediatypeids,omitzero"`
	UsrGrpIDs              []string       `json:"usrgrpids,omitzero"`
	SelectMedias           any            `json:"selectMedias,omitempty"`
	SelectUserGroups       any            `json:"selectUsrgrps,omitempty"`
	SelectRole             any            `json:"selectRole,omitempty"`
	GetAccess              bool           `json:"getAccess,omitempty"`
}

// Get retrieves users matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/user/get
func (s *UserService) Get(ctx context.Context, params UserGetParams) ([]User, error) {
	return getList[User](ctx, s.client, "user.get", params)
}

// UserCountParams contains the filters for user.get counts.
type UserCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	UserIDs                []string       `json:"userids,omitzero"`
	MediaIDs               []string       `json:"mediaids,omitzero"`
	MediaTypeIDs           []string       `json:"mediatypeids,omitzero"`
	UsrGrpIDs              []string       `json:"usrgrpids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *UserService) Count(ctx context.Context, params UserCountParams) (uint64, error) {
	return getCount(ctx, s.client, "user.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *UserService) GetMap(ctx context.Context, params UserGetParams) (map[string]User, error) {
	return getMap[User](ctx, s.client, "user.get", params)
}

// Create creates one or more users and returns the new user ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/user/create
func (s *UserService) Create(ctx context.Context, users ...UserCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "user.create", "userids", users)
}

// Update updates one or more users and returns the affected user ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/user/update
func (s *UserService) Update(ctx context.Context, users ...UserUpdateParams) ([]string, error) {
	return updateIDs(ctx, s.client, "user.update", "userids", users)
}

// Delete deletes users by id and returns the deleted user ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/user/delete
func (s *UserService) Delete(ctx context.Context, userIDs ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "user.delete", "userids", userIDs)
}
