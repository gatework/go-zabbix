package zabbix

import "errors"

// ItemHTTPField is a header or query field returned for an HTTP agent item.
// Slice order and repeated names are significant; a value may be empty.
type ItemHTTPField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ItemHTTPFieldInput is a header or query field sent by an HTTP agent item.
// Slice order and repeated names are significant; a value may be empty.
type ItemHTTPFieldInput struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (p ItemHTTPFieldInput) validate() error {
	if p.Name == "" {
		return errors.New("zabbix: ItemHTTPFieldInput.Name is required")
	}
	return nil
}

// ItemPreprocessing is a preprocessing step returned for an ordinary item.
type ItemPreprocessing struct {
	Type               string `json:"type,omitempty"`
	Params             string `json:"params,omitempty"`
	ErrorHandler       string `json:"error_handler,omitempty"`
	ErrorHandlerParams string `json:"error_handler_params,omitempty"`
}
