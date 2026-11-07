package zabbix

import "errors"

// HostRef identifies an existing object in a write operation.
type HostRef struct {
	HostID string `json:"hostid"`
}

func (p HostRef) validate() error { return requireWriteID(p.HostID, "HostRef.HostID") }

// HostGroupRef identifies an existing object in a write operation.
type HostGroupRef struct {
	GroupID string `json:"groupid"`
}

func (p HostGroupRef) validate() error { return requireWriteID(p.GroupID, "HostGroupRef.GroupID") }

// TemplateGroupRef identifies an existing object in a write operation.
type TemplateGroupRef struct {
	GroupID string `json:"groupid"`
}

func (p TemplateGroupRef) validate() error {
	return requireWriteID(p.GroupID, "TemplateGroupRef.GroupID")
}

// TemplateRef identifies an existing object in a write operation.
type TemplateRef struct {
	TemplateID string `json:"templateid"`
}

func (p TemplateRef) validate() error { return requireWriteID(p.TemplateID, "TemplateRef.TemplateID") }

// UserRef identifies an existing object in a write operation.
type UserRef struct {
	UserID string `json:"userid"`
}

func (p UserRef) validate() error { return requireWriteID(p.UserID, "UserRef.UserID") }

// UserGroupRef identifies an existing object in a write operation.
type UserGroupRef struct {
	UsrGrpID string `json:"usrgrpid"`
}

func (p UserGroupRef) validate() error { return requireWriteID(p.UsrGrpID, "UserGroupRef.UsrGrpID") }

// ServiceRef identifies an existing object in a write operation.
type ServiceRef struct {
	ServiceID string `json:"serviceid"`
}

func (p ServiceRef) validate() error { return requireWriteID(p.ServiceID, "ServiceRef.ServiceID") }

// TriggerRef identifies an existing object in a write operation.
type TriggerRef struct {
	TriggerID string `json:"triggerid"`
}

func (p TriggerRef) validate() error { return requireWriteID(p.TriggerID, "TriggerRef.TriggerID") }

// ScriptRef identifies an existing object in a write operation.
type ScriptRef struct {
	ScriptID string `json:"scriptid"`
}

func (p ScriptRef) validate() error { return requireWriteID(p.ScriptID, "ScriptRef.ScriptID") }

// TagInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/host/object
type TagInput struct {
	Tag   string  `json:"tag"`
	Value *string `json:"value,omitempty"`
}

func (p TagInput) validate() error {
	if p.Tag == "" {
		return errors.New("zabbix: TagInput.Tag is required")
	}
	return nil
}

// SNMPInterfaceInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/object
type SNMPInterfaceInput struct {
	Version        int     `json:"version"`
	Bulk           *int    `json:"bulk,omitempty"`
	Community      *string `json:"community,omitempty"`
	SecurityName   *string `json:"securityname,omitempty"`
	SecurityLevel  *int    `json:"securitylevel,omitempty"`
	AuthPassphrase *string `json:"authpassphrase,omitempty"`
	PrivPassphrase *string `json:"privpassphrase,omitempty"`
	AuthProtocol   *int    `json:"authprotocol,omitempty"`
	PrivProtocol   *int    `json:"privprotocol,omitempty"`
	ContextName    *string `json:"contextname,omitempty"`
	MaxRepetitions *int    `json:"max_repetitions,omitempty"`
}

func (p SNMPInterfaceInput) validate() error {
	if p.Version < 1 || p.Version > 3 {
		return errors.New("zabbix: SNMP version must be between 1 and 3")
	}
	if p.Version != 3 && (p.Community == nil || *p.Community == "") {
		return errors.New("zabbix: SNMP v1/v2 require community")
	}
	return nil
}

// SNMPInterfaceUpdateInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/object
type SNMPInterfaceUpdateInput struct {
	Version        *int    `json:"version,omitempty"`
	Bulk           *int    `json:"bulk,omitempty"`
	Community      *string `json:"community,omitempty"`
	SecurityName   *string `json:"securityname,omitempty"`
	SecurityLevel  *int    `json:"securitylevel,omitempty"`
	AuthPassphrase *string `json:"authpassphrase,omitempty"`
	PrivPassphrase *string `json:"privpassphrase,omitempty"`
	AuthProtocol   *int    `json:"authprotocol,omitempty"`
	PrivProtocol   *int    `json:"privprotocol,omitempty"`
	ContextName    *string `json:"contextname,omitempty"`
	MaxRepetitions *int    `json:"max_repetitions,omitempty"`
}

func (p SNMPInterfaceUpdateInput) validate() error {
	if p.Version != nil && (*p.Version < 1 || *p.Version > 3) {
		return errors.New("zabbix: SNMP version must be between 1 and 3")
	}
	return nil
}

// HostInterfaceInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/object
type HostInterfaceInput struct {
	Main    int                 `json:"main"`
	Type    int                 `json:"type"`
	UseIP   int                 `json:"useip"`
	IP      string              `json:"ip"`
	DNS     string              `json:"dns"`
	Port    string              `json:"port"`
	Details *SNMPInterfaceInput `json:"details,omitempty"`
}

