package zabbix

// Tag is an entity tag attached to a host, item, trigger, problem or event.
// Query operators belong to TagFilter rather than the stored tag.
type Tag struct {
	Tag       string `json:"tag"`
	Value     string `json:"value,omitempty"`
	Automatic string `json:"automatic,omitempty"`
}

// TagFilter matches entity tags in the namespaces that support tag filtering.
// Operator uses Zabbix's numeric values: 0 contains, 1 equals, 2 does not
// contain, 3 does not equal, 4 exists, and 5 does not exist.
type TagFilter struct {
	Tag      string `json:"tag"`
	Value    string `json:"value,omitempty"`
	Operator int    `json:"operator,omitempty"`
}
