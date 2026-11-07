package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
)

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

// SLIResult contains the reporting periods, service IDs, and corresponding SLI
// values. SLI is indexed first by Periods and then by ServiceIDs. If either list
// is empty, SLI is empty.
type SLIResult struct {
	Periods    []SLIPeriod  `json:"periods"`
	ServiceIDs []string     `json:"serviceids"`
	SLI        [][]SLIValue `json:"sli"`
}

// SLIPeriod is a reporting interval with an inclusive start and exclusive end.
type SLIPeriod struct {
	PeriodFrom int64 `json:"period_from"`
	PeriodTo   int64 `json:"period_to"`
}

// SLIValue contains the measured availability and error budget for one service
// in one reporting period. SLI may be -1 when availability cannot be calculated;
// ErrorBudget may be negative when the service has exceeded its error budget.
type SLIValue struct {
	Uptime            int64                 `json:"uptime"`
	Downtime          int64                 `json:"downtime"`
	SLI               float64               `json:"sli"`
	ErrorBudget       int64                 `json:"error_budget"`
	ExcludedDowntimes []SLIExcludedDowntime `json:"excluded_downtimes"`
}

// SLIExcludedDowntime identifies a planned downtime excluded from a reported SLI.
type SLIExcludedDowntime struct {
	Name       string `json:"name"`
	PeriodFrom int64  `json:"period_from"`
	PeriodTo   int64  `json:"period_to"`
}

// GetSLI returns SLI (service level indicator) data for an SLA.
// https://www.zabbix.com/documentation/7.4/en/manual/api/reference/sla/getsli
func (s *SLAService) GetSLI(ctx context.Context, params SLAGetSLIParams) (*SLIResult, error) {
	// CSla::getSli derives serviceids from PHP array keys, which can encode IDs
	// as JSON integers. RawMessage avoids converting those IDs through float64.
	// Pointer scalars and non-nil arrays distinguish missing fields from zero.
	var wire struct {
		Periods []struct {
			PeriodFrom *int64 `json:"period_from"`
			PeriodTo   *int64 `json:"period_to"`
		} `json:"periods"`
		ServiceIDs []json.RawMessage `json:"serviceids"`
		SLI        [][]struct {
			Uptime            *int64   `json:"uptime"`
			Downtime          *int64   `json:"downtime"`
			SLI               *float64 `json:"sli"`
			ErrorBudget       *int64   `json:"error_budget"`
			ExcludedDowntimes []struct {
				Name       *string `json:"name"`
				PeriodFrom *int64  `json:"period_from"`
				PeriodTo   *int64  `json:"period_to"`
			} `json:"excluded_downtimes"`
		} `json:"sli"`
	}
	if err := s.client.Call(ctx, "sla.getsli", params, &wire); err != nil {
		return nil, err
	}
	if wire.Periods == nil || wire.ServiceIDs == nil || wire.SLI == nil {
		return nil, requestError("sla.getsli", ErrorProtocol, errors.New("SLI result requires periods, serviceids, and sli arrays"))
	}
	rows := len(wire.Periods)
	if len(wire.ServiceIDs) == 0 {
		// CSla::calculateSli returns [] when no services match, even if there
		// are reporting periods; it does not return one empty row per period.
		rows = 0
	}
	if len(wire.SLI) != rows || (len(wire.Periods) == 0 && len(wire.ServiceIDs) != 0) {
		return nil, requestError("sla.getsli", ErrorProtocol, errors.New("SLI dimensions do not match periods and services"))
	}
	out := &SLIResult{
		Periods:    make([]SLIPeriod, len(wire.Periods)),
		ServiceIDs: make([]string, len(wire.ServiceIDs)),
		SLI:        make([][]SLIValue, len(wire.SLI)),
	}
	for i, period := range wire.Periods {
		if period.PeriodFrom == nil || period.PeriodTo == nil {
			return nil, requestError("sla.getsli", ErrorProtocol, errors.New("SLI periods require both bounds"))
		}
		out.Periods[i] = SLIPeriod{PeriodFrom: *period.PeriodFrom, PeriodTo: *period.PeriodTo}
	}
	for i, raw := range wire.ServiceIDs {
		raw = bytes.TrimSpace(raw)
		id := string(raw)
		if len(raw) != 0 && raw[0] == '"' {
			if err := json.Unmarshal(raw, &id); err != nil {
				return nil, requestError("sla.getsli", ErrorDecode, err)
			}
		}
		if err := requireWriteID(id, "serviceid"); err != nil {
			return nil, requestError("sla.getsli", ErrorProtocol, err)
		}
		out.ServiceIDs[i] = id
	}
	for i, row := range wire.SLI {
		if len(row) != len(wire.ServiceIDs) {
			return nil, requestError("sla.getsli", ErrorProtocol, errors.New("SLI row does not match services"))
		}
		out.SLI[i] = make([]SLIValue, len(row))
		for j, cell := range row {
			if cell.Uptime == nil || cell.Downtime == nil || cell.SLI == nil || cell.ErrorBudget == nil || cell.ExcludedDowntimes == nil {
				return nil, requestError("sla.getsli", ErrorProtocol, errors.New("SLI value is missing required fields"))
			}
			value := SLIValue{
				Uptime:            *cell.Uptime,
				Downtime:          *cell.Downtime,
				SLI:               *cell.SLI,
				ErrorBudget:       *cell.ErrorBudget,
				ExcludedDowntimes: make([]SLIExcludedDowntime, len(cell.ExcludedDowntimes)),
			}
			for k, downtime := range cell.ExcludedDowntimes {
				if downtime.Name == nil || downtime.PeriodFrom == nil || downtime.PeriodTo == nil {
					return nil, requestError("sla.getsli", ErrorProtocol, errors.New("SLI excluded downtime is missing required fields"))
				}
				value.ExcludedDowntimes[k] = SLIExcludedDowntime{
					Name: *downtime.Name, PeriodFrom: *downtime.PeriodFrom, PeriodTo: *downtime.PeriodTo,
				}
			}
			out.SLI[i][j] = value
		}
	}
	return out, nil
}
