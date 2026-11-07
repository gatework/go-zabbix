package zabbix

import "context"

// ActionService wraps the "action" API namespace.
type ActionService struct{ client *Client }

// Action returns the action service.
func (c *Client) Action() *ActionService { return &ActionService{c} }

// Action object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/action/object
type Action struct {
	ActionID           string                    `json:"actionid,omitempty"`
	Name               string                    `json:"name,omitempty"`
	EventSource        string                    `json:"eventsource,omitempty"`
	Status             string                    `json:"status,omitempty"`
	EscPeriod          string                    `json:"esc_period,omitempty"`
	PauseSuppressed    string                    `json:"pause_suppressed,omitempty"`
	PauseSymptoms      string                    `json:"pause_symptoms,omitempty"`
	NotifyIfCanceled   string                    `json:"notify_if_canceled,omitempty"`
	Operations         []ActionOperation         `json:"operations,omitempty"`
	RecoveryOperations []ActionRecoveryOperation `json:"recovery_operations,omitempty"`
	UpdateOperations   []ActionUpdateOperation   `json:"update_operations,omitempty"`
	Filter             *ActionFilter             `json:"filter,omitempty"`
}

// ActionFilter is the condition filter of an action.
type ActionFilter struct {
	EvalType   string            `json:"evaltype,omitempty"`
	Formula    string            `json:"formula,omitempty"`
	Conditions []ActionCondition `json:"conditions,omitempty"`
}

// ActionCondition is a single condition inside an ActionFilter.
type ActionCondition struct {
	ConditionType string `json:"conditiontype,omitempty"`
	Operator      string `json:"operator,omitempty"`
	Value         string `json:"value,omitempty"`
	Value2        string `json:"value2,omitempty"`
	FormulaID     string `json:"formulaid,omitempty"`
}

// ActionOperation is an operation performed when an action triggers.
type ActionOperation struct {
	OperationID   string           `json:"operationid,omitempty"`
	OperationType string           `json:"operationtype,omitempty"`
	EscPeriod     string           `json:"esc_period,omitempty"`
	EscStepFrom   string           `json:"esc_step_from,omitempty"`
	EscStepTo     string           `json:"esc_step_to,omitempty"`
	EvalType      string           `json:"evaltype,omitempty"`
	OpMessage     *ActionOpMessage `json:"opmessage,omitempty"`
	OpMessageGrp  any              `json:"opmessage_grp,omitempty"`
	OpMessageUsr  any              `json:"opmessage_usr,omitempty"`
	OpCommand     any              `json:"opcommand,omitempty"`
	OpConditions  any              `json:"opconditions,omitempty"`
}

// ActionOpMessage is the message settings of an action operation.
type ActionOpMessage struct {
	DefaultMsg  string `json:"default_msg,omitempty"`
	MediaTypeID string `json:"mediatypeid,omitempty"`
	Subject     string `json:"subject,omitempty"`
	Message     string `json:"message,omitempty"`
}

// ActionRecoveryOperation is an operation performed on problem recovery.
type ActionRecoveryOperation struct {
	OperationID   string `json:"operationid,omitempty"`
	OperationType string `json:"operationtype,omitempty"`
	OpMessage     any    `json:"opmessage,omitempty"`
}

// ActionUpdateOperation is an operation performed on problem update.
type ActionUpdateOperation struct {
	OperationID   string `json:"operationid,omitempty"`
	OperationType string `json:"operationtype,omitempty"`
	OpMessage     any    `json:"opmessage,omitempty"`
}

// ActionGetParams are the parameters for action.get.
type ActionGetParams struct {
	GetParams
	ActionIDs                []string `json:"actionids,omitempty"`
	GroupIDs                 []string `json:"groupids,omitempty"`
	HostIDs                  []string `json:"hostids,omitempty"`
	TriggerIDs               []string `json:"triggerids,omitempty"`
	MediaTypeIDs             []string `json:"mediatypeids,omitempty"`
	UserGroupIDs             []string `json:"usrgrpids,omitempty"`
	UserIDs                  []string `json:"userids,omitempty"`
	SelectOperations         any      `json:"selectOperations,omitempty"`
	SelectRecoveryOperations any      `json:"selectRecoveryOperations,omitempty"`
	SelectUpdateOperations   any      `json:"selectUpdateOperations,omitempty"`
	SelectFilter             any      `json:"selectFilter,omitempty"`
}

// Get retrieves actions matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/action/get
func (s *ActionService) Get(ctx context.Context, params ActionGetParams) ([]Action, error) {
	var out []Action
	err := s.client.Call(ctx, "action.get", params, &out)
	return out, err
}

// Create creates actions and returns the new action ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/action/create
func (s *ActionService) Create(ctx context.Context, items ...Action) ([]string, error) {
	var res struct {
		IDs []string `json:"actionids"`
	}
	err := s.client.Call(ctx, "action.create", items, &res)
	return res.IDs, err
}

// Update updates actions and returns the affected action ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/action/update
func (s *ActionService) Update(ctx context.Context, items ...Action) ([]string, error) {
	var res struct {
		IDs []string `json:"actionids"`
	}
	err := s.client.Call(ctx, "action.update", items, &res)
	return res.IDs, err
}

// Delete deletes actions by id and returns the deleted action ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/action/delete
func (s *ActionService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"actionids"`
	}
	err := s.client.Call(ctx, "action.delete", ids, &res)
	return res.IDs, err
}
