# Zabbix 7.4 read contracts

This matrix covers every existing SDK namespace with a `Get` method: 52 collection methods and four singleton methods. It aligns the existing public query fields and selectors with Zabbix 7.4; it does not add every selector available on the server. `Client.Call` remains the escape hatch for additional API capabilities.

Evidence is the official 7.4 API manual and Zabbix server source pinned to [`0de072729e7e8dffd72f4176ee019ef0efade637`](https://github.com/zabbix/zabbix/tree/0de072729e7e8dffd72f4176ee019ef0efade637). Each namespace below links both sources. Source takes precedence where the manual omits a supported field. A read-only probe against Zabbix 7.4.16 exercised all 52 collection methods with their exposed selectors and all four singletons; populated and empty relation shapes also have HTTP fixture coverage.

## Result shapes

| API | Public contract |
| --- | --- |
| Collection `Get(ctx, XGetParams)` | `[]T`; a successful empty result is a non-nil slice. Count and keyed results cannot be requested through these params. |
| `Count(ctx, XCountParams)` | `uint64`; only where `countOutput` is supported. Its independent params contain filters, without selectors, output, sort, or limit. Decimal strings and JSON integers decode without floating-point conversion; negative, fractional, malformed and overflowing values fail. |
| `GetMap(ctx, XGetParams)` | `map[string]T`; only where `preservekeys` is supported. Accepts a keyed object or the server's empty `[]` and returns a non-nil map. Non-empty arrays, null entries and non-object values fail. |
| Singleton `Get(ctx)` | `*T`; requires an object, returns nil on failure. |
| `Relation[T]` | `Items []T` or `Count *uint64`, for selectors that support count. Both nil means absent, non-nil empty Items means selected empty list, Count pointing to 0 means selected zero count. |
| `SelectedObject[T]` | `Present bool` plus `Value *T`, for object selectors that can return empty `[]`. Absent has Present false; selected empty has Present true and nil Value; an object sets Value, including `{}`. |
| `HostInventory` | `map[string]string`; nil means absent. An object or empty `[]` produces a non-nil map. Other shapes and non-string values fail. |

`Relation`, `SelectedObject`, inventory and interface-details decoding are atomic: failure leaves an existing destination unchanged. Collection helpers return no partial result when an element is invalid. Errors preserve the method and `RequestError` category: invalid selector arguments are `ErrorInvalidParams`, invalid typed results are `ErrorDecode`; envelope/transport failures retain the core category.

There is no universal embedded `GetParams`. Output and selectors accept `"extend"` or `[]string`; only the selectors marked **yes** in the table below accept `"count"`. Main `output: "count"` is invalid: use `Count`. Selection field names and other server constraints are still validated by Zabbix. Query slices and maps use Go 1.24 `omitzero`: nil omits a filter while a non-nil empty collection is transmitted. This prevents an empty ID list from silently becoming an unrestricted query. Optional timestamp pointers preserve an explicit epoch value of 0.

## Namespace capabilities and exposed query fields

The field column lists the actual JSON query names apart from selectors. `GetMap` adds `preservekeys: true` and `Count` adds `countOutput: true` internally.

| Namespace | Count | GetMap | Query fields | Evidence |
| --- | --- | --- | --- | --- |
| `action.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `actionids`, `groupids`, `hostids`, `triggerids`, `mediatypeids`, `usrgrpids`, `userids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/action/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CAction.php) |
| `alert.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `alertids`, `actionids`, `eventids`, `groupids`, `hostids`, `mediatypeids`, `userids`, `time_from`, `time_till` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/alert/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CAlert.php) |
| `auditlog.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `userids`, `time_from`, `time_till` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/auditlog/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CAuditLog.php) |
| `authentication.get` | — | — | singleton; fixed `output: "extend"` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/authentication/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CAuthentication.php) |
| `autoregistration.get` | — | — | singleton; fixed `output: "extend"` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/autoregistration/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CAutoregistration.php) |
| `connector.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `connectorids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/connector/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CConnector.php) |
| `correlation.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `correlationids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/correlation/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CCorrelation.php) |
| `dashboard.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `dashboardids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dashboard/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CDashboard.php) |
| `dcheck.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `dcheckids`, `druleids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dcheck/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CDCheck.php) |
| `dhost.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `dhostids`, `druleids`, `dserviceids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dhost/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CDHost.php) |
| `discoveryrule.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `itemids`, `hostids`, `groupids`, `templateids`, `interfaceids`, `inherited`, `templated`, `monitored` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/discoveryrule/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CDiscoveryRule.php) |
| `drule.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `druleids`, `dhostids`, `dserviceids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/drule/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CDRule.php) |
| `dservice.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `dserviceids`, `dhostids`, `druleids`, `dcheckids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dservice/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CDService.php) |
| `event.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `eventids`, `groupids`, `hostids`, `objectids`, `source`, `object`, `severities`, `acknowledged`, `suppressed`, `value`, `tags`, `evaltype`, `time_from`, `time_till` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/event/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CEvent.php) |
| `graph.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `graphids`, `groupids`, `templateids`, `hostids`, `itemids`, `templated` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graph/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CGraph.php) |
| `graphitem.get` | yes | yes | `output`, `sortfield`, `sortorder`, `limit`, `editable`, `graphids`, `itemids`, `type` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graphitem/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CGraphItem.php) |
| `graphprototype.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `graphids`, `groupids`, `templateids`, `hostids`, `itemids`, `discoveryids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graphprototype/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CGraphPrototype.php) |
| `hanode.get` | yes | yes | `output`, `filter`, `sortfield`, `sortorder`, `limit`, `ha_nodeids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hanode/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CHaNode.php) |
| `history.get` | yes | no | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `history`, `itemids`, `hostids`, `time_from`, `time_till` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/history/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CHistory.php) |
| `host.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `hostids`, `groupids`, `templateids`, `itemids`, `triggerids`, `maintenanceids`, `proxyids`, `tags`, `evaltype`, `with_items`, `with_triggers`, `monitored_hosts` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/host/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CHost.php) |
| `hostgroup.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `groupids`, `hostids`, `graphids`, `triggerids`, `maintenanceids`, `with_hosts`, `with_monitored_hosts`, `with_items`, `with_triggers` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostgroup/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CHostGroup.php) |
| `hostinterface.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `interfaceids`, `hostids`, `triggerids`, `itemids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CHostInterface.php) |
| `hostprototype.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `hostids`, `discoveryids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostprototype/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CHostPrototype.php) |
| `housekeeping.get` | — | — | singleton; fixed `output: "extend"` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/housekeeping/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CHousekeeping.php) |
| `iconmap.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `iconmapids`, `sysmapids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/iconmap/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CIconMap.php) |
| `image.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `imageids`, `sysmapids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/image/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CImage.php) |
| `item.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `itemids`, `hostids`, `groupids`, `templateids`, `interfaceids`, `triggerids`, `tags`, `evaltype`, `templated`, `monitored`, `webitems` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/item/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CItem.php) |
| `itemprototype.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `itemids`, `hostids`, `groupids`, `templateids`, `discoveryids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/itemprototype/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CItemPrototype.php) |
| `maintenance.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `maintenanceids`, `groupids`, `hostids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/maintenance/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CMaintenance.php) |
| `map.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `sysmapids`, `userids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/map/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CMap.php) |
| `mediatype.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `mediatypeids`, `mediaids`, `userids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mediatype/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CMediatype.php) |
| `mfa.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `mfaids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mfa/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CMfa.php) |
| `module.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `moduleids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/module/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CModule.php) |
| `problem.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `eventids`, `groupids`, `hostids`, `objectids`, `source`, `object`, `severities`, `tags`, `evaltype`, `recent`, `acknowledged`, `suppressed`, `time_from`, `time_till` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/problem/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CProblem.php) |
| `proxy.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `proxyids`, `proxy_groupids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxy/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CProxy.php) |
| `proxygroup.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `proxy_groupids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxygroup/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CProxyGroup.php) |
| `regexp.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `regexpids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/regexp/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CRegexp.php) |
| `report.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `reportids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/report/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CReport.php) |
| `role.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `roleids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/role/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CRole.php) |
| `script.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `scriptids`, `hostids`, `groupids`, `usrgrpids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CScript.php) |
| `service.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `serviceids`, `parentids`, `childids`, `tags`, `problem_tags` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/service/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CService.php) |
| `settings.get` | — | — | singleton; fixed `output: "extend"` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/settings/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CSettings.php) |
| `sla.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `slaids`, `serviceids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/sla/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CSla.php) |
| `task.get` | no | yes | `output`, `taskids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/task/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CTask.php) |
| `template.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `templateids`, `groupids`, `parentTemplateids`, `hostids`, `itemids`, `triggerids`, `tags`, `evaltype` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/template/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CTemplate.php) |
| `templatedashboard.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `dashboardids`, `templateids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templatedashboard/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CTemplateDashboard.php) |
| `templategroup.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `groupids`, `templateids`, `graphids`, `triggerids`, `with_items`, `with_triggers` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templategroup/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CTemplateGroup.php) |
| `token.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `tokenids`, `userids`, `valid_at`, `expired_at` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/token/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CToken.php) |
| `trend.get` | yes | no | `output`, `limit`, `itemids`, `time_from`, `time_till` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/trend/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CTrend.php) |
| `trigger.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `triggerids`, `groupids`, `templateids`, `hostids`, `itemids`, `tags`, `evaltype`, `monitored`, `active`, `maintenance`, `min_severity`, `only_true`, `skipDependent`, `expandExpression` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/trigger/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CTrigger.php) |
| `triggerprototype.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `triggerids`, `groupids`, `templateids`, `hostids`, `itemids`, `discoveryids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/triggerprototype/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CTriggerPrototype.php) |
| `user.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `userids`, `mediaids`, `mediatypeids`, `usrgrpids`, `getAccess` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/user/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CUser.php) |
| `userdirectory.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `userdirectoryids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/userdirectory/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CUserDirectory.php) |
| `usergroup.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `usrgrpids`, `userids`, `status` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usergroup/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CUserGroup.php) |
| `usermacro.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `hostids`, `hostmacroids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CUserMacro.php) |
| `valuemap.get` | yes | yes | `output`, `filter`, `search`, `searchByAny`, `searchWildcardsEnabled`, `startSearch`, `excludeSearch`, `sortfield`, `sortorder`, `limit`, `editable`, `valuemapids`, `hostids` | [manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/valuemap/get), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CValueMap.php) |

