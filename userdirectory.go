package zabbix

import "context"

// UserDirectoryService wraps the "userdirectory" API namespace.
type UserDirectoryService struct{ client *Client }

// UserDirectory returns the userdirectory service.
func (c *Client) UserDirectory() *UserDirectoryService { return &UserDirectoryService{c} }

// UserDirectory object (LDAP or SAML identity provider configuration).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/userdirectory/object
type UserDirectory struct {
	UserDirectoryID string `json:"userdirectoryid,omitempty"`
	Name            string `json:"name,omitempty"`
	IDPType         string `json:"idp_type,omitempty"`
	ProvisionStatus string `json:"provision_status,omitempty"`
	Description     string `json:"description,omitempty"`
	Host            string `json:"host,omitempty"`
	Port            string `json:"port,omitempty"`
	BaseDN          string `json:"base_dn,omitempty"`
	SearchAttribute string `json:"search_attribute,omitempty"`
	BindDN          string `json:"bind_dn,omitempty"`
	StartTLS        string `json:"start_tls,omitempty"`
	SearchFilter    string `json:"search_filter,omitempty"`
	GroupBaseDN     string `json:"group_basedn,omitempty"`
	UserUsername    string `json:"user_username,omitempty"`
	UserLastname    string `json:"user_lastname,omitempty"`

	// SAML-specific fields.
	IDPEntityID       string `json:"idp_entityid,omitempty"`
	SSOURL            string `json:"sso_url,omitempty"`
	SLOURL            string `json:"slo_url,omitempty"`
	UsernameAttribute string `json:"username_attribute,omitempty"`
	SPEntityID        string `json:"sp_entityid,omitempty"`
}

// UserDirectoryGetParams are the parameters for userdirectory.get.
// UserDirectoryGetParams queries the list form of userdirectory.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type UserDirectoryGetParams struct {
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
	UserDirectoryIDs       []string       `json:"userdirectoryids,omitzero"`
}

// Get retrieves user directories matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/userdirectory/get
func (s *UserDirectoryService) Get(ctx context.Context, params UserDirectoryGetParams) ([]UserDirectory, error) {
	return getList[UserDirectory](ctx, s.client, "userdirectory.get", params)
}

// UserDirectoryCountParams contains the filters for userdirectory.get counts.
type UserDirectoryCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	UserDirectoryIDs       []string       `json:"userdirectoryids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *UserDirectoryService) Count(ctx context.Context, params UserDirectoryCountParams) (uint64, error) {
	return getCount(ctx, s.client, "userdirectory.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *UserDirectoryService) GetMap(ctx context.Context, params UserDirectoryGetParams) (map[string]UserDirectory, error) {
	return getMap[UserDirectory](ctx, s.client, "userdirectory.get", params)
}

// Create creates user directories and returns the new ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/userdirectory/create
func (s *UserDirectoryService) Create(ctx context.Context, items ...UserDirectoryCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "userdirectory.create", "userdirectoryids", items)
}

// Update updates user directories and returns the affected ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/userdirectory/update
func (s *UserDirectoryService) Update(ctx context.Context, items ...UserDirectoryUpdateParams) ([]string, error) {
	return updateIDs(ctx, s.client, "userdirectory.update", "userdirectoryids", items)
}

// Delete deletes user directories by id and returns the deleted ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/userdirectory/delete
func (s *UserDirectoryService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "userdirectory.delete", "userdirectoryids", ids)
}

// Test checks the connection and configuration of a user directory. The params and
// result shapes depend on the idp_type, so both are passed through untyped.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/userdirectory/test
func (s *UserDirectoryService) Test(ctx context.Context, params any) (any, error) {
	var out any
	err := s.client.Call(ctx, "userdirectory.test", params, &out)
	return out, err
}
