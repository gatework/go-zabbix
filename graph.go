package zabbix

import "context"

// GraphService wraps the "graph" API namespace.
type GraphService struct{ client *Client }

// Graph returns the graph service.
func (c *Client) Graph() *GraphService { return &GraphService{c} }

// Graph object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/graph/object
type Graph struct {
	GraphID        string `json:"graphid,omitempty"`
	Name           string `json:"name,omitempty"`
	Width          string `json:"width,omitempty"`
	Height         string `json:"height,omitempty"`
	YAxisMin       string `json:"yaxismin,omitempty"`
	YAxisMax       string `json:"yaxismax,omitempty"`
	TemplateID     string `json:"templateid,omitempty"`
	ShowWorkPeriod string `json:"show_work_period,omitempty"`
	ShowTriggers   string `json:"show_triggers,omitempty"`
	GraphType      string `json:"graphtype,omitempty"`
	ShowLegend     string `json:"show_legend,omitempty"`
	Show3D         string `json:"show_3d,omitempty"`
	PercentLeft    string `json:"percent_left,omitempty"`
	PercentRight   string `json:"percent_right,omitempty"`
	YMinType       string `json:"ymin_type,omitempty"`
	YMaxType       string `json:"ymax_type,omitempty"`
	YMinItemID     string `json:"ymin_itemid,omitempty"`
	YMaxItemID     string `json:"ymax_itemid,omitempty"`
	Flags          string `json:"flags,omitempty"`
	UUID           string `json:"uuid,omitempty"`

	GItems []GraphGItem `json:"gitems,omitempty"`
}

// GraphGItem is a graph item (a metric line) belonging to a graph.
type GraphGItem struct {
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

// GraphGetParams are the parameters for graph.get.
type GraphGetParams struct {
	GetParams
	GraphIDs         []string `json:"graphids,omitempty"`
	GroupIDs         []string `json:"groupids,omitempty"`
	TemplateIDs      []string `json:"templateids,omitempty"`
	HostIDs          []string `json:"hostids,omitempty"`
	ItemIDs          []string `json:"itemids,omitempty"`
	Templated        *bool    `json:"templated,omitempty"`
	SelectGraphItems any      `json:"selectGraphItems,omitempty"`
	SelectHosts      any      `json:"selectHosts,omitempty"`
	SelectItems      any      `json:"selectItems,omitempty"`
	SelectTemplates  any      `json:"selectTemplates,omitempty"`
}

// Get retrieves graphs matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/graph/get
func (s *GraphService) Get(ctx context.Context, params GraphGetParams) ([]Graph, error) {
	var out []Graph
	err := s.client.Call(ctx, "graph.get", params, &out)
	return out, err
}

// Create creates graphs and returns the new graph ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/graph/create
func (s *GraphService) Create(ctx context.Context, graphs ...Graph) ([]string, error) {
	var res struct {
		IDs []string `json:"graphids"`
	}
	err := s.client.Call(ctx, "graph.create", graphs, &res)
	return res.IDs, err
}

// Update updates graphs and returns the affected graph ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/graph/update
func (s *GraphService) Update(ctx context.Context, graphs ...Graph) ([]string, error) {
	var res struct {
		IDs []string `json:"graphids"`
	}
	err := s.client.Call(ctx, "graph.update", graphs, &res)
	return res.IDs, err
}

// Delete deletes graphs by id and returns the deleted graph ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/graph/delete
func (s *GraphService) Delete(ctx context.Context, graphIDs ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"graphids"`
	}
	err := s.client.Call(ctx, "graph.delete", graphIDs, &res)
	return res.IDs, err
}
