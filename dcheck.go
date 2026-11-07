package zabbix

import "context"

// DCheckService wraps the "dcheck" API namespace.
type DCheckService struct{ client *Client }

// DCheck returns the dcheck service.
func (c *Client) DCheck() *DCheckService { return &DCheckService{c} }

// DCheck object (discovery check).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/dcheck/object
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
type DCheckGetParams struct {
	GetParams
	DCheckIDs []string `json:"dcheckids,omitempty"`
	DRuleIDs  []string `json:"druleids,omitempty"`
}

// Get retrieves discovery checks matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/dcheck/get
func (s *DCheckService) Get(ctx context.Context, params DCheckGetParams) ([]DCheck, error) {
	var out []DCheck
	err := s.client.Call(ctx, "dcheck.get", params, &out)
	return out, err
}
