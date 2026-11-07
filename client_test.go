package zabbix

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

// mockServer spins up an httptest server that decodes each JSON-RPC request and
// delegates to handler, which returns the value to place in the "result" field
// (or an *APIError to place in "error").
func mockServer(t *testing.T, handler func(req rpcRequest) (any, *APIError)) (*Client, *[]rpcRequest) {
	t.Helper()
	var mu sync.Mutex
	var got []rpcRequest

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var req rpcRequest
		if err := json.Unmarshal(body, &req); err != nil {
			t.Errorf("server: bad request json: %v", err)
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		// Preserve the raw Authorization header on the recorded request via params
		// wrapper is unnecessary; record the header separately through a side channel.
		req.Params = authHeader(r) // reuse Params field to smuggle the auth header to assertions
		mu.Lock()
		got = append(got, req)
		mu.Unlock()

		result, apiErr := handler(req)
		resp := map[string]any{"jsonrpc": "2.0", "id": req.ID}
		if apiErr != nil {
			resp["error"] = apiErr
		} else {
			resp["result"] = result
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	t.Cleanup(srv.Close)

	c, err := New(srv.URL)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return c, &got
}

func authHeader(r *http.Request) string { return r.Header.Get("Authorization") }

func TestNewValidatesURL(t *testing.T) {
	if _, err := New("   "); err == nil {
		t.Fatal("expected error for empty url")
	}
}

func TestAPIInfoVersionUnauthenticated(t *testing.T) {
	c, got := mockServer(t, func(req rpcRequest) (any, *APIError) {
		return "7.4.0", nil
	})
	if err := c.SetAPIToken("should-not-be-sent"); err != nil {
		t.Fatal(err)
	}

	v, err := c.APIInfo().Version(context.Background())
	if err != nil {
		t.Fatalf("Version: %v", err)
	}
	if v != "7.4.0" {
		t.Fatalf("version = %q, want 7.4.0", v)
	}
	if h := (*got)[0].Params.(string); h != "" {
		t.Fatalf("apiinfo.version must be unauthenticated, got Authorization %q", h)
	}
}

func TestLoginStoresTokenAndSendsBearer(t *testing.T) {
	const token = "abc123sessiontoken"
	c, got := mockServer(t, func(req rpcRequest) (any, *APIError) {
		if req.Method == "user.login" {
			return token, nil
		}
		return []Host{}, nil
	})

	if err := c.Login(context.Background(), "Admin", "zabbix"); err != nil {
		t.Fatalf("Login: %v", err)
	}
	if state, _ := c.authSnapshot(); state.token != token || state.kind != credentialSession {
		t.Fatal("login did not install session credentials")
	}

	// A subsequent call must carry the bearer token.
	if _, err := c.Host().Get(context.Background(), HostGetParams{}); err != nil {
		t.Fatalf("Host.Get: %v", err)
	}
	last := (*got)[len(*got)-1]
	if h := last.Params.(string); h != "Bearer "+token {
		t.Fatalf("Authorization = %q, want Bearer %s", h, token)
	}
}

func TestAPIErrorReturned(t *testing.T) {
	c, _ := mockServer(t, func(req rpcRequest) (any, *APIError) {
		return nil, &APIError{Code: -32602, Message: "Invalid params.", Data: "Not authorised."}
	})

	_, err := c.Host().Get(context.Background(), HostGetParams{})
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := AsAPIError(err)
	if !ok {
		t.Fatalf("AsAPIError = false, err = %v", err)
	}
	if apiErr.Code != -32602 {
		t.Fatalf("code = %d, want -32602", apiErr.Code)
	}
	if apiErr.Data != "Not authorised." || strings.Contains(err.Error(), "Not authorised") || strings.Contains(apiErr.Error(), "Not authorised") {
		t.Fatal("API data was lost or included in default error text")
	}
}

func TestHostGetDecodes(t *testing.T) {
	c, _ := mockServer(t, func(req rpcRequest) (any, *APIError) {
		return []Host{
			{HostID: "10084", Host: "Zabbix server", Name: "Zabbix server", Status: "0"},
		}, nil
	})

	hosts, err := c.Host().Get(context.Background(), HostGetParams{
		Output: []string{"hostid", "host", "name", "status"},
	})
	if err != nil {
		t.Fatalf("Host.Get: %v", err)
	}
	if len(hosts) != 1 || hosts[0].HostID != "10084" || hosts[0].Host != "Zabbix server" {
		t.Fatalf("unexpected hosts: %+v", hosts)
	}
}

func TestDeleteSendsIDArray(t *testing.T) {
	// Verify params for delete is a bare array of ids, per Zabbix conventions.
	var capturedParams json.RawMessage
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var probe struct {
			Params json.RawMessage `json:"params"`
		}
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &probe)
		capturedParams = probe.Params
		_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":1,"result":{"hostids":["10084"]}}`))
	}))
	t.Cleanup(srv.Close)

	c, _ := New(srv.URL)
	ids, err := c.Host().Delete(context.Background(), "10084")
	if err != nil {
		t.Fatalf("Host.Delete: %v", err)
	}
	if len(ids) != 1 || ids[0] != "10084" {
		t.Fatalf("ids = %v, want [10084]", ids)
	}
	if string(capturedParams) != `["10084"]` {
		t.Fatalf("params = %s, want [\"10084\"]", capturedParams)
	}
}

