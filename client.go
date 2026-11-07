// Package gozabbix is a client library for the Zabbix JSON-RPC API.
//
// It targets Zabbix 7.4 and authenticates using the HTTP `Authorization: Bearer`
// header, which works both for session tokens obtained via user.login and for
// pre-created API tokens.
//
// Author: Resul Çelik 2024
package gozabbix

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	jsonrpcVersion          = "2.0"
	defaultTimeout          = 30 * time.Second
	defaultUserAgent        = "gozabbix"
	defaultMaxResponseBytes = 64 << 20 // 64 MiB
)

// Client is a Zabbix API client. A zero Client is not usable; construct one with
// New. All methods are safe for concurrent use by multiple goroutines.
type Client struct {
	url        string
	userAgent  string
	httpClient *http.Client
	maxBytes   int64
	logger     *slog.Logger

	id atomic.Int64

	mu    sync.RWMutex
	token string // session id (user.login) or pre-created API token
}

// New creates a client for the given Zabbix API endpoint, e.g.
// "https://zabbix.example.com/api_jsonrpc.php". The client performs no network
// I/O until an API method is called; authenticate afterwards with Login, or
// supply a pre-created API token via WithToken.
func New(apiURL string, opts ...Option) (*Client, error) {
	if strings.TrimSpace(apiURL) == "" {
		return nil, errors.New("gozabbix: api url must not be empty")
	}
	c := &Client{
		url:       apiURL,
		userAgent: defaultUserAgent,
		maxBytes:  defaultMaxResponseBytes,
		logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
		httpClient: &http.Client{
			Timeout:   defaultTimeout,
			Transport: &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}},
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	// Keep redirects from forwarding credentials, even with an injected client.
	// Copy client settings so the caller's redirect policy is not modified.
	httpClient := *c.httpClient
	httpClient.CheckRedirect = func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	}
	c.httpClient = &httpClient
	return c, nil
}

// SetToken sets the authentication token (session id or API token) used for
// subsequent requests. Login and Logout call this automatically.
func (c *Client) SetToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = token
}

// Token returns the current authentication token, or "" if none is set.
func (c *Client) Token() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.token
}

// Authenticated reports whether a token is currently set.
func (c *Client) Authenticated() bool { return c.Token() != "" }

// APIError is the structured error object returned by the Zabbix API.
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (e *APIError) Error() string {
	if e.Data == "" {
		return fmt.Sprintf("zabbix api error %d: %s", e.Code, e.Message)
	}
	return fmt.Sprintf("zabbix api error %d: %s (%s)", e.Code, e.Message, e.Data)
}

// AsAPIError extracts an *APIError from err if the error (or one it wraps) is a
// Zabbix API error, distinguishing server-side rejections from transport errors.
func AsAPIError(err error) (*APIError, bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr, true
	}
	return nil, false
}

type rpcRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
	ID      int64  `json:"id"`
}

type rpcResponse struct {
	Jsonrpc string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   json.RawMessage `json:"error"`
	ID      *int64          `json:"id"`
}

// Call invokes an arbitrary authenticated Zabbix API method and decodes the
// result into out (which may be nil). A non-nil out must be a non-nil pointer;
// it is replaced only after successful typed decoding. It is the escape hatch
// for methods this library does not yet wrap with a typed helper.
func (c *Client) Call(ctx context.Context, method string, params, out any) error {
	return c.do(ctx, method, params, out, true)
}

func (c *Client) do(ctx context.Context, method string, params, out any, authenticate bool) error {
	destination := reflect.ValueOf(out)
	if out != nil && (destination.Kind() != reflect.Pointer || destination.IsNil()) {
		return errors.New("gozabbix: output must be nil or a non-nil pointer")
	}
	if params == nil {
		params = []any{} // Zabbix expects a present params field (empty array for parameterless methods).
	}
	reqBody := rpcRequest{
		Jsonrpc: jsonrpcVersion,
		Method:  method,
		Params:  params,
		ID:      c.id.Add(1),
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("gozabbix: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("gozabbix: build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json-rpc")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", c.userAgent)
	if authenticate {
		if tok := c.Token(); tok != "" {
			httpReq.Header.Set("Authorization", "Bearer "+tok)
		}
	}

	c.logger.DebugContext(ctx, "zabbix request", "method", method, "id", reqBody.ID)

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("gozabbix: %s: %w", method, err)
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("gozabbix: %s: unexpected http status %d", method, httpResp.StatusCode)
	}
	limited := &io.LimitedReader{R: httpResp.Body, N: c.maxBytes}
	raw, err := io.ReadAll(limited)
	if err != nil {
		return fmt.Errorf("gozabbix: read response: %w", err)
	}
	if limited.N == 0 {
		// Probe one additional byte: a truncated prefix may itself be valid JSON.
		// Reading separately also avoids overflowing maxBytes + 1 for large limits.
		var extra [1]byte
		n, err := io.ReadFull(httpResp.Body, extra[:])
		if n != 0 {
			return fmt.Errorf("gozabbix: response exceeds %d byte limit", c.maxBytes)
		}
		if err != nil && err != io.EOF {
			return fmt.Errorf("gozabbix: read response: %w", err)
		}
	}

	var resp rpcResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return fmt.Errorf("gozabbix: decode response: %w", err)
	}
	if err := resp.validate(reqBody.ID); err != nil {
		return err
	}
	if out != nil {
		decoded := reflect.New(destination.Type().Elem())
		if err := json.Unmarshal(resp.Result, decoded.Interface()); err != nil {
			return fmt.Errorf("gozabbix: decode result of %s: %w", method, err)
		}
		destination.Elem().Set(decoded.Elem())
	}
	return nil
}

func (resp rpcResponse) validate(requestID int64) error {
	if resp.Jsonrpc != jsonrpcVersion {
		return errors.New("gozabbix: invalid JSON-RPC version")
	}
	if resp.ID == nil || *resp.ID != requestID {
		return errors.New("gozabbix: invalid JSON-RPC response id")
	}
	if (len(resp.Result) == 0) == (len(resp.Error) == 0) {
		return errors.New("gozabbix: response must contain exactly one result or error")
	}
	if len(resp.Error) != 0 {
		var apiError struct {
			Code    *int    `json:"code"`
			Message *string `json:"message"`
			Data    string  `json:"data"`
		}
		if err := json.Unmarshal(resp.Error, &apiError); err != nil || apiError.Code == nil || apiError.Message == nil {
			return errors.New("gozabbix: invalid JSON-RPC error")
		}
		return &APIError{Code: *apiError.Code, Message: *apiError.Message, Data: apiError.Data}
	}
	if bytes.Equal(bytes.TrimSpace(resp.Result), []byte("null")) {
		return errors.New("gozabbix: null JSON-RPC result")
	}
	return nil
}
