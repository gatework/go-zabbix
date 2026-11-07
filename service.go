package zabbix

import "context"

// ServiceService wraps the "service" API namespace.
type ServiceService struct{ client *Client }

// Service returns the service service.
func (c *Client) Service() *ServiceService { return &ServiceService{c} }

// Service object (Zabbix 7.x).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/service/object
type Service struct {
	ServiceID        string                      `json:"serviceid,omitempty"`
	Name             string                      `json:"name,omitempty"`
	Status           string                      `json:"status,omitempty"`
	Algorithm        string                      `json:"algorithm,omitempty"`
	SortOrder        string                      `json:"sortorder,omitempty"`
	Weight           string                      `json:"weight,omitempty"`
	PropagationRule  string                      `json:"propagation_rule,omitempty"`
	PropagationValue string                      `json:"propagation_value,omitempty"`
	Description      string                      `json:"description,omitempty"`
	UUID             string                      `json:"uuid,omitempty"`
	CreatedAt        string                      `json:"created_at,omitempty"`
	ReadOnly         bool                        `json:"readonly,omitempty"`
	StatusRules      Relation[ServiceStatusRule] `json:"status_rules,omitzero"`
	Tags             Relation[Tag]               `json:"tags,omitzero"`
	ProblemTags      Relation[ServiceProblemTag] `json:"problem_tags,omitzero"`
	Parents          Relation[Service]           `json:"parents,omitzero"`
	Children         Relation[Service]           `json:"children,omitzero"`
}

// ServiceStatusRule is a status-calculation rule attached to a service.
type ServiceStatusRule struct {
	Type        string `json:"type,omitempty"`
	LimitValue  string `json:"limit_value,omitempty"`
	LimitStatus string `json:"limit_status,omitempty"`
	NewStatus   string `json:"new_status,omitempty"`
}

// ServiceProblemTag maps problems to a service by tag.
type ServiceProblemTag struct {
	Tag      string `json:"tag,omitempty"`
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value,omitempty"`
}

// ServiceGetParams are the parameters for service.get.
// ServiceGetParams queries the list form of service.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type ServiceGetParams struct {
	Output                 any                 `json:"output,omitempty"`
	Filter                 map[string]any      `json:"filter,omitzero"`
	Search                 map[string]any      `json:"search,omitzero"`
	SearchByAny            bool                `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool                `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool                `json:"startSearch,omitempty"`
	ExcludeSearch          bool                `json:"excludeSearch,omitempty"`
	SortField              []string            `json:"sortfield,omitzero"`
	SortOrder              []string            `json:"sortorder,omitzero"`
	Limit                  int                 `json:"limit,omitempty"`
	Editable               bool                `json:"editable,omitempty"`
	ServiceIDs             []string            `json:"serviceids,omitzero"`
	ParentIDs              []string            `json:"parentids,omitzero"`
	ChildIDs               []string            `json:"childids,omitzero"`
	Tags                   []TagFilter         `json:"tags,omitzero"`
	ProblemTags            []ServiceProblemTag `json:"problem_tags,omitzero"`
	SelectChildren         any                 `json:"selectChildren,omitempty"`
	SelectParents          any                 `json:"selectParents,omitempty"`
	SelectTags             any                 `json:"selectTags,omitempty"`
	SelectProblemTags      any                 `json:"selectProblemTags,omitempty"`
	SelectStatusRules      any                 `json:"selectStatusRules,omitempty"`
}

// Get retrieves services matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/service/get
func (s *ServiceService) Get(ctx context.Context, params ServiceGetParams) ([]Service, error) {
	return getList[Service](ctx, s.client, "service.get", params)
}

// ServiceCountParams contains the filters for service.get counts.
type ServiceCountParams struct {
	Filter                 map[string]any      `json:"filter,omitzero"`
	Search                 map[string]any      `json:"search,omitzero"`
	SearchByAny            bool                `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool                `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool                `json:"startSearch,omitempty"`
	ExcludeSearch          bool                `json:"excludeSearch,omitempty"`
	Editable               bool                `json:"editable,omitempty"`
	ServiceIDs             []string            `json:"serviceids,omitzero"`
	ParentIDs              []string            `json:"parentids,omitzero"`
	ChildIDs               []string            `json:"childids,omitzero"`
	Tags                   []TagFilter         `json:"tags,omitzero"`
	ProblemTags            []ServiceProblemTag `json:"problem_tags,omitzero"`
}

// Count returns the number of objects matching params.
func (s *ServiceService) Count(ctx context.Context, params ServiceCountParams) (uint64, error) {
	return getCount(ctx, s.client, "service.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *ServiceService) GetMap(ctx context.Context, params ServiceGetParams) (map[string]Service, error) {
	return getMap[Service](ctx, s.client, "service.get", params)
}

// Create creates services and returns the new service ids.
func (s *ServiceService) Create(ctx context.Context, items ...ServiceCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "service.create", "serviceids", items)
}

// Update updates services and returns the affected service ids.
func (s *ServiceService) Update(ctx context.Context, items ...ServiceUpdateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "service.update", "serviceids", items)
}

// Delete deletes services by id and returns the deleted service ids.
func (s *ServiceService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "service.delete", "serviceids", ids)
}
