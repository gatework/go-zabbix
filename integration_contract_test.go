//go:build integration

package zabbix

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func integrationValue[T any](value T) *T { return &value }

func integrationCheck(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		// This suite uses a disposable database and synthetic fixture values.
		if apiErr, ok := AsAPIError(err); ok {
			t.Fatalf("%v: %s", err, apiErr.Data)
		}
		t.Fatal(err)
	}
}

func integrationID(t *testing.T, ids []string, err error) string {
	t.Helper()
	integrationCheck(t, err)
	if len(ids) != 1 || ids[0] == "" {
		t.Fatalf("expected one nonempty id, got %v", ids)
	}
	return ids[0]
}

func integrationDelete(t *testing.T, remove func(context.Context, ...string) ([]string, error), id string) {
	t.Helper()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if _, err := remove(ctx, id); err != nil {
			t.Errorf("fixture cleanup: %v", err)
		}
	})
}

func integrationHost(t *testing.T, c *Client, ctx context.Context) (string, string, string) {
	t.Helper()
	name := fmt.Sprintf("go-zabbix-it-%d", time.Now().UnixNano())
	ids, err := c.HostGroup().Create(ctx, HostGroupCreateParams{Name: name})
	groupID := integrationID(t, ids, err)
	integrationDelete(t, c.HostGroup().Delete, groupID)
	ids, err = c.Host().Create(ctx, HostCreateParams{
		Host: name, Groups: []HostGroupRef{{GroupID: groupID}},
		Description: integrationValue("before"),
		Tags:        []TagInput{{Tag: "fixture", Value: integrationValue("before")}},
		Macros:      []HostMacroInput{{Macro: "{$SDK_FIXTURE}", Value: "before"}},
	})
	hostID := integrationID(t, ids, err)
	integrationDelete(t, c.Host().Delete, hostID)
	return hostID, name, groupID
}

func TestIntegrationHostContractRoundTrip(t *testing.T) {
	c, ctx := loginClient(t)
	hostID, _, _ := integrationHost(t, c, ctx)
	params := HostGetParams{
		Output: []string{"hostid", "description"}, HostIDs: []string{hostID},
		SelectTags: "extend", SelectMacros: "extend",
	}
	hosts, err := c.Host().Get(ctx, params)
	integrationCheck(t, err)
	if len(hosts) != 1 || hosts[0].Description != "before" || len(hosts[0].Tags) != 1 || len(hosts[0].Macros) != 1 {
		t.Fatalf("fixture not persisted: %+v", hosts)
	}
	// An omitted field survives a different update.
	_, err = c.Host().Update(ctx, HostUpdateParams{HostID: hostID, Status: integrationValue(0)})
	integrationCheck(t, err)
	hosts, err = c.Host().Get(ctx, params)
	integrationCheck(t, err)
	if len(hosts) != 1 || hosts[0].Description != "before" || len(hosts[0].Tags) != 1 {
		t.Fatal("omitted fields were modified")
	}
	_, err = c.Host().Update(ctx, HostUpdateParams{
		HostID: hostID, Description: integrationValue(""),
		Tags: []TagInput{}, Macros: []HostMacroUpdateInput{},
	})
	integrationCheck(t, err)
	hosts, err = c.Host().Get(ctx, params)
	integrationCheck(t, err)
	if len(hosts) != 1 || hosts[0].Description != "" || len(hosts[0].Tags) != 0 || len(hosts[0].Macros) != 0 {
		t.Fatalf("explicit empty values did not clear fields: %+v", hosts)
	}

	count, err := c.Host().Count(ctx, HostCountParams{HostIDs: []string{hostID}})
	integrationCheck(t, err)
	if count != 1 {
		t.Fatalf("count = %d", count)
	}
	byID, err := c.Host().GetMap(ctx, params)
	integrationCheck(t, err)
	if len(byID) != 1 || byID[hostID].HostID != hostID {
		t.Fatalf("keyed host result = %+v", byID)
	}
	empty, err := c.Host().GetMap(ctx, HostGetParams{HostIDs: []string{}})
	integrationCheck(t, err)
	if empty == nil || len(empty) != 0 {
		t.Fatalf("explicit empty ID selection = %+v", empty)
	}

	ids, err := c.UserMacro().Create(ctx, HostMacroCreateParams{
		HostID: hostID, Macro: "{$SDK_STANDALONE}", Value: "before",
	})
	macroID := integrationID(t, ids, err)
	_, err = c.UserMacro().Update(ctx, HostMacroUpdateParams{HostMacroID: macroID, Value: integrationValue("")})
	integrationCheck(t, err)
	macros, err := c.UserMacro().Get(ctx, UserMacroGetParams{HostIDs: []string{hostID}, Output: "extend"})
	integrationCheck(t, err)
	if len(macros) != 1 || macros[0].HostID != hostID || macros[0].Value != "" {
		t.Fatalf("standalone macro ownership or clear failed: %+v", macros)
	}

	ids, err = c.HostInterface().MassAdd(ctx, HostInterfaceMassAddParams{
		Hosts:      []HostRef{{HostID: hostID}},
		Interfaces: []HostInterfaceInput{{Type: 1, Main: 1, UseIP: 1, IP: "127.0.0.1", Port: "10050"}},
	})
	integrationID(t, ids, err)
	hosts, err = c.Host().Get(ctx, HostGetParams{HostIDs: []string{hostID}, SelectInterfaces: "count"})
	integrationCheck(t, err)
	if len(hosts) != 1 || hosts[0].Interfaces.Count == nil || *hosts[0].Interfaces.Count != 1 {
		t.Fatal("selected interface count not decoded")
	}
	ids, err = c.HostInterface().MassRemove(ctx, HostInterfaceMassRemoveParams{
		HostIDs: []string{hostID}, Interfaces: []HostInterfaceMatchInput{{IP: "127.0.0.2", Port: "10050"}},
	})
	integrationCheck(t, err)
	if len(ids) != 0 {
		t.Fatalf("nonmatching removal returned ids %v", ids)
	}
	ids, err = c.HostInterface().ReplaceHostInterfaces(ctx, HostInterfaceReplaceParams{
		HostID: hostID, Interfaces: []HostInterfaceUpdateInput{},
	})
	integrationCheck(t, err)
	if len(ids) != 0 {
		t.Fatalf("empty replacement returned ids %v", ids)
	}
	interfaces, err := c.HostInterface().Get(ctx, HostInterfaceGetParams{HostIDs: []string{hostID}})
	integrationCheck(t, err)
	if len(interfaces) != 0 {
		t.Fatal("interfaces were not cleared")
	}
}

