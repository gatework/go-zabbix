package zabbix

import (
	"context"
	"strings"
	"testing"
)

func TestSLIQueryPreservesOptionalFilters(t *testing.T) {
	zero := int64(0)
	for _, tc := range []struct {
		name   string
		params SLAGetSLIParams
		want   string
	}{
		{"omitted", SLAGetSLIParams{SLAID: "1"}, `{"slaid":"1"}`},
		{"empty services and zero bounds", SLAGetSLIParams{SLAID: "1", ServiceIDs: []string{}, PeriodFrom: &zero, PeriodTo: &zero}, `{"slaid":"1","serviceids":[],"period_from":0,"period_to":0}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := queryFixtureClient(t, "sla.getsli", tc.want, `{"periods":[],"serviceids":[],"sli":[]}`)
			if _, err := c.SLA().GetSLI(context.Background(), tc.params); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestScriptExecutePreservesEmptyInputAndOutput(t *testing.T) {
	empty := ""
	c := queryFixtureClient(t, "script.execute", `{"scriptid":"1","hostid":"2","manualinput":""}`, `{"response":"success","value":"","debug":[]}`)
	got, err := c.Script().Execute(context.Background(), ScriptExecuteParams{ScriptID: "1", HostID: "2", ManualInput: &empty})
	if err != nil || got == nil || got.Response != "success" || got.Value != "" {
		t.Fatalf("script result = %#v, %v", got, err)
	}
}

func TestScriptExecuteRejectsIncompleteResults(t *testing.T) {
	for _, body := range []string{`{}`, `{"response":"success"}`, `{"response":"success","value":null}`, `{"response":"failed","value":""}`} {
		t.Run(body, func(t *testing.T) {
			c := queryFixtureClient(t, "script.execute", `{"scriptid":"1","hostid":"2"}`, body)
			got, err := c.Script().Execute(context.Background(), ScriptExecuteParams{ScriptID: "1", HostID: "2"})
			if got != nil {
				t.Fatalf("partial script result: %#v", got)
			}
			assertQueryError(t, err, "script.execute", ErrorProtocol)
		})
	}
}

func TestTokenGenerateRequiresCompleteRotationResult(t *testing.T) {
	t.Run("complete result preserves IDs and secrets", func(t *testing.T) {
		c := queryFixtureClient(t, "token.generate", `["0001","9007199254740993"]`, `[{"tokenid":"1","token":"first-secret"},{"tokenid":"9007199254740993","token":"second-secret"}]`)
		got, err := c.APIToken().Generate(context.Background(), "0001", "9007199254740993")
		if err != nil || len(got) != 2 || got[0].Token != "first-secret" || got[1].TokenID != "9007199254740993" {
			t.Fatalf("generated tokens = %#v, %v", got, err)
		}
	})
	for _, body := range []string{
		`[]`,
		`[{}]`,
		`[{"tokenid":"1"}]`,
		`[{"tokenid":"1","token":" "}]`,
		`[{"tokenid":"2","token":"secret-do-not-log"}]`,
	} {
		t.Run(body, func(t *testing.T) {
			c := queryFixtureClient(t, "token.generate", `["1"]`, body)
			got, err := c.APIToken().Generate(context.Background(), "1")
			if got != nil {
				t.Fatalf("partial generated tokens: %#v", got)
			}
			assertQueryError(t, err, "token.generate", ErrorProtocol)
			if strings.Contains(err.Error(), "secret-do-not-log") {
				t.Fatal("error includes generated token")
			}
		})
	}
	t.Run("duplicate response IDs", func(t *testing.T) {
		c := queryFixtureClient(t, "token.generate", `["1","2"]`, `[{"tokenid":"1","token":"a"},{"tokenid":"1","token":"b"}]`)
		got, err := c.APIToken().Generate(context.Background(), "1", "2")
		if got != nil {
			t.Fatalf("partial generated tokens: %#v", got)
		}
		assertQueryError(t, err, "token.generate", ErrorProtocol)
	})
	t.Run("null token entry", func(t *testing.T) {
		c := queryFixtureClient(t, "token.generate", `["1"]`, `[null]`)
		got, err := c.APIToken().Generate(context.Background(), "1")
		if got != nil {
			t.Fatalf("partial generated tokens: %#v", got)
		}
		assertQueryError(t, err, "token.generate", ErrorDecode)
	})
}

func TestTokenGenerateRejectsInvalidIDsBeforeRequest(t *testing.T) {
	for _, ids := range [][]string{nil, {""}, {"1", "01"}} {
		c, requests := mockServer(t, func(rpcRequest) (any, *APIError) { return []TokenGenerated{}, nil })
		_, err := c.APIToken().Generate(context.Background(), ids...)
		assertQueryError(t, err, "token.generate", ErrorInvalidParams)
		if len(*requests) != 0 {
			t.Fatal("invalid token IDs reached the transport")
		}
	}
}

func TestConfigurationExportSeparatesSerializedAndRawResults(t *testing.T) {
	options := map[string][]string{"hosts": {"1"}}
	t.Run("serialized format", func(t *testing.T) {
		c := queryFixtureClient(t, "configuration.export", `{"format":"json","options":{"hosts":["1"]}}`, `"{\"zabbix_export\":{\"version\":\"7.4\"}}"`)
		got, err := c.Configuration().Export(context.Background(), ConfigurationExportParams{Format: "json", Options: options})
		if err != nil || got != `{"zabbix_export":{"version":"7.4"}}` {
			t.Fatalf("serialized export = %q, %v", got, err)
		}
	})
	t.Run("raw object", func(t *testing.T) {
		body := `{"zabbix_export":{"version":"7.4","extension":9007199254740993}}`
		c := queryFixtureClient(t, "configuration.export", `{"format":"raw","options":{"hosts":["1"]}}`, body)
		got, err := c.Configuration().ExportRaw(context.Background(), options)
		if err != nil || string(got) != body {
			t.Fatalf("raw export = %s, %v", got, err)
		}
	})
	t.Run("raw format requires raw method", func(t *testing.T) {
		c, requests := mockServer(t, func(rpcRequest) (any, *APIError) { return map[string]any{}, nil })
		_, err := c.Configuration().Export(context.Background(), ConfigurationExportParams{Format: "raw", Options: options})
		assertQueryError(t, err, "configuration.export", ErrorInvalidParams)
		if len(*requests) != 0 {
			t.Fatal("raw export reached serialized transport")
		}
	})
	for _, body := range []string{`"serialized"`, `[]`, `true`} {
		t.Run("invalid raw "+body, func(t *testing.T) {
			c := queryFixtureClient(t, "configuration.export", `{"format":"raw","options":{"hosts":["1"]}}`, body)
			got, err := c.Configuration().ExportRaw(context.Background(), options)
			if got != nil {
				t.Fatalf("invalid raw export: %s", got)
			}
			assertQueryError(t, err, "configuration.export", ErrorProtocol)
		})
	}
}
