package zabbix

import "context"

// HostPrototypeService wraps the "hostprototype" API namespace.
type HostPrototypeService struct{ client *Client }

// HostPrototype returns the hostprototype service.
func (c *Client) HostPrototype() *HostPrototypeService { return &HostPrototypeService{c} }

// HostPrototype object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostprototype/object
type HostPrototype struct {
	HostID           string `json:"hostid,omitempty"`
	Host             string `json:"host,omitempty"`
	Name             string `json:"name,omitempty"`
	Status           string `json:"status,omitempty"`
	TemplateID       string `json:"templateid,omitempty"`
	Discover         string `json:"discover,omitempty"`
	CustomInterfaces string `json:"custom_interfaces,omitempty"`
	InventoryMode    string `json:"inventory_mode,omitempty"`
	UUID             string `json:"uuid,omitempty"`
	Flags            string `json:"flags,omitempty"`

	GroupLinks      []HostPrototypeGroupLink      `json:"groupLinks,omitempty"`
	GroupPrototypes []HostPrototypeGroupPrototype `json:"groupPrototypes,omitempty"`
	Templates       []Template                    `json:"templates,omitempty"`
	Interfaces      []HostInterface               `json:"interfaces,omitempty"`
	Tags            []Tag                         `json:"tags,omitempty"`
	Macros          []UserMacro                   `json:"macros,omitempty"`
}

// HostPrototypeGroupLink links a host prototype to an existing host group.
type HostPrototypeGroupLink struct {
	GroupPrototypeID string `json:"group_prototypeid,omitempty"`
	HostID           string `json:"hostid,omitempty"`
	GroupID          string `json:"groupid,omitempty"`
	TemplateID       string `json:"templateid,omitempty"`
}

// HostPrototypeGroupPrototype defines a host group created per discovered host.
type HostPrototypeGroupPrototype struct {
	GroupPrototypeID string `json:"group_prototypeid,omitempty"`
	HostID           string `json:"hostid,omitempty"`
	Name             string `json:"name,omitempty"`
	TemplateID       string `json:"templateid,omitempty"`
}

// HostPrototypeGetParams are the parameters for hostprototype.get.
type HostPrototypeGetParams struct {
	GetParams
	HostIDs               []string `json:"hostids,omitempty"`
	DiscoveryIDs          []string `json:"discoveryids,omitempty"`
	InheritedTags         *bool    `json:"inheritedTags,omitempty"`
	SelectGroupLinks      any      `json:"selectGroupLinks,omitempty"`
	SelectGroupPrototypes any      `json:"selectGroupPrototypes,omitempty"`
	SelectDiscoveryRule   any      `json:"selectDiscoveryRule,omitempty"`
	SelectParentHost      any      `json:"selectParentHost,omitempty"`
	SelectInterfaces      any      `json:"selectInterfaces,omitempty"`
	SelectTemplates       any      `json:"selectTemplates,omitempty"`
	SelectTags            any      `json:"selectTags,omitempty"`
	SelectMacros          any      `json:"selectMacros,omitempty"`
}

// Get retrieves host prototypes matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostprototype/get
func (s *HostPrototypeService) Get(ctx context.Context, params HostPrototypeGetParams) ([]HostPrototype, error) {
	var out []HostPrototype
	err := s.client.Call(ctx, "hostprototype.get", params, &out)
	return out, err
}

// Create creates host prototypes and returns the new ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostprototype/create
func (s *HostPrototypeService) Create(ctx context.Context, hosts ...HostPrototype) ([]string, error) {
	var res struct {
		IDs []string `json:"hostids"`
	}
	err := s.client.Call(ctx, "hostprototype.create", hosts, &res)
	return res.IDs, err
}

// Update updates host prototypes and returns the affected ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostprototype/update
func (s *HostPrototypeService) Update(ctx context.Context, hosts ...HostPrototype) ([]string, error) {
	var res struct {
		IDs []string `json:"hostids"`
	}
	err := s.client.Call(ctx, "hostprototype.update", hosts, &res)
	return res.IDs, err
}

// Delete deletes host prototypes by id and returns the deleted ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostprototype/delete
func (s *HostPrototypeService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"hostids"`
	}
	err := s.client.Call(ctx, "hostprototype.delete", ids, &res)
	return res.IDs, err
}
