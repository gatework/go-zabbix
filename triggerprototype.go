package zabbix

import "context"

// TriggerPrototypeService wraps the "triggerprototype" API namespace.
type TriggerPrototypeService struct{ client *Client }

// TriggerPrototype returns the triggerprototype service.
func (c *Client) TriggerPrototype() *TriggerPrototypeService { return &TriggerPrototypeService{c} }

// TriggerPrototype object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/triggerprototype/object
type TriggerPrototype struct {
	TriggerID          string             `json:"triggerid,omitempty"`
	Description        string             `json:"description,omitempty"`
	Expression         string             `json:"expression,omitempty"`
	EventName          string             `json:"event_name,omitempty"`
	OpData             string             `json:"opdata,omitempty"`
	Comments           string             `json:"comments,omitempty"`
	Priority           string             `json:"priority,omitempty"`
	Status             string             `json:"status,omitempty"`
	Value              string             `json:"value,omitempty"`
	State              string             `json:"state,omitempty"`
	Error              string             `json:"error,omitempty"`
	TemplateID         string             `json:"templateid,omitempty"`
	Type               string             `json:"type,omitempty"`
	Flags              string             `json:"flags,omitempty"`
	RecoveryMode       string             `json:"recovery_mode,omitempty"`
	RecoveryExpression string             `json:"recovery_expression,omitempty"`
	CorrelationMode    string             `json:"correlation_mode,omitempty"`
	CorrelationTag     string             `json:"correlation_tag,omitempty"`
	ManualClose        string             `json:"manual_close,omitempty"`
	URL                string             `json:"url,omitempty"`
	URLName            string             `json:"url_name,omitempty"`
	UUID               string             `json:"uuid,omitempty"`
	Discover           string             `json:"discover,omitempty"`
	Tags               []Tag              `json:"tags,omitempty"`
	Hosts              []Host             `json:"hosts,omitempty"`
	Items              []ItemPrototype    `json:"items,omitempty"`
	Dependencies       []TriggerPrototype `json:"dependencies,omitempty"`
}

// TriggerPrototypeGetParams are the parameters for triggerprototype.get.
// TriggerPrototypeGetParams queries the list form of triggerprototype.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type TriggerPrototypeGetParams struct {
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
	DiscoveryIDs           []string       `json:"discoveryids,omitzero"`
	SelectHosts            any            `json:"selectHosts,omitempty"`
	SelectItems            any            `json:"selectItems,omitempty"`
	SelectTags             any            `json:"selectTags,omitempty"`
	SelectDependencies     any            `json:"selectDependencies,omitempty"`
}

// Get retrieves trigger prototypes matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/triggerprototype/get
func (s *TriggerPrototypeService) Get(ctx context.Context, params TriggerPrototypeGetParams) ([]TriggerPrototype, error) {
	return getList[TriggerPrototype](ctx, s.client, "triggerprototype.get", params)
}

// TriggerPrototypeCountParams contains the filters for triggerprototype.get counts.
type TriggerPrototypeCountParams struct {
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
	DiscoveryIDs           []string       `json:"discoveryids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *TriggerPrototypeService) Count(ctx context.Context, params TriggerPrototypeCountParams) (uint64, error) {
	return getCount(ctx, s.client, "triggerprototype.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *TriggerPrototypeService) GetMap(ctx context.Context, params TriggerPrototypeGetParams) (map[string]TriggerPrototype, error) {
	return getMap[TriggerPrototype](ctx, s.client, "triggerprototype.get", params)
}

// Create creates trigger prototypes and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/triggerprototype/create
func (s *TriggerPrototypeService) Create(ctx context.Context, triggers ...TriggerPrototypeCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "triggerprototype.create", "triggerids", triggers)
}

// Update updates trigger prototypes and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/triggerprototype/update
func (s *TriggerPrototypeService) Update(ctx context.Context, triggers ...TriggerPrototypeUpdateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "triggerprototype.update", "triggerids", triggers)
}

// Delete deletes trigger prototypes by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/triggerprototype/delete
func (s *TriggerPrototypeService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "triggerprototype.delete", "triggerids", ids)
}
