package zabbix

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Login obtains a genuine signed frontend cookie through the Zabbix login form.
// API tokens and user.login results cannot substitute for this cookie. SSO/MFA
// callers must authenticate externally and supply cookies through HTTPClient.
// A failed login leaves the previous session intact. Login does not retry.
// Log out first if a previously installed session must be revoked.
func (f *Frontend) Login(ctx context.Context, username, password string) error {
	_, initialGeneration, err := f.snapshot(ctx)
	if err != nil || strings.TrimSpace(username) == "" || password == "" {
		if err != nil {
			return err
		}
		return ErrFrontendInput
	}
	select {
	case f.loginGate <- struct{}{}:
		defer func() { <-f.loginGate }()
	case <-ctx.Done():
		return ctx.Err()
	}
	hc, generation, err := f.snapshot(ctx)
	if err != nil {
		return err
	}
	if generation != initialGeneration {
		return ErrFrontendAuthentication
	}
	candidate, err := freshFrontendClient(hc)
	if err != nil {
		return err
	}
	if err := f.submitLogin(ctx, candidate, username, password); err != nil {
		return err
	}
	f.mu.Lock()
	if generation != f.generation {
		f.mu.Unlock()
		cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), authCleanupTimeout)
		defer cancel()
		return errors.Join(ErrFrontendAuthentication, f.logoutSession(cleanup, candidate))
	}
	f.client = candidate
	f.generation++
	f.mu.Unlock()
	return nil
}

func (f *Frontend) submitLogin(ctx context.Context, hc *http.Client, username, password string) error {
	values := url.Values{"name": {username}, "password": {password}, "enter": {"Sign in"}, "form": {"1"}, "autologin": {"0"}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, f.endpoint("index.php").String(), strings.NewReader(values.Encode()))
	if err != nil {
		return requestError("frontend.login", ErrorInvalidParams, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := hc.Do(req)
	if err != nil {
		return requestError("frontend.login", ErrorTransport, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusSeeOther {
		return ErrFrontendAuthentication
	}
	location, err := resp.Location()
	if err != nil || !f.sameOrigin(location) {
		return ErrFrontendAuthentication
	}
	switch path.Base(location.Path) {
	case "index_mfa.php", "index_sso.php", "index_http.php":
		return ErrFrontendInteractiveLogin
	case "index.php", "setup.php":
		return ErrFrontendAuthentication
	}
	if _, err := f.sessionID(hc); err != nil {
		return err
	}
	return nil
}

// Logout revokes the current frontend session through user.logout, then forgets
// its cookie. A failed revoke preserves the cookie so the caller may retry.
func (f *Frontend) Logout(ctx context.Context) error {
	hc, generation, err := f.snapshot(ctx)
	if err != nil {
		return err
	}
	if err := f.logoutSession(ctx, hc); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	if generation == f.generation {
		next, err := freshFrontendClient(hc)
		if err != nil {
			return err
		}
		f.client = next
		f.generation++
	}
	return nil
}

func (f *Frontend) logoutSession(ctx context.Context, hc *http.Client) error {
	sessionID, err := f.sessionID(hc)
	if err != nil {
		return err
	}
	client, err := New(f.endpoint("api_jsonrpc.php").String(), WithHTTPClient(hc), WithSessionToken(sessionID))
	if err != nil {
		return ErrFrontendInput
	}
	return client.Logout(ctx)
}

// Only read the ID from a server-issued signed cookie; do not construct or sign
// a cookie from an API token. Zabbix itself authenticates it on every request.
func (f *Frontend) sessionID(hc *http.Client) (string, error) {
	for _, cookie := range hc.Jar.Cookies(f.endpoint("chart.php")) {
		if cookie.Name != "zbx_session" {
			continue
		}
		value := cookie.Value
		if len(value) > 16<<10 {
			return "", ErrFrontendAuthentication
		}
		if strings.Contains(value, "%") {
			decoded, err := url.PathUnescape(value)
			if err != nil {
				return "", ErrFrontendAuthentication
			}
			value = decoded
		}
		data, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return "", ErrFrontendAuthentication
		}
		var session struct {
			ID   string `json:"sessionid"`
			Sign string `json:"sign"`
		}
		if json.Unmarshal(data, &session) != nil || session.Sign == "" || validateToken(session.ID) != nil {
			return "", ErrFrontendAuthentication
		}
		return session.ID, nil
	}
	return "", ErrFrontendAuthentication
}

func (f *Frontend) sameOrigin(target *url.URL) bool {
	return target.User == nil && strings.EqualFold(target.Scheme, f.base.Scheme) &&
		strings.EqualFold(target.Host, f.base.Host) && strings.HasPrefix(target.Path, f.base.Path)
}
