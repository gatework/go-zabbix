package zabbix

import "context"

// MaintenanceService wraps the "maintenance" API namespace.
type MaintenanceService struct{ client *Client }

// Maintenance returns the maintenance service.
func (c *Client) Maintenance() *MaintenanceService { return &MaintenanceService{c} }

// Maintenance object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/maintenance/object
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

	// HostGroups and Hosts are selected by selectHostGroups and selectHosts.

	Hosts      []Host      `json:"hosts,omitempty"`
	HostGroups []HostGroup `json:"hostgroups,omitempty"`
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
// MaintenanceGetParams queries the list form of maintenance.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type MaintenanceGetParams struct {
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
	MaintenanceIDs         []string       `json:"maintenanceids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	SelectTimeperiods      any            `json:"selectTimeperiods,omitempty"`
	SelectTags             any            `json:"selectTags,omitempty"`
	SelectHostGroups       any            `json:"selectHostGroups,omitempty"`
	SelectHosts            any            `json:"selectHosts,omitempty"`
}

// Get retrieves maintenances matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/maintenance/get
func (s *MaintenanceService) Get(ctx context.Context, params MaintenanceGetParams) ([]Maintenance, error) {
	return getList[Maintenance](ctx, s.client, "maintenance.get", params)
}

// MaintenanceCountParams contains the filters for maintenance.get counts.
type MaintenanceCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	MaintenanceIDs         []string       `json:"maintenanceids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *MaintenanceService) Count(ctx context.Context, params MaintenanceCountParams) (uint64, error) {
	return getCount(ctx, s.client, "maintenance.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *MaintenanceService) GetMap(ctx context.Context, params MaintenanceGetParams) (map[string]Maintenance, error) {
	return getMap[Maintenance](ctx, s.client, "maintenance.get", params)
}

// Create creates maintenances and returns the new maintenance ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/maintenance/create
func (s *MaintenanceService) Create(ctx context.Context, items ...MaintenanceCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "maintenance.create", "maintenanceids", items)
}

// Update updates maintenances and returns the affected maintenance ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/maintenance/update
func (s *MaintenanceService) Update(ctx context.Context, items ...MaintenanceUpdateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "maintenance.update", "maintenanceids", items)
}

// Delete deletes maintenances by id and returns the deleted maintenance ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/maintenance/delete
func (s *MaintenanceService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "maintenance.delete", "maintenanceids", ids)
}
