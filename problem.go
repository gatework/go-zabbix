package zabbix

import "context"

// ProblemService wraps the "problem" API namespace.
type ProblemService struct{ client *Client }

// Problem returns the problem service.
func (c *Client) Problem() *ProblemService { return &ProblemService{c} }

// Problem object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/problem/object
type Problem struct {
	EventID       string `json:"eventid,omitempty"`
	Source        string `json:"source,omitempty"`
	Object        string `json:"object,omitempty"`
	ObjectID      string `json:"objectid,omitempty"`
	Clock         string `json:"clock,omitempty"`
	Ns            string `json:"ns,omitempty"`
	REventID      string `json:"r_eventid,omitempty"`
	RClock        string `json:"r_clock,omitempty"`
	RNs           string `json:"r_ns,omitempty"`
	CorrelationID string `json:"correlationid,omitempty"`
	UserID        string `json:"userid,omitempty"`
	Name          string `json:"name,omitempty"`
	Acknowledged  string `json:"acknowledged,omitempty"`
	Severity      string `json:"severity,omitempty"`
	Suppressed    string `json:"suppressed,omitempty"`
	OpData        string `json:"opdata,omitempty"`
	Tags          []Tag  `json:"tags,omitempty"`
}

// ProblemGetParams are the parameters for problem.get.
type ProblemGetParams struct {
	GetParams
	EventIDs              []string `json:"eventids,omitempty"`
	GroupIDs              []string `json:"groupids,omitempty"`
	HostIDs               []string `json:"hostids,omitempty"`
	ObjectIDs             []string `json:"objectids,omitempty"`
	Source                *int     `json:"source,omitempty"`
	Object                *int     `json:"object,omitempty"`
	Severities            []int    `json:"severities,omitempty"`
	Tags                  []Tag    `json:"tags,omitempty"`
	EvalType              int      `json:"evaltype,omitempty"`
	Recent                *bool    `json:"recent,omitempty"`
	Acknowledged          *bool    `json:"acknowledged,omitempty"`
	Suppressed            *bool    `json:"suppressed,omitempty"`
	TimeFrom              int64    `json:"time_from,omitempty"`
	TimeTill              int64    `json:"time_till,omitempty"`
	SelectAcknowledges    any      `json:"selectAcknowledges,omitempty"`
	SelectTags            any      `json:"selectTags,omitempty"`
	SelectSuppressionData any      `json:"selectSuppressionData,omitempty"`
}

// Get retrieves problems matching params. Note the problem API is read-only.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/problem/get
func (s *ProblemService) Get(ctx context.Context, params ProblemGetParams) ([]Problem, error) {
	var problems []Problem
	err := s.client.Call(ctx, "problem.get", params, &problems)
	return problems, err
}
