package zabbix

import "context"

// GraphItemService wraps the "graphitem" API namespace.
type GraphItemService struct{ client *Client }

// GraphItem returns the graphitem service.
func (c *Client) GraphItem() *GraphItemService { return &GraphItemService{c} }

// GraphItem object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/graphitem/object
type GraphItem struct {
	GItemID   string `json:"gitemid,omitempty"`
	GraphID   string `json:"graphid,omitempty"`
	ItemID    string `json:"itemid,omitempty"`
	DrawType  string `json:"drawtype,omitempty"`
	SortOrder string `json:"sortorder,omitempty"`
	Color     string `json:"color,omitempty"`
	YAxisSide string `json:"yaxisside,omitempty"`
	CalcFnc   string `json:"calc_fnc,omitempty"`
	Type      string `json:"type,omitempty"`
}

// GraphItemGetParams are the parameters for graphitem.get.
type GraphItemGetParams struct {
	GetParams
	GraphIDs     []string `json:"graphids,omitempty"`
	ItemIDs      []string `json:"itemids,omitempty"`
	Type         *int     `json:"type,omitempty"`
	SelectGraphs any      `json:"selectGraphs,omitempty"`
}

// Get retrieves graph items matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/graphitem/get
func (s *GraphItemService) Get(ctx context.Context, params GraphItemGetParams) ([]GraphItem, error) {
	var out []GraphItem
	err := s.client.Call(ctx, "graphitem.get", params, &out)
	return out, err
}
