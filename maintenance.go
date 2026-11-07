package zabbix

import "context"

// MaintenanceService wraps the "maintenance" API namespace.
type MaintenanceService struct{ client *Client }

// Maintenance returns the maintenance service.
func (c *Client) Maintenance() *MaintenanceService { return &MaintenanceService{c} }

// Maintenance object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/maintenance/object
type Maintenance struct {
	MaintenanceID   string `json:"maintenanceid,omitempty"`
	Name            string `json:"name,omitempty"`
	ActiveSince     string `json:"active_since,omitempty"`
	ActiveTill      string `json:"active_till,omitempty"`
	Description     string `json:"description,omitempty"`
	MaintenanceType string `json:"maintenance_type,omitempty"`
	TagsEvalType    string `json:"tags_evaltype,omitempty"`

	Timeperiods []MaintenanceTimeperiod `json:"timeperiods,omitempty"`
	Tags        []MaintenanceTag        `json:"tags,omitempty"`

	// Groups and Hosts use the "groups"/"hosts" keys on both input and output
	// in Zabbix 7.0+ (input elements are {groupid}/{hostid}).
	Groups []HostGroup `json:"groups,omitempty"`
	Hosts  []Host      `json:"hosts,omitempty"`
}

// MaintenanceTimeperiod is a maintenance time period.
type MaintenanceTimeperiod struct {
	TimeperiodID   string `json:"timeperiodid,omitempty"`
	TimeperiodType string `json:"timeperiod_type,omitempty"`
	Every          string `json:"every,omitempty"`
	Month          string `json:"month,omitempty"`
	DayOfWeek      string `json:"dayofweek,omitempty"`
	Day            string `json:"day,omitempty"`
	StartTime      string `json:"start_time,omitempty"`
	Period         string `json:"period,omitempty"`
	StartDate      string `json:"start_date,omitempty"`
}

// MaintenanceTag is a problem tag filter of a maintenance.
type MaintenanceTag struct {
	Tag      string `json:"tag,omitempty"`
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value,omitempty"`
}

// MaintenanceGetParams are the parameters for maintenance.get.
type MaintenanceGetParams struct {
	GetParams
	MaintenanceIDs    []string `json:"maintenanceids,omitempty"`
	GroupIDs          []string `json:"groupids,omitempty"`
	HostIDs           []string `json:"hostids,omitempty"`
	SelectTimeperiods any      `json:"selectTimeperiods,omitempty"`
	SelectTags        any      `json:"selectTags,omitempty"`
	SelectHostGroups  any      `json:"selectHostGroups,omitempty"`
	SelectHosts       any      `json:"selectHosts,omitempty"`
}

// Get retrieves maintenances matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/maintenance/get
func (s *MaintenanceService) Get(ctx context.Context, params MaintenanceGetParams) ([]Maintenance, error) {
	var out []Maintenance
	err := s.client.Call(ctx, "maintenance.get", params, &out)
	return out, err
}

// Create creates maintenances and returns the new maintenance ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/maintenance/create
func (s *MaintenanceService) Create(ctx context.Context, items ...Maintenance) ([]string, error) {
	var res struct {
		IDs []string `json:"maintenanceids"`
	}
	err := s.client.Call(ctx, "maintenance.create", items, &res)
	return res.IDs, err
}

// Update updates maintenances and returns the affected maintenance ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/maintenance/update
func (s *MaintenanceService) Update(ctx context.Context, items ...Maintenance) ([]string, error) {
	var res struct {
		IDs []string `json:"maintenanceids"`
	}
	err := s.client.Call(ctx, "maintenance.update", items, &res)
	return res.IDs, err
}

// Delete deletes maintenances by id and returns the deleted maintenance ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/maintenance/delete
func (s *MaintenanceService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"maintenanceids"`
	}
	err := s.client.Call(ctx, "maintenance.delete", ids, &res)
	return res.IDs, err
}
