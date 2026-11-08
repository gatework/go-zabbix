# go-zabbix

A standard-library-only Go client for the [Zabbix 7.4 JSON-RPC API](https://www.zabbix.com/documentation/7.4/en/manual/api).
Requires **Go 1.24 or later**. The module is `github.com/gatework/go-zabbix` and
the imported package is `zabbix`.

Version `v0.1.0` changes public contracts from earlier untagged snapshots.
Existing callers should read [MIGRATION.md](MIGRATION.md) before updating.
The [namespace matrix](docs/api-7.4-contracts.md)
records the supported parameters, selectors and result shapes; the
[mutation matrix](docs/api-7.4-mutations.md) records write contracts. Coverage concerns
the typed methods exposed by this library; it does not imply every Zabbix method
or field is wrapped. `Client.Call` handles additional protocol operations.

## Install

```sh
go get github.com/gatework/go-zabbix@v0.1.0
```

Keep the tagged version pinned in application `go.mod` files. Read
[MIGRATION.md](MIGRATION.md) when upgrading from an earlier snapshot.

## Read hosts with an API token

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/gatework/go-zabbix"
)

func main() {
    client, err := zabbix.New(os.Getenv("ZABBIX_URL"),
        zabbix.WithAPIToken(os.Getenv("ZABBIX_API_TOKEN")))
    if err != nil {
        log.Fatal(err)
    }
    defer client.CloseIdleConnections()

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()
    hosts, err := client.Host().Get(ctx, zabbix.HostGetParams{
        Output: []string{"hostid", "name"},
        Limit: 10,
        SelectInterfaces: "count",
    })
    if err != nil {
        log.Fatal(err)
    }
    for _, host := range hosts {
        if host.Interfaces.Count != nil {
            fmt.Println(host.HostID, host.Name, *host.Interfaces.Count)
        }
    }
}
```

Use `Host().Count(ctx, HostCountParams{...})` for a count, or
`Host().GetMap(ctx, HostGetParams{...})` for an ID-keyed map. `Get` always returns
a slice. A count-capable association such as `Interfaces` uses `Relation[T]`:
`Items` contains records and `Count` contains a selected count, including zero.

## Sessions and connection ownership

For username/password authentication, construct a client without a token and
call `client.Login(ctx, username, password)`. Call `client.Logout` with a fresh,
bounded context to revoke the session when finished. `LoginWithUserData` returns
session details. API tokens use `WithAPIToken` or `SetAPIToken`; session IDs use
`WithSessionToken` or `SetSessionToken`. Logout does not revoke an API token.

A client is safe to share across goroutines. A request snapshots its credentials;
a late login or logout cannot overwrite newer authentication state. Local token
setters and `ClearCredentials` do not revoke server sessions. Log out a live
session before replacing it when revocation is required.

Defaults are a 30-second request timeout, TLS 1.2 minimum, a 64 MiB response cap,
and a finite idle connection timeout. `WithHTTPClient` copies a supplied client
and borrows its transport. TLS options clone compatible transports; incompatible
custom transports cause a constructor error. `CloseIdleConnections` closes only
SDK-owned pools. Option order is independent across different settings; the last
option for the same setting wins.

| Option | Purpose |
| --- | --- |
| `WithTimeout(d)` | Override the request timeout; zero disables it |
| `WithHTTPClient(client)` | Borrow caller-managed transport and connection pool |
| `WithAPIToken(token)` / `WithSessionToken(token)` | Select credential kind |
| `WithTLSConfig(config)` | Clone a custom TLS configuration |
| `WithInsecureSkipVerify(true)` | Explicitly disable certificate verification |
| `WithUserAgent(value)` | Set the request User-Agent |
| `WithLogger(logger)` | Log method and request ID at debug level |
| `WithMaxResponseBytes(n)` | Set a positive response-size limit |

## Write only the intended changes

Read models and write parameters are separate. Optional scalar pointers preserve
zero values. Nil collections omit a field; allocated empty collections clear an
association where the operation allows it:

```go
description := ""
_, err := client.Host().Update(ctx, zabbix.HostUpdateParams{
    HostID: "10084",
    Description: &description,
    Tags: []zabbix.TagInput{},
})
```

Nested inputs and ID references have their own types. The server still enforces
permissions and state-dependent restrictions; client validation checks required
identity fields and operation shape.

Ordinary items support HTTP agent configuration and preprocessing. Ordered
`QueryFields` and `Headers` preserve repeated names; `SelectPreprocessing` reads
the configured steps. See the [migration guide](MIGRATION.md#role-service-access-and-http-items)
for partial updates, clearing lists and role service-access tags.

Typed mutations check the result against the request, including ID multiplicity
and operation-specific batch behavior. `Autoregistration.Update` returns only an
error because its server receipt is boolean. A successful receipt does not replace
a read-back when the application needs to verify persisted state.

## Errors and raw calls

`RequestError` records the method and failure category while preserving its cause
for `errors.Is` and `errors.As`. `HTTPError` carries the status; `APIError` carries
the Zabbix code, message and data. Default error strings omit response bodies,
endpoint URLs and upstream error details. Inspect `APIError.Message` and `Data`
only where exposing upstream details is appropriate.

```go
var apiErr *zabbix.APIError
if errors.As(err, &apiErr) {
    // Handle apiErr.Code according to the operation.
}

var result []map[string]any
err = client.Call(ctx, "maintenance.get", map[string]any{"output": "extend"}, &result)
```

`Call` retains the common timeout, response limit, envelope checks and atomic
output decoding. It does not apply typed namespace validation.

## Verification

```sh
go test -race -shuffle=on ./...
go vet ./...
go build ./...
bash scripts/integration.sh -v
```

Unit tests use local HTTP fixtures. The integration script requires Docker and
creates a separate Zabbix 7.4/PostgreSQL stack with random loopback ports. It
prints the server version and image digests, runs read/write round trips, then
removes its containers and volumes. Use only disposable databases for integration
tests. CI runs the minimum Go version, stable Go, lint and the Docker suite.

## License

See [LICENSE](LICENSE).
