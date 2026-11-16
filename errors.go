package zabbix

import (
	"cmp"
	"errors"
	"fmt"
)

// ErrorKind identifies the stage or class of a failed request.
type ErrorKind string

// Request failure categories.
const (
	ErrorKindUnknown        ErrorKind = ""
	ErrorKindInvalidParams  ErrorKind = "invalid parameters"
	ErrorKindTransport      ErrorKind = "transport error"
	ErrorKindHTTP           ErrorKind = "http error"
	ErrorKindAPI            ErrorKind = "api error"
	ErrorKindProtocol       ErrorKind = "protocol error"
	ErrorKindDecode         ErrorKind = "decode error"
	ErrorKindResponseLimit  ErrorKind = "response limit exceeded"
	ErrorKindAuthentication ErrorKind = "authentication error"
)

// AuthenticationError identifies a local credential lifecycle failure.
// Its comparable values support errors.Is; errors.AsType extracts the category.
type AuthenticationError string

// Credential lifecycle failures.
//
//nolint:errname // Typed error constants follow their XxxError type prefix.
const (
	// AuthenticationErrorUnknown is the unclassified zero value.
	AuthenticationErrorUnknown AuthenticationError = ""
	// AuthenticationErrorStateChanged means a newer operation superseded the result.
	AuthenticationErrorStateChanged AuthenticationError = "authentication state changed"
	// AuthenticationErrorSessionRequired means logout requires session credentials.
	AuthenticationErrorSessionRequired AuthenticationError = "session credentials required for logout"
)

func (e AuthenticationError) Error() string {
	return "zabbix: " + cmp.Or(string(e), "authentication error")
}

// MetricError identifies an invalid query, unavailable binding or malformed observation.
// It does not retain upstream payloads or application data.
type MetricError string

// Bounded metric-query failures.
//
//nolint:errname // Typed error constants follow their XxxError type prefix.
const (
	MetricErrorUnknown      MetricError = ""
	MetricErrorInvalidQuery MetricError = "invalid metric query"
	MetricErrorUnavailable  MetricError = "metric object unavailable"
	MetricErrorInvalidData  MetricError = "invalid metric data"
)

func (e MetricError) Error() string { return "zabbix: " + cmp.Or(string(e), "metric error") }

// WebError identifies a web request, session or chart-image failure.
// It does not retain credentials, cookies or response bodies.
type WebError string

// Web-session and chart failures.
//
//nolint:errname // Typed error constants follow their XxxError type prefix.
const (
	WebErrorUnknown                  WebError = ""
	WebErrorInvalidRequest           WebError = "invalid web request"
	WebErrorAuthenticationRequired   WebError = "web authentication required"
	WebErrorInteractiveLoginRequired WebError = "interactive web login required"
	WebErrorInvalidChartImage        WebError = "invalid web chart image"
)

func (e WebError) Error() string { return "zabbix: " + cmp.Or(string(e), "web error") }

// ValidationError describes an invalid local configuration or input field.
// SDK-created values contain field names and constraints, never input values.
type ValidationError string

func (e ValidationError) Error() string { return "zabbix: " + cmp.Or(string(e), "invalid parameters") }

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
	case ErrorKindUnknown:
		detail = "unknown request error"
	case ErrorKindHTTP:
		if httpErr, ok := errors.AsType[*HTTPStatusError](e.Err); ok {
			detail = fmt.Sprintf("unexpected http status %d", httpErr.StatusCode)
		}
	case ErrorKindAPI:
		if apiErr, ok := errors.AsType[*APIError](e.Err); ok {
			detail = fmt.Sprintf("api error %d", apiErr.Code)
		}
	case ErrorKindAuthentication:
		if errors.Is(e.Err, AuthenticationErrorStateChanged) {
			detail = "authentication state changed"
		}
	default:
		// Other categories use their safe metadata label.
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

// HTTPStatusError reports a non-200 HTTP status without retaining a response body.
type HTTPStatusError struct {
	StatusCode int
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("zabbix: unexpected http status %d", e.StatusCode)
}

// APIError retains the structured Zabbix rejection. Message and Data are
// untrusted upstream text and may contain sensitive data. Error prints only Code.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (e *APIError) Error() string { return fmt.Sprintf("zabbix: api error %d", e.Code) }
