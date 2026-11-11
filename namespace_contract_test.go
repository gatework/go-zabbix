package zabbix

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// These fixtures name the protocol result key independently of the Go field.
// A selector that is accepted but silently discarded must fail this test.
func TestEveryExposedSelectorHasAResultField(t *testing.T) {
	tests := []struct {
		method, selector, field string
		count                   bool
		params, model           any
		value                   string
	}{
		{"action.get", "selectOperations", "operations", false, ActionGetParams{}, Action{}, `[{"operationid":"1"}]`},
		{"action.get", "selectRecoveryOperations", "recovery_operations", false, ActionGetParams{}, Action{}, `[{"operationid":"2"}]`},
		{"action.get", "selectUpdateOperations", "update_operations", false, ActionGetParams{}, Action{}, `[{"operationid":"3"}]`},
		{"action.get", "selectFilter", "filter", false, ActionGetParams{}, Action{}, `{"evaltype":"0","eval_formula":"A"}`},
		{"alert.get", "selectHosts", "hosts", false, AlertGetParams{}, Alert{}, `[{"hostid":"1"}]`},
		{"alert.get", "selectMediatypes", "mediatypes", false, AlertGetParams{}, Alert{}, `[{"mediatypeid":"2"}]`},
		{"alert.get", "selectUsers", "users", false, AlertGetParams{}, Alert{}, `[{"userid":"3"}]`},
		{"connector.get", "selectTags", "tags", true, ConnectorGetParams{}, Connector{}, `[{"tag":"env","operator":"0","value":"prod"}]`},
		{"correlation.get", "selectFilter", "filter", false, CorrelationGetParams{}, Correlation{}, `{"evaltype":"0"}`},
		{"correlation.get", "selectOperations", "operations", false, CorrelationGetParams{}, Correlation{}, `[{"type":"1"}]`},
		{"dashboard.get", "selectPages", "pages", false, DashboardGetParams{}, Dashboard{}, `[{"dashboard_pageid":"1"}]`},
		{"dashboard.get", "selectUsers", "users", false, DashboardGetParams{}, Dashboard{}, `[{"userid":"2"}]`},
		{"dashboard.get", "selectUserGroups", "userGroups", false, DashboardGetParams{}, Dashboard{}, `[{"usrgrpid":"3"}]`},
		{"dhost.get", "selectDServices", "dservices", true, DHostGetParams{}, DHost{}, `[{"dserviceid":"1"}]`},
		{"discoveryrule.get", "selectHosts", "hosts", false, DiscoveryRuleGetParams{}, DiscoveryRule{}, `[{"hostid":"1"}]`},
		{"discoveryrule.get", "selectFilter", "filter", false, DiscoveryRuleGetParams{}, DiscoveryRule{}, `{"evaltype":"0","conditions":[{"macro":"{#NAME}","operator":"8","value":".*"}]}`},
		{"discoveryrule.get", "selectLLDMacroPaths", "lld_macro_paths", false, DiscoveryRuleGetParams{}, DiscoveryRule{}, `[{"lld_macro":"{#NAME}","path":"$.name"}]`},
		{"discoveryrule.get", "selectPreprocessing", "preprocessing", false, DiscoveryRuleGetParams{}, DiscoveryRule{}, `[{"type":"1","params":"2"}]`},
		{"drule.get", "selectDChecks", "dchecks", true, DRuleGetParams{}, DRule{}, `[{"dcheckid":"1"}]`},
		{"drule.get", "selectDHosts", "dhosts", true, DRuleGetParams{}, DRule{}, `[{"dhostid":"2"}]`},
		{"dservice.get", "selectDRules", "drules", false, DServiceGetParams{}, DService{}, `[{"druleid":"1"}]`},
		{"dservice.get", "selectDHosts", "dhosts", false, DServiceGetParams{}, DService{}, `[{"dhostid":"2"}]`},
		{"dservice.get", "selectHosts", "hosts", true, DServiceGetParams{}, DService{}, `[{"hostid":"3"}]`},
		{"event.get", "selectHosts", "hosts", false, EventGetParams{}, Event{}, `[{"hostid":"1"}]`},
		{"event.get", "selectTags", "tags", false, EventGetParams{}, Event{}, `[{"tag":"env","value":"prod"}]`},
		{"event.get", "selectAcknowledges", "acknowledges", true, EventGetParams{}, Event{}, `[{"acknowledgeid":"1","action":"4","message":"investigating","username":"Admin"}]`},
		{"graph.get", "selectGraphItems", "gitems", false, GraphGetParams{}, Graph{}, `[{"gitemid":"1"}]`},
		{"graph.get", "selectHosts", "hosts", false, GraphGetParams{}, Graph{}, `[{"hostid":"2"}]`},
		{"graph.get", "selectItems", "items", false, GraphGetParams{}, Graph{}, `[{"itemid":"3"}]`},
		{"graph.get", "selectTemplates", "templates", false, GraphGetParams{}, Graph{}, `[{"templateid":"4"}]`},
		{"graphitem.get", "selectGraphs", "graphs", false, GraphItemGetParams{}, GraphItem{}, `[{"graphid":"1"}]`},
		{"graphprototype.get", "selectGraphItems", "gitems", false, GraphPrototypeGetParams{}, GraphPrototype{}, `[{"gitemid":"1"}]`},
		{"graphprototype.get", "selectHosts", "hosts", false, GraphPrototypeGetParams{}, GraphPrototype{}, `[{"hostid":"2"}]`},
		{"graphprototype.get", "selectItems", "items", false, GraphPrototypeGetParams{}, GraphPrototype{}, `[{"itemid":"3"}]`},
		{"graphprototype.get", "selectDiscoveryRule", "discoveryRule", false, GraphPrototypeGetParams{}, GraphPrototype{}, `{"itemid":"4"}`},
		{"host.get", "selectHostGroups", "hostgroups", false, HostGetParams{}, Host{}, `[{"groupid":"1"}]`},
		{"host.get", "selectInterfaces", "interfaces", true, HostGetParams{}, Host{}, `[{"interfaceid":"2","details":{"version":"3"}}]`},
		{"host.get", "selectItems", "items", true, HostGetParams{}, Host{}, `[{"itemid":"3"}]`},
		{"host.get", "selectTriggers", "triggers", true, HostGetParams{}, Host{}, `[{"triggerid":"4"}]`},
		{"host.get", "selectParentTemplates", "parentTemplates", true, HostGetParams{}, Host{}, `[{"templateid":"5","link_type":"0"}]`},
		{"host.get", "selectTags", "tags", false, HostGetParams{}, Host{}, `[{"tag":"env","value":"prod"}]`},
		{"host.get", "selectMacros", "macros", false, HostGetParams{}, Host{}, `[{"hostmacroid":"6","hostid":"7","macro":"{$KEY}"}]`},
		{"host.get", "selectInventory", "inventory", false, HostGetParams{}, Host{}, `{"name":"node","os":"Linux"}`},
		{"hostgroup.get", "selectHosts", "hosts", true, HostGroupGetParams{}, HostGroup{}, `[{"hostid":"1"}]`},
		{"hostinterface.get", "selectItems", "items", true, HostInterfaceGetParams{}, HostInterface{}, `[{"itemid":"1"}]`},
		{"hostinterface.get", "selectHosts", "hosts", false, HostInterfaceGetParams{}, HostInterface{}, `[{"hostid":"2"}]`},
		{"hostprototype.get", "selectGroupLinks", "groupLinks", false, HostPrototypeGetParams{}, HostPrototype{}, `[{"groupid":"1"}]`},
		{"hostprototype.get", "selectGroupPrototypes", "groupPrototypes", false, HostPrototypeGetParams{}, HostPrototype{}, `[{"name":"group"}]`},
		{"hostprototype.get", "selectDiscoveryRule", "discoveryRule", false, HostPrototypeGetParams{}, HostPrototype{}, `{"itemid":"2"}`},
		{"hostprototype.get", "selectParentHost", "parentHost", false, HostPrototypeGetParams{}, HostPrototype{}, `{"hostid":"3"}`},
		{"hostprototype.get", "selectInterfaces", "interfaces", false, HostPrototypeGetParams{}, HostPrototype{}, `[{"interfaceid":"4"}]`},
		{"hostprototype.get", "selectTemplates", "templates", true, HostPrototypeGetParams{}, HostPrototype{}, `[{"templateid":"5"}]`},
		{"hostprototype.get", "selectTags", "tags", false, HostPrototypeGetParams{}, HostPrototype{}, `[{"tag":"env","value":"prod"}]`},
		{"hostprototype.get", "selectMacros", "macros", false, HostPrototypeGetParams{}, HostPrototype{}, `[{"hostmacroid":"6"}]`},
		{"iconmap.get", "selectMappings", "mappings", false, IconMapGetParams{}, IconMap{}, `[{"iconmappingid":"1"}]`},
		{"item.get", "selectHosts", "hosts", false, ItemGetParams{}, Item{}, `[{"hostid":"1"}]`},
		{"item.get", "selectTriggers", "triggers", true, ItemGetParams{}, Item{}, `[{"triggerid":"2"}]`},
		{"item.get", "selectTags", "tags", false, ItemGetParams{}, Item{}, `[{"tag":"env","value":"prod"}]`},
		{"item.get", "selectPreprocessing", "preprocessing", false, ItemGetParams{}, Item{}, `[{"type":"12","params":"$.value","error_handler":"2","error_handler_params":"fallback"}]`},
		{"itemprototype.get", "selectHosts", "hosts", false, ItemPrototypeGetParams{}, ItemPrototype{}, `[{"hostid":"1"}]`},
		{"itemprototype.get", "selectTags", "tags", false, ItemPrototypeGetParams{}, ItemPrototype{}, `[{"tag":"env","value":"prod"}]`},
		{"itemprototype.get", "selectPreprocessing", "preprocessing", false, ItemPrototypeGetParams{}, ItemPrototype{}, `[{"type":"1","params":"2"}]`},
		{"itemprototype.get", "selectDiscoveryRule", "discoveryRule", false, ItemPrototypeGetParams{}, ItemPrototype{}, `{"itemid":"3"}`},
		{"maintenance.get", "selectTimeperiods", "timeperiods", false, MaintenanceGetParams{}, Maintenance{}, `[{"timeperiodid":"1"}]`},
		{"maintenance.get", "selectTags", "tags", false, MaintenanceGetParams{}, Maintenance{}, `[{"tag":"env","operator":"0","value":"prod"}]`},
		{"maintenance.get", "selectHostGroups", "hostgroups", false, MaintenanceGetParams{}, Maintenance{}, `[{"groupid":"2"}]`},
		{"maintenance.get", "selectHosts", "hosts", false, MaintenanceGetParams{}, Maintenance{}, `[{"hostid":"3"}]`},
		{"map.get", "selectSelements", "selements", false, MapGetParams{}, Map{}, `[{"selementid":"1"}]`},
		{"map.get", "selectLinks", "links", false, MapGetParams{}, Map{}, `[{"linkid":"2"}]`},
		{"map.get", "selectUrls", "urls", false, MapGetParams{}, Map{}, `[{"sysmapurlid":"3","url":"https://example.test"}]`},
		{"map.get", "selectUsers", "users", false, MapGetParams{}, Map{}, `[{"userid":"4"}]`},
		{"map.get", "selectUserGroups", "userGroups", false, MapGetParams{}, Map{}, `[{"usrgrpid":"5"}]`},
		{"map.get", "selectShapes", "shapes", false, MapGetParams{}, Map{}, `[{"sysmap_shapeid":"6","text":"note"}]`},
		{"map.get", "selectLines", "lines", false, MapGetParams{}, Map{}, `[{"sysmap_shapeid":"7","x1":"2"}]`},
		{"mediatype.get", "selectMessageTemplates", "message_templates", false, MediaTypeGetParams{}, MediaType{}, `[{"eventsource":"0","subject":"subject"}]`},
		{"mediatype.get", "selectUsers", "users", false, MediaTypeGetParams{}, MediaType{}, `[{"userid":"1"}]`},
		{"problem.get", "selectAcknowledges", "acknowledges", true, ProblemGetParams{}, Problem{}, `[{"acknowledgeid":"1","suppress_until":"0"}]`},
		{"problem.get", "selectTags", "tags", false, ProblemGetParams{}, Problem{}, `[{"tag":"env","value":"prod"}]`},
		{"problem.get", "selectSuppressionData", "suppression_data", false, ProblemGetParams{}, Problem{}, `[{"maintenanceid":"2","userid":"3","suppress_until":"0"}]`},
		{"proxy.get", "selectHosts", "hosts", true, ProxyGetParams{}, Proxy{}, `[{"hostid":"1","description":"selected host"}]`},
		{"proxy.get", "selectProxyGroup", "proxyGroup", false, ProxyGetParams{}, Proxy{}, `{"proxy_groupid":"2"}`},
		{"proxygroup.get", "selectProxies", "proxies", true, ProxyGroupGetParams{}, ProxyGroup{}, `[{"proxyid":"1","address":"127.0.0.1"}]`},
		{"regexp.get", "selectExpressions", "expressions", false, RegexpGetParams{}, Regexp{}, `[{"expressionid":"1","expression":".*"}]`},
		{"report.get", "selectUsers", "users", false, ReportGetParams{}, Report{}, `[{"userid":"1"}]`},
		{"report.get", "selectUserGroups", "user_groups", false, ReportGetParams{}, Report{}, `[{"usrgrpid":"2"}]`},
		{"role.get", "selectRules", "rules", false, RoleGetParams{}, Role{}, `{"api":["host.get"]}`},
		{"script.get", "selectHosts", "hosts", false, ScriptGetParams{}, Script{}, `[{"hostid":"1"}]`},
		{"script.get", "selectHostGroups", "hostgroups", false, ScriptGetParams{}, Script{}, `[{"groupid":"2"}]`},
		{"service.get", "selectChildren", "children", true, ServiceGetParams{}, Service{}, `[{"serviceid":"1"}]`},
		{"service.get", "selectParents", "parents", true, ServiceGetParams{}, Service{}, `[{"serviceid":"2"}]`},
		{"service.get", "selectTags", "tags", true, ServiceGetParams{}, Service{}, `[{"tag":"env","value":"prod"}]`},
		{"service.get", "selectProblemTags", "problem_tags", true, ServiceGetParams{}, Service{}, `[{"tag":"env","operator":"0","value":"prod"}]`},
		{"service.get", "selectStatusRules", "status_rules", true, ServiceGetParams{}, Service{}, `[{"type":"0","limit_value":"1"}]`},
		{"sla.get", "selectSchedule", "schedule", true, SLAGetParams{}, SLA{}, `[{"period_from":"0","period_to":"86400"}]`},
		{"sla.get", "selectExcludedDowntimes", "excluded_downtimes", true, SLAGetParams{}, SLA{}, `[{"name":"planned","period_from":"0","period_to":"3600"}]`},
		{"sla.get", "selectServiceTags", "service_tags", true, SLAGetParams{}, SLA{}, `[{"tag":"env","operator":"0","value":"prod"}]`},
		{"template.get", "selectTemplateGroups", "templategroups", false, TemplateGetParams{}, Template{}, `[{"groupid":"1"}]`},
		{"template.get", "selectHosts", "hosts", true, TemplateGetParams{}, Template{}, `[{"hostid":"2"}]`},
		{"template.get", "selectItems", "items", true, TemplateGetParams{}, Template{}, `[{"itemid":"3"}]`},
		{"template.get", "selectTriggers", "triggers", true, TemplateGetParams{}, Template{}, `[{"triggerid":"4"}]`},
		{"template.get", "selectTags", "tags", false, TemplateGetParams{}, Template{}, `[{"tag":"env","value":"prod"}]`},
		{"template.get", "selectMacros", "macros", false, TemplateGetParams{}, Template{}, `[{"hostmacroid":"5"}]`},
		{"templatedashboard.get", "selectPages", "pages", false, TemplateDashboardGetParams{}, TemplateDashboard{}, `[{"dashboard_pageid":"1"}]`},
		{"templategroup.get", "selectTemplates", "templates", true, TemplateGroupGetParams{}, TemplateGroup{}, `[{"templateid":"1"}]`},
		{"trigger.get", "selectHosts", "hosts", false, TriggerGetParams{}, Trigger{}, `[{"hostid":"1"}]`},
		{"trigger.get", "selectItems", "items", false, TriggerGetParams{}, Trigger{}, `[{"itemid":"2"}]`},
		{"trigger.get", "selectTags", "tags", false, TriggerGetParams{}, Trigger{}, `[{"tag":"env","value":"prod"}]`},
		{"trigger.get", "selectDependencies", "dependencies", false, TriggerGetParams{}, Trigger{}, `[{"triggerid":"3"}]`},
		{"triggerprototype.get", "selectHosts", "hosts", false, TriggerPrototypeGetParams{}, TriggerPrototype{}, `[{"hostid":"1"}]`},
		{"triggerprototype.get", "selectItems", "items", false, TriggerPrototypeGetParams{}, TriggerPrototype{}, `[{"itemid":"2"}]`},
		{"triggerprototype.get", "selectTags", "tags", false, TriggerPrototypeGetParams{}, TriggerPrototype{}, `[{"tag":"env","value":"prod"}]`},
		{"triggerprototype.get", "selectDependencies", "dependencies", false, TriggerPrototypeGetParams{}, TriggerPrototype{}, `[{"triggerid":"3"}]`},
		{"user.get", "selectMedias", "medias", false, UserGetParams{}, User{}, `[{"mediaid":"1"}]`},
		{"user.get", "selectUsrgrps", "usrgrps", false, UserGetParams{}, User{}, `[{"usrgrpid":"2","name":"operators"}]`},
		{"user.get", "selectRole", "role", false, UserGetParams{}, User{}, `{"roleid":"3"}`},
		{"usergroup.get", "selectUsers", "users", false, UserGroupGetParams{}, UserGroup{}, `[{"userid":"1"}]`},
		{"usergroup.get", "selectTagFilters", "tag_filters", false, UserGroupGetParams{}, UserGroup{}, `[{"groupid":"2","tag":"env","value":"prod"}]`},
		{"usergroup.get", "selectHostGroupRights", "hostgroup_rights", false, UserGroupGetParams{}, UserGroup{}, `[{"id":"3","permission":"2"}]`},
		{"usergroup.get", "selectTemplateGroupRights", "templategroup_rights", false, UserGroupGetParams{}, UserGroup{}, `[{"id":"4","permission":"3"}]`},
		{"valuemap.get", "selectMappings", "mappings", true, ValueMapGetParams{}, ValueMap{}, `[{"type":"0","value":"0","newvalue":"down"}]`},
	}
	for _, tt := range tests {
		t.Run(tt.method+"/"+tt.selector, func(t *testing.T) {
			params := reflect.New(reflect.TypeOf(tt.params)).Elem()
			var selector reflect.Value
			for i := range params.NumField() {
				if strings.Split(params.Type().Field(i).Tag.Get("json"), ",")[0] == tt.selector {
					selector = params.Field(i)
					break
				}
			}
			if !selector.IsValid() {
				t.Fatal("missing public selector field")
			}
			selector.Set(reflect.ValueOf("extend"))
			args, err := queryArguments(tt.method, params.Interface())
			if err != nil {
				t.Fatal(err)
			}
			if len(args) != 1 || string(args[tt.selector]) != `"extend"` {
				t.Fatalf("arguments = %#v", args)
			}
			model := reflect.New(reflect.TypeOf(tt.model)).Interface()
			fixture := []byte(`{"` + tt.field + `":` + tt.value + `}`)
			if err := json.Unmarshal(fixture, model); err != nil {
				t.Fatal(err)
			}
			encoded, err := json.Marshal(model)
			if err != nil {
				t.Fatal(err)
			}
			var got map[string]json.RawMessage
			if err := json.Unmarshal(encoded, &got); err != nil {
				t.Fatal(err)
			}
			value, ok := got[tt.field]
			if !ok {
				t.Fatalf("selected %s was discarded: %s", tt.field, encoded)
			}
			assertQueryJSON(t, value, []byte(tt.value))
			selector.Set(reflect.ValueOf("count"))
			_, err = queryArguments(tt.method, params.Interface())
			if tt.count {
				if err != nil {
					t.Fatal(err)
				}
				model = reflect.New(reflect.TypeOf(tt.model)).Interface()
				if err := json.Unmarshal([]byte(`{"`+tt.field+`":"0"}`), model); err != nil {
					t.Fatal(err)
				}
				encoded, err = json.Marshal(model)
				if err != nil {
					t.Fatal(err)
				}
				assertQueryJSON(t, encoded, []byte(`{"`+tt.field+`":"0"}`))
			} else {
				assertQueryError(t, err, tt.method, ErrorInvalidParams)
			}
		})
	}
}
