package zabbix

import (
	"context"
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

type trackedTransport struct {
	http.RoundTripper
	closed atomic.Int32
}

func (tr *trackedTransport) CloseIdleConnections() { tr.closed.Add(1) }

func TestOptionsAreOrderIndependentAndDoNotMutateInputs(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"jsonrpc":"2.0","id":1,"result":"7.4"}`)
	}))
	defer server.Close()
	hc := server.Client()
	hc.Timeout = 11 * time.Second
	original := hc.Transport.(*http.Transport)
	tlsConfig := original.TLSClientConfig.Clone()
	tlsConfig.InsecureSkipVerify = true
	options := []Option{
		WithHTTPClient(hc), WithTimeout(19 * time.Second), WithTLSConfig(tlsConfig), WithInsecureSkipVerify(false),
	}
	for _, reverse := range []bool{false, true} {
		t.Run(map[bool]string{false: "forward", true: "reverse"}[reverse], func(t *testing.T) {
			ordered := append([]Option(nil), options...)
			if reverse {
				for i, j := 0, len(ordered)-1; i < j; i, j = i+1, j-1 {
					ordered[i], ordered[j] = ordered[j], ordered[i]
				}
			}
			c, err := New(server.URL, ordered...)
			if err != nil {
				t.Fatal(err)
			}
			defer c.CloseIdleConnections()
			if c.httpClient == hc || c.httpClient.Timeout != 19*time.Second || hc.Timeout != 11*time.Second {
				t.Fatal("timeout was lost or caller HTTP client was mutated")
			}
			owned := c.httpClient.Transport.(*http.Transport)
			if owned == original || owned.TLSClientConfig == tlsConfig || owned.TLSClientConfig == original.TLSClientConfig {
				t.Fatal("TLS tuning retained mutable caller-owned configuration")
			}
			if owned.TLSClientConfig.InsecureSkipVerify || !tlsConfig.InsecureSkipVerify || original.TLSClientConfig.InsecureSkipVerify {
				t.Fatal("TLS override order changed behavior or mutated an input")
			}
			if _, err := c.APIInfo().Version(context.Background()); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestTLSOptionsCannotWeakenAnotherClient(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"jsonrpc":"2.0","id":1,"result":"7.4"}`)
	}))
	defer server.Close()
	server.Config.ErrorLog = log.New(io.Discard, "", 0)
	shared := &tls.Config{MinVersion: tls.VersionTLS12}
	verified, err := New(server.URL, WithTLSConfig(shared))
	if err != nil {
		t.Fatal(err)
	}
	defer verified.CloseIdleConnections()
	unverified, err := New(server.URL, WithTLSConfig(shared), WithInsecureSkipVerify(true))
	if err != nil {
		t.Fatal(err)
	}
	defer unverified.CloseIdleConnections()
	if _, err := unverified.APIInfo().Version(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := verified.APIInfo().Version(context.Background()); err == nil {
		t.Fatal("configuring a second client disabled verification on the first")
	}
	if shared.InsecureSkipVerify {
		t.Fatal("caller TLS config changed")
	}
}

func TestBorrowedTransportIsPreservedAndNotClosed(t *testing.T) {
	var calls atomic.Int32
	tr := &trackedTransport{RoundTripper: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls.Add(1)
		_ = req.Body.Close()
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(`{"jsonrpc":"2.0","id":1,"result":true}`)), Header: make(http.Header)}, nil
	})}
	hc := &http.Client{Transport: tr, Timeout: time.Minute}
	c, err := New("https://example.invalid/api_jsonrpc.php", WithHTTPClient(hc))
	if err != nil {
		t.Fatal(err)
	}
	if c.httpClient == hc || c.httpClient.Transport != tr || c.httpClient.Timeout != time.Minute {
		t.Fatal("borrowing changed HTTP client settings or replaced its transport")
	}
	if err := c.Call(context.Background(), "host.get", nil, nil); err != nil {
		t.Fatal(err)
	}
	c.CloseIdleConnections()
	if calls.Load() != 1 || tr.closed.Load() != 0 {
		t.Fatal("borrowed transport was bypassed or closed")
	}
}

func TestCustomTransportRejectsTLSOptions(t *testing.T) {
	tr := &trackedTransport{}
	hc := &http.Client{Transport: tr}
	for _, option := range []Option{WithTLSConfig(&tls.Config{}), WithInsecureSkipVerify(false)} {
		for _, options := range [][]Option{{WithHTTPClient(hc), option}, {option, WithHTTPClient(hc)}} {
			c, err := New("https://example.invalid/", options...)
			if c != nil || err == nil {
				t.Fatal("TLS tuning silently replaced a custom transport")
			}
			if hc.Transport != tr {
				t.Fatal("failed construction mutated a borrowed client")
			}
		}
	}
}

