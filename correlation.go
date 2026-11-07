package zabbix

import "context"

// CorrelationService wraps the "correlation" API namespace.
type CorrelationService struct{ client *Client }

// Correlation returns the correlation service.
func (c *Client) Correlation() *CorrelationService { return &CorrelationService{c} }

// Correlation object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/correlation/object
type Correlation struct {
	CorrelationID string `json:"correlationid,omitempty"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	Status        string `json:"status,omitempty"`

	Filter     *CorrelationFilter     `json:"filter,omitempty"`
	Operations []CorrelationOperation `json:"operations,omitempty"`
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
type CorrelationGetParams struct {
	GetParams
	CorrelationIDs   []string `json:"correlationids,omitempty"`
	SelectFilter     any      `json:"selectFilter,omitempty"`
	SelectOperations any      `json:"selectOperations,omitempty"`
}

// Get retrieves correlations matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/correlation/get
func (s *CorrelationService) Get(ctx context.Context, params CorrelationGetParams) ([]Correlation, error) {
	var out []Correlation
	err := s.client.Call(ctx, "correlation.get", params, &out)
	return out, err
}

// Create creates correlations and returns the new correlation ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/correlation/create
func (s *CorrelationService) Create(ctx context.Context, items ...Correlation) ([]string, error) {
	var res struct {
		IDs []string `json:"correlationids"`
	}
	err := s.client.Call(ctx, "correlation.create", items, &res)
	return res.IDs, err
}

// Update updates correlations and returns the affected correlation ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/correlation/update
func (s *CorrelationService) Update(ctx context.Context, items ...Correlation) ([]string, error) {
	var res struct {
		IDs []string `json:"correlationids"`
	}
	err := s.client.Call(ctx, "correlation.update", items, &res)
	return res.IDs, err
}

// Delete deletes correlations by id and returns the deleted correlation ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/correlation/delete
func (s *CorrelationService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"correlationids"`
	}
	err := s.client.Call(ctx, "correlation.delete", ids, &res)
	return res.IDs, err
}
