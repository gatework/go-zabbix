package zabbix

import "context"

// ItemPrototypeService wraps the "itemprototype" API namespace.
type ItemPrototypeService struct{ client *Client }

// ItemPrototype returns the itemprototype service.
func (c *Client) ItemPrototype() *ItemPrototypeService { return &ItemPrototypeService{c} }

// ItemPrototype object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/itemprototype/object
type ItemPrototype struct {
	ItemID       string `json:"itemid,omitempty"`
	HostID       string `json:"hostid,omitempty"`
	Name         string `json:"name,omitempty"`
	Key          string `json:"key_,omitempty"`
	Type         string `json:"type,omitempty"`
	ValueType    string `json:"value_type,omitempty"`
	Delay        string `json:"delay,omitempty"`
	History      string `json:"history,omitempty"`
	Trends       string `json:"trends,omitempty"`
	Status       string `json:"status,omitempty"`
	Units        string `json:"units,omitempty"`
	ValueMapID   string `json:"valuemapid,omitempty"`
	Description  string `json:"description,omitempty"`
	InterfaceID  string `json:"interfaceid,omitempty"`
	TemplateID   string `json:"templateid,omitempty"`
	Flags        string `json:"flags,omitempty"`
	State        string `json:"state,omitempty"`
	Error        string `json:"error,omitempty"`
	LastClock    string `json:"lastclock,omitempty"`
	LastNS       string `json:"lastns,omitempty"`
	LastValue    string `json:"lastvalue,omitempty"`
	PrevValue    string `json:"prevvalue,omitempty"`
	SNMPOID      string `json:"snmp_oid,omitempty"`
	Params       string `json:"params,omitempty"`
	MasterItemID string `json:"master_itemid,omitempty"`
	Discover     string `json:"discover,omitempty"`
	RuleID       string `json:"ruleid,omitempty"`
	Tags         []Tag  `json:"tags,omitempty"`

	Preprocessing []ItemPrototypePreprocessing  `json:"preprocessing,omitempty"`
	Hosts         []Host                        `json:"hosts,omitempty"`
	DiscoveryRule SelectedObject[DiscoveryRule] `json:"discoveryRule,omitzero"`
}

// ItemPrototypePreprocessing is a preprocessing step for an item prototype.
type ItemPrototypePreprocessing struct {
	Type               string `json:"type,omitempty"`
	Params             string `json:"params,omitempty"`
	ErrorHandler       string `json:"error_handler,omitempty"`
	ErrorHandlerParams string `json:"error_handler_params,omitempty"`
}

// ItemPrototypeGetParams are the parameters for itemprototype.get.
// ItemPrototypeGetParams queries the list form of itemprototype.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type ItemPrototypeGetParams struct {
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
	DiscoveryIDs           []string       `json:"discoveryids,omitzero"`
	SelectHosts            any            `json:"selectHosts,omitempty"`
	SelectTags             any            `json:"selectTags,omitempty"`
	SelectPreprocessing    any            `json:"selectPreprocessing,omitempty"`
	SelectDiscoveryRule    any            `json:"selectDiscoveryRule,omitempty"`
}

// Get retrieves item prototypes matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/itemprototype/get
func (s *ItemPrototypeService) Get(ctx context.Context, params ItemPrototypeGetParams) ([]ItemPrototype, error) {
	return getList[ItemPrototype](ctx, s.client, "itemprototype.get", params)
}

// ItemPrototypeCountParams contains the filters for itemprototype.get counts.
type ItemPrototypeCountParams struct {
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
	DiscoveryIDs           []string       `json:"discoveryids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *ItemPrototypeService) Count(ctx context.Context, params ItemPrototypeCountParams) (uint64, error) {
	return getCount(ctx, s.client, "itemprototype.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *ItemPrototypeService) GetMap(ctx context.Context, params ItemPrototypeGetParams) (map[string]ItemPrototype, error) {
	return getMap[ItemPrototype](ctx, s.client, "itemprototype.get", params)
}

// Create creates item prototypes and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/itemprototype/create
func (s *ItemPrototypeService) Create(ctx context.Context, items ...ItemPrototypeCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "itemprototype.create", "itemids", items)
}

// Update updates item prototypes and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/itemprototype/update
func (s *ItemPrototypeService) Update(ctx context.Context, items ...ItemPrototypeUpdateParams) ([]string, error) {
	return updateIDs(ctx, s.client, "itemprototype.update", "itemids", items)
}

// Delete deletes item prototypes by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/itemprototype/delete
func (s *ItemPrototypeService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "itemprototype.delete", "prototypeids", ids)
}