func TestIntegrationTemplateGroupsAndRights(t *testing.T) {
	c, ctx := loginClient(t)
	_, name, groupID := integrationHost(t, c, ctx)
	ids, err := c.TemplateGroup().Create(ctx, TemplateGroupCreateParams{Name: name})
	templateGroupID := integrationID(t, ids, err)
	integrationDelete(t, c.TemplateGroup().Delete, templateGroupID)
	ids, err = c.Template().Create(ctx, TemplateCreateParams{Host: name + "-template", Groups: []TemplateGroupRef{{GroupID: templateGroupID}}})
	templateID := integrationID(t, ids, err)
	integrationDelete(t, c.Template().Delete, templateID)
	templates, err := c.Template().Get(ctx, TemplateGetParams{
		TemplateIDs: []string{templateID}, Output: []string{"templateid"}, SelectTemplateGroups: "extend",
	})
	integrationCheck(t, err)
	if len(templates) != 1 || len(templates[0].TemplateGroups) != 1 || templates[0].TemplateGroups[0].GroupID != templateGroupID {
		t.Fatalf("template group result lost: %+v", templates)
	}
	ids, err = c.UserGroup().Create(ctx, UserGroupCreateParams{
		Name: name, HostGroupRights: []UserGroupRightInput{{ID: groupID, Permission: 2}},
		TemplateGroupRights: []UserGroupRightInput{{ID: templateGroupID, Permission: 2}},
	})
	userGroupID := integrationID(t, ids, err)
	integrationDelete(t, c.UserGroup().Delete, userGroupID)
	params := UserGroupGetParams{
		UsrGrpIDs: []string{userGroupID}, Output: "extend",
		SelectHostGroupRights: "extend", SelectTemplateGroupRights: "extend",
	}
	groups, err := c.UserGroup().Get(ctx, params)
	integrationCheck(t, err)
	if len(groups) != 1 || len(groups[0].HostGroupRights) != 1 || len(groups[0].TemplateGroupRights) != 1 {
		t.Fatalf("split rights not returned: %+v", groups)
	}
	_, err = c.UserGroup().Update(ctx, UserGroupUpdateParams{
		UsrGrpID: userGroupID, HostGroupRights: []UserGroupRightInput{}, TemplateGroupRights: []UserGroupRightInput{},
	})
	integrationCheck(t, err)
	groups, err = c.UserGroup().Get(ctx, params)
	integrationCheck(t, err)
	if len(groups) != 1 || len(groups[0].HostGroupRights) != 0 || len(groups[0].TemplateGroupRights) != 0 {
		t.Fatal("explicit empty rights did not revoke fixture permissions")
	}
}

