package zabbix

import "context"

// DiscoveryRuleService wraps the "discoveryrule" API namespace (LLD rules).
type DiscoveryRuleService struct{ client *Client }

// DiscoveryRule returns the discoveryrule service.
func (c *Client) DiscoveryRule() *DiscoveryRuleService { return &DiscoveryRuleService{c} }

// DiscoveryRule object (low-level discovery rule).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/discoveryrule/object
type DiscoveryRule struct {
	ItemID              string `json:"itemid,omitempty"`
	HostID              string `json:"hostid,omitempty"`
	Name                string `json:"name,omitempty"`
	Key                 string `json:"key_,omitempty"`
	Type                string `json:"type,omitempty"`
	Delay               string `json:"delay,omitempty"`
	Status              string `json:"status,omitempty"`
	State               string `json:"state,omitempty"`
	Error               string `json:"error,omitempty"`
	Lifetime            string `json:"lifetime,omitempty"`
	LifetimeType        string `json:"lifetime_type,omitempty"`
	EnabledLifetime     string `json:"enabled_lifetime,omitempty"`
	EnabledLifetimeType string `json:"enabled_lifetime_type,omitempty"`
	Description         string `json:"description,omitempty"`
	TemplateID          string `json:"templateid,omitempty"`
	InterfaceID         string `json:"interfaceid,omitempty"`
	Params              string `json:"params,omitempty"`
	SNMPOID             string `json:"snmp_oid,omitempty"`
	MasterItemID        string `json:"master_itemid,omitempty"`
	URL                 string `json:"url,omitempty"`
	IPMISensor          string `json:"ipmi_sensor,omitempty"`
	JMXEndpoint         string `json:"jmx_endpoint,omitempty"`
	Timeout             string `json:"timeout,omitempty"`
	TrapperHosts        string `json:"trapper_hosts,omitempty"`
	Username            string `json:"username,omitempty"`
	Password            string `json:"password,omitempty"`
	AuthType            string `json:"authtype,omitempty"`

	Filter        any                          `json:"filter,omitempty"`
	LLDMacroPaths []DiscoveryRuleLLDMacroPath  `json:"lld_macro_paths,omitempty"`
	Preprocessing []DiscoveryRulePreprocessing `json:"preprocessing,omitempty"`
}

// DiscoveryRuleLLDMacroPath maps an LLD macro to a JSONPath.
type DiscoveryRuleLLDMacroPath struct {
	LLDMacro string `json:"lld_macro,omitempty"`
	Path     string `json:"path,omitempty"`
}

// DiscoveryRulePreprocessing is a preprocessing step for an LLD rule.
type DiscoveryRulePreprocessing struct {
	Type               string `json:"type,omitempty"`
	Params             string `json:"params,omitempty"`
	ErrorHandler       string `json:"error_handler,omitempty"`
	ErrorHandlerParams string `json:"error_handler_params,omitempty"`
}

// DiscoveryRuleGetParams are the parameters for discoveryrule.get.
type DiscoveryRuleGetParams struct {
	GetParams
	ItemIDs             []string `json:"itemids,omitempty"`
	HostIDs             []string `json:"hostids,omitempty"`
	GroupIDs            []string `json:"groupids,omitempty"`
	TemplateIDs         []string `json:"templateids,omitempty"`
	InterfaceIDs        []string `json:"interfaceids,omitempty"`
	Inherited           *bool    `json:"inherited,omitempty"`
	Templated           *bool    `json:"templated,omitempty"`
	Monitored           *bool    `json:"monitored,omitempty"`
	SelectHosts         any      `json:"selectHosts,omitempty"`
	SelectFilter        any      `json:"selectFilter,omitempty"`
	SelectLLDMacroPaths any      `json:"selectLLDMacroPaths,omitempty"`
	SelectPreprocessing any      `json:"selectPreprocessing,omitempty"`
}

// Get retrieves LLD rules matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/discoveryrule/get
func (s *DiscoveryRuleService) Get(ctx context.Context, params DiscoveryRuleGetParams) ([]DiscoveryRule, error) {
	var out []DiscoveryRule
	err := s.client.Call(ctx, "discoveryrule.get", params, &out)
	return out, err
}

// Create creates LLD rules and returns the new ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/discoveryrule/create
func (s *DiscoveryRuleService) Create(ctx context.Context, rules ...DiscoveryRule) ([]string, error) {
	var res struct {
		IDs []string `json:"itemids"`
	}
	err := s.client.Call(ctx, "discoveryrule.create", rules, &res)
	return res.IDs, err
}

// Update updates LLD rules and returns the affected ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/discoveryrule/update
func (s *DiscoveryRuleService) Update(ctx context.Context, rules ...DiscoveryRule) ([]string, error) {
	var res struct {
		IDs []string `json:"itemids"`
	}
	err := s.client.Call(ctx, "discoveryrule.update", rules, &res)
	return res.IDs, err
}

// Delete deletes LLD rules by id and returns the deleted ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/discoveryrule/delete
func (s *DiscoveryRuleService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"itemids"`
	}
	err := s.client.Call(ctx, "discoveryrule.delete", ids, &res)
	return res.IDs, err
}
