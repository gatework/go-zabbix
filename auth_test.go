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
	"sync/atomic"
	"testing"
	"time"
)

func jsonDecodeRequest(r *http.Request, req *rpcRequest) error {
	return json.NewDecoder(r.Body).Decode(req)
}

func writeRPCResult(w http.ResponseWriter, id int64, result any) {
	_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": id, "result": result})
}

func authLatch(t *testing.T) (<-chan struct{}, func()) {
	t.Helper()
	ch := make(chan struct{})
	var once sync.Once
	release := func() { once.Do(func() { close(ch) }) }
	t.Cleanup(release)
	return ch, release
}

func awaitAuthResult(t *testing.T, result <-chan error) error {
	t.Helper()
	select {
	case err := <-result:
		return err
	case <-time.After(2 * time.Second):
		t.Fatal("authentication operation did not finish")
		return nil
	}
}

func TestLoginIsAnonymousWithExistingCredentials(t *testing.T) {
	for _, userData := range []bool{false, true} {
		t.Run(map[bool]string{false: "token", true: "user-data"}[userData], func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				var req rpcRequest
				if err := jsonDecodeRequest(r, &req); err != nil {
					t.Error(err)
					return
				}
				if req.Method == "user.login" {
					if r.Header.Get("Authorization") != "" {
						t.Error("login included old credentials")
					}
					if userData {
						writeRPCResult(w, req.ID, map[string]any{"sessionid": "new-session", "userid": "7", "type": 3, "mfaid": 0})
					} else {
						writeRPCResult(w, req.ID, "new-session")
					}
				} else {
					if r.Header.Get("Authorization") != "Bearer new-session" {
						t.Error("subsequent request did not use new session")
					}
					writeRPCResult(w, req.ID, true)
				}
			}))
			defer server.Close()
			c, err := New(server.URL, WithAPIToken("old-api-token"))
			if err != nil {
				t.Fatal(err)
			}
			if userData {
				session, err := c.LoginWithUserData(context.Background(), "Admin", "password-secret")
				if err != nil || session == nil || session.UserID != "7" || session.SessionID != "new-session" {
					t.Fatalf("login with user data: %v", err)
				}
			} else if err := c.Login(context.Background(), "Admin", "password-secret"); err != nil {
				t.Fatal(err)
			}
			if err := c.Call(context.Background(), "host.get", nil, nil); err != nil {
				t.Fatal(err)
			}
			if calls.Load() != 2 {
				t.Fatal("login performed extra requests")
			}
		})
	}
}

func TestInvalidAuthenticationResultsPreserveCredentials(t *testing.T) {
	cases := []struct {
		name, result string
		userData     bool
		kind         ErrorKind
	}{
		{"empty token", `""`, false, ErrorProtocol},
		{"blank token", `"  "`, false, ErrorProtocol},
		{"wrong token type", `{}`, false, ErrorDecode},
		{"missing session", `{}`, true, ErrorProtocol},
		{"blank session", `{"sessionid":" "}`, true, ErrorProtocol},
		{"wrong session type", `"session-secret"`, true, ErrorDecode},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := responseServer(t, func(req rpcRequest) string {
				return `{"jsonrpc":"2.0","id":1,"result":` + tc.result + `}`
			}, WithAPIToken("original-token"))
			var err error
			if tc.userData {
				_, err = c.LoginWithUserData(context.Background(), "Admin", "password")
			} else {
				err = c.Login(context.Background(), "Admin", "password")
			}
			var requestErr *RequestError
			if !errors.As(err, &requestErr) || requestErr.Kind != tc.kind {
				t.Fatalf("invalid login result: %v", err)
			}
			state, _ := c.authSnapshot()
			if state.kind != credentialAPI || state.token != "original-token" {
				t.Fatal("failed login changed credentials")
			}
		})
	}
}

func TestLogoutRequiresSessionAndSuccessfulResult(t *testing.T) {
	for _, tc := range []struct {
		name      string
		kind      credentialKind
		result    bool
		wantCalls int
		wantError bool
	}{
		{"no credentials", credentialNone, true, 0, false},
		{"api token", credentialAPI, true, 0, true},
		{"successful session", credentialSession, true, 1, false},
		{"false result", credentialSession, false, 1, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				calls.Add(1)
				var req rpcRequest
				if err := jsonDecodeRequest(r, &req); err != nil {
					t.Error(err)
					return
				}
				if req.Method != "user.logout" || r.Header.Get("Authorization") != "Bearer session" {
					t.Error("logout used wrong method or credentials")
				}
				writeRPCResult(w, req.ID, tc.result)
			}))
			defer server.Close()
			var options []Option
			if tc.kind == credentialAPI {
				options = append(options, WithAPIToken("session"))
			}
			if tc.kind == credentialSession {
				options = append(options, WithSessionToken("session"))
			}
			c, err := New(server.URL, options...)
			if err != nil {
				t.Fatal(err)
			}
			err = c.Logout(context.Background())
			if (err != nil) != tc.wantError || int(calls.Load()) != tc.wantCalls {
				t.Fatalf("logout: error=%v calls=%d", err, calls.Load())
			}
			if tc.kind == credentialAPI && !errors.Is(err, ErrNotSession) {
				t.Fatalf("missing API-token error: %v", err)
			}
			wantCredentials := tc.wantError && tc.kind != credentialNone
			if c.HasCredentials() != wantCredentials {
				t.Fatal("logout changed credentials despite failure")
			}
		})
	}
}

