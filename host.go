package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
)

// HostService wraps the "host" API namespace.
type HostService struct{ client *Client }

// Host returns the host service.
func (c *Client) Host() *HostService { return &HostService{c} }

// Host object (Zabbix 7.x).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/host/object
type Host struct {
	HostID            string `json:"hostid,omitempty"`
	Host              string `json:"host,omitempty"`
	Name              string `json:"name,omitempty"`
	Description       string `json:"description,omitempty"`
	Status            string `json:"status,omitempty"`
	Flags             string `json:"flags,omitempty"`
	InventoryMode     string `json:"inventory_mode,omitempty"`
	MonitoredBy       string `json:"monitored_by,omitempty"` // 0 server, 1 proxy, 2 proxy group (7.0+)
	ProxyID           string `json:"proxyid,omitempty"`
	ProxyGroupID      string `json:"proxy_groupid,omitempty"`
	MaintenanceStatus string `json:"maintenance_status,omitempty"`
	MaintenanceID     string `json:"maintenanceid,omitempty"`
	MaintenanceType   string `json:"maintenance_type,omitempty"`
	MaintenanceFrom   string `json:"maintenance_from,omitempty"`
	ActiveAvailable   string `json:"active_available,omitempty"`
	TLSConnect        string `json:"tls_connect,omitempty"`
	TLSAccept         string `json:"tls_accept,omitempty"`
	TLSIssuer         string `json:"tls_issuer,omitempty"`
	TLSSubject        string `json:"tls_subject,omitempty"`
	IPMIAuthType      string `json:"ipmi_authtype,omitempty"`
	IPMIPrivilege     string `json:"ipmi_privilege,omitempty"`
	IPMIUsername      string `json:"ipmi_username,omitempty"`
	IPMIPassword      string `json:"ipmi_password,omitempty"`

	// Interfaces, Tags and Macros use the same key on input and output.
	Interfaces []HostInterface `json:"interfaces,omitempty"`
	Tags       []Tag           `json:"tags,omitempty"`
	Macros     []UserMacro     `json:"macros,omitempty"`

	// Groups/Templates are the create/update INPUT keys.
	Groups    []HostGroup `json:"groups,omitempty"`
	Templates []Template  `json:"templates,omitempty"`

	// HostGroups/ParentTemplates are the host.get OUTPUT keys (Zabbix 6.2+ renamed
	// them from "groups"/"templates"). Populated by selectHostGroups /
	// selectParentTemplates.
	HostGroups      []HostGroup `json:"hostgroups,omitempty"`
	ParentTemplates []Template  `json:"parentTemplates,omitempty"`
}

// HostInterface object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostinterface/object
type HostInterface struct {
	InterfaceID string           `json:"interfaceid,omitempty"`
	HostID      string           `json:"hostid,omitempty"`
	Main        string           `json:"main,omitempty"`
	Type        string           `json:"type,omitempty"`
	UseIP       string           `json:"useip,omitempty"`
	IP          string           `json:"ip,omitempty"`
	DNS         string           `json:"dns,omitempty"`
	Port        string           `json:"port,omitempty"`
	Available   string           `json:"available,omitempty"`
	Details     *InterfaceDetail `json:"details,omitempty"`
}

// UnmarshalJSON handles Zabbix returning the "details" field as an empty array
// ([]) for non-SNMP interfaces and as an object for SNMP interfaces. In the
// array case Details is left nil.
func (h *HostInterface) UnmarshalJSON(data []byte) error {
	type alias HostInterface
	aux := struct {
		Details json.RawMessage `json:"details"`
		*alias
	}{alias: (*alias)(h)}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if raw := bytes.TrimSpace(aux.Details); len(raw) > 0 && raw[0] == '{' {
		var d InterfaceDetail
		if err := json.Unmarshal(raw, &d); err != nil {
			return err
		}
		h.Details = &d
	} else {
		h.Details = nil
	}
	return nil
}

// InterfaceDetail holds SNMP-specific interface settings.
type InterfaceDetail struct {
	Version        string `json:"version,omitempty"`
	Bulk           string `json:"bulk,omitempty"`
	Community      string `json:"community,omitempty"`
	SecurityName   string `json:"securityname,omitempty"`
	SecurityLevel  string `json:"securitylevel,omitempty"`
	AuthPassphrase string `json:"authpassphrase,omitempty"`
	PrivPassphrase string `json:"privpassphrase,omitempty"`
	AuthProtocol   string `json:"authprotocol,omitempty"`
	PrivProtocol   string `json:"privprotocol,omitempty"`
	ContextName    string `json:"contextname,omitempty"`
	MaxRepetitions string `json:"max_repetitions,omitempty"`
}

// UserMacro is a user-defined macro on a host or template.
type UserMacro struct {
	HostMacroID string `json:"hostmacroid,omitempty"`
	Macro       string `json:"macro,omitempty"`
	Value       string `json:"value,omitempty"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
}

// HostGetParams are the parameters for host.get.
type HostGetParams struct {
	GetParams
	HostIDs          []string `json:"hostids,omitempty"`
	GroupIDs         []string `json:"groupids,omitempty"`
	TemplateIDs      []string `json:"templateids,omitempty"`
	ItemIDs          []string `json:"itemids,omitempty"`
	TriggerIDs       []string `json:"triggerids,omitempty"`
	MaintenanceIDs   []string `json:"maintenanceids,omitempty"`
	ProxyIDs         []string `json:"proxyids,omitempty"`
	Tags             []Tag    `json:"tags,omitempty"`
	EvalType         int      `json:"evaltype,omitempty"`
	WithItems        bool     `json:"with_items,omitempty"`
	WithTriggers     bool     `json:"with_triggers,omitempty"`
	MonitoredHosts   bool     `json:"monitored_hosts,omitempty"`
	SelectGroups     any      `json:"selectHostGroups,omitempty"`
	SelectInterfaces any      `json:"selectInterfaces,omitempty"`
	SelectItems      any      `json:"selectItems,omitempty"`
	SelectTriggers   any      `json:"selectTriggers,omitempty"`
	SelectTemplates  any      `json:"selectParentTemplates,omitempty"`
	SelectTags       any      `json:"selectTags,omitempty"`
	SelectMacros     any      `json:"selectMacros,omitempty"`
	SelectInventory  any      `json:"selectInventory,omitempty"`
}

// Get retrieves hosts matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/host/get
func (s *HostService) Get(ctx context.Context, params HostGetParams) ([]Host, error) {
	var hosts []Host
	err := s.client.Call(ctx, "host.get", params, &hosts)
	return hosts, err
}

// Create creates hosts and returns the new host ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/host/create
func (s *HostService) Create(ctx context.Context, hosts ...Host) ([]string, error) {
	var res struct {
		HostIDs []string `json:"hostids"`
	}
	err := s.client.Call(ctx, "host.create", hosts, &res)
	return res.HostIDs, err
}

// Update updates hosts and returns the affected host ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/host/update
func (s *HostService) Update(ctx context.Context, hosts ...Host) ([]string, error) {
	var res struct {
		HostIDs []string `json:"hostids"`
	}
	err := s.client.Call(ctx, "host.update", hosts, &res)
	return res.HostIDs, err
}

// Delete deletes hosts by id and returns the deleted host ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/host/delete
func (s *HostService) Delete(ctx context.Context, hostIDs ...string) ([]string, error) {
	var res struct {
		HostIDs []string `json:"hostids"`
	}
	err := s.client.Call(ctx, "host.delete", hostIDs, &res)
	return res.HostIDs, err
}
