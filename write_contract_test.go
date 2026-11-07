package zabbix

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func writeTestPointer[T any](value T) *T { return &value }

func writeJSON(t *testing.T, value any) map[string]json.RawMessage {
	t.Helper()
	body, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var object map[string]json.RawMessage
	if err := json.Unmarshal(body, &object); err != nil {
		t.Fatal(err)
	}
	return object
}

func TestWriteJSONFixtures(t *testing.T) {
	cases := []struct {
		name   string
		params interface{ validate() error }
		want   string
	}{
		{
			name: "create host keeps empty dns and macro value with exact reference ids",
			params: HostCreateParams{
				Host:       "router",
				Groups:     []HostGroupRef{{GroupID: "9007199254740993"}},
				Interfaces: []HostInterfaceInput{{Main: 1, Type: 1, UseIP: 1, IP: "192.0.2.1", DNS: "", Port: "10050"}},
				Macros:     []HostMacroInput{{Macro: "{$EMPTY}", Value: ""}},
			},
			want: `{"host":"router","groups":[{"groupid":"9007199254740993"}],"interfaces":[{"main":1,"type":1,"useip":1,"ip":"192.0.2.1","dns":"","port":"10050"}],"macros":[{"macro":"{$EMPTY}","value":""}]}`,
		},
		{
			name:   "host update clears text tags and templates while preserving absent fields",
			params: HostUpdateParams{HostID: "42", Description: writeTestPointer(""), Status: writeTestPointer(0), Tags: []TagInput{}, Templates: []TemplateRef{}},
			want:   `{"hostid":"42","description":"","status":0,"tags":[],"templates":[]}`,
		},
		{
			name:   "host update retains an existing proxy assignment",
			params: HostUpdateParams{HostID: "42", MonitoredBy: writeTestPointer(1), Name: writeTestPointer("router")},
			want:   `{"hostid":"42","monitored_by":1,"name":"router"}`,
		},
		{
			name:   "host update retains an existing proxy group assignment",
			params: HostUpdateParams{HostID: "42", MonitoredBy: writeTestPointer(2), Name: writeTestPointer("router")},
			want:   `{"hostid":"42","monitored_by":2,"name":"router"}`,
		},
		{
			name:   "service create allows a value with the default propagation rule",
			params: ServiceCreateParams{Name: "router", Algorithm: 1, PropagationValue: writeTestPointer(0)},
			want:   `{"name":"router","algorithm":1,"sortorder":0,"propagation_value":0}`,
		},
		{
			name:   "service update retains its existing propagation rule",
			params: ServiceUpdateParams{ServiceID: "42", PropagationValue: writeTestPointer(2)},
			want:   `{"serviceid":"42","propagation_value":2}`,
		},
		{
			name:   "user group uses separate host and template permissions",
			params: UserGroupUpdateParams{UsrGrpID: "8", HostGroupRights: []UserGroupRightInput{{ID: "3", Permission: 0}}, TemplateGroupRights: []UserGroupRightInput{}, Users: []UserRef{}},
			want:   `{"usrgrpid":"8","hostgroup_rights":[{"id":"3","permission":0}],"templategroup_rights":[],"users":[]}`,
		},
		{
			name:   "macro create owns host and permits empty value",
			params: HostMacroCreateParams{HostID: "42", Macro: "{$EMPTY}", Value: ""},
			want:   `{"hostid":"42","macro":"{$EMPTY}","value":""}`,
		},
		{
			name:   "user media clearing omits password and role",
			params: UserUpdateParams{UserID: "9", Medias: []UserMediaInput{}, URL: writeTestPointer("")},
			want:   `{"userid":"9","medias":[],"url":""}`,
		},
		{
			name:   "item create sends required zero numeric enums",
			params: ItemCreateParams{HostID: "42", Name: "temperature", Key: "sensor.temp", Type: 0, ValueType: 0, InterfaceID: writeTestPointer("7"), Delay: writeTestPointer("1m")},
			want:   `{"hostid":"42","name":"temperature","key_":"sensor.temp","type":0,"value_type":0,"interfaceid":"7","delay":"1m"}`,
		},
		{
			name:   "sla empty schedule explicitly requests continuous coverage",
			params: SLAUpdateParams{SLAID: "1", Schedule: []SLAScheduleInput{}, ExcludedDowntimes: []SLAExcludedDowntimeInput{}, SLO: writeTestPointer(0.0)},
			want:   `{"slaid":"1","schedule":[],"excluded_downtimes":[],"slo":0}`,
		},
		{
			name:   "interface update sends partial snmp details with zero bulk",
			params: HostInterfaceUpdateParams{InterfaceID: "5", Details: &SNMPInterfaceUpdateInput{Bulk: writeTestPointer(0)}},
			want:   `{"interfaceid":"5","details":{"bulk":0}}`,
		},
		{
			name:   "replacement removes all interfaces explicitly",
			params: HostInterfaceReplaceParams{HostID: "42", Interfaces: []HostInterfaceUpdateInput{}},
			want:   `{"hostid":"42","interfaces":[]}`,
		},
		{
			name:   "interface match includes both address properties",
			params: HostInterfaceMatchInput{IP: "192.0.2.1", DNS: "", Port: "161"},
			want:   `{"ip":"192.0.2.1","dns":"","port":"161"}`,
		},
		{
			name:   "discovery rule filter can remove every condition",
			params: DiscoveryRuleUpdateParams{ItemID: "42", Filter: &DiscoveryFilterInput{EvalType: 0, Conditions: []DiscoveryConditionInput{}}},
			want:   `{"itemid":"42","filter":{"evaltype":0,"conditions":[]}}`,
		},
		{
			name:   "action filter can remove every condition",
			params: ActionUpdateParams{ActionID: "42", Filter: &ActionFilterInput{EvalType: 0, Conditions: []ActionConditionInput{}}},
			want:   `{"actionid":"42","filter":{"evaltype":0,"conditions":[]}}`,
		},
		{
			name:   "action keeps enum zero and restricts nested message references",
			params: ActionCreateParams{Name: "Notify", EventSource: 0, Operations: []ActionOperationInput{{OperationType: 0, OpMessage: &ActionMessageInput{DefaultMsg: writeTestPointer(0), Message: writeTestPointer("")}, OpMessageUsr: []UserRef{{UserID: "4"}}}}},
			want:   `{"name":"Notify","eventsource":0,"operations":[{"operationtype":0,"opmessage":{"default_msg":0,"message":""},"opmessage_usr":[{"userid":"4"}]}]}`,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.params.validate(); err != nil {
				t.Fatal(err)
			}
			var got, want any
			if err := json.Unmarshal(mustWriteJSON(t, tc.params), &got); err != nil {
				t.Fatal(err)
			}
			if err := json.Unmarshal([]byte(tc.want), &want); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("wire request = %s, want %s", mustWriteJSON(t, got), tc.want)
			}
		})
	}
}

