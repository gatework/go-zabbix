//go:build integration

package zabbix

import "testing"

func intPtr(i int) *int { return &i }

// TestIntegrationNamespaceReads calls the read path of every implemented namespace
// against a live Zabbix 7.4 server. The goal is to catch wrong method names or
// wrong outgoing parameter names (which the server rejects with an API error);
// empty result sets are fine on a fresh install.
func TestIntegrationNamespaceReads(t *testing.T) {
	c, ctx := loginClient(t)
	lim := GetParams{Limit: 3}

	checks := []struct {
		name string
		fn   func() error
	}{
		{"usergroup.get", func() error { _, e := c.UserGroup().Get(ctx, UserGroupGetParams{GetParams: lim}); return e }},
		{"role.get", func() error { _, e := c.Role().Get(ctx, RoleGetParams{GetParams: lim}); return e }},
		{"usermacro.get", func() error { _, e := c.UserMacro().Get(ctx, UserMacroGetParams{GetParams: lim}); return e }},
		{"usermacro.getglobal", func() error { _, e := c.UserMacro().GetGlobal(ctx, UserMacroGetParams{GetParams: lim}); return e }},
		{"userdirectory.get", func() error { _, e := c.UserDirectory().Get(ctx, UserDirectoryGetParams{GetParams: lim}); return e }},
		{"token.get", func() error { _, e := c.APIToken().Get(ctx, TokenGetParams{GetParams: lim}); return e }},
		{"mfa.get", func() error { _, e := c.MFA().Get(ctx, MFAGetParams{GetParams: lim}); return e }},
		{"action.get", func() error { _, e := c.Action().Get(ctx, ActionGetParams{GetParams: lim}); return e }},
		{"alert.get", func() error { _, e := c.Alert().Get(ctx, AlertGetParams{GetParams: lim}); return e }},
		{"mediatype.get", func() error { _, e := c.MediaType().Get(ctx, MediaTypeGetParams{GetParams: lim}); return e }},
		{"script.get", func() error { _, e := c.Script().Get(ctx, ScriptGetParams{GetParams: lim}); return e }},
		{"correlation.get", func() error { _, e := c.Correlation().Get(ctx, CorrelationGetParams{GetParams: lim}); return e }},
		{"maintenance.get", func() error { _, e := c.Maintenance().Get(ctx, MaintenanceGetParams{GetParams: lim}); return e }},
		{"service.get", func() error { _, e := c.Service().Get(ctx, ServiceGetParams{GetParams: lim}); return e }},
		{"sla.get", func() error { _, e := c.SLA().Get(ctx, SLAGetParams{GetParams: lim}); return e }},
		{"report.get", func() error { _, e := c.Report().Get(ctx, ReportGetParams{GetParams: lim}); return e }},
		{"dashboard.get", func() error { _, e := c.Dashboard().Get(ctx, DashboardGetParams{GetParams: lim}); return e }},
		{"templatedashboard.get", func() error {
			_, e := c.TemplateDashboard().Get(ctx, TemplateDashboardGetParams{GetParams: lim})
			return e
		}},
		{"module.get", func() error { _, e := c.Module().Get(ctx, ModuleGetParams{GetParams: lim}); return e }},
		{"drule.get", func() error { _, e := c.DRule().Get(ctx, DRuleGetParams{GetParams: lim}); return e }},
		{"dcheck.get", func() error { _, e := c.DCheck().Get(ctx, DCheckGetParams{GetParams: lim}); return e }},
		{"dhost.get", func() error { _, e := c.DHost().Get(ctx, DHostGetParams{GetParams: lim}); return e }},
		{"dservice.get", func() error { _, e := c.DService().Get(ctx, DServiceGetParams{GetParams: lim}); return e }},
		{"discoveryrule.get", func() error { _, e := c.DiscoveryRule().Get(ctx, DiscoveryRuleGetParams{GetParams: lim}); return e }},
		{"itemprototype.get", func() error { _, e := c.ItemPrototype().Get(ctx, ItemPrototypeGetParams{GetParams: lim}); return e }},
		{"triggerprototype.get", func() error {
			_, e := c.TriggerPrototype().Get(ctx, TriggerPrototypeGetParams{GetParams: lim})
			return e
		}},
		{"graphprototype.get", func() error { _, e := c.GraphPrototype().Get(ctx, GraphPrototypeGetParams{GetParams: lim}); return e }},
		{"hostprototype.get", func() error { _, e := c.HostPrototype().Get(ctx, HostPrototypeGetParams{GetParams: lim}); return e }},
		{"templategroup.get", func() error { _, e := c.TemplateGroup().Get(ctx, TemplateGroupGetParams{GetParams: lim}); return e }},
		{"graph.get", func() error { _, e := c.Graph().Get(ctx, GraphGetParams{GetParams: lim}); return e }},
		{"graphitem.get", func() error { _, e := c.GraphItem().Get(ctx, GraphItemGetParams{GetParams: lim}); return e }},
		{"valuemap.get", func() error { _, e := c.ValueMap().Get(ctx, ValueMapGetParams{GetParams: lim}); return e }},
		{"map.get", func() error { _, e := c.Map().Get(ctx, MapGetParams{GetParams: lim}); return e }},
		{"iconmap.get", func() error { _, e := c.IconMap().Get(ctx, IconMapGetParams{GetParams: lim}); return e }},
		{"image.get", func() error { _, e := c.Image().Get(ctx, ImageGetParams{GetParams: lim}); return e }},
		{"hostinterface.get", func() error { _, e := c.HostInterface().Get(ctx, HostInterfaceGetParams{GetParams: lim}); return e }},
		{"proxy.get", func() error { _, e := c.Proxy().Get(ctx, ProxyGetParams{GetParams: lim}); return e }},
		{"proxygroup.get", func() error { _, e := c.ProxyGroup().Get(ctx, ProxyGroupGetParams{GetParams: lim}); return e }},
		{"connector.get", func() error { _, e := c.Connector().Get(ctx, ConnectorGetParams{GetParams: lim}); return e }},
		{"auditlog.get", func() error { _, e := c.AuditLog().Get(ctx, AuditLogGetParams{GetParams: lim}); return e }},
		{"task.get", func() error { _, e := c.Task().Get(ctx, TaskGetParams{Output: "extend"}); return e }},
		{"regexp.get", func() error { _, e := c.Regexp().Get(ctx, RegexpGetParams{GetParams: lim}); return e }},
		{"hanode.get", func() error { _, e := c.HANode().Get(ctx, HANodeGetParams{GetParams: lim}); return e }},
		{"history.get", func() error {
			_, e := c.History().Get(ctx, HistoryGetParams{GetParams: lim, History: intPtr(3)})
			return e
		}},
		{"trend.get", func() error { _, e := c.Trend().Get(ctx, TrendGetParams{GetParams: lim}); return e }},

		// Global singletons (get takes no params).
		{"settings.get", func() error { _, e := c.Settings().Get(ctx); return e }},
		{"housekeeping.get", func() error { _, e := c.Housekeeping().Get(ctx); return e }},
		{"authentication.get", func() error { _, e := c.Authentication().Get(ctx); return e }},
		{"autoregistration.get", func() error { _, e := c.Autoregistration().Get(ctx); return e }},

		// configuration.export is a read-style operation.
		{"configuration.export", func() error {
			_, e := c.Configuration().Export(ctx, ConfigurationExportParams{
				Format:  "json",
				Options: map[string]any{"host_groups": []string{"4"}},
			})
			return e
		}},
	}

	for _, ch := range checks {
		ch := ch
		t.Run(ch.name, func(t *testing.T) {
			if err := ch.fn(); err != nil {
				t.Errorf("%s: %v", ch.name, err)
			}
		})
	}
}

// TestIntegrationSingletonFields sanity-checks that a singleton get decodes some
// expected non-empty field, proving output field names are right.
func TestIntegrationSingletonFields(t *testing.T) {
	c, ctx := loginClient(t)

	s, err := c.Settings().Get(ctx)
	if err != nil {
		t.Fatalf("settings.get: %v", err)
	}
	if s.DefaultLang == "" && s.DefaultTimezone == "" {
		t.Errorf("settings decoded with all-empty fields; output key mismatch? %+v", s)
	}

	hk, err := c.Housekeeping().Get(ctx)
	if err != nil {
		t.Fatalf("housekeeping.get: %v", err)
	}
	if hk.HKEventsMode == "" {
		t.Errorf("housekeeping decoded with empty hk_events_mode; field mismatch? %+v", hk)
	}
}
