package zabbix

import "context"

// GraphService wraps the "graph" API namespace.
type GraphService struct{ client *Client }

// Graph returns the graph service.
func (c *Client) Graph() *GraphService { return &GraphService{c} }

// Graph object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graph/object
type Graph struct {
	GraphID        string      `json:"graphid,omitempty"`
	Name           string      `json:"name,omitempty"`
	Width          string      `json:"width,omitempty"`
	Height         string      `json:"height,omitempty"`
	YAxisMin       string      `json:"yaxismin,omitempty"`
	YAxisMax       string      `json:"yaxismax,omitempty"`
	TemplateID     string      `json:"templateid,omitempty"`
	ShowWorkPeriod string      `json:"show_work_period,omitempty"`
	ShowTriggers   string      `json:"show_triggers,omitempty"`
	GraphType      string      `json:"graphtype,omitempty"`
	ShowLegend     string      `json:"show_legend,omitempty"`
	Show3D         string      `json:"show_3d,omitempty"`
	PercentLeft    string      `json:"percent_left,omitempty"`
	PercentRight   string      `json:"percent_right,omitempty"`
	YMinType       string      `json:"ymin_type,omitempty"`
	YMaxType       string      `json:"ymax_type,omitempty"`
	YMinItemID     string      `json:"ymin_itemid,omitempty"`
	YMaxItemID     string      `json:"ymax_itemid,omitempty"`
	Flags          string      `json:"flags,omitempty"`
	UUID           string      `json:"uuid,omitempty"`
	GItems         []GraphItem `json:"gitems,omitempty"`
	Hosts          []Host      `json:"hosts,omitempty"`
	Items          []Item      `json:"items,omitempty"`
	Templates      []Template  `json:"templates,omitempty"`
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
// GraphGetParams queries the list form of graph.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type GraphGetParams struct {
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
	GraphIDs               []string       `json:"graphids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	TemplateIDs            []string       `json:"templateids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	ItemIDs                []string       `json:"itemids,omitzero"`
	Templated              *bool          `json:"templated,omitempty"`
	SelectGraphItems       any            `json:"selectGraphItems,omitempty"`
	SelectHosts            any            `json:"selectHosts,omitempty"`
	SelectItems            any            `json:"selectItems,omitempty"`
	SelectTemplates        any            `json:"selectTemplates,omitempty"`
}

// Get retrieves graphs matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graph/get
func (s *GraphService) Get(ctx context.Context, params GraphGetParams) ([]Graph, error) {
	return getList[Graph](ctx, s.client, "graph.get", params)
}

// GraphCountParams contains the filters for graph.get counts.
type GraphCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	GraphIDs               []string       `json:"graphids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	TemplateIDs            []string       `json:"templateids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	ItemIDs                []string       `json:"itemids,omitzero"`
	Templated              *bool          `json:"templated,omitempty"`
}

// Count returns the number of objects matching params.
func (s *GraphService) Count(ctx context.Context, params GraphCountParams) (uint64, error) {
	return getCount(ctx, s.client, "graph.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *GraphService) GetMap(ctx context.Context, params GraphGetParams) (map[string]Graph, error) {
	return getMap[Graph](ctx, s.client, "graph.get", params)
}

// Create creates graphs and returns the new graph ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graph/create
func (s *GraphService) Create(ctx context.Context, graphs ...GraphCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "graph.create", "graphids", graphs)
}

// Update updates graphs and returns the affected graph ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graph/update
func (s *GraphService) Update(ctx context.Context, graphs ...GraphUpdateParams) ([]string, error) {
	return updateIDs(ctx, s.client, "graph.update", "graphids", graphs)
}

// Delete deletes graphs by id and returns the deleted graph ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graph/delete
func (s *GraphService) Delete(ctx context.Context, graphIDs ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "graph.delete", "graphids", graphIDs)
}
