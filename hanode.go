package zabbix

import "context"

// HANodeService wraps the "hanode" API namespace (high availability nodes).
type HANodeService struct{ client *Client }

// HANode returns the hanode service.
func (c *Client) HANode() *HANodeService { return &HANodeService{c} }

// HANode is a Zabbix server high availability node.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hanode/object
type HANode struct {
	HANodeID   string `json:"ha_nodeid,omitempty"`
	Name       string `json:"name,omitempty"`
	Address    string `json:"address,omitempty"`
	Port       string `json:"port,omitempty"`
	LastAccess string `json:"lastaccess,omitempty"`
	Status     string `json:"status,omitempty"`
}

// HANodeGetParams are the parameters for hanode.get.
// HANodeGetParams queries the list form of hanode.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type HANodeGetParams struct {
	Output    any            `json:"output,omitempty"`
	Filter    map[string]any `json:"filter,omitzero"`
	SortField []string       `json:"sortfield,omitzero"`
	SortOrder []string       `json:"sortorder,omitzero"`
	Limit     int            `json:"limit,omitempty"`
	HANodeIDs []string       `json:"ha_nodeids,omitzero"`
}

// Get retrieves HA nodes matching params. The hanode API is read-only.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hanode/get
func (s *HANodeService) Get(ctx context.Context, params HANodeGetParams) ([]HANode, error) {
	return getList[HANode](ctx, s.client, "hanode.get", params)
}

// HANodeCountParams contains the filters for hanode.get counts.
type HANodeCountParams struct {
	Filter    map[string]any `json:"filter,omitzero"`
	HANodeIDs []string       `json:"ha_nodeids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *HANodeService) Count(ctx context.Context, params HANodeCountParams) (uint64, error) {
	return getCount(ctx, s.client, "hanode.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *HANodeService) GetMap(ctx context.Context, params HANodeGetParams) (map[string]HANode, error) {
	return getMap[HANode](ctx, s.client, "hanode.get", params)
}
