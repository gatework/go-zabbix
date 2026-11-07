package zabbix

import (
	"context"
	"reflect"
	"strings"
	"testing"
)

func TestSLIResultPreservesIDPrecisionAndValues(t *testing.T) {
	body := `{
		"periods":[{"period_from":0,"period_to":60}],
		"serviceids":[9007199254740993,"18446744073709551615","0007"],
		"sli":[[
			{"uptime":50,"downtime":10,"sli":83.33333333333334,"error_budget":-7,"excluded_downtimes":[{"name":"maintenance","period_from":5,"period_to":9}]},
			{"uptime":0,"downtime":0,"sli":-1,"error_budget":0,"excluded_downtimes":[]},
			{"uptime":60,"downtime":0,"sli":100,"error_budget":0,"excluded_downtimes":[]}
		]]
	}`
	c := queryFixtureClient(t, "sla.getsli", `{"slaid":"1"}`, body)
	got, err := c.SLA().GetSLI(context.Background(), SLAGetSLIParams{SLAID: "1"})
	if err != nil {
		t.Fatal(err)
	}
	want := &SLIResult{
		Periods:    []SLIPeriod{{PeriodFrom: 0, PeriodTo: 60}},
		ServiceIDs: []string{"9007199254740993", "18446744073709551615", "0007"},
		SLI: [][]SLIValue{{
			{Uptime: 50, Downtime: 10, SLI: 83.33333333333334, ErrorBudget: -7, ExcludedDowntimes: []SLIExcludedDowntime{{Name: "maintenance", PeriodFrom: 5, PeriodTo: 9}}},
			{Uptime: 0, Downtime: 0, SLI: -1, ErrorBudget: 0, ExcludedDowntimes: []SLIExcludedDowntime{}},
			{Uptime: 60, Downtime: 0, SLI: 100, ErrorBudget: 0, ExcludedDowntimes: []SLIExcludedDowntime{}},
		}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("SLI result = %#v, want %#v", got, want)
	}
}

func TestSLIResultAllowsEmptyDimensions(t *testing.T) {
	for _, body := range []string{
		`{"periods":[],"serviceids":[],"sli":[]}`,
		`{"periods":[{"period_from":0,"period_to":60}],"serviceids":[],"sli":[]}`,
	} {
		t.Run(body, func(t *testing.T) {
			c := queryFixtureClient(t, "sla.getsli", `{"slaid":"1"}`, body)
			got, err := c.SLA().GetSLI(context.Background(), SLAGetSLIParams{SLAID: "1"})
			if err != nil || got == nil {
				t.Fatalf("empty SLI = %#v, %v", got, err)
			}
			if got.Periods == nil || got.ServiceIDs == nil || got.SLI == nil || len(got.ServiceIDs) != 0 || len(got.SLI) != 0 {
				t.Fatalf("empty SLI collections were lost: %#v", got)
			}
		})
	}
}

func TestSLIResultRejectsMalformedResponses(t *testing.T) {
	const valid = `{"periods":[{"period_from":0,"period_to":60}],"serviceids":[1],"sli":[[{"uptime":0,"downtime":0,"sli":-1,"error_budget":0,"excluded_downtimes":[]}]]}`
	for _, tc := range []struct {
		name string
		body string
		kind ErrorKind
	}{
		{"non-object result", `[]`, ErrorDecode},
		{"missing fields", `{}`, ErrorProtocol},
		{"null periods", strings.Replace(valid, `[{"period_from":0,"period_to":60}]`, `null`, 1), ErrorProtocol},
		{"null service IDs", strings.Replace(valid, `"serviceids":[1]`, `"serviceids":null`, 1), ErrorProtocol},
		{"null period", strings.Replace(valid, `[{"period_from":0,"period_to":60}]`, `[null]`, 1), ErrorProtocol},
		{"missing period bound", strings.Replace(valid, `,"period_to":60`, ``, 1), ErrorProtocol},
		{"missing uptime", strings.Replace(valid, `"uptime":0,`, ``, 1), ErrorProtocol},
		{"null value", strings.Replace(valid, `"sli":-1`, `"sli":null`, 1), ErrorProtocol},
		{"string duration", strings.Replace(valid, `"uptime":0`, `"uptime":"0"`, 1), ErrorDecode},
		{"null excluded downtimes", strings.Replace(valid, `"excluded_downtimes":[]`, `"excluded_downtimes":null`, 1), ErrorProtocol},
		{"incomplete excluded downtime", strings.Replace(valid, `"excluded_downtimes":[]`, `"excluded_downtimes":[{"name":"maintenance","period_from":0}]`, 1), ErrorProtocol},
		{"missing SLI rows", `{"periods":[{"period_from":0,"period_to":60}],"serviceids":[1],"sli":[]}`, ErrorProtocol},
		{"missing SLI columns", `{"periods":[{"period_from":0,"period_to":60}],"serviceids":[1],"sli":[[]]}`, ErrorProtocol},
		{"null SLI row", `{"periods":[{"period_from":0,"period_to":60}],"serviceids":[1],"sli":[null]}`, ErrorProtocol},
		{"null SLI value", `{"periods":[{"period_from":0,"period_to":60}],"serviceids":[1],"sli":[[null]]}`, ErrorProtocol},
		{"rows without services", `{"periods":[{"period_from":0,"period_to":60}],"serviceids":[],"sli":[[]]}`, ErrorProtocol},
		{"services without periods", `{"periods":[],"serviceids":[1],"sli":[]}`, ErrorProtocol},
		{"null SLI matrix", `{"periods":[],"serviceids":[],"sli":null}`, ErrorProtocol},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := queryFixtureClient(t, "sla.getsli", `{"slaid":"1"}`, tc.body)
			got, err := c.SLA().GetSLI(context.Background(), SLAGetSLIParams{SLAID: "1"})
			if got != nil {
				t.Fatalf("partial SLI result: %#v", got)
			}
			assertQueryError(t, err, "sla.getsli", tc.kind)
		})
	}
	for _, id := range []string{`""`, `"0"`, `0`, `-1`, `1.5`, `1e3`, `"1e3"`, `null`, `true`, `{}`, `[]`} {
		t.Run("invalid ID "+id, func(t *testing.T) {
			body := strings.Replace(valid, `"serviceids":[1]`, `"serviceids":[`+id+`]`, 1)
			c := queryFixtureClient(t, "sla.getsli", `{"slaid":"1"}`, body)
			got, err := c.SLA().GetSLI(context.Background(), SLAGetSLIParams{SLAID: "1"})
			if got != nil {
				t.Fatalf("partial SLI result: %#v", got)
			}
			assertQueryError(t, err, "sla.getsli", ErrorProtocol)
		})
	}
}
