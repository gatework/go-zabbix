package zabbix

import "context"

// SettingsService wraps the "settings" API namespace, a global singleton with only
// get and update.
type SettingsService struct{ client *Client }

// Settings returns the settings service.
func (c *Client) Settings() *SettingsService { return &SettingsService{c} }

// Settings object (global GUI/frontend settings). This is a representative subset of
// the full parameter list.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/settings/object
type Settings struct {
	DefaultLang          string `json:"default_lang,omitempty"`
	DefaultTimezone      string `json:"default_timezone,omitempty"`
	DefaultTheme         string `json:"default_theme,omitempty"`
	SearchLimit          string `json:"search_limit,omitempty"`
	SeverityName0        string `json:"severity_name_0,omitempty"`
	SeverityName1        string `json:"severity_name_1,omitempty"`
	SeverityName2        string `json:"severity_name_2,omitempty"`
	SeverityName3        string `json:"severity_name_3,omitempty"`
	SeverityName4        string `json:"severity_name_4,omitempty"`
	SeverityName5        string `json:"severity_name_5,omitempty"`
	BlinkPeriod          string `json:"blink_period,omitempty"`
	WorkPeriod           string `json:"work_period,omitempty"`
	ShowTechnicalErrors  string `json:"show_technical_errors,omitempty"`
	HistoryPeriod        string `json:"history_period,omitempty"`
	PeriodDefault        string `json:"period_default,omitempty"`
	MaxPeriod            string `json:"max_period,omitempty"`
	SNMPTrapLogging      string `json:"snmptrap_logging,omitempty"`
	DefaultInventoryMode string `json:"default_inventory_mode,omitempty"`
	GeomapsTileProvider  string `json:"geomaps_tile_provider,omitempty"`
}

// Get retrieves the global settings.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/settings/get
func (s *SettingsService) Get(ctx context.Context) (*Settings, error) {
	var out Settings
	err := s.client.Call(ctx, "settings.get", map[string]any{"output": "extend"}, &out)
	return &out, err
}

// Update updates the global settings and returns the updated parameter names.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/settings/update
func (s *SettingsService) Update(ctx context.Context, params Settings) ([]string, error) {
	var res []string
	err := s.client.Call(ctx, "settings.update", params, &res)
	return res, err
}
