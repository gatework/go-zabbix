package zabbix

import (
	"errors"
	"fmt"
)

// ErrorKind identifies the stage or class of a failed request.
type ErrorKind string

const (
	ErrorInvalidParams ErrorKind = "invalid parameters"
	ErrorTransport     ErrorKind = "transport error"
	ErrorHTTP          ErrorKind = "HTTP error"
	ErrorAPI           ErrorKind = "API error"
	ErrorProtocol      ErrorKind = "protocol error"
	ErrorDecode        ErrorKind = "decode error"
	ErrorResponseLimit ErrorKind = "response limit exceeded"
	ErrorAuth          ErrorKind = "authentication error"
)

// ErrAuthStateChanged means a newer authentication operation or credential
// setter superseded an in-flight operation. Its result was not installed.
var ErrAuthStateChanged = errors.New("zabbix: authentication state changed")

// ErrNotSession means Logout was called with API-token credentials.
var ErrNotSession = errors.New("zabbix: API token credentials cannot be logged out")

// RequestError attaches the RPC method and failure category while retaining the
// cause for errors.Is/errors.As. Error deliberately omits the cause's text: it
// can contain endpoint URLs, credentials or upstream response data. Inspect Err
// only at a boundary that is allowed to handle those details.
type RequestError struct {
	Method string
	Kind   ErrorKind
	Err    error
}

func (e *RequestError) Error() string {
	detail := string(e.Kind)
	switch e.Kind {
	case ErrorHTTP:
		var httpErr *HTTPError
		if errors.As(e.Err, &httpErr) {
			detail = fmt.Sprintf("unexpected HTTP status %d", httpErr.StatusCode)
		}
	case ErrorAPI:
		if apiErr, ok := AsAPIError(e.Err); ok {
			detail = fmt.Sprintf("API error %d", apiErr.Code)
		}
	case ErrorAuth:
		if errors.Is(e.Err, ErrAuthStateChanged) {
			detail = "authentication state changed"
		}
	}
	if e.Method == "" {
		return "zabbix: " + detail
	}
	return "zabbix: " + e.Method + ": " + detail
}

func (e *RequestError) Unwrap() error { return e.Err }

func requestError(method string, kind ErrorKind, err error) error {
	return &RequestError{Method: method, Kind: kind, Err: err}
}

// HTTPError reports a non-200 HTTP status without retaining a response body.
type HTTPError struct {
	StatusCode int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("zabbix: unexpected HTTP status %d", e.StatusCode)
}

// APIError retains the structured Zabbix rejection. Message and Data are
// untrusted upstream text and may contain sensitive data. Error prints only Code.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (e *APIError) Error() string { return fmt.Sprintf("zabbix API error %d", e.Code) }

// AsAPIError extracts a Zabbix rejection through any wrapping error.
func AsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}
