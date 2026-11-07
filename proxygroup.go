package zabbix

import "context"

// ProxyGroupService wraps the "proxygroup" API namespace.
type ProxyGroupService struct{ client *Client }

// ProxyGroup returns the proxygroup service.
func (c *Client) ProxyGroup() *ProxyGroupService { return &ProxyGroupService{c} }

// ProxyGroup object (Zabbix 7.0+).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxygroup/object
type ProxyGroup struct {
	ProxyGroupID  string `json:"proxy_groupid,omitempty"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	FailoverDelay string `json:"failover_delay,omitempty"`
	MinOnline     string `json:"min_online,omitempty"`
	State         string `json:"state,omitempty"`

	Proxies []ProxyGroupProxy `json:"proxies,omitempty"`
}

// ProxyGroupProxy is a proxy that is a member of a proxy group.
type ProxyGroupProxy struct {
	ProxyID string `json:"proxyid,omitempty"`
	Name    string `json:"name,omitempty"`
}

// ProxyGroupGetParams are the parameters for proxygroup.get.
type ProxyGroupGetParams struct {
	GetParams
	ProxyGroupIDs []string `json:"proxy_groupids,omitempty"`
	SelectProxies any      `json:"selectProxies,omitempty"`
}

// Get retrieves proxy groups matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxygroup/get
func (s *ProxyGroupService) Get(ctx context.Context, params ProxyGroupGetParams) ([]ProxyGroup, error) {
	var out []ProxyGroup
	err := s.client.Call(ctx, "proxygroup.get", params, &out)
	return out, err
}

// Create creates proxy groups and returns the new ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxygroup/create
func (s *ProxyGroupService) Create(ctx context.Context, groups ...ProxyGroup) ([]string, error) {
	var res struct {
		IDs []string `json:"proxy_groupids"`
	}
	err := s.client.Call(ctx, "proxygroup.create", groups, &res)
	return res.IDs, err
}

// Update updates proxy groups and returns the affected ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxygroup/update
func (s *ProxyGroupService) Update(ctx context.Context, groups ...ProxyGroup) ([]string, error) {
	var res struct {
		IDs []string `json:"proxy_groupids"`
	}
	err := s.client.Call(ctx, "proxygroup.update", groups, &res)
	return res.IDs, err
}

// Delete deletes proxy groups by id and returns the deleted ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxygroup/delete
func (s *ProxyGroupService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"proxy_groupids"`
	}
	err := s.client.Call(ctx, "proxygroup.delete", ids, &res)
	return res.IDs, err
}
