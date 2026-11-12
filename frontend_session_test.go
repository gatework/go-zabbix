package zabbix_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gatework/go-zabbix"
)

func TestFrontend_ClearSessionSupersedesLogin(t *testing.T) {
	t.Parallel()
	started, release := make(chan struct{}), make(chan struct{})
	var revoked atomic.Int32
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/zabbix/index.php" {
			close(started)
			<-release
			http.SetCookie(w, &http.Cookie{Name: "zbx_session", Value: frontendCookie(), Path: "/zabbix/"})
			http.Redirect(w, r, "/zabbix/zabbix.php", http.StatusSeeOther)
			return
		}
		var request struct {
			Method string `json:"method"`
			ID     int64  `json:"id"`
		}
		if r.URL.Path != "/zabbix/api_jsonrpc.php" || json.NewDecoder(r.Body).Decode(&request) != nil || request.Method != "user.logout" {
			t.Error("expected orphan session cleanup")
		}
		revoked.Add(1)
		_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": request.ID, "result": true})
	}))
	t.Cleanup(server.Close)
	// Release the handler before server.Close, including on an early test failure.
	t.Cleanup(func() {
		select {
		case <-release:
		default:
			close(release)
		}
	})
	client, err := zabbix.NewFrontend(server.URL+"/zabbix/", zabbix.FrontendOptions{HTTPClient: server.Client()})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- client.Login(ctx, "monitor", "secret") }()
	select {
	case <-started:
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	if err := client.ClearSession(); err != nil {
		t.Fatal(err)
	}
	close(release)
	select {
	case err := <-done:
		if !errors.Is(err, zabbix.ErrFrontendAuthentication) {
			t.Fatalf("superseded login: %v", err)
		}
	case <-ctx.Done():
		t.Fatal(ctx.Err())
	}
	if revoked.Load() != 1 {
		t.Fatalf("cleanup requests = %d", revoked.Load())
	}
	if _, err := client.Chart(ctx, chartRequest()); !errors.Is(err, zabbix.ErrFrontendAuthentication) {
		t.Fatalf("late login restored cleared cookies: %v", err)
	}
}

func TestFrontend_CookiesArePrivateAndFailedLogoutRetainsSession(t *testing.T) {
	t.Parallel()
	picture := chartPNG(t, 940, 280)
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/zabbix/api_jsonrpc.php" {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		if _, err := r.Cookie("zbx_session"); err != nil {
			t.Error("private cookie disappeared")
		}
		w.Header().Set("Content-Type", "image/png")
		_, _ = w.Write(picture)
	}))
	t.Cleanup(server.Close)
	hc := server.Client()
	hc.Jar, _ = cookiejar.New(nil)
	base, _ := url.Parse(server.URL + "/zabbix/")
	hc.Jar.SetCookies(base, []*http.Cookie{{Name: "zbx_session", Value: frontendCookie()}})
	client, err := zabbix.NewFrontend(base.String(), zabbix.FrontendOptions{HTTPClient: hc})
	if err != nil {
		t.Fatal(err)
	}
	// Mutating the caller jar after construction cannot erase the private copy.
	hc.Jar.SetCookies(base, []*http.Cookie{{Name: "zbx_session", MaxAge: -1}})
	if err := client.Logout(t.Context()); err == nil {
		t.Fatal("failed revoke accepted")
	}
	if _, err := client.Chart(t.Context(), chartRequest()); err != nil {
		t.Fatal(err)
	}
	if len(hc.Jar.Cookies(base)) != 0 {
		t.Fatal("private operations mutated caller cookies")
	}
	if err := client.ClearSession(); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Chart(t.Context(), chartRequest()); !errors.Is(err, zabbix.ErrFrontendAuthentication) {
		t.Fatal(err)
	}
}

func TestClient_ZeroValueRequestsReturnErrors(t *testing.T) {
	t.Parallel()
	var client zabbix.Client
	for name, invoke := range map[string]func() error{
		"raw":     func() error { return client.Call(t.Context(), "host.get", nil, nil) },
		"version": func() error { _, err := client.APIInfo().Version(t.Context()); return err },
		"login":   func() error { return client.Login(t.Context(), "monitor", "secret") },
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var requestErr *zabbix.RequestError
			if err := invoke(); !errors.As(err, &requestErr) || requestErr.Kind != zabbix.ErrorInvalidParams {
				t.Fatalf("zero client = %v", err)
			}
		})
	}
}
