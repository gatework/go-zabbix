package zabbix

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func queryFixtureClient(t *testing.T, method, wantParams, result string) *Client {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Method string          `json:"method"`
			Params json.RawMessage `json:"params"`
			ID     json.RawMessage `json:"id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Error(err)
			http.Error(w, "bad JSON", 400)
			return
		}
		if req.Method != method {
			t.Errorf("method = %s, want %s", req.Method, method)
		}
		assertQueryJSON(t, req.Params, []byte(wantParams))
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": req.ID, "result": json.RawMessage(result)}); err != nil {
			t.Error(err)
		}
	}))
	t.Cleanup(srv.Close)
	c, err := New(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func assertQueryJSON(t *testing.T, got, want []byte) {
	t.Helper()
	var actual, expected any
	if err := json.Unmarshal(got, &actual); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(want, &expected); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("JSON = %s, want %s", got, want)
	}
}

func assertQueryError(t *testing.T, err error, method string, kind ErrorKind) {
	t.Helper()
	var request *RequestError
	if !errors.As(err, &request) || request.Method != method || request.Kind != kind {
		t.Fatalf("error = %#v, want %s %s", err, method, kind)
	}
}

func TestHostQueryPreservesEmptyFiltersAndDecodesSelections(t *testing.T) {
	c := queryFixtureClient(t, "host.get", `{"hostids":[],"output":["hostid"],"selectInterfaces":"count","selectItems":["itemid"],"selectInventory":"extend"}`,
		`[{"hostid":"101","interfaces":"0","items":[{"itemid":"12"}],"inventory":[]}]`)
	hosts, err := c.Host().Get(context.Background(), HostGetParams{
		HostIDs: []string{}, Output: []string{"hostid"}, SelectInterfaces: "count", SelectItems: []string{"itemid"}, SelectInventory: "extend",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hosts) != 1 || hosts[0].HostID != "101" || hosts[0].Interfaces.Count == nil || *hosts[0].Interfaces.Count != 0 || hosts[0].Interfaces.Items != nil {
		t.Fatalf("hosts = %#v", hosts)
	}
	if len(hosts[0].Items.Items) != 1 || hosts[0].Items.Items[0].ItemID != "12" || hosts[0].Inventory == nil || len(hosts[0].Inventory) != 0 {
		t.Fatalf("selected data = %#v", hosts[0])
	}
	if hosts[0].Triggers.Items != nil || hosts[0].Triggers.Count != nil {
		t.Fatal("unselected relation must remain absent")
	}
}

func TestCountAndMapContracts(t *testing.T) {
	ctx := context.Background()
	t.Run("count keeps exact uint64", func(t *testing.T) {
		c := queryFixtureClient(t, "host.get", `{"hostids":["101"],"countOutput":true}`, `"18446744073709551615"`)
		got, err := c.Host().Count(ctx, HostCountParams{HostIDs: []string{"101"}})
		if err != nil || got != math.MaxUint64 {
			t.Fatalf("count = %d, %v", got, err)
		}
	})
	t.Run("map retains keys", func(t *testing.T) {
		c := queryFixtureClient(t, "host.get", `{"output":["host"],"preservekeys":true}`, `{"101":{"host":"server"}}`)
		got, err := c.Host().GetMap(ctx, HostGetParams{Output: []string{"host"}})
		if err != nil || got["101"].Host != "server" {
			t.Fatalf("map = %#v, %v", got, err)
		}
	})
	for _, body := range []string{`[]`, `{}`} {
		t.Run("empty map "+body, func(t *testing.T) {
			c := queryFixtureClient(t, "host.get", `{"preservekeys":true}`, body)
			got, err := c.Host().GetMap(ctx, HostGetParams{})
			if err != nil || got == nil || len(got) != 0 {
				t.Fatalf("map = %#v, %v", got, err)
			}
		})
	}
	for _, body := range []string{`[{"hostid":"101"}]`, `{"101":null}`, `{"101":[]}`, `"2"`} {
		t.Run("bad map "+body, func(t *testing.T) {
			c := queryFixtureClient(t, "host.get", `{"preservekeys":true}`, body)
			got, err := c.Host().GetMap(ctx, HostGetParams{})
			if got != nil {
				t.Fatalf("partial result: %#v", got)
			}
			assertQueryError(t, err, "host.get", ErrorDecode)
		})
	}
	for _, body := range []string{`"-1"`, `-1`, `1.5`, `"1e3"`, `"18446744073709551616"`, `[]`, `{}`} {
		t.Run("bad count "+body, func(t *testing.T) {
			c := queryFixtureClient(t, "host.get", `{"countOutput":true}`, body)
			got, err := c.Host().Count(ctx, HostCountParams{})
			if got != 0 {
				t.Fatalf("partial count: %d", got)
			}
			assertQueryError(t, err, "host.get", ErrorDecode)
		})
	}
	t.Run("empty list", func(t *testing.T) {
		c := queryFixtureClient(t, "host.get", `{}`, `[]`)
		got, err := c.Host().Get(ctx, HostGetParams{})
		if err != nil || got == nil || len(got) != 0 {
			t.Fatalf("list = %#v, %v", got, err)
		}
	})
	for _, body := range []string{`{}`, `"2"`, `[null]`, `[{"hostid":"101"},{"interfaces":"bad"}]`} {
		t.Run("bad list "+body, func(t *testing.T) {
			c := queryFixtureClient(t, "host.get", `{}`, body)
			got, err := c.Host().Get(ctx, HostGetParams{})
			if got != nil {
				t.Fatalf("partial result: %#v", got)
			}
			assertQueryError(t, err, "host.get", ErrorDecode)
		})
	}
}

func TestSelectionValidationPrecedesRequest(t *testing.T) {
	tests := []struct {
		name   string
		params HostGetParams
	}{
		{"output count", HostGetParams{Output: "count"}},
		{"unavailable relation count", HostGetParams{SelectTags: "count"}},
		{"integer selector", HostGetParams{SelectInterfaces: 2}},
		{"object selector", HostGetParams{SelectItems: map[string]string{"output": "extend"}}},
		{"null selector", HostGetParams{SelectItems: []string(nil)}},
		{"empty field", HostGetParams{SelectItems: []string{""}}},
		{"negative limit", HostGetParams{Limit: -1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, requests := mockServer(t, func(rpcRequest) (any, *APIError) { return []Host{}, nil })
			_, err := c.Host().Get(context.Background(), tt.params)
			assertQueryError(t, err, "host.get", ErrorInvalidParams)
			if len(*requests) != 0 {
				t.Fatal("invalid selection reached the transport")
			}
		})
	}
}

func TestNamespaceSpecificQueryArguments(t *testing.T) {
	zero := int64(0)
	floating := 0
	t.Run("history zero timestamps and value type", func(t *testing.T) {
		c := queryFixtureClient(t, "history.get", `{"history":0,"hostids":[],"time_from":0,"time_till":0}`, `[{"itemid":"4","value":"0.25","ns":"7"}]`)
		got, err := c.History().Get(context.Background(), HistoryGetParams{History: &floating, HostIDs: []string{}, TimeFrom: &zero, TimeTill: &zero})
		if err != nil || len(got) != 1 || got[0].Value != "0.25" || got[0].NS != "7" {
			t.Fatalf("history = %#v, %v", got, err)
		}
	})
	t.Run("HA status is a filter", func(t *testing.T) {
		c := queryFixtureClient(t, "hanode.get", `{"filter":{"status":0}}`, `[{"ha_nodeid":"abc","status":"0"}]`)
		got, err := c.HANode().Get(context.Background(), HANodeGetParams{Filter: map[string]any{"status": 0}})
		if err != nil || len(got) != 1 || got[0].Status != "0" {
			t.Fatalf("nodes = %#v, %v", got, err)
		}
	})
	t.Run("global entity is selected by method", func(t *testing.T) {
		c := queryFixtureClient(t, "usermacro.get", `{"globalmacro":true,"globalmacroids":["7"]}`, `[{"globalmacroid":"7","macro":"{$GLOBAL}","value":"v"}]`)
		got, err := c.UserMacro().GetGlobal(context.Background(), GlobalMacroGetParams{GlobalMacroIDs: []string{"7"}})
		if err != nil || len(got) != 1 || got[0].GlobalMacroID != "7" {
			t.Fatalf("macros = %#v, %v", got, err)
		}
	})
	t.Run("global count", func(t *testing.T) {
		c := queryFixtureClient(t, "usermacro.get", `{"globalmacro":true,"countOutput":true}`, `"0"`)
		got, err := c.UserMacro().CountGlobal(context.Background(), GlobalMacroCountParams{})
		if err != nil || got != 0 {
			t.Fatalf("count = %d, %v", got, err)
		}
	})
	t.Run("global map", func(t *testing.T) {
		c := queryFixtureClient(t, "usermacro.get", `{"globalmacro":true,"preservekeys":true}`, `{"7":{"globalmacroid":"7"}}`)
		got, err := c.UserMacro().GetGlobalMap(context.Background(), GlobalMacroGetParams{})
		if err != nil || got["7"].GlobalMacroID != "7" {
			t.Fatalf("map = %#v, %v", got, err)
		}
	})
}

func TestScriptQueriesUseObjectArrays(t *testing.T) {
	empty := ""
	t.Run("hosts", func(t *testing.T) {
		c := queryFixtureClient(t, "script.getscriptsbyhosts", `[{"hostid":"101","scriptid":"9","manualinput":""}]`, `{"101":[{"scriptid":"9","name":"check","confirmation":"ready"}]}`)
		got, err := c.Script().GetScriptsByHosts(context.Background(), []ScriptHostParams{{HostID: "101", ScriptID: "9", ManualInput: &empty}})
		if err != nil || len(got["101"]) != 1 || got["101"][0].Confirmation != "ready" {
			t.Fatalf("scripts = %#v, %v", got, err)
		}
	})
	t.Run("events", func(t *testing.T) {
		c := queryFixtureClient(t, "script.getscriptsbyevents", `[{"eventid":"3"}]`, `{"3":[]}`)
		got, err := c.Script().GetScriptsByEvents(context.Background(), []ScriptEventParams{{EventID: "3"}})
		if err != nil || got["3"] == nil || len(got["3"]) != 0 {
			t.Fatalf("scripts = %#v, %v", got, err)
		}
	})
	t.Run("empty host request", func(t *testing.T) {
		c := queryFixtureClient(t, "script.getscriptsbyhosts", `[]`, `[]`)
		got, err := c.Script().GetScriptsByHosts(context.Background(), nil)
		if err != nil || got == nil || len(got) != 0 {
			t.Fatalf("scripts = %#v, %v", got, err)
		}
	})
	t.Run("malformed nested response", func(t *testing.T) {
		c := queryFixtureClient(t, "script.getscriptsbyevents", `[{"eventid":"3"}]`, `{"3":{"9":{"scriptid":"9"}}}`)
		got, err := c.Script().GetScriptsByEvents(context.Background(), []ScriptEventParams{{EventID: "3"}})
		if got != nil {
			t.Fatal("partial scripts returned")
		}
		assertQueryError(t, err, "script.getscriptsbyevents", ErrorDecode)
	})
}