func mustWriteJSON(t *testing.T, value any) []byte {
	t.Helper()
	body, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	return body
}

// This contract is checked for every concrete mutation model and all its nested
// inputs. It protects clear/zero/omit semantics when new fields are added.
func TestEveryWriteFieldPreservesPresence(t *testing.T) {
	seen := make(map[reflect.Type]bool)
	var check func(reflect.Type)
	check = func(typ reflect.Type) {
		for typ.Kind() == reflect.Pointer || typ.Kind() == reflect.Slice {
			typ = typ.Elem()
		}
		if typ.Kind() != reflect.Struct || seen[typ] {
			return
		}
		seen[typ] = true
		base := writeJSON(t, reflect.New(typ).Elem().Interface())
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			key, _, _ := strings.Cut(field.Tag.Get("json"), ",")
			if key == "" || key == "-" {
				t.Fatalf("%s.%s has no wire property", typ, field.Name)
			}
			check(field.Type)
			switch field.Type.Kind() {
			case reflect.Pointer:
				if _, present := base[key]; present && strings.Contains(field.Tag.Get("json"), "omit") {
					t.Errorf("%s.%s: nil pointer was not omitted", typ, field.Name)
				}
				value := reflect.New(typ).Elem()
				value.Field(i).Set(reflect.New(field.Type.Elem()))
				if _, present := writeJSON(t, value.Interface())[key]; !present {
					t.Errorf("%s.%s: explicit zero was omitted", typ, field.Name)
				}
			case reflect.Slice:
				if strings.Contains(field.Tag.Get("json"), "omit") {
					if _, present := base[key]; present {
						t.Errorf("%s.%s: nil collection was not omitted", typ, field.Name)
					}
				}
				value := reflect.New(typ).Elem()
				value.Field(i).Set(reflect.MakeSlice(field.Type, 0, 0))
				if got := string(writeJSON(t, value.Interface())[key]); got != "[]" {
					t.Errorf("%s.%s: explicit empty = %q, want []", typ, field.Name, got)
				}
			case reflect.Interface:
				// Only protocol-selected values and module/task configuration are dynamic.
				allowed := map[string]bool{"UserMediaInput.SendTo": true, "WidgetFieldInput.Value": true, "TaskCreateParams.Request": true, "ModuleCreateParams.Config": true, "ModuleUpdateParams.Config": true}
				if !allowed[typ.Name()+"."+field.Name] {
					t.Errorf("unexpected dynamic field %s.%s", typ, field.Name)
				}
			default:
				if _, present := base[key]; !present {
					t.Errorf("%s.%s: required zero was omitted", typ, field.Name)
				}
			}
		}
	}
	for _, params := range allMutationParams() {
		check(reflect.TypeOf(params))
	}
}

