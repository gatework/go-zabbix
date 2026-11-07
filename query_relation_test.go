package zabbix

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestRelationPresenceAndAtomicDecoding(t *testing.T) {
	for _, tc := range []struct {
		body          string
		absent, empty bool
		count         *uint64
	}{
		{body: `{}`, absent: true},
		{body: `{"interfaces":[]}`, empty: true},
		{body: `{"interfaces":"0"}`, count: func() *uint64 { v := uint64(0); return &v }()},
		{body: `{"interfaces":2}`, count: func() *uint64 { v := uint64(2); return &v }()},
	} {
		t.Run(tc.body, func(t *testing.T) {
			var host Host
			if err := json.Unmarshal([]byte(tc.body), &host); err != nil {
				t.Fatal(err)
			}
			got := host.Interfaces
			if tc.absent && (got.Items != nil || got.Count != nil) {
				t.Fatalf("absent = %#v", got)
			}
			if tc.empty && (got.Items == nil || len(got.Items) != 0 || got.Count != nil) {
				t.Fatalf("empty = %#v", got)
			}
			if tc.count != nil && (got.Count == nil || *got.Count != *tc.count || got.Items != nil) {
				t.Fatalf("count = %#v", got)
			}
			if tc.absent {
				encoded, err := json.Marshal(host)
				if err != nil {
					t.Fatal(err)
				}
				assertQueryJSON(t, encoded, []byte(`{}`))
			}
		})
	}
	original := Relation[Host]{Items: []Host{{HostID: "keep"}}}
	for _, bad := range []string{`null`, `{}`, `true`, `[null]`, `[{"hostid":"new"},{"hostid":4}]`, `"-1"`, `1.5`, `""`, `"18446744073709551616"`} {
		got := original
		if err := json.Unmarshal([]byte(bad), &got); err == nil {
			t.Fatalf("accepted %s", bad)
		}
		if !reflect.DeepEqual(got, original) {
			t.Fatalf("mutated after %s: %#v", bad, got)
		}
	}
	zero := uint64(0)
	if _, err := json.Marshal(Relation[Host]{Items: []Host{}, Count: &zero}); err == nil {
		t.Fatal("ambiguous relation marshaled")
	}
}

func TestSelectedObjectsAndInventoryDecodeAtomically(t *testing.T) {
	var prototype HostPrototype
	if err := json.Unmarshal([]byte(`{"discoveryRule":[],"parentHost":{"hostid":"1"}}`), &prototype); err != nil {
		t.Fatal(err)
	}
	if !prototype.DiscoveryRule.Present || prototype.DiscoveryRule.Value != nil || !prototype.ParentHost.Present || prototype.ParentHost.Value.HostID != "1" {
		t.Fatalf("prototype = %#v", prototype)
	}
	original := SelectedObject[Host]{Present: true, Value: &Host{HostID: "keep"}}
	for _, bad := range []string{`null`, `"1"`, `0`, `[{}]`, `{"hostid":5}`} {
		got := original
		if err := json.Unmarshal([]byte(bad), &got); err == nil {
			t.Fatalf("accepted %s", bad)
		}
		if !reflect.DeepEqual(got, original) {
			t.Fatalf("mutated after %s", bad)
		}
	}
	inventory := HostInventory{"name": "keep"}
	for _, bad := range []string{`null`, `"value"`, `[{}]`, `{"name":null}`, `{"name":3}`} {
		if err := json.Unmarshal([]byte(bad), &inventory); err == nil {
			t.Fatalf("accepted %s", bad)
		}
		if inventory["name"] != "keep" || len(inventory) != 1 {
			t.Fatalf("mutated after %s", bad)
		}
	}
	if err := json.Unmarshal([]byte(`{}`), &inventory); err != nil || inventory == nil || len(inventory) != 0 {
		t.Fatalf("empty inventory = %#v, %v", inventory, err)
	}
}

func TestHostInterfaceDetailsRejectUnsupportedShapesAtomically(t *testing.T) {
	original := HostInterface{InterfaceID: "keep", Details: &InterfaceDetail{Version: "3"}}
	for _, bad := range []string{`null`, `true`, `0`, `"text"`, `[{}]`, `{"version":3}`} {
		got := original
		if err := json.Unmarshal([]byte(`{"interfaceid":"new","details":`+bad+`}`), &got); err == nil {
			t.Fatalf("accepted %s", bad)
		}
		if !reflect.DeepEqual(got, original) {
			t.Fatalf("mutated after %s", bad)
		}
	}
	for _, valid := range []string{`[]`, `{}`} {
		got := original
		if err := json.Unmarshal([]byte(`{"interfaceid":"new","details":`+valid+`}`), &got); err != nil {
			t.Fatal(err)
		}
		if got.InterfaceID != "new" {
			t.Fatalf("interface = %#v", got)
		}
		if valid == `[]` && got.Details != nil {
			t.Fatal("non-SNMP interface retained previous details")
		}
	}
}

func TestComputedReadFieldsRetainNumbersAndBooleans(t *testing.T) {
	var user User
	if err := json.Unmarshal([]byte(`{"gui_access":0,"debug_mode":"1","users_status":0,"role":{"roleid":"3"}}`), &user); err != nil {
		t.Fatal(err)
	}
	if user.GUIAccess.String() != "0" || user.DebugMode.String() != "1" || user.UsersStatus.String() != "0" || user.Role.Value.RoleID != "3" {
		t.Fatalf("user = %#v", user)
	}
	var role Role
	if err := json.Unmarshal([]byte(`{"rules":{"modules":[{"moduleid":43,"status":"1"}],"services.read.tag":{"tag":"env","value":"prod"},"api":["host.get"]}}`), &role); err != nil {
		t.Fatal(err)
	}
	if role.Rules.Value == nil || role.Rules.Value.Modules[0].ModuleID.String() != "43" || role.Rules.Value.ServicesReadTag.Value != "prod" {
		t.Fatalf("rules = %#v", role.Rules)
	}
	var service Service
	if err := json.Unmarshal([]byte(`{"serviceid":"4","readonly":true,"children":"0"}`), &service); err != nil {
		t.Fatal(err)
	}
	if !service.ReadOnly || service.Children.Count == nil || *service.Children.Count != 0 {
		t.Fatalf("service = %#v", service)
	}
}
