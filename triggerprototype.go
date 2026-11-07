package zabbix

import "context"

// TriggerPrototypeService wraps the "triggerprototype" API namespace.
type TriggerPrototypeService struct{ client *Client }

// TriggerPrototype returns the triggerprototype service.
func (c *Client) TriggerPrototype() *TriggerPrototypeService { return &TriggerPrototypeService{c} }

// TriggerPrototype object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/triggerprototype/object
type TriggerPrototype struct {
	TriggerID          string `json:"triggerid,omitempty"`
	Description        string `json:"description,omitempty"`
	Expression         string `json:"expression,omitempty"`
	EventName          string `json:"event_name,omitempty"`
	OpData             string `json:"opdata,omitempty"`
	Comments           string `json:"comments,omitempty"`
	Priority           string `json:"priority,omitempty"`
	Status             string `json:"status,omitempty"`
	Value              string `json:"value,omitempty"`
	State              string `json:"state,omitempty"`
	Error              string `json:"error,omitempty"`
	TemplateID         string `json:"templateid,omitempty"`
	Type               string `json:"type,omitempty"`
	Flags              string `json:"flags,omitempty"`
	RecoveryMode       string `json:"recovery_mode,omitempty"`
	RecoveryExpression string `json:"recovery_expression,omitempty"`
	CorrelationMode    string `json:"correlation_mode,omitempty"`
	CorrelationTag     string `json:"correlation_tag,omitempty"`
	ManualClose        string `json:"manual_close,omitempty"`
	URL                string `json:"url,omitempty"`
	URLName            string `json:"url_name,omitempty"`
	UUID               string `json:"uuid,omitempty"`
	Discover           string `json:"discover,omitempty"`
	Tags               []Tag  `json:"tags,omitempty"`
}

// TriggerPrototypeGetParams are the parameters for triggerprototype.get.
type TriggerPrototypeGetParams struct {
	GetParams
	TriggerIDs         []string `json:"triggerids,omitempty"`
	GroupIDs           []string `json:"groupids,omitempty"`
	TemplateIDs        []string `json:"templateids,omitempty"`
	HostIDs            []string `json:"hostids,omitempty"`
	ItemIDs            []string `json:"itemids,omitempty"`
	DiscoveryIDs       []string `json:"discoveryids,omitempty"`
	SelectHosts        any      `json:"selectHosts,omitempty"`
	SelectItems        any      `json:"selectItems,omitempty"`
	SelectTags         any      `json:"selectTags,omitempty"`
	SelectDependencies any      `json:"selectDependencies,omitempty"`
}

// Get retrieves trigger prototypes matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/triggerprototype/get
func (s *TriggerPrototypeService) Get(ctx context.Context, params TriggerPrototypeGetParams) ([]TriggerPrototype, error) {
	var out []TriggerPrototype
	err := s.client.Call(ctx, "triggerprototype.get", params, &out)
	return out, err
}

// Create creates trigger prototypes and returns the new ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/triggerprototype/create
func (s *TriggerPrototypeService) Create(ctx context.Context, triggers ...TriggerPrototype) ([]string, error) {
	var res struct {
		IDs []string `json:"triggerids"`
	}
	err := s.client.Call(ctx, "triggerprototype.create", triggers, &res)
	return res.IDs, err
}

// Update updates trigger prototypes and returns the affected ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/triggerprototype/update
func (s *TriggerPrototypeService) Update(ctx context.Context, triggers ...TriggerPrototype) ([]string, error) {
	var res struct {
		IDs []string `json:"triggerids"`
	}
	err := s.client.Call(ctx, "triggerprototype.update", triggers, &res)
	return res.IDs, err
}

// Delete deletes trigger prototypes by id and returns the deleted ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/triggerprototype/delete
func (s *TriggerPrototypeService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"triggerids"`
	}
	err := s.client.Call(ctx, "triggerprototype.delete", ids, &res)
	return res.IDs, err
}