func TestWriteModelsDoNotContainResponseObjects(t *testing.T) {
	seen := make(map[reflect.Type]bool)
	var check func(reflect.Type)
	check = func(typ reflect.Type) {
		for typ.Kind() == reflect.Pointer || typ.Kind() == reflect.Slice {
			typ = typ.Elem()
		}
		if typ.Kind() != reflect.Struct || seen[typ] {
			return
		}
		seen[typ] = true
		if !strings.HasSuffix(typ.Name(), "Params") && !strings.HasSuffix(typ.Name(), "Input") && !strings.HasSuffix(typ.Name(), "Ref") {
			t.Errorf("response model %s leaked into a write contract", typ)
		}
		for i := 0; i < typ.NumField(); i++ {
			check(typ.Field(i).Type)
		}
	}
	for _, params := range allMutationParams() {
		check(reflect.TypeOf(params))
	}
	for _, params := range allMutationParams() {
		typ := reflect.TypeOf(params)
		for _, forbidden := range []string{"Flags", "LastAccess", "LastClock", "LastNs", "LastValue", "PrevValue", "CreatedAt", "CreatorUserID", "AttemptClock", "AttemptFailed", "AttemptIP", "TsProvisioned", "Error", "Readonly", "ActiveAvailable", "MaintenanceID", "MaintenanceStatus", "MaintenanceType", "MaintenanceFrom"} {
			if (forbidden == "MaintenanceType" || forbidden == "MaintenanceID") && strings.HasPrefix(typ.Name(), "Maintenance") {
				continue
			}
			if _, found := typ.FieldByName(forbidden); found {
				t.Errorf("%s exposes read-only %s", typ, forbidden)
			}
		}
		if strings.HasSuffix(typ.Name(), "UpdateParams") {
			for _, constant := range map[string][]string{"ItemUpdateParams": {"HostID"}, "ItemPrototypeUpdateParams": {"HostID", "RuleID"}, "DiscoveryRuleUpdateParams": {"HostID"}, "HostInterfaceUpdateParams": {"HostID"}, "ValueMapUpdateParams": {"HostID"}, "HostMacroUpdateParams": {"HostID"}, "TokenUpdateParams": {"UserID"}, "ImageUpdateParams": {"ImageType"}, "ActionUpdateParams": {"EventSource"}, "TemplateDashboardUpdateParams": {"TemplateID"}}[typ.Name()] {
				if _, found := typ.FieldByName(constant); found {
					t.Errorf("%s exposes create-only %s", typ, constant)
				}
			}
		}
	}
}

