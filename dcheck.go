package zabbix

import "context"

// DCheckService wraps the "dcheck" API namespace.
type DCheckService struct{ client *Client }

// DCheck returns the dcheck service.
func (c *Client) DCheck() *DCheckService { return &DCheckService{c} }

// DCheck object (discovery check).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dcheck/object
type DCheck struct {
	DCheckID             string `json:"dcheckid,omitempty"`
	DRuleID              string `json:"druleid,omitempty"`
	Type                 string `json:"type,omitempty"`
	Key                  string `json:"key_,omitempty"`
	SNMPCommunity        string `json:"snmp_community,omitempty"`
	Ports                string `json:"ports,omitempty"`
	SNMPv3SecurityName   string `json:"snmpv3_securityname,omitempty"`
	SNMPv3SecurityLevel  string `json:"snmpv3_securitylevel,omitempty"`
	SNMPv3AuthPassphrase string `json:"snmpv3_authpassphrase,omitempty"`
	SNMPv3PrivPassphrase string `json:"snmpv3_privpassphrase,omitempty"`
	SNMPv3AuthProtocol   string `json:"snmpv3_authprotocol,omitempty"`
	SNMPv3PrivProtocol   string `json:"snmpv3_privprotocol,omitempty"`
	SNMPv3ContextName    string `json:"snmpv3_contextname,omitempty"`
	Uniq                 string `json:"uniq,omitempty"`
	HostSource           string `json:"host_source,omitempty"`
	NameSource           string `json:"name_source,omitempty"`
	AllowRedirect        string `json:"allow_redirect,omitempty"`
}

// DCheckGetParams are the parameters for dcheck.get.
// DCheckGetParams queries the list form of dcheck.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type DCheckGetParams struct {
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
	Editable               bool           `json:"editable,omitempty"`
	DCheckIDs              []string       `json:"dcheckids,omitzero"`
	DRuleIDs               []string       `json:"druleids,omitzero"`
}

// Get retrieves discovery checks matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dcheck/get
func (s *DCheckService) Get(ctx context.Context, params DCheckGetParams) ([]DCheck, error) {
	return getList[DCheck](ctx, s.client, "dcheck.get", params)
}

// DCheckCountParams contains the filters for dcheck.get counts.
type DCheckCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	DCheckIDs              []string       `json:"dcheckids,omitzero"`
	DRuleIDs               []string       `json:"druleids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *DCheckService) Count(ctx context.Context, params DCheckCountParams) (uint64, error) {
	return getCount(ctx, s.client, "dcheck.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *DCheckService) GetMap(ctx context.Context, params DCheckGetParams) (map[string]DCheck, error) {
	return getMap[DCheck](ctx, s.client, "dcheck.get", params)
}
