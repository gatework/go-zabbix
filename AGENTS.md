# Repository Guidelines

## Project Structure & Module Organization

This standard-library-only Go library targets Zabbix 7.4 and requires Go 1.24+.
The module is `github.com/gatework/go-zabbix`; keep implementations in the root
`zabbix` package. Namespace files such as `host.go` and `item.go` expose API
operations; `client.go`, `auth.go`, and `options.go` handle shared client behavior.
Tests sit beside source as `*_test.go`; `example_test.go` contains compile-checked
usage. `docs/` holds protocol matrices and verification notes.

## Build, Test, and Development Commands

Run these from the repository root:

- `go build ./...` — compile the library.
- `go vet ./...` — check for suspicious code.
- `go test -race -shuffle=on -count=1 ./...` — run uncached unit tests with race detection.
- `go test -tags=integration -run '^$' ./...` — compile integration tests without running them.
- `golangci-lint run ./...` — run lint checks; use the version pinned in `.github/workflows/ci.yml`.
- `bash scripts/integration.sh -v` — create a disposable Zabbix/PostgreSQL stack using `docker-compose.yaml`, run live tests, and remove it. Requires Docker Compose.

CI checks both Go 1.24 and stable Go.

## Coding Style & Naming Conventions

Use `gofmt` for tab indentation and `goimports` for imports; both are configured
in `.golangci.yml`. Use MixedCaps identifiers and uppercase initialisms such as
`HostID`, `HTTPClient`, and `APIToken`. Follow existing namespace and parameter
names, such as `HostGetParams` and `HostUpdateParams`. Document exported APIs.

## Testing Guidelines

Use standard-library `testing` and `net/http/httptest`. Name tests
`Test<Behavior>` and live tests `TestIntegration<Behavior>` behind the
`integration` build tag. Add focused request/response regressions for behavior
changes. Assert decoded fields and persisted changes, including omitted, empty,
and zero values. CI defines no numeric coverage threshold. Run live mutations
only against disposable fixtures.

## Protocol & Configuration Rules

Read [CLAUDE.md](CLAUDE.md) for protocol and ownership rules and
[GLOSSARY.md](GLOSSARY.md) for domain terminology. Before namespace changes,
consult the relevant matrix in `docs/` and its versioned Zabbix 7.4 reference;
update the matrix alongside regressions. Preserve decimal-string IDs and separate
read models from write inputs. Keep credentials and payloads out of logs.

## Commit & Pull Request Guidelines

Recent history favors concise imperative subjects with `fix:`, `test:`, `docs:`,
or `refactor:` prefixes. Keep commits focused. PRs should explain the behavior
change, link relevant issues, and report checks performed. Public API changes
must update `MIGRATION.md`, README examples, and compile-checked examples.
