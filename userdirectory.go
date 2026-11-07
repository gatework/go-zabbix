package zabbix

import "context"

// UserDirectoryService wraps the "userdirectory" API namespace.
type UserDirectoryService struct{ client *Client }

// UserDirectory returns the userdirectory service.
func (c *Client) UserDirectory() *UserDirectoryService { return &UserDirectoryService{c} }

// UserDirectory object (LDAP or SAML identity provider configuration).
// https://www.zabbix.com/documentation/current/en/manual/api/reference/userdirectory/object
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
type UserDirectoryGetParams struct {
	GetParams
	UserDirectoryIDs []string `json:"userdirectoryids,omitempty"`
}

// Get retrieves user directories matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/userdirectory/get
func (s *UserDirectoryService) Get(ctx context.Context, params UserDirectoryGetParams) ([]UserDirectory, error) {
	var out []UserDirectory
	err := s.client.Call(ctx, "userdirectory.get", params, &out)
	return out, err
}

// Create creates user directories and returns the new ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/userdirectory/create
func (s *UserDirectoryService) Create(ctx context.Context, items ...UserDirectory) ([]string, error) {
	var res struct {
		IDs []string `json:"userdirectoryids"`
	}
	err := s.client.Call(ctx, "userdirectory.create", items, &res)
	return res.IDs, err
}

// Update updates user directories and returns the affected ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/userdirectory/update
func (s *UserDirectoryService) Update(ctx context.Context, items ...UserDirectory) ([]string, error) {
	var res struct {
		IDs []string `json:"userdirectoryids"`
	}
	err := s.client.Call(ctx, "userdirectory.update", items, &res)
	return res.IDs, err
}

// Delete deletes user directories by id and returns the deleted ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/userdirectory/delete
func (s *UserDirectoryService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"userdirectoryids"`
	}
	err := s.client.Call(ctx, "userdirectory.delete", ids, &res)
	return res.IDs, err
}

// Test checks the connection and configuration of a user directory. The params and
// result shapes depend on the idp_type, so both are passed through untyped.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/userdirectory/test
func (s *UserDirectoryService) Test(ctx context.Context, params any) (any, error) {
	var out any
	err := s.client.Call(ctx, "userdirectory.test", params, &out)
	return out, err
}
