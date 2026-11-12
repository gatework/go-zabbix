# v0.3.0 monitoring extraction

`Client.Monitoring()` now owns the reusable operations previously implemented
in go-admin: host/item lookup, history/trend query limits and traffic processing.
No go-admin model, database or authorization dependency is introduced.

| Previous go-admin responsibility/type | SDK replacement |
| --- | --- |
| Provider host/item search and lookup | `Monitoring().Hosts`, `Items`, `Host`, `Item` |
| `ZabbixHost`, `ZabbixHostList` | `MonitoredHost`, `MonitoredHostList` |
| `ZabbixItem`, `ZabbixItemList` | `MonitoredItem`, `MonitoredItemList` |
| `ZabbixSample`, `ZabbixTrend` | `MetricSample`, `MetricTrend` |
| Provider history/trends | `Monitoring().History`, `Trends` |
| Service rate/peak/gap helpers | `Monitoring().Traffic` or `ReadTraffic` |

Use SDK types directly; retain application DTO conversion at the HTTP boundary.
SDK timestamps are seconds and rates are nullable bits/second. go-admin keeps
its existing millisecond tuples, A/Z endpoint model and permissions. Source
errors keep their error chains; the provider maps them to safe application
errors. `ErrMonitoringInput`, `ErrMonitoringBinding` and `ErrMonitoringData`
separate invalid requests, unavailable bindings and malformed observations.

The 1h/6h/24h/7d algorithms retain their frozen go-admin behavior. Longer trend
windows now remain bounded to 720 plot points; partial results never prove the
whole-period peak. Query limits and completeness are documented in
[monitoring.md](docs/monitoring.md).

`NewFrontend` is independent of `New`: supply a frontend directory URL and a
genuine web login/cookie session. `Chart` accepts item IDs or one graph ID and
an explicit period. See [the 7.4 frontend contract](docs/frontend-7.4-contracts.md).
Existing API authentication and namespace wire contracts are unchanged.
Requests made on a zero `Client` now return `ErrorInvalidParams` instead of
panicking. Construct clients with `New` as before.

Pin `github.com/gatework/go-zabbix@v0.3.0` to adopt these additions. Verify downstream
consumers with `GOWORK=off` and no local `replace`.

# v0.2.0 standard adoption

- Go 1.26 or later is required.
- Rename reads of `User.TsProvisioned` to `User.TSProvisioned`. The JSON field remains `ts_provisioned`; no compatibility alias is retained.
- Request/response shapes, mutation omission rules and transport ownership remain unchanged.

# Migrating to the 7.4 contract API

This is a grouped, source-breaking revision of the pre-v1 library. The module
remains `github.com/gatework/go-zabbix`; the package identifier is now `zabbix`.
Callers of `github.com/Deepreo/gozabbix` must change their import and dependency.
Pin `github.com/gatework/go-zabbix@v0.2.0` when adopting the current API.

## Authentication and connection ownership

| Previous call | Replacement |
| --- | --- |
| `WithToken(apiToken)` | `WithAPIToken(apiToken)` |
| `WithToken(sessionID)` | `WithSessionToken(sessionID)` |
| `SetToken(value)` | `SetAPIToken(value)` or `SetSessionToken(value)`; check the error |
| `SetToken("")` | `ClearCredentials()` |
| `Authenticated()` | `HasCredentials()` checks local presence, not server validity |
| `User().Login(ctx, user, password)` | `Client.Login(ctx, user, password)` returns an error |
| `User().LoginWithUserData(...)` | `Client.LoginWithUserData(...)` |
| `User().Logout(ctx)` | `Client.Logout(ctx)` |

The token getter is removed. Keep credentials at the point of configuration, or
use the `Session` returned by `LoginWithUserData` when session details are needed.
Logging out an API-token client is an error and leaves the token unchanged.
A superseded authentication operation cannot clear or overwrite newer credentials.
Call `Logout` before replacing a live session if that earlier session must be
revoked. Credential setters and `ClearCredentials` change local state only.

`Session.GUIAccess` and `Session.DebugMode` now use the embedded `User` fields,
both `json.Number`. Use `.String()` or `.Int64()` rather than comparing an `int`
or `string` directly. Keyed session literals must put these fields inside
`User: zabbix.User{...}`. This accepts the server's stored strings and computed
numeric defaults without shadowing either field. A valid session created by
`LoginWithUserData` is logged out if the remaining user data cannot be decoded;
that cleanup is bounded independently of caller cancellation and leaves existing
credentials unchanged.

Different construction options no longer depend on their order. A supplied
`http.Client` is copied. Its transport remains shared unless TLS customization
requires an owned clone. Incompatible TLS options return a constructor error
instead of replacing a custom transport. `CloseIdleConnections` only closes
connections owned by this SDK client.

## Query results and selectors

Move fields out of the removed `GetParams` embedding:

```go
params := zabbix.HostGetParams{
    Output: []string{"hostid", "name"},
    Limit: 10,
}
```

Use `Count(ctx, HostCountParams{...})` for a count and
`GetMap(ctx, HostGetParams{...})` for ID-keyed records. These methods exist only
where the namespace supports the corresponding wire result. `Get` always returns
a record slice. `CountOutput`, `Preservekeys` and other switches that change the
entity or result shape are removed.

Selectors accept `"extend"`, a `[]string` field list, and `"count"` only where
the 7.4 method permits it. Count-capable associations use `Relation[T]`:
`Items` holds records and `Count` is a `*uint64`. A non-nil count pointer can
represent zero; both fields nil means the association was not selected.
An empty selected array is distinct from an omitted association.