func (p HostInterfaceInput) validate() error {
	if p.Port == "" {
		return errors.New("zabbix: HostInterfaceInput.Port is required")
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

// HostInterfaceUpdateInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostinterface/object
type HostInterfaceUpdateInput struct {
	InterfaceID *string                   `json:"interfaceid,omitempty"`
	Main        *int                      `json:"main,omitempty"`
	Type        *int                      `json:"type,omitempty"`
	UseIP       *int                      `json:"useip,omitempty"`
	IP          *string                   `json:"ip,omitempty"`
	DNS         *string                   `json:"dns,omitempty"`
	Port        *string                   `json:"port,omitempty"`
	Details     *SNMPInterfaceUpdateInput `json:"details,omitempty"`
}

func (p HostInterfaceUpdateInput) validate() error {
	if p.InterfaceID != nil {
		if err := requireWriteID(*p.InterfaceID, "HostInterfaceUpdateInput.InterfaceID"); err != nil {
			return err
		}
	}
	if p.Details != nil {
		if err := p.Details.validate(); err != nil {
			return err
		}
	}
	if p.Type != nil && (*p.Type < 1 || *p.Type > 4) {
		return errors.New("zabbix: interface type must be between 1 and 4")
	}
	if p.InterfaceID == nil && (p.Type == nil || p.Main == nil || p.UseIP == nil || p.IP == nil || p.DNS == nil || p.Port == nil || *p.Port == "") {
		return errors.New("zabbix: new host interfaces require type, main, useip, ip, dns and port")
	}
	if p.InterfaceID == nil {
		if *p.UseIP != 0 && *p.UseIP != 1 {
			return errors.New("zabbix: interface useip must be 0 or 1")
		}
		if *p.UseIP == 1 && *p.IP == "" || *p.UseIP == 0 && *p.DNS == "" {
			return errors.New("zabbix: the selected interface address must not be empty")
		}
		if *p.Type == 2 {
			if p.Details == nil || p.Details.Version == nil {
				return errors.New("zabbix: new SNMP interfaces require details and version")
			}
			if *p.Details.Version < 3 && (p.Details.Community == nil || *p.Details.Community == "") {
				return errors.New("zabbix: SNMP versions 1 and 2 require community")
			}
		}
	}
	return nil
}

// HostMacroInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/object
type HostMacroInput struct {
	Macro       string  `json:"macro"`
	Value       string  `json:"value"`
	Type        *int    `json:"type,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (p HostMacroInput) validate() error {
	if p.Macro == "" {
		return errors.New("zabbix: HostMacroInput.Macro is required")
	}
	return nil
}

// HostMacroUpdateInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usermacro/object
type HostMacroUpdateInput struct {
	HostMacroID *string `json:"hostmacroid,omitempty"`
	Macro       *string `json:"macro,omitempty"`
	Value       *string `json:"value,omitempty"`
	Type        *int    `json:"type,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (p HostMacroUpdateInput) validate() error {
	if p.HostMacroID == nil && (p.Macro == nil || *p.Macro == "" || p.Value == nil) {
		return errors.New("zabbix: new host macros require macro and value")
	}
	return nil
}

// GraphItemInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graphitem/object
type GraphItemInput struct {
	ItemID    string `json:"itemid"`
	DrawType  *int   `json:"drawtype,omitempty"`
	SortOrder *int   `json:"sortorder,omitempty"`
	Color     string `json:"color"`
	YAxisSide *int   `json:"yaxisside,omitempty"`
	CalcFnc   *int   `json:"calc_fnc,omitempty"`
	Type      *int   `json:"type,omitempty"`
}

func (p GraphItemInput) validate() error {
	if err := requireWriteID(p.ItemID, "GraphItemInput.ItemID"); err != nil {
		return err
	}
	if p.Color == "" {
		return errors.New("zabbix: GraphItemInput.Color is required")
	}
	return nil
}

// GraphItemUpdateInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/graphitem/object
type GraphItemUpdateInput struct {
	GItemID   *string `json:"gitemid,omitempty"`
	ItemID    *string `json:"itemid,omitempty"`
	DrawType  *int    `json:"drawtype,omitempty"`
	SortOrder *int    `json:"sortorder,omitempty"`
	Color     *string `json:"color,omitempty"`
	YAxisSide *int    `json:"yaxisside,omitempty"`
	CalcFnc   *int    `json:"calc_fnc,omitempty"`
	Type      *int    `json:"type,omitempty"`
}

func (p GraphItemUpdateInput) validate() error {
	if p.GItemID == nil && (p.ItemID == nil || *p.ItemID == "" || p.Color == nil || *p.Color == "") {
		return errors.New("zabbix: new graph items require itemid and color")
	}
	return nil
}

// HostPrototypeGroupInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostprototype/object
type HostPrototypeGroupInput struct {
	Name string `json:"name"`
}

func (p HostPrototypeGroupInput) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: HostPrototypeGroupInput.Name is required")
	}
	return nil
}

// HostPrototypeGroupUpdateInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostprototype/object
type HostPrototypeGroupUpdateInput struct {
	GroupPrototypeID *string `json:"group_prototypeid,omitempty"`
	Name             *string `json:"name,omitempty"`
}

func (p HostPrototypeGroupUpdateInput) validate() error {
	return nil
}

// HostPrototypeInterfaceInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/hostprototype/object
type HostPrototypeInterfaceInput struct {
	Main    int                 `json:"main"`
	Type    int                 `json:"type"`
	UseIP   int                 `json:"useip"`
	IP      *string             `json:"ip,omitempty"`
	DNS     *string             `json:"dns,omitempty"`
	Port    string              `json:"port"`
	Details *SNMPInterfaceInput `json:"details,omitempty"`
}

