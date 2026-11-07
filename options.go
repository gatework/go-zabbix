package zabbix

import (
	"crypto/tls"
	"log/slog"
	"net/http"
	"time"
)

// Option configures a Client at construction time. Options are applied in order,
// so WithHTTPClient should come before options that tune the default transport
// (WithTimeout, WithTLSConfig, WithInsecureSkipVerify) if you want those applied
// to your custom client.
type Option func(*Client)

// WithHTTPClient replaces the underlying *http.Client entirely. When set, you are
// responsible for its Timeout and TLS configuration.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		if hc != nil {
			c.httpClient = hc
		}
	}
}

// WithTimeout sets the per-request timeout on the client's http.Client.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.httpClient.Timeout = d }
}

// WithToken sets a pre-created API token (or an existing session id), so the
// client is authenticated without calling Login.
func WithToken(token string) Option {
	return func(c *Client) { c.SetToken(token) }
}

// WithUserAgent overrides the User-Agent header sent with every request.
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		if ua != "" {
			c.userAgent = ua
		}
	}
}

// WithLogger installs an slog.Logger; requests are logged at debug level. By
// default logs are discarded.
func WithLogger(l *slog.Logger) Option {
	return func(c *Client) {
		if l != nil {
			c.logger = l
		}
	}
}

// WithMaxResponseBytes caps how many bytes are read from a response body,
// guarding against unbounded memory use. The default is 64 MiB.
func WithMaxResponseBytes(n int64) Option {
	return func(c *Client) {
		if n > 0 {
			c.maxBytes = n
		}
	}
}

// WithTLSConfig sets the TLS configuration on the default transport. It is a
// no-op if a custom http.Client with a non-*http.Transport is in use.
func WithTLSConfig(cfg *tls.Config) Option {
	return func(c *Client) {
		if cfg != nil {
			defaultTransport(c).TLSClientConfig = cfg
		}
	}
}

// WithInsecureSkipVerify disables TLS certificate verification on the default
// transport. Use only for self-signed certificates in trusted environments.
func WithInsecureSkipVerify(skip bool) Option {
	return func(c *Client) {
		tr := defaultTransport(c)
		if tr.TLSClientConfig == nil {
			tr.TLSClientConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		}
		tr.TLSClientConfig.InsecureSkipVerify = skip
	}
}

// defaultTransport returns the client's *http.Transport, installing a fresh one
// if the current transport is not a *http.Transport.
func defaultTransport(c *Client) *http.Transport {
	if tr, ok := c.httpClient.Transport.(*http.Transport); ok {
		return tr
	}
	tr := &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}}
	c.httpClient.Transport = tr
	return tr
}
