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
func (w *WebClient) Login(ctx context.Context, username, password string) error {
	_, initialGeneration, err := w.snapshot(ctx)
	if err != nil || strings.TrimSpace(username) == "" || password == "" {
		if err != nil {
			return err
		}
		return WebErrorInvalidRequest
	}
	select {
	case w.loginGate <- struct{}{}:
		defer func() { <-w.loginGate }()
	case <-ctx.Done():
		return ctx.Err()
	}
	hc, generation, err := w.snapshot(ctx)
	if err != nil {
		return err
	}
	if generation != initialGeneration {
		return WebErrorAuthenticationRequired
	}
	candidate, err := freshWebClient(hc)
	if err != nil {
		return err
	}
	if err := w.submitLogin(ctx, candidate, username, password); err != nil {
		return err
	}
	w.mu.Lock()
	if generation != w.generation {
		w.mu.Unlock()
		cleanup, cancel := context.WithTimeout(context.WithoutCancel(ctx), authCleanupTimeout)
		defer cancel()
		return errors.Join(WebErrorAuthenticationRequired, w.logoutSession(cleanup, candidate))
	}
	w.client = candidate
	w.generation++
	w.mu.Unlock()
	return nil
}

func (w *WebClient) submitLogin(ctx context.Context, hc *http.Client, username, password string) error {
	values := url.Values{"name": {username}, "password": {password}, "enter": {"Sign in"}, "form": {"1"}, "autologin": {"0"}}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, w.endpoint("index.php").String(), strings.NewReader(values.Encode()))
	if err != nil {
		return requestError("web.login", ErrorKindInvalidParams, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := hc.Do(req)
	if err != nil {
		return requestError("web.login", ErrorKindTransport, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusSeeOther {
		return WebErrorAuthenticationRequired
	}
	location, err := resp.Location()
	if err != nil || !w.sameOrigin(location) {
		return WebErrorAuthenticationRequired
	}
	switch path.Base(location.Path) {
	case "index_mfa.php", "index_sso.php", "index_http.php":
		return WebErrorInteractiveLoginRequired
	case "index.php", "setup.php":
		return WebErrorAuthenticationRequired
	}
	if _, err := w.sessionID(hc); err != nil {
		return err
	}
	return nil
}

// Logout revokes the current frontend session through user.logout, then forgets
// its cookie. A failed revoke preserves the cookie so the caller may retry.
func (w *WebClient) Logout(ctx context.Context) error {
	hc, generation, err := w.snapshot(ctx)
	if err != nil {
		return err
	}
	if err := w.logoutSession(ctx, hc); err != nil {
		return err
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if generation == w.generation {
		next, err := freshWebClient(hc)
		if err != nil {
			return err
		}
		w.client = next
		w.generation++
	}
	return nil
}

func (w *WebClient) logoutSession(ctx context.Context, hc *http.Client) error {
	sessionID, err := w.sessionID(hc)
	if err != nil {
		return err
	}
	client, err := New(w.endpoint("api_jsonrpc.php").String(), WithHTTPClient(hc), WithSessionToken(sessionID))
	if err != nil {
		return WebErrorInvalidRequest
	}
	return client.Logout(ctx)
}

// Only read the ID from a server-issued signed cookie; do not construct or sign
// a cookie from an API token. Zabbix itself authenticates it on every request.
func (w *WebClient) sessionID(hc *http.Client) (string, error) {
	for _, cookie := range hc.Jar.Cookies(w.endpoint("chart.php")) {
		if cookie.Name != "zbx_session" {
			continue
		}
		value := cookie.Value
		if len(value) > 16<<10 {
			return "", WebErrorAuthenticationRequired
		}
		if strings.Contains(value, "%") {
			decoded, err := url.PathUnescape(value)
			if err != nil {
				return "", WebErrorAuthenticationRequired
			}
			value = decoded
		}
		data, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			return "", WebErrorAuthenticationRequired
		}
		var session struct {
			ID   string `json:"sessionid"`
			Sign string `json:"sign"`
		}
		if json.Unmarshal(data, &session) != nil || session.Sign == "" || validateToken(session.ID) != nil {
			return "", WebErrorAuthenticationRequired
		}
		return session.ID, nil
	}
	return "", WebErrorAuthenticationRequired
}

func (w *WebClient) sameOrigin(target *url.URL) bool {
	return target.User == nil && strings.EqualFold(target.Scheme, w.base.Scheme) &&
		strings.EqualFold(target.Host, w.base.Host) && strings.HasPrefix(target.Path, w.base.Path)
}
