package zabbix

import "context"

// AlertService wraps the "alert" API namespace.
type AlertService struct{ client *Client }

// Alert returns the alert service.
func (c *Client) Alert() *AlertService { return &AlertService{c} }

// Alert object. The alert API is read-only.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/alert/object
type Alert struct {
	AlertID       string      `json:"alertid,omitempty"`
	ActionID      string      `json:"actionid,omitempty"`
	EventID       string      `json:"eventid,omitempty"`
	UserID        string      `json:"userid,omitempty"`
	Clock         string      `json:"clock,omitempty"`
	MediaTypeID   string      `json:"mediatypeid,omitempty"`
	SendTo        string      `json:"sendto,omitempty"`
	Subject       string      `json:"subject,omitempty"`
	Message       string      `json:"message,omitempty"`
	Status        string      `json:"status,omitempty"`
	Retries       string      `json:"retries,omitempty"`
	Error         string      `json:"error,omitempty"`
	EscStep       string      `json:"esc_step,omitempty"`
	AlertType     string      `json:"alerttype,omitempty"`
	PEventID      string      `json:"p_eventid,omitempty"`
	AcknowledgeID string      `json:"acknowledgeid,omitempty"`
	Hosts         []Host      `json:"hosts,omitempty"`
	MediaTypes    []MediaType `json:"mediatypes,omitempty"`
	Users         []User      `json:"users,omitempty"`
}

// AlertGetParams are the parameters for alert.get.
// AlertGetParams queries the list form of alert.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type AlertGetParams struct {
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
	AlertIDs               []string       `json:"alertids,omitzero"`
	ActionIDs              []string       `json:"actionids,omitzero"`
	EventIDs               []string       `json:"eventids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	MediaTypeIDs           []string       `json:"mediatypeids,omitzero"`
	UserIDs                []string       `json:"userids,omitzero"`
	TimeFrom               *int64         `json:"time_from,omitempty"`
	TimeTill               *int64         `json:"time_till,omitempty"`
	SelectHosts            any            `json:"selectHosts,omitempty"`
	SelectMediatypes       any            `json:"selectMediatypes,omitempty"`
	SelectUsers            any            `json:"selectUsers,omitempty"`
}

// Get retrieves alerts matching params. Note the alert API is read-only.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/alert/get
func (s *AlertService) Get(ctx context.Context, params AlertGetParams) ([]Alert, error) {
	return getList[Alert](ctx, s.client, "alert.get", params)
}

// AlertCountParams contains the filters for alert.get counts.
type AlertCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	AlertIDs               []string       `json:"alertids,omitzero"`
	ActionIDs              []string       `json:"actionids,omitzero"`
	EventIDs               []string       `json:"eventids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	MediaTypeIDs           []string       `json:"mediatypeids,omitzero"`
	UserIDs                []string       `json:"userids,omitzero"`
	TimeFrom               *int64         `json:"time_from,omitempty"`
	TimeTill               *int64         `json:"time_till,omitempty"`
}

// Count returns the number of objects matching params.
func (s *AlertService) Count(ctx context.Context, params AlertCountParams) (uint64, error) {
	return getCount(ctx, s.client, "alert.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *AlertService) GetMap(ctx context.Context, params AlertGetParams) (map[string]Alert, error) {
	return getMap[Alert](ctx, s.client, "alert.get", params)
}
