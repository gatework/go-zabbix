package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
)

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

// Export returns serialized configuration in json, xml, or yaml format.
// Use ExportRaw to receive the dynamic configuration object in raw format.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/configuration/export
func (s *ConfigurationService) Export(ctx context.Context, params ConfigurationExportParams) (string, error) {
	switch params.Format {
	case "json", "xml", "yaml":
	default:
		return "", requestError("configuration.export", ErrorInvalidParams, errors.New("format must be json, xml, or yaml; use ExportRaw for raw"))
	}
	var out string
	err := s.client.Call(ctx, "configuration.export", params, &out)
	return out, err
}

// ExportRaw returns the dynamic configuration object as JSON without decoding its
// format-dependent fields. options has the same shape as Export's Options field.
// Zabbix 7.4 CConfiguration::validateExport accepts raw, and CRawExportWriter::write
// returns the configuration object rather than a serialized string.
func (s *ConfigurationService) ExportRaw(ctx context.Context, options any) (json.RawMessage, error) {
	params := ConfigurationExportParams{Format: "raw", Options: options}
	var raw json.RawMessage
	if err := s.client.Call(ctx, "configuration.export", params, &raw); err != nil {
		return nil, err
	}
	if result := bytes.TrimSpace(raw); len(result) == 0 || result[0] != '{' {
		return nil, requestError("configuration.export", ErrorProtocol, errors.New("raw export result must be an object"))
	}
	return raw, nil
}

// Import imports the given configuration and reports whether it succeeded.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/configuration/import
func (s *ConfigurationService) Import(ctx context.Context, params ConfigurationImportParams) (bool, error) {
	var out bool
	err := s.client.Call(ctx, "configuration.import", params, &out)
	return out, err
}

// ImportCompare compares the given configuration against the current one and
// returns the differences.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/configuration/importcompare
func (s *ConfigurationService) ImportCompare(ctx context.Context, params ConfigurationImportParams) (any, error) {
	var out any
	err := s.client.Call(ctx, "configuration.importcompare", params, &out)
	return out, err
}