`usermacro.get` is split by entity in the SDK. `UserMacro.Get/Count/GetMap` accept host/template macro parameters and return `UserMacro`. `GetGlobal/CountGlobal/GetGlobalMap` accept `GlobalMacroGetParams` or `GlobalMacroCountParams` and return `GlobalMacro`; these methods set `globalmacro: true` internally. Host params cannot select the global entity. Global params replace `hostids` and `hostmacroids` with `globalmacroids`.

## Selector response matrix

The response field is the Go field and its JSON wire key. A `Relation[T]` is used exactly where count is accepted; plain slices retain a list-only contract. The namespace evidence above covers each selector in this table. `namespace_contract_test.go` verifies all 119 selectors have an outgoing field and a populated decoded result field, and checks their count capability.

| Namespace | Selector | Count | Response field and type |
| --- | --- | --- | --- |
| `action` | `selectOperations` | no | `Action.Operations` (`operations`): `[]ActionOperation` |
| `action` | `selectRecoveryOperations` | no | `Action.RecoveryOperations` (`recovery_operations`): `[]ActionRecoveryOperation` |
| `action` | `selectUpdateOperations` | no | `Action.UpdateOperations` (`update_operations`): `[]ActionUpdateOperation` |
| `action` | `selectFilter` | no | `Action.Filter` (`filter`): `SelectedObject[ActionFilter]` |
| `alert` | `selectHosts` | no | `Alert.Hosts` (`hosts`): `[]Host` |
| `alert` | `selectMediatypes` | no | `Alert.MediaTypes` (`mediatypes`): `[]MediaType` |
| `alert` | `selectUsers` | no | `Alert.Users` (`users`): `[]User` |
| `connector` | `selectTags` | yes | `Connector.Tags` (`tags`): `Relation[ConnectorTag]` |
| `correlation` | `selectFilter` | no | `Correlation.Filter` (`filter`): `SelectedObject[CorrelationFilter]` |
| `correlation` | `selectOperations` | no | `Correlation.Operations` (`operations`): `[]CorrelationOperation` |
| `dashboard` | `selectPages` | no | `Dashboard.Pages` (`pages`): `[]DashboardPage` |
| `dashboard` | `selectUsers` | no | `Dashboard.Users` (`users`): `[]DashboardUser` |
| `dashboard` | `selectUserGroups` | no | `Dashboard.UserGroups` (`userGroups`): `[]DashboardUserGroup` |
| `dhost` | `selectDServices` | yes | `DHost.DServices` (`dservices`): `Relation[DService]` |
| `discoveryrule` | `selectHosts` | no | `DiscoveryRule.Hosts` (`hosts`): `[]Host` |
| `discoveryrule` | `selectFilter` | no | `DiscoveryRule.Filter` (`filter`): `SelectedObject[DiscoveryRuleFilter]` |
| `discoveryrule` | `selectLLDMacroPaths` | no | `DiscoveryRule.LLDMacroPaths` (`lld_macro_paths`): `[]DiscoveryRuleLLDMacroPath` |
| `discoveryrule` | `selectPreprocessing` | no | `DiscoveryRule.Preprocessing` (`preprocessing`): `[]DiscoveryRulePreprocessing` |
| `drule` | `selectDChecks` | yes | `DRule.DChecks` (`dchecks`): `Relation[DCheck]` |
| `drule` | `selectDHosts` | yes | `DRule.DHosts` (`dhosts`): `Relation[DHost]` |
| `dservice` | `selectDRules` | no | `DService.DRules` (`drules`): `[]DRule` |
| `dservice` | `selectDHosts` | no | `DService.DHosts` (`dhosts`): `[]DHost` |
| `dservice` | `selectHosts` | yes | `DService.Hosts` (`hosts`): `Relation[Host]` |
| `event` | `selectHosts` | no | `Event.Hosts` (`hosts`): `[]Host` |
| `event` | `selectTags` | no | `Event.Tags` (`tags`): `[]Tag` |
| `event` | `selectAcknowledges` | yes | `Event.Acknowledges` (`acknowledges`): `Relation[EventUpdate]` |
| `graph` | `selectGraphItems` | no | `Graph.GItems` (`gitems`): `[]GraphItem` |
| `graph` | `selectHosts` | no | `Graph.Hosts` (`hosts`): `[]Host` |
| `graph` | `selectItems` | no | `Graph.Items` (`items`): `[]Item` |
| `graph` | `selectTemplates` | no | `Graph.Templates` (`templates`): `[]Template` |
| `graphitem` | `selectGraphs` | no | `GraphItem.Graphs` (`graphs`): `[]Graph` |
| `graphprototype` | `selectGraphItems` | no | `GraphPrototype.GItems` (`gitems`): `[]GraphItem` |
| `graphprototype` | `selectHosts` | no | `GraphPrototype.Hosts` (`hosts`): `[]Host` |
| `graphprototype` | `selectItems` | no | `GraphPrototype.Items` (`items`): `[]ItemPrototype` |
| `graphprototype` | `selectDiscoveryRule` | no | `GraphPrototype.DiscoveryRule` (`discoveryRule`): `SelectedObject[DiscoveryRule]` |
| `host` | `selectHostGroups` | no | `Host.HostGroups` (`hostgroups`): `[]HostGroup` |
| `host` | `selectInterfaces` | yes | `Host.Interfaces` (`interfaces`): `Relation[HostInterface]` |
| `host` | `selectItems` | yes | `Host.Items` (`items`): `Relation[Item]` |
| `host` | `selectTriggers` | yes | `Host.Triggers` (`triggers`): `Relation[Trigger]` |
| `host` | `selectParentTemplates` | yes | `Host.ParentTemplates` (`parentTemplates`): `Relation[Template]` |
| `host` | `selectTags` | no | `Host.Tags` (`tags`): `[]Tag` |
| `host` | `selectMacros` | no | `Host.Macros` (`macros`): `[]UserMacro` |
| `host` | `selectInventory` | no | `Host.Inventory` (`inventory`): `HostInventory` |
| `hostgroup` | `selectHosts` | yes | `HostGroup.Hosts` (`hosts`): `Relation[Host]` |
| `hostinterface` | `selectItems` | yes | `HostInterface.Items` (`items`): `Relation[Item]` |
| `hostinterface` | `selectHosts` | no | `HostInterface.Hosts` (`hosts`): `[]Host` |
| `hostprototype` | `selectGroupLinks` | no | `HostPrototype.GroupLinks` (`groupLinks`): `[]HostPrototypeGroupLink` |
| `hostprototype` | `selectGroupPrototypes` | no | `HostPrototype.GroupPrototypes` (`groupPrototypes`): `[]HostPrototypeGroupPrototype` |
| `hostprototype` | `selectDiscoveryRule` | no | `HostPrototype.DiscoveryRule` (`discoveryRule`): `SelectedObject[DiscoveryRule]` |
| `hostprototype` | `selectParentHost` | no | `HostPrototype.ParentHost` (`parentHost`): `SelectedObject[Host]` |
| `hostprototype` | `selectInterfaces` | no | `HostPrototype.Interfaces` (`interfaces`): `[]HostInterface` |
| `hostprototype` | `selectTemplates` | yes | `HostPrototype.Templates` (`templates`): `Relation[Template]` |
| `hostprototype` | `selectTags` | no | `HostPrototype.Tags` (`tags`): `[]Tag` |
| `hostprototype` | `selectMacros` | no | `HostPrototype.Macros` (`macros`): `[]UserMacro` |
| `iconmap` | `selectMappings` | no | `IconMap.Mappings` (`mappings`): `[]IconMapMapping` |
| `item` | `selectHosts` | no | `Item.Hosts` (`hosts`): `[]Host` |
| `item` | `selectTriggers` | yes | `Item.Triggers` (`triggers`): `Relation[Trigger]` |
| `item` | `selectTags` | no | `Item.Tags` (`tags`): `[]Tag` |
| `itemprototype` | `selectHosts` | no | `ItemPrototype.Hosts` (`hosts`): `[]Host` |
| `itemprototype` | `selectTags` | no | `ItemPrototype.Tags` (`tags`): `[]Tag` |
| `itemprototype` | `selectPreprocessing` | no | `ItemPrototype.Preprocessing` (`preprocessing`): `[]ItemPrototypePreprocessing` |
| `itemprototype` | `selectDiscoveryRule` | no | `ItemPrototype.DiscoveryRule` (`discoveryRule`): `SelectedObject[DiscoveryRule]` |
| `maintenance` | `selectTimeperiods` | no | `Maintenance.Timeperiods` (`timeperiods`): `[]MaintenanceTimeperiod` |
| `maintenance` | `selectTags` | no | `Maintenance.Tags` (`tags`): `[]MaintenanceTag` |
| `maintenance` | `selectHostGroups` | no | `Maintenance.HostGroups` (`hostgroups`): `[]HostGroup` |
| `maintenance` | `selectHosts` | no | `Maintenance.Hosts` (`hosts`): `[]Host` |
| `map` | `selectSelements` | no | `Map.Selements` (`selements`): `[]MapSelement` |
| `map` | `selectLinks` | no | `Map.Links` (`links`): `[]MapLink` |
| `map` | `selectUrls` | no | `Map.URLs` (`urls`): `[]MapURL` |
| `map` | `selectUsers` | no | `Map.Users` (`users`): `[]MapUser` |
| `map` | `selectUserGroups` | no | `Map.UserGroups` (`userGroups`): `[]MapUserGroup` |
| `map` | `selectShapes` | no | `Map.Shapes` (`shapes`): `[]MapShape` |
| `map` | `selectLines` | no | `Map.Lines` (`lines`): `[]MapLine` |
| `mediatype` | `selectMessageTemplates` | no | `MediaType.MessageTemplates` (`message_templates`): `[]MediaTypeMessageTemplate` |
| `mediatype` | `selectUsers` | no | `MediaType.Users` (`users`): `[]User` |
| `problem` | `selectAcknowledges` | yes | `Problem.Acknowledges` (`acknowledges`): `Relation[EventUpdate]` |
| `problem` | `selectTags` | no | `Problem.Tags` (`tags`): `[]Tag` |
| `problem` | `selectSuppressionData` | no | `Problem.SuppressionData` (`suppression_data`): `[]SuppressionData` |
| `proxy` | `selectHosts` | yes | `Proxy.Hosts` (`hosts`): `Relation[Host]` |
| `proxy` | `selectProxyGroup` | no | `Proxy.ProxyGroup` (`proxyGroup`): `SelectedObject[ProxyGroup]` |
| `proxygroup` | `selectProxies` | yes | `ProxyGroup.Proxies` (`proxies`): `Relation[Proxy]` |
| `regexp` | `selectExpressions` | no | `Regexp.Expressions` (`expressions`): `[]RegexpExpression` |
| `report` | `selectUsers` | no | `Report.Users` (`users`): `[]ReportUser` |
| `report` | `selectUserGroups` | no | `Report.UserGroups` (`user_groups`): `[]ReportUserGroup` |
| `role` | `selectRules` | no | `Role.Rules` (`rules`): `SelectedObject[RoleRules]` |
| `script` | `selectHosts` | no | `Script.Hosts` (`hosts`): `[]Host` |
| `script` | `selectHostGroups` | no | `Script.HostGroups` (`hostgroups`): `[]HostGroup` |
| `service` | `selectChildren` | yes | `Service.Children` (`children`): `Relation[Service]` |
| `service` | `selectParents` | yes | `Service.Parents` (`parents`): `Relation[Service]` |
| `service` | `selectTags` | yes | `Service.Tags` (`tags`): `Relation[Tag]` |
| `service` | `selectProblemTags` | yes | `Service.ProblemTags` (`problem_tags`): `Relation[ServiceProblemTag]` |
| `service` | `selectStatusRules` | yes | `Service.StatusRules` (`status_rules`): `Relation[ServiceStatusRule]` |
| `sla` | `selectSchedule` | yes | `SLA.Schedule` (`schedule`): `Relation[SLASchedule]` |
| `sla` | `selectExcludedDowntimes` | yes | `SLA.ExcludedDowntimes` (`excluded_downtimes`): `Relation[SLAExcludedDowntime]` |
| `sla` | `selectServiceTags` | yes | `SLA.ServiceTags` (`service_tags`): `Relation[SLAServiceTag]` |
| `template` | `selectTemplateGroups` | no | `Template.TemplateGroups` (`templategroups`): `[]TemplateGroup` |
| `template` | `selectHosts` | yes | `Template.Hosts` (`hosts`): `Relation[Host]` |
| `template` | `selectItems` | yes | `Template.Items` (`items`): `Relation[Item]` |
| `template` | `selectTriggers` | yes | `Template.Triggers` (`triggers`): `Relation[Trigger]` |
| `template` | `selectTags` | no | `Template.Tags` (`tags`): `[]Tag` |
| `template` | `selectMacros` | no | `Template.Macros` (`macros`): `[]UserMacro` |
| `templatedashboard` | `selectPages` | no | `TemplateDashboard.Pages` (`pages`): `[]DashboardPage` |
| `templategroup` | `selectTemplates` | yes | `TemplateGroup.Templates` (`templates`): `Relation[Template]` |
| `trigger` | `selectHosts` | no | `Trigger.Hosts` (`hosts`): `[]Host` |
| `trigger` | `selectItems` | no | `Trigger.Items` (`items`): `[]Item` |
| `trigger` | `selectTags` | no | `Trigger.Tags` (`tags`): `[]Tag` |
| `trigger` | `selectDependencies` | no | `Trigger.Dependencies` (`dependencies`): `[]Trigger` |
| `triggerprototype` | `selectHosts` | no | `TriggerPrototype.Hosts` (`hosts`): `[]Host` |
| `triggerprototype` | `selectItems` | no | `TriggerPrototype.Items` (`items`): `[]ItemPrototype` |
| `triggerprototype` | `selectTags` | no | `TriggerPrototype.Tags` (`tags`): `[]Tag` |
| `triggerprototype` | `selectDependencies` | no | `TriggerPrototype.Dependencies` (`dependencies`): `[]TriggerPrototype` |
| `user` | `selectMedias` | no | `User.Medias` (`medias`): `[]Media` |
| `user` | `selectUsrgrps` | no | `User.UserGroups` (`usrgrps`): `[]UserGroup` |
| `user` | `selectRole` | no | `User.Role` (`role`): `SelectedObject[Role]` |
| `usergroup` | `selectUsers` | no | `UserGroup.Users` (`users`): `[]User` |
| `usergroup` | `selectTagFilters` | no | `UserGroup.TagFilters` (`tag_filters`): `[]UserGroupTagFilter` |
| `usergroup` | `selectHostGroupRights` | no | `UserGroup.HostGroupRights` (`hostgroup_rights`): `[]UserGroupRight` |
| `usergroup` | `selectTemplateGroupRights` | no | `UserGroup.TemplateGroupRights` (`template_group_rights`): `[]UserGroupRight` |
| `valuemap` | `selectMappings` | yes | `ValueMap.Mappings` (`mappings`): `Relation[ValueMapMapping]` |

