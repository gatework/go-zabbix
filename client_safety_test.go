package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func responseServer(t *testing.T, response func(rpcRequest) string, opts ...Option) *Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req rpcRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Error(err)
			return
		}
		_, _ = io.WriteString(w, response(req))
	}))
	t.Cleanup(server.Close)
	client, err := New(server.URL, opts...)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestResponseLimitRejectsTruncatedValidJSON(t *testing.T) {
	for _, extra := range []string{"", " ", "x"} {
		t.Run(fmt.Sprintf("extra=%q", extra), func(t *testing.T) {
			const body = `{"jsonrpc":"2.0","id":1,"result":[]}`
			client := responseServer(t, func(rpcRequest) string { return body + extra }, WithMaxResponseBytes(int64(len(body))))
			var out []string
			err := client.Call(context.Background(), "host.get", nil, &out)
			if extra == "" && (err != nil || out == nil) {
				t.Fatalf("exact-limit empty result = %#v, error = %v", out, err)
			}
			if extra != "" && (err == nil || !strings.Contains(err.Error(), "limit")) {
				t.Fatalf("overflow must be distinguished from JSON decoding, got %v", err)
			}
		})
	}
}

func TestResponseEnvelopeValidation(t *testing.T) {
	cases := map[string]string{
		"invalid JSON":      `{`,
		"trailing JSON":     `{"jsonrpc":"2.0","id":%d,"result":[]} {}`,
		"wrong version":     `{"jsonrpc":"1.0","id":%d,"result":[]}`,
		"missing version":   `{"id":%d,"result":[]}`,
		"missing id":        `{"jsonrpc":"2.0","result":[]}`,
		"null id":           `{"jsonrpc":"2.0","id":null,"result":[]}`,
		"wrong id":          `{"jsonrpc":"2.0","id":-42,"result":[]}`,
		"string id":         `{"jsonrpc":"2.0","id":"1","result":[]}`,
		"fractional id":     `{"jsonrpc":"2.0","id":1.5,"result":[]}`,
		"missing result":    `{"jsonrpc":"2.0","id":%d}`,
		"null result":       `{"jsonrpc":"2.0","id":%d,"result":null}`,
		"null error":        `{"jsonrpc":"2.0","id":%d,"error":null}`,
		"empty error":       `{"jsonrpc":"2.0","id":%d,"error":{}}`,
		"incomplete error":  `{"jsonrpc":"2.0","id":%d,"error":{"code":-1}}`,
		"result and error":  `{"jsonrpc":"2.0","id":%d,"result":[],"error":{"code":-1,"message":"bad"}}`,
		"result null error": `{"jsonrpc":"2.0","id":%d,"result":[],"error":null}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			client := responseServer(t, func(req rpcRequest) string {
				return strings.ReplaceAll(body, "%d", strconv.FormatInt(req.ID, 10))
			})
			out := []string{"unchanged"}
			err := client.Call(context.Background(), "host.get", nil, &out)
			if err == nil {
				t.Fatal("invalid response envelope was accepted")
			}
			if _, ok := AsAPIError(err); ok {
				t.Fatal("malformed envelope was accepted as a valid API error")
			}
			if !reflect.DeepEqual(out, []string{"unchanged"}) {
				t.Fatalf("invalid envelope changed output: %#v", out)
			}
			if err := client.Call(context.Background(), "host.get", nil, nil); err == nil {
				t.Fatal("nil output bypassed envelope validation")
			}
		})
	}
}

func TestTypedDecodeDoesNotPartiallyOverwriteOutput(t *testing.T) {
	client := responseServer(t, func(req rpcRequest) string {
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":[{"hostid":"changed"},{"hostid":42}]}`, req.ID)
	})
	out := []Host{{HostID: "original"}}
	if err := client.Call(context.Background(), "host.get", nil, &out); err == nil {
		t.Fatal("incompatible host data was accepted")
	}
	if !reflect.DeepEqual(out, []Host{{HostID: "original"}}) {
		t.Fatalf("failed typed decoding changed output: %#v", out)
	}
}

func TestCallOutputContract(t *testing.T) {
	var requests atomic.Int32
	client := responseServer(t, func(req rpcRequest) string {
		requests.Add(1)
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":true}`, req.ID)
	})
	var nilPointer *bool
	for _, out := range []any{false, nilPointer} {
		if err := client.Call(context.Background(), "host.get", nil, out); err == nil {
			t.Errorf("invalid output %T was accepted", out)
		}
	}
	if requests.Load() != 0 {
		t.Fatal("invalid output sent a network request")
	}
	if err := client.Call(context.Background(), "host.get", nil, nil); err != nil {
		t.Fatalf("nil output must remain supported: %v", err)
	}
}

