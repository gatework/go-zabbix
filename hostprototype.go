package zabbix

import "context"

// HostPrototypeService wraps the "hostprototype" API namespace.
type HostPrototypeService struct{ client *Client }

// HostPrototype returns the hostprototype service.
func (c *Client) HostPrototype() *HostPrototypeService { return &HostPrototypeService{c} }

// HostPrototype object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostprototype/object
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
	Templates       Relation[Template]            `json:"templates,omitzero"`
	Interfaces      []HostInterface               `json:"interfaces,omitempty"`
	Tags            []Tag                         `json:"tags,omitempty"`
	Macros          []UserMacro                   `json:"macros,omitempty"`
	DiscoveryRule   SelectedObject[DiscoveryRule] `json:"discoveryRule,omitzero"`
	ParentHost      SelectedObject[Host]          `json:"parentHost,omitzero"`
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
// HostPrototypeGetParams queries the list form of hostprototype.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type HostPrototypeGetParams struct {
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
	DiscoveryIDs           []string       `json:"discoveryids,omitzero"`
	SelectGroupLinks       any            `json:"selectGroupLinks,omitempty"`
	SelectGroupPrototypes  any            `json:"selectGroupPrototypes,omitempty"`
	SelectDiscoveryRule    any            `json:"selectDiscoveryRule,omitempty"`
	SelectParentHost       any            `json:"selectParentHost,omitempty"`
	SelectInterfaces       any            `json:"selectInterfaces,omitempty"`
	SelectTemplates        any            `json:"selectTemplates,omitempty"`
	SelectTags             any            `json:"selectTags,omitempty"`
	SelectMacros           any            `json:"selectMacros,omitempty"`
}

// Get retrieves host prototypes matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostprototype/get
func (s *HostPrototypeService) Get(ctx context.Context, params HostPrototypeGetParams) ([]HostPrototype, error) {
	return getList[HostPrototype](ctx, s.client, "hostprototype.get", params)
}

// HostPrototypeCountParams contains the filters for hostprototype.get counts.
type HostPrototypeCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	HostIDs                []string       `json:"hostids,omitzero"`
	DiscoveryIDs           []string       `json:"discoveryids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *HostPrototypeService) Count(ctx context.Context, params HostPrototypeCountParams) (uint64, error) {
	return getCount(ctx, s.client, "hostprototype.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *HostPrototypeService) GetMap(ctx context.Context, params HostPrototypeGetParams) (map[string]HostPrototype, error) {
	return getMap[HostPrototype](ctx, s.client, "hostprototype.get", params)
}

// Create creates host prototypes and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostprototype/create
func (s *HostPrototypeService) Create(ctx context.Context, hosts ...HostPrototypeCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "hostprototype.create", "hostids", hosts)
}

// Update updates host prototypes and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostprototype/update
func (s *HostPrototypeService) Update(ctx context.Context, hosts ...HostPrototypeUpdateParams) ([]string, error) {
	return updateIDs(ctx, s.client, "hostprototype.update", "hostids", hosts)
}

// Delete deletes host prototypes by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostprototype/delete
func (s *HostPrototypeService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "hostprototype.delete", "hostids", ids)
}