`image.get` additionally exposes boolean `select_image`, returning the base64 string in `Image.Image`; it is not a string/field-list selector and does not accept count. `user.get` exposes boolean `getAccess`; `GUIAccess`, `DebugMode`, and `UsersStatus` use `json.Number` because computed defaults are JSON numbers and stored values are strings. `Service.ReadOnly` is a JSON boolean. `Role.Rules` is a selected typed `RoleRules` object; its module IDs use `json.Number` because PHP array keys can produce numeric IDs. These representations are verified against the 7.4 source and fixtures.

## Corrected protocol differences

- `TemplateGetParams.SelectTemplateGroups` decodes into `Template.TemplateGroups` using `templategroups`; host groups and parent templates use `hostgroups` and `parentTemplates`. Response models no longer expose the obsolete `groups` write keys.
- `HANodeGetParams.Filter["status"]` replaces an unsupported top-level status option. `Connector` data type and `Map` expandproblem likewise belong in property filters, not top-level query options.
- `TrendGetParams` contains only output, item IDs, time bounds and limit. Zabbix 7.4 does not support trend sorting, searching, editable or preservekeys; server SQL provides no ordering guarantee. `History` also has no GetMap variant.
- `Task` has Get and GetMap, with output and task IDs, but no Count, filters, sorting, limit or selectLastValues.
- `UserGroup` rights are split into `HostGroupRights` and `TemplateGroupRights`, with `selectHostGroupRights` and `selectTemplateGroupRights`. The old rights/selectRights pair is absent in 7.4.
- `Tag` is a returned tag, while `TagFilter` carries the integer query operator. Connector/service/SLA/maintenance tag objects retain their distinct protocol meanings.
- `HostInterface.Details` accepts an object or empty `[]` for non-SNMP interfaces. It rejects null, scalar values and non-empty arrays, and does not leave partial updates behind.
- `TokenGetParams.ValidAt` and `ExpiredAt` replace the unsupported `valid` field. Timestamp pointers allow 0 to be transmitted.
- Kept source-supported fields omitted by the manual: `valuemap.get hostids`, `itemprototype.get groupids`, `triggerprototype.get itemids`. Removed unsupported fields include hostgroup.real_hosts, hostprototype.inheritedTags and itemprototype.interfaceids.

## Script lookup arrays

`GetScriptsByHosts(ctx, []ScriptHostParams)` and `GetScriptsByEvents(ctx, []ScriptEventParams)` send arrays of objects containing `hostid` or `eventid`, optional `scriptid`, and optional `manualinput`. `*string` preserves an explicitly empty manual input. They return `map[string][]Script`, accept the server's empty `[]` for an empty lookup, and reject incorrect nested shapes. These lookups do not execute scripts. [Hosts manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/getscriptsbyhosts), [events manual](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/getscriptsbyevents), [source](https://raw.githubusercontent.com/zabbix/zabbix/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services/CScript.php).

## Validation boundary

Unit fixtures cover emitted JSON, list/count/map variants, large counts, empty and absent selectors, malformed/partial responses, global macro separation, Script lookup arrays, HA status filtering, explicit zero timestamps, user access fields and all exposed selector return fields. `integration_query_contract_test.go` exercises all 51 Count and 50 GetMap collection capabilities plus the global macro variants against a configured real server. Populated CRUD/query integration tests provide separate evidence for mutation behavior. A passing empty-collection probe proves method/parameter acceptance and empty shape decoding; it does not prove populated fields for that namespace.
