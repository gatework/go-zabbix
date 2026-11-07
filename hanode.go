package zabbix

import "context"

// HANodeService wraps the "hanode" API namespace (high availability nodes).
type HANodeService struct{ client *Client }

// HANode returns the hanode service.
func (c *Client) HANode() *HANodeService { return &HANodeService{c} }

// HANode is a Zabbix server high availability node.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hanode/object
type HANode struct {
	HANodeID   string `json:"ha_nodeid,omitempty"`
	Name       string `json:"name,omitempty"`
	Address    string `json:"address,omitempty"`
	Port       string `json:"port,omitempty"`
	LastAccess string `json:"lastaccess,omitempty"`
	Status     string `json:"status,omitempty"`
}

// HANodeGetParams are the parameters for hanode.get.
type HANodeGetParams struct {
	GetParams
	HANodeIDs []string `json:"ha_nodeids,omitempty"`
	Status    []int    `json:"status,omitempty"`
}

// Get retrieves HA nodes matching params. The hanode API is read-only.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hanode/get
func (s *HANodeService) Get(ctx context.Context, params HANodeGetParams) ([]HANode, error) {
	var out []HANode
	err := s.client.Call(ctx, "hanode.get", params, &out)
	return out, err
}
