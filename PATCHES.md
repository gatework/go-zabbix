# Local response safety patch

- Upstream: https://github.com/Deepreo/gozabbix
- Upstream revision: `432897ae0c316a0f7f662c157ef521bd76548b34`
- Published module: `github.com/gatework/go-zabbix`
- Local branch: `motor-circuit-readonly`
- License: upstream GPL-3.0 `LICENSE`, unchanged.
- Production changes: `client.go` only. The existing public API, service helpers,
  context handling and Bearer authentication are retained.

## Scope

`New(apiURL, ...Option)`, `WithToken`, `WithTimeout`, `WithMaxResponseBytes`,
`WithHTTPClient`, and `Call(ctx, method, params, out)` keep their existing
signatures. No compatibility API from another SDK is added.

The application selects its 20-second request timeout and 4 MiB response limit
through existing options:

```go
client, err := gozabbix.New(apiURL,
    gozabbix.WithToken(token),
    gozabbix.WithTimeout(20*time.Second),
    gozabbix.WithMaxResponseBytes(4<<20),
)
```

Upstream library defaults remain 30 seconds and 64 MiB. Normal TLS verification
and the default TLS 1.2 minimum remain in place. Custom HTTP transports and their
trust configuration remain supported. `APIInfo().Version(ctx)` keeps the
upstream unauthenticated behavior; a client with `WithToken` needs no version
probe before calls.

## Focused changes

1. Successful response reads probe one byte beyond the configured limit. This
   detects overflow even when the truncated prefix is valid JSON. The probe is
   separate to avoid an integer overflow for a maximum int64 limit.
2. Responses must have JSON-RPC version 2.0, a matching numeric request ID, and
   exactly one result or error. Missing/null results, null or malformed errors,
   missing/mismatched IDs and invalid JSON fail explicitly. Valid upstream API
   errors retain the `*APIError` and `AsAPIError` contract.
3. Non-nil output values must be non-nil pointers. Typed decoding uses a fresh
   value and replaces the caller's output only after success, preventing partial
   overwrites. A nil output remains supported and still validates the envelope.
   Successful decoding replaces the old value instead of merging prior fields.
4. Redirects are refused even with an injected HTTP client. A shallow client
   copy preserves the caller's original redirect policy and reuses its transport,
   jar and timeout.
5. Non-200 HTTP errors report status and method only; response bodies are never
   included in those errors. Existing debug logging remains metadata-only.
   Structured API-error message/data are still upstream text, so applications
   remain responsible for public-facing error messages.

## Tests and provenance

`client_safety_test.go` uses loopback `httptest` servers only. The new regressions
were run against the upstream code first and failed for silent response
truncation, missing protocol checks, partial output mutation, redirects and
credential-bearing HTTP error bodies. Existing TLS/context/Bearer behavior is
covered without duplicating its implementation.

The first patched full offline run passed before the module path was renamed:

```text
docker exec circuit-gozabbix-sdk-dev go test -race ./...
ok github.com/Deepreo/gozabbix 1.157s
```

The container uses `golang:1.27.2-bookworm` and mounts this clone at `/workspace`.
No `integration` build tag is enabled. The upstream integration test files
perform live login/logout and CRUD writes and are intentionally excluded.
No live endpoint, credentials, module cache or host Go installation is used or
modified by this patch.

The downstream snapshot manifest records source hashes and the upstream commit.
The local patch digest covers the binary diff of `client.go` against the upstream
revision followed by the new `client_safety_test.go` contents; this document is
excluded to avoid a self-referential digest. Final command results and the digest
are recorded in the task report.

Final local patch SHA-256: `7cf31d800939312b2c2876af2bb1d8b1145e09dc01f40bee157e3170a023d746`.

Final fresh verification: `go test -race -count=1 -v ./...` passed all 21
top-level tests (1.175s); `go vet ./...`, `gofmt -l .`, and
`git diff --check` passed. The three integration-tagged test files were excluded.
