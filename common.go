package zabbix

// GetParams holds the query parameters common to every Zabbix `*.get` method.
// Namespace-specific parameter structs embed it and add their own filters.
//
// Output accepts either the string "extend"/"count" or a []string of field
// names. Filter and Search map field names to a value (or []string of values).
// Any select* option accepts "extend", "count", or a []string of fields.
//
// See https://www.zabbix.com/documentation/current/en/manual/api/reference_commentary
type GetParams struct {
	Output                 any            `json:"output,omitempty"`
	Filter                 map[string]any `json:"filter,omitempty"`
	Search                 map[string]any `json:"search,omitempty"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	SortField              []string       `json:"sortfield,omitempty"`
	SortOrder              []string       `json:"sortorder,omitempty"`
	Limit                  int            `json:"limit,omitempty"`
	Preservekeys           bool           `json:"preservekeys,omitempty"`
	CountOutput            bool           `json:"countOutput,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
}

// Tag is an entity tag (name/value pair) attached to hosts, items, triggers,
// problems and events.
type Tag struct {
	Tag      string `json:"tag"`
	Value    string `json:"value,omitempty"`
	Operator int    `json:"operator,omitempty"` // used in *.get tag filters
}
