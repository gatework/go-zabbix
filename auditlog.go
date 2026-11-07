package zabbix

import "context"

// AuditLogService wraps the "auditlog" API namespace.
type AuditLogService struct{ client *Client }

// AuditLog returns the auditlog service.
func (c *Client) AuditLog() *AuditLogService { return &AuditLogService{c} }

// AuditLog is a single audit log record.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/auditlog/object
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
type AuditLogGetParams struct {
	GetParams
	UserIDs  []string `json:"userids,omitempty"`
	TimeFrom int64    `json:"time_from,omitempty"`
	TimeTill int64    `json:"time_till,omitempty"`
}

// Get retrieves audit log records matching params. The auditlog API is read-only.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/auditlog/get
func (s *AuditLogService) Get(ctx context.Context, params AuditLogGetParams) ([]AuditLog, error) {
	var out []AuditLog
	err := s.client.Call(ctx, "auditlog.get", params, &out)
	return out, err
}
