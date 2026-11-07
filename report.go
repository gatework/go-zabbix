package zabbix

import "context"

// ReportService wraps the "report" API namespace.
type ReportService struct{ client *Client }

// Report returns the report service.
func (c *Client) Report() *ReportService { return &ReportService{c} }

// Report object (Zabbix 7.x scheduled report).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/report/object
type Report struct {
	ReportID    string `json:"reportid,omitempty"`
	UserID      string `json:"userid,omitempty"`
	Name        string `json:"name,omitempty"`
	DashboardID string `json:"dashboardid,omitempty"`
	Period      string `json:"period,omitempty"`
	Cycle       string `json:"cycle,omitempty"`
	StartTime   string `json:"start_time,omitempty"`
	ActiveSince string `json:"active_since,omitempty"`
	ActiveTill  string `json:"active_till,omitempty"`
	State       string `json:"state,omitempty"`
	LastSent    string `json:"lastsent,omitempty"`
	Info        string `json:"info,omitempty"`
	Status      string `json:"status,omitempty"`
	Description string `json:"description,omitempty"`
	Subject     string `json:"subject,omitempty"`
	Message     string `json:"message,omitempty"`

	Users      []ReportUser      `json:"users,omitempty"`
	UserGroups []ReportUserGroup `json:"user_groups,omitempty"`
}

// ReportUser is a user recipient of a scheduled report.
type ReportUser struct {
	UserID       string `json:"userid,omitempty"`
	AccessUserID string `json:"access_userid,omitempty"`
	Exclude      string `json:"exclude,omitempty"`
}

// ReportUserGroup is a user group recipient of a scheduled report.
type ReportUserGroup struct {
	UsrGrpID     string `json:"usrgrpid,omitempty"`
	AccessUserID string `json:"access_userid,omitempty"`
}

// ReportGetParams are the parameters for report.get.
type ReportGetParams struct {
	GetParams
	ReportIDs        []string `json:"reportids,omitempty"`
	SelectUsers      any      `json:"selectUsers,omitempty"`
	SelectUserGroups any      `json:"selectUserGroups,omitempty"`
}

// Get retrieves reports matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/report/get
func (s *ReportService) Get(ctx context.Context, params ReportGetParams) ([]Report, error) {
	var out []Report
	err := s.client.Call(ctx, "report.get", params, &out)
	return out, err
}

// Create creates reports and returns the new report ids.
func (s *ReportService) Create(ctx context.Context, items ...Report) ([]string, error) {
	var res struct {
		IDs []string `json:"reportids"`
	}
	err := s.client.Call(ctx, "report.create", items, &res)
	return res.IDs, err
}

// Update updates reports and returns the affected report ids.
func (s *ReportService) Update(ctx context.Context, items ...Report) ([]string, error) {
	var res struct {
		IDs []string `json:"reportids"`
	}
	err := s.client.Call(ctx, "report.update", items, &res)
	return res.IDs, err
}

// Delete deletes reports by id and returns the deleted report ids.
func (s *ReportService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"reportids"`
	}
	err := s.client.Call(ctx, "report.delete", ids, &res)
	return res.IDs, err
}