func TestLogoutCannotClearNewerCredentials(t *testing.T) {
	for _, sameSession := range []bool{false, true} {
		t.Run(map[bool]string{false: "new API token", true: "same session new generation"}[sameSession], func(t *testing.T) {
			started := make(chan struct{})
			release, unblock := authLatch(t)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var req rpcRequest
				if err := jsonDecodeRequest(r, &req); err != nil {
					t.Error(err)
					return
				}
				if req.Method == "user.logout" {
					if r.Header.Get("Authorization") != "Bearer original-session" {
						t.Error("logout did not snapshot original session")
					}
					close(started)
					<-release
				}
				writeRPCResult(w, req.ID, true)
			}))
			t.Cleanup(server.Close)
			t.Cleanup(unblock)
			c, err := New(server.URL, WithSessionToken("original-session"))
			if err != nil {
				t.Fatal(err)
			}
			done := make(chan error, 1)
			go func() { done <- c.Logout(context.Background()) }()
			<-started
			if sameSession {
				err = c.SetSessionToken("original-session")
			} else {
				err = c.SetAPIToken("new-api")
			}
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			if err := c.Call(ctx, "host.get", nil, nil); err != nil {
				t.Fatal("pending logout blocked normal requests:", err)
			}
			unblock()
			if err := awaitAuthResult(t, done); !errors.Is(err, ErrAuthStateChanged) {
				t.Fatalf("stale logout: %v", err)
			}
			state, _ := c.authSnapshot()
			if !c.HasCredentials() || sameSession && state.kind != credentialSession || !sameSession && state.token != "new-api" {
				t.Fatal("old logout cleared newer credentials")
			}
		})
	}
}

func TestLaterLoginReservationSupersedesEarlierLogin(t *testing.T) {
	firstStarted, secondStarted := make(chan struct{}), make(chan struct{})
	firstRelease, releaseFirst := authLatch(t)
	secondRelease, releaseSecond := authLatch(t)
	cleaned := make(chan string, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Method string      `json:"method"`
			ID     int64       `json:"id"`
			Params loginParams `json:"params"`
		}
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
			return
		}
		var method struct {
			Method string `json:"method"`
		}
		if err := json.Unmarshal(raw, &method); err != nil {
			t.Error(err)
			return
		}
		if method.Method == "user.login" {
			if err := json.Unmarshal(raw, &req); err != nil {
				t.Error(err)
				return
			}
			if req.Params.Username == "first" {
				close(firstStarted)
				<-firstRelease
			} else {
				close(secondStarted)
				<-secondRelease
			}
			writeRPCResult(w, req.ID, req.Params.Username+"-session")
		} else {
			var req rpcRequest
			if err := json.Unmarshal(raw, &req); err != nil {
				t.Error(err)
				return
			}
			cleaned <- r.Header.Get("Authorization")
			writeRPCResult(w, req.ID, true)
		}
	}))
	t.Cleanup(server.Close)
	t.Cleanup(releaseFirst)
	t.Cleanup(releaseSecond)
	c, err := New(server.URL, WithAPIToken("original"))
	if err != nil {
		t.Fatal(err)
	}
	first, second := make(chan error, 1), make(chan error, 1)
	go func() { first <- c.Login(context.Background(), "first", "password") }()
	<-firstStarted
	go func() { second <- c.Login(context.Background(), "second", "password") }()
	<-secondStarted
	releaseFirst()
	if err := awaitAuthResult(t, first); !errors.Is(err, ErrAuthStateChanged) {
		t.Fatalf("older login committed while newer login pending: %v", err)
	}
	if got := <-cleaned; got != "Bearer first-session" {
		t.Fatalf("cleanup logged out wrong session: %q", got)
	}
	state, _ := c.authSnapshot()
	if state.token != "original" {
		t.Fatal("pending login replaced original credentials")
	}
	releaseSecond()
	if err := awaitAuthResult(t, second); err != nil {
		t.Fatal(err)
	}
	state, _ = c.authSnapshot()
	if state.token != "second-session" || state.kind != credentialSession {
		t.Fatal("newer login was not installed")
	}
}