func TestContextCancellation(t *testing.T) {
	// The handler blocks until the client's context is cancelled (or the test
	// releases it on cleanup, so Server.Close never waits indefinitely).
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-release:
		}
	}))
	defer srv.Close()
	defer close(release) // LIFO: runs before srv.Close so the handler unblocks

	c, _ := New(srv.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := c.Host().Get(ctx, HostGetParams{})
	if err == nil {
		t.Fatal("expected context deadline error")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want context.DeadlineExceeded", err)
	}
}

func TestNon200Status(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "gateway down", http.StatusBadGateway)
	}))
	t.Cleanup(srv.Close)

	c, _ := New(srv.URL)
	_, err := c.Host().Get(context.Background(), HostGetParams{})
	if err == nil || !strings.Contains(err.Error(), "502") {
		t.Fatalf("err = %v, want http 502", err)
	}
}

func TestLogoutClearsToken(t *testing.T) {
	c, _ := mockServer(t, func(req rpcRequest) (any, *APIError) {
		return true, nil
	})
	if err := c.SetSessionToken("tok"); err != nil {
		t.Fatal(err)
	}
	if err := c.Logout(context.Background()); err != nil {
		t.Fatalf("Logout: %v", err)
	}
	if c.HasCredentials() {
		t.Fatal("token should be cleared after logout")
	}
}

func TestEventAcknowledge(t *testing.T) {
	c, got := mockServer(t, func(req rpcRequest) (any, *APIError) {
		return map[string][]string{"eventids": {"42"}}, nil
	})
	if err := c.SetAPIToken("tok"); err != nil {
		t.Fatal(err)
	}

	message := "looking into it"
	ids, err := c.Event().Acknowledge(context.Background(), EventAcknowledgeParams{
		EventIDs: []string{"42"},
		Action:   AckActionAcknowledge | AckActionAddMessage,
		Message:  &message,
	})
	if err != nil {
		t.Fatalf("Acknowledge: %v", err)
	}
	if len(ids) != 1 || ids[0] != "42" {
		t.Fatalf("ids = %v, want [42]", ids)
	}
	if (*got)[0].Method != "event.acknowledge" {
		t.Fatalf("method = %q", (*got)[0].Method)
	}
}

func TestWithAPITokenOption(t *testing.T) {
	c, err := New("https://example.invalid/api_jsonrpc.php", WithAPIToken("pretoken"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	state, _ := c.authSnapshot()
	if !c.HasCredentials() || state.token != "pretoken" || state.kind != credentialAPI {
		t.Fatalf("WithAPIToken did not set token")
	}
}
