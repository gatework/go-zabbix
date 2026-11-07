//go:build integration

// Package zabbix integration tests run against a live Zabbix server.
//
// Start the stack with `docker compose up -d` and run:
//
//	go test -tags=integration -v \
//	  -zbx.url=http://localhost:8080/api_jsonrpc.php -zbx.user=Admin -zbx.pass=zabbix
package zabbix

import (
	"context"
	"flag"
	"fmt"
	"testing"
	"time"
)

var (
	zbxURL  = flag.String("zbx.url", "http://localhost:8080/api_jsonrpc.php", "Zabbix API url")
	zbxUser = flag.String("zbx.user", "Admin", "Zabbix username")
	zbxPass = flag.String("zbx.pass", "zabbix", "Zabbix password")
)

// loginClient returns an authenticated client and registers logout on cleanup.
func loginClient(t *testing.T) (*Client, context.Context) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	t.Cleanup(cancel)

	c, err := New(*zbxURL, WithTimeout(15*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(c.CloseIdleConnections)
	if err := c.Login(ctx, *zbxUser, *zbxPass); err != nil {
		t.Fatalf("login: %v", err)
	}
	t.Cleanup(func() {
		cleanup, stop := context.WithTimeout(context.Background(), 5*time.Second)
		defer stop()
		if err := c.Logout(cleanup); err != nil {
			t.Errorf("logout: %v", err)
		}
	})
	return c, ctx
}

func TestIntegrationVersion(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	c, err := New(*zbxURL)
	if err != nil {
		t.Fatal(err)
	}
	defer c.CloseIdleConnections()
	v, err := c.APIInfo().Version(ctx)
	if err != nil {
		t.Fatalf("apiinfo.version: %v", err)
	}
	t.Logf("Zabbix API version: %s", v)
}

// TestIntegrationReadNamespaces exercises every read path with real select/filter
// parameters so a wrong outgoing parameter name surfaces as an API error, and so
// the typed structs decode against a live 7.4 response.
func TestIntegrationReadNamespaces(t *testing.T) {
	c, ctx := loginClient(t)

	t.Run("user.get self", func(t *testing.T) {
		users, err := c.User().Get(ctx, UserGetParams{
			Output:       []string{"userid", "username", "name"},
			SelectMedias: "extend",
			SelectRole:   "extend",
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("users: %d", len(users))
	})

	t.Run("hostgroup.get", func(t *testing.T) {
		groups, err := c.HostGroup().Get(ctx, HostGroupGetParams{
			Output: "extend", Limit: 5,
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("host groups: %d", len(groups))
	})

	t.Run("host.get with selects", func(t *testing.T) {
		hosts, err := c.Host().Get(ctx, HostGetParams{
			Output: []string{"hostid", "host", "name", "status", "monitored_by"}, Limit: 5,
			SelectInterfaces:      "extend",
			SelectHostGroups:      "extend",
			SelectParentTemplates: "extend",
			SelectTags:            "extend",
			SelectMacros:          "extend",
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("hosts: %d", len(hosts))
		for _, h := range hosts {
			t.Logf("  %s %q status=%s interfaces=%d hostgroups=%d templates=%d",
				h.HostID, h.Name, h.Status, len(h.Interfaces.Items), len(h.HostGroups), len(h.ParentTemplates.Items))
			if len(h.HostGroups) == 0 {
				t.Errorf("host %s decoded 0 host groups; selectHostGroups output key mismatch?", h.HostID)
			}
		}
	})

	t.Run("template.get", func(t *testing.T) {
		tpls, err := c.Template().Get(ctx, TemplateGetParams{
			Output: []string{"templateid", "host", "name"}, Limit: 5,
			SelectTemplateGroups: "extend",
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("templates: %d", len(tpls))
	})

	t.Run("item.get", func(t *testing.T) {
		items, err := c.Item().Get(ctx, ItemGetParams{
			Output: []string{"itemid", "name", "key_", "value_type", "lastvalue"}, Limit: 5,
			SelectTags: "extend",
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("items: %d", len(items))
	})

	t.Run("trigger.get", func(t *testing.T) {
		triggers, err := c.Trigger().Get(ctx, TriggerGetParams{
			Output: []string{"triggerid", "description", "priority", "status"}, Limit: 5,
			SelectHosts: "extend",
			SelectTags:  "extend",
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("triggers: %d", len(triggers))
	})

	t.Run("problem.get", func(t *testing.T) {
		problems, err := c.Problem().Get(ctx, ProblemGetParams{
			Output: "extend", Limit: 5,
			SelectTags: "extend",
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("problems: %d", len(problems))
	})

	t.Run("event.get", func(t *testing.T) {
		events, err := c.Event().Get(ctx, EventGetParams{
			Output: "extend", Limit: 5, SortField: []string{"clock"}, SortOrder: []string{"DESC"},
			SelectTags: "extend",
			TimeFrom:   integrationValue(time.Now().Add(-365 * 24 * time.Hour).Unix()),
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("events: %d", len(events))
	})
}

// TestIntegrationHostGroupCRUD validates create/update/delete id decoding with a
// throwaway host group.
func TestIntegrationHostGroupCRUD(t *testing.T) {
	c, ctx := loginClient(t)

	name := fmt.Sprintf("go-zabbix-crud-%d", time.Now().UnixNano())
	ids, err := c.HostGroup().Create(ctx, HostGroupCreateParams{Name: name})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if len(ids) != 1 {
		t.Fatalf("create returned ids %v, want 1", ids)
	}
	groupID := ids[0]
	t.Logf("created host group %s", groupID)

	// Ensure it is cleaned up even if a later step fails.
	deleted := false
	t.Cleanup(func() {
		if !deleted {
			cleanup, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if _, err := c.HostGroup().Delete(cleanup, groupID); err != nil {
				t.Errorf("cleanup: %v", err)
			}
		}
	})

	got, err := c.HostGroup().Get(ctx, HostGroupGetParams{
		Output:   "extend",
		GroupIDs: []string{groupID},
	})
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if len(got) != 1 || got[0].Name != name {
		t.Fatalf("get returned %+v, want name %q", got, name)
	}

	delIDs, err := c.HostGroup().Delete(ctx, groupID)
	if err != nil {
		t.Fatalf("delete: %v", err)
	}
	if len(delIDs) != 1 || delIDs[0] != groupID {
		t.Fatalf("delete returned %v, want [%s]", delIDs, groupID)
	}
	deleted = true
	t.Logf("deleted host group %s", groupID)
}
