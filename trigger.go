package zabbix

import "context"

// TriggerService wraps the "trigger" API namespace.
type TriggerService struct{ client *Client }

// Trigger returns the trigger service.
func (c *Client) Trigger() *TriggerService { return &TriggerService{c} }

// Trigger object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/trigger/object
type Trigger struct {
	TriggerID          string    `json:"triggerid,omitempty"`
	Description        string    `json:"description,omitempty"`
	Expression         string    `json:"expression,omitempty"`
	EventName          string    `json:"event_name,omitempty"`
	OpData             string    `json:"opdata,omitempty"`
	Comments           string    `json:"comments,omitempty"`
	Priority           string    `json:"priority,omitempty"`
	Status             string    `json:"status,omitempty"`
	Value              string    `json:"value,omitempty"`
	State              string    `json:"state,omitempty"`
	Error              string    `json:"error,omitempty"`
	TemplateID         string    `json:"templateid,omitempty"`
	Type               string    `json:"type,omitempty"`
	Flags              string    `json:"flags,omitempty"`
	RecoveryMode       string    `json:"recovery_mode,omitempty"`
	RecoveryExpression string    `json:"recovery_expression,omitempty"`
	CorrelationMode    string    `json:"correlation_mode,omitempty"`
	CorrelationTag     string    `json:"correlation_tag,omitempty"`
	ManualClose        string    `json:"manual_close,omitempty"`
	URL                string    `json:"url,omitempty"`
	URLName            string    `json:"url_name,omitempty"`
	UUID               string    `json:"uuid,omitempty"`
	Tags               []Tag     `json:"tags,omitempty"`
	Hosts              []Host    `json:"hosts,omitempty"`
	Items              []Item    `json:"items,omitempty"`
	Dependencies       []Trigger `json:"dependencies,omitempty"`
}

// TriggerGetParams are the parameters for trigger.get.
// TriggerGetParams queries the list form of trigger.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type TriggerGetParams struct {
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
	TriggerIDs             []string       `json:"triggerids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	TemplateIDs            []string       `json:"templateids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	ItemIDs                []string       `json:"itemids,omitzero"`
	Tags                   []TagFilter    `json:"tags,omitzero"`
	EvalType               int            `json:"evaltype,omitempty"`
	Monitored              bool           `json:"monitored,omitempty"`
	Active                 bool           `json:"active,omitempty"`
	Maintenance            *bool          `json:"maintenance,omitempty"`
	MinSeverity            int            `json:"min_severity,omitempty"`
	OnlyTrue               bool           `json:"only_true,omitempty"`
	SkipDependent          bool           `json:"skipDependent,omitempty"`
	ExpandExpression       bool           `json:"expandExpression,omitempty"`
	SelectHosts            any            `json:"selectHosts,omitempty"`
	SelectItems            any            `json:"selectItems,omitempty"`
	SelectTags             any            `json:"selectTags,omitempty"`
	SelectDependencies     any            `json:"selectDependencies,omitempty"`
}

// Get retrieves triggers matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/trigger/get
func (s *TriggerService) Get(ctx context.Context, params TriggerGetParams) ([]Trigger, error) {
	return getList[Trigger](ctx, s.client, "trigger.get", params)
}

// TriggerCountParams contains the filters for trigger.get counts.
type TriggerCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	TriggerIDs             []string       `json:"triggerids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	TemplateIDs            []string       `json:"templateids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	ItemIDs                []string       `json:"itemids,omitzero"`
	Tags                   []TagFilter    `json:"tags,omitzero"`
	EvalType               int            `json:"evaltype,omitempty"`
	Monitored              bool           `json:"monitored,omitempty"`
	Active                 bool           `json:"active,omitempty"`
	Maintenance            *bool          `json:"maintenance,omitempty"`
	MinSeverity            int            `json:"min_severity,omitempty"`
	OnlyTrue               bool           `json:"only_true,omitempty"`
	SkipDependent          bool           `json:"skipDependent,omitempty"`
}

// Count returns the number of objects matching params.
func (s *TriggerService) Count(ctx context.Context, params TriggerCountParams) (uint64, error) {
	return getCount(ctx, s.client, "trigger.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *TriggerService) GetMap(ctx context.Context, params TriggerGetParams) (map[string]Trigger, error) {
	return getMap[Trigger](ctx, s.client, "trigger.get", params)
}

// Create creates triggers and returns the new trigger ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/trigger/create
func (s *TriggerService) Create(ctx context.Context, triggers ...TriggerCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "trigger.create", "triggerids", triggers)
}

// Update updates triggers and returns the affected trigger ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/trigger/update
func (s *TriggerService) Update(ctx context.Context, triggers ...TriggerUpdateParams) ([]string, error) {
	return updateIDs(ctx, s.client, "trigger.update", "triggerids", triggers)
}

// Delete deletes triggers by id and returns the deleted trigger ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/trigger/delete
func (s *TriggerService) Delete(ctx context.Context, triggerIDs ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "trigger.delete", "triggerids", triggerIDs)
}
