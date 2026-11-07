//go:build integration

package zabbix

import (
	"context"
	"testing"
)

// TestIntegrationCRUD validates create/delete id-result decoding (the `<x>ids`
// response keys) for a representative set of writable namespaces against a live
// Zabbix 7.4 server. Each sub-test creates a throwaway object, asserts the
// returned id is non-empty, and deletes it (with cleanup as a safety net).
func TestIntegrationCRUD(t *testing.T) {
	c, ctx := loginClient(t)

	// mustDelete registers a best-effort cleanup so a failed assertion never leaks
	// server-side objects.
	del := func(t *testing.T, fn func(context.Context, ...string) ([]string, error), id string) {
		t.Cleanup(func() { _, _ = fn(context.Background(), id) })
	}
	check := func(t *testing.T, ids []string, err error) string {
		t.Helper()
		if err != nil {
			t.Fatalf("create: %v", err)
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
	}

	t.Run("templategroup", func(t *testing.T) {
		ids, err := c.TemplateGroup().Create(ctx, TemplateGroup{Name: "gozabbix-it-tg"})
		id := check(t, ids, err)
		del(t, c.TemplateGroup().Delete, id)
		d, err := c.TemplateGroup().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("usergroup", func(t *testing.T) {
		ids, err := c.UserGroup().Create(ctx, UserGroup{Name: "gozabbix-it-ug"})
		id := check(t, ids, err)
		del(t, c.UserGroup().Delete, id)
		d, err := c.UserGroup().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("role", func(t *testing.T) {
		ids, err := c.Role().Create(ctx, Role{Name: "gozabbix-it-role", Type: "1"})
		id := check(t, ids, err)
		del(t, c.Role().Delete, id)
		d, err := c.Role().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("regexp", func(t *testing.T) {
		ids, err := c.Regexp().Create(ctx, Regexp{
			Name:       "gozabbix-it-re",
			TestString: "abc",
			Expressions: []RegexpExpression{
				// expression_type 0 (character string included) requires an empty delimiter.
				{Expression: "abc", ExpressionType: "0", CaseSensitive: "0"},
			},
		})
		id := check(t, ids, err)
		del(t, c.Regexp().Delete, id)
		d, err := c.Regexp().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("service", func(t *testing.T) {
		ids, err := c.Service().Create(ctx, Service{Name: "gozabbix-it-svc", Algorithm: "1", SortOrder: "0"})
		id := check(t, ids, err)
		del(t, c.Service().Delete, id)
		d, err := c.Service().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("connector", func(t *testing.T) {
		ids, err := c.Connector().Create(ctx, Connector{Name: "gozabbix-it-conn", URL: "https://example.invalid/collector"})
		id := check(t, ids, err)
		del(t, c.Connector().Delete, id)
		d, err := c.Connector().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("proxy", func(t *testing.T) {
		ids, err := c.Proxy().Create(ctx, Proxy{Name: "gozabbix-it-proxy", OperatingMode: "0"})
		id := check(t, ids, err)
		del(t, c.Proxy().Delete, id)
		d, err := c.Proxy().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("valuemap", func(t *testing.T) {
		// value maps are host/template scoped; attach to the default Zabbix server host.
		ids, err := c.ValueMap().Create(ctx, ValueMap{
			HostID: "10084",
			Name:   "gozabbix-it-vm",
			Mappings: []ValueMapMapping{
				{Type: "0", Value: "1", NewValue: "up"},
			},
		})
		id := check(t, ids, err)
		del(t, c.ValueMap().Delete, id)
		d, err := c.ValueMap().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("map", func(t *testing.T) {
		ids, err := c.Map().Create(ctx, Map{Name: "gozabbix-it-map", Width: "800", Height: "600"})
		id := check(t, ids, err)
		del(t, c.Map().Delete, id)
		d, err := c.Map().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("dashboard", func(t *testing.T) {
		ids, err := c.Dashboard().Create(ctx, Dashboard{
			Name:  "gozabbix-it-dash",
			Pages: []DashboardPage{{}},
		})
		id := check(t, ids, err)
		del(t, c.Dashboard().Delete, id)
		d, err := c.Dashboard().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("usermacro-global", func(t *testing.T) {
		ids, err := c.UserMacro().CreateGlobal(ctx, GlobalMacro{Macro: "{$GOZABBIX_IT}", Value: "1"})
		id := check(t, ids, err)
		t.Cleanup(func() { _, _ = c.UserMacro().DeleteGlobal(context.Background(), id) })
		d, err := c.UserMacro().DeleteGlobal(ctx, id)
		roundtripDelete(t, d, err, id)
	})

	t.Run("token", func(t *testing.T) {
		ids, err := c.APIToken().Create(ctx, Token{Name: "gozabbix-it-token"})
		id := check(t, ids, err)
		t.Cleanup(func() { _, _ = c.APIToken().Delete(context.Background(), id) })
		d, err := c.APIToken().Delete(ctx, id)
		roundtripDelete(t, d, err, id)
	})
}
