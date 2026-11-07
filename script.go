package zabbix

import (
	"context"
	"encoding/json"
	"errors"
)

// ScriptService wraps the "script" API namespace.
type ScriptService struct{ client *Client }

// Script returns the script service.
func (c *Client) Script() *ScriptService { return &ScriptService{c} }

// Script object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/object
type Script struct {
	ScriptID                 string `json:"scriptid,omitempty"`
	Name                     string `json:"name,omitempty"`
	Command                  string `json:"command,omitempty"`
	HostAccess               string `json:"host_access,omitempty"`
	UsrGrpID                 string `json:"usrgrpid,omitempty"`
	GroupID                  string `json:"groupid,omitempty"`
	Description              string `json:"description,omitempty"`
	Confirmation             string `json:"confirmation,omitempty"`
	Type                     string `json:"type,omitempty"`
	ExecuteOn                string `json:"execute_on,omitempty"`
	Timeout                  string `json:"timeout,omitempty"`
	Scope                    string `json:"scope,omitempty"`
	Port                     string `json:"port,omitempty"`
	AuthType                 string `json:"authtype,omitempty"`
	Username                 string `json:"username,omitempty"`
	Password                 string `json:"password,omitempty"`
	PublicKey                string `json:"publickey,omitempty"`
	PrivateKey               string `json:"privatekey,omitempty"`
	MenuPath                 string `json:"menu_path,omitempty"`
	URL                      string `json:"url,omitempty"`
	NewWindow                string `json:"new_window,omitempty"`
	ManualInput              string `json:"manualinput,omitempty"`
	ManualInputPrompt        string `json:"manualinput_prompt,omitempty"`
	ManualInputValidator     string `json:"manualinput_validator,omitempty"`
	ManualInputValidatorType string `json:"manualinput_validator_type,omitempty"`
	ManualInputDefaultValue  string `json:"manualinput_default_value,omitempty"`

	Parameters []ScriptParameter `json:"parameters,omitempty"`
	Hosts      []Host            `json:"hosts,omitempty"`
	HostGroups []HostGroup       `json:"hostgroups,omitempty"`
}

// ScriptParameter is a user-defined parameter passed to a script.
type ScriptParameter struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// ScriptGetParams are the parameters for script.get.
// ScriptGetParams queries the list form of script.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type ScriptGetParams struct {
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
	ScriptIDs              []string       `json:"scriptids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	UsrGrpIDs              []string       `json:"usrgrpids,omitzero"`
	SelectHosts            any            `json:"selectHosts,omitempty"`
	SelectHostGroups       any            `json:"selectHostGroups,omitempty"`
}

// Get retrieves scripts matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/get
func (s *ScriptService) Get(ctx context.Context, params ScriptGetParams) ([]Script, error) {
	return getList[Script](ctx, s.client, "script.get", params)
}

// ScriptCountParams contains the filters for script.get counts.
type ScriptCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	ScriptIDs              []string       `json:"scriptids,omitzero"`
	HostIDs                []string       `json:"hostids,omitzero"`
	GroupIDs               []string       `json:"groupids,omitzero"`
	UsrGrpIDs              []string       `json:"usrgrpids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *ScriptService) Count(ctx context.Context, params ScriptCountParams) (uint64, error) {
	return getCount(ctx, s.client, "script.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *ScriptService) GetMap(ctx context.Context, params ScriptGetParams) (map[string]Script, error) {
	return getMap[Script](ctx, s.client, "script.get", params)
}

// Create creates scripts and returns the new script ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/create
func (s *ScriptService) Create(ctx context.Context, items ...ScriptCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "script.create", "scriptids", items)
}

// Update updates scripts and returns the affected script ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/update
func (s *ScriptService) Update(ctx context.Context, items ...ScriptUpdateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "script.update", "scriptids", items)
}

// Delete deletes scripts by id and returns the deleted script ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/delete
func (s *ScriptService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "script.delete", "scriptids", ids)
}

// ScriptExecuteParams are the parameters for script.execute.
type ScriptExecuteParams struct {
	ScriptID    string  `json:"scriptid"`
	HostID      string  `json:"hostid,omitempty"`
	EventID     string  `json:"eventid,omitempty"`
	ManualInput *string `json:"manualinput,omitempty"`
}

// ScriptExecuteResult is the result of script.execute.
type ScriptExecuteResult struct {
	Response string `json:"response,omitempty"`
	Value    string `json:"value,omitempty"`
	Debug    any    `json:"debug,omitempty"`
}

// Execute runs a script on a host or event and returns its result.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/execute
func (s *ScriptService) Execute(ctx context.Context, params ScriptExecuteParams) (*ScriptExecuteResult, error) {
	// CScript::execute always returns response=success and a string value. An
	// empty output is valid, so pointers distinguish it from a missing field.
	var out struct {
		Response *string `json:"response"`
		Value    *string `json:"value"`
		Debug    any     `json:"debug"`
	}
	if err := s.client.Call(ctx, "script.execute", params, &out); err != nil {
		return nil, err
	}
	if out.Response == nil || *out.Response != "success" || out.Value == nil {
		return nil, requestError("script.execute", ErrorProtocol, errors.New("script result requires success response and string value"))
	}
	return &ScriptExecuteResult{Response: *out.Response, Value: *out.Value, Debug: out.Debug}, nil
}

// ScriptHostParams identifies a host and optional script/manual input to resolve.
type ScriptHostParams struct {
	HostID      string  `json:"hostid"`
	ScriptID    string  `json:"scriptid,omitempty"`
	ManualInput *string `json:"manualinput,omitempty"`
}

// ScriptEventParams identifies an event and optional script/manual input to resolve.
type ScriptEventParams struct {
	EventID     string  `json:"eventid"`
	ScriptID    string  `json:"scriptid,omitempty"`
	ManualInput *string `json:"manualinput,omitempty"`
}

// GetScriptsByHosts returns available scripts keyed by the requested host IDs.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/getscriptsbyhosts
func (s *ScriptService) GetScriptsByHosts(ctx context.Context, params []ScriptHostParams) (map[string][]Script, error) {
	if params == nil {
		params = []ScriptHostParams{}
	}
	for _, param := range params {
		if param.HostID == "" {
			return nil, requestError("script.getscriptsbyhosts", ErrorInvalidParams, errors.New("hostid is required"))
		}
	}
	return s.getScripts(ctx, "script.getscriptsbyhosts", params)
}

// GetScriptsByEvents returns available scripts keyed by the requested event IDs.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/script/getscriptsbyevents
func (s *ScriptService) GetScriptsByEvents(ctx context.Context, params []ScriptEventParams) (map[string][]Script, error) {
	if params == nil {
		params = []ScriptEventParams{}
	}
	for _, param := range params {
		if param.EventID == "" {
			return nil, requestError("script.getscriptsbyevents", ErrorInvalidParams, errors.New("eventid is required"))
		}
	}
	return s.getScripts(ctx, "script.getscriptsbyevents", params)
}

func (s *ScriptService) getScripts(ctx context.Context, method string, params any) (map[string][]Script, error) {
	var raw json.RawMessage
	if err := s.client.Call(ctx, method, params, &raw); err != nil {
		return nil, err
	}
	groups, err := decodeObjectGroups[Script](raw)
	if err != nil {
		return nil, requestError(method, ErrorDecode, err)
	}
	return groups, nil
}
