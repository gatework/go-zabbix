package zabbix

import "context"

// ProxyService wraps the "proxy" API namespace.
type ProxyService struct{ client *Client }

// Proxy returns the proxy service.
func (c *Client) Proxy() *ProxyService { return &ProxyService{c} }

// Proxy object (Zabbix 7.0+).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/object
type Proxy struct {
	ProxyID              string `json:"proxyid,omitempty"`
	Name                 string `json:"name,omitempty"`
	ProxyGroupID         string `json:"proxy_groupid,omitempty"`
	LocalAddress         string `json:"local_address,omitempty"`
	LocalPort            string `json:"local_port,omitempty"`
	OperatingMode        string `json:"operating_mode,omitempty"`
	AllowedAddresses     string `json:"allowed_addresses,omitempty"`
	Address              string `json:"address,omitempty"`
	Port                 string `json:"port,omitempty"`
	Description          string `json:"description,omitempty"`
	TLSConnect           string `json:"tls_connect,omitempty"`
	TLSAccept            string `json:"tls_accept,omitempty"`
	TLSIssuer            string `json:"tls_issuer,omitempty"`
	TLSSubject           string `json:"tls_subject,omitempty"`
	CustomTimeouts       string `json:"custom_timeouts,omitempty"`
	TimeoutZabbixAgent   string `json:"timeout_zabbix_agent,omitempty"`
	TimeoutSimpleCheck   string `json:"timeout_simple_check,omitempty"`
	TimeoutSNMPAgent     string `json:"timeout_snmp_agent,omitempty"`
	TimeoutExternalCheck string `json:"timeout_external_check,omitempty"`
	TimeoutDBMonitor     string `json:"timeout_db_monitor,omitempty"`
	TimeoutHTTPAgent     string `json:"timeout_http_agent,omitempty"`
	TimeoutSSHAgent      string `json:"timeout_ssh_agent,omitempty"`
	TimeoutTelnetAgent   string `json:"timeout_telnet_agent,omitempty"`
	TimeoutScript        string `json:"timeout_script,omitempty"`
	TimeoutBrowser       string `json:"timeout_browser,omitempty"`
	LastAccess           string `json:"lastaccess,omitempty"`
	Version              string `json:"version,omitempty"`
	Compatibility        string `json:"compatibility,omitempty"`
	State                string `json:"state,omitempty"`

	// Hosts is used both as the assignment INPUT ([{hostid}]) and the get OUTPUT.
	Hosts []ProxyHost `json:"hosts,omitempty"`
}

// ProxyHost is a host monitored by a proxy.
type ProxyHost struct {
	HostID string `json:"hostid,omitempty"`
	Host   string `json:"host,omitempty"`
}

// ProxyGetParams are the parameters for proxy.get.
type ProxyGetParams struct {
	GetParams
	ProxyIDs         []string `json:"proxyids,omitempty"`
	ProxyGroupIDs    []string `json:"proxy_groupids,omitempty"`
	SelectHosts      any      `json:"selectHosts,omitempty"`
	SelectProxyGroup any      `json:"selectProxyGroup,omitempty"`
}

// Get retrieves proxies matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/get
func (s *ProxyService) Get(ctx context.Context, params ProxyGetParams) ([]Proxy, error) {
	var out []Proxy
	err := s.client.Call(ctx, "proxy.get", params, &out)
	return out, err
}

// Create creates proxies and returns the new ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/create
func (s *ProxyService) Create(ctx context.Context, proxies ...Proxy) ([]string, error) {
	var res struct {
		IDs []string `json:"proxyids"`
	}
	err := s.client.Call(ctx, "proxy.create", proxies, &res)
	return res.IDs, err
}

// Update updates proxies and returns the affected ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/update
func (s *ProxyService) Update(ctx context.Context, proxies ...Proxy) ([]string, error) {
	var res struct {
		IDs []string `json:"proxyids"`
	}
	err := s.client.Call(ctx, "proxy.update", proxies, &res)
	return res.IDs, err
}

// Delete deletes proxies by id and returns the deleted ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/proxy/delete
func (s *ProxyService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"proxyids"`
	}
	err := s.client.Call(ctx, "proxy.delete", ids, &res)
	return res.IDs, err
}
