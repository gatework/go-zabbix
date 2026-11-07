package zabbix

import "context"

// APIInfoService wraps the "apiinfo" API namespace.
type APIInfoService struct{ client *Client }

// APIInfo returns the apiinfo service.
func (c *Client) APIInfo() *APIInfoService { return &APIInfoService{c} }

// Version returns the Zabbix API version, e.g. "7.4.0".
//
// Per the Zabbix API, apiinfo.version must be called without authentication, so
// this method never sends the Authorization header.
//
// https://www.zabbix.com/documentation/current/en/manual/api/reference/apiinfo/version
func (s *APIInfoService) Version(ctx context.Context) (string, error) {
	var version string
	err := s.client.do(ctx, "apiinfo.version", []any{}, &version, false)
	return version, err
}
