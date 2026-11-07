# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`zabbix` is a Go client library for the [Zabbix JSON-RPC API](https://www.zabbix.com/documentation/current/en/manual/api), targeting **Zabbix 7.4**. It is imported as `github.com/gatework/go-zabbix` and has **no third-party dependencies** — only the Go standard library (requires the Go toolchain pinned in [go.mod](go.mod)).

## Commands

```bash
go build ./...
go vet ./...
gofmt -l .                      # must print nothing (CI fails otherwise)

# Unit tests (mocked HTTP, no network)
go test ./...
go test -race -count=1 ./...     # how CI runs them
go test -run TestLoginStoresTokenAndSendsBearer -v

# Integration tests (live Zabbix; behind the `integration` build tag)
docker compose up -d
go test -tags=integration -v \
  -zbx.url=http://localhost:8080/api_jsonrpc.php -zbx.user=Admin -zbx.pass=zabbix
```

Default credentials for the compose stack are `Admin` / `zabbix`; the web UI + API is at `http://localhost:8080/api_jsonrpc.php`.

## Architecture

The library is a thin, uniform wrapper over the JSON-RPC API. Core plumbing lives in [client.go](client.go); each Zabbix API namespace gets its own file.

- **`Client`** ([client.go](client.go)) holds the endpoint URL, an `*http.Client`, an atomic request-id counter, and a mutex-guarded auth token. Construct with `New(url, ...Option)`, which performs **no network I/O** — authenticate afterwards with `Login`, or pass a pre-created API token via `WithToken`. The token is used as an `Authorization: Bearer` header (Zabbix 6.4+), which is why both session tokens and API tokens work through one path.

- **`(*Client).do(ctx, method, params, out, authenticate)`** ([client.go](client.go)) is the single choke point for every request. It builds the JSON-RPC envelope, attaches the bearer header when `authenticate` is true and a token is set, POSTs with a context, reads the body under an `io.LimitReader` cap, and decodes the `result` into `out`. A non-nil `error` object in the response becomes an `*APIError` (distinguishable from transport errors via `AsAPIError`). The public `Call` method is `do` with auth enabled — it is the **escape hatch** for any namespace not yet typed.

- **`Result` is `json.RawMessage`**, not `interface{}`. Each namespace method passes a concrete `out` pointer to `do`, so there is no generic re-marshal step. The one polymorphic case is `user.login` (returns a bare string, or an object when `userData` is set), handled explicitly in [user.go](user.go).

- **Namespace pattern** — each API family is a `XService` struct wrapping `*Client`, exposed via an accessor (`client.Host()` → `*HostService`). Methods build typed params and call `client.Call`. **When adding a namespace**, copy an existing file (e.g. [host.go](host.go)): define the object struct(s) with Zabbix doc URLs in comments, a `XGetParams` struct embedding `GetParams` from [common.go](common.go), and `Get/Create/Update/Delete` methods. `Get` returns `[]Object`; `Create/Update/Delete` return `[]string` of affected ids decoded from the `{"<obj>ids": [...]}` result.

## Conventions

- One file per namespace (`<namespace>.go`). The full Zabbix 7.4 surface is implemented (60+ namespaces); every file follows the pattern above. A few namespaces deviate for real API reasons — e.g. `task.get` does NOT embed `GetParams` because the server rejects `limit`/`sort`/`filter` there, and singletons (`settings`, `authentication`, `housekeeping`, `autoregistration`) expose `Get(ctx) (*T, error)` + `Update`. When in doubt, validate against the live stack via the integration tests before assuming a param is accepted.
- Zabbix returns all object fields as **strings** (even numeric ids/flags) — keep struct fields `string` with `omitempty` tags mirroring the object reference, and link the relevant Zabbix doc URL above each type/method.
- Query params embed the shared `GetParams` ([common.go](common.go)); `select*` and `output` fields are `any` because Zabbix accepts either `"extend"` or a `[]string`.
- Delete/version/logout params must serialize to a present JSON array — `do` substitutes `[]any{}` for nil params so the `params` field is never omitted.
- Tests mock the server with `httptest`; a blocking-handler test must release its handler on cleanup (see `TestContextCancellation`) so `Server.Close` never hangs.
