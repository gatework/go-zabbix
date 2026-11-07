package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

var errSensitiveJSON = errors.New("payload-password-secret")

type failingJSON struct{}

func (failingJSON) MarshalJSON() ([]byte, error) { return nil, errSensitiveJSON }
func (*failingJSON) UnmarshalJSON([]byte) error  { return errSensitiveJSON }

func TestRequestErrorHidesTransportTextAndPreservesCauses(t *testing.T) {
	var logs bytes.Buffer
	sensitive := errors.New("transport-password-secret")
	cause := errors.Join(sensitive, context.Canceled)
	c, err := New("https://example.invalid/api_jsonrpc.php?key=query-secret",
		WithAPIToken("token-secret"),
		WithLogger(slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelDebug}))),
		WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			_ = req.Body.Close()
			return nil, cause
		})}))
	if err != nil {
		t.Fatal(err)
	}
	err = c.Call(context.Background(), "host.get", map[string]string{"password": "request-password-secret"}, nil)
	var requestErr *RequestError
	var urlErr *url.Error
	if !errors.As(err, &requestErr) || requestErr.Method != "host.get" || requestErr.Kind != ErrorTransport {
		t.Fatalf("missing request context: %v", err)
	}
	if !errors.Is(err, sensitive) || !errors.Is(err, context.Canceled) || !errors.As(err, &urlErr) {
		t.Fatal("transport error chain was lost")
	}
	if strings.Contains(err.Error(), "secret") || strings.Contains(err.Error(), "example.invalid") || strings.Contains(logs.String(), "secret") || strings.Contains(logs.String(), "example.invalid") {
		t.Fatal("default error or logs disclosed endpoint, credential or payload")
	}
	if !strings.Contains(urlErr.URL, "query-secret") {
		t.Fatal("explicit error inspection lost original diagnostic detail")
	}
}

func TestLocalRequestErrorsDoNotSendRequests(t *testing.T) {
	var calls int
	c, err := New("https://example.invalid/", WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		calls++
		_ = req.Body.Close()
		return nil, errors.New("unexpected transport call")
	})}))
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name   string
		ctx    context.Context
		params any
		cause  error
	}{
		{"nil context", nil, nil, nil},
		{"marshal failure", context.Background(), failingJSON{}, errSensitiveJSON},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := c.Call(tc.ctx, "host.get", tc.params, nil)
			var requestErr *RequestError
			if !errors.As(err, &requestErr) || requestErr.Kind != ErrorInvalidParams || requestErr.Method != "host.get" {
				t.Fatalf("wrong local error: %v", err)
			}
			if tc.cause != nil && !errors.Is(err, tc.cause) {
				t.Fatal("marshal error cause was lost")
			}
			if strings.Contains(err.Error(), "secret") {
				t.Fatal("marshal error disclosed payload")
			}
		})
	}
	if calls != 0 {
		t.Fatal("invalid input caused network I/O")
	}
}

func TestRequestErrorPreservesDecodeAndAPIDataWithoutPrintingIt(t *testing.T) {
	for _, apiFailure := range []bool{false, true} {
		t.Run(map[bool]string{false: "decode", true: "api"}[apiFailure], func(t *testing.T) {
			c := responseServer(t, func(req rpcRequest) string {
				if apiFailure {
					return `{"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":"message-secret","data":"data-secret"}}`
				}
				return `{"jsonrpc":"2.0","id":1,"result":{"password":"body-secret"}}`
			})
			out := failingJSON{}
			err := c.Call(context.Background(), "host.get", nil, &out)
			var requestErr *RequestError
			if !errors.As(err, &requestErr) {
				t.Fatalf("missing request error: %v", err)
			}
			if strings.Contains(err.Error(), "secret") {
				t.Fatal("error printed response content")
			}
			if apiFailure {
				apiErr, ok := AsAPIError(err)
				if !ok || requestErr.Kind != ErrorAPI || apiErr.Code != -32602 || apiErr.Message != "message-secret" || apiErr.Data != "data-secret" {
					t.Fatal("structured API error was lost")
				}
				if strings.Contains(apiErr.Error(), "secret") {
					t.Fatal("API error default text printed server data")
				}
			} else if requestErr.Kind != ErrorDecode || !errors.Is(err, errSensitiveJSON) {
				t.Fatal("decode cause was lost")
			}
		})
	}
}

func TestHTTPErrorIsTypedAndBodyIsClosed(t *testing.T) {
	var closed bool
	c, err := New("https://example.invalid/?key=query-secret", WithHTTPClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		_ = req.Body.Close()
		return &http.Response{StatusCode: http.StatusTooManyRequests, Header: make(http.Header), Body: onCloseBody{
			ReadCloser: io.NopCloser(strings.NewReader("response-secret")), close: func() { closed = true },
		}}, nil
	})}))
	if err != nil {
		t.Fatal(err)
	}
	err = c.Call(context.Background(), "host.get", nil, nil)
	var httpErr *HTTPError
	var requestErr *RequestError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusTooManyRequests || !errors.As(err, &requestErr) || requestErr.Kind != ErrorHTTP {
		t.Fatalf("wrong HTTP error: %v", err)
	}
	if !closed || strings.Contains(err.Error(), "secret") {
		t.Fatal("HTTP error leaked body or failed to close it")
	}
}

func TestDecodeErrorKeepsJSONTypeInformation(t *testing.T) {
	c := responseServer(t, func(req rpcRequest) string { return `{"jsonrpc":"2.0","id":1,"result":42}` })
	out := "unchanged"
	err := c.Call(context.Background(), "host.get", nil, &out)
	var typeErr *json.UnmarshalTypeError
	if !errors.As(err, &typeErr) || out != "unchanged" {
		t.Fatal("decode error lost its type or changed output")
	}
}
