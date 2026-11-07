package zabbix

import "errors"

// Write parameters are independent of returned objects. Optional scalar pointers
// distinguish an omitted property from an explicit zero or empty value. Optional
// collections use omitzero: nil leaves a relation unchanged; an allocated empty
// collection is sent so the server can clear it (where that operation permits).
// Server-dependent restrictions, including inherited/provisioned object rules,
// remain enforced by Zabbix.

// HostPrototypeCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostprototype/object
type HostPrototypeCreateParams struct {
	Host             string                        `json:"host"`
	Name             *string                       `json:"name,omitempty"`
	Status           *int                          `json:"status,omitempty"`
	Discover         *int                          `json:"discover,omitempty"`
	CustomInterfaces *int                          `json:"custom_interfaces,omitempty"`
	InventoryMode    *int                          `json:"inventory_mode,omitempty"`
	UUID             *string                       `json:"uuid,omitempty"`
	GroupLinks       []HostGroupRef                `json:"groupLinks"`
	GroupPrototypes  []HostPrototypeGroupInput     `json:"groupPrototypes,omitzero"`
	Templates        []TemplateRef                 `json:"templates,omitzero"`
	Interfaces       []HostPrototypeInterfaceInput `json:"interfaces,omitzero"`
	Tags             []TagInput                    `json:"tags,omitzero"`
	Macros           []HostMacroInput              `json:"macros,omitzero"`
	RuleID           string                        `json:"ruleid"`
}