func TestWriteValidationRejectsMissingIdentityAndConditionalPayloads(t *testing.T) {
	cases := []interface{ validate() error }{
		HostUpdateParams{}, ItemUpdateParams{ItemID: "-1"}, UserUpdateParams{UserID: "1.5"},
		HostCreateParams{Host: "r"}, HostCreateParams{Host: "r", Groups: []HostGroupRef{{}}},
		HostCreateParams{Host: "r", Groups: []HostGroupRef{{GroupID: "1"}}, MonitoredBy: writeTestPointer(1)},
		HostCreateParams{Host: "r", Groups: []HostGroupRef{{GroupID: "1"}}, MonitoredBy: writeTestPointer(2)},
		HostMacroCreateParams{Macro: "{$X}", Value: "x"},
		HostInterfaceCreateParams{HostID: "1", Type: 2, Main: 1, UseIP: 1, IP: "192.0.2.1", Port: "161"},
		HostInterfaceUpdateInput{}, SNMPInterfaceInput{Version: 2},
		HostInterfaceUpdateInput{InterfaceID: writeTestPointer("")},
		HostInterfaceReplaceParams{HostID: "1"}, HostInterfaceMatchInput{},
		ActionCreateParams{Name: "x", EventSource: 0},
		ActionOperationInput{OperationType: 0}, ActionOperationInput{OperationType: 1},
		ActionOperationInput{OperationType: 4}, ActionOperationInput{OperationType: 6},
		ActionOperationInput{OperationType: 10},
		UserMediaInput{MediaTypeID: "1", SendTo: 42},
		GraphItemUpdateInput{}, DashboardWidgetUpdateInput{},
		HostMacroUpdateInput{Macro: writeTestPointer("{$X}")},
		ServiceCreateParams{Name: "x", PropagationRule: writeTestPointer(0)},
		ServiceUpdateParams{ServiceID: "1", PropagationRule: writeTestPointer(0)},
		MaintenanceCreateParams{Name: "x", Timeperiods: []MaintenancePeriodInput{{}}},
		TaskCreateParams{Type: 7, Request: map[string]any{}},
	}
	for _, params := range cases {
		if err := params.validate(); err == nil {
			t.Errorf("%T accepted incomplete/invalid input: %#v", params, params)
		}
	}
	for _, params := range []interface{ validate() error }{
		HostUpdateParams{HostID: "9007199254740993"},
		ActionHostInput{HostID: "0"},
		MapLinkInput{SelementID1: "new-a", SelementID2: "new-b"},
		UserMediaInput{MediaTypeID: "1", SendTo: []string{"test@example.invalid"}},
		SNMPInterfaceUpdateInput{Bulk: writeTestPointer(0)},
	} {
		if err := params.validate(); err != nil {
			t.Errorf("%T rejected valid input: %v", params, err)
		}
	}
}

func TestEveryUpdateRequiresItsWireID(t *testing.T) {
	for _, params := range allMutationParams() {
		typ := reflect.TypeOf(params)
		if !strings.HasSuffix(typ.Name(), "UpdateParams") {
			continue
		}
		switch typ.Name() {
		case "SettingsUpdateParams", "AuthenticationUpdateParams", "AutoregistrationUpdateParams", "HousekeepingUpdateParams":
			continue
		}
		if err := params.validate(); err == nil {
			t.Errorf("%s accepted an update without its identity", typ)
		}
		identities := 0
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if field.Type.Kind() == reflect.String && !strings.Contains(field.Tag.Get("json"), ",") {
				identities++
			}
		}
		if identities != 1 {
			t.Errorf("%s has %d mandatory scalar identities, want one", typ, identities)
		}
	}
}

