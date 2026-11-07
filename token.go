package zabbix

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
)

// APITokenService wraps the "token" API namespace.
type APITokenService struct{ client *Client }

// APIToken returns the token administration service.
func (c *Client) APIToken() *APITokenService { return &APITokenService{c} }

// Token object (API token).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/token/object
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
// TokenGetParams queries the list form of token.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type TokenGetParams struct {
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
	TokenIDs               []string       `json:"tokenids,omitzero"`
	UserIDs                []string       `json:"userids,omitzero"`
	ValidAt                *int64         `json:"valid_at,omitempty"`
	ExpiredAt              *int64         `json:"expired_at,omitempty"`
}

// Get retrieves tokens matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/token/get
func (s *APITokenService) Get(ctx context.Context, params TokenGetParams) ([]Token, error) {
	return getList[Token](ctx, s.client, "token.get", params)
}

// TokenCountParams contains the filters for token.get counts.
type TokenCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	TokenIDs               []string       `json:"tokenids,omitzero"`
	UserIDs                []string       `json:"userids,omitzero"`
	ValidAt                *int64         `json:"valid_at,omitempty"`
	ExpiredAt              *int64         `json:"expired_at,omitempty"`
}

// Count returns the number of objects matching params.
func (s *APITokenService) Count(ctx context.Context, params TokenCountParams) (uint64, error) {
	return getCount(ctx, s.client, "token.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *APITokenService) GetMap(ctx context.Context, params TokenGetParams) (map[string]Token, error) {
	return getMap[Token](ctx, s.client, "token.get", params)
}

// Create creates tokens and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/token/create
func (s *APITokenService) Create(ctx context.Context, items ...TokenCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "token.create", "tokenids", items)
}

// Update updates tokens and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/token/update
func (s *APITokenService) Update(ctx context.Context, items ...TokenUpdateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "token.update", "tokenids", items)
}

// Delete deletes tokens by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/token/delete
func (s *APITokenService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "token.delete", "tokenids", ids)
}

// Generate generates the authorization secrets for the given tokens and returns the
// generated token strings.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/token/generate
func (s *APITokenService) Generate(ctx context.Context, tokenIDs ...string) ([]TokenGenerated, error) {
	if len(tokenIDs) == 0 {
		return nil, requestError("token.generate", ErrorInvalidParams, errors.New("token IDs must not be empty"))
	}
	expected := make(map[string]struct{}, len(tokenIDs))
	for _, id := range tokenIDs {
		if err := requireWriteID(id, "tokenid"); err != nil {
			return nil, requestError("token.generate", ErrorInvalidParams, err)
		}
		canonical := strings.TrimLeft(id, "0")
		if _, exists := expected[canonical]; exists {
			return nil, requestError("token.generate", ErrorInvalidParams, errors.New("token IDs must be unique"))
		}
		expected[canonical] = struct{}{}
	}
	var raw json.RawMessage
	if err := s.client.Call(ctx, "token.generate", tokenIDs, &raw); err != nil {
		return nil, err
	}
	out, err := decodeObjectList[TokenGenerated](raw)
	if err != nil {
		return nil, requestError("token.generate", ErrorDecode, err)
	}
	// Zabbix 7.4 CToken::generate returns one ID/secret pair for each requested
	// token. Do not report rotation success if any new secret is missing.
	if len(out) != len(tokenIDs) {
		return nil, requestError("token.generate", ErrorProtocol, errors.New("generated token count does not match request"))
	}
	for _, generated := range out {
		canonical := strings.TrimLeft(generated.TokenID, "0")
		if _, exists := expected[canonical]; !exists || strings.TrimSpace(generated.Token) == "" {
			return nil, requestError("token.generate", ErrorProtocol, errors.New("generated tokens require requested IDs and non-empty secrets"))
		}
		delete(expected, canonical)
	}
	return out, nil
}
