package zabbix

import "context"

// AuthenticationService wraps the "authentication" API namespace, a global singleton
// with only get and update.
type AuthenticationService struct{ client *Client }

// Authentication returns the authentication service.
func (c *Client) Authentication() *AuthenticationService { return &AuthenticationService{c} }

// Authentication object (global authentication settings).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/authentication/object
type Authentication struct {
	AuthenticationType   string `json:"authentication_type,omitempty"`
	HTTPAuthEnabled      string `json:"http_auth_enabled,omitempty"`
	HTTPLoginForm        string `json:"http_login_form,omitempty"`
	HTTPStripDomains     string `json:"http_strip_domains,omitempty"`
	HTTPCaseSensitive    string `json:"http_case_sensitive,omitempty"`
	LDAPAuthEnabled      string `json:"ldap_auth_enabled,omitempty"`
	LDAPUserDirectoryID  string `json:"ldap_userdirectoryid,omitempty"`
	LDAPCaseSensitive    string `json:"ldap_case_sensitive,omitempty"`
	SAMLAuthEnabled      string `json:"saml_auth_enabled,omitempty"`
	PasswdMinLength      string `json:"passwd_min_length,omitempty"`
	PasswdCheckRules     string `json:"passwd_check_rules,omitempty"`
	JITProvisionInterval string `json:"jit_provision_interval,omitempty"`
	DisabledUsrGrpID     string `json:"disabled_usrgrpid,omitempty"`
	MFAStatus            string `json:"mfa_status,omitempty"`
	MFAID                string `json:"mfaid,omitempty"`
}

// Get retrieves the global authentication settings.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/authentication/get
func (s *AuthenticationService) Get(ctx context.Context) (*Authentication, error) {
	return getSingleton[Authentication](ctx, s.client, "authentication.get")
}

// Update updates the global authentication settings and returns the updated
// parameter names.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/authentication/update
func (s *AuthenticationService) Update(ctx context.Context, params Authentication) ([]string, error) {
	var res []string
	err := s.client.Call(ctx, "authentication.update", params, &res)
	return res, err
}
