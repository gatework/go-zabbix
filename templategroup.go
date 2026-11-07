package zabbix

import "context"

// TemplateGroupService wraps the "templategroup" API namespace.
type TemplateGroupService struct{ client *Client }

// TemplateGroup returns the templategroup service.
func (c *Client) TemplateGroup() *TemplateGroupService { return &TemplateGroupService{c} }

// TemplateGroupGetParams are the parameters for templategroup.get.
type TemplateGroupGetParams struct {
	GetParams
	GroupIDs        []string `json:"groupids,omitempty"`
	TemplateIDs     []string `json:"templateids,omitempty"`
	GraphIDs        []string `json:"graphids,omitempty"`
	TriggerIDs      []string `json:"triggerids,omitempty"`
	WithItems       bool     `json:"with_items,omitempty"`
	WithTriggers    bool     `json:"with_triggers,omitempty"`
	SelectTemplates any      `json:"selectTemplates,omitempty"`
}

// Get retrieves template groups matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/templategroup/get
func (s *TemplateGroupService) Get(ctx context.Context, params TemplateGroupGetParams) ([]TemplateGroup, error) {
	var out []TemplateGroup
	err := s.client.Call(ctx, "templategroup.get", params, &out)
	return out, err
}

// Create creates template groups and returns the new group ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/templategroup/create
func (s *TemplateGroupService) Create(ctx context.Context, groups ...TemplateGroup) ([]string, error) {
	var res struct {
		IDs []string `json:"groupids"`
	}
	err := s.client.Call(ctx, "templategroup.create", groups, &res)
	return res.IDs, err
}

// Update updates template groups and returns the affected group ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/templategroup/update
func (s *TemplateGroupService) Update(ctx context.Context, groups ...TemplateGroup) ([]string, error) {
	var res struct {
		IDs []string `json:"groupids"`
	}
	err := s.client.Call(ctx, "templategroup.update", groups, &res)
	return res.IDs, err
}

// Delete deletes template groups by id and returns the deleted group ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/templategroup/delete
func (s *TemplateGroupService) Delete(ctx context.Context, groupIDs ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"groupids"`
	}
	err := s.client.Call(ctx, "templategroup.delete", groupIDs, &res)
	return res.IDs, err
}
