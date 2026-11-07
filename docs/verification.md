# 7.4 contract verification

Local verification on 2026-11-07 used macOS arm64, Go 1.24.13 and Go 1.27.1.
The module has no third-party Go dependencies. GitHub Actions repeats the minimum
and stable Go checks and runs the disposable Docker suite on Linux; its result
must be read for the published commit separately from these local results.

| Check | Local result |
| --- | --- |
| Go 1.24.13 `go test -race -shuffle=on -count=1 ./...` | Passed |
| Go 1.24.13 integration-tag compilation | Passed |
| `go vet ./...` and `go build ./...` | Passed |
| golangci-lint v2.12.2 | 0 issues |
| `bash scripts/integration.sh -v` | Passed, including unit tests and live regressions with the race detector |
| Temporary stack cleanup | Containers, network and volumes removed by the driver |

## Live fixture

The disposable stack reported Zabbix API **7.4.16** and these image digests:

- PostgreSQL: `sha256:3540c0362c809df084ffac0ae69fe794515f9de92692f02178bcdce0ef0302ac`
- Zabbix server: `sha256:1839765578f8db2569518c53237ea3d03c6e0f23ab7597a3763e9683f99b89f1`
- Zabbix web: `sha256:4e2ee7da48a3fcec5524bda190889d4d0c928cf1f0b78822fded5cdff4b1e87c`

Compose targets the 7.4 image series. The driver prints the actual version and
digests on each run; later image updates must not be mistaken for this fixture.
Ports bind to loopback and are selected independently for each stack.

## Evidence covered

The live suite exercises the exposed collection/singleton reads, all 51 Count
and 50 GetMap capabilities and both global macro variants. HTTP fixtures cover
all 120 exposed selector mappings, including populated relation objects, counts,
empty selections, malformed shapes and atomic decoding. Separate live probes
also queried every exposed selector against the server. Empty live reads alone
do not prove every populated field or every server state.

Live mutation regressions verify host zero/empty/omitted updates; standalone
host-macro ownership; template groups and both user-group permission collections;
interface mass-add, no-match removal and empty replacement; discovery-rule and
item-prototype deletion keys; event severity zero and indefinite suppression;
service/host partial updates; API-token generation and use; SLI empty selection
and epoch zero; raw configuration export; and populated host/event URL-script
lookups with macro expansion. The script lookup tests do not execute scripts.
Additional CRUD smoke checks cover a representative set of writable namespaces.
All fixture-created objects are registered for bounded cleanup.

The second-round regressions additionally verify user-data login with default
GUI access and debug-enabled groups; role service-access tag clearing while
preserving lists, modes and API rules; all HTTP item fields and preprocessing,
partial updates and empty-list clearing; the autoregistration boolean receipt;
and populated numeric-ID MFA deletion and interface mass removal. The HTTP item
is disabled: this proves configuration persistence, not an external HTTP poll.
Autoregistration submits its existing TLS mode to exercise the receipt without
creating or erasing a write-only PSK.

HTTP tests separately cover malformed token-generation/script-execution results,
large numeric IDs, bounded bodies, TLS/client ownership, cancellation and
concurrent credential replacement. The integration suite does not claim live
coverage of every mutation, enum combination, inherited/provisioned restriction,
or external Zabbix deployment.

Mutation HTTP fixtures also reject count, identity, multiplicity and
singleton-field mismatches while preserving the server's Task, acknowledgement
and bulk-interface semantics. Login cleanup regressions prove bounded
compensation after user-data decode failure. CI compiles integration-tag tests
under both minimum and stable Go versions as well as running the live fixture
on stable Go.

See [query contracts](api-7.4-contracts.md),
[mutation contracts](api-7.4-mutations.md) and [migration](../MIGRATION.md).
