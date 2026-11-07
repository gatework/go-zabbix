# Zabbix 7.4 mutation contracts

The SDK separates returned objects from mutation inputs. This table covers the
existing typed mutations; it is not a list of every Zabbix write operation.
Fields are documented on the concrete types in [write_params.go](../write_params.go),
[write_types.go](../write_types.go) and [write_operations.go](../write_operations.go).
Each row links the versioned method documentation. Result keys and conditional
rules were also checked against [the pinned 7.4 server source](https://github.com/zabbix/zabbix/tree/0de072729e7e8dffd72f4176ee019ef0efade637/ui/include/classes/api/services).

## Shared rules

- Optional scalar pointers distinguish omission from explicit zero, false or empty text.
- Nil replacement collections omit the property. Non-nil empty collections send `[]`; server rules decide whether clearing is permitted.
- Required IDs and returned IDs remain decimal strings, without floating-point conversion.
- Nested `Input` and `Ref` types contain writable values and identity references, not response metadata.
- Empty batches, missing required IDs and malformed result envelopes fail locally. State-dependent permissions and inherited/provisioned restrictions remain server checks.
- Partial updates need not repeat unchanged server state: a host can retain its proxy, and a service can change its propagation value while retaining its rule.

## Method matrix

| Method | Go method | Input | Successful result |
| --- | --- | --- | --- |
| [`action.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/action/create) | [`Action.Create`](../action.go) | `...ActionCreateParams` | `actionids` IDs |
| [`action.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/action/delete) | [`Action.Delete`](../action.go) | `...string` | `actionids` IDs |
| [`action.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/action/update) | [`Action.Update`](../action.go) | `...ActionUpdateParams` | `actionids` IDs |
| [`authentication.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/authentication/update) | [`Authentication.Update`](../authentication.go) | `AuthenticationUpdateParams` | updated field names |
| [`autoregistration.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/autoregistration/update) | [`Autoregistration.Update`](../autoregistration.go) | `AutoregistrationUpdateParams` | updated field names |
| [`connector.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/connector/create) | [`Connector.Create`](../connector.go) | `...ConnectorCreateParams` | `connectorids` IDs |
| [`connector.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/connector/delete) | [`Connector.Delete`](../connector.go) | `...string` | `connectorids` IDs |
| [`connector.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/connector/update) | [`Connector.Update`](../connector.go) | `...ConnectorUpdateParams` | `connectorids` IDs |
| [`correlation.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/correlation/create) | [`Correlation.Create`](../correlation.go) | `...CorrelationCreateParams` | `correlationids` IDs |
| [`correlation.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/correlation/delete) | [`Correlation.Delete`](../correlation.go) | `...string` | `correlationids` IDs |
| [`correlation.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/correlation/update) | [`Correlation.Update`](../correlation.go) | `...CorrelationUpdateParams` | `correlationids` IDs |
| [`dashboard.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dashboard/create) | [`Dashboard.Create`](../dashboard.go) | `...DashboardCreateParams` | `dashboardids` IDs |
| [`dashboard.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dashboard/delete) | [`Dashboard.Delete`](../dashboard.go) | `...string` | `dashboardids` IDs |
| [`dashboard.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dashboard/update) | [`Dashboard.Update`](../dashboard.go) | `...DashboardUpdateParams` | `dashboardids` IDs |
| [`discoveryrule.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/discoveryrule/create) | [`DiscoveryRule.Create`](../discoveryrule.go) | `...DiscoveryRuleCreateParams` | `itemids` IDs |
| [`discoveryrule.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/discoveryrule/delete) | [`DiscoveryRule.Delete`](../discoveryrule.go) | `...string` | `ruleids` IDs |
| [`discoveryrule.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/discoveryrule/update) | [`DiscoveryRule.Update`](../discoveryrule.go) | `...DiscoveryRuleUpdateParams` | `itemids` IDs |
| [`drule.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/drule/create) | [`DRule.Create`](../drule.go) | `...DRuleCreateParams` | `druleids` IDs |
| [`drule.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/drule/delete) | [`DRule.Delete`](../drule.go) | `...string` | `druleids` IDs |
| [`drule.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/drule/update) | [`DRule.Update`](../drule.go) | `...DRuleUpdateParams` | `druleids` IDs |
| [`event.acknowledge`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/event/acknowledge) | [`Event.Acknowledge`](../event.go) | `EventAcknowledgeParams` | `eventids` IDs |
| [`graph.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graph/create) | [`Graph.Create`](../graph.go) | `...GraphCreateParams` | `graphids` IDs |
| [`graph.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graph/delete) | [`Graph.Delete`](../graph.go) | `...string` | `graphids` IDs |
| [`graph.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graph/update) | [`Graph.Update`](../graph.go) | `...GraphUpdateParams` | `graphids` IDs |
| [`graphprototype.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graphprototype/create) | [`GraphPrototype.Create`](../graphprototype.go) | `...GraphPrototypeCreateParams` | `graphids` IDs |
| [`graphprototype.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graphprototype/delete) | [`GraphPrototype.Delete`](../graphprototype.go) | `...string` | `graphids` IDs |
| [`graphprototype.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graphprototype/update) | [`GraphPrototype.Update`](../graphprototype.go) | `...GraphPrototypeUpdateParams` | `graphids` IDs |
| [`host.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/host/create) | [`Host.Create`](../host.go) | `...HostCreateParams` | `hostids` IDs |
| [`host.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/host/delete) | [`Host.Delete`](../host.go) | `...string` | `hostids` IDs |
| [`host.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/host/update) | [`Host.Update`](../host.go) | `...HostUpdateParams` | `hostids` IDs |
| [`hostgroup.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostgroup/create) | [`HostGroup.Create`](../hostgroup.go) | `...HostGroupCreateParams` | `groupids` IDs |
| [`hostgroup.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostgroup/delete) | [`HostGroup.Delete`](../hostgroup.go) | `...string` | `groupids` IDs |
| [`hostgroup.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostgroup/update) | [`HostGroup.Update`](../hostgroup.go) | `...HostGroupUpdateParams` | `groupids` IDs |
| [`hostinterface.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/create) | [`HostInterface.Create`](../hostinterface.go) | `...HostInterfaceCreateParams` | `interfaceids` IDs |
| [`hostinterface.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/delete) | [`HostInterface.Delete`](../hostinterface.go) | `...string` | `interfaceids` IDs |
| [`hostinterface.massadd`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/massadd) | [`HostInterface.MassAdd`](../hostinterface.go) | `HostInterfaceMassAddParams` | `interfaceids` IDs |
| [`hostinterface.massremove`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/massremove) | [`HostInterface.MassRemove`](../hostinterface.go) | `HostInterfaceMassRemoveParams` | `interfaceids` IDs |
| [`hostinterface.replacehostinterfaces`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/replacehostinterfaces) | [`HostInterface.ReplaceHostInterfaces`](../hostinterface.go) | `HostInterfaceReplaceParams` | `interfaceids` IDs |
| [`hostinterface.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/update) | [`HostInterface.Update`](../hostinterface.go) | `...HostInterfaceUpdateParams` | `interfaceids` IDs |
| [`hostprototype.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostprototype/create) | [`HostPrototype.Create`](../hostprototype.go) | `...HostPrototypeCreateParams` | `hostids` IDs |
| [`hostprototype.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostprototype/delete) | [`HostPrototype.Delete`](../hostprototype.go) | `...string` | `hostids` IDs |
| [`hostprototype.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostprototype/update) | [`HostPrototype.Update`](../hostprototype.go) | `...HostPrototypeUpdateParams` | `hostids` IDs |
| [`housekeeping.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/housekeeping/update) | [`Housekeeping.Update`](../housekeeping.go) | `HousekeepingUpdateParams` | updated field names |
| [`iconmap.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/iconmap/create) | [`IconMap.Create`](../iconmap.go) | `...IconMapCreateParams` | `iconmapids` IDs |
| [`iconmap.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/iconmap/delete) | [`IconMap.Delete`](../iconmap.go) | `...string` | `iconmapids` IDs |
| [`iconmap.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/iconmap/update) | [`IconMap.Update`](../iconmap.go) | `...IconMapUpdateParams` | `iconmapids` IDs |
| [`image.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/image/create) | [`Image.Create`](../image.go) | `...ImageCreateParams` | `imageids` IDs |
| [`image.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/image/delete) | [`Image.Delete`](../image.go) | `...string` | `imageids` IDs |
| [`image.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/image/update) | [`Image.Update`](../image.go) | `...ImageUpdateParams` | `imageids` IDs |
| [`item.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/item/create) | [`Item.Create`](../item.go) | `...ItemCreateParams` | `itemids` IDs |
| [`item.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/item/delete) | [`Item.Delete`](../item.go) | `...string` | `itemids` IDs |
| [`item.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/item/update) | [`Item.Update`](../item.go) | `...ItemUpdateParams` | `itemids` IDs |
| [`itemprototype.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/itemprototype/create) | [`ItemPrototype.Create`](../itemprototype.go) | `...ItemPrototypeCreateParams` | `itemids` IDs |
| [`itemprototype.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/itemprototype/delete) | [`ItemPrototype.Delete`](../itemprototype.go) | `...string` | `prototypeids` IDs |
| [`itemprototype.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/itemprototype/update) | [`ItemPrototype.Update`](../itemprototype.go) | `...ItemPrototypeUpdateParams` | `itemids` IDs |
| [`maintenance.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/maintenance/create) | [`Maintenance.Create`](../maintenance.go) | `...MaintenanceCreateParams` | `maintenanceids` IDs |
| [`maintenance.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/maintenance/delete) | [`Maintenance.Delete`](../maintenance.go) | `...string` | `maintenanceids` IDs |
| [`maintenance.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/maintenance/update) | [`Maintenance.Update`](../maintenance.go) | `...MaintenanceUpdateParams` | `maintenanceids` IDs |
| [`map.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/map/create) | [`Map.Create`](../map.go) | `...MapCreateParams` | `sysmapids` IDs |
| [`map.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/map/delete) | [`Map.Delete`](../map.go) | `...string` | `sysmapids` IDs |
| [`map.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/map/update) | [`Map.Update`](../map.go) | `...MapUpdateParams` | `sysmapids` IDs |
| [`mediatype.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mediatype/create) | [`MediaType.Create`](../mediatype.go) | `...MediaTypeCreateParams` | `mediatypeids` IDs |
| [`mediatype.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mediatype/delete) | [`MediaType.Delete`](../mediatype.go) | `...string` | `mediatypeids` IDs |
| [`mediatype.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mediatype/update) | [`MediaType.Update`](../mediatype.go) | `...MediaTypeUpdateParams` | `mediatypeids` IDs |
| [`mfa.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mfa/create) | [`MFA.Create`](../mfa.go) | `...MFACreateParams` | `mfaids` IDs |
| [`mfa.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mfa/delete) | [`MFA.Delete`](../mfa.go) | `...string` | `mfaids` IDs |
| [`mfa.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mfa/update) | [`MFA.Update`](../mfa.go) | `...MFAUpdateParams` | `mfaids` IDs |
| [`module.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/module/create) | [`Module.Create`](../module.go) | `...ModuleCreateParams` | `moduleids` IDs |
| [`module.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/module/delete) | [`Module.Delete`](../module.go) | `...string` | `moduleids` IDs |
| [`module.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/module/update) | [`Module.Update`](../module.go) | `...ModuleUpdateParams` | `moduleids` IDs |
| [`proxy.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxy/create) | [`Proxy.Create`](../proxy.go) | `...ProxyCreateParams` | `proxyids` IDs |
| [`proxy.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxy/delete) | [`Proxy.Delete`](../proxy.go) | `...string` | `proxyids` IDs |
| [`proxy.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxy/update) | [`Proxy.Update`](../proxy.go) | `...ProxyUpdateParams` | `proxyids` IDs |
| [`proxygroup.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxygroup/create) | [`ProxyGroup.Create`](../proxygroup.go) | `...ProxyGroupCreateParams` | `proxy_groupids` IDs |
| [`proxygroup.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxygroup/delete) | [`ProxyGroup.Delete`](../proxygroup.go) | `...string` | `proxy_groupids` IDs |
| [`proxygroup.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxygroup/update) | [`ProxyGroup.Update`](../proxygroup.go) | `...ProxyGroupUpdateParams` | `proxy_groupids` IDs |
| [`regexp.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/regexp/create) | [`Regexp.Create`](../regexp.go) | `...RegexpCreateParams` | `regexpids` IDs |
| [`regexp.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/regexp/delete) | [`Regexp.Delete`](../regexp.go) | `...string` | `regexpids` IDs |
| [`regexp.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/regexp/update) | [`Regexp.Update`](../regexp.go) | `...RegexpUpdateParams` | `regexpids` IDs |
| [`report.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/report/create) | [`Report.Create`](../report.go) | `...ReportCreateParams` | `reportids` IDs |
| [`report.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/report/delete) | [`Report.Delete`](../report.go) | `...string` | `reportids` IDs |
| [`report.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/report/update) | [`Report.Update`](../report.go) | `...ReportUpdateParams` | `reportids` IDs |
| [`role.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/role/create) | [`Role.Create`](../role.go) | `...RoleCreateParams` | `roleids` IDs |
| [`role.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/role/delete) | [`Role.Delete`](../role.go) | `...string` | `roleids` IDs |
| [`role.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/role/update) | [`Role.Update`](../role.go) | `...RoleUpdateParams` | `roleids` IDs |
| [`script.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/create) | [`Script.Create`](../script.go) | `...ScriptCreateParams` | `scriptids` IDs |
| [`script.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/delete) | [`Script.Delete`](../script.go) | `...string` | `scriptids` IDs |
| [`script.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/update) | [`Script.Update`](../script.go) | `...ScriptUpdateParams` | `scriptids` IDs |
| [`service.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/service/create) | [`Service.Create`](../service.go) | `...ServiceCreateParams` | `serviceids` IDs |
| [`service.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/service/delete) | [`Service.Delete`](../service.go) | `...string` | `serviceids` IDs |
| [`service.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/service/update) | [`Service.Update`](../service.go) | `...ServiceUpdateParams` | `serviceids` IDs |
| [`settings.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/settings/update) | [`Settings.Update`](../settings.go) | `SettingsUpdateParams` | updated field names |
| [`sla.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/sla/create) | [`SLA.Create`](../sla.go) | `...SLACreateParams` | `slaids` IDs |
| [`sla.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/sla/delete) | [`SLA.Delete`](../sla.go) | `...string` | `slaids` IDs |
| [`sla.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/sla/update) | [`SLA.Update`](../sla.go) | `...SLAUpdateParams` | `slaids` IDs |
| [`task.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/task/create) | [`Task.Create`](../task.go) | `...TaskCreateParams` | `taskids` IDs |
| [`template.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/template/create) | [`Template.Create`](../template.go) | `...TemplateCreateParams` | `templateids` IDs |
| [`template.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/template/delete) | [`Template.Delete`](../template.go) | `...string` | `templateids` IDs |
| [`template.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/template/update) | [`Template.Update`](../template.go) | `...TemplateUpdateParams` | `templateids` IDs |
| [`templatedashboard.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templatedashboard/create) | [`TemplateDashboard.Create`](../templatedashboard.go) | `...TemplateDashboardCreateParams` | `dashboardids` IDs |
| [`templatedashboard.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templatedashboard/delete) | [`TemplateDashboard.Delete`](../templatedashboard.go) | `...string` | `dashboardids` IDs |
| [`templatedashboard.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templatedashboard/update) | [`TemplateDashboard.Update`](../templatedashboard.go) | `...TemplateDashboardUpdateParams` | `dashboardids` IDs |
| [`templategroup.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templategroup/create) | [`TemplateGroup.Create`](../templategroup.go) | `...TemplateGroupCreateParams` | `groupids` IDs |
| [`templategroup.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templategroup/delete) | [`TemplateGroup.Delete`](../templategroup.go) | `...string` | `groupids` IDs |
| [`templategroup.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templategroup/update) | [`TemplateGroup.Update`](../templategroup.go) | `...TemplateGroupUpdateParams` | `groupids` IDs |
| [`token.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/token/create) | [`APIToken.Create`](../token.go) | `...TokenCreateParams` | `tokenids` IDs |
| [`token.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/token/delete) | [`APIToken.Delete`](../token.go) | `...string` | `tokenids` IDs |
| [`token.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/token/update) | [`APIToken.Update`](../token.go) | `...TokenUpdateParams` | `tokenids` IDs |
| [`trigger.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/trigger/create) | [`Trigger.Create`](../trigger.go) | `...TriggerCreateParams` | `triggerids` IDs |
| [`trigger.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/trigger/delete) | [`Trigger.Delete`](../trigger.go) | `...string` | `triggerids` IDs |
| [`trigger.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/trigger/update) | [`Trigger.Update`](../trigger.go) | `...TriggerUpdateParams` | `triggerids` IDs |
| [`triggerprototype.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/triggerprototype/create) | [`TriggerPrototype.Create`](../triggerprototype.go) | `...TriggerPrototypeCreateParams` | `triggerids` IDs |
| [`triggerprototype.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/triggerprototype/delete) | [`TriggerPrototype.Delete`](../triggerprototype.go) | `...string` | `triggerids` IDs |
| [`triggerprototype.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/triggerprototype/update) | [`TriggerPrototype.Update`](../triggerprototype.go) | `...TriggerPrototypeUpdateParams` | `triggerids` IDs |
| [`user.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/user/create) | [`User.Create`](../user.go) | `...UserCreateParams` | `userids` IDs |
| [`user.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/user/delete) | [`User.Delete`](../user.go) | `...string` | `userids` IDs |
| [`user.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/user/update) | [`User.Update`](../user.go) | `...UserUpdateParams` | `userids` IDs |
| [`userdirectory.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/userdirectory/create) | [`UserDirectory.Create`](../userdirectory.go) | `...UserDirectoryCreateParams` | `userdirectoryids` IDs |
| [`userdirectory.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/userdirectory/delete) | [`UserDirectory.Delete`](../userdirectory.go) | `...string` | `userdirectoryids` IDs |
| [`userdirectory.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/userdirectory/update) | [`UserDirectory.Update`](../userdirectory.go) | `...UserDirectoryUpdateParams` | `userdirectoryids` IDs |
| [`usergroup.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usergroup/create) | [`UserGroup.Create`](../usergroup.go) | `...UserGroupCreateParams` | `usrgrpids` IDs |
| [`usergroup.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usergroup/delete) | [`UserGroup.Delete`](../usergroup.go) | `...string` | `usrgrpids` IDs |
| [`usergroup.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usergroup/update) | [`UserGroup.Update`](../usergroup.go) | `...UserGroupUpdateParams` | `usrgrpids` IDs |
| [`usermacro.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/create) | [`UserMacro.Create`](../usermacro.go) | `...HostMacroCreateParams` | `hostmacroids` IDs |
| [`usermacro.createglobal`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/createglobal) | [`UserMacro.CreateGlobal`](../usermacro.go) | `...GlobalMacroCreateParams` | `globalmacroids` IDs |
| [`usermacro.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/delete) | [`UserMacro.Delete`](../usermacro.go) | `...string` | `hostmacroids` IDs |
| [`usermacro.deleteglobal`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/deleteglobal) | [`UserMacro.DeleteGlobal`](../usermacro.go) | `...string` | `globalmacroids` IDs |
| [`usermacro.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/update) | [`UserMacro.Update`](../usermacro.go) | `...HostMacroUpdateParams` | `hostmacroids` IDs |
| [`usermacro.updateglobal`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/updateglobal) | [`UserMacro.UpdateGlobal`](../usermacro.go) | `...GlobalMacroUpdateParams` | `globalmacroids` IDs |
| [`valuemap.create`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/valuemap/create) | [`ValueMap.Create`](../valuemap.go) | `...ValueMapCreateParams` | `valuemapids` IDs |
| [`valuemap.delete`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/valuemap/delete) | [`ValueMap.Delete`](../valuemap.go) | `...string` | `valuemapids` IDs |
| [`valuemap.update`](https://www.zabbix.com/documentation/7.4/en/manual/api/reference/valuemap/update) | [`ValueMap.Update`](../valuemap.go) | `...ValueMapUpdateParams` | `valuemapids` IDs |

## Protocol details covered by regressions

- `itemprototype.delete` returns `prototypeids`; `discoveryrule.delete` returns `ruleids`.
- `hostinterface.massadd` in 7.4.16 nests `interfaceids` one extra time. Only that operation accepts the known nested shape, alongside the documented flat array.
- `hostinterface.massremove` with no match and `replacehostinterfaces` with an empty replacement return a legitimate empty ID list.
- `event.acknowledge` can return integer IDs after PHP array-key normalization. They are converted to decimal strings without a float intermediate. Explicit severity `0` and `suppress_until: 0` are preserved.
- Host macros contain their owning `HostID`. Nested host macros and global macros have separate inputs.
- User-group host and template permissions are independent `hostgroup_rights` and `templategroup_rights` collections.

The HTTP fixtures verify outgoing payloads, zero/empty distinctions, malformed
responses and exact IDs. Docker integration tests create their own entities and
verify persisted updates, selected associations and returned deletion IDs. They
exercise representative mutations; they do not claim a successful live write
for every enum combination or every namespace.
