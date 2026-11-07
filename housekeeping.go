package zabbix

import "context"

// HousekeepingService wraps the "housekeeping" API namespace, a global singleton with
// only get and update.
type HousekeepingService struct{ client *Client }

// Housekeeping returns the housekeeping service.
func (c *Client) Housekeeping() *HousekeepingService { return &HousekeepingService{c} }

// Housekeeping object (global housekeeping/retention settings).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/housekeeping/object
type Housekeeping struct {
	HKEventsMode      string `json:"hk_events_mode,omitempty"`
	HKEventsTrigger   string `json:"hk_events_trigger,omitempty"`
	HKEventsInternal  string `json:"hk_events_internal,omitempty"`
	HKEventsDiscovery string `json:"hk_events_discovery,omitempty"`
	HKEventsAutoreg   string `json:"hk_events_autoreg,omitempty"`
	HKServicesMode    string `json:"hk_services_mode,omitempty"`
	HKServices        string `json:"hk_services,omitempty"`
	HKAuditMode       string `json:"hk_audit_mode,omitempty"`
	HKAudit           string `json:"hk_audit,omitempty"`
	HKSessionsMode    string `json:"hk_sessions_mode,omitempty"`
	HKSessions        string `json:"hk_sessions,omitempty"`
	HKHistoryMode     string `json:"hk_history_mode,omitempty"`
	HKHistoryGlobal   string `json:"hk_history_global,omitempty"`
	HKHistory         string `json:"hk_history,omitempty"`
	HKTrendsMode      string `json:"hk_trends_mode,omitempty"`
	HKTrendsGlobal    string `json:"hk_trends_global,omitempty"`
	HKTrends          string `json:"hk_trends,omitempty"`
	CompressionStatus string `json:"compression_status,omitempty"`
	CompressOlder     string `json:"compress_older,omitempty"`
}

// Get retrieves the global housekeeping settings.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/housekeeping/get
func (s *HousekeepingService) Get(ctx context.Context) (*Housekeeping, error) {
	var out Housekeeping
	err := s.client.Call(ctx, "housekeeping.get", map[string]any{"output": "extend"}, &out)
	return &out, err
}

// Update updates the global housekeeping settings and returns the updated parameter
// names.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/housekeeping/update
func (s *HousekeepingService) Update(ctx context.Context, params Housekeeping) ([]string, error) {
	var res []string
	err := s.client.Call(ctx, "housekeeping.update", params, &res)
	return res, err
}