func TestIntegrationDiscoveryDeletionResults(t *testing.T) {
	c, ctx := loginClient(t)
	hostID, _, _ := integrationHost(t, c, ctx)
	ids, err := c.DiscoveryRule().Create(ctx, DiscoveryRuleCreateParams{
		HostID: hostID, Name: "SDK discovery", Key: "sdk.discovery", Type: 2,
	})
	ruleID := integrationID(t, ids, err)
	ids, err = c.ItemPrototype().Create(ctx, ItemPrototypeCreateParams{
		HostID: hostID, RuleID: ruleID, Name: "SDK item {#ID}", Key: "sdk.item[{#ID}]", Type: 2, ValueType: 3,
	})
	prototypeID := integrationID(t, ids, err)
	ids, err = c.ItemPrototype().Delete(ctx, prototypeID)
	if got := integrationID(t, ids, err); got != prototypeID {
		t.Fatalf("deleted prototype id = %s", got)
	}
	ids, err = c.DiscoveryRule().Delete(ctx, ruleID)
	if got := integrationID(t, ids, err); got != ruleID {
		t.Fatalf("deleted rule id = %s", got)
	}
}

func TestIntegrationLoginWithUserData(t *testing.T) {
	c, err := New(*zbxURL, WithTimeout(15*time.Second))
	integrationCheck(t, err)
	defer c.CloseIdleConnections()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	session, err := c.LoginWithUserData(ctx, *zbxUser, *zbxPass)
	integrationCheck(t, err)
	if session.SessionID == "" || session.UserID == "" || !c.HasCredentials() {
		t.Fatal("incomplete login session")
	}
	integrationCheck(t, c.Logout(ctx))
	if c.HasCredentials() {
		t.Fatal("logout retained session credentials")
	}
}

func TestIntegrationPartialStateDependentUpdates(t *testing.T) {
	c, ctx := loginClient(t)
	name := fmt.Sprintf("go-zabbix-partial-%d", time.Now().UnixNano())
	ids, err := c.Service().Create(ctx, ServiceCreateParams{
		Name: name, Algorithm: 0, SortOrder: 0,
		PropagationRule: integrationValue(1), PropagationValue: integrationValue(1),
	})
	serviceID := integrationID(t, ids, err)
	integrationDelete(t, c.Service().Delete, serviceID)
	_, err = c.Service().Update(ctx, ServiceUpdateParams{
		ServiceID: serviceID, PropagationValue: integrationValue(2),
	})
	integrationCheck(t, err)
	services, err := c.Service().Get(ctx, ServiceGetParams{ServiceIDs: []string{serviceID}, Output: "extend"})
	integrationCheck(t, err)
	if len(services) != 1 || services[0].PropagationRule != "1" || services[0].PropagationValue != "2" {
		t.Fatalf("partial propagation update = %+v", services)
	}
	ids, err = c.Proxy().Create(ctx, ProxyCreateParams{Name: name, OperatingMode: 0})
	proxyID := integrationID(t, ids, err)
	integrationDelete(t, c.Proxy().Delete, proxyID)
	hostID, _, _ := integrationHost(t, c, ctx)
	_, err = c.Host().Update(ctx, HostUpdateParams{
		HostID: hostID, MonitoredBy: integrationValue(1), ProxyID: &proxyID,
	})
	integrationCheck(t, err)
	_, err = c.Host().Update(ctx, HostUpdateParams{
		HostID: hostID, MonitoredBy: integrationValue(1), Name: &name,
	})
	integrationCheck(t, err)
	hosts, err := c.Host().Get(ctx, HostGetParams{HostIDs: []string{hostID}, Output: "extend"})
	integrationCheck(t, err)
	if len(hosts) != 1 || hosts[0].ProxyID != proxyID || hosts[0].Name != name {
		t.Fatalf("partial proxy host update = %+v", hosts)
	}
}
