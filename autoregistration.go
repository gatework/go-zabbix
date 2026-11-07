package zabbix

import "context"

// AutoregistrationService wraps the "autoregistration" API namespace, a global
// singleton with only get and update.
type AutoregistrationService struct{ client *Client }

// Autoregistration returns the autoregistration service.
func (c *Client) Autoregistration() *AutoregistrationService { return &AutoregistrationService{c} }

// Autoregistration object (global active-agent autoregistration settings). On update
// the TLSPSKIdentity and TLSPSK fields may also be supplied.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/autoregistration/object
type Autoregistration struct {
	TLSAccept      string `json:"tls_accept,omitempty"`
	TLSPSKIdentity string `json:"tls_psk_identity,omitempty"`
	TLSPSK         string `json:"tls_psk,omitempty"`
}

// Get retrieves the global autoregistration settings.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/autoregistration/get
func (s *AutoregistrationService) Get(ctx context.Context) (*Autoregistration, error) {
	var out Autoregistration
	err := s.client.Call(ctx, "autoregistration.get", map[string]any{"output": "extend"}, &out)
	return &out, err
}

// Update updates the global autoregistration settings and returns the updated
// parameter names.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/autoregistration/update
func (s *AutoregistrationService) Update(ctx context.Context, params Autoregistration) ([]string, error) {
	var res []string
	err := s.client.Call(ctx, "autoregistration.update", params, &res)
	return res, err
}
