package zabbix

import "context"

// DashboardService wraps the "dashboard" API namespace.
type DashboardService struct{ client *Client }

// Dashboard returns the dashboard service.
func (c *Client) Dashboard() *DashboardService { return &DashboardService{c} }

// Dashboard object (Zabbix 7.x).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dashboard/object
type Dashboard struct {
	DashboardID   string `json:"dashboardid,omitempty"`
	Name          string `json:"name,omitempty"`
	UserID        string `json:"userid,omitempty"`
	Private       string `json:"private,omitempty"`
	DisplayPeriod string `json:"display_period,omitempty"`
	AutoStart     string `json:"auto_start,omitempty"`
	UUID          string `json:"uuid,omitempty"`
	TemplateID    string `json:"templateid,omitempty"`

	Pages      []DashboardPage      `json:"pages,omitempty"`
	Users      []DashboardUser      `json:"users,omitempty"`
	UserGroups []DashboardUserGroup `json:"userGroups,omitempty"`
}

// DashboardPage is a single page of a dashboard. Shared with template dashboards.
type DashboardPage struct {
	DashboardPageID string            `json:"dashboard_pageid,omitempty"`
	Name            string            `json:"name,omitempty"`
	DisplayPeriod   string            `json:"display_period,omitempty"`
	Widgets         []DashboardWidget `json:"widgets,omitempty"`
}

// DashboardWidget is a widget placed on a dashboard page.
type DashboardWidget struct {
	WidgetID string                 `json:"widgetid,omitempty"`
	Type     string                 `json:"type,omitempty"`
	Name     string                 `json:"name,omitempty"`
	X        string                 `json:"x,omitempty"`
	Y        string                 `json:"y,omitempty"`
	Width    string                 `json:"width,omitempty"`
	Height   string                 `json:"height,omitempty"`
	ViewMode string                 `json:"view_mode,omitempty"`
	Fields   []DashboardWidgetField `json:"fields,omitempty"`
}

// DashboardWidgetField is a configuration field of a dashboard widget.
type DashboardWidgetField struct {
	Type  string `json:"type,omitempty"`
	Name  string `json:"name,omitempty"`
	Value any    `json:"value,omitempty"`
}

// DashboardUser grants a user access to a dashboard.
type DashboardUser struct {
	UserID     string `json:"userid,omitempty"`
	Permission string `json:"permission,omitempty"`
}

// DashboardUserGroup grants a user group access to a dashboard.
type DashboardUserGroup struct {
	UsrGrpID   string `json:"usrgrpid,omitempty"`
	Permission string `json:"permission,omitempty"`
}

// DashboardGetParams are the parameters for dashboard.get.
// DashboardGetParams queries the list form of dashboard.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type DashboardGetParams struct {
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
	SelectPages            any            `json:"selectPages,omitempty"`
	SelectUsers            any            `json:"selectUsers,omitempty"`
	SelectUserGroups       any            `json:"selectUserGroups,omitempty"`
}

// Get retrieves dashboards matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dashboard/get
func (s *DashboardService) Get(ctx context.Context, params DashboardGetParams) ([]Dashboard, error) {
	return getList[Dashboard](ctx, s.client, "dashboard.get", params)
}

// DashboardCountParams contains the filters for dashboard.get counts.
type DashboardCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	DashboardIDs           []string       `json:"dashboardids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *DashboardService) Count(ctx context.Context, params DashboardCountParams) (uint64, error) {
	return getCount(ctx, s.client, "dashboard.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *DashboardService) GetMap(ctx context.Context, params DashboardGetParams) (map[string]Dashboard, error) {
	return getMap[Dashboard](ctx, s.client, "dashboard.get", params)
}

// Create creates dashboards and returns the new dashboard ids.
func (s *DashboardService) Create(ctx context.Context, items ...Dashboard) ([]string, error) {
	var res struct {
		IDs []string `json:"dashboardids"`
	}
	err := s.client.Call(ctx, "dashboard.create", items, &res)
	return res.IDs, err
}

// Update updates dashboards and returns the affected dashboard ids.
func (s *DashboardService) Update(ctx context.Context, items ...Dashboard) ([]string, error) {
	var res struct {
		IDs []string `json:"dashboardids"`
	}
	err := s.client.Call(ctx, "dashboard.update", items, &res)
	return res.IDs, err
}

// Delete deletes dashboards by id and returns the deleted dashboard ids.
func (s *DashboardService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"dashboardids"`
	}
	err := s.client.Call(ctx, "dashboard.delete", ids, &res)
	return res.IDs, err
}
