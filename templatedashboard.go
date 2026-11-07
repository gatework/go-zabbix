package zabbix

import "context"

// TemplateDashboardService wraps the "templatedashboard" API namespace.
type TemplateDashboardService struct{ client *Client }

// TemplateDashboard returns the templatedashboard service.
func (c *Client) TemplateDashboard() *TemplateDashboardService {
	return &TemplateDashboardService{c}
}

// TemplateDashboard object (Zabbix 7.x). Reuses DashboardPage / DashboardWidget
// / DashboardWidgetField defined in dashboard.go.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/templatedashboard/object
type TemplateDashboard struct {
	DashboardID   string `json:"dashboardid,omitempty"`
	TemplateID    string `json:"templateid,omitempty"`
	Name          string `json:"name,omitempty"`
	DisplayPeriod string `json:"display_period,omitempty"`
	AutoStart     string `json:"auto_start,omitempty"`
	UUID          string `json:"uuid,omitempty"`

	Pages []DashboardPage `json:"pages,omitempty"`
}

// TemplateDashboardGetParams are the parameters for templatedashboard.get.
type TemplateDashboardGetParams struct {
	GetParams
	DashboardIDs []string `json:"dashboardids,omitempty"`
	TemplateIDs  []string `json:"templateids,omitempty"`
	SelectPages  any      `json:"selectPages,omitempty"`
}

// Get retrieves template dashboards matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/templatedashboard/get
func (s *TemplateDashboardService) Get(ctx context.Context, params TemplateDashboardGetParams) ([]TemplateDashboard, error) {
	var out []TemplateDashboard
	err := s.client.Call(ctx, "templatedashboard.get", params, &out)
	return out, err
}

// Create creates template dashboards and returns the new dashboard ids.
func (s *TemplateDashboardService) Create(ctx context.Context, items ...TemplateDashboard) ([]string, error) {
	var res struct {
		IDs []string `json:"dashboardids"`
	}
	err := s.client.Call(ctx, "templatedashboard.create", items, &res)
	return res.IDs, err
}

// Update updates template dashboards and returns the affected dashboard ids.
func (s *TemplateDashboardService) Update(ctx context.Context, items ...TemplateDashboard) ([]string, error) {
	var res struct {
		IDs []string `json:"dashboardids"`
	}
	err := s.client.Call(ctx, "templatedashboard.update", items, &res)
	return res.IDs, err
}

// Delete deletes template dashboards by id and returns the deleted dashboard ids.
func (s *TemplateDashboardService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"dashboardids"`
	}
	err := s.client.Call(ctx, "templatedashboard.delete", ids, &res)
	return res.IDs, err
}
