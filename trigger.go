package zabbix

import "context"

// TriggerService wraps the "trigger" API namespace.
type TriggerService struct{ client *Client }

// Trigger returns the trigger service.
func (c *Client) Trigger() *TriggerService { return &TriggerService{c} }

// Trigger object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/trigger/object
type Trigger struct {
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
	Tags               []Tag  `json:"tags,omitempty"`
}

// TriggerGetParams are the parameters for trigger.get.
type TriggerGetParams struct {
	GetParams
	TriggerIDs         []string `json:"triggerids,omitempty"`
	GroupIDs           []string `json:"groupids,omitempty"`
	TemplateIDs        []string `json:"templateids,omitempty"`
	HostIDs            []string `json:"hostids,omitempty"`
	ItemIDs            []string `json:"itemids,omitempty"`
	Tags               []Tag    `json:"tags,omitempty"`
	EvalType           int      `json:"evaltype,omitempty"`
	Monitored          bool     `json:"monitored,omitempty"`
	Active             bool     `json:"active,omitempty"`
	Maintenance        *bool    `json:"maintenance,omitempty"`
	MinSeverity        int      `json:"min_severity,omitempty"`
	OnlyTrue           bool     `json:"only_true,omitempty"`
	SkipDependent      bool     `json:"skipDependent,omitempty"`
	ExpandExpression   bool     `json:"expandExpression,omitempty"`
	SelectHosts        any      `json:"selectHosts,omitempty"`
	SelectItems        any      `json:"selectItems,omitempty"`
	SelectTags         any      `json:"selectTags,omitempty"`
	SelectDependencies any      `json:"selectDependencies,omitempty"`
}

// Get retrieves triggers matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/trigger/get
func (s *TriggerService) Get(ctx context.Context, params TriggerGetParams) ([]Trigger, error) {
	var triggers []Trigger
	err := s.client.Call(ctx, "trigger.get", params, &triggers)
	return triggers, err
}

// Create creates triggers and returns the new trigger ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/trigger/create
func (s *TriggerService) Create(ctx context.Context, triggers ...Trigger) ([]string, error) {
	var res struct {
		TriggerIDs []string `json:"triggerids"`
	}
	err := s.client.Call(ctx, "trigger.create", triggers, &res)
	return res.TriggerIDs, err
}

// Update updates triggers and returns the affected trigger ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/trigger/update
func (s *TriggerService) Update(ctx context.Context, triggers ...Trigger) ([]string, error) {
	var res struct {
		TriggerIDs []string `json:"triggerids"`
	}
	err := s.client.Call(ctx, "trigger.update", triggers, &res)
	return res.TriggerIDs, err
}

// Delete deletes triggers by id and returns the deleted trigger ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/trigger/delete
func (s *TriggerService) Delete(ctx context.Context, triggerIDs ...string) ([]string, error) {
	var res struct {
		TriggerIDs []string `json:"triggerids"`
	}
	err := s.client.Call(ctx, "trigger.delete", triggerIDs, &res)
	return res.TriggerIDs, err
}