// Enumerate the public mutation family so newly added contracts join the common tests.
func allMutationParams() []interface{ validate() error } {
	return []interface{ validate() error }{
		ActionCreateParams{},
		ActionUpdateParams{},
		AuthenticationUpdateParams{},
		AutoregistrationUpdateParams{},
		ConnectorCreateParams{},
		ConnectorUpdateParams{},
		CorrelationCreateParams{},
		CorrelationUpdateParams{},
		DRuleCreateParams{},
		DRuleUpdateParams{},
		DashboardCreateParams{},
		DashboardUpdateParams{},
		DiscoveryRuleCreateParams{},
		DiscoveryRuleUpdateParams{},
		EventAcknowledgeParams{},
		GlobalMacroCreateParams{},
		GlobalMacroUpdateParams{},
		GraphCreateParams{},
		GraphPrototypeCreateParams{},
		GraphPrototypeUpdateParams{},
		GraphUpdateParams{},
		HostCreateParams{},
		HostGroupCreateParams{},
		HostGroupUpdateParams{},
		HostInterfaceCreateParams{},
		HostInterfaceMassAddParams{},
		HostInterfaceMassRemoveParams{},
		HostInterfaceReplaceParams{},
		HostInterfaceUpdateParams{},
		HostMacroCreateParams{},
		HostMacroUpdateParams{},
		HostPrototypeCreateParams{},
		HostPrototypeUpdateParams{},
		HostUpdateParams{},
		HousekeepingUpdateParams{},
		IconMapCreateParams{},
		IconMapUpdateParams{},
		ImageCreateParams{},
		ImageUpdateParams{},
		ItemCreateParams{},
		ItemPrototypeCreateParams{},
		ItemPrototypeUpdateParams{},
		ItemUpdateParams{},
		MFACreateParams{},
		MFAUpdateParams{},
		MaintenanceCreateParams{},
		MaintenanceUpdateParams{},
		MapCreateParams{},
		MapUpdateParams{},
		MediaTypeCreateParams{},
		MediaTypeUpdateParams{},
		ModuleCreateParams{},
		ModuleUpdateParams{},
		ProxyCreateParams{},
		ProxyGroupCreateParams{},
		ProxyGroupUpdateParams{},
		ProxyUpdateParams{},
		RegexpCreateParams{},
		RegexpUpdateParams{},
		ReportCreateParams{},
		ReportUpdateParams{},
		RoleCreateParams{},
		RoleUpdateParams{},
		SLACreateParams{},
		SLAUpdateParams{},
		ScriptCreateParams{},
		ScriptUpdateParams{},
		ServiceCreateParams{},
		ServiceUpdateParams{},
		SettingsUpdateParams{},
		TaskCreateParams{},
		TemplateCreateParams{},
		TemplateDashboardCreateParams{},
		TemplateDashboardUpdateParams{},
		TemplateGroupCreateParams{},
		TemplateGroupUpdateParams{},
		TemplateUpdateParams{},
		TokenCreateParams{},
		TokenUpdateParams{},
		TriggerCreateParams{},
		TriggerPrototypeCreateParams{},
		TriggerPrototypeUpdateParams{},
		TriggerUpdateParams{},
		UserCreateParams{},
		UserDirectoryCreateParams{},
		UserDirectoryUpdateParams{},
		UserGroupCreateParams{},
		UserGroupUpdateParams{},
		UserUpdateParams{},
		ValueMapCreateParams{},
		ValueMapUpdateParams{},
	}
}
