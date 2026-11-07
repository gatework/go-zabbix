package zabbix

import "context"

// HostGroupService wraps the "hostgroup" API namespace.
type HostGroupService struct{ client *Client }

// HostGroup returns the host group service.
func (c *Client) HostGroup() *HostGroupService { return &HostGroupService{c} }

// HostGroup object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostgroup/object
type HostGroup struct {
	GroupID string         `json:"groupid,omitempty"`
	Name    string         `json:"name,omitempty"`
	Flags   string         `json:"flags,omitempty"`
	UUID    string         `json:"uuid,omitempty"`
	Hosts   Relation[Host] `json:"hosts,omitzero"`
}

// HostGroupGetParams are the parameters for hostgroup.get.
// HostGroupGetParams queries the list form of hostgroup.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type HostGroupGetParams struct {
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
	GroupIDs               []string       `json:"groupids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	GraphIDs               []string       `json:"graphids,omitzero"`
	TriggerIDs             []string       `json:"triggerids,omitzero"`
	MaintenanceIDs         []string       `json:"maintenanceids,omitzero"`
	WithHosts              bool           `json:"with_hosts,omitempty"`
	WithMonitoredHosts     bool           `json:"with_monitored_hosts,omitempty"`
	WithItems              bool           `json:"with_items,omitempty"`
	WithTriggers           bool           `json:"with_triggers,omitempty"`
	SelectHosts            any            `json:"selectHosts,omitempty"`
}

// Get retrieves host groups matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostgroup/get
func (s *HostGroupService) Get(ctx context.Context, params HostGroupGetParams) ([]HostGroup, error) {
	return getList[HostGroup](ctx, s.client, "hostgroup.get", params)
}

// HostGroupCountParams contains the filters for hostgroup.get counts.
type HostGroupCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	GraphIDs               []string       `json:"graphids,omitzero"`
	TriggerIDs             []string       `json:"triggerids,omitzero"`
	MaintenanceIDs         []string       `json:"maintenanceids,omitzero"`
	WithHosts              bool           `json:"with_hosts,omitempty"`
	WithMonitoredHosts     bool           `json:"with_monitored_hosts,omitempty"`
	WithItems              bool           `json:"with_items,omitempty"`
	WithTriggers           bool           `json:"with_triggers,omitempty"`
}

// Count returns the number of objects matching params.
func (s *HostGroupService) Count(ctx context.Context, params HostGroupCountParams) (uint64, error) {
	return getCount(ctx, s.client, "hostgroup.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *HostGroupService) GetMap(ctx context.Context, params HostGroupGetParams) (map[string]HostGroup, error) {
	return getMap[HostGroup](ctx, s.client, "hostgroup.get", params)
}

// Create creates host groups and returns the new group ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostgroup/create
func (s *HostGroupService) Create(ctx context.Context, groups ...HostGroupCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "hostgroup.create", "groupids", groups)
}

// Update updates host groups and returns the affected group ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostgroup/update
func (s *HostGroupService) Update(ctx context.Context, groups ...HostGroupUpdateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "hostgroup.update", "groupids", groups)
}

// Delete deletes host groups by id and returns the deleted group ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostgroup/delete
func (s *HostGroupService) Delete(ctx context.Context, groupIDs ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "hostgroup.delete", "groupids", groupIDs)
}
