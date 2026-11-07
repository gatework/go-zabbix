package zabbix

import "context"

// HostInterfaceService wraps the "hostinterface" API namespace.
type HostInterfaceService struct{ client *Client }

// HostInterface returns the hostinterface service.
func (c *Client) HostInterface() *HostInterfaceService { return &HostInterfaceService{c} }

// HostInterfaceGetParams are the parameters for hostinterface.get.
type HostInterfaceGetParams struct {
	GetParams
	InterfaceIDs []string `json:"interfaceids,omitempty"`
	HostIDs      []string `json:"hostids,omitempty"`
	TriggerIDs   []string `json:"triggerids,omitempty"`
	ItemIDs      []string `json:"itemids,omitempty"`
	SelectItems  any      `json:"selectItems,omitempty"`
	SelectHosts  any      `json:"selectHosts,omitempty"`
}

// Get retrieves host interfaces matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostinterface/get
func (s *HostInterfaceService) Get(ctx context.Context, params HostInterfaceGetParams) ([]HostInterface, error) {
	var out []HostInterface
	err := s.client.Call(ctx, "hostinterface.get", params, &out)
	return out, err
}

// Create creates host interfaces and returns the new interface ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostinterface/create
func (s *HostInterfaceService) Create(ctx context.Context, interfaces ...HostInterface) ([]string, error) {
	var res struct {
		IDs []string `json:"interfaceids"`
	}
	err := s.client.Call(ctx, "hostinterface.create", interfaces, &res)
	return res.IDs, err
}

// Update updates host interfaces and returns the affected interface ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostinterface/update
func (s *HostInterfaceService) Update(ctx context.Context, interfaces ...HostInterface) ([]string, error) {
	var res struct {
		IDs []string `json:"interfaceids"`
	}
	err := s.client.Call(ctx, "hostinterface.update", interfaces, &res)
	return res.IDs, err
}

// Delete deletes host interfaces by id and returns the deleted interface ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostinterface/delete
func (s *HostInterfaceService) Delete(ctx context.Context, interfaceIDs ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"interfaceids"`
	}
	err := s.client.Call(ctx, "hostinterface.delete", interfaceIDs, &res)
	return res.IDs, err
}

// HostInterfaceMassAddParams are the parameters for hostinterface.massadd.
type HostInterfaceMassAddParams struct {
	Hosts      any             `json:"hosts"`
	Interfaces []HostInterface `json:"interfaces"`
}

// MassAdd adds interfaces to the given hosts and returns the new interface ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostinterface/massadd
func (s *HostInterfaceService) MassAdd(ctx context.Context, params HostInterfaceMassAddParams) ([]string, error) {
	var res struct {
		IDs []string `json:"interfaceids"`
	}
	err := s.client.Call(ctx, "hostinterface.massadd", params, &res)
	return res.IDs, err
}

// HostInterfaceMassRemoveParams are the parameters for hostinterface.massremove.
type HostInterfaceMassRemoveParams struct {
	HostIDs    []string        `json:"hostids"`
	Interfaces []HostInterface `json:"interfaces"`
}

// MassRemove removes the given interfaces from the hosts and returns the
// affected interface ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostinterface/massremove
func (s *HostInterfaceService) MassRemove(ctx context.Context, params HostInterfaceMassRemoveParams) ([]string, error) {
	var res struct {
		IDs []string `json:"interfaceids"`
	}
	err := s.client.Call(ctx, "hostinterface.massremove", params, &res)
	return res.IDs, err
}

// ReplaceHostInterfaces replaces all interfaces on a host with the given set and
// returns the resulting interface ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/hostinterface/replacehostinterfaces
func (s *HostInterfaceService) ReplaceHostInterfaces(ctx context.Context, hostID string, interfaces []HostInterface) ([]string, error) {
	params := struct {
		HostID     string          `json:"hostid"`
		Interfaces []HostInterface `json:"interfaces"`
	}{HostID: hostID, Interfaces: interfaces}
	var res struct {
		IDs []string `json:"interfaceids"`
	}
	err := s.client.Call(ctx, "hostinterface.replacehostinterfaces", params, &res)
	return res.IDs, err
}
