package zabbix

import "context"

// ConnectorService wraps the "connector" API namespace.
type ConnectorService struct{ client *Client }

// Connector returns the connector service.
func (c *Client) Connector() *ConnectorService { return &ConnectorService{c} }

// Connector object (Zabbix 7.0+).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/connector/object
type Connector struct {
	ConnectorID     string                 `json:"connectorid,omitempty"`
	Name            string                 `json:"name,omitempty"`
	Protocol        string                 `json:"protocol,omitempty"`
	DataType        string                 `json:"data_type,omitempty"`
	URL             string                 `json:"url,omitempty"`
	MaxRecords      string                 `json:"max_records,omitempty"`
	MaxSenders      string                 `json:"max_senders,omitempty"`
	MaxAttempts     string                 `json:"max_attempts,omitempty"`
	AttemptInterval string                 `json:"attempt_interval,omitempty"`
	Timeout         string                 `json:"timeout,omitempty"`
	HTTPProxy       string                 `json:"http_proxy,omitempty"`
	AuthType        string                 `json:"authtype,omitempty"`
	Username        string                 `json:"username,omitempty"`
	Password        string                 `json:"password,omitempty"`
	Token           string                 `json:"token,omitempty"`
	VerifyPeer      string                 `json:"verify_peer,omitempty"`
	VerifyHost      string                 `json:"verify_host,omitempty"`
	SSLCertFile     string                 `json:"ssl_cert_file,omitempty"`
	SSLKeyFile      string                 `json:"ssl_key_file,omitempty"`
	SSLKeyPassword  string                 `json:"ssl_key_password,omitempty"`
	Description     string                 `json:"description,omitempty"`
	Status          string                 `json:"status,omitempty"`
	TagsEvalType    string                 `json:"tags_evaltype,omitempty"`
	ItemValueType   string                 `json:"item_value_type,omitempty"`
	Tags            Relation[ConnectorTag] `json:"tags,omitzero"`
}

// ConnectorTag is a tag filter attached to a connector.
type ConnectorTag struct {
	Tag      string `json:"tag,omitempty"`
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value,omitempty"`
}

// ConnectorGetParams are the parameters for connector.get.
// ConnectorGetParams queries the list form of connector.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type ConnectorGetParams struct {
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
	ConnectorIDs           []string       `json:"connectorids,omitzero"`
	SelectTags             any            `json:"selectTags,omitempty"`
}

// Get retrieves connectors matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/connector/get
func (s *ConnectorService) Get(ctx context.Context, params ConnectorGetParams) ([]Connector, error) {
	return getList[Connector](ctx, s.client, "connector.get", params)
}

// ConnectorCountParams contains the filters for connector.get counts.
type ConnectorCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	ConnectorIDs           []string       `json:"connectorids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *ConnectorService) Count(ctx context.Context, params ConnectorCountParams) (uint64, error) {
	return getCount(ctx, s.client, "connector.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *ConnectorService) GetMap(ctx context.Context, params ConnectorGetParams) (map[string]Connector, error) {
	return getMap[Connector](ctx, s.client, "connector.get", params)
}

// Create creates connectors and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/connector/create
func (s *ConnectorService) Create(ctx context.Context, connectors ...Connector) ([]string, error) {
	var res struct {
		IDs []string `json:"connectorids"`
	}
	err := s.client.Call(ctx, "connector.create", connectors, &res)
	return res.IDs, err
}

// Update updates connectors and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/connector/update
func (s *ConnectorService) Update(ctx context.Context, connectors ...Connector) ([]string, error) {
	var res struct {
		IDs []string `json:"connectorids"`
	}
	err := s.client.Call(ctx, "connector.update", connectors, &res)
	return res.IDs, err
}

// Delete deletes connectors by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/connector/delete
func (s *ConnectorService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"connectorids"`
	}
	err := s.client.Call(ctx, "connector.delete", ids, &res)
	return res.IDs, err
}