func (p HostPrototypeInterfaceInput) validate() error {
	if p.Port == "" {
		return errors.New("zabbix: HostPrototypeInterfaceInput.Port is required")
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
	return nil
}

// DiscoveryCheckInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dcheck/object
type DiscoveryCheckInput struct {
	Type                 int     `json:"type"`
	Key                  *string `json:"key_,omitempty"`
	SNMPCommunity        *string `json:"snmp_community,omitempty"`
	Ports                *string `json:"ports,omitempty"`
	SNMPv3SecurityName   *string `json:"snmpv3_securityname,omitempty"`
	SNMPv3SecurityLevel  *int    `json:"snmpv3_securitylevel,omitempty"`
	SNMPv3AuthPassphrase *string `json:"snmpv3_authpassphrase,omitempty"`
	SNMPv3PrivPassphrase *string `json:"snmpv3_privpassphrase,omitempty"`
	SNMPv3AuthProtocol   *int    `json:"snmpv3_authprotocol,omitempty"`
	SNMPv3PrivProtocol   *int    `json:"snmpv3_privprotocol,omitempty"`
	SNMPv3ContextName    *string `json:"snmpv3_contextname,omitempty"`
	Uniq                 *int    `json:"uniq,omitempty"`
	HostSource           *int    `json:"host_source,omitempty"`
	NameSource           *int    `json:"name_source,omitempty"`
	AllowRedirect        *int    `json:"allow_redirect,omitempty"`
}

func (p DiscoveryCheckInput) validate() error {
	return nil
}

// DiscoveryCheckUpdateInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dcheck/object
type DiscoveryCheckUpdateInput struct {
	DCheckID             *string `json:"dcheckid,omitempty"`
	Type                 int     `json:"type"`
	Key                  *string `json:"key_,omitempty"`
	SNMPCommunity        *string `json:"snmp_community,omitempty"`
	Ports                *string `json:"ports,omitempty"`
	SNMPv3SecurityName   *string `json:"snmpv3_securityname,omitempty"`
	SNMPv3SecurityLevel  *int    `json:"snmpv3_securitylevel,omitempty"`
	SNMPv3AuthPassphrase *string `json:"snmpv3_authpassphrase,omitempty"`
	SNMPv3PrivPassphrase *string `json:"snmpv3_privpassphrase,omitempty"`
	SNMPv3AuthProtocol   *int    `json:"snmpv3_authprotocol,omitempty"`
	SNMPv3PrivProtocol   *int    `json:"snmpv3_privprotocol,omitempty"`
	SNMPv3ContextName    *string `json:"snmpv3_contextname,omitempty"`
	Uniq                 *int    `json:"uniq,omitempty"`
	HostSource           *int    `json:"host_source,omitempty"`
	NameSource           *int    `json:"name_source,omitempty"`
	AllowRedirect        *int    `json:"allow_redirect,omitempty"`
}

func (p DiscoveryCheckUpdateInput) validate() error {
	return nil
}

// ItemPreprocessingInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/item/object
type ItemPreprocessingInput struct {
	Type               int     `json:"type"`
	Params             *string `json:"params,omitempty"`
	ErrorHandler       *int    `json:"error_handler,omitempty"`
	ErrorHandlerParams *string `json:"error_handler_params,omitempty"`
}

func (p ItemPreprocessingInput) validate() error {
	return nil
}

// DiscoveryPreprocessingInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/discoveryrule/object
type DiscoveryPreprocessingInput struct {
	Type               int     `json:"type"`
	Params             *string `json:"params,omitempty"`
	ErrorHandler       *int    `json:"error_handler,omitempty"`
	ErrorHandlerParams *string `json:"error_handler_params,omitempty"`
}

func (p DiscoveryPreprocessingInput) validate() error {
	return nil
}

// LLDMacroPathInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/discoveryrule/object
type LLDMacroPathInput struct {
	LLDMacro string `json:"lld_macro"`
	Path     string `json:"path"`
}

func (p LLDMacroPathInput) validate() error {
	if p.LLDMacro == "" {
		return errors.New("zabbix: LLDMacroPathInput.LLDMacro is required")
	}
	if p.Path == "" {
		return errors.New("zabbix: LLDMacroPathInput.Path is required")
	}
	return nil
}

// DiscoveryFilterInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/discoveryrule/object
type DiscoveryFilterInput struct {
	Conditions []DiscoveryConditionInput `json:"conditions"`
	EvalType   int                       `json:"evaltype"`
	Formula    *string                   `json:"formula,omitempty"`
}

func (p DiscoveryFilterInput) validate() error {
	if p.Conditions == nil {
		return errors.New("zabbix: DiscoveryFilterInput.Conditions is required")
	}
	for _, value := range p.Conditions {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// DiscoveryConditionInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/discoveryrule/object
type DiscoveryConditionInput struct {
	Macro     string  `json:"macro"`
	Value     *string `json:"value,omitempty"`
	FormulaID *string `json:"formulaid,omitempty"`
	Operator  *int    `json:"operator,omitempty"`
}

func (p DiscoveryConditionInput) validate() error {
	if p.Macro == "" {
		return errors.New("zabbix: DiscoveryConditionInput.Macro is required")
	}
	return nil
}

// UserGroupRightInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usergroup/object
type UserGroupRightInput struct {
	Permission int    `json:"permission"`
	ID         string `json:"id"`
}

func (p UserGroupRightInput) validate() error {
	if err := requireWriteID(p.ID, "UserGroupRightInput.ID"); err != nil {
		return err
	}
	return nil
}

// UserGroupTagInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/usergroup/object
type UserGroupTagInput struct {
	GroupID string  `json:"groupid"`
	Tag     *string `json:"tag,omitempty"`
	Value   *string `json:"value,omitempty"`
}

func (p UserGroupTagInput) validate() error {
	if err := requireWriteID(p.GroupID, "UserGroupTagInput.GroupID"); err != nil {
		return err
	}
	return nil
}

// UserMediaInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/user/object
type UserMediaInput struct {
	MediaTypeID string  `json:"mediatypeid"`
	SendTo      any     `json:"sendto"`
	Active      *int    `json:"active,omitempty"`
	Severity    *int    `json:"severity,omitempty"`
	Period      *string `json:"period,omitempty"`
}

func (p UserMediaInput) validate() error {
	if err := requireWriteID(p.MediaTypeID, "UserMediaInput.MediaTypeID"); err != nil {
		return err
	}
	if p.SendTo == nil {
		return errors.New("zabbix: UserMediaInput.SendTo is required")
	}
	if !validMediaRecipient(p.SendTo) {
		return errors.New("zabbix: media sendto must be a string or string slice")
	}
	return nil
}

// ValueMappingInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/valuemap/object
type ValueMappingInput struct {
	Type     *int    `json:"type,omitempty"`
	Value    *string `json:"value,omitempty"`
	NewValue string  `json:"newvalue"`
}

func (p ValueMappingInput) validate() error {
	if p.NewValue == "" {
		return errors.New("zabbix: ValueMappingInput.NewValue is required")
	}
	return nil
}

// ServiceStatusRuleInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/service/object
type ServiceStatusRuleInput struct {
	Type        int `json:"type"`
	LimitValue  int `json:"limit_value"`
	LimitStatus int `json:"limit_status"`
	NewStatus   int `json:"new_status"`
}

func (p ServiceStatusRuleInput) validate() error {
	return nil
}

// ServiceProblemTagInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/service/object
type ServiceProblemTagInput struct {
	Tag      string  `json:"tag"`
	Operator *int    `json:"operator,omitempty"`
	Value    *string `json:"value,omitempty"`
}

func (p ServiceProblemTagInput) validate() error {
	if p.Tag == "" {
		return errors.New("zabbix: ServiceProblemTagInput.Tag is required")
	}
	return nil
}

// MaintenancePeriodInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/maintenance/object
type MaintenancePeriodInput struct {
	TimeperiodType *int   `json:"timeperiod_type,omitempty"`
	Every          *int   `json:"every,omitempty"`
	Month          *int   `json:"month,omitempty"`
	DayOfWeek      *int   `json:"dayofweek,omitempty"`
	Day            *int   `json:"day,omitempty"`
	StartTime      *int   `json:"start_time,omitempty"`
	Period         *int   `json:"period,omitempty"`
	StartDate      *int64 `json:"start_date,omitempty"`
}

func (p MaintenancePeriodInput) validate() error {
	return nil
}

// MaintenanceTagInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/maintenance/object
type MaintenanceTagInput struct {
	Tag      string  `json:"tag"`
	Operator *int    `json:"operator,omitempty"`
	Value    *string `json:"value,omitempty"`
}

func (p MaintenanceTagInput) validate() error {
	if p.Tag == "" {
		return errors.New("zabbix: MaintenanceTagInput.Tag is required")
	}
	return nil
}

// ConnectorTagInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/connector/object
type ConnectorTagInput struct {
	Tag      string  `json:"tag"`
	Operator *int    `json:"operator,omitempty"`
	Value    *string `json:"value,omitempty"`
}

func (p ConnectorTagInput) validate() error {
	if p.Tag == "" {
		return errors.New("zabbix: ConnectorTagInput.Tag is required")
	}
	return nil
}

// SLATagInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/sla/object
type SLATagInput struct {
	Tag      string  `json:"tag"`
	Operator *int    `json:"operator,omitempty"`
	Value    *string `json:"value,omitempty"`
}

func (p SLATagInput) validate() error {
	if p.Tag == "" {
		return errors.New("zabbix: SLATagInput.Tag is required")
	}
	return nil
}

// SLAScheduleInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/sla/object
type SLAScheduleInput struct {
	PeriodFrom int `json:"period_from"`
	PeriodTo   int `json:"period_to"`
}

func (p SLAScheduleInput) validate() error {
	return nil
}

// SLAExcludedDowntimeInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/sla/object
type SLAExcludedDowntimeInput struct {
	Name       string `json:"name"`
	PeriodFrom int    `json:"period_from"`
	PeriodTo   int    `json:"period_to"`
}

func (p SLAExcludedDowntimeInput) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: SLAExcludedDowntimeInput.Name is required")
	}
	return nil
}

// RegexpExpressionInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/regexp/object
type RegexpExpressionInput struct {
	Expression     string  `json:"expression"`
	ExpressionType int     `json:"expression_type"`
	ExpDelimiter   *string `json:"exp_delimiter,omitempty"`
	CaseSensitive  *int    `json:"case_sensitive,omitempty"`
}

func (p RegexpExpressionInput) validate() error {
	if p.Expression == "" {
		return errors.New("zabbix: RegexpExpressionInput.Expression is required")
	}
	return nil
}

// IconMappingInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/iconmap/object
type IconMappingInput struct {
	IconID        string `json:"iconid"`
	InventoryLink int    `json:"inventory_link"`
	Expression    string `json:"expression"`
}

func (p IconMappingInput) validate() error {
	if err := requireWriteID(p.IconID, "IconMappingInput.IconID"); err != nil {
		return err
	}
	if p.Expression == "" {
		return errors.New("zabbix: IconMappingInput.Expression is required")
	}
	return nil
}

// ReportUserInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/report/object
type ReportUserInput struct {
	UserID       string  `json:"userid"`
	AccessUserID *string `json:"access_userid,omitempty"`
	Exclude      *int    `json:"exclude,omitempty"`
}

func (p ReportUserInput) validate() error {
	if err := requireWriteID(p.UserID, "ReportUserInput.UserID"); err != nil {
		return err
	}
	return nil
}

// ReportUserGroupInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/report/object
type ReportUserGroupInput struct {
	UsrGrpID     string  `json:"usrgrpid"`
	AccessUserID *string `json:"access_userid,omitempty"`
}

func (p ReportUserGroupInput) validate() error {
	if err := requireWriteID(p.UsrGrpID, "ReportUserGroupInput.UsrGrpID"); err != nil {
		return err
	}
	return nil
}

// ScriptParameterInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/object
type ScriptParameterInput struct {
	Name  string  `json:"name"`
	Value *string `json:"value,omitempty"`
}

func (p ScriptParameterInput) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: ScriptParameterInput.Name is required")
	}
	return nil
}

// MediaMessageTemplateInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mediatype/object
type MediaMessageTemplateInput struct {
	EventSource int     `json:"eventsource"`
	Recovery    int     `json:"recovery"`
	Subject     *string `json:"subject,omitempty"`
	Message     *string `json:"message,omitempty"`
}

func (p MediaMessageTemplateInput) validate() error {
	return nil
}

// MediaParameterInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mediatype/object
type MediaParameterInput struct {
	Name      *string `json:"name,omitempty"`
	Value     string  `json:"value"`
	SortOrder *int    `json:"sortorder,omitempty"`
}

func (p MediaParameterInput) validate() error {
	if (p.Name == nil) == (p.SortOrder == nil) {
		return errors.New("zabbix: media parameters require exactly one name or sortorder")
	}
	return nil
}

// DashboardPageInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dashboard/object
type DashboardPageInput struct {
	Name          *string                `json:"name,omitempty"`
	DisplayPeriod *int                   `json:"display_period,omitempty"`
	Widgets       []DashboardWidgetInput `json:"widgets,omitzero"`
}

func (p DashboardPageInput) validate() error {
	for _, value := range p.Widgets {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// DashboardPageUpdateInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dashboard/object
type DashboardPageUpdateInput struct {
	DashboardPageID *string                      `json:"dashboard_pageid,omitempty"`
	Name            *string                      `json:"name,omitempty"`
	DisplayPeriod   *int                         `json:"display_period,omitempty"`
	Widgets         []DashboardWidgetUpdateInput `json:"widgets,omitzero"`
}

func (p DashboardPageUpdateInput) validate() error {
	for _, value := range p.Widgets {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// DashboardWidgetInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dashboard/object
type DashboardWidgetInput struct {
	Type     string             `json:"type"`
	Name     *string            `json:"name,omitempty"`
	X        *int               `json:"x,omitempty"`
	Y        *int               `json:"y,omitempty"`
	Width    *int               `json:"width,omitempty"`
	Height   *int               `json:"height,omitempty"`
	ViewMode *int               `json:"view_mode,omitempty"`
	Fields   []WidgetFieldInput `json:"fields,omitzero"`
}

func (p DashboardWidgetInput) validate() error {
	if p.Type == "" {
		return errors.New("zabbix: DashboardWidgetInput.Type is required")
	}
	for _, value := range p.Fields {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// DashboardWidgetUpdateInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dashboard/object
type DashboardWidgetUpdateInput struct {
	WidgetID *string            `json:"widgetid,omitempty"`
	Type     *string            `json:"type,omitempty"`
	Name     *string            `json:"name,omitempty"`
	X        *int               `json:"x,omitempty"`
	Y        *int               `json:"y,omitempty"`
	Width    *int               `json:"width,omitempty"`
	Height   *int               `json:"height,omitempty"`
	ViewMode *int               `json:"view_mode,omitempty"`
	Fields   []WidgetFieldInput `json:"fields,omitzero"`
}

func (p DashboardWidgetUpdateInput) validate() error {
	for _, value := range p.Fields {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if p.WidgetID == nil && (p.Type == nil || *p.Type == "") {
		return errors.New("zabbix: new widgets require type")
	}
	return nil
}

// TemplateDashboardPageInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templatedashboard/object
type TemplateDashboardPageInput struct {
	Name          *string                        `json:"name,omitempty"`
	DisplayPeriod *int                           `json:"display_period,omitempty"`
	Widgets       []TemplateDashboardWidgetInput `json:"widgets,omitzero"`
}

func (p TemplateDashboardPageInput) validate() error {
	for _, value := range p.Widgets {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// TemplateDashboardPageUpdateInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templatedashboard/object
type TemplateDashboardPageUpdateInput struct {
	DashboardPageID *string                              `json:"dashboard_pageid,omitempty"`
	Name            *string                              `json:"name,omitempty"`
	DisplayPeriod   *int                                 `json:"display_period,omitempty"`
	Widgets         []TemplateDashboardWidgetUpdateInput `json:"widgets,omitzero"`
}

func (p TemplateDashboardPageUpdateInput) validate() error {
	for _, value := range p.Widgets {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// TemplateDashboardWidgetInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templatedashboard/object
type TemplateDashboardWidgetInput struct {
	Type     string             `json:"type"`
	Name     *string            `json:"name,omitempty"`
	X        *int               `json:"x,omitempty"`
	Y        *int               `json:"y,omitempty"`
	Width    *int               `json:"width,omitempty"`
	Height   *int               `json:"height,omitempty"`
	ViewMode *int               `json:"view_mode,omitempty"`
	Fields   []WidgetFieldInput `json:"fields,omitzero"`
}

func (p TemplateDashboardWidgetInput) validate() error {
	if p.Type == "" {
		return errors.New("zabbix: TemplateDashboardWidgetInput.Type is required")
	}
	for _, value := range p.Fields {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// TemplateDashboardWidgetUpdateInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/templatedashboard/object
type TemplateDashboardWidgetUpdateInput struct {
	WidgetID *string            `json:"widgetid,omitempty"`
	Type     *string            `json:"type,omitempty"`
	Name     *string            `json:"name,omitempty"`
	X        *int               `json:"x,omitempty"`
	Y        *int               `json:"y,omitempty"`
	Width    *int               `json:"width,omitempty"`
	Height   *int               `json:"height,omitempty"`
	ViewMode *int               `json:"view_mode,omitempty"`
	Fields   []WidgetFieldInput `json:"fields,omitzero"`
}

func (p TemplateDashboardWidgetUpdateInput) validate() error {
	for _, value := range p.Fields {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if p.WidgetID == nil && (p.Type == nil || *p.Type == "") {
		return errors.New("zabbix: new widgets require type")
	}
	return nil
}

// WidgetFieldInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dashboard/object
type WidgetFieldInput struct {
	Type  int    `json:"type"`
	Name  string `json:"name"`
	Value any    `json:"value"`
}

func (p WidgetFieldInput) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: WidgetFieldInput.Name is required")
	}
	if p.Value == nil {
		return errors.New("zabbix: WidgetFieldInput.Value is required")
	}
	return nil
}

// DashboardUserInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dashboard/object
type DashboardUserInput struct {
	UserID     string `json:"userid"`
	Permission int    `json:"permission"`
}

func (p DashboardUserInput) validate() error {
	if err := requireWriteID(p.UserID, "DashboardUserInput.UserID"); err != nil {
		return err
	}
	return nil
}

// DashboardUserGroupInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/dashboard/object
type DashboardUserGroupInput struct {
	UsrGrpID   string `json:"usrgrpid"`
	Permission int    `json:"permission"`
}

func (p DashboardUserGroupInput) validate() error {
	if err := requireWriteID(p.UsrGrpID, "DashboardUserGroupInput.UsrGrpID"); err != nil {
		return err
	}
	return nil
}

// CorrelationFilterInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/correlation/object
type CorrelationFilterInput struct {
	EvalType   int                         `json:"evaltype"`
	Formula    *string                     `json:"formula,omitempty"`
	Conditions []CorrelationConditionInput `json:"conditions"`
}

func (p CorrelationFilterInput) validate() error {
	if len(p.Conditions) == 0 {
		return errors.New("zabbix: CorrelationFilterInput.Conditions must not be empty")
	}
	for _, value := range p.Conditions {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// CorrelationConditionInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/correlation/object
type CorrelationConditionInput struct {
	Type      int     `json:"type"`
	Tag       *string `json:"tag,omitempty"`
	OldTag    *string `json:"oldtag,omitempty"`
	NewTag    *string `json:"newtag,omitempty"`
	Value     *string `json:"value,omitempty"`
	FormulaID *string `json:"formulaid,omitempty"`
	Operator  *int    `json:"operator,omitempty"`
	GroupID   *string `json:"groupid,omitempty"`
}

func (p CorrelationConditionInput) validate() error {
	return nil
}

// CorrelationOperationInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/correlation/object
type CorrelationOperationInput struct {
	Type int `json:"type"`
}

func (p CorrelationOperationInput) validate() error {
	return nil
}

// ActionFilterInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/action/object
type ActionFilterInput struct {
	EvalType   int                    `json:"evaltype"`
	Formula    *string                `json:"formula,omitempty"`
	Conditions []ActionConditionInput `json:"conditions"`
}

func (p ActionFilterInput) validate() error {
	if p.Conditions == nil {
		return errors.New("zabbix: ActionFilterInput.Conditions is required")
	}
	for _, value := range p.Conditions {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ActionConditionInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/action/object
type ActionConditionInput struct {
	ConditionType int     `json:"conditiontype"`
	Operator      *int    `json:"operator,omitempty"`
	Value         string  `json:"value"`
	Value2        *string `json:"value2,omitempty"`
	FormulaID     *string `json:"formulaid,omitempty"`
}

func (p ActionConditionInput) validate() error {
	return nil
}

// ActionMessageInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/action/object
type ActionMessageInput struct {
	DefaultMsg  *int    `json:"default_msg,omitempty"`
	MediaTypeID *string `json:"mediatypeid,omitempty"`
	Subject     *string `json:"subject,omitempty"`
	Message     *string `json:"message,omitempty"`
}

func (p ActionMessageInput) validate() error {
	return nil
}

// ActionOperationConditionInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/action/object
type ActionOperationConditionInput struct {
	ConditionType int    `json:"conditiontype"`
	Value         string `json:"value"`
	Operator      *int   `json:"operator,omitempty"`
}

func (p ActionOperationConditionInput) validate() error {
	return nil
}

// ActionHostInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/action/object
type ActionHostInput struct {
	HostID string `json:"hostid"`
}

func (p ActionHostInput) validate() error {
	if p.HostID != "0" {
		if err := requireWriteID(p.HostID, "ActionHostInput.HostID"); err != nil {
			return err
		}
	}
	return nil
}

// ActionOperationInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/action/object
type ActionOperationInput struct {
	OperationType int                             `json:"operationtype"`
	EscPeriod     *string                         `json:"esc_period,omitempty"`
	EscStepFrom   *int                            `json:"esc_step_from,omitempty"`
	EscStepTo     *int                            `json:"esc_step_to,omitempty"`
	EvalType      *int                            `json:"evaltype,omitempty"`
	OpCommand     *ScriptRef                      `json:"opcommand,omitempty"`
	OpCommandGrp  []HostGroupRef                  `json:"opcommand_grp,omitzero"`
	OpCommandHst  []ActionHostInput               `json:"opcommand_hst,omitzero"`
	OpConditions  []ActionOperationConditionInput `json:"opconditions,omitzero"`
	OpGroup       []HostGroupRef                  `json:"opgroup,omitzero"`
	OpMessage     *ActionMessageInput             `json:"opmessage,omitempty"`
	OpMessageGrp  []UserGroupRef                  `json:"opmessage_grp,omitzero"`
	OpMessageUsr  []UserRef                       `json:"opmessage_usr,omitzero"`
	OpTemplate    []TemplateRef                   `json:"optemplate,omitzero"`
	OpInventory   *ActionInventoryInput           `json:"opinventory,omitempty"`
	OpTag         []TagInput                      `json:"optag,omitzero"`
}

func (p ActionOperationInput) validate() error {
	if p.OpCommand != nil {
		if err := p.OpCommand.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.OpCommandGrp {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.OpCommandHst {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.OpConditions {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.OpGroup {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if p.OpMessage != nil {
		if err := p.OpMessage.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.OpMessageGrp {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.OpMessageUsr {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.OpTemplate {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if p.OpInventory != nil {
		if err := p.OpInventory.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.OpTag {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if p.OperationType == 0 && (p.OpMessage == nil || len(p.OpMessageGrp)+len(p.OpMessageUsr) == 0) {
		return errors.New("zabbix: message operations require message settings and recipients")
	}
	if p.OperationType == 1 && p.OpCommand == nil {
		return errors.New("zabbix: script operations require a script")
	}
	if (p.OperationType == 4 || p.OperationType == 5) && len(p.OpGroup) == 0 {
		return errors.New("zabbix: host group operations require groups")
	}
	if (p.OperationType == 6 || p.OperationType == 7) && len(p.OpTemplate) == 0 {
		return errors.New("zabbix: template operations require templates")
	}
	if p.OperationType == 10 && p.OpInventory == nil {
		return errors.New("zabbix: inventory operations require inventory mode")
	}
	if p.OperationType < 0 || p.OperationType == 11 || p.OperationType == 12 || p.OperationType > 14 {
		return errors.New("zabbix: unsupported action operation type")
	}
	return nil
}

// ActionRecoveryOperationInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/action/object
type ActionRecoveryOperationInput struct {
	OperationType int                 `json:"operationtype"`
	OpCommand     *ScriptRef          `json:"opcommand,omitempty"`
	OpCommandGrp  []HostGroupRef      `json:"opcommand_grp,omitzero"`
	OpCommandHst  []ActionHostInput   `json:"opcommand_hst,omitzero"`
	OpMessage     *ActionMessageInput `json:"opmessage,omitempty"`
	OpMessageGrp  []UserGroupRef      `json:"opmessage_grp,omitzero"`
	OpMessageUsr  []UserRef           `json:"opmessage_usr,omitzero"`
}

func (p ActionRecoveryOperationInput) validate() error {
	if p.OpCommand != nil {
		if err := p.OpCommand.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.OpCommandGrp {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.OpCommandHst {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if p.OpMessage != nil {
		if err := p.OpMessage.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.OpMessageGrp {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.OpMessageUsr {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if p.OperationType == 0 && (p.OpMessage == nil || len(p.OpMessageGrp)+len(p.OpMessageUsr) == 0) {
		return errors.New("zabbix: message operations require message settings and recipients")
	}
	if p.OperationType == 1 && p.OpCommand == nil {
		return errors.New("zabbix: script operations require a script")
	}
	if p.OperationType != 0 && p.OperationType != 1 && p.OperationType != 11 {
		return errors.New("zabbix: unsupported recovery operation type")
	}
	return nil
}

// ActionUpdateOperationInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/action/object
type ActionUpdateOperationInput struct {
	OperationType int                 `json:"operationtype"`
	OpCommand     *ScriptRef          `json:"opcommand,omitempty"`
	OpCommandGrp  []HostGroupRef      `json:"opcommand_grp,omitzero"`
	OpCommandHst  []ActionHostInput   `json:"opcommand_hst,omitzero"`
	OpMessage     *ActionMessageInput `json:"opmessage,omitempty"`
	OpMessageGrp  []UserGroupRef      `json:"opmessage_grp,omitzero"`
	OpMessageUsr  []UserRef           `json:"opmessage_usr,omitzero"`
}

func (p ActionUpdateOperationInput) validate() error {
	if p.OpCommand != nil {
		if err := p.OpCommand.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.OpCommandGrp {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.OpCommandHst {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if p.OpMessage != nil {
		if err := p.OpMessage.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.OpMessageGrp {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.OpMessageUsr {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if p.OperationType == 0 && (p.OpMessage == nil || len(p.OpMessageGrp)+len(p.OpMessageUsr) == 0) {
		return errors.New("zabbix: message operations require message settings and recipients")
	}
	if p.OperationType == 1 && p.OpCommand == nil {
		return errors.New("zabbix: script operations require a script")
	}
	if p.OperationType != 0 && p.OperationType != 1 && p.OperationType != 12 {
		return errors.New("zabbix: unsupported update operation type")
	}
	return nil
}

// ActionInventoryInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/host/object
type ActionInventoryInput struct {
	InventoryMode int `json:"inventory_mode"`
}

func (p ActionInventoryInput) validate() error {
	return nil
}

// MapElementInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/map/object
type MapElementInput struct {
	SelementID     *string         `json:"selementid,omitempty"`
	ElementType    int             `json:"elementtype"`
	IconIDOff      string          `json:"iconid_off"`
	IconIDOn       *string         `json:"iconid_on,omitempty"`
	Label          *string         `json:"label,omitempty"`
	LabelLocation  *int            `json:"label_location,omitempty"`
	X              *int            `json:"x,omitempty"`
	Y              *int            `json:"y,omitempty"`
	ElementSubtype *int            `json:"elementsubtype,omitempty"`
	AreaType       *int            `json:"areatype,omitempty"`
	Width          *int            `json:"width,omitempty"`
	Height         *int            `json:"height,omitempty"`
	ViewType       *int            `json:"viewtype,omitempty"`
	UseIconMap     *int            `json:"use_iconmap,omitempty"`
	EvalType       *int            `json:"evaltype,omitempty"`
	Elements       []MapElementRef `json:"elements,omitzero"`
}

func (p MapElementInput) validate() error {
	if err := requireWriteID(p.IconIDOff, "MapElementInput.IconIDOff"); err != nil {
		return err
	}
	for _, value := range p.Elements {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// MapLinkInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/map/object
type MapLinkInput struct {
	LinkID      *string `json:"linkid,omitempty"`
	SelementID1 string  `json:"selementid1"`
	SelementID2 string  `json:"selementid2"`
	DrawType    *int    `json:"drawtype,omitempty"`
	Color       *string `json:"color,omitempty"`
	Label       *string `json:"label,omitempty"`
}

func (p MapLinkInput) validate() error {
	if p.SelementID1 == "" {
		return errors.New("zabbix: MapLinkInput.SelementID1 is required")
	}
	if p.SelementID2 == "" {
		return errors.New("zabbix: MapLinkInput.SelementID2 is required")
	}
	return nil
}

// MapURLInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/map/object
type MapURLInput struct {
	Name        string `json:"name"`
	URL         string `json:"url"`
	ElementType *int   `json:"elementtype,omitempty"`
}

func (p MapURLInput) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: MapURLInput.Name is required")
	}
	if p.URL == "" {
		return errors.New("zabbix: MapURLInput.URL is required")
	}
	return nil
}

// MapUserInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/map/object
type MapUserInput struct {
	UserID     string `json:"userid"`
	Permission int    `json:"permission"`
}

func (p MapUserInput) validate() error {
	if err := requireWriteID(p.UserID, "MapUserInput.UserID"); err != nil {
		return err
	}
	return nil
}

// MapUserGroupInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/map/object
type MapUserGroupInput struct {
	UsrGrpID   string `json:"usrgrpid"`
	Permission int    `json:"permission"`
}

func (p MapUserGroupInput) validate() error {
	if err := requireWriteID(p.UsrGrpID, "MapUserGroupInput.UsrGrpID"); err != nil {
		return err
	}
	return nil
}

// MapShapeInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/map/object
type MapShapeInput struct {
	Type            int     `json:"type"`
	X               *int    `json:"x,omitempty"`
	Y               *int    `json:"y,omitempty"`
	Width           *int    `json:"width,omitempty"`
	Height          *int    `json:"height,omitempty"`
	Text            *string `json:"text,omitempty"`
	Font            *int    `json:"font,omitempty"`
	FontSize        *int    `json:"font_size,omitempty"`
	FontColor       *string `json:"font_color,omitempty"`
	TextHalign      *int    `json:"text_halign,omitempty"`
	TextValign      *int    `json:"text_valign,omitempty"`
	BorderType      *int    `json:"border_type,omitempty"`
	BorderWidth     *int    `json:"border_width,omitempty"`
	BorderColor     *string `json:"border_color,omitempty"`
	BackgroundColor *string `json:"background_color,omitempty"`
	ZIndex          *int    `json:"zindex,omitempty"`
}

func (p MapShapeInput) validate() error {
	return nil
}

// MapLineInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/map/object
type MapLineInput struct {
	X1        *int    `json:"x1,omitempty"`
	Y1        *int    `json:"y1,omitempty"`
	X2        *int    `json:"x2,omitempty"`
	Y2        *int    `json:"y2,omitempty"`
	LineType  *int    `json:"line_type,omitempty"`
	LineWidth *int    `json:"line_width,omitempty"`
	LineColor *string `json:"line_color,omitempty"`
	ZIndex    *int    `json:"zindex,omitempty"`
}

func (p MapLineInput) validate() error {
	return nil
}

// MapShapeUpdateInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/map/object
type MapShapeUpdateInput struct {
	SysmapShapeID   *string `json:"sysmap_shapeid,omitempty"`
	Type            *int    `json:"type,omitempty"`
	X               *int    `json:"x,omitempty"`
	Y               *int    `json:"y,omitempty"`
	Width           *int    `json:"width,omitempty"`
	Height          *int    `json:"height,omitempty"`
	Text            *string `json:"text,omitempty"`
	Font            *int    `json:"font,omitempty"`
	FontSize        *int    `json:"font_size,omitempty"`
	FontColor       *string `json:"font_color,omitempty"`
	TextHalign      *int    `json:"text_halign,omitempty"`
	TextValign      *int    `json:"text_valign,omitempty"`
	BorderType      *int    `json:"border_type,omitempty"`
	BorderWidth     *int    `json:"border_width,omitempty"`
	BorderColor     *string `json:"border_color,omitempty"`
	BackgroundColor *string `json:"background_color,omitempty"`
	ZIndex          *int    `json:"zindex,omitempty"`
}

func (p MapShapeUpdateInput) validate() error {
	return nil
}

// MapLineUpdateInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/map/object
type MapLineUpdateInput struct {
	SysmapShapeID *string `json:"sysmap_shapeid,omitempty"`
	X1            *int    `json:"x1,omitempty"`
	Y1            *int    `json:"y1,omitempty"`
	X2            *int    `json:"x2,omitempty"`
	Y2            *int    `json:"y2,omitempty"`
	LineType      *int    `json:"line_type,omitempty"`
	LineWidth     *int    `json:"line_width,omitempty"`
	LineColor     *string `json:"line_color,omitempty"`
	ZIndex        *int    `json:"zindex,omitempty"`
}

func (p MapLineUpdateInput) validate() error {
	return nil
}

// MapElementRef contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/map/object
type MapElementRef struct {
	HostID    *string `json:"hostid,omitempty"`
	GroupID   *string `json:"groupid,omitempty"`
	SysmapID  *string `json:"sysmapid,omitempty"`
	TriggerID *string `json:"triggerid,omitempty"`
}

func (p MapElementRef) validate() error {
	return nil
}

// RoleRulesInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/role/object
type RoleRulesInput struct {
	UI                   []RoleUIInput     `json:"ui,omitzero"`
	UIDefaultAccess      *int              `json:"ui.default_access,omitempty"`
	ServicesReadMode     *int              `json:"services.read.mode,omitempty"`
	ServicesReadList     []ServiceRef      `json:"services.read.list,omitzero"`
	ServicesReadTag      *TagInput         `json:"services.read.tag,omitempty"`
	ServicesWriteMode    *int              `json:"services.write.mode,omitempty"`
	ServicesWriteList    []ServiceRef      `json:"services.write.list,omitzero"`
	ServicesWriteTag     *TagInput         `json:"services.write.tag,omitempty"`
	Modules              []RoleModuleInput `json:"modules,omitzero"`
	ModulesDefaultAccess *int              `json:"modules.default_access,omitempty"`
	APIAccess            *int              `json:"api.access,omitempty"`
	APIMode              *int              `json:"api.mode,omitempty"`
	API                  []string          `json:"api,omitzero"`
	Actions              []RoleActionInput `json:"actions,omitzero"`
	ActionsDefaultAccess *int              `json:"actions.default_access,omitempty"`
}

func (p RoleRulesInput) validate() error {
	for _, value := range p.UI {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.ServicesReadList {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if p.ServicesReadTag != nil {
		if err := p.ServicesReadTag.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.ServicesWriteList {
		if err := value.validate(); err != nil {
			return err
		}
	}
	if p.ServicesWriteTag != nil {
		if err := p.ServicesWriteTag.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Modules {
		if err := value.validate(); err != nil {
			return err
		}
	}
	for _, value := range p.Actions {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// RoleUIInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/role/object
type RoleUIInput struct {
	Name   string `json:"name"`
	Status *int   `json:"status,omitempty"`
}

func (p RoleUIInput) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: RoleUIInput.Name is required")
	}
	return nil
}

// RoleModuleInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/role/object
type RoleModuleInput struct {
	ModuleID string `json:"moduleid"`
	Status   *int   `json:"status,omitempty"`
}

func (p RoleModuleInput) validate() error {
	if err := requireWriteID(p.ModuleID, "RoleModuleInput.ModuleID"); err != nil {
		return err
	}
	return nil
}

// RoleActionInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/role/object
type RoleActionInput struct {
	Name   string `json:"name"`
	Status *int   `json:"status,omitempty"`
}

func (p RoleActionInput) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: RoleActionInput.Name is required")
	}
	return nil
}

// ProvisioningGroupInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/userdirectory/object
type ProvisioningGroupInput struct {
	Name       string         `json:"name"`
	RoleID     string         `json:"roleid"`
	UserGroups []UserGroupRef `json:"user_groups"`
}

func (p ProvisioningGroupInput) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: ProvisioningGroupInput.Name is required")
	}
	if err := requireWriteID(p.RoleID, "ProvisioningGroupInput.RoleID"); err != nil {
		return err
	}
	if len(p.UserGroups) == 0 {
		return errors.New("zabbix: ProvisioningGroupInput.UserGroups must not be empty")
	}
	for _, value := range p.UserGroups {
		if err := value.validate(); err != nil {
			return err
		}
	}
	return nil
}

// ProvisioningMediaInput contains writable Zabbix 7.4 fields.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/userdirectory/object
type ProvisioningMediaInput struct {
	Name        string  `json:"name"`
	MediaTypeID string  `json:"mediatypeid"`
	Attribute   string  `json:"attribute"`
	Active      *int    `json:"active,omitempty"`
	Severity    *int    `json:"severity,omitempty"`
	Period      *string `json:"period,omitempty"`
}

func (p ProvisioningMediaInput) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: ProvisioningMediaInput.Name is required")
	}
	if err := requireWriteID(p.MediaTypeID, "ProvisioningMediaInput.MediaTypeID"); err != nil {
		return err
	}
	if p.Attribute == "" {
		return errors.New("zabbix: ProvisioningMediaInput.Attribute is required")
	}
	return nil
}

// requireWriteID accepts decimal wire IDs without converting through a machine integer.
func requireWriteID(id, field string) error {
	nonzero := false
	for _, c := range id {
		if c < '0' || c > '9' {
			return errors.New("zabbix: " + field + " must be a positive decimal ID")
		}
		nonzero = nonzero || c != '0'
	}
	if !nonzero {
		return errors.New("zabbix: " + field + " must be a positive decimal ID")
	}
	return nil
}

// The recipient representation is selected by the media type, not by a generic
// serialization fallback. Empty values are checked by the server's media rules.
func validMediaRecipient(value any) bool {
	switch value.(type) {
	case string, []string:
		return true
	default:
		return false
	}
}
