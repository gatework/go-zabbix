// Package zabbix is a client library for the Zabbix JSON-RPC API.
//
// It targets Zabbix 7.4 and authenticates using the HTTP `Authorization: Bearer`
// header, which works both for session tokens obtained via user.login and for
// pre-created API tokens.
//
// Author: Resul Çelik 2024
package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	jsonrpcVersion          = "2.0"
	defaultTimeout          = 30 * time.Second
	defaultIdleTimeout      = 90 * time.Second
	defaultUserAgent        = "gozabbix"
	defaultMaxResponseBytes = 64 << 20 // 64 MiB
)

// Client is a Zabbix API client. A zero Client is not usable; construct one with
// New. Its methods are safe for concurrent use. Each request snapshots its
// credentials; changing credentials does not alter requests already in flight.
type Client struct {
	url            string
	userAgent      string
	httpClient     *http.Client
	ownedTransport *http.Transport
	maxBytes       int64
	logger         *slog.Logger
	id             atomic.Int64

	mu          sync.RWMutex
	credentials credentials
	generation  uint64
}

// New constructs a client without network I/O. The endpoint must have an HTTP
// or HTTPS scheme and a hostname, with no user information or fragment. Query
// parameters are allowed but are never included in default errors or logs.
func New(apiURL string, opts ...Option) (*Client, error) {
	endpoint, err := url.Parse(apiURL)
	if err != nil || endpoint.Hostname() == "" || endpoint.User != nil || strings.Contains(apiURL, "#") {
		return nil, errors.New("zabbix: invalid API endpoint")
	}
	endpoint.Scheme = strings.ToLower(endpoint.Scheme)
	if endpoint.Scheme != "http" && endpoint.Scheme != "https" {
		return nil, errors.New("zabbix: API endpoint must use HTTP or HTTPS")
	}
	cfg := clientConfig{userAgent: defaultUserAgent, maxBytes: defaultMaxResponseBytes}
	for _, opt := range opts {
		if opt == nil {
			return nil, errors.New("zabbix: option must not be nil")
		}
		opt(&cfg)
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	hc, owned, err := cfg.buildHTTPClient()
	if err != nil {
		return nil, err
	}
	logger := cfg.logger
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	return &Client{
		url: endpoint.String(), userAgent: cfg.userAgent, httpClient: hc,
		ownedTransport: owned, maxBytes: cfg.maxBytes, logger: logger,
		credentials: cfg.credentials,
	}, nil
}

// CloseIdleConnections releases idle connections in transports owned by this
// client. It does not interrupt active requests or close a borrowed transport,
// and the client remains usable. The owner of a borrowed transport handles its
// cleanup. Default and TLS-adjusted transports have a finite idle timeout.
func (c *Client) CloseIdleConnections() {
	if c.ownedTransport != nil {
		c.ownedTransport.CloseIdleConnections()
	}
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
	token := ""
	if authenticate {
		state, _ := c.authSnapshot()
		token = state.token
	}
	return c.doWithToken(ctx, method, params, out, token)
}

// doWithToken binds a request to an explicit credential snapshot. Login uses an
// empty token; Logout and compensating cleanup use the exact session they own.
func (c *Client) doWithToken(ctx context.Context, method string, params, out any, token string) error {
	if ctx == nil {
		return requestError(method, ErrorInvalidParams, errors.New("context must not be nil"))
	}
	destination := reflect.ValueOf(out)
	if out != nil && (destination.Kind() != reflect.Pointer || destination.IsNil()) {
		return requestError(method, ErrorInvalidParams, errors.New("output must be nil or a non-nil pointer"))
	}
	if params == nil {
		params = []any{}
	}
	reqBody := rpcRequest{Jsonrpc: jsonrpcVersion, Method: method, Params: params, ID: c.id.Add(1)}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return requestError(method, ErrorInvalidParams, err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, bytes.NewReader(body))
	if err != nil {
		return requestError(method, ErrorInvalidParams, err)
	}
	httpReq.Header.Set("Content-Type", "application/json-rpc")
	httpReq.Header.Set("Accept", "application/json")
	httpReq.Header.Set("User-Agent", c.userAgent)
	if token != "" {
		httpReq.Header.Set("Authorization", "Bearer "+token)
	}
	c.logger.DebugContext(ctx, "zabbix request", "method", method, "id", reqBody.ID)
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return requestError(method, ErrorTransport, err)
	}
	defer httpResp.Body.Close()
	if httpResp.StatusCode != http.StatusOK {
		return requestError(method, ErrorHTTP, &HTTPError{StatusCode: httpResp.StatusCode})
	}
	limited := &io.LimitedReader{R: httpResp.Body, N: c.maxBytes}
	raw, err := io.ReadAll(limited)
	if err != nil {
		return requestError(method, ErrorTransport, err)
	}
	if limited.N == 0 {
		// Probe separately so a valid truncated prefix is rejected without overflowing
		// maxBytes + 1 when the configured limit is the largest int64.
		var extra [1]byte
		n, err := io.ReadFull(httpResp.Body, extra[:])
		if n != 0 {
			return requestError(method, ErrorResponseLimit, fmt.Errorf("response exceeds %d byte limit", c.maxBytes))
		}
		if err != nil && !errors.Is(err, io.EOF) {
			return requestError(method, ErrorTransport, err)
		}
	}
	var resp rpcResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return requestError(method, ErrorProtocol, err)
	}
	if err := resp.validate(reqBody.ID); err != nil {
		kind := ErrorProtocol
		if _, ok := AsAPIError(err); ok {
			kind = ErrorAPI
		}
		return requestError(method, kind, err)
	}
	if out != nil {
		decoded := reflect.New(destination.Type().Elem())
		if err := json.Unmarshal(resp.Result, decoded.Interface()); err != nil {
			return requestError(method, ErrorDecode, err)
		}
		destination.Elem().Set(decoded.Elem())
	}
	return nil
}

func (resp rpcResponse) validate(requestID int64) error {
	if resp.Jsonrpc != jsonrpcVersion {
		return errors.New("zabbix: invalid JSON-RPC version")
	}
	if resp.ID == nil || *resp.ID != requestID {
		return errors.New("zabbix: invalid JSON-RPC response id")
	}
	if (len(resp.Result) == 0) == (len(resp.Error) == 0) {
		return errors.New("zabbix: response must contain exactly one result or error")
	}
	if len(resp.Error) != 0 {
		var apiError struct {
			Code    *int    `json:"code"`
			Message *string `json:"message"`
			Data    string  `json:"data"`
		}
		if err := json.Unmarshal(resp.Error, &apiError); err != nil || apiError.Code == nil || apiError.Message == nil {
			return errors.New("zabbix: invalid JSON-RPC error")
		}
		return &APIError{Code: *apiError.Code, Message: *apiError.Message, Data: apiError.Data}
	}
	if bytes.Equal(bytes.TrimSpace(resp.Result), []byte("null")) {
		return errors.New("zabbix: null JSON-RPC result")
	}
	return nil
}
