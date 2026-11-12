# Changelog

## v0.3.0 - 2026-11-12

- Move bounded host/item queries and traffic processing from go-admin into `Client.Monitoring` and `ReadTraffic`, preserving units, gaps, peaks and partial-data signals.
- Add a separate Zabbix 7.4 frontend client for real cookie login, fixed-period item/graph PNGs and session logout, with bounded image validation.
- Bound long trend plots, return errors for zero-client requests, and add migration fixtures, executable examples and frontend fuzz/lifecycle regressions.

## v0.2.0 - 2026-11-11

- Adopt the published go-standards v1.0.2 lint, formatting, architecture and make gates without adding runtime dependencies.
- Require Go 1.26 or later and verify both supported Go release lines in CI, including disposable Zabbix 7.4 round trips.
- Rename User.TsProvisioned to User.TSProvisioned; the ts_provisioned wire key and caller-owned transport behavior are unchanged.

## v0.1.0

- First tagged Zabbix 7.4 JSON-RPC SDK release with typed query and mutation contracts.
