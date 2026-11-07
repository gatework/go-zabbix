package zabbix

import "context"

// ProxyService wraps the "proxy" API namespace.
type ProxyService struct{ client *Client }

// Proxy returns the proxy service.
func (c *Client) Proxy() *ProxyService { return &ProxyService{c} }

// Proxy object (Zabbix 7.0+).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxy/object
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

	// Hosts is populated by selectHosts.
	Hosts      Relation[Host]             `json:"hosts,omitzero"`
	ProxyGroup SelectedObject[ProxyGroup] `json:"proxyGroup,omitzero"`
}

// ProxyHost is a host monitored by a proxy.
type ProxyHost struct {
	HostID string `json:"hostid,omitempty"`
	Host   string `json:"host,omitempty"`
}

// ProxyGetParams are the parameters for proxy.get.
// ProxyGetParams queries the list form of proxy.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type ProxyGetParams struct {
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
	ProxyIDs               []string       `json:"proxyids,omitzero"`
	ProxyGroupIDs          []string       `json:"proxy_groupids,omitzero"`
	SelectHosts            any            `json:"selectHosts,omitempty"`
	SelectProxyGroup       any            `json:"selectProxyGroup,omitempty"`
}

// Get retrieves proxies matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxy/get
func (s *ProxyService) Get(ctx context.Context, params ProxyGetParams) ([]Proxy, error) {
	return getList[Proxy](ctx, s.client, "proxy.get", params)
}

// ProxyCountParams contains the filters for proxy.get counts.
type ProxyCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	ProxyIDs               []string       `json:"proxyids,omitzero"`
	ProxyGroupIDs          []string       `json:"proxy_groupids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *ProxyService) Count(ctx context.Context, params ProxyCountParams) (uint64, error) {
	return getCount(ctx, s.client, "proxy.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *ProxyService) GetMap(ctx context.Context, params ProxyGetParams) (map[string]Proxy, error) {
	return getMap[Proxy](ctx, s.client, "proxy.get", params)
}

// Create creates proxies and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxy/create
func (s *ProxyService) Create(ctx context.Context, proxies ...ProxyCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "proxy.create", "proxyids", proxies)
}

// Update updates proxies and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxy/update
func (s *ProxyService) Update(ctx context.Context, proxies ...ProxyUpdateParams) ([]string, error) {
	return updateIDs(ctx, s.client, "proxy.update", "proxyids", proxies)
}

// Delete deletes proxies by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxy/delete
func (s *ProxyService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "proxy.delete", "proxyids", ids)
}
