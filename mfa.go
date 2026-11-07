package zabbix

import "context"

// MFAService wraps the "mfa" API namespace.
type MFAService struct{ client *Client }

// MFA returns the mfa service.
func (c *Client) MFA() *MFAService { return &MFAService{c} }

// MFA object (multi-factor authentication method).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/mfa/object
type MFA struct {
	MFAID        string `json:"mfaid,omitempty"`
	Type         string `json:"type,omitempty"`
	Name         string `json:"name,omitempty"`
	HashFunction string `json:"hash_function,omitempty"`
	CodeLength   string `json:"code_length,omitempty"`
	APIHostname  string `json:"api_hostname,omitempty"`
	ClientID     string `json:"clientid,omitempty"`
}

// MFAGetParams are the parameters for mfa.get.
type MFAGetParams struct {
	GetParams
	MFAIDs []string `json:"mfaids,omitempty"`
}

// Get retrieves MFA methods matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/mfa/get
func (s *MFAService) Get(ctx context.Context, params MFAGetParams) ([]MFA, error) {
	var out []MFA
	err := s.client.Call(ctx, "mfa.get", params, &out)
	return out, err
}

// Create creates MFA methods and returns the new ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/mfa/create
func (s *MFAService) Create(ctx context.Context, items ...MFA) ([]string, error) {
	var res struct {
		IDs []string `json:"mfaids"`
	}
	err := s.client.Call(ctx, "mfa.create", items, &res)
	return res.IDs, err
}

// Update updates MFA methods and returns the affected ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/mfa/update
func (s *MFAService) Update(ctx context.Context, items ...MFA) ([]string, error) {
	var res struct {
		IDs []string `json:"mfaids"`
	}
	err := s.client.Call(ctx, "mfa.update", items, &res)
	return res.IDs, err
}

// Delete deletes MFA methods by id and returns the deleted ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/mfa/delete
func (s *MFAService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"mfaids"`
	}
	err := s.client.Call(ctx, "mfa.delete", ids, &res)
	return res.IDs, err
}
