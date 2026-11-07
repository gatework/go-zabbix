package zabbix

import "context"

// TemplateService wraps the "template" API namespace.
type TemplateService struct{ client *Client }

// Template returns the template service.
func (c *Client) Template() *TemplateService { return &TemplateService{c} }

// Template object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/template/object
type Template struct {
	TemplateID  string `json:"templateid,omitempty"`
	Host        string `json:"host,omitempty"` // technical name
	Name        string `json:"name,omitempty"` // visible name
	Description string `json:"description,omitempty"`
	UUID        string `json:"uuid,omitempty"`
	Vendor      string `json:"vendor_name,omitempty"`
	Version     string `json:"vendor_version,omitempty"`

	Tags           []Tag             `json:"tags,omitempty"`
	Macros         []UserMacro       `json:"macros,omitempty"`
	TemplateGroups []TemplateGroup   `json:"templategroups,omitempty"`
	Hosts          Relation[Host]    `json:"hosts,omitzero"`
	Items          Relation[Item]    `json:"items,omitzero"`
	Triggers       Relation[Trigger] `json:"triggers,omitzero"`
	LinkType       string            `json:"link_type,omitempty"`
}

// TemplateGroup is a template-group reference (Zabbix 6.2+ split templates into
// their own group type).
type TemplateGroup struct {
	GroupID   string             `json:"groupid,omitempty"`
	Name      string             `json:"name,omitempty"`
	Templates Relation[Template] `json:"templates,omitzero"`
}

// TemplateGetParams are the parameters for template.get.
// TemplateGetParams queries the list form of template.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type TemplateGetParams struct {
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
	TemplateIDs            []string       `json:"templateids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	ParentTemplateIDs      []string       `json:"parentTemplateids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	ItemIDs                []string       `json:"itemids,omitzero"`
	TriggerIDs             []string       `json:"triggerids,omitzero"`
	Tags                   []TagFilter    `json:"tags,omitzero"`
	EvalType               int            `json:"evaltype,omitempty"`
	SelectTemplateGroups   any            `json:"selectTemplateGroups,omitempty"`
	SelectHosts            any            `json:"selectHosts,omitempty"`
	SelectItems            any            `json:"selectItems,omitempty"`
	SelectTriggers         any            `json:"selectTriggers,omitempty"`
	SelectTags             any            `json:"selectTags,omitempty"`
	SelectMacros           any            `json:"selectMacros,omitempty"`
}

// Get retrieves templates matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/template/get
func (s *TemplateService) Get(ctx context.Context, params TemplateGetParams) ([]Template, error) {
	return getList[Template](ctx, s.client, "template.get", params)
}

// TemplateCountParams contains the filters for template.get counts.
type TemplateCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	TemplateIDs            []string       `json:"templateids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	ParentTemplateIDs      []string       `json:"parentTemplateids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	ItemIDs                []string       `json:"itemids,omitzero"`
	TriggerIDs             []string       `json:"triggerids,omitzero"`
	Tags                   []TagFilter    `json:"tags,omitzero"`
	EvalType               int            `json:"evaltype,omitempty"`
}

// Count returns the number of objects matching params.
func (s *TemplateService) Count(ctx context.Context, params TemplateCountParams) (uint64, error) {
	return getCount(ctx, s.client, "template.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *TemplateService) GetMap(ctx context.Context, params TemplateGetParams) (map[string]Template, error) {
	return getMap[Template](ctx, s.client, "template.get", params)
}

// Create creates templates and returns the new template ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/template/create
func (s *TemplateService) Create(ctx context.Context, templates ...Template) ([]string, error) {
	var res struct {
		TemplateIDs []string `json:"templateids"`
	}
	err := s.client.Call(ctx, "template.create", templates, &res)
	return res.TemplateIDs, err
}

// Update updates templates and returns the affected template ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/template/update
func (s *TemplateService) Update(ctx context.Context, templates ...Template) ([]string, error) {
	var res struct {
		TemplateIDs []string `json:"templateids"`
	}
	err := s.client.Call(ctx, "template.update", templates, &res)
	return res.TemplateIDs, err
}

// Delete deletes templates by id and returns the deleted template ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/template/delete
func (s *TemplateService) Delete(ctx context.Context, templateIDs ...string) ([]string, error) {
	var res struct {
		TemplateIDs []string `json:"templateids"`
	}
	err := s.client.Call(ctx, "template.delete", templateIDs, &res)
	return res.TemplateIDs, err
}