func (p HostPrototypeCreateParams) validate() error {
	if p.Host == "" {
		return errors.New("zabbix: HostPrototypeCreateParams.Host is required")
	}
	if len(p.GroupLinks) == 0 {
		return errors.New("zabbix: HostPrototypeCreateParams.GroupLinks must not be empty")
	}
	for _, value := range p.GroupLinks {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.GroupPrototypes {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Templates {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Interfaces {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Macros {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if err := requireWriteID(p.RuleID, "HostPrototypeCreateParams.RuleID"); err != nil {
		return err
	}
	return nil
}

// HostPrototypeUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostprototype/object
type HostPrototypeUpdateParams struct {
	HostID           string                          `json:"hostid"`
	Host             *string                         `json:"host,omitempty"`
	Name             *string                         `json:"name,omitempty"`
	Status           *int                            `json:"status,omitempty"`
	Discover         *int                            `json:"discover,omitempty"`
	CustomInterfaces *int                            `json:"custom_interfaces,omitempty"`
	InventoryMode    *int                            `json:"inventory_mode,omitempty"`
	UUID             *string                         `json:"uuid,omitempty"`
	GroupLinks       []HostGroupRef                  `json:"groupLinks,omitzero"`
	GroupPrototypes  []HostPrototypeGroupUpdateInput `json:"groupPrototypes,omitzero"`
	Templates        []TemplateRef                   `json:"templates,omitzero"`
	Interfaces       []HostPrototypeInterfaceInput   `json:"interfaces,omitzero"`
	Tags             []TagInput                      `json:"tags,omitzero"`
	Macros           []HostMacroUpdateInput          `json:"macros,omitzero"`
}

func (p HostPrototypeUpdateParams) validate() error {
	if err := requireWriteID(p.HostID, "HostPrototypeUpdateParams.HostID"); err != nil {
		return err
	}
	for _, value := range p.GroupLinks {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.GroupPrototypes {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Templates {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Interfaces {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Macros {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// UserGroupCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usergroup/object
type UserGroupCreateParams struct {
	Name                string                `json:"name"`
	GUIAccess           *int                  `json:"gui_access,omitempty"`
	UsersStatus         *int                  `json:"users_status,omitempty"`
	DebugMode           *int                  `json:"debug_mode,omitempty"`
	UserDirectoryID     *string               `json:"userdirectoryid,omitempty"`
	MFAStatus           *int                  `json:"mfa_status,omitempty"`
	MFAID               *string               `json:"mfaid,omitempty"`
	TagFilters          []UserGroupTagInput   `json:"tag_filters,omitzero"`
	Users               []UserRef             `json:"users,omitzero"`
	HostGroupRights     []UserGroupRightInput `json:"hostgroup_rights,omitzero"`
	TemplateGroupRights []UserGroupRightInput `json:"templategroup_rights,omitzero"`
}

func (p UserGroupCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: UserGroupCreateParams.Name is required")
	}
	for _, value := range p.TagFilters {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Users {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.HostGroupRights {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.TemplateGroupRights {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// UserGroupUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usergroup/object
type UserGroupUpdateParams struct {
	UsrGrpID            string                `json:"usrgrpid"`
	Name                *string               `json:"name,omitempty"`
	GUIAccess           *int                  `json:"gui_access,omitempty"`
	UsersStatus         *int                  `json:"users_status,omitempty"`
	DebugMode           *int                  `json:"debug_mode,omitempty"`
	UserDirectoryID     *string               `json:"userdirectoryid,omitempty"`
	MFAStatus           *int                  `json:"mfa_status,omitempty"`
	MFAID               *string               `json:"mfaid,omitempty"`
	TagFilters          []UserGroupTagInput   `json:"tag_filters,omitzero"`
	Users               []UserRef             `json:"users,omitzero"`
	HostGroupRights     []UserGroupRightInput `json:"hostgroup_rights,omitzero"`
	TemplateGroupRights []UserGroupRightInput `json:"templategroup_rights,omitzero"`
}

func (p UserGroupUpdateParams) validate() error {
	if err := requireWriteID(p.UsrGrpID, "UserGroupUpdateParams.UsrGrpID"); err != nil {
		return err
	}
	for _, value := range p.TagFilters {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Users {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.HostGroupRights {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.TemplateGroupRights {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// TokenCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/token/object
type TokenCreateParams struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	UserID      *string `json:"userid,omitempty"`
	Status      *int    `json:"status,omitempty"`
	ExpiresAt   *int64  `json:"expires_at,omitempty"`
}

func (p TokenCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: TokenCreateParams.Name is required")
	}
	return nil
}

// TokenUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/token/object
type TokenUpdateParams struct {
	TokenID     string  `json:"tokenid"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *int    `json:"status,omitempty"`
	ExpiresAt   *int64  `json:"expires_at,omitempty"`
}

func (p TokenUpdateParams) validate() error {
	if err := requireWriteID(p.TokenID, "TokenUpdateParams.TokenID"); err != nil {
		return err
	}
	return nil
}

// TemplateDashboardCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templatedashboard/object
type TemplateDashboardCreateParams struct {
	TemplateID    string                       `json:"templateid"`
	Name          string                       `json:"name"`
	DisplayPeriod *int                         `json:"display_period,omitempty"`
	AutoStart     *int                         `json:"auto_start,omitempty"`
	UUID          *string                      `json:"uuid,omitempty"`
	Pages         []TemplateDashboardPageInput `json:"pages"`
}

func (p TemplateDashboardCreateParams) validate() error {
	if err := requireWriteID(p.TemplateID, "TemplateDashboardCreateParams.TemplateID"); err != nil {
		return err
	}
	if p.Name == "" {
		return errors.New("zabbix: TemplateDashboardCreateParams.Name is required")
	}
	if len(p.Pages) == 0 {
		return errors.New("zabbix: TemplateDashboardCreateParams.Pages must not be empty")
	}
	for _, value := range p.Pages {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// TemplateDashboardUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templatedashboard/object
type TemplateDashboardUpdateParams struct {
	DashboardID   string                             `json:"dashboardid"`
	Name          *string                            `json:"name,omitempty"`
	DisplayPeriod *int                               `json:"display_period,omitempty"`
	AutoStart     *int                               `json:"auto_start,omitempty"`
	UUID          *string                            `json:"uuid,omitempty"`
	Pages         []TemplateDashboardPageUpdateInput `json:"pages,omitzero"`
}

func (p TemplateDashboardUpdateParams) validate() error {
	if err := requireWriteID(p.DashboardID, "TemplateDashboardUpdateParams.DashboardID"); err != nil {
		return err
	}
	for _, value := range p.Pages {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ServiceCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/service/object
type ServiceCreateParams struct {
	Name             string                   `json:"name"`
	Algorithm        int                      `json:"algorithm"`
	SortOrder        int                      `json:"sortorder"`
	Weight           *int                     `json:"weight,omitempty"`
	PropagationRule  *int                     `json:"propagation_rule,omitempty"`
	PropagationValue *int                     `json:"propagation_value,omitempty"`
	Description      *string                  `json:"description,omitempty"`
	UUID             *string                  `json:"uuid,omitempty"`
	StatusRules      []ServiceStatusRuleInput `json:"status_rules,omitzero"`
	Tags             []TagInput               `json:"tags,omitzero"`
	ProblemTags      []ServiceProblemTagInput `json:"problem_tags,omitzero"`
	Parents          []ServiceRef             `json:"parents,omitzero"`
	Children         []ServiceRef             `json:"children,omitzero"`
}

func (p ServiceCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: ServiceCreateParams.Name is required")
	}
	for _, value := range p.StatusRules {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.ProblemTags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Parents {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Children {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if (p.PropagationRule == nil) != (p.PropagationValue == nil) {
		return errors.New("zabbix: propagation_rule and propagation_value must be supplied together")
	}
	return nil
}

// ServiceUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/service/object
type ServiceUpdateParams struct {
	ServiceID        string                   `json:"serviceid"`
	Name             *string                  `json:"name,omitempty"`
	Algorithm        *int                     `json:"algorithm,omitempty"`
	SortOrder        *int                     `json:"sortorder,omitempty"`
	Weight           *int                     `json:"weight,omitempty"`
	PropagationRule  *int                     `json:"propagation_rule,omitempty"`
	PropagationValue *int                     `json:"propagation_value,omitempty"`
	Description      *string                  `json:"description,omitempty"`
	UUID             *string                  `json:"uuid,omitempty"`
	StatusRules      []ServiceStatusRuleInput `json:"status_rules,omitzero"`
	Tags             []TagInput               `json:"tags,omitzero"`
	ProblemTags      []ServiceProblemTagInput `json:"problem_tags,omitzero"`
	Parents          []ServiceRef             `json:"parents,omitzero"`
	Children         []ServiceRef             `json:"children,omitzero"`
}

func (p ServiceUpdateParams) validate() error {
	if err := requireWriteID(p.ServiceID, "ServiceUpdateParams.ServiceID"); err != nil {
		return err
	}
	for _, value := range p.StatusRules {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.ProblemTags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Parents {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Children {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if (p.PropagationRule == nil) != (p.PropagationValue == nil) {
		return errors.New("zabbix: propagation_rule and propagation_value must be supplied together")
	}
	return nil
}

// UserCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/user/object
type UserCreateParams struct {
	Username    string           `json:"username"`
	Passwd      *string          `json:"passwd,omitempty"`
	RoleID      *string          `json:"roleid,omitempty"`
	Autologin   *int             `json:"autologin,omitempty"`
	Autologout  *string          `json:"autologout,omitempty"`
	Lang        *string          `json:"lang,omitempty"`
	Name        *string          `json:"name,omitempty"`
	Refresh     *string          `json:"refresh,omitempty"`
	RowsPerPage *int             `json:"rows_per_page,omitempty"`
	Surname     *string          `json:"surname,omitempty"`
	Theme       *string          `json:"theme,omitempty"`
	URL         *string          `json:"url,omitempty"`
	Timezone    *string          `json:"timezone,omitempty"`
	UserGroups  []UserGroupRef   `json:"usrgrps,omitzero"`
	Medias      []UserMediaInput `json:"medias,omitzero"`
}

func (p UserCreateParams) validate() error {
	if p.Username == "" {
		return errors.New("zabbix: UserCreateParams.Username is required")
	}
	for _, value := range p.UserGroups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Medias {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// UserUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/user/object
type UserUpdateParams struct {
	UserID          string           `json:"userid"`
	Username        *string          `json:"username,omitempty"`
	Passwd          *string          `json:"passwd,omitempty"`
	RoleID          *string          `json:"roleid,omitempty"`
	Autologin       *int             `json:"autologin,omitempty"`
	Autologout      *string          `json:"autologout,omitempty"`
	Lang            *string          `json:"lang,omitempty"`
	Name            *string          `json:"name,omitempty"`
	Refresh         *string          `json:"refresh,omitempty"`
	RowsPerPage     *int             `json:"rows_per_page,omitempty"`
	Surname         *string          `json:"surname,omitempty"`
	Theme           *string          `json:"theme,omitempty"`
	URL             *string          `json:"url,omitempty"`
	Timezone        *string          `json:"timezone,omitempty"`
	UserGroups      []UserGroupRef   `json:"usrgrps,omitzero"`
	Medias          []UserMediaInput `json:"medias,omitzero"`
	CurrentPassword *string          `json:"current_passwd,omitempty"`
}

func (p UserUpdateParams) validate() error {
	if err := requireWriteID(p.UserID, "UserUpdateParams.UserID"); err != nil {
		return err
	}
	for _, value := range p.UserGroups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Medias {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// TaskCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/task/object
type TaskCreateParams struct {
	Type    int `json:"type"`
	Request any `json:"request"`
}

func (p TaskCreateParams) validate() error {
	if p.Request == nil {
		return errors.New("zabbix: TaskCreateParams.Request is required")
	}
	if p.Type != 1 && p.Type != 2 && p.Type != 6 {
		return errors.New("zabbix: task type must be 1, 2, or 6")
	}
	return nil
}

// ValueMapCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/valuemap/object
type ValueMapCreateParams struct {
	HostID   string              `json:"hostid"`
	Name     string              `json:"name"`
	UUID     *string             `json:"uuid,omitempty"`
	Mappings []ValueMappingInput `json:"mappings"`
}

func (p ValueMapCreateParams) validate() error {
	if err := requireWriteID(p.HostID, "ValueMapCreateParams.HostID"); err != nil {
		return err
	}
	if p.Name == "" {
		return errors.New("zabbix: ValueMapCreateParams.Name is required")
	}
	if len(p.Mappings) == 0 {
		return errors.New("zabbix: ValueMapCreateParams.Mappings must not be empty")
	}
	for _, value := range p.Mappings {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ValueMapUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/valuemap/object
type ValueMapUpdateParams struct {
	ValueMapID string              `json:"valuemapid"`
	Name       *string             `json:"name,omitempty"`
	UUID       *string             `json:"uuid,omitempty"`
	Mappings   []ValueMappingInput `json:"mappings,omitzero"`
}

func (p ValueMapUpdateParams) validate() error {
	if err := requireWriteID(p.ValueMapID, "ValueMapUpdateParams.ValueMapID"); err != nil {
		return err
	}
	for _, value := range p.Mappings {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// MFACreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mfa/object
type MFACreateParams struct {
	Type         *int    `json:"type,omitempty"`
	Name         string  `json:"name"`
	HashFunction *int    `json:"hash_function,omitempty"`
	CodeLength   *int    `json:"code_length,omitempty"`
	APIHostname  *string `json:"api_hostname,omitempty"`
	ClientID     *string `json:"clientid,omitempty"`
	ClientSecret *string `json:"client_secret,omitempty"`
}

func (p MFACreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: MFACreateParams.Name is required")
	}
	return nil
}

// MFAUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mfa/object
type MFAUpdateParams struct {
	MFAID        string  `json:"mfaid"`
	Type         *int    `json:"type,omitempty"`
	Name         *string `json:"name,omitempty"`
	HashFunction *int    `json:"hash_function,omitempty"`
	CodeLength   *int    `json:"code_length,omitempty"`
	APIHostname  *string `json:"api_hostname,omitempty"`
	ClientID     *string `json:"clientid,omitempty"`
	ClientSecret *string `json:"client_secret,omitempty"`
}

func (p MFAUpdateParams) validate() error {
	if err := requireWriteID(p.MFAID, "MFAUpdateParams.MFAID"); err != nil {
		return err
	}
	return nil
}

// AutoregistrationUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/autoregistration/object
type AutoregistrationUpdateParams struct {
	TLSAccept      *int    `json:"tls_accept,omitempty"`
	TLSPSKIdentity *string `json:"tls_psk_identity,omitempty"`
	TLSPSK         *string `json:"tls_psk,omitempty"`
}

func (p AutoregistrationUpdateParams) validate() error {
	return nil
}

// ActionCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/action/object
type ActionCreateParams struct {
	Name               string                         `json:"name"`
	EventSource        int                            `json:"eventsource"`
	Status             *int                           `json:"status,omitempty"`
	EscPeriod          *string                        `json:"esc_period,omitempty"`
	PauseSuppressed    *int                           `json:"pause_suppressed,omitempty"`
	PauseSymptoms      *int                           `json:"pause_symptoms,omitempty"`
	NotifyIfCanceled   *int                           `json:"notify_if_canceled,omitempty"`
	Operations         []ActionOperationInput         `json:"operations,omitzero"`
	RecoveryOperations []ActionRecoveryOperationInput `json:"recovery_operations,omitzero"`
	UpdateOperations   []ActionUpdateOperationInput   `json:"update_operations,omitzero"`
	Filter             *ActionFilterInput             `json:"filter,omitempty"`
}

func (p ActionCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: ActionCreateParams.Name is required")
	}
	for _, value := range p.Operations {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.RecoveryOperations {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.UpdateOperations {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if p.Filter != nil {
		if err := p.Filter.validate(); err != nil {
			return err
		}
	}
	if len(p.Operations) == 0 && len(p.RecoveryOperations) == 0 && len(p.UpdateOperations) == 0 {
		return errors.New("zabbix: action requires at least one operation")
	}
	return nil
}

// ActionUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/action/object
type ActionUpdateParams struct {
	ActionID           string                         `json:"actionid"`
	Name               *string                        `json:"name,omitempty"`
	Status             *int                           `json:"status,omitempty"`
	EscPeriod          *string                        `json:"esc_period,omitempty"`
	PauseSuppressed    *int                           `json:"pause_suppressed,omitempty"`
	PauseSymptoms      *int                           `json:"pause_symptoms,omitempty"`
	NotifyIfCanceled   *int                           `json:"notify_if_canceled,omitempty"`
	Operations         []ActionOperationInput         `json:"operations,omitzero"`
	RecoveryOperations []ActionRecoveryOperationInput `json:"recovery_operations,omitzero"`
	UpdateOperations   []ActionUpdateOperationInput   `json:"update_operations,omitzero"`
	Filter             *ActionFilterInput             `json:"filter,omitempty"`
}

func (p ActionUpdateParams) validate() error {
	if err := requireWriteID(p.ActionID, "ActionUpdateParams.ActionID"); err != nil {
		return err
	}
	for _, value := range p.Operations {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.RecoveryOperations {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.UpdateOperations {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if p.Filter != nil {
		if err := p.Filter.validate(); err != nil {
			return err
		}
	}
	return nil
}

// TemplateGroupCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templategroup/object
type TemplateGroupCreateParams struct {
	Name string `json:"name"`
}

func (p TemplateGroupCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: TemplateGroupCreateParams.Name is required")
	}
	return nil
}

// TemplateGroupUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templategroup/object
type TemplateGroupUpdateParams struct {
	GroupID string  `json:"groupid"`
	Name    *string `json:"name,omitempty"`
}

func (p TemplateGroupUpdateParams) validate() error {
	if err := requireWriteID(p.GroupID, "TemplateGroupUpdateParams.GroupID"); err != nil {
		return err
	}
	return nil
}

// DiscoveryRuleCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/discoveryrule/object
type DiscoveryRuleCreateParams struct {
	HostID              string                        `json:"hostid"`
	Name                string                        `json:"name"`
	Key                 string                        `json:"key_"`
	Type                int                           `json:"type"`
	Delay               *string                       `json:"delay,omitempty"`
	Status              *int                          `json:"status,omitempty"`
	Lifetime            *string                       `json:"lifetime,omitempty"`
	LifetimeType        *int                          `json:"lifetime_type,omitempty"`
	EnabledLifetime     *string                       `json:"enabled_lifetime,omitempty"`
	EnabledLifetimeType *int                          `json:"enabled_lifetime_type,omitempty"`
	Description         *string                       `json:"description,omitempty"`
	InterfaceID         *string                       `json:"interfaceid,omitempty"`
	Params              *string                       `json:"params,omitempty"`
	SNMPOID             *string                       `json:"snmp_oid,omitempty"`
	MasterItemID        *string                       `json:"master_itemid,omitempty"`
	URL                 *string                       `json:"url,omitempty"`
	IPMISensor          *string                       `json:"ipmi_sensor,omitempty"`
	JMXEndpoint         *string                       `json:"jmx_endpoint,omitempty"`
	Timeout             *string                       `json:"timeout,omitempty"`
	TrapperHosts        *string                       `json:"trapper_hosts,omitempty"`
	Username            *string                       `json:"username,omitempty"`
	Password            *string                       `json:"password,omitempty"`
	AuthType            *int                          `json:"authtype,omitempty"`
	Filter              *DiscoveryFilterInput         `json:"filter,omitempty"`
	LLDMacroPaths       []LLDMacroPathInput           `json:"lld_macro_paths,omitzero"`
	Preprocessing       []DiscoveryPreprocessingInput `json:"preprocessing,omitzero"`
}

func (p DiscoveryRuleCreateParams) validate() error {
	if err := requireWriteID(p.HostID, "DiscoveryRuleCreateParams.HostID"); err != nil {
		return err
	}
	if p.Name == "" {
		return errors.New("zabbix: DiscoveryRuleCreateParams.Name is required")
	}
	if p.Key == "" {
		return errors.New("zabbix: DiscoveryRuleCreateParams.Key is required")
	}
	if p.Filter != nil {
		if err := p.Filter.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.LLDMacroPaths {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Preprocessing {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// DiscoveryRuleUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/discoveryrule/object
type DiscoveryRuleUpdateParams struct {
	ItemID              string                        `json:"itemid"`
	Name                *string                       `json:"name,omitempty"`
	Key                 *string                       `json:"key_,omitempty"`
	Type                *int                          `json:"type,omitempty"`
	Delay               *string                       `json:"delay,omitempty"`
	Status              *int                          `json:"status,omitempty"`
	Lifetime            *string                       `json:"lifetime,omitempty"`
	LifetimeType        *int                          `json:"lifetime_type,omitempty"`
	EnabledLifetime     *string                       `json:"enabled_lifetime,omitempty"`
	EnabledLifetimeType *int                          `json:"enabled_lifetime_type,omitempty"`
	Description         *string                       `json:"description,omitempty"`
	InterfaceID         *string                       `json:"interfaceid,omitempty"`
	Params              *string                       `json:"params,omitempty"`
	SNMPOID             *string                       `json:"snmp_oid,omitempty"`
	MasterItemID        *string                       `json:"master_itemid,omitempty"`
	URL                 *string                       `json:"url,omitempty"`
	IPMISensor          *string                       `json:"ipmi_sensor,omitempty"`
	JMXEndpoint         *string                       `json:"jmx_endpoint,omitempty"`
	Timeout             *string                       `json:"timeout,omitempty"`
	TrapperHosts        *string                       `json:"trapper_hosts,omitempty"`
	Username            *string                       `json:"username,omitempty"`
	Password            *string                       `json:"password,omitempty"`
	AuthType            *int                          `json:"authtype,omitempty"`
	Filter              *DiscoveryFilterInput         `json:"filter,omitempty"`
	LLDMacroPaths       []LLDMacroPathInput           `json:"lld_macro_paths,omitzero"`
	Preprocessing       []DiscoveryPreprocessingInput `json:"preprocessing,omitzero"`
}

func (p DiscoveryRuleUpdateParams) validate() error {
	if err := requireWriteID(p.ItemID, "DiscoveryRuleUpdateParams.ItemID"); err != nil {
		return err
	}
	if p.Filter != nil {
		if err := p.Filter.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.LLDMacroPaths {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Preprocessing {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// GraphPrototypeCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graphprototype/object
type GraphPrototypeCreateParams struct {
	Name           string           `json:"name"`
	Width          int              `json:"width"`
	Height         int              `json:"height"`
	YAxisMin       *float64         `json:"yaxismin,omitempty"`
	YAxisMax       *float64         `json:"yaxismax,omitempty"`
	ShowWorkPeriod *int             `json:"show_work_period,omitempty"`
	ShowTriggers   *int             `json:"show_triggers,omitempty"`
	GraphType      *int             `json:"graphtype,omitempty"`
	ShowLegend     *int             `json:"show_legend,omitempty"`
	Show3D         *int             `json:"show_3d,omitempty"`
	PercentLeft    *float64         `json:"percent_left,omitempty"`
	PercentRight   *float64         `json:"percent_right,omitempty"`
	YMinType       *int             `json:"ymin_type,omitempty"`
	YMaxType       *int             `json:"ymax_type,omitempty"`
	YMinItemID     *string          `json:"ymin_itemid,omitempty"`
	YMaxItemID     *string          `json:"ymax_itemid,omitempty"`
	Discover       *int             `json:"discover,omitempty"`
	UUID           *string          `json:"uuid,omitempty"`
	GItems         []GraphItemInput `json:"gitems"`
}

func (p GraphPrototypeCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: GraphPrototypeCreateParams.Name is required")
	}
	if len(p.GItems) == 0 {
		return errors.New("zabbix: GraphPrototypeCreateParams.GItems must not be empty")
	}
	for _, value := range p.GItems {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// GraphPrototypeUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graphprototype/object
type GraphPrototypeUpdateParams struct {
	GraphID        string                 `json:"graphid"`
	Name           *string                `json:"name,omitempty"`
	Width          *int                   `json:"width,omitempty"`
	Height         *int                   `json:"height,omitempty"`
	YAxisMin       *float64               `json:"yaxismin,omitempty"`
	YAxisMax       *float64               `json:"yaxismax,omitempty"`
	ShowWorkPeriod *int                   `json:"show_work_period,omitempty"`
	ShowTriggers   *int                   `json:"show_triggers,omitempty"`
	GraphType      *int                   `json:"graphtype,omitempty"`
	ShowLegend     *int                   `json:"show_legend,omitempty"`
	Show3D         *int                   `json:"show_3d,omitempty"`
	PercentLeft    *float64               `json:"percent_left,omitempty"`
	PercentRight   *float64               `json:"percent_right,omitempty"`
	YMinType       *int                   `json:"ymin_type,omitempty"`
	YMaxType       *int                   `json:"ymax_type,omitempty"`
	YMinItemID     *string                `json:"ymin_itemid,omitempty"`
	YMaxItemID     *string                `json:"ymax_itemid,omitempty"`
	Discover       *int                   `json:"discover,omitempty"`
	UUID           *string                `json:"uuid,omitempty"`
	GItems         []GraphItemUpdateInput `json:"gitems,omitzero"`
}

func (p GraphPrototypeUpdateParams) validate() error {
	if err := requireWriteID(p.GraphID, "GraphPrototypeUpdateParams.GraphID"); err != nil {
		return err
	}
	for _, value := range p.GItems {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// RoleCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/role/object
type RoleCreateParams struct {
	Name  string          `json:"name"`
	Type  int             `json:"type"`
	Rules *RoleRulesInput `json:"rules,omitempty"`
}

func (p RoleCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: RoleCreateParams.Name is required")
	}
	if p.Rules != nil {
		if err := p.Rules.validate(); err != nil {
			return err
		}
	}
	return nil
}

// RoleUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/role/object
type RoleUpdateParams struct {
	RoleID string          `json:"roleid"`
	Name   *string         `json:"name,omitempty"`
	Type   *int            `json:"type,omitempty"`
	Rules  *RoleRulesInput `json:"rules,omitempty"`
}

func (p RoleUpdateParams) validate() error {
	if err := requireWriteID(p.RoleID, "RoleUpdateParams.RoleID"); err != nil {
		return err
	}
	if p.Rules != nil {
		if err := p.Rules.validate(); err != nil {
			return err
		}
	}
	return nil
}

// HostMacroCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/object
type HostMacroCreateParams struct {
	Macro       string  `json:"macro"`
	Value       string  `json:"value"`
	Type        *int    `json:"type,omitempty"`
	Description *string `json:"description,omitempty"`
	HostID      string  `json:"hostid"`
}

func (p HostMacroCreateParams) validate() error {
	if p.Macro == "" {
		return errors.New("zabbix: HostMacroCreateParams.Macro is required")
	}
	if err := requireWriteID(p.HostID, "HostMacroCreateParams.HostID"); err != nil {
		return err
	}
	return nil
}

// HostMacroUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/object
type HostMacroUpdateParams struct {
	HostMacroID string  `json:"hostmacroid"`
	Macro       *string `json:"macro,omitempty"`
	Value       *string `json:"value,omitempty"`
	Type        *int    `json:"type,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (p HostMacroUpdateParams) validate() error {
	if err := requireWriteID(p.HostMacroID, "HostMacroUpdateParams.HostMacroID"); err != nil {
		return err
	}
	return nil
}

// GlobalMacroCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/object
type GlobalMacroCreateParams struct {
	Macro       string  `json:"macro"`
	Value       string  `json:"value"`
	Type        *int    `json:"type,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (p GlobalMacroCreateParams) validate() error {
	if p.Macro == "" {
		return errors.New("zabbix: GlobalMacroCreateParams.Macro is required")
	}
	return nil
}

// GlobalMacroUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/object
type GlobalMacroUpdateParams struct {
	GlobalMacroID string  `json:"globalmacroid"`
	Macro         *string `json:"macro,omitempty"`
	Value         *string `json:"value,omitempty"`
	Type          *int    `json:"type,omitempty"`
	Description   *string `json:"description,omitempty"`
}

func (p GlobalMacroUpdateParams) validate() error {
	if err := requireWriteID(p.GlobalMacroID, "GlobalMacroUpdateParams.GlobalMacroID"); err != nil {
		return err
	}
	return nil
}

// HostGroupCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostgroup/object
type HostGroupCreateParams struct {
	Name string  `json:"name"`
	UUID *string `json:"uuid,omitempty"`
}

func (p HostGroupCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: HostGroupCreateParams.Name is required")
	}
	return nil
}

// HostGroupUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostgroup/object
type HostGroupUpdateParams struct {
	GroupID string  `json:"groupid"`
	Name    *string `json:"name,omitempty"`
	UUID    *string `json:"uuid,omitempty"`
}

func (p HostGroupUpdateParams) validate() error {
	if err := requireWriteID(p.GroupID, "HostGroupUpdateParams.GroupID"); err != nil {
		return err
	}
	return nil
}

// TriggerPrototypeCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/triggerprototype/object
type TriggerPrototypeCreateParams struct {
	Description        string     `json:"description"`
	Expression         string     `json:"expression"`
	EventName          *string    `json:"event_name,omitempty"`
	OpData             *string    `json:"opdata,omitempty"`
	Comments           *string    `json:"comments,omitempty"`
	Priority           *int       `json:"priority,omitempty"`
	Status             *int       `json:"status,omitempty"`
	Type               *int       `json:"type,omitempty"`
	RecoveryMode       *int       `json:"recovery_mode,omitempty"`
	RecoveryExpression *string    `json:"recovery_expression,omitempty"`
	CorrelationMode    *int       `json:"correlation_mode,omitempty"`
	CorrelationTag     *string    `json:"correlation_tag,omitempty"`
	ManualClose        *int       `json:"manual_close,omitempty"`
	URL                *string    `json:"url,omitempty"`
	URLName            *string    `json:"url_name,omitempty"`
	UUID               *string    `json:"uuid,omitempty"`
	Discover           *int       `json:"discover,omitempty"`
	Tags               []TagInput `json:"tags,omitzero"`
}

func (p TriggerPrototypeCreateParams) validate() error {
	if p.Description == "" {
		return errors.New("zabbix: TriggerPrototypeCreateParams.Description is required")
	}
	if p.Expression == "" {
		return errors.New("zabbix: TriggerPrototypeCreateParams.Expression is required")
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// TriggerPrototypeUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/triggerprototype/object
type TriggerPrototypeUpdateParams struct {
	TriggerID          string     `json:"triggerid"`
	Description        *string    `json:"description,omitempty"`
	Expression         *string    `json:"expression,omitempty"`
	EventName          *string    `json:"event_name,omitempty"`
	OpData             *string    `json:"opdata,omitempty"`
	Comments           *string    `json:"comments,omitempty"`
	Priority           *int       `json:"priority,omitempty"`
	Status             *int       `json:"status,omitempty"`
	Type               *int       `json:"type,omitempty"`
	RecoveryMode       *int       `json:"recovery_mode,omitempty"`
	RecoveryExpression *string    `json:"recovery_expression,omitempty"`
	CorrelationMode    *int       `json:"correlation_mode,omitempty"`
	CorrelationTag     *string    `json:"correlation_tag,omitempty"`
	ManualClose        *int       `json:"manual_close,omitempty"`
	URL                *string    `json:"url,omitempty"`
	URLName            *string    `json:"url_name,omitempty"`
	UUID               *string    `json:"uuid,omitempty"`
	Discover           *int       `json:"discover,omitempty"`
	Tags               []TagInput `json:"tags,omitzero"`
}

func (p TriggerPrototypeUpdateParams) validate() error {
	if err := requireWriteID(p.TriggerID, "TriggerPrototypeUpdateParams.TriggerID"); err != nil {
		return err
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// GraphCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graph/object
type GraphCreateParams struct {
	Name           string           `json:"name"`
	Width          int              `json:"width"`
	Height         int              `json:"height"`
	YAxisMin       *float64         `json:"yaxismin,omitempty"`
	YAxisMax       *float64         `json:"yaxismax,omitempty"`
	ShowWorkPeriod *int             `json:"show_work_period,omitempty"`
	ShowTriggers   *int             `json:"show_triggers,omitempty"`
	GraphType      *int             `json:"graphtype,omitempty"`
	ShowLegend     *int             `json:"show_legend,omitempty"`
	Show3D         *int             `json:"show_3d,omitempty"`
	PercentLeft    *float64         `json:"percent_left,omitempty"`
	PercentRight   *float64         `json:"percent_right,omitempty"`
	YMinType       *int             `json:"ymin_type,omitempty"`
	YMaxType       *int             `json:"ymax_type,omitempty"`
	YMinItemID     *string          `json:"ymin_itemid,omitempty"`
	YMaxItemID     *string          `json:"ymax_itemid,omitempty"`
	UUID           *string          `json:"uuid,omitempty"`
	GItems         []GraphItemInput `json:"gitems"`
}

func (p GraphCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: GraphCreateParams.Name is required")
	}
	if len(p.GItems) == 0 {
		return errors.New("zabbix: GraphCreateParams.GItems must not be empty")
	}
	for _, value := range p.GItems {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// GraphUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graph/object
type GraphUpdateParams struct {
	GraphID        string                 `json:"graphid"`
	Name           *string                `json:"name,omitempty"`
	Width          *int                   `json:"width,omitempty"`
	Height         *int                   `json:"height,omitempty"`
	YAxisMin       *float64               `json:"yaxismin,omitempty"`
	YAxisMax       *float64               `json:"yaxismax,omitempty"`
	ShowWorkPeriod *int                   `json:"show_work_period,omitempty"`
	ShowTriggers   *int                   `json:"show_triggers,omitempty"`
	GraphType      *int                   `json:"graphtype,omitempty"`
	ShowLegend     *int                   `json:"show_legend,omitempty"`
	Show3D         *int                   `json:"show_3d,omitempty"`
	PercentLeft    *float64               `json:"percent_left,omitempty"`
	PercentRight   *float64               `json:"percent_right,omitempty"`
	YMinType       *int                   `json:"ymin_type,omitempty"`
	YMaxType       *int                   `json:"ymax_type,omitempty"`
	YMinItemID     *string                `json:"ymin_itemid,omitempty"`
	YMaxItemID     *string                `json:"ymax_itemid,omitempty"`
	UUID           *string                `json:"uuid,omitempty"`
	GItems         []GraphItemUpdateInput `json:"gitems,omitzero"`
}

func (p GraphUpdateParams) validate() error {
	if err := requireWriteID(p.GraphID, "GraphUpdateParams.GraphID"); err != nil {
		return err
	}
	for _, value := range p.GItems {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ProxyCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxy/object
type ProxyCreateParams struct {
	Name                 string    `json:"name"`
	ProxyGroupID         *string   `json:"proxy_groupid,omitempty"`
	LocalAddress         *string   `json:"local_address,omitempty"`
	LocalPort            *string   `json:"local_port,omitempty"`
	OperatingMode        int       `json:"operating_mode"`
	AllowedAddresses     *string   `json:"allowed_addresses,omitempty"`
	Address              *string   `json:"address,omitempty"`
	Port                 *string   `json:"port,omitempty"`
	Description          *string   `json:"description,omitempty"`
	TLSConnect           *int      `json:"tls_connect,omitempty"`
	TLSAccept            *int      `json:"tls_accept,omitempty"`
	TLSIssuer            *string   `json:"tls_issuer,omitempty"`
	TLSSubject           *string   `json:"tls_subject,omitempty"`
	CustomTimeouts       *int      `json:"custom_timeouts,omitempty"`
	TimeoutZabbixAgent   *string   `json:"timeout_zabbix_agent,omitempty"`
	TimeoutSimpleCheck   *string   `json:"timeout_simple_check,omitempty"`
	TimeoutSNMPAgent     *string   `json:"timeout_snmp_agent,omitempty"`
	TimeoutExternalCheck *string   `json:"timeout_external_check,omitempty"`
	TimeoutDBMonitor     *string   `json:"timeout_db_monitor,omitempty"`
	TimeoutHTTPAgent     *string   `json:"timeout_http_agent,omitempty"`
	TimeoutSSHAgent      *string   `json:"timeout_ssh_agent,omitempty"`
	TimeoutTelnetAgent   *string   `json:"timeout_telnet_agent,omitempty"`
	TimeoutScript        *string   `json:"timeout_script,omitempty"`
	TimeoutBrowser       *string   `json:"timeout_browser,omitempty"`
	Hosts                []HostRef `json:"hosts,omitzero"`
	TLSPSKIdentity       *string   `json:"tls_psk_identity,omitempty"`
	TLSPSK               *string   `json:"tls_psk,omitempty"`
}

func (p ProxyCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: ProxyCreateParams.Name is required")
	}
	for _, value := range p.Hosts {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ProxyUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxy/object
type ProxyUpdateParams struct {
	ProxyID              string    `json:"proxyid"`
	Name                 *string   `json:"name,omitempty"`
	ProxyGroupID         *string   `json:"proxy_groupid,omitempty"`
	LocalAddress         *string   `json:"local_address,omitempty"`
	LocalPort            *string   `json:"local_port,omitempty"`
	OperatingMode        *int      `json:"operating_mode,omitempty"`
	AllowedAddresses     *string   `json:"allowed_addresses,omitempty"`
	Address              *string   `json:"address,omitempty"`
	Port                 *string   `json:"port,omitempty"`
	Description          *string   `json:"description,omitempty"`
	TLSConnect           *int      `json:"tls_connect,omitempty"`
	TLSAccept            *int      `json:"tls_accept,omitempty"`
	TLSIssuer            *string   `json:"tls_issuer,omitempty"`
	TLSSubject           *string   `json:"tls_subject,omitempty"`
	CustomTimeouts       *int      `json:"custom_timeouts,omitempty"`
	TimeoutZabbixAgent   *string   `json:"timeout_zabbix_agent,omitempty"`
	TimeoutSimpleCheck   *string   `json:"timeout_simple_check,omitempty"`
	TimeoutSNMPAgent     *string   `json:"timeout_snmp_agent,omitempty"`
	TimeoutExternalCheck *string   `json:"timeout_external_check,omitempty"`
	TimeoutDBMonitor     *string   `json:"timeout_db_monitor,omitempty"`
	TimeoutHTTPAgent     *string   `json:"timeout_http_agent,omitempty"`
	TimeoutSSHAgent      *string   `json:"timeout_ssh_agent,omitempty"`
	TimeoutTelnetAgent   *string   `json:"timeout_telnet_agent,omitempty"`
	TimeoutScript        *string   `json:"timeout_script,omitempty"`
	TimeoutBrowser       *string   `json:"timeout_browser,omitempty"`
	Hosts                []HostRef `json:"hosts,omitzero"`
	TLSPSKIdentity       *string   `json:"tls_psk_identity,omitempty"`
	TLSPSK               *string   `json:"tls_psk,omitempty"`
}

func (p ProxyUpdateParams) validate() error {
	if err := requireWriteID(p.ProxyID, "ProxyUpdateParams.ProxyID"); err != nil {
		return err
	}
	for _, value := range p.Hosts {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// MaintenanceCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/maintenance/object
type MaintenanceCreateParams struct {
	Name            string                   `json:"name"`
	ActiveSince     int64                    `json:"active_since"`
	ActiveTill      int64                    `json:"active_till"`
	Description     *string                  `json:"description,omitempty"`
	MaintenanceType *int                     `json:"maintenance_type,omitempty"`
	TagsEvalType    *int                     `json:"tags_evaltype,omitempty"`
	Timeperiods     []MaintenancePeriodInput `json:"timeperiods"`
	Tags            []MaintenanceTagInput    `json:"tags,omitzero"`
	Groups          []HostGroupRef           `json:"groups,omitzero"`
	Hosts           []HostRef                `json:"hosts,omitzero"`
}

func (p MaintenanceCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: MaintenanceCreateParams.Name is required")
	}
	if len(p.Timeperiods) == 0 {
		return errors.New("zabbix: MaintenanceCreateParams.Timeperiods must not be empty")
	}
	for _, value := range p.Timeperiods {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Groups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Hosts {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if len(p.Groups) == 0 && len(p.Hosts) == 0 {
		return errors.New("zabbix: maintenance requires a host or host group")
	}
	return nil
}

// MaintenanceUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/maintenance/object
type MaintenanceUpdateParams struct {
	MaintenanceID   string                   `json:"maintenanceid"`
	Name            *string                  `json:"name,omitempty"`
	ActiveSince     *int64                   `json:"active_since,omitempty"`
	ActiveTill      *int64                   `json:"active_till,omitempty"`
	Description     *string                  `json:"description,omitempty"`
	MaintenanceType *int                     `json:"maintenance_type,omitempty"`
	TagsEvalType    *int                     `json:"tags_evaltype,omitempty"`
	Timeperiods     []MaintenancePeriodInput `json:"timeperiods,omitzero"`
	Tags            []MaintenanceTagInput    `json:"tags,omitzero"`
	Groups          []HostGroupRef           `json:"groups,omitzero"`
	Hosts           []HostRef                `json:"hosts,omitzero"`
}

func (p MaintenanceUpdateParams) validate() error {
	if err := requireWriteID(p.MaintenanceID, "MaintenanceUpdateParams.MaintenanceID"); err != nil {
		return err
	}
	for _, value := range p.Timeperiods {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Groups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Hosts {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// UserDirectoryCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/userdirectory/object
type UserDirectoryCreateParams struct {
	Name              *string                  `json:"name,omitempty"`
	IDPType           int                      `json:"idp_type"`
	ProvisionStatus   *int                     `json:"provision_status,omitempty"`
	Description       *string                  `json:"description,omitempty"`
	Host              *string                  `json:"host,omitempty"`
	Port              *int                     `json:"port,omitempty"`
	BaseDN            *string                  `json:"base_dn,omitempty"`
	SearchAttribute   *string                  `json:"search_attribute,omitempty"`
	BindDN            *string                  `json:"bind_dn,omitempty"`
	StartTLS          *int                     `json:"start_tls,omitempty"`
	SearchFilter      *string                  `json:"search_filter,omitempty"`
	GroupBaseDN       *string                  `json:"group_basedn,omitempty"`
	UserUsername      *string                  `json:"user_username,omitempty"`
	UserLastname      *string                  `json:"user_lastname,omitempty"`
	IDPEntityID       *string                  `json:"idp_entityid,omitempty"`
	SSOURL            *string                  `json:"sso_url,omitempty"`
	SLOURL            *string                  `json:"slo_url,omitempty"`
	UsernameAttribute *string                  `json:"username_attribute,omitempty"`
	SPEntityID        *string                  `json:"sp_entityid,omitempty"`
	BindPassword      *string                  `json:"bind_password,omitempty"`
	GroupName         *string                  `json:"group_name,omitempty"`
	ProvisionGroups   []ProvisioningGroupInput `json:"provision_groups,omitzero"`
	ProvisionMedia    []ProvisioningMediaInput `json:"provision_media,omitzero"`
}

func (p UserDirectoryCreateParams) validate() error {
	for _, value := range p.ProvisionGroups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.ProvisionMedia {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// UserDirectoryUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/userdirectory/object
type UserDirectoryUpdateParams struct {
	UserDirectoryID   string                   `json:"userdirectoryid"`
	Name              *string                  `json:"name,omitempty"`
	IDPType           *int                     `json:"idp_type,omitempty"`
	ProvisionStatus   *int                     `json:"provision_status,omitempty"`
	Description       *string                  `json:"description,omitempty"`
	Host              *string                  `json:"host,omitempty"`
	Port              *int                     `json:"port,omitempty"`
	BaseDN            *string                  `json:"base_dn,omitempty"`
	SearchAttribute   *string                  `json:"search_attribute,omitempty"`
	BindDN            *string                  `json:"bind_dn,omitempty"`
	StartTLS          *int                     `json:"start_tls,omitempty"`
	SearchFilter      *string                  `json:"search_filter,omitempty"`
	GroupBaseDN       *string                  `json:"group_basedn,omitempty"`
	UserUsername      *string                  `json:"user_username,omitempty"`
	UserLastname      *string                  `json:"user_lastname,omitempty"`
	IDPEntityID       *string                  `json:"idp_entityid,omitempty"`
	SSOURL            *string                  `json:"sso_url,omitempty"`
	SLOURL            *string                  `json:"slo_url,omitempty"`
	UsernameAttribute *string                  `json:"username_attribute,omitempty"`
	SPEntityID        *string                  `json:"sp_entityid,omitempty"`
	BindPassword      *string                  `json:"bind_password,omitempty"`
	GroupName         *string                  `json:"group_name,omitempty"`
	ProvisionGroups   []ProvisioningGroupInput `json:"provision_groups,omitzero"`
	ProvisionMedia    []ProvisioningMediaInput `json:"provision_media,omitzero"`
}

func (p UserDirectoryUpdateParams) validate() error {
	if err := requireWriteID(p.UserDirectoryID, "UserDirectoryUpdateParams.UserDirectoryID"); err != nil {
		return err
	}
	for _, value := range p.ProvisionGroups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.ProvisionMedia {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// HostCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/host/object
type HostCreateParams struct {
	Host           string               `json:"host"`
	Name           *string              `json:"name,omitempty"`
	Description    *string              `json:"description,omitempty"`
	Status         *int                 `json:"status,omitempty"`
	InventoryMode  *int                 `json:"inventory_mode,omitempty"`
	MonitoredBy    *int                 `json:"monitored_by,omitempty"`
	ProxyID        *string              `json:"proxyid,omitempty"`
	ProxyGroupID   *string              `json:"proxy_groupid,omitempty"`
	TLSConnect     *int                 `json:"tls_connect,omitempty"`
	TLSAccept      *int                 `json:"tls_accept,omitempty"`
	TLSIssuer      *string              `json:"tls_issuer,omitempty"`
	TLSSubject     *string              `json:"tls_subject,omitempty"`
	IPMIAuthType   *int                 `json:"ipmi_authtype,omitempty"`
	IPMIPrivilege  *int                 `json:"ipmi_privilege,omitempty"`
	IPMIUsername   *string              `json:"ipmi_username,omitempty"`
	IPMIPassword   *string              `json:"ipmi_password,omitempty"`
	Interfaces     []HostInterfaceInput `json:"interfaces,omitzero"`
	Tags           []TagInput           `json:"tags,omitzero"`
	Macros         []HostMacroInput     `json:"macros,omitzero"`
	Groups         []HostGroupRef       `json:"groups"`
	Templates      []TemplateRef        `json:"templates,omitzero"`
	TLSPSKIdentity *string              `json:"tls_psk_identity,omitempty"`
	TLSPSK         *string              `json:"tls_psk,omitempty"`
}

func (p HostCreateParams) validate() error {
	if p.Host == "" {
		return errors.New("zabbix: HostCreateParams.Host is required")
	}
	for _, value := range p.Interfaces {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Macros {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if len(p.Groups) == 0 {
		return errors.New("zabbix: HostCreateParams.Groups must not be empty")
	}
	for _, value := range p.Groups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Templates {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if p.MonitoredBy != nil && *p.MonitoredBy == 1 && (p.ProxyID == nil || *p.ProxyID == "") {
		return errors.New("zabbix: proxyid is required when monitored_by is proxy")
	}
	if p.MonitoredBy != nil && *p.MonitoredBy == 2 && (p.ProxyGroupID == nil || *p.ProxyGroupID == "") {
		return errors.New("zabbix: proxy_groupid is required when monitored_by is proxy group")
	}
	return nil
}

// HostUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/host/object
type HostUpdateParams struct {
	HostID         string                     `json:"hostid"`
	Host           *string                    `json:"host,omitempty"`
	Name           *string                    `json:"name,omitempty"`
	Description    *string                    `json:"description,omitempty"`
	Status         *int                       `json:"status,omitempty"`
	InventoryMode  *int                       `json:"inventory_mode,omitempty"`
	MonitoredBy    *int                       `json:"monitored_by,omitempty"`
	ProxyID        *string                    `json:"proxyid,omitempty"`
	ProxyGroupID   *string                    `json:"proxy_groupid,omitempty"`
	TLSConnect     *int                       `json:"tls_connect,omitempty"`
	TLSAccept      *int                       `json:"tls_accept,omitempty"`
	TLSIssuer      *string                    `json:"tls_issuer,omitempty"`
	TLSSubject     *string                    `json:"tls_subject,omitempty"`
	IPMIAuthType   *int                       `json:"ipmi_authtype,omitempty"`
	IPMIPrivilege  *int                       `json:"ipmi_privilege,omitempty"`
	IPMIUsername   *string                    `json:"ipmi_username,omitempty"`
	IPMIPassword   *string                    `json:"ipmi_password,omitempty"`
	Interfaces     []HostInterfaceUpdateInput `json:"interfaces,omitzero"`
	Tags           []TagInput                 `json:"tags,omitzero"`
	Macros         []HostMacroUpdateInput     `json:"macros,omitzero"`
	Groups         []HostGroupRef             `json:"groups,omitzero"`
	Templates      []TemplateRef              `json:"templates,omitzero"`
	TLSPSKIdentity *string                    `json:"tls_psk_identity,omitempty"`
	TLSPSK         *string                    `json:"tls_psk,omitempty"`
	TemplatesClear []TemplateRef              `json:"templates_clear,omitzero"`
}

func (p HostUpdateParams) validate() error {
	if err := requireWriteID(p.HostID, "HostUpdateParams.HostID"); err != nil {
		return err
	}
	for _, value := range p.Interfaces {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Macros {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Groups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Templates {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.TemplatesClear {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if p.MonitoredBy != nil && *p.MonitoredBy == 1 && (p.ProxyID == nil || *p.ProxyID == "") {
		return errors.New("zabbix: proxyid is required when monitored_by is proxy")
	}
	if p.MonitoredBy != nil && *p.MonitoredBy == 2 && (p.ProxyGroupID == nil || *p.ProxyGroupID == "") {
		return errors.New("zabbix: proxy_groupid is required when monitored_by is proxy group")
	}
	return nil
}

// TemplateCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/template/object
type TemplateCreateParams struct {
	Host        string             `json:"host"`
	Name        *string            `json:"name,omitempty"`
	Description *string            `json:"description,omitempty"`
	UUID        *string            `json:"uuid,omitempty"`
	Vendor      *string            `json:"vendor_name,omitempty"`
	Version     *string            `json:"vendor_version,omitempty"`
	Groups      []TemplateGroupRef `json:"groups"`
	Tags        []TagInput         `json:"tags,omitzero"`
	Macros      []HostMacroInput   `json:"macros,omitzero"`
	Templates   []TemplateRef      `json:"templates,omitzero"`
}

func (p TemplateCreateParams) validate() error {
	if p.Host == "" {
		return errors.New("zabbix: TemplateCreateParams.Host is required")
	}
	if len(p.Groups) == 0 {
		return errors.New("zabbix: TemplateCreateParams.Groups must not be empty")
	}
	for _, value := range p.Groups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Macros {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Templates {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// TemplateUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/template/object
type TemplateUpdateParams struct {
	TemplateID     string                 `json:"templateid"`
	Host           *string                `json:"host,omitempty"`
	Name           *string                `json:"name,omitempty"`
	Description    *string                `json:"description,omitempty"`
	UUID           *string                `json:"uuid,omitempty"`
	Vendor         *string                `json:"vendor_name,omitempty"`
	Version        *string                `json:"vendor_version,omitempty"`
	Groups         []TemplateGroupRef     `json:"groups,omitzero"`
	Tags           []TagInput             `json:"tags,omitzero"`
	Macros         []HostMacroUpdateInput `json:"macros,omitzero"`
	Templates      []TemplateRef          `json:"templates,omitzero"`
	TemplatesClear []TemplateRef          `json:"templates_clear,omitzero"`
}

func (p TemplateUpdateParams) validate() error {
	if err := requireWriteID(p.TemplateID, "TemplateUpdateParams.TemplateID"); err != nil {
		return err
	}
	for _, value := range p.Groups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Macros {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Templates {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.TemplatesClear {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ConnectorCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/connector/object
type ConnectorCreateParams struct {
	Name            string              `json:"name"`
	Protocol        *int                `json:"protocol,omitempty"`
	DataType        *int                `json:"data_type,omitempty"`
	URL             string              `json:"url"`
	MaxRecords      *int                `json:"max_records,omitempty"`
	MaxSenders      *int                `json:"max_senders,omitempty"`
	MaxAttempts     *int                `json:"max_attempts,omitempty"`
	AttemptInterval *string             `json:"attempt_interval,omitempty"`
	Timeout         *string             `json:"timeout,omitempty"`
	HTTPProxy       *string             `json:"http_proxy,omitempty"`
	AuthType        *int                `json:"authtype,omitempty"`
	Username        *string             `json:"username,omitempty"`
	Password        *string             `json:"password,omitempty"`
	Token           *string             `json:"token,omitempty"`
	VerifyPeer      *int                `json:"verify_peer,omitempty"`
	VerifyHost      *int                `json:"verify_host,omitempty"`
	SSLCertFile     *string             `json:"ssl_cert_file,omitempty"`
	SSLKeyFile      *string             `json:"ssl_key_file,omitempty"`
	SSLKeyPassword  *string             `json:"ssl_key_password,omitempty"`
	Description     *string             `json:"description,omitempty"`
	Status          *int                `json:"status,omitempty"`
	TagsEvalType    *int                `json:"tags_evaltype,omitempty"`
	ItemValueType   *int                `json:"item_value_type,omitempty"`
	Tags            []ConnectorTagInput `json:"tags,omitzero"`
}

func (p ConnectorCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: ConnectorCreateParams.Name is required")
	}
	if p.URL == "" {
		return errors.New("zabbix: ConnectorCreateParams.URL is required")
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ConnectorUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/connector/object
type ConnectorUpdateParams struct {
	ConnectorID     string              `json:"connectorid"`
	Name            *string             `json:"name,omitempty"`
	Protocol        *int                `json:"protocol,omitempty"`
	DataType        *int                `json:"data_type,omitempty"`
	URL             *string             `json:"url,omitempty"`
	MaxRecords      *int                `json:"max_records,omitempty"`
	MaxSenders      *int                `json:"max_senders,omitempty"`
	MaxAttempts     *int                `json:"max_attempts,omitempty"`
	AttemptInterval *string             `json:"attempt_interval,omitempty"`
	Timeout         *string             `json:"timeout,omitempty"`
	HTTPProxy       *string             `json:"http_proxy,omitempty"`
	AuthType        *int                `json:"authtype,omitempty"`
	Username        *string             `json:"username,omitempty"`
	Password        *string             `json:"password,omitempty"`
	Token           *string             `json:"token,omitempty"`
	VerifyPeer      *int                `json:"verify_peer,omitempty"`
	VerifyHost      *int                `json:"verify_host,omitempty"`
	SSLCertFile     *string             `json:"ssl_cert_file,omitempty"`
	SSLKeyFile      *string             `json:"ssl_key_file,omitempty"`
	SSLKeyPassword  *string             `json:"ssl_key_password,omitempty"`
	Description     *string             `json:"description,omitempty"`
	Status          *int                `json:"status,omitempty"`
	TagsEvalType    *int                `json:"tags_evaltype,omitempty"`
	ItemValueType   *int                `json:"item_value_type,omitempty"`
	Tags            []ConnectorTagInput `json:"tags,omitzero"`
}

func (p ConnectorUpdateParams) validate() error {
	if err := requireWriteID(p.ConnectorID, "ConnectorUpdateParams.ConnectorID"); err != nil {
		return err
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// SettingsUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/settings/object
type SettingsUpdateParams struct {
	DefaultLang          *string `json:"default_lang,omitempty"`
	DefaultTimezone      *string `json:"default_timezone,omitempty"`
	DefaultTheme         *string `json:"default_theme,omitempty"`
	SearchLimit          *int    `json:"search_limit,omitempty"`
	SeverityName0        *string `json:"severity_name_0,omitempty"`
	SeverityName1        *string `json:"severity_name_1,omitempty"`
	SeverityName2        *string `json:"severity_name_2,omitempty"`
	SeverityName3        *string `json:"severity_name_3,omitempty"`
	SeverityName4        *string `json:"severity_name_4,omitempty"`
	SeverityName5        *string `json:"severity_name_5,omitempty"`
	BlinkPeriod          *string `json:"blink_period,omitempty"`
	WorkPeriod           *string `json:"work_period,omitempty"`
	ShowTechnicalErrors  *int    `json:"show_technical_errors,omitempty"`
	HistoryPeriod        *string `json:"history_period,omitempty"`
	PeriodDefault        *string `json:"period_default,omitempty"`
	MaxPeriod            *string `json:"max_period,omitempty"`
	SNMPTrapLogging      *int    `json:"snmptrap_logging,omitempty"`
	DefaultInventoryMode *int    `json:"default_inventory_mode,omitempty"`
	GeomapsTileProvider  *string `json:"geomaps_tile_provider,omitempty"`
}

func (p SettingsUpdateParams) validate() error {
	return nil
}

// HostInterfaceCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/object
type HostInterfaceCreateParams struct {
	HostID  string              `json:"hostid"`
	Main    int                 `json:"main"`
	Type    int                 `json:"type"`
	UseIP   int                 `json:"useip"`
	IP      string              `json:"ip"`
	DNS     string              `json:"dns"`
	Port    string              `json:"port"`
	Details *SNMPInterfaceInput `json:"details,omitempty"`
}

func (p HostInterfaceCreateParams) validate() error {
	if err := requireWriteID(p.HostID, "HostInterfaceCreateParams.HostID"); err != nil {
		return err
	}
	if p.Port == "" {
		return errors.New("zabbix: HostInterfaceCreateParams.Port is required")
	}
	if p.Details != nil {
		if err := p.Details.validate(); err != nil {
			return err
		}
	}
	if p.Type < 1 || p.Type > 4 {
		return errors.New("zabbix: interface type must be between 1 and 4")
	}
	if p.Type == 2 && p.Details == nil {
		return errors.New("zabbix: SNMP interfaces require details")
	}
	if p.UseIP != 0 && p.UseIP != 1 {
		return errors.New("zabbix: interface useip must be 0 or 1")
	}
	if p.UseIP == 1 && p.IP == "" || p.UseIP == 0 && p.DNS == "" {
		return errors.New("zabbix: the selected interface address must not be empty")
	}
	return nil
}

// HostInterfaceUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/object
type HostInterfaceUpdateParams struct {
	InterfaceID string                    `json:"interfaceid"`
	Main        *int                      `json:"main,omitempty"`
	Type        *int                      `json:"type,omitempty"`
	UseIP       *int                      `json:"useip,omitempty"`
	IP          *string                   `json:"ip,omitempty"`
	DNS         *string                   `json:"dns,omitempty"`
	Port        *string                   `json:"port,omitempty"`
	Details     *SNMPInterfaceUpdateInput `json:"details,omitempty"`
}

func (p HostInterfaceUpdateParams) validate() error {
	if err := requireWriteID(p.InterfaceID, "HostInterfaceUpdateParams.InterfaceID"); err != nil {
		return err
	}
	if p.Details != nil {
		if err := p.Details.validate(); err != nil {
			return err
		}
	}
	if p.Type != nil && (*p.Type < 1 || *p.Type > 4) {
		return errors.New("zabbix: interface type must be between 1 and 4")
	}
	return nil
}

// AuthenticationUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/authentication/object
type AuthenticationUpdateParams struct {
	AuthenticationType   *int    `json:"authentication_type,omitempty"`
	HTTPAuthEnabled      *int    `json:"http_auth_enabled,omitempty"`
	HTTPLoginForm        *int    `json:"http_login_form,omitempty"`
	HTTPStripDomains     *string `json:"http_strip_domains,omitempty"`
	HTTPCaseSensitive    *int    `json:"http_case_sensitive,omitempty"`
	LDAPAuthEnabled      *int    `json:"ldap_auth_enabled,omitempty"`
	LDAPUserDirectoryID  *string `json:"ldap_userdirectoryid,omitempty"`
	LDAPCaseSensitive    *int    `json:"ldap_case_sensitive,omitempty"`
	SAMLAuthEnabled      *int    `json:"saml_auth_enabled,omitempty"`
	PasswdMinLength      *int    `json:"passwd_min_length,omitempty"`
	PasswdCheckRules     *int    `json:"passwd_check_rules,omitempty"`
	JITProvisionInterval *string `json:"jit_provision_interval,omitempty"`
	DisabledUsrGrpID     *string `json:"disabled_usrgrpid,omitempty"`
	MFAStatus            *int    `json:"mfa_status,omitempty"`
	MFAID                *string `json:"mfaid,omitempty"`
}

func (p AuthenticationUpdateParams) validate() error {
	return nil
}

// MediaTypeCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mediatype/object
type MediaTypeCreateParams struct {
	Name               string                      `json:"name"`
	Type               int                         `json:"type"`
	Status             *int                        `json:"status,omitempty"`
	SMTPServer         *string                     `json:"smtp_server,omitempty"`
	SMTPPort           *int                        `json:"smtp_port,omitempty"`
	SMTPHelo           *string                     `json:"smtp_helo,omitempty"`
	SMTPEmail          *string                     `json:"smtp_email,omitempty"`
	SMTPSecurity       *int                        `json:"smtp_security,omitempty"`
	SMTPVerifyPeer     *int                        `json:"smtp_verify_peer,omitempty"`
	SMTPVerifyHost     *int                        `json:"smtp_verify_host,omitempty"`
	SMTPAuthentication *int                        `json:"smtp_authentication,omitempty"`
	Username           *string                     `json:"username,omitempty"`
	Passwd             *string                     `json:"passwd,omitempty"`
	ExecPath           *string                     `json:"exec_path,omitempty"`
	GSMModem           *string                     `json:"gsm_modem,omitempty"`
	Script             *string                     `json:"script,omitempty"`
	Timeout            *string                     `json:"timeout,omitempty"`
	ProcessTags        *int                        `json:"process_tags,omitempty"`
	ShowEventMenu      *int                        `json:"show_event_menu,omitempty"`
	EventMenuURL       *string                     `json:"event_menu_url,omitempty"`
	EventMenuName      *string                     `json:"event_menu_name,omitempty"`
	Description        *string                     `json:"description,omitempty"`
	MaxSessions        *int                        `json:"maxsessions,omitempty"`
	MaxAttempts        *int                        `json:"maxattempts,omitempty"`
	AttemptInterval    *string                     `json:"attempt_interval,omitempty"`
	MessageFormat      *int                        `json:"message_format,omitempty"`
	Parameters         []MediaParameterInput       `json:"parameters,omitzero"`
	MessageTemplates   []MediaMessageTemplateInput `json:"message_templates,omitzero"`
}

func (p MediaTypeCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: MediaTypeCreateParams.Name is required")
	}
	for _, value := range p.Parameters {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.MessageTemplates {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// MediaTypeUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mediatype/object
type MediaTypeUpdateParams struct {
	MediaTypeID        string                      `json:"mediatypeid"`
	Name               *string                     `json:"name,omitempty"`
	Type               *int                        `json:"type,omitempty"`
	Status             *int                        `json:"status,omitempty"`
	SMTPServer         *string                     `json:"smtp_server,omitempty"`
	SMTPPort           *int                        `json:"smtp_port,omitempty"`
	SMTPHelo           *string                     `json:"smtp_helo,omitempty"`
	SMTPEmail          *string                     `json:"smtp_email,omitempty"`
	SMTPSecurity       *int                        `json:"smtp_security,omitempty"`
	SMTPVerifyPeer     *int                        `json:"smtp_verify_peer,omitempty"`
	SMTPVerifyHost     *int                        `json:"smtp_verify_host,omitempty"`
	SMTPAuthentication *int                        `json:"smtp_authentication,omitempty"`
	Username           *string                     `json:"username,omitempty"`
	Passwd             *string                     `json:"passwd,omitempty"`
	ExecPath           *string                     `json:"exec_path,omitempty"`
	GSMModem           *string                     `json:"gsm_modem,omitempty"`
	Script             *string                     `json:"script,omitempty"`
	Timeout            *string                     `json:"timeout,omitempty"`
	ProcessTags        *int                        `json:"process_tags,omitempty"`
	ShowEventMenu      *int                        `json:"show_event_menu,omitempty"`
	EventMenuURL       *string                     `json:"event_menu_url,omitempty"`
	EventMenuName      *string                     `json:"event_menu_name,omitempty"`
	Description        *string                     `json:"description,omitempty"`
	MaxSessions        *int                        `json:"maxsessions,omitempty"`
	MaxAttempts        *int                        `json:"maxattempts,omitempty"`
	AttemptInterval    *string                     `json:"attempt_interval,omitempty"`
	MessageFormat      *int                        `json:"message_format,omitempty"`
	Parameters         []MediaParameterInput       `json:"parameters,omitzero"`
	MessageTemplates   []MediaMessageTemplateInput `json:"message_templates,omitzero"`
}

func (p MediaTypeUpdateParams) validate() error {
	if err := requireWriteID(p.MediaTypeID, "MediaTypeUpdateParams.MediaTypeID"); err != nil {
		return err
	}
	for _, value := range p.Parameters {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.MessageTemplates {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ProxyGroupCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxygroup/object
type ProxyGroupCreateParams struct {
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	FailoverDelay *string `json:"failover_delay,omitempty"`
	MinOnline     *string `json:"min_online,omitempty"`
}

func (p ProxyGroupCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: ProxyGroupCreateParams.Name is required")
	}
	return nil
}

// ProxyGroupUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/proxygroup/object
type ProxyGroupUpdateParams struct {
	ProxyGroupID  string  `json:"proxy_groupid"`
	Name          *string `json:"name,omitempty"`
	Description   *string `json:"description,omitempty"`
	FailoverDelay *string `json:"failover_delay,omitempty"`
	MinOnline     *string `json:"min_online,omitempty"`
}

func (p ProxyGroupUpdateParams) validate() error {
	if err := requireWriteID(p.ProxyGroupID, "ProxyGroupUpdateParams.ProxyGroupID"); err != nil {
		return err
	}
	return nil
}

// DRuleCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/drule/object
type DRuleCreateParams struct {
	Name           string                `json:"name"`
	IPRange        string                `json:"iprange"`
	Delay          *string               `json:"delay,omitempty"`
	Status         *int                  `json:"status,omitempty"`
	ProxyID        *string               `json:"proxyid,omitempty"`
	ConcurrencyMax *int                  `json:"concurrency_max,omitempty"`
	DChecks        []DiscoveryCheckInput `json:"dchecks"`
}

func (p DRuleCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: DRuleCreateParams.Name is required")
	}
	if p.IPRange == "" {
		return errors.New("zabbix: DRuleCreateParams.IPRange is required")
	}
	if len(p.DChecks) == 0 {
		return errors.New("zabbix: DRuleCreateParams.DChecks must not be empty")
	}
	for _, value := range p.DChecks {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// DRuleUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/drule/object
type DRuleUpdateParams struct {
	DRuleID        string                      `json:"druleid"`
	Name           *string                     `json:"name,omitempty"`
	IPRange        *string                     `json:"iprange,omitempty"`
	Delay          *string                     `json:"delay,omitempty"`
	Status         *int                        `json:"status,omitempty"`
	ProxyID        *string                     `json:"proxyid,omitempty"`
	ConcurrencyMax *int                        `json:"concurrency_max,omitempty"`
	DChecks        []DiscoveryCheckUpdateInput `json:"dchecks,omitzero"`
}

func (p DRuleUpdateParams) validate() error {
	if err := requireWriteID(p.DRuleID, "DRuleUpdateParams.DRuleID"); err != nil {
		return err
	}
	for _, value := range p.DChecks {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// SLACreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/sla/object
type SLACreateParams struct {
	Name              string                     `json:"name"`
	Period            int                        `json:"period"`
	SLO               float64                    `json:"slo"`
	EffectiveDate     *int                       `json:"effective_date,omitempty"`
	Timezone          string                     `json:"timezone"`
	Status            *int                       `json:"status,omitempty"`
	Description       *string                    `json:"description,omitempty"`
	ServiceTags       []SLATagInput              `json:"service_tags"`
	Schedule          []SLAScheduleInput         `json:"schedule,omitzero"`
	ExcludedDowntimes []SLAExcludedDowntimeInput `json:"excluded_downtimes,omitzero"`
}

func (p SLACreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: SLACreateParams.Name is required")
	}
	if p.Timezone == "" {
		return errors.New("zabbix: SLACreateParams.Timezone is required")
	}
	if len(p.ServiceTags) == 0 {
		return errors.New("zabbix: SLACreateParams.ServiceTags must not be empty")
	}
	for _, value := range p.ServiceTags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Schedule {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.ExcludedDowntimes {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// SLAUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/sla/object
type SLAUpdateParams struct {
	SLAID             string                     `json:"slaid"`
	Name              *string                    `json:"name,omitempty"`
	Period            *int                       `json:"period,omitempty"`
	SLO               *float64                   `json:"slo,omitempty"`
	EffectiveDate     *int                       `json:"effective_date,omitempty"`
	Timezone          *string                    `json:"timezone,omitempty"`
	Status            *int                       `json:"status,omitempty"`
	Description       *string                    `json:"description,omitempty"`
	ServiceTags       []SLATagInput              `json:"service_tags,omitzero"`
	Schedule          []SLAScheduleInput         `json:"schedule,omitzero"`
	ExcludedDowntimes []SLAExcludedDowntimeInput `json:"excluded_downtimes,omitzero"`
}

func (p SLAUpdateParams) validate() error {
	if err := requireWriteID(p.SLAID, "SLAUpdateParams.SLAID"); err != nil {
		return err
	}
	for _, value := range p.ServiceTags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Schedule {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.ExcludedDowntimes {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// RegexpCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/regexp/object
type RegexpCreateParams struct {
	Name        string                  `json:"name"`
	TestString  *string                 `json:"test_string,omitempty"`
	Expressions []RegexpExpressionInput `json:"expressions"`
}

func (p RegexpCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: RegexpCreateParams.Name is required")
	}
	if len(p.Expressions) == 0 {
		return errors.New("zabbix: RegexpCreateParams.Expressions must not be empty")
	}
	for _, value := range p.Expressions {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// RegexpUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/regexp/object
type RegexpUpdateParams struct {
	RegexpID    string                  `json:"regexpid"`
	Name        *string                 `json:"name,omitempty"`
	TestString  *string                 `json:"test_string,omitempty"`
	Expressions []RegexpExpressionInput `json:"expressions,omitzero"`
}

func (p RegexpUpdateParams) validate() error {
	if err := requireWriteID(p.RegexpID, "RegexpUpdateParams.RegexpID"); err != nil {
		return err
	}
	for _, value := range p.Expressions {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// HousekeepingUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/housekeeping/object
type HousekeepingUpdateParams struct {
	HKEventsMode      *int    `json:"hk_events_mode,omitempty"`
	HKEventsTrigger   *string `json:"hk_events_trigger,omitempty"`
	HKEventsInternal  *string `json:"hk_events_internal,omitempty"`
	HKEventsDiscovery *string `json:"hk_events_discovery,omitempty"`
	HKEventsAutoreg   *string `json:"hk_events_autoreg,omitempty"`
	HKServicesMode    *int    `json:"hk_services_mode,omitempty"`
	HKServices        *string `json:"hk_services,omitempty"`
	HKAuditMode       *int    `json:"hk_audit_mode,omitempty"`
	HKAudit           *string `json:"hk_audit,omitempty"`
	HKSessionsMode    *int    `json:"hk_sessions_mode,omitempty"`
	HKSessions        *string `json:"hk_sessions,omitempty"`
	HKHistoryMode     *int    `json:"hk_history_mode,omitempty"`
	HKHistoryGlobal   *int    `json:"hk_history_global,omitempty"`
	HKHistory         *string `json:"hk_history,omitempty"`
	HKTrendsMode      *int    `json:"hk_trends_mode,omitempty"`
	HKTrendsGlobal    *int    `json:"hk_trends_global,omitempty"`
	HKTrends          *string `json:"hk_trends,omitempty"`
	CompressionStatus *int    `json:"compression_status,omitempty"`
	CompressOlder     *string `json:"compress_older,omitempty"`
}

func (p HousekeepingUpdateParams) validate() error {
	return nil
}

// DashboardCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dashboard/object
type DashboardCreateParams struct {
	Name          string                    `json:"name"`
	UserID        *string                   `json:"userid,omitempty"`
	Private       *int                      `json:"private,omitempty"`
	DisplayPeriod *int                      `json:"display_period,omitempty"`
	AutoStart     *int                      `json:"auto_start,omitempty"`
	Pages         []DashboardPageInput      `json:"pages"`
	Users         []DashboardUserInput      `json:"users,omitzero"`
	UserGroups    []DashboardUserGroupInput `json:"userGroups,omitzero"`
}

func (p DashboardCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: DashboardCreateParams.Name is required")
	}
	if len(p.Pages) == 0 {
		return errors.New("zabbix: DashboardCreateParams.Pages must not be empty")
	}
	for _, value := range p.Pages {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Users {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.UserGroups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// DashboardUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dashboard/object
type DashboardUpdateParams struct {
	DashboardID   string                     `json:"dashboardid"`
	Name          *string                    `json:"name,omitempty"`
	UserID        *string                    `json:"userid,omitempty"`
	Private       *int                       `json:"private,omitempty"`
	DisplayPeriod *int                       `json:"display_period,omitempty"`
	AutoStart     *int                       `json:"auto_start,omitempty"`
	Pages         []DashboardPageUpdateInput `json:"pages,omitzero"`
	Users         []DashboardUserInput       `json:"users,omitzero"`
	UserGroups    []DashboardUserGroupInput  `json:"userGroups,omitzero"`
}

func (p DashboardUpdateParams) validate() error {
	if err := requireWriteID(p.DashboardID, "DashboardUpdateParams.DashboardID"); err != nil {
		return err
	}
	for _, value := range p.Pages {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Users {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.UserGroups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ItemPrototypeCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/itemprototype/object
type ItemPrototypeCreateParams struct {
	HostID        string                   `json:"hostid"`
	Name          string                   `json:"name"`
	Key           string                   `json:"key_"`
	Type          int                      `json:"type"`
	ValueType     int                      `json:"value_type"`
	Delay         *string                  `json:"delay,omitempty"`
	History       *string                  `json:"history,omitempty"`
	Trends        *string                  `json:"trends,omitempty"`
	Status        *int                     `json:"status,omitempty"`
	Units         *string                  `json:"units,omitempty"`
	ValueMapID    *string                  `json:"valuemapid,omitempty"`
	Description   *string                  `json:"description,omitempty"`
	InterfaceID   *string                  `json:"interfaceid,omitempty"`
	SNMPOID       *string                  `json:"snmp_oid,omitempty"`
	Params        *string                  `json:"params,omitempty"`
	MasterItemID  *string                  `json:"master_itemid,omitempty"`
	Discover      *int                     `json:"discover,omitempty"`
	RuleID        string                   `json:"ruleid"`
	Tags          []TagInput               `json:"tags,omitzero"`
	Preprocessing []ItemPreprocessingInput `json:"preprocessing,omitzero"`
}

func (p ItemPrototypeCreateParams) validate() error {
	if err := requireWriteID(p.HostID, "ItemPrototypeCreateParams.HostID"); err != nil {
		return err
	}
	if p.Name == "" {
		return errors.New("zabbix: ItemPrototypeCreateParams.Name is required")
	}
	if p.Key == "" {
		return errors.New("zabbix: ItemPrototypeCreateParams.Key is required")
	}
	if err := requireWriteID(p.RuleID, "ItemPrototypeCreateParams.RuleID"); err != nil {
		return err
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Preprocessing {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ItemPrototypeUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/itemprototype/object
type ItemPrototypeUpdateParams struct {
	ItemID        string                   `json:"itemid"`
	Name          *string                  `json:"name,omitempty"`
	Key           *string                  `json:"key_,omitempty"`
	Type          *int                     `json:"type,omitempty"`
	ValueType     *int                     `json:"value_type,omitempty"`
	Delay         *string                  `json:"delay,omitempty"`
	History       *string                  `json:"history,omitempty"`
	Trends        *string                  `json:"trends,omitempty"`
	Status        *int                     `json:"status,omitempty"`
	Units         *string                  `json:"units,omitempty"`
	ValueMapID    *string                  `json:"valuemapid,omitempty"`
	Description   *string                  `json:"description,omitempty"`
	InterfaceID   *string                  `json:"interfaceid,omitempty"`
	SNMPOID       *string                  `json:"snmp_oid,omitempty"`
	Params        *string                  `json:"params,omitempty"`
	MasterItemID  *string                  `json:"master_itemid,omitempty"`
	Discover      *int                     `json:"discover,omitempty"`
	Tags          []TagInput               `json:"tags,omitzero"`
	Preprocessing []ItemPreprocessingInput `json:"preprocessing,omitzero"`
}

func (p ItemPrototypeUpdateParams) validate() error {
	if err := requireWriteID(p.ItemID, "ItemPrototypeUpdateParams.ItemID"); err != nil {
		return err
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Preprocessing {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// MapCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/map/object
type MapCreateParams struct {
	Name               string              `json:"name"`
	Width              int                 `json:"width"`
	Height             int                 `json:"height"`
	BackgroundID       *string             `json:"backgroundid,omitempty"`
	LabelType          *int                `json:"label_type,omitempty"`
	LabelLocation      *int                `json:"label_location,omitempty"`
	Highlight          *int                `json:"highlight,omitempty"`
	ExpandProblem      *int                `json:"expandproblem,omitempty"`
	MarkElements       *int                `json:"markelements,omitempty"`
	ShowUnack          *int                `json:"show_unack,omitempty"`
	GridSize           *int                `json:"grid_size,omitempty"`
	GridShow           *int                `json:"grid_show,omitempty"`
	GridAlign          *int                `json:"grid_align,omitempty"`
	LabelFormat        *int                `json:"label_format,omitempty"`
	LabelTypeHost      *int                `json:"label_type_host,omitempty"`
	LabelTypeHostGroup *int                `json:"label_type_hostgroup,omitempty"`
	LabelTypeTrigger   *int                `json:"label_type_trigger,omitempty"`
	LabelTypeMap       *int                `json:"label_type_map,omitempty"`
	LabelTypeImage     *int                `json:"label_type_image,omitempty"`
	LabelStringHost    *string             `json:"label_string_host,omitempty"`
	IconMapID          *string             `json:"iconmapid,omitempty"`
	ExpandMacros       *int                `json:"expand_macros,omitempty"`
	SeverityMin        *int                `json:"severity_min,omitempty"`
	UserID             *string             `json:"userid,omitempty"`
	Private            *int                `json:"private,omitempty"`
	ShowSuppressed     *int                `json:"show_suppressed,omitempty"`
	Selements          []MapElementInput   `json:"selements,omitzero"`
	Links              []MapLinkInput      `json:"links,omitzero"`
	URLs               []MapURLInput       `json:"urls,omitzero"`
	Users              []MapUserInput      `json:"users,omitzero"`
	UserGroups         []MapUserGroupInput `json:"userGroups,omitzero"`
	Shapes             []MapShapeInput     `json:"shapes,omitzero"`
	Lines              []MapLineInput      `json:"lines,omitzero"`
}

func (p MapCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: MapCreateParams.Name is required")
	}
	for _, value := range p.Selements {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Links {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.URLs {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Users {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.UserGroups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Shapes {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Lines {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// MapUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/map/object
type MapUpdateParams struct {
	SysmapID           string                `json:"sysmapid"`
	Name               *string               `json:"name,omitempty"`
	Width              *int                  `json:"width,omitempty"`
	Height             *int                  `json:"height,omitempty"`
	BackgroundID       *string               `json:"backgroundid,omitempty"`
	LabelType          *int                  `json:"label_type,omitempty"`
	LabelLocation      *int                  `json:"label_location,omitempty"`
	Highlight          *int                  `json:"highlight,omitempty"`
	ExpandProblem      *int                  `json:"expandproblem,omitempty"`
	MarkElements       *int                  `json:"markelements,omitempty"`
	ShowUnack          *int                  `json:"show_unack,omitempty"`
	GridSize           *int                  `json:"grid_size,omitempty"`
	GridShow           *int                  `json:"grid_show,omitempty"`
	GridAlign          *int                  `json:"grid_align,omitempty"`
	LabelFormat        *int                  `json:"label_format,omitempty"`
	LabelTypeHost      *int                  `json:"label_type_host,omitempty"`
	LabelTypeHostGroup *int                  `json:"label_type_hostgroup,omitempty"`
	LabelTypeTrigger   *int                  `json:"label_type_trigger,omitempty"`
	LabelTypeMap       *int                  `json:"label_type_map,omitempty"`
	LabelTypeImage     *int                  `json:"label_type_image,omitempty"`
	LabelStringHost    *string               `json:"label_string_host,omitempty"`
	IconMapID          *string               `json:"iconmapid,omitempty"`
	ExpandMacros       *int                  `json:"expand_macros,omitempty"`
	SeverityMin        *int                  `json:"severity_min,omitempty"`
	UserID             *string               `json:"userid,omitempty"`
	Private            *int                  `json:"private,omitempty"`
	ShowSuppressed     *int                  `json:"show_suppressed,omitempty"`
	Selements          []MapElementInput     `json:"selements,omitzero"`
	Links              []MapLinkInput        `json:"links,omitzero"`
	URLs               []MapURLInput         `json:"urls,omitzero"`
	Users              []MapUserInput        `json:"users,omitzero"`
	UserGroups         []MapUserGroupInput   `json:"userGroups,omitzero"`
	Shapes             []MapShapeUpdateInput `json:"shapes,omitzero"`
	Lines              []MapLineUpdateInput  `json:"lines,omitzero"`
}

func (p MapUpdateParams) validate() error {
	if err := requireWriteID(p.SysmapID, "MapUpdateParams.SysmapID"); err != nil {
		return err
	}
	for _, value := range p.Selements {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Links {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.URLs {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Users {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.UserGroups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Shapes {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Lines {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// TriggerCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/trigger/object
type TriggerCreateParams struct {
	Description        string     `json:"description"`
	Expression         string     `json:"expression"`
	EventName          *string    `json:"event_name,omitempty"`
	OpData             *string    `json:"opdata,omitempty"`
	Comments           *string    `json:"comments,omitempty"`
	Priority           *int       `json:"priority,omitempty"`
	Status             *int       `json:"status,omitempty"`
	Type               *int       `json:"type,omitempty"`
	RecoveryMode       *int       `json:"recovery_mode,omitempty"`
	RecoveryExpression *string    `json:"recovery_expression,omitempty"`
	CorrelationMode    *int       `json:"correlation_mode,omitempty"`
	CorrelationTag     *string    `json:"correlation_tag,omitempty"`
	ManualClose        *int       `json:"manual_close,omitempty"`
	URL                *string    `json:"url,omitempty"`
	URLName            *string    `json:"url_name,omitempty"`
	UUID               *string    `json:"uuid,omitempty"`
	Tags               []TagInput `json:"tags,omitzero"`
}

func (p TriggerCreateParams) validate() error {
	if p.Description == "" {
		return errors.New("zabbix: TriggerCreateParams.Description is required")
	}
	if p.Expression == "" {
		return errors.New("zabbix: TriggerCreateParams.Expression is required")
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// TriggerUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/trigger/object
type TriggerUpdateParams struct {
	TriggerID          string     `json:"triggerid"`
	Description        *string    `json:"description,omitempty"`
	Expression         *string    `json:"expression,omitempty"`
	EventName          *string    `json:"event_name,omitempty"`
	OpData             *string    `json:"opdata,omitempty"`
	Comments           *string    `json:"comments,omitempty"`
	Priority           *int       `json:"priority,omitempty"`
	Status             *int       `json:"status,omitempty"`
	Type               *int       `json:"type,omitempty"`
	RecoveryMode       *int       `json:"recovery_mode,omitempty"`
	RecoveryExpression *string    `json:"recovery_expression,omitempty"`
	CorrelationMode    *int       `json:"correlation_mode,omitempty"`
	CorrelationTag     *string    `json:"correlation_tag,omitempty"`
	ManualClose        *int       `json:"manual_close,omitempty"`
	URL                *string    `json:"url,omitempty"`
	URLName            *string    `json:"url_name,omitempty"`
	UUID               *string    `json:"uuid,omitempty"`
	Tags               []TagInput `json:"tags,omitzero"`
}

func (p TriggerUpdateParams) validate() error {
	if err := requireWriteID(p.TriggerID, "TriggerUpdateParams.TriggerID"); err != nil {
		return err
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ItemCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/item/object
type ItemCreateParams struct {
	HostID       string     `json:"hostid"`
	Name         string     `json:"name"`
	Key          string     `json:"key_"`
	Type         int        `json:"type"`
	ValueType    int        `json:"value_type"`
	Delay        *string    `json:"delay,omitempty"`
	History      *string    `json:"history,omitempty"`
	Trends       *string    `json:"trends,omitempty"`
	Status       *int       `json:"status,omitempty"`
	Units        *string    `json:"units,omitempty"`
	ValueMapID   *string    `json:"valuemapid,omitempty"`
	Description  *string    `json:"description,omitempty"`
	InterfaceID  *string    `json:"interfaceid,omitempty"`
	SNMPOID      *string    `json:"snmp_oid,omitempty"`
	Params       *string    `json:"params,omitempty"`
	MasterItemID *string    `json:"master_itemid,omitempty"`
	Tags         []TagInput `json:"tags,omitzero"`
}

func (p ItemCreateParams) validate() error {
	if err := requireWriteID(p.HostID, "ItemCreateParams.HostID"); err != nil {
		return err
	}
	if p.Name == "" {
		return errors.New("zabbix: ItemCreateParams.Name is required")
	}
	if p.Key == "" {
		return errors.New("zabbix: ItemCreateParams.Key is required")
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ItemUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/item/object
type ItemUpdateParams struct {
	ItemID       string     `json:"itemid"`
	Name         *string    `json:"name,omitempty"`
	Key          *string    `json:"key_,omitempty"`
	Type         *int       `json:"type,omitempty"`
	ValueType    *int       `json:"value_type,omitempty"`
	Delay        *string    `json:"delay,omitempty"`
	History      *string    `json:"history,omitempty"`
	Trends       *string    `json:"trends,omitempty"`
	Status       *int       `json:"status,omitempty"`
	Units        *string    `json:"units,omitempty"`
	ValueMapID   *string    `json:"valuemapid,omitempty"`
	Description  *string    `json:"description,omitempty"`
	InterfaceID  *string    `json:"interfaceid,omitempty"`
	SNMPOID      *string    `json:"snmp_oid,omitempty"`
	Params       *string    `json:"params,omitempty"`
	MasterItemID *string    `json:"master_itemid,omitempty"`
	Tags         []TagInput `json:"tags,omitzero"`
}

func (p ItemUpdateParams) validate() error {
	if err := requireWriteID(p.ItemID, "ItemUpdateParams.ItemID"); err != nil {
		return err
	}
	for _, value := range p.Tags {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ImageCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/image/object
type ImageCreateParams struct {
	ImageType int    `json:"imagetype"`
	Name      string `json:"name"`
	Image     string `json:"image"`
}

func (p ImageCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: ImageCreateParams.Name is required")
	}
	if p.Image == "" {
		return errors.New("zabbix: ImageCreateParams.Image is required")
	}
	return nil
}

// ImageUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/image/object
type ImageUpdateParams struct {
	ImageID string  `json:"imageid"`
	Name    *string `json:"name,omitempty"`
	Image   *string `json:"image,omitempty"`
}

func (p ImageUpdateParams) validate() error {
	if err := requireWriteID(p.ImageID, "ImageUpdateParams.ImageID"); err != nil {
		return err
	}
	return nil
}

// ReportCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/report/object
type ReportCreateParams struct {
	UserID      string                 `json:"userid"`
	Name        string                 `json:"name"`
	DashboardID string                 `json:"dashboardid"`
	Period      *int                   `json:"period,omitempty"`
	Cycle       *int                   `json:"cycle,omitempty"`
	StartTime   *int                   `json:"start_time,omitempty"`
	ActiveSince *string                `json:"active_since,omitempty"`
	ActiveTill  *string                `json:"active_till,omitempty"`
	Status      *int                   `json:"status,omitempty"`
	Description *string                `json:"description,omitempty"`
	Subject     *string                `json:"subject,omitempty"`
	Message     *string                `json:"message,omitempty"`
	Users       []ReportUserInput      `json:"users,omitzero"`
	UserGroups  []ReportUserGroupInput `json:"user_groups,omitzero"`
}

func (p ReportCreateParams) validate() error {
	if err := requireWriteID(p.UserID, "ReportCreateParams.UserID"); err != nil {
		return err
	}
	if p.Name == "" {
		return errors.New("zabbix: ReportCreateParams.Name is required")
	}
	if err := requireWriteID(p.DashboardID, "ReportCreateParams.DashboardID"); err != nil {
		return err
	}
	for _, value := range p.Users {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.UserGroups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if len(p.Users) == 0 && len(p.UserGroups) == 0 {
		return errors.New("zabbix: report requires a user or user group")
	}
	return nil
}

// ReportUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/report/object
type ReportUpdateParams struct {
	ReportID    string                 `json:"reportid"`
	UserID      *string                `json:"userid,omitempty"`
	Name        *string                `json:"name,omitempty"`
	DashboardID *string                `json:"dashboardid,omitempty"`
	Period      *int                   `json:"period,omitempty"`
	Cycle       *int                   `json:"cycle,omitempty"`
	StartTime   *int                   `json:"start_time,omitempty"`
	ActiveSince *string                `json:"active_since,omitempty"`
	ActiveTill  *string                `json:"active_till,omitempty"`
	Status      *int                   `json:"status,omitempty"`
	Description *string                `json:"description,omitempty"`
	Subject     *string                `json:"subject,omitempty"`
	Message     *string                `json:"message,omitempty"`
	Users       []ReportUserInput      `json:"users,omitzero"`
	UserGroups  []ReportUserGroupInput `json:"user_groups,omitzero"`
}

func (p ReportUpdateParams) validate() error {
	if err := requireWriteID(p.ReportID, "ReportUpdateParams.ReportID"); err != nil {
		return err
	}
	for _, value := range p.Users {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.UserGroups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// CorrelationCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/correlation/object
type CorrelationCreateParams struct {
	Name        string                      `json:"name"`
	Description *string                     `json:"description,omitempty"`
	Status      *int                        `json:"status,omitempty"`
	Filter      *CorrelationFilterInput     `json:"filter"`
	Operations  []CorrelationOperationInput `json:"operations"`
}

func (p CorrelationCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: CorrelationCreateParams.Name is required")
	}
	if p.Filter == nil {
		return errors.New("zabbix: CorrelationCreateParams.Filter is required")
	}
	if p.Filter != nil {
		if err := p.Filter.validate(); err != nil {
			return err
		}
	}
	if len(p.Operations) == 0 {
		return errors.New("zabbix: CorrelationCreateParams.Operations must not be empty")
	}
	for _, value := range p.Operations {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// CorrelationUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/correlation/object
type CorrelationUpdateParams struct {
	CorrelationID string                      `json:"correlationid"`
	Name          *string                     `json:"name,omitempty"`
	Description   *string                     `json:"description,omitempty"`
	Status        *int                        `json:"status,omitempty"`
	Filter        *CorrelationFilterInput     `json:"filter,omitempty"`
	Operations    []CorrelationOperationInput `json:"operations,omitzero"`
}

func (p CorrelationUpdateParams) validate() error {
	if err := requireWriteID(p.CorrelationID, "CorrelationUpdateParams.CorrelationID"); err != nil {
		return err
	}
	if p.Filter != nil {
		if err := p.Filter.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Operations {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// IconMapCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/iconmap/object
type IconMapCreateParams struct {
	Name          string             `json:"name"`
	DefaultIconID string             `json:"default_iconid"`
	Mappings      []IconMappingInput `json:"mappings"`
}

func (p IconMapCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: IconMapCreateParams.Name is required")
	}
	if err := requireWriteID(p.DefaultIconID, "IconMapCreateParams.DefaultIconID"); err != nil {
		return err
	}
	if len(p.Mappings) == 0 {
		return errors.New("zabbix: IconMapCreateParams.Mappings must not be empty")
	}
	for _, value := range p.Mappings {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// IconMapUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/iconmap/object
type IconMapUpdateParams struct {
	IconMapID     string             `json:"iconmapid"`
	Name          *string            `json:"name,omitempty"`
	DefaultIconID *string            `json:"default_iconid,omitempty"`
	Mappings      []IconMappingInput `json:"mappings,omitzero"`
}

func (p IconMapUpdateParams) validate() error {
	if err := requireWriteID(p.IconMapID, "IconMapUpdateParams.IconMapID"); err != nil {
		return err
	}
	for _, value := range p.Mappings {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ScriptCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/object
type ScriptCreateParams struct {
	Name                     string                 `json:"name"`
	Command                  *string                `json:"command,omitempty"`
	HostAccess               *int                   `json:"host_access,omitempty"`
	UsrGrpID                 *string                `json:"usrgrpid,omitempty"`
	GroupID                  *string                `json:"groupid,omitempty"`
	Description              *string                `json:"description,omitempty"`
	Confirmation             *string                `json:"confirmation,omitempty"`
	Type                     int                    `json:"type"`
	ExecuteOn                *int                   `json:"execute_on,omitempty"`
	Timeout                  *string                `json:"timeout,omitempty"`
	Scope                    int                    `json:"scope"`
	Port                     *string                `json:"port,omitempty"`
	AuthType                 *int                   `json:"authtype,omitempty"`
	Username                 *string                `json:"username,omitempty"`
	Password                 *string                `json:"password,omitempty"`
	PublicKey                *string                `json:"publickey,omitempty"`
	PrivateKey               *string                `json:"privatekey,omitempty"`
	MenuPath                 *string                `json:"menu_path,omitempty"`
	URL                      *string                `json:"url,omitempty"`
	NewWindow                *int                   `json:"new_window,omitempty"`
	ManualInput              *int                   `json:"manualinput,omitempty"`
	ManualInputPrompt        *string                `json:"manualinput_prompt,omitempty"`
	ManualInputValidator     *string                `json:"manualinput_validator,omitempty"`
	ManualInputValidatorType *int                   `json:"manualinput_validator_type,omitempty"`
	ManualInputDefaultValue  *string                `json:"manualinput_default_value,omitempty"`
	Parameters               []ScriptParameterInput `json:"parameters,omitzero"`
}

func (p ScriptCreateParams) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: ScriptCreateParams.Name is required")
	}
	for _, value := range p.Parameters {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ScriptUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/object
type ScriptUpdateParams struct {
	ScriptID                 string                 `json:"scriptid"`
	Name                     *string                `json:"name,omitempty"`
	Command                  *string                `json:"command,omitempty"`
	HostAccess               *int                   `json:"host_access,omitempty"`
	UsrGrpID                 *string                `json:"usrgrpid,omitempty"`
	GroupID                  *string                `json:"groupid,omitempty"`
	Description              *string                `json:"description,omitempty"`
	Confirmation             *string                `json:"confirmation,omitempty"`
	Type                     *int                   `json:"type,omitempty"`
	ExecuteOn                *int                   `json:"execute_on,omitempty"`
	Timeout                  *string                `json:"timeout,omitempty"`
	Scope                    *int                   `json:"scope,omitempty"`
	Port                     *string                `json:"port,omitempty"`
	AuthType                 *int                   `json:"authtype,omitempty"`
	Username                 *string                `json:"username,omitempty"`
	Password                 *string                `json:"password,omitempty"`
	PublicKey                *string                `json:"publickey,omitempty"`
	PrivateKey               *string                `json:"privatekey,omitempty"`
	MenuPath                 *string                `json:"menu_path,omitempty"`
	URL                      *string                `json:"url,omitempty"`
	NewWindow                *int                   `json:"new_window,omitempty"`
	ManualInput              *int                   `json:"manualinput,omitempty"`
	ManualInputPrompt        *string                `json:"manualinput_prompt,omitempty"`
	ManualInputValidator     *string                `json:"manualinput_validator,omitempty"`
	ManualInputValidatorType *int                   `json:"manualinput_validator_type,omitempty"`
	ManualInputDefaultValue  *string                `json:"manualinput_default_value,omitempty"`
	Parameters               []ScriptParameterInput `json:"parameters,omitzero"`
}

func (p ScriptUpdateParams) validate() error {
	if err := requireWriteID(p.ScriptID, "ScriptUpdateParams.ScriptID"); err != nil {
		return err
	}
	for _, value := range p.Parameters {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ModuleCreateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/module/object
type ModuleCreateParams struct {
	ID           string `json:"id"`
	RelativePath string `json:"relative_path"`
	Status       *int   `json:"status,omitempty"`
	Config       any    `json:"config,omitzero"`
}

func (p ModuleCreateParams) validate() error {
	if p.ID == "" {
		return errors.New("zabbix: ModuleCreateParams.ID is required")
	}
	if p.RelativePath == "" {
		return errors.New("zabbix: ModuleCreateParams.RelativePath is required")
	}
	return nil
}

// ModuleUpdateParams contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/module/object
type ModuleUpdateParams struct {
	ModuleID     string  `json:"moduleid"`
	ID           *string `json:"id,omitempty"`
	RelativePath *string `json:"relative_path,omitempty"`
	Status       *int    `json:"status,omitempty"`
	Config       any     `json:"config,omitzero"`
}

func (p ModuleUpdateParams) validate() error {
	if err := requireWriteID(p.ModuleID, "ModuleUpdateParams.ModuleID"); err != nil {
		return err
	}
	return nil
}
