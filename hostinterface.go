package zabbix

import (
	"context"
	"errors"
)

// HostInterfaceService wraps the "hostinterface" API namespace.
type HostInterfaceService struct{ client *Client }

// HostInterface returns the hostinterface service.
func (c *Client) HostInterface() *HostInterfaceService { return &HostInterfaceService{c} }

// HostInterfaceGetParams are the parameters for hostinterface.get.
// HostInterfaceGetParams queries the list form of hostinterface.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type HostInterfaceGetParams struct {
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
	InterfaceIDs           []string       `json:"interfaceids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	TriggerIDs             []string       `json:"triggerids,omitzero"`
	ItemIDs                []string       `json:"itemids,omitzero"`
	SelectItems            any            `json:"selectItems,omitempty"`
	SelectHosts            any            `json:"selectHosts,omitempty"`
}

// Get retrieves host interfaces matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/get
func (s *HostInterfaceService) Get(ctx context.Context, params HostInterfaceGetParams) ([]HostInterface, error) {
	return getList[HostInterface](ctx, s.client, "hostinterface.get", params)
}

// HostInterfaceCountParams contains the filters for hostinterface.get counts.
type HostInterfaceCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	InterfaceIDs           []string       `json:"interfaceids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	TriggerIDs             []string       `json:"triggerids,omitzero"`
	ItemIDs                []string       `json:"itemids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *HostInterfaceService) Count(ctx context.Context, params HostInterfaceCountParams) (uint64, error) {
	return getCount(ctx, s.client, "hostinterface.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *HostInterfaceService) GetMap(ctx context.Context, params HostInterfaceGetParams) (map[string]HostInterface, error) {
	return getMap[HostInterface](ctx, s.client, "hostinterface.get", params)
}

// Create creates host interfaces and returns the new interface ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/create
func (s *HostInterfaceService) Create(ctx context.Context, interfaces ...HostInterfaceCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "hostinterface.create", "interfaceids", interfaces)
}

// Update updates host interfaces and returns the affected interface ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/update
func (s *HostInterfaceService) Update(ctx context.Context, interfaces ...HostInterfaceUpdateParams) ([]string, error) {
	return updateIDs(ctx, s.client, "hostinterface.update", "interfaceids", interfaces)
}

// Delete deletes host interfaces by id and returns the deleted interface ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/delete
func (s *HostInterfaceService) Delete(ctx context.Context, interfaceIDs ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "hostinterface.delete", "interfaceids", interfaceIDs)
}

// HostInterfaceMassAddParams are the parameters for hostinterface.massadd.
type HostInterfaceMassAddParams struct {
	Hosts      []HostRef            `json:"hosts"`
	Interfaces []HostInterfaceInput `json:"interfaces"`
}

func (p HostInterfaceMassAddParams) validate() error {
	if len(p.Hosts) == 0 || len(p.Interfaces) == 0 {
		return errors.New("zabbix: hostinterface.massadd requires hosts and interfaces")
	}
	for _, host := range p.Hosts {
		if err := host.validate(); err != nil {
			return err
		}
	}
	for _, item := range p.Interfaces {
		if err := item.validate(); err != nil {
			return err
		}
	}
	return nil
}

// MassAdd adds interfaces to the given hosts and returns the new interface ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/massadd
func (s *HostInterfaceService) MassAdd(ctx context.Context, params HostInterfaceMassAddParams) ([]string, error) {
	return mutateObjectIDs(ctx, s.client, "hostinterface.massadd", "interfaceids", params, false)
}

// HostInterfaceMassRemoveParams are the parameters for hostinterface.massremove.
type HostInterfaceMassRemoveParams struct {
	HostIDs    []string                  `json:"hostids"`
	Interfaces []HostInterfaceMatchInput `json:"interfaces"`
}

func (p HostInterfaceMassRemoveParams) validate() error {
	if len(p.HostIDs) == 0 || len(p.Interfaces) == 0 {
		return errors.New("zabbix: hostinterface.massremove requires host IDs and interfaces")
	}
	for _, id := range p.HostIDs {
		if err := requireWriteID(id, "HostInterfaceMassRemoveParams.HostIDs"); err != nil {
			return err
		}
	}
	for _, item := range p.Interfaces {
		if err := item.validate(); err != nil {
			return err
		}
	}
	return nil
}

// MassRemove removes the given interfaces from the hosts and returns the
// affected interface ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/massremove
func (s *HostInterfaceService) MassRemove(ctx context.Context, params HostInterfaceMassRemoveParams) ([]string, error) {
	return mutateObjectIDs(ctx, s.client, "hostinterface.massremove", "interfaceids", params, true)
}

// ReplaceHostInterfaces replaces all interfaces on a host with the given set and
// returns the resulting interface ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/replacehostinterfaces
func (s *HostInterfaceService) ReplaceHostInterfaces(ctx context.Context, params HostInterfaceReplaceParams) ([]string, error) {
	return mutateObjectIDs(ctx, s.client, "hostinterface.replacehostinterfaces", "interfaceids", params, true)
}
