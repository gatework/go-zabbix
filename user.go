package zabbix

import "context"

// UserService wraps the "user" API namespace.
type UserService struct{ client *Client }

// User returns the user service.
func (c *Client) User() *UserService { return &UserService{c} }

// User object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/user/object
type User struct {
	UserID          string   `json:"userid,omitempty"`
	Username        string   `json:"username,omitempty"`
	Passwd          string   `json:"passwd,omitempty"`
	RoleID          string   `json:"roleid,omitempty"`
	AttemptClock    string   `json:"attempt_clock,omitempty"`
	AttemptFailed   string   `json:"attempt_failed,omitempty"`
	AttemptIP       string   `json:"attempt_ip,omitempty"`
	Autologin       string   `json:"autologin,omitempty"`
	Autologout      string   `json:"autologout,omitempty"`
	Lang            string   `json:"lang,omitempty"`
	Name            string   `json:"name,omitempty"`
	Refresh         string   `json:"refresh,omitempty"`
	RowsPerPage     string   `json:"rows_per_page,omitempty"`
	Surname         string   `json:"surname,omitempty"`
	Theme           string   `json:"theme,omitempty"`
	TsProvisioned   string   `json:"ts_provisioned,omitempty"`
	URL             string   `json:"url,omitempty"`
	UserDirectoryID string   `json:"userdirectoryid,omitempty"`
	Timezone        string   `json:"timezone,omitempty"`
	UserGroups      []UsrGrp `json:"usrgrps,omitempty"`
	Medias          []Media  `json:"medias,omitempty"`
}

// UsrGrp is a user-group reference used when creating or updating users.
type UsrGrp struct {
	UsrGrpID string `json:"usrgrpid"`
}

// Media object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/user/object
type Media struct {
	MediaTypeID          string `json:"mediatypeid,omitempty"`
	SendTo               any    `json:"sendto,omitempty"` // string or []string depending on media type
	Active               string `json:"active,omitempty"`
	Severity             string `json:"severity,omitempty"`
	Period               string `json:"period,omitempty"`
	UserDirectoryMediaID string `json:"userdirectory_mediaid,omitempty"`
}

// UserGetParams are the parameters for user.get.
type UserGetParams struct {
	GetParams
	UserIDs       []string `json:"userids,omitempty"`
	MediaIDs      []string `json:"mediaids,omitempty"`
	MediaTypeIDs  []string `json:"mediatypeids,omitempty"`
	UsrGrpIDs     []string `json:"usrgrpids,omitempty"`
	SelectMedias  any      `json:"selectMedias,omitempty"`
	SelectUsrGrps any      `json:"selectUsrgrps,omitempty"`
	SelectRole    any      `json:"selectRole,omitempty"`
	GetAccess     bool     `json:"getAccess,omitempty"`
}

// Get retrieves users matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/user/get
func (s *UserService) Get(ctx context.Context, params UserGetParams) ([]User, error) {
	var users []User
	err := s.client.Call(ctx, "user.get", params, &users)
	return users, err
}

// Create creates one or more users and returns the new user ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/user/create
func (s *UserService) Create(ctx context.Context, users ...User) ([]string, error) {
	var res struct {
		UserIDs []string `json:"userids"`
	}
	err := s.client.Call(ctx, "user.create", users, &res)
	return res.UserIDs, err
}

// Update updates one or more users and returns the affected user ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/user/update
func (s *UserService) Update(ctx context.Context, users ...User) ([]string, error) {
	var res struct {
		UserIDs []string `json:"userids"`
	}
	err := s.client.Call(ctx, "user.update", users, &res)
	return res.UserIDs, err
}

// Delete deletes users by id and returns the deleted user ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/user/delete
func (s *UserService) Delete(ctx context.Context, userIDs ...string) ([]string, error) {
	var res struct {
		UserIDs []string `json:"userids"`
	}
	err := s.client.Call(ctx, "user.delete", userIDs, &res)
	return res.UserIDs, err
}
