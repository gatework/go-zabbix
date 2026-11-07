package zabbix

import "context"

// RegexpService wraps the "regexp" API namespace (global regular expressions).
type RegexpService struct{ client *Client }

// Regexp returns the regexp service.
func (c *Client) Regexp() *RegexpService { return &RegexpService{c} }

// Regexp is a global regular expression (Zabbix 7.0+).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/regexp/object
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
type RegexpGetParams struct {
	GetParams
	RegexpIDs         []string `json:"regexpids,omitempty"`
	SelectExpressions any      `json:"selectExpressions,omitempty"`
}

// Get retrieves global regular expressions matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/regexp/get
func (s *RegexpService) Get(ctx context.Context, params RegexpGetParams) ([]Regexp, error) {
	var out []Regexp
	err := s.client.Call(ctx, "regexp.get", params, &out)
	return out, err
}

// Create creates global regular expressions and returns the new ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/regexp/create
func (s *RegexpService) Create(ctx context.Context, regexps ...Regexp) ([]string, error) {
	var res struct {
		IDs []string `json:"regexpids"`
	}
	err := s.client.Call(ctx, "regexp.create", regexps, &res)
	return res.IDs, err
}

// Update updates global regular expressions and returns the affected ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/regexp/update
func (s *RegexpService) Update(ctx context.Context, regexps ...Regexp) ([]string, error) {
	var res struct {
		IDs []string `json:"regexpids"`
	}
	err := s.client.Call(ctx, "regexp.update", regexps, &res)
	return res.IDs, err
}

// Delete deletes global regular expressions by id and returns the deleted ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/regexp/delete
func (s *RegexpService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"regexpids"`
	}
	err := s.client.Call(ctx, "regexp.delete", ids, &res)
	return res.IDs, err
}
