//go:build integration

package zabbix

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
)

// TestIntegrationScriptLookupContracts only looks up URL scripts. It never
// executes a script or opens the resulting URL. Both scopes share one event
// fixture so the configuration-cache wait happens once.
func TestIntegrationScriptLookupContracts(t *testing.T) {
	c, ctx := loginClient(t)
	hostID, eventID := integrationProblemEvent(t, c, ctx)
	input := "contract-input"
	cases := []struct {
		name      string
		scope     int
		id, macro string
		lookup    func(string, *string) (map[string][]Script, error)
	}{
		{"hosts", 2, hostID, "{HOST.ID}", func(scriptID string, value *string) (map[string][]Script, error) {
			return c.Script().GetScriptsByHosts(ctx, []ScriptHostParams{{HostID: hostID, ScriptID: scriptID, ManualInput: value}})
		}},
		{"events", 4, eventID, "{EVENT.ID}", func(scriptID string, value *string) (map[string][]Script, error) {
			return c.Script().GetScriptsByEvents(ctx, []ScriptEventParams{{EventID: eventID, ScriptID: scriptID, ManualInput: value}})
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			name := fmt.Sprintf("sdk-query-url-%s-%d", tc.name, time.Now().UnixNano())
			scriptID := createIntegrationURLScript(t, c, ctx, name, tc.scope, tc.macro)
			result, err := tc.lookup(scriptID, &input)
			if err != nil {
				t.Fatalf("lookup: %v", err)
			}
			scripts, ok := result[tc.id]
			if len(result) != 1 || !ok || len(scripts) != 1 {
				t.Fatalf("lookup result = %+v, want one script under %s", result, tc.id)
			}
			script := scripts[0]
			wantURL := "https://example.invalid/lookup/" + tc.id + "?input=" + input
			if script.ScriptID != scriptID || script.Name != name || script.Type != "6" || script.Scope != strconv.Itoa(tc.scope) {
				t.Fatalf("script identity/scope = %+v", script)
			}
			if script.URL != wantURL || script.Confirmation != "input="+input || script.ManualInputPrompt != "Value for "+tc.id {
				t.Fatalf("resolved script = %+v, want URL %s", script, wantURL)
			}
			if script.ManualInput != "1" || script.ManualInputValidatorType != "0" || script.ManualInputValidator != ".*" || script.NewWindow != "1" {
				t.Fatalf("script input settings = %+v", script)
			}
			empty := ""
			result, err = tc.lookup(scriptID, &empty)
			if err != nil {
				t.Fatalf("empty manualinput: %v", err)
			}
			if len(result[tc.id]) != 1 || result[tc.id][0].URL != "https://example.invalid/lookup/"+tc.id+"?input=" {
				t.Fatalf("explicit empty manualinput was not resolved: %+v", result)
			}
			result, err = tc.lookup("not-an-id", nil)
			var request *RequestError
			var api *APIError
			if result != nil || !errors.As(err, &request) || request.Kind != ErrorAPI || !errors.As(err, &api) {
				t.Fatalf("invalid lookup = %#v, %v; want nil result and API error", result, err)
			}
		})
	}
	t.Run("empty request", func(t *testing.T) {
		result, err := c.Script().GetScriptsByHosts(ctx, nil)
		if err != nil || result == nil || len(result) != 0 {
			t.Fatalf("empty host lookup = %#v, %v", result, err)
		}
		result, err = c.Script().GetScriptsByEvents(ctx, nil)
		if err != nil || result == nil || len(result) != 0 {
			t.Fatalf("empty event lookup = %#v, %v", result, err)
		}
	})
}

// The fixture deliberately uses Call with the dedicated write DTO. This keeps
// its creation independent from the query methods that this test verifies.
func createIntegrationURLScript(t *testing.T, c *Client, ctx context.Context, name string, scope int, macro string) string {
	t.Helper()
	url := "https://example.invalid/lookup/" + macro + "?input={MANUALINPUT}"
	confirmation := "input={MANUALINPUT}"
	prompt := "Value for " + macro
	validator := ".*"
	manualInput, validatorType, newWindow := 1, 0, 1
	var created struct {
		IDs []string `json:"scriptids"`
	}
	err := c.Call(ctx, "script.create", []ScriptCreateParams{{
		Name: name, Type: 6, Scope: scope, URL: &url, Confirmation: &confirmation,
		ManualInput: &manualInput, ManualInputPrompt: &prompt, ManualInputValidator: &validator,
		ManualInputValidatorType: &validatorType, NewWindow: &newWindow,
	}}, &created)
	if err != nil {
		t.Fatalf("create URL script: %v; cause %v", err, errors.Unwrap(err))
	}
	if len(created.IDs) != 1 || created.IDs[0] == "" {
		t.Fatalf("create URL script IDs = %v", created.IDs)
	}
	id := created.IDs[0]
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		var deleted struct {
			IDs []string `json:"scriptids"`
		}
		if err := c.Call(cleanup, "script.delete", []string{id}, &deleted); err != nil {
			t.Errorf("delete URL script: %v", err)
			return
		}
		if len(deleted.IDs) != 1 || deleted.IDs[0] != id {
			t.Errorf("delete URL script IDs = %v, want [%s]", deleted.IDs, id)
		}
	})
	return id
}
