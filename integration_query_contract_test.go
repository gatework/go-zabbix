//go:build integration

package zabbix

import (
	"errors"
	"testing"
)

// TestIntegrationQueryResultVariants verifies every newly exposed Count/GetMap
// against the server. Counts have no projection, sorting or limiting options.
func TestIntegrationQueryResultVariants(t *testing.T) {
	c, ctx := loginClient(t)
	checks := []struct {
		name string
		run  func() error
	}{
		{"action.Count", func() error { _, err := c.Action().Count(ctx, ActionCountParams{}); return err }},
		{"action.GetMap", func() error { return nonNilQueryMap(c.Action().GetMap(ctx, ActionGetParams{Limit: 3})) }},
		{"alert.Count", func() error { _, err := c.Alert().Count(ctx, AlertCountParams{}); return err }},
		{"alert.GetMap", func() error { return nonNilQueryMap(c.Alert().GetMap(ctx, AlertGetParams{Limit: 3})) }},
		{"auditlog.Count", func() error { _, err := c.AuditLog().Count(ctx, AuditLogCountParams{}); return err }},
		{"auditlog.GetMap", func() error { return nonNilQueryMap(c.AuditLog().GetMap(ctx, AuditLogGetParams{Limit: 3})) }},
		{"connector.Count", func() error { _, err := c.Connector().Count(ctx, ConnectorCountParams{}); return err }},
		{"connector.GetMap", func() error { return nonNilQueryMap(c.Connector().GetMap(ctx, ConnectorGetParams{Limit: 3})) }},
		{"correlation.Count", func() error { _, err := c.Correlation().Count(ctx, CorrelationCountParams{}); return err }},
		{"correlation.GetMap", func() error { return nonNilQueryMap(c.Correlation().GetMap(ctx, CorrelationGetParams{Limit: 3})) }},
		{"dashboard.Count", func() error { _, err := c.Dashboard().Count(ctx, DashboardCountParams{}); return err }},
		{"dashboard.GetMap", func() error { return nonNilQueryMap(c.Dashboard().GetMap(ctx, DashboardGetParams{Limit: 3})) }},
		{"dcheck.Count", func() error { _, err := c.DCheck().Count(ctx, DCheckCountParams{}); return err }},
		{"dcheck.GetMap", func() error { return nonNilQueryMap(c.DCheck().GetMap(ctx, DCheckGetParams{Limit: 3})) }},
		{"dhost.Count", func() error { _, err := c.DHost().Count(ctx, DHostCountParams{}); return err }},
		{"dhost.GetMap", func() error { return nonNilQueryMap(c.DHost().GetMap(ctx, DHostGetParams{Limit: 3})) }},
		{"discoveryrule.Count", func() error { _, err := c.DiscoveryRule().Count(ctx, DiscoveryRuleCountParams{}); return err }},
		{"discoveryrule.GetMap", func() error { return nonNilQueryMap(c.DiscoveryRule().GetMap(ctx, DiscoveryRuleGetParams{Limit: 3})) }},
		{"drule.Count", func() error { _, err := c.DRule().Count(ctx, DRuleCountParams{}); return err }},
		{"drule.GetMap", func() error { return nonNilQueryMap(c.DRule().GetMap(ctx, DRuleGetParams{Limit: 3})) }},
		{"dservice.Count", func() error { _, err := c.DService().Count(ctx, DServiceCountParams{}); return err }},
		{"dservice.GetMap", func() error { return nonNilQueryMap(c.DService().GetMap(ctx, DServiceGetParams{Limit: 3})) }},
		{"event.Count", func() error { _, err := c.Event().Count(ctx, EventCountParams{}); return err }},
		{"event.GetMap", func() error { return nonNilQueryMap(c.Event().GetMap(ctx, EventGetParams{Limit: 3})) }},
		{"graph.Count", func() error { _, err := c.Graph().Count(ctx, GraphCountParams{}); return err }},
		{"graph.GetMap", func() error { return nonNilQueryMap(c.Graph().GetMap(ctx, GraphGetParams{Limit: 3})) }},
		{"graphitem.Count", func() error { _, err := c.GraphItem().Count(ctx, GraphItemCountParams{}); return err }},
		{"graphitem.GetMap", func() error { return nonNilQueryMap(c.GraphItem().GetMap(ctx, GraphItemGetParams{Limit: 3})) }},
		{"graphprototype.Count", func() error { _, err := c.GraphPrototype().Count(ctx, GraphPrototypeCountParams{}); return err }},
		{"graphprototype.GetMap", func() error { return nonNilQueryMap(c.GraphPrototype().GetMap(ctx, GraphPrototypeGetParams{Limit: 3})) }},
		{"hanode.Count", func() error { _, err := c.HANode().Count(ctx, HANodeCountParams{}); return err }},
		{"hanode.GetMap", func() error { return nonNilQueryMap(c.HANode().GetMap(ctx, HANodeGetParams{Limit: 3})) }},
		{"history.Count", func() error { _, err := c.History().Count(ctx, HistoryCountParams{}); return err }},
		{"host.Count", func() error { _, err := c.Host().Count(ctx, HostCountParams{}); return err }},
		{"host.GetMap", func() error { return nonNilQueryMap(c.Host().GetMap(ctx, HostGetParams{Limit: 3})) }},
		{"hostgroup.Count", func() error { _, err := c.HostGroup().Count(ctx, HostGroupCountParams{}); return err }},
		{"hostgroup.GetMap", func() error { return nonNilQueryMap(c.HostGroup().GetMap(ctx, HostGroupGetParams{Limit: 3})) }},
		{"hostinterface.Count", func() error { _, err := c.HostInterface().Count(ctx, HostInterfaceCountParams{}); return err }},
		{"hostinterface.GetMap", func() error { return nonNilQueryMap(c.HostInterface().GetMap(ctx, HostInterfaceGetParams{Limit: 3})) }},
		{"hostprototype.Count", func() error { _, err := c.HostPrototype().Count(ctx, HostPrototypeCountParams{}); return err }},
		{"hostprototype.GetMap", func() error { return nonNilQueryMap(c.HostPrototype().GetMap(ctx, HostPrototypeGetParams{Limit: 3})) }},
		{"iconmap.Count", func() error { _, err := c.IconMap().Count(ctx, IconMapCountParams{}); return err }},
		{"iconmap.GetMap", func() error { return nonNilQueryMap(c.IconMap().GetMap(ctx, IconMapGetParams{Limit: 3})) }},
		{"image.Count", func() error { _, err := c.Image().Count(ctx, ImageCountParams{}); return err }},
		{"image.GetMap", func() error { return nonNilQueryMap(c.Image().GetMap(ctx, ImageGetParams{Limit: 3})) }},
		{"item.Count", func() error { _, err := c.Item().Count(ctx, ItemCountParams{}); return err }},
		{"item.GetMap", func() error { return nonNilQueryMap(c.Item().GetMap(ctx, ItemGetParams{Limit: 3})) }},
		{"itemprototype.Count", func() error { _, err := c.ItemPrototype().Count(ctx, ItemPrototypeCountParams{}); return err }},
		{"itemprototype.GetMap", func() error { return nonNilQueryMap(c.ItemPrototype().GetMap(ctx, ItemPrototypeGetParams{Limit: 3})) }},
		{"maintenance.Count", func() error { _, err := c.Maintenance().Count(ctx, MaintenanceCountParams{}); return err }},
		{"maintenance.GetMap", func() error { return nonNilQueryMap(c.Maintenance().GetMap(ctx, MaintenanceGetParams{Limit: 3})) }},
		{"map.Count", func() error { _, err := c.Map().Count(ctx, MapCountParams{}); return err }},
		{"map.GetMap", func() error { return nonNilQueryMap(c.Map().GetMap(ctx, MapGetParams{Limit: 3})) }},
		{"mediatype.Count", func() error { _, err := c.MediaType().Count(ctx, MediaTypeCountParams{}); return err }},
		{"mediatype.GetMap", func() error { return nonNilQueryMap(c.MediaType().GetMap(ctx, MediaTypeGetParams{Limit: 3})) }},
		{"mfa.Count", func() error { _, err := c.MFA().Count(ctx, MFACountParams{}); return err }},
		{"mfa.GetMap", func() error { return nonNilQueryMap(c.MFA().GetMap(ctx, MFAGetParams{Limit: 3})) }},
		{"module.Count", func() error { _, err := c.Module().Count(ctx, ModuleCountParams{}); return err }},
		{"module.GetMap", func() error { return nonNilQueryMap(c.Module().GetMap(ctx, ModuleGetParams{Limit: 3})) }},
		{"problem.Count", func() error { _, err := c.Problem().Count(ctx, ProblemCountParams{}); return err }},
		{"problem.GetMap", func() error { return nonNilQueryMap(c.Problem().GetMap(ctx, ProblemGetParams{Limit: 3})) }},
		{"proxy.Count", func() error { _, err := c.Proxy().Count(ctx, ProxyCountParams{}); return err }},
		{"proxy.GetMap", func() error { return nonNilQueryMap(c.Proxy().GetMap(ctx, ProxyGetParams{Limit: 3})) }},
		{"proxygroup.Count", func() error { _, err := c.ProxyGroup().Count(ctx, ProxyGroupCountParams{}); return err }},
		{"proxygroup.GetMap", func() error { return nonNilQueryMap(c.ProxyGroup().GetMap(ctx, ProxyGroupGetParams{Limit: 3})) }},
		{"regexp.Count", func() error { _, err := c.Regexp().Count(ctx, RegexpCountParams{}); return err }},
		{"regexp.GetMap", func() error { return nonNilQueryMap(c.Regexp().GetMap(ctx, RegexpGetParams{Limit: 3})) }},
		{"report.Count", func() error { _, err := c.Report().Count(ctx, ReportCountParams{}); return err }},
		{"report.GetMap", func() error { return nonNilQueryMap(c.Report().GetMap(ctx, ReportGetParams{Limit: 3})) }},
		{"role.Count", func() error { _, err := c.Role().Count(ctx, RoleCountParams{}); return err }},
		{"role.GetMap", func() error { return nonNilQueryMap(c.Role().GetMap(ctx, RoleGetParams{Limit: 3})) }},
		{"script.Count", func() error { _, err := c.Script().Count(ctx, ScriptCountParams{}); return err }},
		{"script.GetMap", func() error { return nonNilQueryMap(c.Script().GetMap(ctx, ScriptGetParams{Limit: 3})) }},
		{"service.Count", func() error { _, err := c.Service().Count(ctx, ServiceCountParams{}); return err }},
		{"service.GetMap", func() error { return nonNilQueryMap(c.Service().GetMap(ctx, ServiceGetParams{Limit: 3})) }},
		{"sla.Count", func() error { _, err := c.SLA().Count(ctx, SLACountParams{}); return err }},
		{"sla.GetMap", func() error { return nonNilQueryMap(c.SLA().GetMap(ctx, SLAGetParams{Limit: 3})) }},
		{"task.GetMap", func() error { return nonNilQueryMap(c.Task().GetMap(ctx, TaskGetParams{TaskIDs: []string{}})) }},
		{"template.Count", func() error { _, err := c.Template().Count(ctx, TemplateCountParams{}); return err }},
		{"template.GetMap", func() error { return nonNilQueryMap(c.Template().GetMap(ctx, TemplateGetParams{Limit: 3})) }},
		{"templatedashboard.Count", func() error { _, err := c.TemplateDashboard().Count(ctx, TemplateDashboardCountParams{}); return err }},
		{"templatedashboard.GetMap", func() error {
			return nonNilQueryMap(c.TemplateDashboard().GetMap(ctx, TemplateDashboardGetParams{Limit: 3}))
		}},
		{"templategroup.Count", func() error { _, err := c.TemplateGroup().Count(ctx, TemplateGroupCountParams{}); return err }},
		{"templategroup.GetMap", func() error { return nonNilQueryMap(c.TemplateGroup().GetMap(ctx, TemplateGroupGetParams{Limit: 3})) }},
		{"token.Count", func() error { _, err := c.APIToken().Count(ctx, TokenCountParams{}); return err }},
		{"token.GetMap", func() error { return nonNilQueryMap(c.APIToken().GetMap(ctx, TokenGetParams{Limit: 3})) }},
		{"trend.Count", func() error { _, err := c.Trend().Count(ctx, TrendCountParams{}); return err }},
		{"trigger.Count", func() error { _, err := c.Trigger().Count(ctx, TriggerCountParams{}); return err }},
		{"trigger.GetMap", func() error { return nonNilQueryMap(c.Trigger().GetMap(ctx, TriggerGetParams{Limit: 3})) }},
		{"triggerprototype.Count", func() error { _, err := c.TriggerPrototype().Count(ctx, TriggerPrototypeCountParams{}); return err }},
		{"triggerprototype.GetMap", func() error {
			return nonNilQueryMap(c.TriggerPrototype().GetMap(ctx, TriggerPrototypeGetParams{Limit: 3}))
		}},
		{"user.Count", func() error { _, err := c.User().Count(ctx, UserCountParams{}); return err }},
		{"user.GetMap", func() error { return nonNilQueryMap(c.User().GetMap(ctx, UserGetParams{Limit: 3})) }},
		{"userdirectory.Count", func() error { _, err := c.UserDirectory().Count(ctx, UserDirectoryCountParams{}); return err }},
		{"userdirectory.GetMap", func() error { return nonNilQueryMap(c.UserDirectory().GetMap(ctx, UserDirectoryGetParams{Limit: 3})) }},
		{"usergroup.Count", func() error { _, err := c.UserGroup().Count(ctx, UserGroupCountParams{}); return err }},
		{"usergroup.GetMap", func() error { return nonNilQueryMap(c.UserGroup().GetMap(ctx, UserGroupGetParams{Limit: 3})) }},
		{"usermacro.Count", func() error { _, err := c.UserMacro().Count(ctx, UserMacroCountParams{}); return err }},
		{"usermacro.GetMap", func() error { return nonNilQueryMap(c.UserMacro().GetMap(ctx, UserMacroGetParams{Limit: 3})) }},
		{"valuemap.Count", func() error { _, err := c.ValueMap().Count(ctx, ValueMapCountParams{}); return err }},
		{"valuemap.GetMap", func() error { return nonNilQueryMap(c.ValueMap().GetMap(ctx, ValueMapGetParams{Limit: 3})) }},
		{"usermacro.CountGlobal", func() error { _, err := c.UserMacro().CountGlobal(ctx, GlobalMacroCountParams{}); return err }},
		{"usermacro.GetGlobalMap", func() error { return nonNilQueryMap(c.UserMacro().GetGlobalMap(ctx, GlobalMacroGetParams{Limit: 3})) }},
	}
	for _, check := range checks {
		t.Run(check.name, func(t *testing.T) {
			if err := check.run(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func nonNilQueryMap[T any](result map[string]T, err error) error {
	if err != nil {
		return err
	}
	if result == nil {
		return errors.New("successful GetMap returned a nil map")
	}
	return nil
}
