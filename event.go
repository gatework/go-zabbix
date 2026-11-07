package zabbix

import (
	"context"
	"errors"
)

// EventService wraps the "event" API namespace.
type EventService struct{ client *Client }

// Event returns the event service.
func (c *Client) Event() *EventService { return &EventService{c} }

// Event acknowledge action bit flags, combined into EventAcknowledgeParams.Action.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/event/acknowledge
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
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/event/object
type Event struct {
	EventID       string                `json:"eventid,omitempty"`
	Source        string                `json:"source,omitempty"`
	Object        string                `json:"object,omitempty"`
	ObjectID      string                `json:"objectid,omitempty"`
	Acknowledged  string                `json:"acknowledged,omitempty"`
	Clock         string                `json:"clock,omitempty"`
	NS            string                `json:"ns,omitempty"`
	Name          string                `json:"name,omitempty"`
	Value         string                `json:"value,omitempty"`
	Severity      string                `json:"severity,omitempty"`
	REventID      string                `json:"r_eventid,omitempty"`
	CEventID      string                `json:"c_eventid,omitempty"`
	CorrelationID string                `json:"correlationid,omitempty"`
	UserID        string                `json:"userid,omitempty"`
	Suppressed    string                `json:"suppressed,omitempty"`
	OpData        string                `json:"opdata,omitempty"`
	Tags          []Tag                 `json:"tags,omitempty"`
	Hosts         []Host                `json:"hosts,omitempty"`
	Acknowledges  Relation[EventUpdate] `json:"acknowledges,omitzero"`
}

// EventGetParams are the parameters for event.get.
// EventGetParams queries the list form of event.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type EventGetParams struct {
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
	Acknowledged           *bool          `json:"acknowledged,omitempty"`
	Suppressed             *bool          `json:"suppressed,omitempty"`
	Value                  []int          `json:"value,omitzero"`
	Tags                   []TagFilter    `json:"tags,omitzero"`
	EvalType               int            `json:"evaltype,omitempty"`
	TimeFrom               *int64         `json:"time_from,omitempty"`
	TimeTill               *int64         `json:"time_till,omitempty"`
	SelectHosts            any            `json:"selectHosts,omitempty"`
	SelectTags             any            `json:"selectTags,omitempty"`
	SelectAcknowledges     any            `json:"selectAcknowledges,omitempty"`
}

// EventUpdate is an acknowledgement, message, severity change, suppression, or
// rank change selected by event.get or problem.get using selectAcknowledges.
type EventUpdate struct {
	AcknowledgeID string `json:"acknowledgeid,omitempty"`
	UserID        string `json:"userid,omitempty"`
	EventID       string `json:"eventid,omitempty"`
	Clock         string `json:"clock,omitempty"`
	Message       string `json:"message,omitempty"`
	Action        string `json:"action,omitempty"`
	OldSeverity   string `json:"old_severity,omitempty"`
	NewSeverity   string `json:"new_severity,omitempty"`
	SuppressUntil string `json:"suppress_until,omitempty"`
	TaskID        string `json:"taskid,omitempty"`
	Username      string `json:"username,omitempty"`
	Name          string `json:"name,omitempty"`
	Surname       string `json:"surname,omitempty"`
}

// Get retrieves events matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/event/get
func (s *EventService) Get(ctx context.Context, params EventGetParams) ([]Event, error) {
	return getList[Event](ctx, s.client, "event.get", params)
}

// EventCountParams contains the filters for event.get counts.
type EventCountParams struct {
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
	Acknowledged           *bool          `json:"acknowledged,omitempty"`
	Suppressed             *bool          `json:"suppressed,omitempty"`
	Value                  []int          `json:"value,omitzero"`
	Tags                   []TagFilter    `json:"tags,omitzero"`
	EvalType               int            `json:"evaltype,omitempty"`
	TimeFrom               *int64         `json:"time_from,omitempty"`
	TimeTill               *int64         `json:"time_till,omitempty"`
}

