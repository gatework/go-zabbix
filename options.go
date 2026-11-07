package zabbix

import (
	"crypto/tls"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"reflect"
	"time"
)

// Option configures a Client before any resources are created. Options for
// different settings are independent of order; the last value for a setting
// wins. WithAPIToken and WithSessionToken configure the same credential setting.
type Option func(*clientConfig)

type clientConfig struct {
	httpClient  *http.Client
	httpSet     bool
	timeout     time.Duration
	timeoutSet  bool
	tlsConfig   *tls.Config
	tlsSet      bool
	skipVerify  bool
	skipSet     bool
	userAgent   string
	maxBytes    int64
	logger      *slog.Logger
	credentials credentials
}

// WithHTTPClient borrows hc's transport and jar without modifying hc. Its timeout
// is retained unless WithTimeout is supplied. The SDK refuses redirects.
// TLS options require a *http.Transport and create a private transport copy.
func WithHTTPClient(hc *http.Client) Option {
	return func(cfg *clientConfig) { cfg.httpClient, cfg.httpSet = hc, true }
}

// WithTimeout overrides the HTTP request timeout. Zero disables the timeout;
// callers should then provide a context deadline. Negative values are invalid.
func WithTimeout(d time.Duration) Option {
	return func(cfg *clientConfig) { cfg.timeout, cfg.timeoutSet = d, true }
}

// WithAPIToken supplies a pre-created API token. It is not verified by New.
func WithAPIToken(token string) Option {
	return func(cfg *clientConfig) { cfg.credentials = credentials{kind: credentialAPI, token: token} }
}

// WithSessionToken supplies a session token obtained from user.login.
// It is not verified by New and can be invalidated with Logout.
func WithSessionToken(token string) Option {
	return func(cfg *clientConfig) { cfg.credentials = credentials{kind: credentialSession, token: token} }
}

// WithUserAgent overrides the User-Agent header. Empty disables that header.
func WithUserAgent(ua string) Option {
	return func(cfg *clientConfig) { cfg.userAgent = ua }
}

// WithLogger installs a logger for request method and ID at debug level.
// A nil logger discards logs. URLs, credentials and payloads are never logged.
func WithLogger(logger *slog.Logger) Option {
	return func(cfg *clientConfig) { cfg.logger = logger }
}

// WithMaxResponseBytes sets a positive response-body limit. The default is 64 MiB.
func WithMaxResponseBytes(n int64) Option {
	return func(cfg *clientConfig) { cfg.maxBytes = n }
}

// WithTLSConfig uses a copy of cfg. A nil configuration is invalid. For an
// injected HTTP client, its effective transport must be an *http.Transport.
func WithTLSConfig(cfg *tls.Config) Option {
	return func(options *clientConfig) { options.tlsConfig, options.tlsSet = cfg, true }
}

// WithInsecureSkipVerify overrides certificate verification in the private TLS
// configuration, regardless of option order. Disabling verification is unsafe
// unless the connection is otherwise authenticated.
func WithInsecureSkipVerify(skip bool) Option {
	return func(cfg *clientConfig) { cfg.skipVerify, cfg.skipSet = skip, true }
}

func (cfg clientConfig) validate() error {
	switch {
	case cfg.httpSet && cfg.httpClient == nil:
		return errors.New("zabbix: HTTP client must not be nil")
	case cfg.httpClient != nil && isNilTransport(cfg.httpClient.Transport):
		return errors.New("zabbix: HTTP transport must not contain a typed nil")
	case cfg.timeoutSet && cfg.timeout < 0:
		return errors.New("zabbix: timeout must not be negative")
	case cfg.httpClient != nil && !cfg.timeoutSet && cfg.httpClient.Timeout < 0:
		return errors.New("zabbix: HTTP client timeout must not be negative")
	case cfg.maxBytes <= 0:
		return errors.New("zabbix: response limit must be positive")
	case cfg.tlsSet && cfg.tlsConfig == nil:
		return errors.New("zabbix: TLS configuration must not be nil")
	case !validHeaderValue(cfg.userAgent):
		return errors.New("zabbix: invalid user agent")
	}
	if cfg.credentials.kind != credentialNone {
		return validateToken(cfg.credentials.token)
	}
	return nil
}

// A nil interface selects http.DefaultTransport; a typed nil would instead be
// invoked by http.Client. Only nil-capable concrete kinds support IsNil.
func isNilTransport(transport http.RoundTripper) bool {
	if transport == nil {
		return false
	}
	value := reflect.ValueOf(transport)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}

func (cfg clientConfig) buildHTTPClient() (*http.Client, *http.Transport, error) {
	hc := &http.Client{Timeout: defaultTimeout}
	var owned *http.Transport
	if cfg.httpClient != nil {
		*hc = *cfg.httpClient
	} else {
		owned = newTransport()
		hc.Transport = owned
	}
	if cfg.timeoutSet {
		hc.Timeout = cfg.timeout
	}
	if cfg.tlsSet || cfg.skipSet {
		transport := hc.Transport
		if transport == nil {
			transport = http.DefaultTransport
		}
		tr, ok := transport.(*http.Transport)
		if !ok || tr == nil {
			return nil, nil, errors.New("zabbix: TLS options require an HTTP transport")
		}
		if owned == nil {
			owned = tr.Clone()
		}
		if cfg.tlsSet {
			owned.TLSClientConfig = cfg.tlsConfig.Clone()
		}
		if owned.TLSClientConfig == nil {
			owned.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		}
		if cfg.skipSet {
			owned.TLSClientConfig.InsecureSkipVerify = cfg.skipVerify
		}
		if owned.IdleConnTimeout <= 0 {
			owned.IdleConnTimeout = defaultIdleTimeout
		}
		hc.Transport = owned
	}
	hc.CheckRedirect = func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }
	return hc, owned, nil
}

func newTransport() *http.Transport {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       defaultIdleTimeout,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
		TLSClientConfig:       &tls.Config{MinVersion: tls.VersionTLS12},
	}
}

func validHeaderValue(value string) bool {
	for i := 0; i < len(value); i++ {
		if value[i] == 0x7f || value[i] < 0x20 && value[i] != '\t' {
			return false
		}
	}
	return true
}
