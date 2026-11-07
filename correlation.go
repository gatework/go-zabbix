package zabbix

import "context"

// CorrelationService wraps the "correlation" API namespace.
type CorrelationService struct{ client *Client }

// Correlation returns the correlation service.
func (c *Client) Correlation() *CorrelationService { return &CorrelationService{c} }

// Correlation object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/correlation/object
type Correlation struct {
	CorrelationID string                            `json:"correlationid,omitempty"`
	Name          string                            `json:"name,omitempty"`
	Description   string                            `json:"description,omitempty"`
	Status        string                            `json:"status,omitempty"`
	Filter        SelectedObject[CorrelationFilter] `json:"filter,omitzero"`
	Operations    []CorrelationOperation            `json:"operations,omitempty"`
}

// CorrelationFilter is the condition filter of a correlation.
type CorrelationFilter struct {
	EvalType    string                 `json:"evaltype,omitempty"`
	Formula     string                 `json:"formula,omitempty"`
	Conditions  []CorrelationCondition `json:"conditions,omitempty"`
	EvalFormula string                 `json:"eval_formula,omitempty"`
}

// CorrelationCondition is a single condition inside a CorrelationFilter.
type CorrelationCondition struct {
	Type      string `json:"type,omitempty"`
	Tag       string `json:"tag,omitempty"`
	OldTag    string `json:"oldtag,omitempty"`
	NewTag    string `json:"newtag,omitempty"`
	Value     string `json:"value,omitempty"`
	FormulaID string `json:"formulaid,omitempty"`
	Operator  string `json:"operator,omitempty"`
}

// CorrelationOperation is an operation performed when a correlation matches.
type CorrelationOperation struct {
	Type string `json:"type,omitempty"`
}

// CorrelationGetParams are the parameters for correlation.get.
// CorrelationGetParams queries the list form of correlation.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type CorrelationGetParams struct {
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
	CorrelationIDs         []string       `json:"correlationids,omitzero"`
	SelectFilter           any            `json:"selectFilter,omitempty"`
	SelectOperations       any            `json:"selectOperations,omitempty"`
}

// Get retrieves correlations matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/correlation/get
func (s *CorrelationService) Get(ctx context.Context, params CorrelationGetParams) ([]Correlation, error) {
	return getList[Correlation](ctx, s.client, "correlation.get", params)
}

// CorrelationCountParams contains the filters for correlation.get counts.
type CorrelationCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	CorrelationIDs         []string       `json:"correlationids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *CorrelationService) Count(ctx context.Context, params CorrelationCountParams) (uint64, error) {
	return getCount(ctx, s.client, "correlation.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *CorrelationService) GetMap(ctx context.Context, params CorrelationGetParams) (map[string]Correlation, error) {
	return getMap[Correlation](ctx, s.client, "correlation.get", params)
}

// Create creates correlations and returns the new correlation ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/correlation/create
func (s *CorrelationService) Create(ctx context.Context, items ...CorrelationCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "correlation.create", "correlationids", items)
}

// Update updates correlations and returns the affected correlation ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/correlation/update
func (s *CorrelationService) Update(ctx context.Context, items ...CorrelationUpdateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "correlation.update", "correlationids", items)
}

// Delete deletes correlations by id and returns the deleted correlation ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/correlation/delete
func (s *CorrelationService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "correlation.delete", "correlationids", ids)
}
