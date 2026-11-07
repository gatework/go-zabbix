package zabbix

import "context"

// HostGroupService wraps the "hostgroup" API namespace.
type HostGroupService struct{ client *Client }

// HostGroup returns the host group service.
func (c *Client) HostGroup() *HostGroupService { return &HostGroupService{c} }

// HostGroup object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostgroup/object
type HostGroup struct {
	GroupID string `json:"groupid,omitempty"`
	Name    string `json:"name,omitempty"`
	Flags   string `json:"flags,omitempty"`
	UUID    string `json:"uuid,omitempty"`
}

// HostGroupGetParams are the parameters for hostgroup.get.
type HostGroupGetParams struct {
	GetParams
	GroupIDs           []string `json:"groupids,omitempty"`
	HostIDs            []string `json:"hostids,omitempty"`
	GraphIDs           []string `json:"graphids,omitempty"`
	TriggerIDs         []string `json:"triggerids,omitempty"`
	MaintenanceIDs     []string `json:"maintenanceids,omitempty"`
	WithHosts          bool     `json:"with_hosts,omitempty"`
	WithMonitoredHosts bool     `json:"with_monitored_hosts,omitempty"`
	WithItems          bool     `json:"with_items,omitempty"`
	WithTriggers       bool     `json:"with_triggers,omitempty"`
	RealHosts          bool     `json:"real_hosts,omitempty"`
	SelectHosts        any      `json:"selectHosts,omitempty"`
}

// Get retrieves host groups matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostgroup/get
func (s *HostGroupService) Get(ctx context.Context, params HostGroupGetParams) ([]HostGroup, error) {
	var groups []HostGroup
	err := s.client.Call(ctx, "hostgroup.get", params, &groups)
	return groups, err
}

// Create creates host groups and returns the new group ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostgroup/create
func (s *HostGroupService) Create(ctx context.Context, groups ...HostGroup) ([]string, error) {
	var res struct {
		GroupIDs []string `json:"groupids"`
	}
	err := s.client.Call(ctx, "hostgroup.create", groups, &res)
	return res.GroupIDs, err
}

// Update updates host groups and returns the affected group ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostgroup/update
func (s *HostGroupService) Update(ctx context.Context, groups ...HostGroup) ([]string, error) {
	var res struct {
		GroupIDs []string `json:"groupids"`
	}
	err := s.client.Call(ctx, "hostgroup.update", groups, &res)
	return res.GroupIDs, err
}

// Delete deletes host groups by id and returns the deleted group ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostgroup/delete
func (s *HostGroupService) Delete(ctx context.Context, groupIDs ...string) ([]string, error) {
	var res struct {
		GroupIDs []string `json:"groupids"`
	}
	err := s.client.Call(ctx, "hostgroup.delete", groupIDs, &res)
	return res.GroupIDs, err
}
