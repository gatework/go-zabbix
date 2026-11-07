package zabbix

import "context"

// ConfigurationService wraps the "configuration" API namespace.
type ConfigurationService struct{ client *Client }

// Configuration returns the configuration service.
func (c *Client) Configuration() *ConfigurationService { return &ConfigurationService{c} }

// ConfigurationExportParams are the parameters for configuration.export.
type ConfigurationExportParams struct {
	Format      string `json:"format"`
	Prettyprint bool   `json:"prettyprint,omitempty"`
	Options     any    `json:"options"`
}

// ConfigurationImportParams are the parameters for configuration.import and
// configuration.importcompare.
type ConfigurationImportParams struct {
	Format string `json:"format"`
	Source string `json:"source"`
	Rules  any    `json:"rules"`
}

// Export exports configuration in the requested format and returns the serialized
// configuration as a string.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/configuration/export
func (s *ConfigurationService) Export(ctx context.Context, params ConfigurationExportParams) (string, error) {
	var out string
	err := s.client.Call(ctx, "configuration.export", params, &out)
	return out, err
}

// Import imports the given configuration and reports whether it succeeded.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/configuration/import
func (s *ConfigurationService) Import(ctx context.Context, params ConfigurationImportParams) (bool, error) {
	var out bool
	err := s.client.Call(ctx, "configuration.import", params, &out)
	return out, err
}

// ImportCompare compares the given configuration against the current one and
// returns the differences.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/configuration/importcompare
func (s *ConfigurationService) ImportCompare(ctx context.Context, params ConfigurationImportParams) (any, error) {
	var out any
	err := s.client.Call(ctx, "configuration.importcompare", params, &out)
	return out, err
}
