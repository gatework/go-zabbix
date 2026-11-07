package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
)

// Session describes the authenticated user returned by user.login with userData.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/user/login
type Session struct {
	User
	Type          int    `json:"type"`
	UserIP        string `json:"userip"`
	DebugMode     int    `json:"debug_mode"`
	GUIAccess     string `json:"gui_access"`
	MfaID         int    `json:"mfaid"`
	Deprovisioned bool   `json:"deprovisioned"`
	AuthType      int    `json:"auth_type"`
	SessionID     string `json:"sessionid"`
	Secret        string `json:"secret"`
}

type loginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserData bool   `json:"userData,omitempty"`
}

// Login authenticates with a username and password, stores the returned session
// token on the client, and returns it.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/user/login
func (s *UserService) Login(ctx context.Context, username, password string) (string, error) {
	var raw json.RawMessage
	if err := s.client.do(ctx, "user.login",
		loginParams{Username: username, Password: password}, &raw, true); err != nil {
		return "", err
	}
	var token string
	if err := json.Unmarshal(raw, &token); err != nil {
		return "", fmt.Errorf("gozabbix: decode login token: %w", err)
	}
	s.client.SetToken(token)
	return token, nil
}

// LoginWithUserData authenticates and returns the full session details. The
// session token is stored on the client.
func (s *UserService) LoginWithUserData(ctx context.Context, username, password string) (*Session, error) {
	var session Session
	if err := s.client.do(ctx, "user.login",
		loginParams{Username: username, Password: password, UserData: true}, &session, true); err != nil {
		return nil, err
	}
	s.client.SetToken(session.SessionID)
	return &session, nil
}

// Logout invalidates the current session token on the server and clears it from
// the client. It is a no-op for pre-created API tokens (which cannot be logged
// out); callers using WithToken should not call Logout.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/user/logout
func (s *UserService) Logout(ctx context.Context) error {
	var ok bool
	if err := s.client.do(ctx, "user.logout", []any{}, &ok, true); err != nil {
		return err
	}
	s.client.SetToken("")
	return nil
}

// Login is a convenience shortcut for User().Login that discards the token
// (which is stored on the client regardless).
func (c *Client) Login(ctx context.Context, username, password string) error {
	_, err := c.User().Login(ctx, username, password)
	return err
}
