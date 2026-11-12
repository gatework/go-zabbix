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
	// ErrFrontendInput indicates an invalid frontend endpoint, session or chart request.
	ErrFrontendInput = errors.New("zabbix: invalid frontend request")
	// ErrFrontendAuthentication indicates that a frontend login or session was rejected.
	ErrFrontendAuthentication = errors.New("zabbix: frontend authentication required")
	// ErrFrontendInteractiveLogin indicates that SSO or MFA needs an external login flow.
	ErrFrontendInteractiveLogin = errors.New("zabbix: interactive frontend login required")
	// ErrFrontendImage indicates that the response is not a usable chart PNG.
	ErrFrontendImage = errors.New("zabbix: invalid frontend chart image")
)

// FrontendOptions configures web access independently from API bearer credentials.
type FrontendOptions struct {
	// HTTPClient supplies a borrowed transport. Its cookies for the frontend are
	// copied once into a private jar; subsequent session changes are independent.
	HTTPClient *http.Client
	// Timezone must match the frontend user's timezone; nil means UTC.
	Timezone *time.Location
	// MaxImageBytes defaults to 8 MiB. It must not exceed 64 MiB.
	MaxImageBytes int64
}

// Frontend accesses the Zabbix 7.4 web chart endpoints using a signed cookie.
// Construct it with NewFrontend; zero-value operations return errors.
// A Frontend is safe for concurrent use and never retains login passwords.
type Frontend struct {
	base           *url.URL
	location       *time.Location
	maxBytes       int64
	ownedTransport *http.Transport
	loginGate      chan struct{}

	mu         sync.RWMutex // protects the current client/jar and session generation
	client     *http.Client
	generation uint64
}

// NewFrontend constructs a client without network I/O. baseURL is the frontend
// directory, for example https://monitor.example/zabbix/, not api_jsonrpc.php.
// The borrowed client's timeout is preserved when positive; otherwise it is 30s.
func NewFrontend(baseURL string, options FrontendOptions) (*Frontend, error) {
	base, err := url.Parse(baseURL)
	if err != nil || base.Hostname() == "" || base.User != nil || base.RawQuery != "" || base.ForceQuery ||
		strings.Contains(baseURL, "#") || (base.Scheme != "https" && base.Scheme != "http") {
		return nil, ErrFrontendInput
	}
	if strings.HasSuffix(base.Path, ".php") || strings.Contains(base.Path, "\\") {
		return nil, ErrFrontendInput
	}
	base.Path = strings.TrimRight(base.Path, "/") + "/"
	base.RawPath = ""
	if options.MaxImageBytes == 0 {
		options.MaxImageBytes = 8 << 20
	}
	if options.MaxImageBytes < 1 || options.MaxImageBytes > 64<<20 {
		return nil, ErrFrontendInput
	}
	cfg := clientConfig{httpClient: options.HTTPClient, maxBytes: options.MaxImageBytes}
	if err := cfg.validate(); err != nil {
		return nil, ErrFrontendInput
	}
	hc, owned, err := cfg.buildHTTPClient()
	if err != nil {
		return nil, ErrFrontendInput
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, ErrFrontendInput
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
	return &Frontend{
		base: base, location: location, maxBytes: options.MaxImageBytes,
		ownedTransport: owned, client: hc, loginGate: make(chan struct{}, 1),
	}, nil
}

// CloseIdleConnections closes only connections owned by this Frontend.
func (f *Frontend) CloseIdleConnections() {
	if f != nil && f.ownedTransport != nil {
		f.ownedTransport.CloseIdleConnections()
	}
}

// ClearSession forgets local cookies. It does not revoke a session on the server;
// use Logout when revocation is required. An in-flight login cannot restore it.
func (f *Frontend) ClearSession() error {
	if f == nil || f.base == nil {
		return ErrFrontendInput
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	next, err := freshFrontendClient(f.client)
	if err != nil {
		return err
	}
	f.client = next
	f.generation++
	return nil
}

func (f *Frontend) snapshot(ctx context.Context) (*http.Client, uint64, error) {
	if ctx == nil || f == nil || f.base == nil {
		return nil, 0, ErrFrontendInput
	}
	if err := ctx.Err(); err != nil {
		return nil, 0, err
	}
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.client, f.generation, nil
}

func (f *Frontend) endpoint(name string) *url.URL {
	return f.base.ResolveReference(&url.URL{Path: name})
}

func freshFrontendClient(previous *http.Client) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, ErrFrontendInput
	}
	next := *previous
	next.Jar = jar
	return &next, nil
}
