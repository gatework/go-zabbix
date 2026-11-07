# Working on go-zabbix

The module is `github.com/gatework/go-zabbix`; its package is `zabbix`. It targets
Zabbix 7.4 and the minimum Go version in `go.mod`, with standard-library
dependencies only. Keep namespace navigation in this single package.

## Protocol changes

Before changing a namespace, read its row in
[the 7.4 contract matrix](docs/api-7.4-contracts.md) and the linked versioned
Zabbix reference. Keep the matrix and a request/response regression in sync.
The `current` documentation may describe a different Zabbix release.

- Query parameters expose only that method's capabilities. `Get` returns records;
  `Count` and `GetMap` handle the alternate result shapes where supported.
- Every exposed selector needs a corresponding response field. Count-capable
  associations use `Relation[T]`; distinguish an omitted association, an empty
  list, and a zero count.
- Keep reading models separate from create/update inputs. An omitted update
  field leaves server state unchanged; an explicit empty collection replaces
  the association with an empty collection. Use scalar pointers and collection
  `omitzero` tags for these distinct meanings.
- Preserve IDs as decimal strings. Other values follow the documented wire
  type; Zabbix does not return every field as a string.
- Keep raw `Call` available for protocol capabilities outside the typed surface.

## Client ownership and failures

- `New` performs no network I/O. Options configure a private construction value;
  borrowed clients and transports must remain unchanged.
- A shared caller transport keeps its pool. Cleanup closes only SDK-owned idle
  connections. API tokens and login sessions have different lifecycles.
- Keep response limits, JSON-RPC envelope validation and atomic output decoding
  in the common request path. A failed decode must not partially replace output.
- Preserve error chains for `errors.Is/As`. Default error strings and debug logs
  contain method/status metadata, not credentials, URL queries or payloads.

## Verification

Run the build, vet, lint and race checks in `.github/workflows/ci.yml`.
Tests must also compile and pass on the minimum Go version. HTTP tests use
`httptest`; blocking handlers need cleanup that releases the handler before
closing the test server.

Use `bash scripts/integration.sh -v` for real protocol changes. It creates and
removes an isolated Zabbix 7.4 stack. Integration tests may create, update and
delete objects, so point them only at a disposable fixture database. Assertions
must check returned fields and persisted changes, rather than only a nil error.

For public interface changes, update [the migration guide](MIGRATION.md), README
examples and compile-checked examples. Domain names are defined in
[the glossary](GLOSSARY.md).
