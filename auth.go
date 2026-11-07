package zabbix

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

const authCleanupTimeout = 5 * time.Second

type credentialKind uint8

const (
	credentialNone credentialKind = iota
	credentialAPI
	credentialSession
)

type credentials struct {
	kind  credentialKind
	token string
}

// SetAPIToken installs an unverified API token and supersedes pending
// authentication operations. It does not revoke previously installed credentials.
func (c *Client) SetAPIToken(token string) error {
	return c.setCredentials(credentialAPI, token)
}

// SetSessionToken installs an unverified session token and supersedes pending
// authentication operations. It does not revoke previously installed credentials.
func (c *Client) SetSessionToken(token string) error {
	return c.setCredentials(credentialSession, token)
}

func (c *Client) setCredentials(kind credentialKind, token string) error {
	if err := validateToken(token); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.generation++
	c.credentials = credentials{kind: kind, token: token}
	return nil
}

// ClearCredentials forgets local credentials and supersedes pending
// authentication operations. Use Logout to invalidate a session on the server.
func (c *Client) ClearCredentials() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.generation++
	c.credentials = credentials{}
}

// HasCredentials reports local credential presence, not server-side validity.
func (c *Client) HasCredentials() bool {
	state, _ := c.authSnapshot()
	return state.kind != credentialNone
}

func (c *Client) authSnapshot() (credentials, uint64) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.credentials, c.generation
}

func validateToken(token string) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("zabbix: credential token must not be blank")
	}
	for i := 0; i < len(token); i++ {
		if token[i] <= 0x20 || token[i] == 0x7f {
			return errors.New("zabbix: credential token contains whitespace or control characters")
		}
	}
	return nil
}

// Session describes the authenticated user returned by user.login with userData.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/user/login
type Session struct {
	User
	Type          int    `json:"type"`
	UserIP        string `json:"userip"`
	MFAID         int    `json:"mfaid"`
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

// Login authenticates without sending existing credentials and installs the new
// session. Existing credentials remain usable until login succeeds. A newer
// authentication operation or credential setter supersedes this operation;
// an orphaned successful session is logged out with a bounded cleanup context.
// Call Logout first if a previously installed session must be revoked.
func (c *Client) Login(ctx context.Context, username, password string) error {
	_, err := c.login(ctx, username, password, false)
	return err
}

// LoginWithUserData logs in and returns the full session. It has the same
// credential replacement and concurrency behavior as Login. If user-data
// decoding fails after a valid session ID is returned, that session is logged
// out with a bounded cleanup context and existing credentials are preserved.
func (c *Client) LoginWithUserData(ctx context.Context, username, password string) (*Session, error) {
	return c.login(ctx, username, password, true)
}

func (c *Client) login(ctx context.Context, username, password string, userData bool) (*Session, error) {
	if ctx == nil {
		return nil, requestError("user.login", ErrorInvalidParams, errors.New("context must not be nil"))
	}
	c.mu.Lock()
	c.generation++
	generation := c.generation
	c.mu.Unlock()

	params := loginParams{Username: username, Password: password, UserData: userData}
	var raw json.RawMessage
	if err := c.doWithToken(ctx, "user.login", params, &raw, ""); err != nil {
		return nil, err
	}
	var sessionID string
	if userData {
		// Recover the session independently of the user fields so a failed full
		// decode can still clean up the session created by this login.
		var result struct {
			SessionID string `json:"sessionid"`
		}
		if err := json.Unmarshal(raw, &result); err != nil {
			return nil, requestError("user.login", ErrorDecode, err)
		}
		sessionID = result.SessionID
	} else if err := json.Unmarshal(raw, &sessionID); err != nil {
		return nil, requestError("user.login", ErrorDecode, err)
	}
	if err := validateToken(sessionID); err != nil {
		return nil, requestError("user.login", ErrorProtocol, err)
	}
	session := &Session{SessionID: sessionID}
	if userData {
		if err := json.Unmarshal(raw, session); err != nil {
			cleanupErr := c.cleanupSession(ctx, sessionID)
			return nil, requestError("user.login", ErrorDecode, errors.Join(err, cleanupErr))
		}
	}
	c.mu.Lock()
	if c.generation == generation {
		c.credentials = credentials{kind: credentialSession, token: session.SessionID}
		c.mu.Unlock()
		return session, nil
	}
	c.mu.Unlock()
	cleanupErr := c.cleanupSession(ctx, sessionID)
	return nil, requestError("user.login", ErrorAuth, errors.Join(ErrAuthStateChanged, cleanupErr))
}

func (c *Client) cleanupSession(ctx context.Context, token string) error {
	cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), authCleanupTimeout)
	defer cancel()
	return c.logoutSession(cleanupCtx, token)
}

// Logout invalidates the current session and clears it only if no newer
// authentication operation or credential setter has superseded it. API tokens
// return ErrNotSession without a request. With no credentials, Logout succeeds
// and supersedes any pending login. A false server result leaves credentials set.
func (c *Client) Logout(ctx context.Context) error {
	if ctx == nil {
		return requestError("user.logout", ErrorInvalidParams, errors.New("context must not be nil"))
	}
	c.mu.Lock()
	if c.credentials.kind == credentialAPI {
		c.mu.Unlock()
		return requestError("user.logout", ErrorAuth, ErrNotSession)
	}
	c.generation++
	generation := c.generation
	state := c.credentials
	c.mu.Unlock()
	if state.kind == credentialNone {
		return nil
	}
	if err := c.logoutSession(ctx, state.token); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.generation != generation {
		return requestError("user.logout", ErrorAuth, ErrAuthStateChanged)
	}
	c.credentials = credentials{}
	return nil
}

func (c *Client) logoutSession(ctx context.Context, token string) error {
	var ok bool
	if err := c.doWithToken(ctx, "user.logout", []any{}, &ok, token); err != nil {
		return err
	}
	if !ok {
		return requestError("user.logout", ErrorProtocol, errors.New("logout result is not true"))
	}
	return nil
}