func TestRedirectsAreRefusedWithoutMutatingInjectedClient(t *testing.T) {
	for _, injected := range []bool{false, true} {
		t.Run(fmt.Sprintf("injected=%v", injected), func(t *testing.T) {
			var followed atomic.Int32
			destination := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				followed.Add(1)
				_, _ = io.WriteString(w, `{"jsonrpc":"2.0","id":1,"result":[]}`)
			}))
			defer destination.Close()
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, destination.URL, http.StatusTemporaryRedirect)
			}))
			defer server.Close()
			var redirectCalls atomic.Int32
			httpClient := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error {
				redirectCalls.Add(1)
				return nil
			}}
			opts := []Option{WithToken("token-secret")}
			if injected {
				opts = append(opts, WithHTTPClient(httpClient))
			}
			client, err := New(server.URL, opts...)
			if err != nil {
				t.Fatal(err)
			}
			var out []string
			if err := client.Call(context.Background(), "host.get", nil, &out); err == nil {
				t.Error("redirect was accepted")
			}
			if followed.Load() != 0 || redirectCalls.Load() != 0 {
				t.Error("SDK followed a redirect or used the caller redirect hook")
			}
			before := redirectCalls.Load()
			if err := httpClient.CheckRedirect(nil, nil); err != nil || redirectCalls.Load() != before+1 {
				t.Fatal("SDK mutated caller-owned redirect policy")
			}
		})
	}
}

func TestHTTPErrorBodyIsNotExposed(t *testing.T) {
	for _, status := range []int{401, 403, 429, 500} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
				_, _ = io.WriteString(w, "password-secret token-secret")
			}))
			defer server.Close()
			client, err := New(server.URL)
			if err != nil {
				t.Fatal(err)
			}
			err = client.Call(context.Background(), "host.get", nil, nil)
			if err == nil || !strings.Contains(err.Error(), strconv.Itoa(status)) || strings.Contains(err.Error(), "secret") {
				t.Fatalf("HTTP status error missing or contains response body: %v", err)
			}
		})
	}
}

func TestDefaultTLSVerificationAndInjectedTrust(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"jsonrpc":"2.0","id":1,"result":"7.0.25"}`)
	}))
	defer server.Close()
	server.Config.ErrorLog = log.New(io.Discard, "", 0)
	client, err := New(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.APIInfo().Version(context.Background()); err == nil {
		t.Fatal("default TLS accepted an untrusted certificate")
	}
	client, err = New(server.URL, WithHTTPClient(server.Client()))
	if err != nil {
		t.Fatal(err)
	}
	if version, err := client.APIInfo().Version(context.Background()); err != nil || version != "7.0.25" {
		t.Fatalf("injected trusted client was not used: version=%q, error=%v", version, err)
	}
}

func TestExistingBearerOptionsAndPayloadFreeLogging(t *testing.T) {
	var logs bytes.Buffer
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		var req map[string]json.RawMessage
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Error(err)
			return
		}
		if string(req["method"]) != `"host.get"` || req["auth"] != nil || r.Header.Get("Authorization") != "Bearer token-secret" {
			t.Error("existing explicit Bearer request contract changed")
		}
		_, _ = fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%s,"result":{"counter":9007199254740993,"name":"response-secret"}}`, req["id"])
	}))
	defer server.Close()
	client, err := New(server.URL, WithToken("token-secret"), WithTimeout(20*time.Second), WithMaxResponseBytes(4<<20),
		WithLogger(slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug}))))
	if err != nil {
		t.Fatal(err)
	}
	if requests.Load() != 0 {
		t.Fatal("construction performed a version probe")
	}
	var out struct {
		Counter int64 `json:"counter"`
	}
	if err := client.Call(context.Background(), "host.get", map[string]string{"search": "request-secret"}, &out); err != nil {
		t.Fatal(err)
	}
	if out.Counter != 9007199254740993 || requests.Load() != 1 {
		t.Fatal("typed result lost precision or sent extra requests")
	}
	if logs.Len() == 0 || strings.Contains(logs.String(), "secret") {
		t.Fatal("SDK logging lost metadata or exposed credentials/payloads")
	}
}

func TestCancellationDuringResponseRead(t *testing.T) {
	started := make(chan struct{})
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(io.Discard, r.Body)
		_, _ = io.WriteString(w, `{"jsonrpc":"2.0","result":`)
		w.(http.Flusher).Flush()
		close(started)
		select {
		case <-r.Context().Done():
		case <-release:
		}
	}))
	defer server.Close()
	defer close(release)
	client, err := New(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan error, 1)
	go func() { done <- client.Call(ctx, "host.get", nil, nil) }()
	<-started
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("body cancellation error = %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("response reading ignored caller cancellation")
	}
}

func TestConcurrentCallsKeepDistinctResponseIDs(t *testing.T) {
	var ids sync.Map
	client := responseServer(t, func(req rpcRequest) string {
		if _, duplicate := ids.LoadOrStore(req.ID, true); duplicate {
			t.Errorf("duplicate id %d", req.ID)
		}
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":[]}`, req.ID)
	}, WithToken("token"))
	var wg sync.WaitGroup
	for i := 0; i < 24; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var out []Host
			if err := client.Call(context.Background(), "host.get", nil, &out); err != nil || out == nil {
				t.Errorf("concurrent call = %#v, %v", out, err)
			}
		}()
	}
	wg.Wait()
}
