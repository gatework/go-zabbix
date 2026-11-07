package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
)

// HostService wraps the "host" API namespace.
type HostService struct{ client *Client }

// Host returns the host service.
func (c *Client) Host() *HostService { return &HostService{c} }

// Host object (Zabbix 7.x).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/host/object
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

	// Related objects are populated by the corresponding selectors.
	Interfaces Relation[HostInterface] `json:"interfaces,omitzero"`
	Tags       []Tag                   `json:"tags,omitempty"`
	Macros     []UserMacro             `json:"macros,omitempty"`

	// Populated by selectHostGroups and selectParentTemplates.
	HostGroups      []HostGroup        `json:"hostgroups,omitempty"`
	ParentTemplates Relation[Template] `json:"parentTemplates,omitzero"`
	Items           Relation[Item]     `json:"items,omitzero"`
	Triggers        Relation[Trigger]  `json:"triggers,omitzero"`
	Inventory       HostInventory      `json:"inventory,omitzero"`
}

// HostInterface object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/object
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
	Items       Relation[Item]   `json:"items,omitzero"`
	Hosts       []Host           `json:"hosts,omitempty"`
}

// UnmarshalJSON handles Zabbix returning the "details" field as an empty array
// ([]) for non-SNMP interfaces and as an object for SNMP interfaces. In the
// array case Details is left nil.
func (h *HostInterface) UnmarshalJSON(data []byte) error {
	type alias HostInterface
	var next alias
	aux := struct {
		Details json.RawMessage `json:"details"`
		*alias
	}{alias: &next}
	if raw := bytes.TrimSpace(data); len(raw) == 0 || raw[0] != '{' {
		return errors.New("host interface must be an object")
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	raw := bytes.TrimSpace(aux.Details)
	switch {
	case len(raw) == 0 || isEmptyArray(raw):
	case raw[0] == '{':
		var d InterfaceDetail
		if err := json.Unmarshal(raw, &d); err != nil {
			return err
		}
		next.Details = &d
	default:
		return errors.New("interface details must be an object or an empty array")
	}
	*h = HostInterface(next)
	return nil
}

// HostInventory contains the fields requested with selectInventory. A nil map
// means inventory was not selected; an empty non-nil map is an empty result.
type HostInventory map[string]string

func (h *HostInventory) UnmarshalJSON(data []byte) error {
	raw := bytes.TrimSpace(data)
	next := make(HostInventory)
	if isEmptyArray(raw) {
		*h = next
		return nil
	}
	if len(raw) == 0 || raw[0] != '{' {
		return errors.New("host inventory must be an object or an empty array")
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	for key, value := range fields {
		if raw := bytes.TrimSpace(value); len(raw) == 0 || raw[0] != '"' {
			return errors.New("host inventory values must be strings")
		}
		var text string
		if err := json.Unmarshal(value, &text); err != nil {
			return err
		}
		next[key] = text
	}
	*h = next
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
	HostID      string `json:"hostid,omitempty"`
}

// HostGetParams are the parameters for host.get.
// HostGetParams queries the list form of host.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type HostGetParams struct {
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
	HostIDs                []string       `json:"hostids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	TemplateIDs            []string       `json:"templateids,omitzero"`
	ItemIDs                []string       `json:"itemids,omitzero"`
	TriggerIDs             []string       `json:"triggerids,omitzero"`
	MaintenanceIDs         []string       `json:"maintenanceids,omitzero"`
	ProxyIDs               []string       `json:"proxyids,omitzero"`
	Tags                   []TagFilter    `json:"tags,omitzero"`
	EvalType               int            `json:"evaltype,omitempty"`
	WithItems              bool           `json:"with_items,omitempty"`
	WithTriggers           bool           `json:"with_triggers,omitempty"`
	MonitoredHosts         bool           `json:"monitored_hosts,omitempty"`
	SelectHostGroups       any            `json:"selectHostGroups,omitempty"`
	SelectInterfaces       any            `json:"selectInterfaces,omitempty"`
	SelectItems            any            `json:"selectItems,omitempty"`
	SelectTriggers         any            `json:"selectTriggers,omitempty"`
	SelectParentTemplates  any            `json:"selectParentTemplates,omitempty"`
	SelectTags             any            `json:"selectTags,omitempty"`
	SelectMacros           any            `json:"selectMacros,omitempty"`
	SelectInventory        any            `json:"selectInventory,omitempty"`
}

// Get retrieves hosts matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/host/get
func (s *HostService) Get(ctx context.Context, params HostGetParams) ([]Host, error) {
	return getList[Host](ctx, s.client, "host.get", params)
}

// HostCountParams contains the filters for host.get counts.
type HostCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	HostIDs                []string       `json:"hostids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	TemplateIDs            []string       `json:"templateids,omitzero"`
	ItemIDs                []string       `json:"itemids,omitzero"`
	TriggerIDs             []string       `json:"triggerids,omitzero"`
	MaintenanceIDs         []string       `json:"maintenanceids,omitzero"`
	ProxyIDs               []string       `json:"proxyids,omitzero"`
	Tags                   []TagFilter    `json:"tags,omitzero"`
	EvalType               int            `json:"evaltype,omitempty"`
	WithItems              bool           `json:"with_items,omitempty"`
	WithTriggers           bool           `json:"with_triggers,omitempty"`
	MonitoredHosts         bool           `json:"monitored_hosts,omitempty"`
}

// Count returns the number of objects matching params.
func (s *HostService) Count(ctx context.Context, params HostCountParams) (uint64, error) {
	return getCount(ctx, s.client, "host.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *HostService) GetMap(ctx context.Context, params HostGetParams) (map[string]Host, error) {
	return getMap[Host](ctx, s.client, "host.get", params)
}

// Create creates hosts and returns the new host ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/host/create
func (s *HostService) Create(ctx context.Context, hosts ...Host) ([]string, error) {
	var res struct {
		HostIDs []string `json:"hostids"`
	}
	err := s.client.Call(ctx, "host.create", hosts, &res)
	return res.HostIDs, err
}

// Update updates hosts and returns the affected host ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/host/update
func (s *HostService) Update(ctx context.Context, hosts ...Host) ([]string, error) {
	var res struct {
		HostIDs []string `json:"hostids"`
	}
	err := s.client.Call(ctx, "host.update", hosts, &res)
	return res.HostIDs, err
}

// Delete deletes hosts by id and returns the deleted host ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/host/delete
func (s *HostService) Delete(ctx context.Context, hostIDs ...string) ([]string, error) {
	var res struct {
		HostIDs []string `json:"hostids"`
	}
	err := s.client.Call(ctx, "host.delete", hostIDs, &res)
	return res.HostIDs, err
}