Single-object selectors that may return an empty array use `SelectedObject[T]`.
`Present` distinguishes an omitted field; when present, a nil `Value` represents
the server's empty-array result. Host inventory uses `HostInventory`, retaining
the difference between an omitted inventory and an empty selected inventory.

Read the [contract matrix](docs/api-7.4-contracts.md) for supported parameters
and results. In particular, `trend.get` has no sorting or search contract.
A truncated trend response cannot be described as the newest records. Keep range
limits and completeness decisions in the consuming business logic.

## Writes: omitted, empty and zero

Use dedicated `CreateParams` and `UpdateParams` instead of sending a read model
back to the server. Required identifiers remain strings. Optional scalar fields
are pointers; optional replacement collections use `omitzero`.

| Input | Update meaning |
| --- | --- |
| nil scalar pointer | Leave the field unchanged |
| Pointer to `""`, `0`, or `false` | Send that explicit value |
| nil collection | Leave the association unchanged |
| Non-nil empty collection | Clear the association, where the operation permits it |

For example, clearing a host description:

```go
description := ""
params := zabbix.HostUpdateParams{HostID: "10084", Description: &description}
```

Nested mutation inputs exclude read-only metadata. Use their dedicated
`Input`/`Ref` types rather than converting an entire response object. Standalone
host-macro creation includes `HostID`; nested host macros and global macros have
their own inputs. Event acknowledgement can explicitly send severity zero
and suppression deadlines.

Mutation helpers validate the operation-specific result key. In 7.4,
`itemprototype.delete` returns `prototypeids` and `discoveryrule.delete` returns
`ruleids`. The observed nested `hostinterface.massadd` response is normalized to
its ID list. Removing nonmatching interfaces or replacing them with an empty set
may legitimately return an empty ID list. Event acknowledgement, MFA deletion,
and interface mass removal normalize integer response IDs without passing
through floating point. Other ID results retain their string-only contract.
Create receipts must match the number of submitted objects; update/delete
receipts must match the submitted ID multiplicities, regardless of order or
leading zeros. Documented Task and bulk-interface exceptions are preserved.
These checks validate the receipt; use a read-back to verify persisted state. See the
[mutation matrix](docs/api-7.4-mutations.md) for method-specific inputs and results.

`Autoregistration.Update` now returns only `error`, matching the server's boolean
success result. Replace `_, err := client.Autoregistration().Update(...)` with
`err := client.Autoregistration().Update(...)`. Authentication, housekeeping and
settings still return field names, which must match the submitted properties
(including the server-derived tile properties for a named geomaps provider).

## Role service access and HTTP items

`RoleRules.ServicesReadTag` and `ServicesWriteTag` now use `ServiceAccessTag`;
their writable counterparts use `*ServiceAccessTagInput`. These permission
conditions differ from ordinary entity tags: `&ServiceAccessTagInput{Tag: ""}`
clears a condition. A nil pointer preserves it, and a nonempty tag with an empty
value matches any value. `TagInput` still requires a nonempty entity tag name.

Ordinary `Item` reads and `ItemCreateParams`/`ItemUpdateParams` now expose HTTP
agent configuration, ordered `QueryFields`/`Headers`, and `Preprocessing`.
Use `ItemHTTPFieldInput{Name: "X-Source", Value: "monitoring"}` for a header.
Lists preserve repeated names and empty values. A nil list omits it; an allocated
empty list clears it. HTTP item creation (type 19) requires URL and interval;
partial updates can rely on stored properties. Body syntax, macros and
state-dependent constraints remain server-validated.

Use `ItemGetParams{SelectPreprocessing: "extend"}` or a field list to read steps;
the selector does not accept `"count"`. Supplied `ItemPreprocessingInput` steps
replace the entire step configuration, so required parameters and error-handler
properties must be present. A fallback value can explicitly be empty; a custom
error message cannot. This shared validation also applies to item prototypes.

## Specialized operations

- `Script.GetScriptsByHosts` and `GetScriptsByEvents` now take per-host/event
  parameter arrays and return `map[string][]Script`. Use `ScriptHostParams` or
  `ScriptEventParams`; a `*string` manual input can explicitly send empty text.
  `ScriptExecuteParams.ManualInput` is also a `*string`.
- `SLA.GetSLI` returns `*SLIResult` instead of `any`. `Periods` and `ServiceIDs`
  index the two dimensions of `SLI`. IDs remain exact strings even when the
  server encodes integers larger than 2^53. `SLAGetSLIParams.PeriodFrom` and
  `PeriodTo` are pointers; nil and Unix timestamp zero are distinct. Nil
  `ServiceIDs` selects all matching services, while an allocated empty slice
  selects none.
- `Configuration.Export` returns a string for `json`, `xml`, or `yaml`.
  For the `raw` object format, use `ExportRaw(ctx, options)`, which returns
  `json.RawMessage` without losing dynamic numeric values.
- Token generation must return a secret for every requested token ID. Script
  execution must return a successful response with a value (empty text is
  valid). Incomplete results return errors instead of false success.

## Errors and raw calls

Use `errors.As` to inspect `*RequestError`, `*HTTPError` or `*APIError`, and
`errors.Is` for cancellation and timeout causes. Default error text contains
request metadata; upstream API details remain available through `APIError`
for deliberate inspection. Application-facing errors still require appropriate
redaction of those upstream fields.

`Call(ctx, method, params, out)` remains the raw protocol escape hatch. Its
signature and atomic-output contract are unchanged; it does not supply typed
namespace validation or business-level completeness guarantees.
