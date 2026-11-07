package zabbix

import "context"

// APITokenService wraps the "token" API namespace.
type APITokenService struct{ client *Client }

// APIToken returns the token service. (Named APIToken to avoid colliding with the
// Client.Token accessor for the current session token.)
func (c *Client) APIToken() *APITokenService { return &APITokenService{c} }

// Token object (API token).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/token/object
type Token struct {
	TokenID       string `json:"tokenid,omitempty"`
	Name          string `json:"name,omitempty"`
	Description   string `json:"description,omitempty"`
	UserID        string `json:"userid,omitempty"`
	LastAccess    string `json:"lastaccess,omitempty"`
	Status        string `json:"status,omitempty"`
	ExpiresAt     string `json:"expires_at,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	CreatorUserID string `json:"creator_userid,omitempty"`
}

// TokenGenerated is the result of token.generate: the token id and its secret string.
type TokenGenerated struct {
	TokenID string `json:"tokenid,omitempty"`
	Token   string `json:"token,omitempty"`
}

// TokenGetParams are the parameters for token.get.
type TokenGetParams struct {
	GetParams
	TokenIDs []string `json:"tokenids,omitempty"`
	UserIDs  []string `json:"userids,omitempty"`
	Valid    *bool    `json:"valid,omitempty"`
}

// Get retrieves tokens matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/token/get
func (s *APITokenService) Get(ctx context.Context, params TokenGetParams) ([]Token, error) {
	var out []Token
	err := s.client.Call(ctx, "token.get", params, &out)
	return out, err
}

// Create creates tokens and returns the new ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/token/create
func (s *APITokenService) Create(ctx context.Context, items ...Token) ([]string, error) {
	var res struct {
		IDs []string `json:"tokenids"`
	}
	err := s.client.Call(ctx, "token.create", items, &res)
	return res.IDs, err
}

// Update updates tokens and returns the affected ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/token/update
func (s *APITokenService) Update(ctx context.Context, items ...Token) ([]string, error) {
	var res struct {
		IDs []string `json:"tokenids"`
	}
	err := s.client.Call(ctx, "token.update", items, &res)
	return res.IDs, err
}

// Delete deletes tokens by id and returns the deleted ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/token/delete
func (s *APITokenService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"tokenids"`
	}
	err := s.client.Call(ctx, "token.delete", ids, &res)
	return res.IDs, err
}

// Generate generates the authorization secrets for the given tokens and returns the
// generated token strings.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/token/generate
func (s *APITokenService) Generate(ctx context.Context, tokenIDs ...string) ([]TokenGenerated, error) {
	var out []TokenGenerated
	err := s.client.Call(ctx, "token.generate", tokenIDs, &out)
	return out, err
}
