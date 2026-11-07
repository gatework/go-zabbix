package zabbix

import "context"

// ItemService wraps the "item" API namespace.
type ItemService struct{ client *Client }

// Item returns the item service.
func (c *Client) Item() *ItemService { return &ItemService{c} }

// Item object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/item/object
type Item struct {
	ItemID       string            `json:"itemid,omitempty"`
	HostID       string            `json:"hostid,omitempty"`
	Name         string            `json:"name,omitempty"`
	Key          string            `json:"key_,omitempty"`
	Type         string            `json:"type,omitempty"`
	ValueType    string            `json:"value_type,omitempty"`
	Delay        string            `json:"delay,omitempty"`
	History      string            `json:"history,omitempty"`
	Trends       string            `json:"trends,omitempty"`
	Status       string            `json:"status,omitempty"`
	Units        string            `json:"units,omitempty"`
	ValueMapID   string            `json:"valuemapid,omitempty"`
	Description  string            `json:"description,omitempty"`
	InterfaceID  string            `json:"interfaceid,omitempty"`
	TemplateID   string            `json:"templateid,omitempty"`
	Flags        string            `json:"flags,omitempty"`
	State        string            `json:"state,omitempty"`
	Error        string            `json:"error,omitempty"`
	LastClock    string            `json:"lastclock,omitempty"`
	LastNS       string            `json:"lastns,omitempty"`
	LastValue    string            `json:"lastvalue,omitempty"`
	PrevValue    string            `json:"prevvalue,omitempty"`
	SNMPOID      string            `json:"snmp_oid,omitempty"`
	Params       string            `json:"params,omitempty"`
	MasterItemID string            `json:"master_itemid,omitempty"`
	Tags         []Tag             `json:"tags,omitempty"`
	Hosts        []Host            `json:"hosts,omitempty"`
	Triggers     Relation[Trigger] `json:"triggers,omitzero"`
}

// ItemGetParams are the parameters for item.get.
// ItemGetParams queries the list form of item.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type ItemGetParams struct {
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
	TriggerIDs             []string       `json:"triggerids,omitzero"`
	Tags                   []TagFilter    `json:"tags,omitzero"`
	EvalType               int            `json:"evaltype,omitempty"`
	Templated              *bool          `json:"templated,omitempty"`
	Monitored              *bool          `json:"monitored,omitempty"`
	WebItems               bool           `json:"webitems,omitempty"`
	SelectHosts            any            `json:"selectHosts,omitempty"`
	SelectTriggers         any            `json:"selectTriggers,omitempty"`
	SelectTags             any            `json:"selectTags,omitempty"`
}

// Get retrieves items matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/item/get
func (s *ItemService) Get(ctx context.Context, params ItemGetParams) ([]Item, error) {
	return getList[Item](ctx, s.client, "item.get", params)
}

// ItemCountParams contains the filters for item.get counts.
type ItemCountParams struct {
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
	TriggerIDs             []string       `json:"triggerids,omitzero"`
	Tags                   []TagFilter    `json:"tags,omitzero"`
	EvalType               int            `json:"evaltype,omitempty"`
	Templated              *bool          `json:"templated,omitempty"`
	Monitored              *bool          `json:"monitored,omitempty"`
	WebItems               bool           `json:"webitems,omitempty"`
}

// Count returns the number of objects matching params.
func (s *ItemService) Count(ctx context.Context, params ItemCountParams) (uint64, error) {
	return getCount(ctx, s.client, "item.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *ItemService) GetMap(ctx context.Context, params ItemGetParams) (map[string]Item, error) {
	return getMap[Item](ctx, s.client, "item.get", params)
}

// Create creates items and returns the new item ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/item/create
func (s *ItemService) Create(ctx context.Context, items ...ItemCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "item.create", "itemids", items)
}

// Update updates items and returns the affected item ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/item/update
func (s *ItemService) Update(ctx context.Context, items ...ItemUpdateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "item.update", "itemids", items)
}

// Delete deletes items by id and returns the deleted item ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/item/delete
func (s *ItemService) Delete(ctx context.Context, itemIDs ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "item.delete", "itemids", itemIDs)
}
