package zabbix

import "context"

// RegexpService wraps the "regexp" API namespace (global regular expressions).
type RegexpService struct{ client *Client }

// Regexp returns the regexp service.
func (c *Client) Regexp() *RegexpService { return &RegexpService{c} }

// Regexp is a global regular expression (Zabbix 7.0+).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/regexp/object
type Regexp struct {
	RegexpID   string `json:"regexpid,omitempty"`
	Name       string `json:"name,omitempty"`
	TestString string `json:"test_string,omitempty"`

	Expressions []RegexpExpression `json:"expressions,omitempty"`
}

// RegexpExpression is one expression that makes up a global regular expression.
type RegexpExpression struct {
	ExpressionID   string `json:"expressionid,omitempty"`
	RegexpID       string `json:"regexpid,omitempty"`
	Expression     string `json:"expression,omitempty"`
	ExpressionType string `json:"expression_type,omitempty"`
	ExpDelimiter   string `json:"exp_delimiter,omitempty"`
	CaseSensitive  string `json:"case_sensitive,omitempty"`
}

// RegexpGetParams are the parameters for regexp.get.
// RegexpGetParams queries the list form of regexp.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type RegexpGetParams struct {
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
	RegexpIDs              []string       `json:"regexpids,omitzero"`
	SelectExpressions      any            `json:"selectExpressions,omitempty"`
}

// Get retrieves global regular expressions matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/regexp/get
func (s *RegexpService) Get(ctx context.Context, params RegexpGetParams) ([]Regexp, error) {
	return getList[Regexp](ctx, s.client, "regexp.get", params)
}

// RegexpCountParams contains the filters for regexp.get counts.
type RegexpCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	RegexpIDs              []string       `json:"regexpids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *RegexpService) Count(ctx context.Context, params RegexpCountParams) (uint64, error) {
	return getCount(ctx, s.client, "regexp.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *RegexpService) GetMap(ctx context.Context, params RegexpGetParams) (map[string]Regexp, error) {
	return getMap[Regexp](ctx, s.client, "regexp.get", params)
}

// Create creates global regular expressions and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/regexp/create
func (s *RegexpService) Create(ctx context.Context, regexps ...RegexpCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "regexp.create", "regexpids", regexps)
}

// Update updates global regular expressions and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/regexp/update
func (s *RegexpService) Update(ctx context.Context, regexps ...RegexpUpdateParams) ([]string, error) {
	return updateIDs(ctx, s.client, "regexp.update", "regexpids", regexps)
}

// Delete deletes global regular expressions by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/regexp/delete
func (s *RegexpService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "regexp.delete", "regexpids", ids)
}
