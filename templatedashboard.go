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
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templatedashboard/object
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
// TemplateDashboardGetParams queries the list form of templatedashboard.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type TemplateDashboardGetParams struct {
	Output                 any            `json:"output,omitempty"`
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	SortField              []string       `json:"sortfield,omitzero"`
	SortOrder              []string       `json:"sortorder,omitzero"`
	Limit                  int            `json:"limit,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	DashboardIDs           []string       `json:"dashboardids,omitzero"`
	TemplateIDs            []string       `json:"templateids,omitzero"`
	SelectPages            any            `json:"selectPages,omitempty"`
}

// Get retrieves template dashboards matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templatedashboard/get
func (s *TemplateDashboardService) Get(ctx context.Context, params TemplateDashboardGetParams) ([]TemplateDashboard, error) {
	return getList[TemplateDashboard](ctx, s.client, "templatedashboard.get", params)
}

// TemplateDashboardCountParams contains the filters for templatedashboard.get counts.
type TemplateDashboardCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	DashboardIDs           []string       `json:"dashboardids,omitzero"`
	TemplateIDs            []string       `json:"templateids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *TemplateDashboardService) Count(ctx context.Context, params TemplateDashboardCountParams) (uint64, error) {
	return getCount(ctx, s.client, "templatedashboard.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *TemplateDashboardService) GetMap(ctx context.Context, params TemplateDashboardGetParams) (map[string]TemplateDashboard, error) {
	return getMap[TemplateDashboard](ctx, s.client, "templatedashboard.get", params)
}

// Create creates template dashboards and returns the new dashboard ids.
func (s *TemplateDashboardService) Create(ctx context.Context, items ...TemplateDashboardCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "templatedashboard.create", "dashboardids", items)
}

// Update updates template dashboards and returns the affected dashboard ids.
func (s *TemplateDashboardService) Update(ctx context.Context, items ...TemplateDashboardUpdateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "templatedashboard.update", "dashboardids", items)
}

// Delete deletes template dashboards by id and returns the deleted dashboard ids.
func (s *TemplateDashboardService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "templatedashboard.delete", "dashboardids", ids)
}