// Count returns the number of objects matching params.
func (s *EventService) Count(ctx context.Context, params EventCountParams) (uint64, error) {
	return getCount(ctx, s.client, "event.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *EventService) GetMap(ctx context.Context, params EventGetParams) (map[string]Event, error) {
	return getMap[Event](ctx, s.client, "event.get", params)
}

// EventAcknowledgeParams are the parameters for event.acknowledge. Action is a
// bitmask of the AckAction* constants.
type EventAcknowledgeParams struct {
	EventIDs      []string `json:"eventids"`
	Action        int      `json:"action"`
	Message       *string  `json:"message,omitempty"`
	Severity      *int     `json:"severity,omitempty"`
	SuppressUntil *int64   `json:"suppress_until,omitempty"`
	CauseEventID  *string  `json:"cause_eventid,omitempty"`
}

func (p EventAcknowledgeParams) validate() error {
	if len(p.EventIDs) == 0 {
		return errors.New("zabbix: event.acknowledge requires event IDs")
	}
	for _, id := range p.EventIDs {
		if err := requireWriteID(id, "EventAcknowledgeParams.EventIDs"); err != nil {
			return err
		}
	}
	const mask = AckActionClose | AckActionAcknowledge | AckActionAddMessage |
		AckActionChangeSeverity | AckActionUnacknowledge | AckActionSuppress |
		AckActionUnsuppress | AckActionChangeToCause | AckActionChangeToSymptom
	if p.Action < 0 || p.Action&^mask != 0 {
		return errors.New("zabbix: event.acknowledge action contains unsupported bits")
	}
	if p.Action&AckActionAcknowledge != 0 && p.Action&AckActionUnacknowledge != 0 {
		return errors.New("zabbix: cannot acknowledge and unacknowledge the same event")
	}
	if p.Action&AckActionSuppress != 0 && p.Action&AckActionUnsuppress != 0 {
		return errors.New("zabbix: cannot suppress and unsuppress the same event")
	}
	if p.Action&AckActionClose != 0 && p.Action&(AckActionSuppress|AckActionUnsuppress) != 0 {
		return errors.New("zabbix: cannot close an event while changing its suppression")
	}
	if p.Action&AckActionChangeToCause != 0 && p.Action&AckActionChangeToSymptom != 0 {
		return errors.New("zabbix: cannot change the same event to both cause and symptom")
	}
	if p.Action&AckActionAddMessage != 0 && (p.Message == nil || *p.Message == "") {
		return errors.New("zabbix: adding an event message requires a nonempty message")
	}
	if p.Action&AckActionChangeSeverity != 0 && p.Severity == nil {
		return errors.New("zabbix: changing event severity requires severity")
	}
	if p.Severity != nil && (*p.Severity < 0 || *p.Severity > 5) {
		return errors.New("zabbix: event severity must be between 0 and 5")
	}
	if p.Action&AckActionSuppress != 0 && p.SuppressUntil == nil {
		return errors.New("zabbix: suppressing an event requires suppress_until")
	}
	if p.SuppressUntil != nil && *p.SuppressUntil < 0 {
		return errors.New("zabbix: suppress_until must not be negative")
	}
	if p.Action&AckActionChangeToSymptom != 0 {
		if p.CauseEventID == nil {
			return errors.New("zabbix: changing an event to a symptom requires cause_eventid")
		}
		if err := requireWriteID(*p.CauseEventID, "EventAcknowledgeParams.CauseEventID"); err != nil {
			return err
		}
	} else if p.CauseEventID != nil {
		return errors.New("zabbix: cause_eventid is only accepted when changing an event to a symptom")
	}
	return nil
}

// Acknowledge updates problem events (acknowledge, close, add message, change
// severity, suppress, ...) and returns the affected event ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/event/acknowledge
func (s *EventService) Acknowledge(ctx context.Context, params EventAcknowledgeParams) ([]string, error) {
	return mutateObjectIDs(ctx, s.client, "event.acknowledge", "eventids", params, false)
}
