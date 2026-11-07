//go:build integration

package zabbix

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestIntegrationRoleServiceAccessTags(t *testing.T) {
	c, ctx := loginClient(t)
	name := fmt.Sprintf("go-zabbix-role-access-%d", time.Now().UnixNano())
	ids, err := c.Service().Create(ctx, ServiceCreateParams{Name: name + "-read", Algorithm: 0})
	readServiceID := integrationID(t, ids, err)
	integrationDelete(t, c.Service().Delete, readServiceID)
	ids, err = c.Service().Create(ctx, ServiceCreateParams{Name: name + "-write", Algorithm: 0})
	writeServiceID := integrationID(t, ids, err)
	integrationDelete(t, c.Service().Delete, writeServiceID)
	ids, err = c.Role().Create(ctx, RoleCreateParams{
		Name: name, Type: 1,
		Rules: &RoleRulesInput{
			ServicesReadMode:  integrationValue(0),
			ServicesReadList:  []ServiceRef{{ServiceID: readServiceID}},
			ServicesReadTag:   &ServiceAccessTagInput{Tag: "sdk-read", Value: integrationValue(name)},
			ServicesWriteMode: integrationValue(0),
			ServicesWriteList: []ServiceRef{{ServiceID: writeServiceID}},
			ServicesWriteTag:  &ServiceAccessTagInput{Tag: "sdk-write", Value: integrationValue(name)},
			APIAccess:         integrationValue(1),
			APIMode:           integrationValue(1),
			API:               []string{"service.get"},
		},
	})
	roleID := integrationID(t, ids, err)
	integrationDelete(t, c.Role().Delete, roleID)
	getRole := func(t *testing.T) Role {
		t.Helper()
		roles, err := c.Role().Get(ctx, RoleGetParams{RoleIDs: []string{roleID}, Output: "extend", SelectRules: "extend"})
		integrationCheck(t, err)
		if len(roles) != 1 || roles[0].RoleID != roleID || roles[0].Rules.Value == nil {
			t.Fatalf("selected role = %#v", roles)
		}
		return roles[0]
	}
	role := getRole(t)
	wantRules := *role.Rules.Value
	if role.Name != name || role.Type != "1" || wantRules.ServicesReadMode != "0" || wantRules.ServicesWriteMode != "0" {
		t.Fatalf("created role modes = %#v, rules = %#v", role, wantRules)
	}
	if len(wantRules.ServicesReadList) != 1 || wantRules.ServicesReadList[0].ServiceID != readServiceID ||
		len(wantRules.ServicesWriteList) != 1 || wantRules.ServicesWriteList[0].ServiceID != writeServiceID {
		t.Fatalf("created role service lists = %#v", wantRules)
	}
	if wantRules.ServicesReadTag == nil || *wantRules.ServicesReadTag != (ServiceAccessTag{Tag: "sdk-read", Value: name}) ||
		wantRules.ServicesWriteTag == nil || *wantRules.ServicesWriteTag != (ServiceAccessTag{Tag: "sdk-write", Value: name}) {
		t.Fatalf("created role service tags = %#v", wantRules)
	}
	if wantRules.APIAccess != "1" || wantRules.APIMode != "1" || !reflect.DeepEqual(wantRules.API, []string{"service.get"}) {
		t.Fatalf("created role API rules = %#v", wantRules)
	}

	name += "-renamed"
	ids, err = c.Role().Update(ctx, RoleUpdateParams{RoleID: roleID, Name: &name})
	if updatedID := integrationID(t, ids, err); updatedID != roleID {
		t.Fatalf("updated role ID = %s, want %s", updatedID, roleID)
	}
	role = getRole(t)
	if role.Name != name || !reflect.DeepEqual(*role.Rules.Value, wantRules) {
		t.Fatalf("omitting rules changed service access: role = %#v, rules = %#v", role, role.Rules.Value)
	}

	for _, tc := range []struct {
		name  string
		rules RoleRulesInput
	}{
		{"clear read tag", RoleRulesInput{ServicesReadTag: &ServiceAccessTagInput{Tag: ""}}},
		{"clear write tag", RoleRulesInput{ServicesWriteTag: &ServiceAccessTagInput{Tag: "", Value: integrationValue("")}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ids, err := c.Role().Update(ctx, RoleUpdateParams{RoleID: roleID, Rules: &tc.rules})
			if updatedID := integrationID(t, ids, err); updatedID != roleID {
				t.Fatalf("updated role ID = %s, want %s", updatedID, roleID)
			}
			if tc.rules.ServicesReadTag != nil {
				wantRules.ServicesReadTag = &ServiceAccessTag{}
			} else {
				wantRules.ServicesWriteTag = &ServiceAccessTag{}
			}
			role := getRole(t)
			if role.Name != name || role.Type != "1" || !reflect.DeepEqual(*role.Rules.Value, wantRules) {
				t.Fatalf("service access update changed unrelated fields: role = %#v, rules = %#v, want rules = %#v", role, role.Rules.Value, wantRules)
			}
		})
	}
}
