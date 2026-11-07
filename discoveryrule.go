package zabbix

import "context"

// DiscoveryRuleService wraps the "discoveryrule" API namespace (LLD rules).
type DiscoveryRuleService struct{ client *Client }

// DiscoveryRule returns the discoveryrule service.
func (c *Client) DiscoveryRule() *DiscoveryRuleService { return &DiscoveryRuleService{c} }

// DiscoveryRule object (low-level discovery rule).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/discoveryrule/object
type DiscoveryRule struct {
	ItemID              string                              `json:"itemid,omitempty"`
	HostID              string                              `json:"hostid,omitempty"`
	Name                string                              `json:"name,omitempty"`
	Key                 string                              `json:"key_,omitempty"`
	Type                string                              `json:"type,omitempty"`
	Delay               string                              `json:"delay,omitempty"`
	Status              string                              `json:"status,omitempty"`
	State               string                              `json:"state,omitempty"`
	Error               string                              `json:"error,omitempty"`
	Lifetime            string                              `json:"lifetime,omitempty"`
	LifetimeType        string                              `json:"lifetime_type,omitempty"`
	EnabledLifetime     string                              `json:"enabled_lifetime,omitempty"`
	EnabledLifetimeType string                              `json:"enabled_lifetime_type,omitempty"`
	Description         string                              `json:"description,omitempty"`
	TemplateID          string                              `json:"templateid,omitempty"`
	InterfaceID         string                              `json:"interfaceid,omitempty"`
	Params              string                              `json:"params,omitempty"`
	SNMPOID             string                              `json:"snmp_oid,omitempty"`
	MasterItemID        string                              `json:"master_itemid,omitempty"`
	URL                 string                              `json:"url,omitempty"`
	IPMISensor          string                              `json:"ipmi_sensor,omitempty"`
	JMXEndpoint         string                              `json:"jmx_endpoint,omitempty"`
	Timeout             string                              `json:"timeout,omitempty"`
	TrapperHosts        string                              `json:"trapper_hosts,omitempty"`
	Username            string                              `json:"username,omitempty"`
	Password            string                              `json:"password,omitempty"`
	AuthType            string                              `json:"authtype,omitempty"`
	Filter              SelectedObject[DiscoveryRuleFilter] `json:"filter,omitzero"`
	LLDMacroPaths       []DiscoveryRuleLLDMacroPath         `json:"lld_macro_paths,omitempty"`
	Preprocessing       []DiscoveryRulePreprocessing        `json:"preprocessing,omitempty"`
	Hosts               []Host                              `json:"hosts,omitempty"`
}

// DiscoveryRuleLLDMacroPath maps an LLD macro to a JSONPath.
type DiscoveryRuleLLDMacroPath struct {
	LLDMacro string `json:"lld_macro,omitempty"`
	Path     string `json:"path,omitempty"`
}

// DiscoveryRuleFilter is a low-level discovery filter selected by selectFilter.
type DiscoveryRuleFilter struct {
	EvalType    string                         `json:"evaltype,omitempty"`
	Formula     string                         `json:"formula,omitempty"`
	EvalFormula string                         `json:"eval_formula,omitempty"`
	Conditions  []DiscoveryRuleFilterCondition `json:"conditions,omitempty"`
}

// DiscoveryRuleFilterCondition tests one low-level discovery macro.
type DiscoveryRuleFilterCondition struct {
	Macro     string `json:"macro,omitempty"`
	Value     string `json:"value,omitempty"`
	FormulaID string `json:"formulaid,omitempty"`
	Operator  string `json:"operator,omitempty"`
}

// DiscoveryRulePreprocessing is a preprocessing step for an LLD rule.
type DiscoveryRulePreprocessing struct {
	Type               string `json:"type,omitempty"`
	Params             string `json:"params,omitempty"`
	ErrorHandler       string `json:"error_handler,omitempty"`
	ErrorHandlerParams string `json:"error_handler_params,omitempty"`
}

// DiscoveryRuleGetParams are the parameters for discoveryrule.get.
// DiscoveryRuleGetParams queries the list form of discoveryrule.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type DiscoveryRuleGetParams struct {
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
	ItemIDs                []string       `json:"itemids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	TemplateIDs            []string       `json:"templateids,omitzero"`
	InterfaceIDs           []string       `json:"interfaceids,omitzero"`
	Inherited              *bool          `json:"inherited,omitempty"`
	Templated              *bool          `json:"templated,omitempty"`
	Monitored              *bool          `json:"monitored,omitempty"`
	SelectHosts            any            `json:"selectHosts,omitempty"`
	SelectFilter           any            `json:"selectFilter,omitempty"`
	SelectLLDMacroPaths    any            `json:"selectLLDMacroPaths,omitempty"`
	SelectPreprocessing    any            `json:"selectPreprocessing,omitempty"`
}

// Get retrieves LLD rules matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/discoveryrule/get
func (s *DiscoveryRuleService) Get(ctx context.Context, params DiscoveryRuleGetParams) ([]DiscoveryRule, error) {
	return getList[DiscoveryRule](ctx, s.client, "discoveryrule.get", params)
}

// DiscoveryRuleCountParams contains the filters for discoveryrule.get counts.
type DiscoveryRuleCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	ItemIDs                []string       `json:"itemids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	TemplateIDs            []string       `json:"templateids,omitzero"`
	InterfaceIDs           []string       `json:"interfaceids,omitzero"`
	Inherited              *bool          `json:"inherited,omitempty"`
	Templated              *bool          `json:"templated,omitempty"`
	Monitored              *bool          `json:"monitored,omitempty"`
}

// Count returns the number of objects matching params.
func (s *DiscoveryRuleService) Count(ctx context.Context, params DiscoveryRuleCountParams) (uint64, error) {
	return getCount(ctx, s.client, "discoveryrule.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *DiscoveryRuleService) GetMap(ctx context.Context, params DiscoveryRuleGetParams) (map[string]DiscoveryRule, error) {
	return getMap[DiscoveryRule](ctx, s.client, "discoveryrule.get", params)
}

// Create creates LLD rules and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/discoveryrule/create
func (s *DiscoveryRuleService) Create(ctx context.Context, rules ...DiscoveryRuleCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "discoveryrule.create", "itemids", rules)
}

// Update updates LLD rules and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/discoveryrule/update
func (s *DiscoveryRuleService) Update(ctx context.Context, rules ...DiscoveryRuleUpdateParams) ([]string, error) {
	return updateIDs(ctx, s.client, "discoveryrule.update", "itemids", rules)
}

// Delete deletes LLD rules by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/discoveryrule/delete
func (s *DiscoveryRuleService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "discoveryrule.delete", "ruleids", ids)
}
