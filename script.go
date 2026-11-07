package zabbix

import "context"

// ScriptService wraps the "script" API namespace.
type ScriptService struct{ client *Client }

// Script returns the script service.
func (c *Client) Script() *ScriptService { return &ScriptService{c} }

// Script object.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/script/object
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
}

// ScriptParameter is a user-defined parameter passed to a script.
type ScriptParameter struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// ScriptGetParams are the parameters for script.get.
type ScriptGetParams struct {
	GetParams
	ScriptIDs    []string `json:"scriptids,omitempty"`
	HostIDs      []string `json:"hostids,omitempty"`
	GroupIDs     []string `json:"groupids,omitempty"`
	UsrGrpIDs    []string `json:"usrgrpids,omitempty"`
	SelectHosts  any      `json:"selectHosts,omitempty"`
	SelectGroups any      `json:"selectGroups,omitempty"`
}

// Get retrieves scripts matching params.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/script/get
func (s *ScriptService) Get(ctx context.Context, params ScriptGetParams) ([]Script, error) {
	var out []Script
	err := s.client.Call(ctx, "script.get", params, &out)
	return out, err
}

// Create creates scripts and returns the new script ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/script/create
func (s *ScriptService) Create(ctx context.Context, items ...Script) ([]string, error) {
	var res struct {
		IDs []string `json:"scriptids"`
	}
	err := s.client.Call(ctx, "script.create", items, &res)
	return res.IDs, err
}

// Update updates scripts and returns the affected script ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/script/update
func (s *ScriptService) Update(ctx context.Context, items ...Script) ([]string, error) {
	var res struct {
		IDs []string `json:"scriptids"`
	}
	err := s.client.Call(ctx, "script.update", items, &res)
	return res.IDs, err
}

// Delete deletes scripts by id and returns the deleted script ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/script/delete
func (s *ScriptService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	var res struct {
		IDs []string `json:"scriptids"`
	}
	err := s.client.Call(ctx, "script.delete", ids, &res)
	return res.IDs, err
}

// ScriptExecuteParams are the parameters for script.execute.
type ScriptExecuteParams struct {
	ScriptID    string `json:"scriptid"`
	HostID      string `json:"hostid,omitempty"`
	EventID     string `json:"eventid,omitempty"`
	ManualInput string `json:"manualinput,omitempty"`
}

// ScriptExecuteResult is the result of script.execute.
type ScriptExecuteResult struct {
	Response string `json:"response,omitempty"`
	Value    string `json:"value,omitempty"`
	Debug    any    `json:"debug,omitempty"`
}

// Execute runs a script on a host or event and returns its result.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/script/execute
func (s *ScriptService) Execute(ctx context.Context, params ScriptExecuteParams) (*ScriptExecuteResult, error) {
	var out ScriptExecuteResult
	err := s.client.Call(ctx, "script.execute", params, &out)
	return &out, err
}

// GetScriptsByHosts returns the scripts available for the given host ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/script/getscriptsbyhosts
func (s *ScriptService) GetScriptsByHosts(ctx context.Context, hostIDs []string) (any, error) {
	var out any
	err := s.client.Call(ctx, "script.getscriptsbyhosts", map[string]any{"hostids": hostIDs}, &out)
	return out, err
}

// GetScriptsByEvents returns the scripts available for the given event ids.
// https://www.zabbix.com/documentation/current/en/manual/api/reference/script/getscriptsbyevents
func (s *ScriptService) GetScriptsByEvents(ctx context.Context, eventIDs []string) (any, error) {
	var out any
	err := s.client.Call(ctx, "script.getscriptsbyevents", map[string]any{"eventids": eventIDs}, &out)
	return out, err
}