type mapTransport map[string]string
type sliceTransport []string
type chanTransport chan struct{}
type valueTransport struct{}

func (mapTransport) RoundTrip(*http.Request) (*http.Response, error)   { panic("unexpected request") }
func (sliceTransport) RoundTrip(*http.Request) (*http.Response, error) { panic("unexpected request") }
func (chanTransport) RoundTrip(*http.Request) (*http.Response, error)  { panic("unexpected request") }
func (valueTransport) RoundTrip(*http.Request) (*http.Response, error) { panic("unexpected request") }

func TestNewRejectsTypedNilTransports(t *testing.T) {
	for _, tr := range []http.RoundTripper{
		(*http.Transport)(nil), roundTripFunc(nil), mapTransport(nil), sliceTransport(nil), chanTransport(nil),
	} {
		t.Run(reflect.TypeOf(tr).String(), func(t *testing.T) {
			hc := &http.Client{Transport: tr, Timeout: time.Minute}
			before := *hc
			c, err := New("https://example.invalid/", WithHTTPClient(hc))
			if c != nil || err == nil {
				t.Fatal("constructor accepted a typed-nil transport")
			}
			if !reflect.DeepEqual(*hc, before) {
				t.Fatal("failed construction modified the borrowed HTTP client")
			}
		})
	}
	for _, tr := range []http.RoundTripper{nil, valueTransport{}} {
		hc := &http.Client{Transport: tr}
		c, err := New("https://example.invalid/", WithHTTPClient(hc))
		if err != nil || c.httpClient.Transport != tr || hc.Transport != tr {
			t.Fatal("nil default or non-nilable custom transport was changed or rejected")
		}
	}
}

func TestOwnedIdleConnectionsCanBeReleased(t *testing.T) {
	closed := make(chan struct{}, 2)
	server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req rpcRequest
		if err := jsonDecodeRequest(r, &req); err != nil {
			t.Error(err)
			return
		}
		writeRPCResult(w, req.ID, true)
	}))
	server.Config.ConnState = func(_ net.Conn, state http.ConnState) {
		if state == http.StateClosed {
			closed <- struct{}{}
		}
	}
	server.Start()
	defer server.Close()
	c, err := New(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer c.CloseIdleConnections()
	if c.ownedTransport == nil || c.ownedTransport.IdleConnTimeout <= 0 {
		t.Fatal("owned transport has no idle bound")
	}
	if err := c.Call(context.Background(), "host.get", nil, nil); err != nil {
		t.Fatal(err)
	}
	c.CloseIdleConnections()
	select {
	case <-closed:
	case <-time.After(time.Second):
		t.Fatal("owned idle connection remained open")
	}
	if err := c.Call(context.Background(), "host.get", nil, nil); err != nil {
		t.Fatal("closing idle connections made client unusable:", err)
	}
}

func TestNewValidatesFinalConfiguration(t *testing.T) {
	invalidURLs := []string{"", "relative/path", "ftp://example.invalid", "https:///missing", "https://name:password-secret@example.invalid", "https://example.invalid/#fragment-secret", "https://example.invalid/#", "https://password-secret@%zz"}
	for _, endpoint := range invalidURLs {
		if c, err := New(endpoint); c != nil || err == nil || strings.Contains(err.Error(), "secret") {
			t.Fatalf("invalid URL accepted or exposed: client=%v error=%v", c != nil, err)
		}
	}
	invalidOptions := []Option{nil, WithHTTPClient(nil), WithTLSConfig(nil), WithTimeout(-time.Second), WithMaxResponseBytes(0), WithMaxResponseBytes(-1), WithAPIToken(" "), WithSessionToken("token\nsecret"), WithUserAgent("agent\nsecret")}
	for i, option := range invalidOptions {
		if c, err := New("https://example.invalid/?key=query-secret", option); c != nil || err == nil || strings.Contains(err.Error(), "secret") {
			t.Fatalf("invalid option %d accepted or exposed: client=%v error=%v", i, c != nil, err)
		}
	}
	c, err := New("https://example.invalid/", WithTimeout(-1), WithTimeout(0), WithMaxResponseBytes(0), WithMaxResponseBytes(1024), WithAPIToken(""), WithSessionToken("final-session"))
	if err != nil {
		t.Fatal("last value did not win:", err)
	}
	state, _ := c.authSnapshot()
	if c.httpClient.Timeout != 0 || c.maxBytes != 1024 || state.kind != credentialSession || state.token != "final-session" {
		t.Fatal("final configuration was not applied")
	}
	defaults, err := New("https://example.invalid/")
	if err != nil {
		t.Fatal(err)
	}
	if defaults.httpClient.Timeout != 30*time.Second || defaults.maxBytes != 64<<20 {
		t.Fatal("default limits changed")
	}
}
