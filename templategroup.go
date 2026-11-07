package zabbix

import "context"

// TemplateGroupService wraps the "templategroup" API namespace.
type TemplateGroupService struct{ client *Client }

// TemplateGroup returns the templategroup service.
func (c *Client) TemplateGroup() *TemplateGroupService { return &TemplateGroupService{c} }

// TemplateGroupGetParams are the parameters for templategroup.get.
// TemplateGroupGetParams queries the list form of templategroup.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type TemplateGroupGetParams struct {
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
	GroupIDs               []string       `json:"groupids,omitzero"`
	TemplateIDs            []string       `json:"templateids,omitzero"`
	GraphIDs               []string       `json:"graphids,omitzero"`
	TriggerIDs             []string       `json:"triggerids,omitzero"`
	WithItems              bool           `json:"with_items,omitempty"`
	WithTriggers           bool           `json:"with_triggers,omitempty"`
	SelectTemplates        any            `json:"selectTemplates,omitempty"`
}

// Get retrieves template groups matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templategroup/get
func (s *TemplateGroupService) Get(ctx context.Context, params TemplateGroupGetParams) ([]TemplateGroup, error) {
	return getList[TemplateGroup](ctx, s.client, "templategroup.get", params)
}

// TemplateGroupCountParams contains the filters for templategroup.get counts.
type TemplateGroupCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	TemplateIDs            []string       `json:"templateids,omitzero"`
	GraphIDs               []string       `json:"graphids,omitzero"`
	TriggerIDs             []string       `json:"triggerids,omitzero"`
	WithItems              bool           `json:"with_items,omitempty"`
	WithTriggers           bool           `json:"with_triggers,omitempty"`
}

// Count returns the number of objects matching params.
func (s *TemplateGroupService) Count(ctx context.Context, params TemplateGroupCountParams) (uint64, error) {
	return getCount(ctx, s.client, "templategroup.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *TemplateGroupService) GetMap(ctx context.Context, params TemplateGroupGetParams) (map[string]TemplateGroup, error) {
	return getMap[TemplateGroup](ctx, s.client, "templategroup.get", params)
}

// Create creates template groups and returns the new group ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templategroup/create
func (s *TemplateGroupService) Create(ctx context.Context, groups ...TemplateGroupCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "templategroup.create", "groupids", groups)
}

// Update updates template groups and returns the affected group ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templategroup/update
func (s *TemplateGroupService) Update(ctx context.Context, groups ...TemplateGroupUpdateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "templategroup.update", "groupids", groups)
}

// Delete deletes template groups by id and returns the deleted group ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templategroup/delete
func (s *TemplateGroupService) Delete(ctx context.Context, groupIDs ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "templategroup.delete", "groupids", groupIDs)
}
