package zabbix

import "context"

// AlertService wraps the "alert" API namespace.
type AlertService struct{ client *Client }

// Alert returns the alert service.
func (c *Client) Alert() *AlertService { return &AlertService{c} }

// Alert object. The alert API is read-only.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/alert/object
type Alert struct {
	AlertID       string `json:"alertid,omitempty"`
	ActionID      string `json:"actionid,omitempty"`
	EventID       string `json:"eventid,omitempty"`
	UserID        string `json:"userid,omitempty"`
	Clock         string `json:"clock,omitempty"`
	MediaTypeID   string `json:"mediatypeid,omitempty"`
	SendTo        string `json:"sendto,omitempty"`
	Subject       string `json:"subject,omitempty"`
	Message       string `json:"message,omitempty"`
	Status        string `json:"status,omitempty"`
	Retries       string `json:"retries,omitempty"`
	Error         string `json:"error,omitempty"`
	EscStep       string `json:"esc_step,omitempty"`
	AlertType     string `json:"alerttype,omitempty"`
	PEventID      string `json:"p_eventid,omitempty"`
	AcknowledgeID string `json:"acknowledgeid,omitempty"`
}

// AlertGetParams are the parameters for alert.get.
type AlertGetParams struct {
	GetParams
	AlertIDs         []string `json:"alertids,omitempty"`
	ActionIDs        []string `json:"actionids,omitempty"`
	EventIDs         []string `json:"eventids,omitempty"`
	GroupIDs         []string `json:"groupids,omitempty"`
	HostIDs          []string `json:"hostids,omitempty"`
	MediaTypeIDs     []string `json:"mediatypeids,omitempty"`
	UserIDs          []string `json:"userids,omitempty"`
	TimeFrom         int64    `json:"time_from,omitempty"`
	TimeTill         int64    `json:"time_till,omitempty"`
	SelectHosts      any      `json:"selectHosts,omitempty"`
	SelectMediatypes any      `json:"selectMediatypes,omitempty"`
	SelectUsers      any      `json:"selectUsers,omitempty"`
}

// Get retrieves alerts matching params. Note the alert API is read-only.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/alert/get
func (s *AlertService) Get(ctx context.Context, params AlertGetParams) ([]Alert, error) {
	var out []Alert
	err := s.client.Call(ctx, "alert.get", params, &out)
	return out, err
}
