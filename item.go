package zabbix

import "context"

// ItemService wraps the "item" API namespace.
type ItemService struct{ client *Client }

// Item returns the item service.
func (c *Client) Item() *ItemService { return &ItemService{c} }

// Item object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/item/object
type Item struct {
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
	Tags        []Tag  `json:"tags,omitempty"`
}

// ItemGetParams are the parameters for item.get.
type ItemGetParams struct {
	GetParams
	ItemIDs        []string `json:"itemids,omitempty"`
	HostIDs        []string `json:"hostids,omitempty"`
	GroupIDs       []string `json:"groupids,omitempty"`
	TemplateIDs    []string `json:"templateids,omitempty"`
	InterfaceIDs   []string `json:"interfaceids,omitempty"`
	TriggerIDs     []string `json:"triggerids,omitempty"`
	Tags           []Tag    `json:"tags,omitempty"`
	EvalType       int      `json:"evaltype,omitempty"`
	Templated      *bool    `json:"templated,omitempty"`
	Monitored      *bool    `json:"monitored,omitempty"`
	WebItems       bool     `json:"webitems,omitempty"`
	SelectHosts    any      `json:"selectHosts,omitempty"`
	SelectTriggers any      `json:"selectTriggers,omitempty"`
	SelectTags     any      `json:"selectTags,omitempty"`
}

// Get retrieves items matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/item/get
func (s *ItemService) Get(ctx context.Context, params ItemGetParams) ([]Item, error) {
	var items []Item
	err := s.client.Call(ctx, "item.get", params, &items)
	return items, err
}

// Create creates items and returns the new item ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/item/create
func (s *ItemService) Create(ctx context.Context, items ...Item) ([]string, error) {
	var res struct {
		ItemIDs []string `json:"itemids"`
	}
	err := s.client.Call(ctx, "item.create", items, &res)
	return res.ItemIDs, err
}

// Update updates items and returns the affected item ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/item/update
func (s *ItemService) Update(ctx context.Context, items ...Item) ([]string, error) {
	var res struct {
		ItemIDs []string `json:"itemids"`
	}
	err := s.client.Call(ctx, "item.update", items, &res)
	return res.ItemIDs, err
}

// Delete deletes items by id and returns the deleted item ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/item/delete
func (s *ItemService) Delete(ctx context.Context, itemIDs ...string) ([]string, error) {
	var res struct {
		ItemIDs []string `json:"itemids"`
	}
	err := s.client.Call(ctx, "item.delete", itemIDs, &res)
	return res.ItemIDs, err
}
