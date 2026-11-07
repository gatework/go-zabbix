package zabbix

import "context"

// MFAService wraps the "mfa" API namespace.
type MFAService struct{ client *Client }

// MFA returns the mfa service.
func (c *Client) MFA() *MFAService { return &MFAService{c} }

// MFA object (multi-factor authentication method).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mfa/object
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
// MFAGetParams queries the list form of mfa.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type MFAGetParams struct {
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
	MFAIDs                 []string       `json:"mfaids,omitzero"`
}

// Get retrieves MFA methods matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mfa/get
func (s *MFAService) Get(ctx context.Context, params MFAGetParams) ([]MFA, error) {
	return getList[MFA](ctx, s.client, "mfa.get", params)
}

// MFACountParams contains the filters for mfa.get counts.
type MFACountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	MFAIDs                 []string       `json:"mfaids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *MFAService) Count(ctx context.Context, params MFACountParams) (uint64, error) {
	return getCount(ctx, s.client, "mfa.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *MFAService) GetMap(ctx context.Context, params MFAGetParams) (map[string]MFA, error) {
	return getMap[MFA](ctx, s.client, "mfa.get", params)
}

// Create creates MFA methods and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mfa/create
func (s *MFAService) Create(ctx context.Context, items ...MFA) ([]string, error) {
	var res struct {
		IDs []string `json:"mfaids"`
	}
	err := s.client.Call(ctx, "mfa.create", items, &res)
	return res.IDs, err
}

// Update updates MFA methods and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mfa/update
func (s *MFAService) Update(ctx context.Context, items ...MFA) ([]string, error) {
	var res struct {
		IDs []string `json:"mfaids"`
	}
	err := s.client.Call(ctx, "mfa.update", items, &res)
	return res.IDs, err
}

// Delete deletes MFA methods by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mfa/delete
func (s *MFAService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"mfaids"`
	}
	err := s.client.Call(ctx, "mfa.delete", ids, &res)
	return res.IDs, err
}
