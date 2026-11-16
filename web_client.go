package zabbix

import (
	"context"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	// ErrWebInput indicates an invalid frontend endpoint, session or chart request.
	ErrWebInput = errors.New("zabbix: invalid web request")
	// ErrWebAuthentication indicates that a frontend login or session was rejected.
	ErrWebAuthentication = errors.New("zabbix: web authentication required")
	// ErrWebInteractiveLogin indicates that SSO or MFA needs an external login flow.
	ErrWebInteractiveLogin = errors.New("zabbix: interactive web login required")
	// ErrWebImage indicates that the response is not a usable chart PNG.
	ErrWebImage = errors.New("zabbix: invalid web chart image")
)

// WebOptions configures web access independently from API bearer credentials.
type WebOptions struct {
	// HTTPClient supplies a borrowed transport. Its cookies for the frontend are
	// copied once into a private jar; subsequent session changes are independent.
	HTTPClient *http.Client
	// Timezone must match the frontend user's timezone; nil means UTC.
	Timezone *time.Location
	// MaxImageBytes defaults to 8 MiB. It must not exceed 64 MiB.
	MaxImageBytes int64
}

// WebClient accesses the Zabbix 7.4 web chart endpoints using a signed cookie.
// Construct it with NewWebClient; zero-value operations return errors.
// A WebClient is safe for concurrent use and never retains login passwords.
type WebClient struct {
	base           *url.URL
	location       *time.Location
	maxBytes       int64
	ownedTransport *http.Transport
	loginGate      chan struct{}

	mu         sync.RWMutex // protects the current client/jar and session generation
	client     *http.Client
	generation uint64
}

// NewWebClient constructs a client without network I/O. baseURL is the frontend
// directory, for example https://monitor.example/zabbix/, not api_jsonrpc.php.
// The borrowed client's timeout is preserved when positive; otherwise it is 30s.
func NewWebClient(baseURL string, options WebOptions) (*WebClient, error) {
	base, err := url.Parse(baseURL)
	if err != nil || base.Hostname() == "" || base.User != nil || base.RawQuery != "" || base.ForceQuery ||
		strings.Contains(baseURL, "#") || (base.Scheme != "https" && base.Scheme != "http") {
		return nil, ErrWebInput
	}
	if strings.HasSuffix(base.Path, ".php") || strings.Contains(base.Path, "\\") {
		return nil, ErrWebInput
	}
	base.Path = strings.TrimRight(base.Path, "/") + "/"
	base.RawPath = ""
	if options.MaxImageBytes == 0 {
		options.MaxImageBytes = 8 << 20
	}
	if options.MaxImageBytes < 1 || options.MaxImageBytes > 64<<20 {
		return nil, ErrWebInput
	}
	cfg := clientConfig{httpClient: options.HTTPClient, maxBytes: options.MaxImageBytes}
	if err := cfg.validate(); err != nil {
		return nil, ErrWebInput
	}
	hc, owned, err := cfg.buildHTTPClient()
	if err != nil {
		return nil, ErrWebInput
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, ErrWebInput
	}
	if hc.Jar != nil {
		jar.SetCookies(base, hc.Jar.Cookies(base))
	}
	hc.Jar = jar
	if hc.Timeout <= 0 {
		hc.Timeout = defaultTimeout
	}
	location := options.Timezone
	if location == nil {
		location = time.UTC
	}
	return &WebClient{
		base: base, location: location, maxBytes: options.MaxImageBytes,
		ownedTransport: owned, client: hc, loginGate: make(chan struct{}, 1),
	}, nil
}

// CloseIdleConnections closes only connections owned by this WebClient.
func (w *WebClient) CloseIdleConnections() {
	if w != nil && w.ownedTransport != nil {
		w.ownedTransport.CloseIdleConnections()
	}
}

// ClearSession forgets local cookies. It does not revoke a session on the server;
// use Logout when revocation is required. An in-flight login cannot restore it.
func (w *WebClient) ClearSession() error {
	if w == nil || w.base == nil {
		return ErrWebInput
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	next, err := freshWebClient(w.client)
	if err != nil {
		return err
	}
	w.client = next
	w.generation++
	return nil
}

func (w *WebClient) snapshot(ctx context.Context) (*http.Client, uint64, error) {
	if ctx == nil || w == nil || w.base == nil {
		return nil, 0, ErrWebInput
	}
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.client, w.generation, nil
}

func (w *WebClient) endpoint(name string) *url.URL {
	return w.base.ResolveReference(&url.URL{Path: name})
}

func freshWebClient(previous *http.Client) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, ErrWebInput
	}
	next := *previous
	next.Jar = jar
	return &next, nil
}
