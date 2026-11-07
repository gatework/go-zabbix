package zabbix

import "context"

// ItemPrototypeService wraps the "itemprototype" API namespace.
type ItemPrototypeService struct{ client *Client }

// ItemPrototype returns the itemprototype service.
func (c *Client) ItemPrototype() *ItemPrototypeService { return &ItemPrototypeService{c} }

// ItemPrototype object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/itemprototype/object
type ItemPrototype struct {
	ItemID      string `json:"itemid,omitempty"`
	HostID      string `json:"hostid,omitempty"`
	Name        string `json:"name,omitempty"`
	Key         string `json:"key_,omitempty"`
	Type        string `json:"type,omitempty"`
	ValueType   string `json:"value_type,omitempty"`
	Delay       string `json:"delay,omitempty"`
	History     string `json:"history,omitempty"`
	Trends      string `json:"trends,omitempty"`
	Status      string `json:"status,omitempty"`
	Units       string `json:"units,omitempty"`
	ValueMapID  string `json:"valuemapid,omitempty"`
	Description string `json:"description,omitempty"`
	InterfaceID string `json:"interfaceid,omitempty"`
	TemplateID  string `json:"templateid,omitempty"`
	Flags       string `json:"flags,omitempty"`
	State       string `json:"state,omitempty"`
	Error       string `json:"error,omitempty"`
	LastClock   string `json:"lastclock,omitempty"`
	LastNs      string `json:"lastns,omitempty"`
	LastValue   string `json:"lastvalue,omitempty"`
	PrevValue   string `json:"prevvalue,omitempty"`
	SNMPOID     string `json:"snmp_oid,omitempty"`
	Params      string `json:"params,omitempty"`
	Master      string `json:"master_itemid,omitempty"`
	Discover    string `json:"discover,omitempty"`
	RuleID      string `json:"ruleid,omitempty"`
	Tags        []Tag  `json:"tags,omitempty"`

	Preprocessing []ItemPrototypePreprocessing `json:"preprocessing,omitempty"`
}

// ItemPrototypePreprocessing is a preprocessing step for an item prototype.
type ItemPrototypePreprocessing struct {
	Type               string `json:"type,omitempty"`
	Params             string `json:"params,omitempty"`
	ErrorHandler       string `json:"error_handler,omitempty"`
	ErrorHandlerParams string `json:"error_handler_params,omitempty"`
}

// ItemPrototypeGetParams are the parameters for itemprototype.get.
type ItemPrototypeGetParams struct {
	GetParams
	ItemIDs             []string `json:"itemids,omitempty"`
	HostIDs             []string `json:"hostids,omitempty"`
	GroupIDs            []string `json:"groupids,omitempty"`
	TemplateIDs         []string `json:"templateids,omitempty"`
	DiscoveryIDs        []string `json:"discoveryids,omitempty"`
	InterfaceIDs        []string `json:"interfaceids,omitempty"`
	SelectHosts         any      `json:"selectHosts,omitempty"`
	SelectTags          any      `json:"selectTags,omitempty"`
	SelectPreprocessing any      `json:"selectPreprocessing,omitempty"`
	SelectDiscoveryRule any      `json:"selectDiscoveryRule,omitempty"`
}

// Get retrieves item prototypes matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/itemprototype/get
func (s *ItemPrototypeService) Get(ctx context.Context, params ItemPrototypeGetParams) ([]ItemPrototype, error) {
	var out []ItemPrototype
	err := s.client.Call(ctx, "itemprototype.get", params, &out)
	return out, err
}

// Create creates item prototypes and returns the new ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/itemprototype/create
func (s *ItemPrototypeService) Create(ctx context.Context, items ...ItemPrototype) ([]string, error) {
	var res struct {
		IDs []string `json:"itemids"`
	}
	err := s.client.Call(ctx, "itemprototype.create", items, &res)
	return res.IDs, err
}

// Update updates item prototypes and returns the affected ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/itemprototype/update
func (s *ItemPrototypeService) Update(ctx context.Context, items ...ItemPrototype) ([]string, error) {
	var res struct {
		IDs []string `json:"itemids"`
	}
	err := s.client.Call(ctx, "itemprototype.update", items, &res)
	return res.IDs, err
}

// Delete deletes item prototypes by id and returns the deleted ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/itemprototype/delete
func (s *ItemPrototypeService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"itemids"`
	}
	err := s.client.Call(ctx, "itemprototype.delete", ids, &res)
	return res.IDs, err
}
