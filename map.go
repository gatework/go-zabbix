package zabbix

import "context"

// MapService wraps the "map" API namespace.
type MapService struct{ client *Client }

// Map returns the map service.
func (c *Client) Map() *MapService { return &MapService{c} }

// Map object (network map / sysmap).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/map/object
type Map struct {
	SysmapID           string `json:"sysmapid,omitempty"`
	Name               string `json:"name,omitempty"`
	Width              string `json:"width,omitempty"`
	Height             string `json:"height,omitempty"`
	BackgroundID       string `json:"backgroundid,omitempty"`
	LabelType          string `json:"label_type,omitempty"`
	LabelLocation      string `json:"label_location,omitempty"`
	Highlight          string `json:"highlight,omitempty"`
	ExpandProblem      string `json:"expandproblem,omitempty"`
	MarkElements       string `json:"markelements,omitempty"`
	ShowUnack          string `json:"show_unack,omitempty"`
	GridSize           string `json:"grid_size,omitempty"`
	GridShow           string `json:"grid_show,omitempty"`
	GridAlign          string `json:"grid_align,omitempty"`
	LabelFormat        string `json:"label_format,omitempty"`
	LabelTypeHost      string `json:"label_type_host,omitempty"`
	LabelTypeHostGroup string `json:"label_type_hostgroup,omitempty"`
	LabelTypeTrigger   string `json:"label_type_trigger,omitempty"`
	LabelTypeMap       string `json:"label_type_map,omitempty"`
	LabelTypeImage     string `json:"label_type_image,omitempty"`
	LabelStringHost    string `json:"label_string_host,omitempty"`
	IconMapID          string `json:"iconmapid,omitempty"`
	ExpandMacros       string `json:"expand_macros,omitempty"`
	SeverityMin        string `json:"severity_min,omitempty"`
	UserID             string `json:"userid,omitempty"`
	Private            string `json:"private,omitempty"`
	ShowSuppressed     string `json:"show_suppressed,omitempty"`

	Selements  []MapSelement  `json:"selements,omitempty"`
	Links      []MapLink      `json:"links,omitempty"`
	Urls       []MapURL       `json:"urls,omitempty"`
	Users      []MapUser      `json:"users,omitempty"`
	UserGroups []MapUserGroup `json:"userGroups,omitempty"`
	Shapes     any            `json:"shapes,omitempty"`
	Lines      any            `json:"lines,omitempty"`
}

// MapSelement is a map element (icon) placed on a map.
type MapSelement struct {
	SelementID     string `json:"selementid,omitempty"`
	ElementType    string `json:"elementtype,omitempty"`
	IconIDOff      string `json:"iconid_off,omitempty"`
	IconIDOn       string `json:"iconid_on,omitempty"`
	Label          string `json:"label,omitempty"`
	LabelLocation  string `json:"label_location,omitempty"`
	X              string `json:"x,omitempty"`
	Y              string `json:"y,omitempty"`
	ElementSubtype string `json:"elementsubtype,omitempty"`
	AreaType       string `json:"areatype,omitempty"`
	Width          string `json:"width,omitempty"`
	Height         string `json:"height,omitempty"`
	ViewType       string `json:"viewtype,omitempty"`
	UseIconMap     string `json:"use_iconmap,omitempty"`
	EvalType       string `json:"evaltype,omitempty"`
}

// MapLink is a connection between two map elements.
type MapLink struct {
	LinkID      string `json:"linkid,omitempty"`
	SelementID1 string `json:"selementid1,omitempty"`
	SelementID2 string `json:"selementid2,omitempty"`
	DrawType    string `json:"drawtype,omitempty"`
	Color       string `json:"color,omitempty"`
	Label       string `json:"label,omitempty"`
}

// MapURL is a URL associated with a map.
type MapURL struct {
	SysmapURLID string `json:"sysmapurlid,omitempty"`
	Name        string `json:"name,omitempty"`
	URL         string `json:"url,omitempty"`
	ElementType string `json:"elementtype,omitempty"`
}

// MapUser is a user share of a map.
type MapUser struct {
	SysmapUserID string `json:"sysmapuserid,omitempty"`
	UserID       string `json:"userid,omitempty"`
	Permission   string `json:"permission,omitempty"`
}

// MapUserGroup is a user-group share of a map.
type MapUserGroup struct {
	SysmapUsrGrpID string `json:"sysmapusrgrpid,omitempty"`
	UsrGrpID       string `json:"usrgrpid,omitempty"`
	Permission     string `json:"permission,omitempty"`
}

// MapGetParams are the parameters for map.get.
type MapGetParams struct {
	GetParams
	SysmapIDs        []string `json:"sysmapids,omitempty"`
	UserIDs          []string `json:"userids,omitempty"`
	Expandproblem    *bool    `json:"expandproblem,omitempty"`
	SelectSelements  any      `json:"selectSelements,omitempty"`
	SelectLinks      any      `json:"selectLinks,omitempty"`
	SelectUrls       any      `json:"selectUrls,omitempty"`
	SelectUsers      any      `json:"selectUsers,omitempty"`
	SelectUserGroups any      `json:"selectUserGroups,omitempty"`
	SelectShapes     any      `json:"selectShapes,omitempty"`
	SelectLines      any      `json:"selectLines,omitempty"`
}

// Get retrieves maps matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/map/get
func (s *MapService) Get(ctx context.Context, params MapGetParams) ([]Map, error) {
	var out []Map
	err := s.client.Call(ctx, "map.get", params, &out)
	return out, err
}

// Create creates maps and returns the new map ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/map/create
func (s *MapService) Create(ctx context.Context, maps ...Map) ([]string, error) {
	var res struct {
		IDs []string `json:"sysmapids"`
	}
	err := s.client.Call(ctx, "map.create", maps, &res)
	return res.IDs, err
}

// Update updates maps and returns the affected map ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/map/update
func (s *MapService) Update(ctx context.Context, maps ...Map) ([]string, error) {
	var res struct {
		IDs []string `json:"sysmapids"`
	}
	err := s.client.Call(ctx, "map.update", maps, &res)
	return res.IDs, err
}

// Delete deletes maps by id and returns the deleted map ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/map/delete
func (s *MapService) Delete(ctx context.Context, sysmapIDs ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"sysmapids"`
	}
	err := s.client.Call(ctx, "map.delete", sysmapIDs, &res)
	return res.IDs, err
}
