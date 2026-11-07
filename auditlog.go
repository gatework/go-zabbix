package zabbix

import "context"

// AuditLogService wraps the "auditlog" API namespace.
type AuditLogService struct{ client *Client }

// AuditLog returns the auditlog service.
func (c *Client) AuditLog() *AuditLogService { return &AuditLogService{c} }

// AuditLog is a single audit log record.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/auditlog/object
type AuditLog struct {
	AuditID      string `json:"auditid,omitempty"`
	UserID       string `json:"userid,omitempty"`
	Username     string `json:"username,omitempty"`
	Clock        string `json:"clock,omitempty"`
	IP           string `json:"ip,omitempty"`
	Action       string `json:"action,omitempty"`
	ResourceType string `json:"resourcetype,omitempty"`
	ResourceID   string `json:"resourceid,omitempty"`
	ResourceCUID string `json:"resource_cuid,omitempty"`
	ResourceName string `json:"resourcename,omitempty"`
	RecordSetID  string `json:"recordsetid,omitempty"`
	Details      string `json:"details,omitempty"`
}

// AuditLogGetParams are the parameters for auditlog.get.
// AuditLogGetParams queries the list form of auditlog.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type AuditLogGetParams struct {
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
	UserIDs                []string       `json:"userids,omitzero"`
	TimeFrom               *int64         `json:"time_from,omitempty"`
	TimeTill               *int64         `json:"time_till,omitempty"`
}

// Get retrieves audit log records matching params. The auditlog API is read-only.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/auditlog/get
func (s *AuditLogService) Get(ctx context.Context, params AuditLogGetParams) ([]AuditLog, error) {
	return getList[AuditLog](ctx, s.client, "auditlog.get", params)
}

// AuditLogCountParams contains the filters for auditlog.get counts.
type AuditLogCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	UserIDs                []string       `json:"userids,omitzero"`
	TimeFrom               *int64         `json:"time_from,omitempty"`
	TimeTill               *int64         `json:"time_till,omitempty"`
}

// Count returns the number of objects matching params.
func (s *AuditLogService) Count(ctx context.Context, params AuditLogCountParams) (uint64, error) {
	return getCount(ctx, s.client, "auditlog.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *AuditLogService) GetMap(ctx context.Context, params AuditLogGetParams) (map[string]AuditLog, error) {
	return getMap[AuditLog](ctx, s.client, "auditlog.get", params)
}
