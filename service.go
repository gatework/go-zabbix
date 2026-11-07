package zabbix

import "context"

// ServiceService wraps the "service" API namespace.
type ServiceService struct{ client *Client }

// Service returns the service service.
func (c *Client) Service() *ServiceService { return &ServiceService{c} }

// Service object (Zabbix 7.x).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/service/object
type Service struct {
	ServiceID        string `json:"serviceid,omitempty"`
	Name             string `json:"name,omitempty"`
	Status           string `json:"status,omitempty"`
	Algorithm        string `json:"algorithm,omitempty"`
	SortOrder        string `json:"sortorder,omitempty"`
	Weight           string `json:"weight,omitempty"`
	PropagationRule  string `json:"propagation_rule,omitempty"`
	PropagationValue string `json:"propagation_value,omitempty"`
	Description      string `json:"description,omitempty"`
	UUID             string `json:"uuid,omitempty"`
	CreatedAt        string `json:"created_at,omitempty"`
	Readonly         string `json:"readonly,omitempty"`

	StatusRules []ServiceStatusRule `json:"status_rules,omitempty"`
	Tags        []Tag               `json:"tags,omitempty"`
	ProblemTags []ServiceProblemTag `json:"problem_tags,omitempty"`
	Parents     any                 `json:"parents,omitempty"`
	Children    any                 `json:"children,omitempty"`
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
type ServiceGetParams struct {
	GetParams
	ServiceIDs        []string            `json:"serviceids,omitempty"`
	ParentIDs         []string            `json:"parentids,omitempty"`
	ChildIDs          []string            `json:"childids,omitempty"`
	Tags              []Tag               `json:"tags,omitempty"`
	ProblemTags       []ServiceProblemTag `json:"problem_tags,omitempty"`
	SelectChildren    any                 `json:"selectChildren,omitempty"`
	SelectParents     any                 `json:"selectParents,omitempty"`
	SelectTags        any                 `json:"selectTags,omitempty"`
	SelectProblemTags any                 `json:"selectProblemTags,omitempty"`
	SelectStatusRules any                 `json:"selectStatusRules,omitempty"`
}

// Get retrieves services matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/service/get
func (s *ServiceService) Get(ctx context.Context, params ServiceGetParams) ([]Service, error) {
	var out []Service
	err := s.client.Call(ctx, "service.get", params, &out)
	return out, err
}

// Create creates services and returns the new service ids.
func (s *ServiceService) Create(ctx context.Context, items ...Service) ([]string, error) {
	var res struct {
		IDs []string `json:"serviceids"`
	}
	err := s.client.Call(ctx, "service.create", items, &res)
	return res.IDs, err
}

// Update updates services and returns the affected service ids.
func (s *ServiceService) Update(ctx context.Context, items ...Service) ([]string, error) {
	var res struct {
		IDs []string `json:"serviceids"`
	}
	err := s.client.Call(ctx, "service.update", items, &res)
	return res.IDs, err
}

// Delete deletes services by id and returns the deleted service ids.
func (s *ServiceService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"serviceids"`
	}
	err := s.client.Call(ctx, "service.delete", ids, &res)
	return res.IDs, err
}
