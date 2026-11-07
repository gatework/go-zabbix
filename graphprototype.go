package zabbix

import "context"

// GraphPrototypeService wraps the "graphprototype" API namespace.
type GraphPrototypeService struct{ client *Client }

// GraphPrototype returns the graphprototype service.
func (c *Client) GraphPrototype() *GraphPrototypeService { return &GraphPrototypeService{c} }

// GraphPrototype object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graphprototype/object
type GraphPrototype struct {
	GraphID        string                        `json:"graphid,omitempty"`
	Name           string                        `json:"name,omitempty"`
	Width          string                        `json:"width,omitempty"`
	Height         string                        `json:"height,omitempty"`
	YAxisMin       string                        `json:"yaxismin,omitempty"`
	YAxisMax       string                        `json:"yaxismax,omitempty"`
	TemplateID     string                        `json:"templateid,omitempty"`
	ShowWorkPeriod string                        `json:"show_work_period,omitempty"`
	ShowTriggers   string                        `json:"show_triggers,omitempty"`
	GraphType      string                        `json:"graphtype,omitempty"`
	ShowLegend     string                        `json:"show_legend,omitempty"`
	Show3D         string                        `json:"show_3d,omitempty"`
	PercentLeft    string                        `json:"percent_left,omitempty"`
	PercentRight   string                        `json:"percent_right,omitempty"`
	YMinType       string                        `json:"ymin_type,omitempty"`
	YMaxType       string                        `json:"ymax_type,omitempty"`
	YMinItemID     string                        `json:"ymin_itemid,omitempty"`
	YMaxItemID     string                        `json:"ymax_itemid,omitempty"`
	Discover       string                        `json:"discover,omitempty"`
	UUID           string                        `json:"uuid,omitempty"`
	Flags          string                        `json:"flags,omitempty"`
	GItems         []GraphItem                   `json:"gitems,omitempty"`
	Hosts          []Host                        `json:"hosts,omitempty"`
	Items          []ItemPrototype               `json:"items,omitempty"`
	DiscoveryRule  SelectedObject[DiscoveryRule] `json:"discoveryRule,omitzero"`
}

// GraphPrototypeGItem is a graph item (line) of a graph prototype.
type GraphPrototypeGItem struct {
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

// GraphPrototypeGetParams are the parameters for graphprototype.get.
// GraphPrototypeGetParams queries the list form of graphprototype.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type GraphPrototypeGetParams struct {
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
	DiscoveryIDs           []string       `json:"discoveryids,omitzero"`
	SelectGraphItems       any            `json:"selectGraphItems,omitempty"`
	SelectHosts            any            `json:"selectHosts,omitempty"`
	SelectItems            any            `json:"selectItems,omitempty"`
	SelectDiscoveryRule    any            `json:"selectDiscoveryRule,omitempty"`
}

// Get retrieves graph prototypes matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graphprototype/get
func (s *GraphPrototypeService) Get(ctx context.Context, params GraphPrototypeGetParams) ([]GraphPrototype, error) {
	return getList[GraphPrototype](ctx, s.client, "graphprototype.get", params)
}

// GraphPrototypeCountParams contains the filters for graphprototype.get counts.
type GraphPrototypeCountParams struct {
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
	DiscoveryIDs           []string       `json:"discoveryids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *GraphPrototypeService) Count(ctx context.Context, params GraphPrototypeCountParams) (uint64, error) {
	return getCount(ctx, s.client, "graphprototype.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *GraphPrototypeService) GetMap(ctx context.Context, params GraphPrototypeGetParams) (map[string]GraphPrototype, error) {
	return getMap[GraphPrototype](ctx, s.client, "graphprototype.get", params)
}

// Create creates graph prototypes and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graphprototype/create
func (s *GraphPrototypeService) Create(ctx context.Context, graphs ...GraphPrototype) ([]string, error) {
	var res struct {
		IDs []string `json:"graphids"`
	}
	err := s.client.Call(ctx, "graphprototype.create", graphs, &res)
	return res.IDs, err
}

// Update updates graph prototypes and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graphprototype/update
func (s *GraphPrototypeService) Update(ctx context.Context, graphs ...GraphPrototype) ([]string, error) {
	var res struct {
		IDs []string `json:"graphids"`
	}
	err := s.client.Call(ctx, "graphprototype.update", graphs, &res)
	return res.IDs, err
}

// Delete deletes graph prototypes by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graphprototype/delete
func (s *GraphPrototypeService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"graphids"`
	}
	err := s.client.Call(ctx, "graphprototype.delete", ids, &res)
	return res.IDs, err
}
