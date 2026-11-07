package zabbix

import "context"

// MediaTypeService wraps the "mediatype" API namespace.
type MediaTypeService struct{ client *Client }

// MediaType returns the mediatype service.
func (c *Client) MediaType() *MediaTypeService { return &MediaTypeService{c} }

// MediaType object.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mediatype/object
type MediaType struct {
	MediaTypeID        string `json:"mediatypeid,omitempty"`
	Name               string `json:"name,omitempty"`
	Type               string `json:"type,omitempty"`
	Status             string `json:"status,omitempty"`
	SMTPServer         string `json:"smtp_server,omitempty"`
	SMTPPort           string `json:"smtp_port,omitempty"`
	SMTPHelo           string `json:"smtp_helo,omitempty"`
	SMTPEmail          string `json:"smtp_email,omitempty"`
	SMTPSecurity       string `json:"smtp_security,omitempty"`
	SMTPVerifyPeer     string `json:"smtp_verify_peer,omitempty"`
	SMTPVerifyHost     string `json:"smtp_verify_host,omitempty"`
	SMTPAuthentication string `json:"smtp_authentication,omitempty"`
	Username           string `json:"username,omitempty"`
	Passwd             string `json:"passwd,omitempty"`
	ExecPath           string `json:"exec_path,omitempty"`
	GSMModem           string `json:"gsm_modem,omitempty"`
	ContentType        string `json:"content_type,omitempty"`
	Script             string `json:"script,omitempty"`
	Timeout            string `json:"timeout,omitempty"`
	ProcessTags        string `json:"process_tags,omitempty"`
	ShowEventMenu      string `json:"show_event_menu,omitempty"`
	EventMenuURL       string `json:"event_menu_url,omitempty"`
	EventMenuName      string `json:"event_menu_name,omitempty"`
	Description        string `json:"description,omitempty"`
	MaxSessions        string `json:"maxsessions,omitempty"`
	MaxAttempts        string `json:"maxattempts,omitempty"`
	AttemptInterval    string `json:"attempt_interval,omitempty"`
	MessageFormat      string `json:"message_format,omitempty"`

	Parameters       []MediaTypeParameter       `json:"parameters,omitempty"`
	MessageTemplates []MediaTypeMessageTemplate `json:"message_templates,omitempty"`
	Users            []User                     `json:"users,omitempty"`
}

// MediaTypeParameter is a webhook/script parameter of a media type.
type MediaTypeParameter struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// MediaTypeMessageTemplate is a per-eventsource message template of a media type.
type MediaTypeMessageTemplate struct {
	EventSource string `json:"eventsource,omitempty"`
	Recovery    string `json:"recovery,omitempty"`
	Subject     string `json:"subject,omitempty"`
	Message     string `json:"message,omitempty"`
}

// MediaTypeGetParams are the parameters for mediatype.get.
// MediaTypeGetParams queries the list form of mediatype.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type MediaTypeGetParams struct {
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
	MediaTypeIDs           []string       `json:"mediatypeids,omitzero"`
	MediaIDs               []string       `json:"mediaids,omitzero"`
	UserIDs                []string       `json:"userids,omitzero"`
	SelectMessageTemplates any            `json:"selectMessageTemplates,omitempty"`
	SelectUsers            any            `json:"selectUsers,omitempty"`
}

// Get retrieves media types matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mediatype/get
func (s *MediaTypeService) Get(ctx context.Context, params MediaTypeGetParams) ([]MediaType, error) {
	return getList[MediaType](ctx, s.client, "mediatype.get", params)
}

// MediaTypeCountParams contains the filters for mediatype.get counts.
type MediaTypeCountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	MediaTypeIDs           []string       `json:"mediatypeids,omitzero"`
	MediaIDs               []string       `json:"mediaids,omitzero"`
	UserIDs                []string       `json:"userids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *MediaTypeService) Count(ctx context.Context, params MediaTypeCountParams) (uint64, error) {
	return getCount(ctx, s.client, "mediatype.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *MediaTypeService) GetMap(ctx context.Context, params MediaTypeGetParams) (map[string]MediaType, error) {
	return getMap[MediaType](ctx, s.client, "mediatype.get", params)
}

// Create creates media types and returns the new media type ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mediatype/create
func (s *MediaTypeService) Create(ctx context.Context, items ...MediaTypeCreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "mediatype.create", "mediatypeids", items)
}

// Update updates media types and returns the affected media type ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mediatype/update
func (s *MediaTypeService) Update(ctx context.Context, items ...MediaTypeUpdateParams) ([]string, error) {
	return updateIDs(ctx, s.client, "mediatype.update", "mediatypeids", items)
}

// Delete deletes media types by id and returns the deleted media type ids.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/mediatype/delete
func (s *MediaTypeService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "mediatype.delete", "mediatypeids", ids)
}