type onCloseBody struct {
	io.ReadCloser
	close func()
}

func (b onCloseBody) Close() error {
	err := b.ReadCloser.Close()
	b.close()
	return err
}

func TestStaleLoginCleanupIsBoundedAndPreservesErrors(t *testing.T) {
	for _, cleanupFails := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "cleanup failure"}[cleanupFails], func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var c *Client
			var calls int
			tr := roundTripFunc(func(req *http.Request) (*http.Response, error) {
				calls++
				var rpc rpcRequest
				if err := jsonDecodeRequest(req, &rpc); err != nil {
					return nil, err
				}
				_ = req.Body.Close()
				result := any(true)
				if rpc.Method == "user.login" {
					result = "orphan-session"
				}
				raw, err := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": rpc.ID, "result": result})
				if err != nil {
					return nil, err
				}
				resp := &http.Response{StatusCode: http.StatusOK, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(string(raw)))}
				if rpc.Method == "user.login" {
					resp.Body = onCloseBody{ReadCloser: resp.Body, close: func() {
						if err := c.SetAPIToken("new-api-token"); err != nil {
							t.Error(err)
						}
						cancel()
					}}
				} else {
					if req.Header.Get("Authorization") != "Bearer orphan-session" {
						t.Error("cleanup did not use stale session")
					}
					deadline, ok := req.Context().Deadline()
					if req.Context().Err() != nil || !ok || time.Until(deadline) > authCleanupTimeout || time.Until(deadline) <= 0 {
						t.Error("cleanup reused cancellation or has no finite deadline")
					}
					if cleanupFails {
						resp.StatusCode = http.StatusBadGateway
					}
				}
				return resp, nil
			})
			var err error
			c, err = New("https://example.invalid/", WithHTTPClient(&http.Client{Transport: tr}))
			if err != nil {
				t.Fatal(err)
			}
			err = c.Login(ctx, "Admin", "password-secret")
			if !errors.Is(err, ErrAuthStateChanged) {
				t.Fatalf("missing stale-state error: %v", err)
			}
			var httpErr *HTTPError
			if errors.As(err, &httpErr) != cleanupFails || cleanupFails && httpErr.StatusCode != http.StatusBadGateway {
				t.Fatalf("cleanup cause was lost: %v", err)
			}
			state, _ := c.authSnapshot()
			if state.kind != credentialAPI || state.token != "new-api-token" || calls != 2 {
				t.Fatal("cleanup changed new credentials or retried requests")
			}
		})
	}
}

func TestEmptyLogoutSupersedesPendingLogin(t *testing.T) {
	started := make(chan struct{})
	release, unblock := authLatch(t)
	var cleanup atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req rpcRequest
		if err := jsonDecodeRequest(r, &req); err != nil {
			t.Error(err)
			return
		}
		if req.Method == "user.login" {
			close(started)
			<-release
			writeRPCResult(w, req.ID, "orphan")
		} else {
			cleanup.Add(1)
			writeRPCResult(w, req.ID, true)
		}
	}))
	t.Cleanup(server.Close)
	t.Cleanup(unblock)
	c, err := New(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- c.Login(context.Background(), "Admin", "password") }()
	<-started
	if err := c.Logout(context.Background()); err != nil {
		t.Fatal(err)
	}
	unblock()
	if err := awaitAuthResult(t, done); !errors.Is(err, ErrAuthStateChanged) {
		t.Fatalf("pending login survived logout: %v", err)
	}
	if c.HasCredentials() || cleanup.Load() != 1 {
		t.Fatal("pending session was not cleaned up")
	}
}

func TestInvalidCredentialSetterDoesNotChangeState(t *testing.T) {
	c, err := New("https://example.invalid/", WithAPIToken("original"))
	if err != nil {
		t.Fatal(err)
	}
	before, generation := c.authSnapshot()
	for _, setter := range []func(string) error{c.SetAPIToken, c.SetSessionToken} {
		for _, token := range []string{"", " ", "token\nsecret"} {
			if err := setter(token); err == nil || strings.Contains(err.Error(), "secret") {
				t.Fatal("invalid credentials accepted or exposed")
			}
		}
	}
	after, gotGeneration := c.authSnapshot()
	if after != before || generation != gotGeneration {
		t.Fatal("invalid setter changed authentication state")
	}
	c.ClearCredentials()
	if c.HasCredentials() {
		t.Fatal("clear did not forget credentials")
	}
}
