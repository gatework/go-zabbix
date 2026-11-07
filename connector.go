package zabbix

import "context"

// ConnectorService wraps the "connector" API namespace.
type ConnectorService struct{ client *Client }

// Connector returns the connector service.
func (c *Client) Connector() *ConnectorService { return &ConnectorService{c} }

// Connector object (Zabbix 7.0+).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/connector/object
type Connector struct {
	ConnectorID     string `json:"connectorid,omitempty"`
	Name            string `json:"name,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
	DataType        string `json:"data_type,omitempty"`
	URL             string `json:"url,omitempty"`
	MaxRecords      string `json:"max_records,omitempty"`
	MaxSenders      string `json:"max_senders,omitempty"`
	MaxAttempts     string `json:"max_attempts,omitempty"`
	AttemptInterval string `json:"attempt_interval,omitempty"`
	Timeout         string `json:"timeout,omitempty"`
	HTTPProxy       string `json:"http_proxy,omitempty"`
	AuthType        string `json:"authtype,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	Token           string `json:"token,omitempty"`
	VerifyPeer      string `json:"verify_peer,omitempty"`
	VerifyHost      string `json:"verify_host,omitempty"`
	SSLCertFile     string `json:"ssl_cert_file,omitempty"`
	SSLKeyFile      string `json:"ssl_key_file,omitempty"`
	SSLKeyPassword  string `json:"ssl_key_password,omitempty"`
	Description     string `json:"description,omitempty"`
	Status          string `json:"status,omitempty"`
	TagsEvalType    string `json:"tags_evaltype,omitempty"`
	ItemValueType   string `json:"item_value_type,omitempty"`

	Tags []ConnectorTag `json:"tags,omitempty"`
}

// ConnectorTag is a tag filter attached to a connector.
type ConnectorTag struct {
	Tag      string `json:"tag,omitempty"`
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value,omitempty"`
}

// ConnectorGetParams are the parameters for connector.get.
type ConnectorGetParams struct {
	GetParams
	ConnectorIDs []string `json:"connectorids,omitempty"`
	DataType     *int     `json:"data_type,omitempty"`
	SelectTags   any      `json:"selectTags,omitempty"`
}

// Get retrieves connectors matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/connector/get
func (s *ConnectorService) Get(ctx context.Context, params ConnectorGetParams) ([]Connector, error) {
	var out []Connector
	err := s.client.Call(ctx, "connector.get", params, &out)
	return out, err
}

// Create creates connectors and returns the new ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/connector/create
func (s *ConnectorService) Create(ctx context.Context, connectors ...Connector) ([]string, error) {
	var res struct {
		IDs []string `json:"connectorids"`
	}
	err := s.client.Call(ctx, "connector.create", connectors, &res)
	return res.IDs, err
}

// Update updates connectors and returns the affected ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/connector/update
func (s *ConnectorService) Update(ctx context.Context, connectors ...Connector) ([]string, error) {
	var res struct {
		IDs []string `json:"connectorids"`
	}
	err := s.client.Call(ctx, "connector.update", connectors, &res)
	return res.IDs, err
}

// Delete deletes connectors by id and returns the deleted ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/connector/delete
func (s *ConnectorService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"connectorids"`
	}
	err := s.client.Call(ctx, "connector.delete", ids, &res)
	return res.IDs, err
}
