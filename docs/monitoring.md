# Device, circuit and alert monitoring

`Client.Monitoring()` owns device/item lookup, bounded history and trends, and
traffic normalization. `Frontend` owns signed web sessions and PNG retrieval.
Both live in go-zabbix and depend only on the standard library. Applications own
their database records, permissions, circuit endpoints and notification delivery.

## Bind the physical interface

Search `Monitoring().Hosts(ctx, keyword)`, then
`Monitoring().Items(ctx, hostID, interfaceKeyword)`. Persist the selected host ID
and inbound/outbound item IDs along with the application's interface label.
`HostInterface` is an SNMP/agent access endpoint, not an Ethernet port. Item names
and SNMP ifIndex alone are not stable identities; revalidate ownership after
discovery or configuration changes. `Traffic` revalidates host and item identity
on each read, including enabled status, numeric type, units and unsupported state.

CPU, memory, availability and interface status use the existing typed `Host`,
`Item`, `History`, `Problem` and `Trigger` APIs. The higher-level `Monitoring`
rate helpers deliberately accept only numeric `bps` and `Bps` items. Counters
must already be preprocessed by Zabbix into rates; the SDK does not derive a
rate from an unlabelled cumulative counter.

## Fetch a period

```go
series, err := client.Monitoring().Traffic(ctx, zabbix.TrafficQuery{
    HostID: "10084", ItemID: "42001",
    From: from.Unix(), Till: till.Unix(), UseTrends: false,
})
```

Keep `from` and `till` fixed for the complete operation. Timestamps in the SDK
are Unix seconds; go-admin converts them to milliseconds at its HTTP DTO boundary.
`ReadTraffic(ctx, source, query)` accepts a small `TrafficSource` interface for
applications that need their own configuration or error mapping around the SDK.

| Operation | Bound and meaning |
| --- | --- |
| Hosts | 50 records; a 51st detects `HasMore` |
| Items | 200 records; a 201st detects `HasMore` |
| History | 12,001 received samples; over 12,000 means `Partial` |
| Trends | 201 received hourly aggregates; over 200 means `Partial` |
| Plot | At most 720 points; missing/invalid values remain nil |

`Bps` values are multiplied by eight. Negative, invalid and non-finite values
are missing evidence, while an observed zero remains zero. History buckets
average observations and insert gap markers using the larger of five minutes
and three times the median observed cadence. Trend peaks use `value_max`, not
`value_avg`. Periods longer than 720 hours use count-weighted trend buckets.
The peak is computed before plot downsampling and includes the truncation probe.

`Partial` means the received observations cannot prove the complete period's
peak. `trend.get` does not support sorting, so a truncated result cannot be
called the newest 201 hours. Use the lower-level typed queries to partition
longer windows when a complete export is required.

`CurrentBPS` uses the item's last observation only when it is at most five
minutes old relative to `Till`; it is independent of trend aggregation.
`LastSampleAt` is retained for older samples, and `Stale` reports known old data.
Unknown/future last clocks do not supply a current value. `SampleCount == 0`
means no valid received rate observations. On a sample-fetch failure, a series
may retain validated identity/latest-sample metadata; always handle the error.

## Attach the Zabbix image to an alert

1. Obtain problems/events with `Problem().Get` or `Event().Get`. For trigger
   events, `objectid` identifies the trigger; `Trigger().Get` with
   `SelectItems` identifies the involved items. A trigger can reference several
   metrics, so use the application's explicit circuit bindings to select the
   intended inbound/outbound pair.
2. Enforce the application's device/circuit permissions and validate that each
   selected item belongs to the expected host. Fix the image window around the
   event timestamp, for example the preceding hour. A future post-event window
   requires the application to delay the attachment job until data exists.
3. Use `Frontend.Login` with a GUI-enabled account, or supply a genuine signed
   CookieJar obtained by an external SSO/MFA flow. API tokens cannot substitute
   for web-session cookies.
4. Call `Frontend.Chart` with both item IDs and the fixed `From`/`To` window,
   or use `GraphID` for an existing configured graph. Send the returned PNG bytes
   as an attachment through the application's notification adapter.
5. Keep event deduplication, retries, delivery records and failed-image policy
   in the application. Never report an image-fetch failure as an observed zero
   rate. The SDK neither sends messages nor marks alerts delivered.

Reuse a `Frontend` across a notification worker's requests. It does not store
passwords or retry login automatically. On authentication failure, the worker
can explicitly renew the session and retry the **same** fixed window. On
shutdown, call `Logout` with a fresh bounded context; `ClearSession` only forgets
local cookies. See [the frontend contract](frontend-7.4-contracts.md) for timezone,
SSO/MFA and image-validation limits.

## go-admin ownership

go-admin's provider keeps environment configuration, HTTPS policy, shared
transport ownership, request budgets and mapping SDK errors to application
errors. Circuit services retain DataScope, A/Z endpoints, period selection and
the existing UI JSON shape. Query construction, truncation, rate conversion,
gap detection and peak calculation belong exclusively to go-zabbix.

The frozen `testdata/traffic_go_admin.json` was produced by go-admin commit
`2d9198625edccc4fb8f2812f26c891c4654bf225` before its algorithms were removed.
It covers 16 combinations of history/trend, bps/Bps and 1h/6h/24h/7d periods.
Do not regenerate it from the replacement implementation. Longer-window
bounding is an intentional extension tested separately.
