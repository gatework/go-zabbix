package zabbix

import "context"

// HistoryService wraps the "history" API namespace.
type HistoryService struct{ client *Client }

// History returns the history service.
func (c *Client) History() *HistoryService { return &HistoryService{c} }

// HistoryValue is a single historical value. Because Zabbix stores values of
// several types (float, unsigned int, string, text, log) the Value is kept as a
// string regardless of the underlying value type. Log-specific fields are only
// populated for log history (value type 2).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/history/object
type HistoryValue struct {
	ItemID     string `json:"itemid,omitempty"`
	Clock      string `json:"clock,omitempty"`
	Value      string `json:"value,omitempty"`
	NS         string `json:"ns,omitempty"`
	LogEventID string `json:"logeventid,omitempty"`
	Severity   string `json:"severity,omitempty"`
	Source     string `json:"source,omitempty"`
	Timestamp  string `json:"timestamp,omitempty"`
}

// HistoryGetParams are the parameters for history.get. History selects the value
// type of the items being queried (0 float, 1 char, 2 log, 3 unsigned int,
// 4 text); the Zabbix server defaults to 3 when it is omitted. It is a pointer so
// that value type 0 (float) — the most common item type — can be requested
// explicitly rather than being dropped by omitempty.
// HistoryGetParams queries the list form of history.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type HistoryGetParams struct {
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
	History                *int           `json:"history,omitempty"`
	ItemIDs                []string       `json:"itemids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	TimeFrom               *int64         `json:"time_from,omitempty"`
	TimeTill               *int64         `json:"time_till,omitempty"`
}

// Get retrieves historical values matching params. The history API is read-only.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/history/get
func (s *HistoryService) Get(ctx context.Context, params HistoryGetParams) ([]HistoryValue, error) {
	return getList[HistoryValue](ctx, s.client, "history.get", params)
}

// HistoryCountParams contains the filters for history.get counts.
type HistoryCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	History                *int           `json:"history,omitempty"`
	ItemIDs                []string       `json:"itemids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	TimeFrom               *int64         `json:"time_from,omitempty"`
	TimeTill               *int64         `json:"time_till,omitempty"`
}

// Count returns the number of objects matching params.
func (s *HistoryService) Count(ctx context.Context, params HistoryCountParams) (uint64, error) {
	return getCount(ctx, s.client, "history.get", params)
}
