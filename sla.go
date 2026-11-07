package zabbix

import "context"

// SLAService wraps the "sla" API namespace.
type SLAService struct{ client *Client }

// SLA returns the sla service.
func (c *Client) SLA() *SLAService { return &SLAService{c} }

// SLA object (Zabbix 7.x).
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/sla/object
type SLA struct {
	SLAID             string                        `json:"slaid,omitempty"`
	Name              string                        `json:"name,omitempty"`
	Period            string                        `json:"period,omitempty"`
	SLO               string                        `json:"slo,omitempty"`
	EffectiveDate     string                        `json:"effective_date,omitempty"`
	Timezone          string                        `json:"timezone,omitempty"`
	Status            string                        `json:"status,omitempty"`
	Description       string                        `json:"description,omitempty"`
	ServiceTags       Relation[SLAServiceTag]       `json:"service_tags,omitzero"`
	Schedule          Relation[SLASchedule]         `json:"schedule,omitzero"`
	ExcludedDowntimes Relation[SLAExcludedDowntime] `json:"excluded_downtimes,omitzero"`
}

// SLAServiceTag associates services with an SLA by tag.
type SLAServiceTag struct {
	Tag      string `json:"tag,omitempty"`
	Operator string `json:"operator,omitempty"`
	Value    string `json:"value,omitempty"`
}

// SLASchedule defines a weekly schedule period for an SLA.
type SLASchedule struct {
	PeriodFrom string `json:"period_from,omitempty"`
	PeriodTo   string `json:"period_to,omitempty"`
}

// SLAExcludedDowntime is a planned downtime excluded from SLA calculation.
type SLAExcludedDowntime struct {
	Name       string `json:"name,omitempty"`
	PeriodFrom string `json:"period_from,omitempty"`
	PeriodTo   string `json:"period_to,omitempty"`
}

// SLAGetParams are the parameters for sla.get.
// SLAGetParams queries the list form of sla.get.
// Output and selectors accept "extend" or []string; only documented selectors accept "count".
type SLAGetParams struct {
	Output                  any            `json:"output,omitempty"`
	Filter                  map[string]any `json:"filter,omitzero"`
	Search                  map[string]any `json:"search,omitzero"`
	SearchByAny             bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled  bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch             bool           `json:"startSearch,omitempty"`
	ExcludeSearch           bool           `json:"excludeSearch,omitempty"`
	SortField               []string       `json:"sortfield,omitzero"`
	SortOrder               []string       `json:"sortorder,omitzero"`
	Limit                   int            `json:"limit,omitempty"`
	Editable                bool           `json:"editable,omitempty"`
	SLAIDs                  []string       `json:"slaids,omitzero"`
	ServiceIDs              []string       `json:"serviceids,omitzero"`
	SelectSchedule          any            `json:"selectSchedule,omitempty"`
	SelectExcludedDowntimes any            `json:"selectExcludedDowntimes,omitempty"`
	SelectServiceTags       any            `json:"selectServiceTags,omitempty"`
}

// Get retrieves SLAs matching params.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/sla/get
func (s *SLAService) Get(ctx context.Context, params SLAGetParams) ([]SLA, error) {
	return getList[SLA](ctx, s.client, "sla.get", params)
}

// SLACountParams contains the filters for sla.get counts.
type SLACountParams struct {
	Filter                 map[string]any `json:"filter,omitzero"`
	Search                 map[string]any `json:"search,omitzero"`
	SearchByAny            bool           `json:"searchByAny,omitempty"`
	SearchWildcardsEnabled bool           `json:"searchWildcardsEnabled,omitempty"`
	StartSearch            bool           `json:"startSearch,omitempty"`
	ExcludeSearch          bool           `json:"excludeSearch,omitempty"`
	Editable               bool           `json:"editable,omitempty"`
	SLAIDs                 []string       `json:"slaids,omitzero"`
	ServiceIDs             []string       `json:"serviceids,omitzero"`
}

// Count returns the number of objects matching params.
func (s *SLAService) Count(ctx context.Context, params SLACountParams) (uint64, error) {
	return getCount(ctx, s.client, "sla.get", params)
}

// GetMap retrieves objects keyed by their server IDs. Empty results are a non-nil map.
func (s *SLAService) GetMap(ctx context.Context, params SLAGetParams) (map[string]SLA, error) {
	return getMap[SLA](ctx, s.client, "sla.get", params)
}

// Create creates SLAs and returns the new SLA ids.
func (s *SLAService) Create(ctx context.Context, items ...SLACreateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "sla.create", "slaids", items)
}

// Update updates SLAs and returns the affected SLA ids.
func (s *SLAService) Update(ctx context.Context, items ...SLAUpdateParams) ([]string, error) {
	return mutateIDs(ctx, s.client, "sla.update", "slaids", items)
}

// Delete deletes SLAs by id and returns the deleted SLA ids.
func (s *SLAService) Delete(ctx context.Context, ids ...string) ([]string, error) {
	return deleteIDs(ctx, s.client, "sla.delete", "slaids", ids)
}

// SLAGetSLIParams are the parameters for sla.getsli.
// Nil ServiceIDs selects all matching services; an empty slice selects none.
// Nil period bounds use server defaults; pointers preserve Unix timestamp zero.
// Zabbix 7.4 CSla::getSli defaults these fields to null and forwards serviceids.
type SLAGetSLIParams struct {
	SLAID      string   `json:"slaid"`
	ServiceIDs []string `json:"serviceids,omitzero"`
	Periods    int      `json:"periods,omitempty"`
	PeriodFrom *int64   `json:"period_from,omitempty"`
	PeriodTo   *int64   `json:"period_to,omitempty"`
}

// GetSLI returns SLI (service level indicator) data for an SLA.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/sla/getsli
func (s *SLAService) GetSLI(ctx context.Context, params SLAGetSLIParams) (any, error) {
	var out any
	err := s.client.Call(ctx, "sla.getsli", params, &out)
	return out, err
}
