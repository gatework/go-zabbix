package zabbix

import "context"

// EventService wraps the "event" API namespace.
type EventService struct{ client *Client }

// Event returns the event service.
func (c *Client) Event() *EventService { return &EventService{c} }

// Event acknowledge action bit flags, combined into EventAcknowledgeParams.Action.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/event/acknowledge
const (
	AckActionClose           = 1 << 0 // close problem
	AckActionAcknowledge     = 1 << 1 // acknowledge event
	AckActionAddMessage      = 1 << 2 // add message
	AckActionChangeSeverity  = 1 << 3 // change severity
	AckActionUnacknowledge   = 1 << 4 // unacknowledge event
	AckActionSuppress        = 1 << 5 // suppress event
	AckActionUnsuppress      = 1 << 6 // unsuppress event
	AckActionChangeToCause   = 1 << 7 // change event rank to cause
	AckActionChangeToSymptom = 1 << 8 // change event rank to symptom
)

// Event object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/event/object
type Event struct {
	EventID       string `json:"eventid,omitempty"`
	Source        string `json:"source,omitempty"`
	Object        string `json:"object,omitempty"`
	ObjectID      string `json:"objectid,omitempty"`
	Acknowledged  string `json:"acknowledged,omitempty"`
	Clock         string `json:"clock,omitempty"`
	Ns            string `json:"ns,omitempty"`
	Name          string `json:"name,omitempty"`
	Value         string `json:"value,omitempty"`
	Severity      string `json:"severity,omitempty"`
	REventID      string `json:"r_eventid,omitempty"`
	CEventID      string `json:"c_eventid,omitempty"`
	CorrelationID string `json:"correlationid,omitempty"`
	UserID        string `json:"userid,omitempty"`
	Suppressed    string `json:"suppressed,omitempty"`
	OpData        string `json:"opdata,omitempty"`
	Tags          []Tag  `json:"tags,omitempty"`
}

// EventGetParams are the parameters for event.get.
type EventGetParams struct {
	GetParams
	EventIDs           []string `json:"eventids,omitempty"`
	GroupIDs           []string `json:"groupids,omitempty"`
	HostIDs            []string `json:"hostids,omitempty"`
	ObjectIDs          []string `json:"objectids,omitempty"`
	Source             *int     `json:"source,omitempty"`
	Object             *int     `json:"object,omitempty"`
	Severities         []int    `json:"severities,omitempty"`
	Acknowledged       *bool    `json:"acknowledged,omitempty"`
	Suppressed         *bool    `json:"suppressed,omitempty"`
	Value              []int    `json:"value,omitempty"`
	Tags               []Tag    `json:"tags,omitempty"`
	EvalType           int      `json:"evaltype,omitempty"`
	TimeFrom           int64    `json:"time_from,omitempty"`
	TimeTill           int64    `json:"time_till,omitempty"`
	SelectHosts        any      `json:"selectHosts,omitempty"`
	SelectTags         any      `json:"selectTags,omitempty"`
	SelectAcknowledges any      `json:"selectAcknowledges,omitempty"`
}

// Get retrieves events matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/event/get
func (s *EventService) Get(ctx context.Context, params EventGetParams) ([]Event, error) {
	var events []Event
	err := s.client.Call(ctx, "event.get", params, &events)
	return events, err
}

// EventAcknowledgeParams are the parameters for event.acknowledge. Action is a
// bitmask of the AckAction* constants.
type EventAcknowledgeParams struct {
	EventIDs []string `json:"eventids"`
	Action   int      `json:"action"`
	Message  string   `json:"message,omitempty"`
	Severity int      `json:"severity,omitempty"`
}

// Acknowledge updates problem events (acknowledge, close, add message, change
// severity, suppress, ...) and returns the affected event ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/event/acknowledge
func (s *EventService) Acknowledge(ctx context.Context, params EventAcknowledgeParams) ([]string, error) {
	var res struct {
		EventIDs []string `json:"eventids"`
	}
	err := s.client.Call(ctx, "event.acknowledge", params, &res)
	return res.EventIDs, err
}
