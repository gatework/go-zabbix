package zabbix

import "context"

// ProblemService wraps the "problem" API namespace.
type ProblemService struct{ client *Client }

// Problem returns the problem service.
func (c *Client) Problem() *ProblemService { return &ProblemService{c} }

// Problem object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/problem/object
type Problem struct {
	EventID         string                `json:"eventid,omitempty"`
	Source          string                `json:"source,omitempty"`
	Object          string                `json:"object,omitempty"`
	ObjectID        string                `json:"objectid,omitempty"`
	Clock           string                `json:"clock,omitempty"`
	NS              string                `json:"ns,omitempty"`
	REventID        string                `json:"r_eventid,omitempty"`
	RClock          string                `json:"r_clock,omitempty"`
	RNS             string                `json:"r_ns,omitempty"`
	CorrelationID   string                `json:"correlationid,omitempty"`
	UserID          string                `json:"userid,omitempty"`
	Name            string                `json:"name,omitempty"`
	Acknowledged    string                `json:"acknowledged,omitempty"`
	Severity        string                `json:"severity,omitempty"`
	Suppressed      string                `json:"suppressed,omitempty"`
	OpData          string                `json:"opdata,omitempty"`
	Tags            []Tag                 `json:"tags,omitempty"`
	Acknowledges    Relation[EventUpdate] `json:"acknowledges,omitzero"`
	SuppressionData []SuppressionData     `json:"suppression_data,omitempty"`
}

// ProblemGetParams are the parameters for problem.get.
// ProblemGetParams queries the list form of problem.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type ProblemGetParams struct {
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
	EventIDs               []string       `json:"eventids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	ObjectIDs              []string       `json:"objectids,omitzero"`
	Source                 *int           `json:"source,omitempty"`
	Object                 *int           `json:"object,omitempty"`
	Severities             []int          `json:"severities,omitzero"`
	Tags                   []TagFilter    `json:"tags,omitzero"`
	EvalType               int            `json:"evaltype,omitempty"`
	Recent                 *bool          `json:"recent,omitempty"`
	Acknowledged           *bool          `json:"acknowledged,omitempty"`
	Suppressed             *bool          `json:"suppressed,omitempty"`
	TimeFrom               *int64         `json:"time_from,omitempty"`
	TimeTill               *int64         `json:"time_till,omitempty"`
	SelectAcknowledges     any            `json:"selectAcknowledges,omitempty"`
	SelectTags             any            `json:"selectTags,omitempty"`
	SelectSuppressionData  any            `json:"selectSuppressionData,omitempty"`
}

// SuppressionData identifies a maintenance or manual suppression of a problem.
type SuppressionData struct {
	MaintenanceID string `json:"maintenanceid,omitempty"`
	UserID        string `json:"userid,omitempty"`
	SuppressUntil string `json:"suppress_until,omitempty"`
}

// Get retrieves problems matching params. Note the problem API is read-only.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/problem/get
func (s *ProblemService) Get(ctx context.Context, params ProblemGetParams) ([]Problem, error) {
	return getList[Problem](ctx, s.client, "problem.get", params)
}

// ProblemCountParams contains the filters for problem.get counts.
type ProblemCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	EventIDs               []string       `json:"eventids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	ObjectIDs              []string       `json:"objectids,omitzero"`
	Source                 *int           `json:"source,omitempty"`
	Object                 *int           `json:"object,omitempty"`
	Severities             []int          `json:"severities,omitzero"`
	Tags                   []TagFilter    `json:"tags,omitzero"`
	EvalType               int            `json:"evaltype,omitempty"`
	Recent                 *bool          `json:"recent,omitempty"`
	Acknowledged           *bool          `json:"acknowledged,omitempty"`
	Suppressed             *bool          `json:"suppressed,omitempty"`
	TimeFrom               *int64         `json:"time_from,omitempty"`
	TimeTill               *int64         `json:"time_till,omitempty"`
}

// Count returns the number of objects matching params.
func (s *ProblemService) Count(ctx context.Context, params ProblemCountParams) (uint64, error) {
	return getCount(ctx, s.client, "problem.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *ProblemService) GetMap(ctx context.Context, params ProblemGetParams) (map[string]Problem, error) {
	return getMap[Problem](ctx, s.client, "problem.get", params)
}
