package zabbix

import "context"

// GraphItemService wraps the "graphitem" API namespace.
type GraphItemService struct{ client *Client }

// GraphItem returns the graphitem service.
func (c *Client) GraphItem() *GraphItemService { return &GraphItemService{c} }

// GraphItem object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graphitem/object
type GraphItem struct {
	GItemID   string  `json:"gitemid,omitempty"`
	GraphID   string  `json:"graphid,omitempty"`
	ItemID    string  `json:"itemid,omitempty"`
	DrawType  string  `json:"drawtype,omitempty"`
	SortOrder string  `json:"sortorder,omitempty"`
	Color     string  `json:"color,omitempty"`
	YAxisSide string  `json:"yaxisside,omitempty"`
	CalcFnc   string  `json:"calc_fnc,omitempty"`
	Type      string  `json:"type,omitempty"`
	Graphs    []Graph `json:"graphs,omitempty"`
}

// GraphItemGetParams are the parameters for graphitem.get.
// GraphItemGetParams queries the list form of graphitem.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type GraphItemGetParams struct {
	Output       any      `json:"output,omitempty"`
	SortField    []string `json:"sortfield,omitzero"`
	SortOrder    []string `json:"sortorder,omitzero"`
	Limit        int      `json:"limit,omitempty"`
	Editable     bool     `json:"editable,omitempty"`
	GraphIDs     []string `json:"graphids,omitzero"`
	ItemIDs      []string `json:"itemids,omitzero"`
	Type         *int     `json:"type,omitempty"`
	SelectGraphs any      `json:"selectGraphs,omitempty"`
}

// Get retrieves graph items matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graphitem/get
func (s *GraphItemService) Get(ctx context.Context, params GraphItemGetParams) ([]GraphItem, error) {
	return getList[GraphItem](ctx, s.client, "graphitem.get", params)
}

// GraphItemCountParams contains the filters for graphitem.get counts.
type GraphItemCountParams struct {
	Editable bool     `json:"editable,omitempty"`
	GraphIDs []string `json:"graphids,omitzero"`
	ItemIDs  []string `json:"itemids,omitzero"`
	Type     *int     `json:"type,omitempty"`
}

// Count returns the number of objects matching params.
func (s *GraphItemService) Count(ctx context.Context, params GraphItemCountParams) (uint64, error) {
	return getCount(ctx, s.client, "graphitem.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *GraphItemService) GetMap(ctx context.Context, params GraphItemGetParams) (map[string]GraphItem, error) {
	return getMap[GraphItem](ctx, s.client, "graphitem.get", params)
}
