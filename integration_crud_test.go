//go:build integration

package zabbix

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestIntegrationCRUD validates create/delete id-result decoding (the `<x>ids`
// response keys) for a representative set of writable namespaces against a live
// Zabbix 7.4 server. Each sub-test creates a throwaway object, asserts the
// returned id is non-empty, and deletes it (with cleanup as a safety net).
func TestIntegrationCRUD(t *testing.T) {
	c, ctx := loginClient(t)

	// Register bounded cleanup before checking subsequent operations.
	deleted := make(map[string]bool)
	suffix := fmt.Sprintf("-%d", time.Now().UnixNano())
	del := func(t *testing.T, fn func(context.Context, ...string) ([]string, error), id string) {
		t.Cleanup(func() {
			if deleted[t.Name()] {
				return
			}
			cleanup, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if _, err := fn(cleanup, id); err != nil {
				t.Errorf("cleanup: %v", err)
			}
		})
	}
	check := func(t *testing.T, ids []string, err error) string {
		t.Helper()
		if err != nil {
			integrationCheck(t, err)
		}
		if len(ids) != 1 || ids[0] == "" {
			t.Fatalf("create returned ids %v, want exactly one non-empty id", ids)
		}
		return ids[0]
	}
	roundtripDelete := func(t *testing.T, delIDs []string, err error, id string) {
		t.Helper()
		if err != nil {
			t.Fatalf("delete: %v", err)
		}
		if len(delIDs) != 1 || delIDs[0] != id {
			t.Fatalf("delete returned %v, want [%s]", delIDs, id)
		}
		deleted[t.Name()] = true
	}

	t.Run("templategroup", func(t *testing.T) {
		ids, err := c.TemplateGroup().Create(ctx, TemplateGroupCreateParams{Name: "gozabbix-it-tg" + suffix})
		id := check(t, ids, err)
		del(t, c.TemplateGroup().Delete, id)
		d, err := c.TemplateGroup().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("usergroup", func(t *testing.T) {
		ids, err := c.UserGroup().Create(ctx, UserGroupCreateParams{Name: "gozabbix-it-ug" + suffix})
		id := check(t, ids, err)
		del(t, c.UserGroup().Delete, id)
		d, err := c.UserGroup().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("role", func(t *testing.T) {
		ids, err := c.Role().Create(ctx, RoleCreateParams{Name: "gozabbix-it-role" + suffix, Type: 1})
		id := check(t, ids, err)
		del(t, c.Role().Delete, id)
		d, err := c.Role().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("regexp", func(t *testing.T) {
		ids, err := c.Regexp().Create(ctx, RegexpCreateParams{
			Name:       "gozabbix-it-re" + suffix,
			TestString: integrationValue("abc"),
			Expressions: []RegexpExpressionInput{
				// expression_type 0 (character string included) requires an empty delimiter.
				{Expression: "abc", ExpressionType: 0, CaseSensitive: integrationValue(0)},
			},
		})
		id := check(t, ids, err)
		del(t, c.Regexp().Delete, id)
		d, err := c.Regexp().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("service", func(t *testing.T) {
		ids, err := c.Service().Create(ctx, ServiceCreateParams{Name: "gozabbix-it-svc" + suffix, Algorithm: 1, SortOrder: 0})
		id := check(t, ids, err)
		del(t, c.Service().Delete, id)
		d, err := c.Service().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("connector", func(t *testing.T) {
		ids, err := c.Connector().Create(ctx, ConnectorCreateParams{Name: "gozabbix-it-conn" + suffix, URL: "https://example.invalid/collector"})
		id := check(t, ids, err)
		del(t, c.Connector().Delete, id)
		d, err := c.Connector().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("proxy", func(t *testing.T) {
		ids, err := c.Proxy().Create(ctx, ProxyCreateParams{Name: "gozabbix-it-proxy" + suffix, OperatingMode: 0})
		id := check(t, ids, err)
		del(t, c.Proxy().Delete, id)
		d, err := c.Proxy().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("valuemap", func(t *testing.T) {
		// Value maps belong to a disposable fixture host.
		hostID, _, _ := integrationHost(t, c, ctx)
		ids, err := c.ValueMap().Create(ctx, ValueMapCreateParams{
			HostID: hostID,
			Name:   "gozabbix-it-vm" + suffix,
			Mappings: []ValueMappingInput{
				{Type: integrationValue(0), Value: integrationValue("1"), NewValue: "up"},
			},
		})
		id := check(t, ids, err)
		del(t, c.ValueMap().Delete, id)
		d, err := c.ValueMap().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("map", func(t *testing.T) {
		ids, err := c.Map().Create(ctx, MapCreateParams{Name: "gozabbix-it-map" + suffix, Width: 800, Height: 600})
		id := check(t, ids, err)
		del(t, c.Map().Delete, id)
		d, err := c.Map().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("dashboard", func(t *testing.T) {
		ids, err := c.Dashboard().Create(ctx, DashboardCreateParams{
			Name:  "gozabbix-it-dash" + suffix,
			Pages: []DashboardPageInput{{}},
		})
		id := check(t, ids, err)
		del(t, c.Dashboard().Delete, id)
		d, err := c.Dashboard().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("usermacro-global", func(t *testing.T) {
		ids, err := c.UserMacro().CreateGlobal(ctx, GlobalMacroCreateParams{Macro: fmt.Sprintf("{$GOZABBIX_IT_%d}", time.Now().UnixNano()), Value: "1"})
		id := check(t, ids, err)
		del(t, c.UserMacro().DeleteGlobal, id)
		d, err := c.UserMacro().DeleteGlobal(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("token", func(t *testing.T) {
		ids, err := c.APIToken().Create(ctx, TokenCreateParams{Name: "gozabbix-it-token" + suffix})
		id := check(t, ids, err)
		del(t, c.APIToken().Delete, id)
		d, err := c.APIToken().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})
}
