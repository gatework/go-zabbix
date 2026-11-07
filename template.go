package zabbix

import "context"

// TemplateService wraps the "template" API namespace.
type TemplateService struct{ client *Client }

// Template returns the template service.
func (c *Client) Template() *TemplateService { return &TemplateService{c} }

// Template object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/template/object
type Template struct {
	TemplateID  string `json:"templateid,omitempty"`
	Host        string `json:"host,omitempty"` // technical name
	Name        string `json:"name,omitempty"` // visible name
	Description string `json:"description,omitempty"`
	UUID        string `json:"uuid,omitempty"`
	Vendor      string `json:"vendor_name,omitempty"`
	Version     string `json:"vendor_version,omitempty"`

	Groups []TemplateGroup `json:"groups,omitempty"`
	Tags   []Tag           `json:"tags,omitempty"`
	Macros []UserMacro     `json:"macros,omitempty"`
}

// TemplateGroup is a template-group reference (Zabbix 6.2+ split templates into
// their own group type).
type TemplateGroup struct {
	GroupID string `json:"groupid,omitempty"`
	Name    string `json:"name,omitempty"`
}

// TemplateGetParams are the parameters for template.get.
type TemplateGetParams struct {
	GetParams
	TemplateIDs       []string `json:"templateids,omitempty"`
	GroupIDs          []string `json:"groupids,omitempty"`
	ParentTemplateIDs []string `json:"parentTemplateids,omitempty"`
	HostIDs           []string `json:"hostids,omitempty"`
	ItemIDs           []string `json:"itemids,omitempty"`
	TriggerIDs        []string `json:"triggerids,omitempty"`
	Tags              []Tag    `json:"tags,omitempty"`
	EvalType          int      `json:"evaltype,omitempty"`
	SelectGroups      any      `json:"selectTemplateGroups,omitempty"`
	SelectHosts       any      `json:"selectHosts,omitempty"`
	SelectItems       any      `json:"selectItems,omitempty"`
	SelectTriggers    any      `json:"selectTriggers,omitempty"`
	SelectTags        any      `json:"selectTags,omitempty"`
	SelectMacros      any      `json:"selectMacros,omitempty"`
}

// Get retrieves templates matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/template/get
func (s *TemplateService) Get(ctx context.Context, params TemplateGetParams) ([]Template, error) {
	var templates []Template
	err := s.client.Call(ctx, "template.get", params, &templates)
	return templates, err
}

// Create creates templates and returns the new template ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/template/create
func (s *TemplateService) Create(ctx context.Context, templates ...Template) ([]string, error) {
	var res struct {
		TemplateIDs []string `json:"templateids"`
	}
	err := s.client.Call(ctx, "template.create", templates, &res)
	return res.TemplateIDs, err
}

// Update updates templates and returns the affected template ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/template/update
func (s *TemplateService) Update(ctx context.Context, templates ...Template) ([]string, error) {
	var res struct {
		TemplateIDs []string `json:"templateids"`
	}
	err := s.client.Call(ctx, "template.update", templates, &res)
	return res.TemplateIDs, err
}

// Delete deletes templates by id and returns the deleted template ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/template/delete
func (s *TemplateService) Delete(ctx context.Context, templateIDs ...string) ([]string, error) {
	var res struct {
		TemplateIDs []string `json:"templateids"`
	}
	err := s.client.Call(ctx, "template.delete", templateIDs, &res)
	return res.TemplateIDs, err
}
