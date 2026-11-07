# gozabbix

A dependency-free Go client library for the [Zabbix JSON-RPC API](https://www.zabbix.com/documentation/current/en/manual/api), targeting **Zabbix 7.4**.

It authenticates with the HTTP `Authorization: Bearer` header (Zabbix 6.4+), which
works for both session tokens obtained via `user.login` and pre-created API tokens.

## Features

- Context-aware: every call takes a `context.Context` for timeouts and cancellation.
- Safe defaults: 30s request timeout, TLS 1.2 minimum, 64 MiB response cap.
- Configurable via functional options (custom `http.Client`, TLS, User-Agent, logger).
- Typed helpers for the most common namespaces plus a `Call` escape hatch for the rest.
- Thread-safe; a single `*Client` may be shared across goroutines.

## Install

```bash
go get github.com/gatework/go-zabbix
```

## Usage

### Session login

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gatework/go-zabbix"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	c, err := zabbix.New("https://zabbix.example.com/api_jsonrpc.php",
		zabbix.WithTimeout(15*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Login(ctx, "Admin", "zabbix"); err != nil {
		log.Fatal(err)
	}
	defer c.User().Logout(context.Background())

	hosts, err := c.Host().Get(ctx, zabbix.HostGetParams{
		GetParams: zabbix.GetParams{
			Output: []string{"hostid", "host", "name", "status"},
			Limit:  10,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, h := range hosts {
		fmt.Printf("%s\t%s\n", h.HostID, h.Name)
	}
}
```

### Pre-created API token

For service accounts, create an API token in the Zabbix UI and skip `Login`:

```go
c, _ := zabbix.New(url, zabbix.WithToken("your-api-token"))
version, _ := c.APIInfo().Version(ctx) // apiinfo.version is always unauthenticated
```

### Handling API errors

Transport failures and Zabbix API errors are distinguishable:

```go
_, err := c.Host().Get(ctx, zabbix.HostGetParams{})
if apiErr, ok := gozabbix.AsAPIError(err); ok {
	fmt.Println("zabbix rejected the request:", apiErr.Code, apiErr.Message)
}
```

### Calling unimplemented methods

Any method not yet wrapped is reachable via `Call`:

```go
var out []map[string]any
err := c.Call(ctx, "maintenance.get", map[string]any{"output": "extend"}, &out)
```

## Implemented namespaces

Typed helpers cover the full Zabbix 7.4 API surface (60+ namespaces):

- **Hosts & templates**: `host`, `hostgroup`, `hostinterface`, `hostprototype`,
  `template`, `templategroup`, `templatedashboard`, `item`, `itemprototype`,
  `trigger`, `triggerprototype`, `graph`, `graphprototype`, `graphitem`, `valuemap`
- **Monitoring & data**: `problem`, `event`, `history`, `trend`, `service`, `sla`,
  `dashboard`, `map`, `iconmap`, `image`
- **Discovery**: `drule`, `dcheck`, `dhost`, `dservice`, `discoveryrule`
- **Users & auth**: `user`, `usergroup`, `role`, `usermacro`, `userdirectory`,
  `token`, `mfa`, `authentication`, `settings`, `housekeeping`, `autoregistration`
- **Alerting & automation**: `action`, `alert`, `mediatype`, `script`, `correlation`,
  `maintenance`, `report`, `connector`
- **Infrastructure**: `proxy`, `proxygroup`, `hanode`, `module`, `regexp`, `task`,
  `auditlog`, `configuration` (import/export), `apiinfo`

Any method not wrapped by a typed helper is reachable through `Client.Call`.

## Options

| Option | Purpose |
| --- | --- |
| `WithTimeout(d)` | Per-request timeout (default 30s) |
| `WithHTTPClient(hc)` | Replace the underlying `*http.Client` |
| `WithToken(t)` | Authenticate with a pre-created API token |
| `WithTLSConfig(cfg)` | Custom TLS configuration |
| `WithInsecureSkipVerify(true)` | Accept self-signed certificates |
| `WithUserAgent(ua)` | Override the User-Agent header |
| `WithLogger(l)` | Debug-level request logging via `slog` |
| `WithMaxResponseBytes(n)` | Cap response body size (default 64 MiB) |

## Testing

Unit tests use a mocked HTTP server and require no network:

```bash
go test ./...
```

Integration tests run against a live Zabbix (see [docker-compose.yaml](docker-compose.yaml)):

```bash
docker compose up -d
go test -tags=integration -v \
  -zbx.url=http://localhost:8080/api_jsonrpc.php -zbx.user=Admin -zbx.pass=zabbix
```

## License

See [LICENSE](LICENSE).
