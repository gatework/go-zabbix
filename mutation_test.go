package zabbix

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func TestMutationRequiresAffectedIDs(t *testing.T) {
	for _, result := range []string{`{}`, `{"groupids":null}`, `{"groupids":[]}`, `{"groupids":[""]}`, `{"groupids":["0"]}`, `{"groupids":["1.5"]}`, `{"groupids":[23]}`} {
		t.Run(result, func(t *testing.T) {
			c := responseServer(t, func(req rpcRequest) string {
				return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":%s}`, req.ID, result)
			})
			ids, err := c.HostGroup().Create(context.Background(), HostGroupCreateParams{Name: "fixture"})
			var requestErr *RequestError
			if err == nil || ids != nil || !errors.As(err, &requestErr) || requestErr.Method != "hostgroup.create" {
				t.Fatalf("malformed mutation result: ids=%v error=%v", ids, err)
			}
		})
	}
}

func TestMutationRejectsInvalidInputBeforeSending(t *testing.T) {
	c := responseServer(t, func(req rpcRequest) string {
		t.Errorf("invalid input reached upstream: %s", req.Method)
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"groupids":["1"]}}`, req.ID)
	})
	for name, call := range map[string]func() error{
		"empty create": func() error { _, err := c.HostGroup().Create(context.Background()); return err },
		"missing name": func() error {
			_, err := c.HostGroup().Create(context.Background(), HostGroupCreateParams{})
			return err
		},
		"missing update id": func() error {
			_, err := c.HostGroup().Update(context.Background(), HostGroupUpdateParams{})
			return err
		},
		"empty delete":    func() error { _, err := c.HostGroup().Delete(context.Background()); return err },
		"blank delete id": func() error { _, err := c.HostGroup().Delete(context.Background(), " "); return err },
		"zero delete id":  func() error { _, err := c.HostGroup().Delete(context.Background(), "0"); return err },
		"negative delete id": func() error {
			_, err := c.HostGroup().Delete(context.Background(), "-1")
			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			var requestErr *RequestError
			if err := call(); !errors.As(err, &requestErr) || requestErr.Kind != ErrorInvalidParams {
				t.Fatalf("want invalid input error, got %v", err)
			}
		})
	}
}

func TestMutationPreservesLargeWireIDs(t *testing.T) {
	const id = "9007199254740993"
	c := responseServer(t, func(req rpcRequest) string {
		params, _ := req.Params.([]any)
		if len(params) != 1 || params[0] != id {
			t.Errorf("delete request IDs = %#v", req.Params)
		}
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"groupids":[%q]}}`, req.ID, id)
	})
	ids, err := c.HostGroup().Delete(context.Background(), id)
	if err != nil || len(ids) != 1 || ids[0] != id {
		t.Fatalf("large wire ID = %v, error = %v", ids, err)
	}
}

func TestSingletonUpdateRequiresFieldsAndPreservesExplicitZero(t *testing.T) {
	c := responseServer(t, func(req rpcRequest) string {
		params := writeJSON(t, req.Params)
		if req.Method != "authentication.update" || len(params) != 1 || string(params["http_auth_enabled"]) != "0" {
			t.Errorf("unexpected singleton request: %s %#v", req.Method, params)
		}
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":["http_auth_enabled"]}`, req.ID)
	})
	for name, update := range map[string]func() error{
		"authentication": func() error {
			_, err := c.Authentication().Update(context.Background(), AuthenticationUpdateParams{})
			return err
		},
		"autoregistration": func() error {
			_, err := c.Autoregistration().Update(context.Background(), AutoregistrationUpdateParams{})
			return err
		},
		"housekeeping": func() error {
			_, err := c.Housekeeping().Update(context.Background(), HousekeepingUpdateParams{})
			return err
		},
		"settings": func() error {
			_, err := c.Settings().Update(context.Background(), SettingsUpdateParams{})
			return err
		},
	} {
		t.Run(name, func(t *testing.T) {
			var requestErr *RequestError
			if err := update(); !errors.As(err, &requestErr) || requestErr.Kind != ErrorInvalidParams {
				t.Fatalf("empty singleton update = %v", err)
			}
		})
	}
	fields, err := c.Authentication().Update(context.Background(), AuthenticationUpdateParams{HTTPAuthEnabled: writeTestPointer(0)})
	if err != nil || len(fields) != 1 || fields[0] != "http_auth_enabled" {
		t.Fatalf("singleton updated fields = %v, error = %v", fields, err)
	}
}

func TestSingletonUpdateRejectsMalformedFields(t *testing.T) {
	for _, result := range []string{`null`, `[]`, `[""]`, `[7]`, `{}`} {
		t.Run(result, func(t *testing.T) {
			c := responseServer(t, func(req rpcRequest) string {
				return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":%s}`, req.ID, result)
			})
			fields, err := c.Authentication().Update(context.Background(), AuthenticationUpdateParams{HTTPAuthEnabled: writeTestPointer(0)})
			var requestErr *RequestError
			if err == nil || fields != nil || !errors.As(err, &requestErr) {
				t.Fatalf("malformed updated fields = %v, error = %v", fields, err)
			}
		})
	}
}

func TestEventAcknowledgePreservesZeroSeverityAndIndefiniteSuppression(t *testing.T) {
	c := responseServer(t, func(req rpcRequest) string {
		params := writeJSON(t, req.Params)
		if req.Method != "event.acknowledge" || string(params["severity"]) != "0" || string(params["suppress_until"]) != "0" {
			t.Errorf("event update lost explicit zero: %s %#v", req.Method, params)
		}
		if _, ok := params["message"]; ok {
			t.Error("omitted message was sent")
		}
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"eventids":["42"]}}`, req.ID)
	})
	ids, err := c.Event().Acknowledge(context.Background(), EventAcknowledgeParams{
		EventIDs: []string{"42"}, Action: AckActionChangeSeverity | AckActionSuppress,
		Severity: writeTestPointer(0), SuppressUntil: writeTestPointer(int64(0)),
	})
	if err != nil || len(ids) != 1 || ids[0] != "42" {
		t.Fatalf("acknowledged IDs = %v, error = %v", ids, err)
	}
}

func TestEventAcknowledgeRejectsIncompleteOrConflictingActions(t *testing.T) {
	c := responseServer(t, func(req rpcRequest) string {
		t.Errorf("invalid acknowledgement reached upstream: %#v", req.Params)
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"eventids":["42"]}}`, req.ID)
	})
	for name, params := range map[string]EventAcknowledgeParams{
		"missing severity":            {Action: AckActionChangeSeverity},
		"missing suppression":         {Action: AckActionSuppress},
		"missing cause":               {Action: AckActionChangeToSymptom},
		"missing message":             {Action: AckActionAddMessage},
		"empty message":               {Action: AckActionAddMessage, Message: writeTestPointer("")},
		"invalid severity":            {Action: AckActionChangeSeverity, Severity: writeTestPointer(6)},
		"invalid cause":               {Action: AckActionChangeToSymptom, CauseEventID: writeTestPointer("0")},
		"unexpected cause":            {Action: AckActionAcknowledge, CauseEventID: writeTestPointer("43")},
		"conflicting acknowledgement": {Action: AckActionAcknowledge | AckActionUnacknowledge},
		"conflicting suppression":     {Action: AckActionSuppress | AckActionUnsuppress},
		"closing while suppressing":   {Action: AckActionClose | AckActionSuppress},
		"conflicting rank":            {Action: AckActionChangeToCause | AckActionChangeToSymptom},
		"unsupported action":          {Action: 1 << 9},
	} {
		t.Run(name, func(t *testing.T) {
			params.EventIDs = []string{"42"}
			_, err := c.Event().Acknowledge(context.Background(), params)
			var requestErr *RequestError
			if !errors.As(err, &requestErr) || requestErr.Kind != ErrorInvalidParams || requestErr.Method != "event.acknowledge" {
				t.Fatalf("invalid event update = %v", err)
			}
		})
	}
}

func TestHostInterfaceEmptyResultsAreOnlyAllowedForRemovalAndReplacement(t *testing.T) {
	for _, operation := range []string{"massadd", "massremove", "replacehostinterfaces"} {
		for _, result := range []string{`{"interfaceids":[]}`, `{}`, `{"interfaceids":null}`} {
			t.Run(operation+"/"+result, func(t *testing.T) {
				c := responseServer(t, func(req rpcRequest) string {
					if req.Method != "hostinterface."+operation {
						t.Errorf("method = %s", req.Method)
					}
					return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":%s}`, req.ID, result)
				})
				var ids []string
				var err error
				switch operation {
				case "massadd":
					ids, err = c.HostInterface().MassAdd(context.Background(), HostInterfaceMassAddParams{
						Hosts: []HostRef{{HostID: "1"}}, Interfaces: []HostInterfaceInput{{Type: 1, Main: 1, UseIP: 1, IP: "192.0.2.1", Port: "10050"}},
					})
				case "massremove":
					ids, err = c.HostInterface().MassRemove(context.Background(), HostInterfaceMassRemoveParams{
						HostIDs: []string{"1"}, Interfaces: []HostInterfaceMatchInput{{IP: "192.0.2.1", Port: "10050"}},
					})
				case "replacehostinterfaces":
					ids, err = c.HostInterface().ReplaceHostInterfaces(context.Background(), HostInterfaceReplaceParams{HostID: "1", Interfaces: []HostInterfaceUpdateInput{}})
				}
				wantSuccess := operation != "massadd" && result == `{"interfaceids":[]}`
				if wantSuccess {
					if err != nil || ids == nil || len(ids) != 0 {
						t.Fatalf("empty affected IDs = %#v, error = %v", ids, err)
					}
				} else if err == nil || ids != nil {
					t.Fatalf("malformed affected IDs accepted: %#v, error = %v", ids, err)
				}
			})
		}
	}
}

func TestHostInterfaceMassAddDoesNotAcceptMalformedNestedResults(t *testing.T) {
	for _, result := range []string{`{"interfaceids":{}}`, `{"interfaceids":{"interfaceids":null}}`, `{"interfaceids":{"interfaceids":[]}}`, `{"interfaceids":{"interfaceids":[0]}}`} {
		t.Run(result, func(t *testing.T) {
			c := responseServer(t, func(req rpcRequest) string {
				return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":%s}`, req.ID, result)
			})
			ids, err := c.HostInterface().MassAdd(context.Background(), HostInterfaceMassAddParams{
				Hosts: []HostRef{{HostID: "1"}}, Interfaces: []HostInterfaceInput{{Type: 1, Main: 1, UseIP: 1, IP: "192.0.2.1", Port: "10050"}},
			})
			if err == nil || ids != nil {
				t.Fatalf("malformed nested IDs accepted: %#v, error = %v", ids, err)
			}
		})
	}
}

func TestHostUpdateSendsOnlyRequestedChanges(t *testing.T) {
	c := responseServer(t, func(req rpcRequest) string {
		var want any
		_ = json.Unmarshal([]byte(`[{"hostid":"42","description":"","status":0,"tags":[]}]`), &want)
		got, _ := json.Marshal(req.Params)
		expected, _ := json.Marshal(want)
		if req.Method != "host.update" || string(got) != string(expected) {
			t.Errorf("host update request = %s %s, want %s", req.Method, got, expected)
		}
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"hostids":["42"]}}`, req.ID)
	})
	ids, err := c.Host().Update(context.Background(), HostUpdateParams{HostID: "42", Description: writeTestPointer(""), Status: writeTestPointer(0), Tags: []TagInput{}})
	if err != nil || len(ids) != 1 || ids[0] != "42" {
		t.Fatalf("updated host IDs = %v, error = %v", ids, err)
	}
}

func TestItemPrototypeDeleteDecodesPrototypeIDs(t *testing.T) {
	c := responseServer(t, func(req rpcRequest) string {
		if req.Method != "itemprototype.delete" {
			t.Errorf("method = %s", req.Method)
		}
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"prototypeids":["12"]}}`, req.ID)
	})
	ids, err := c.ItemPrototype().Delete(context.Background(), "12")
	if err != nil || len(ids) != 1 || ids[0] != "12" {
		t.Fatalf("prototype IDs = %v, error = %v", ids, err)
	}
}

func TestDiscoveryRuleDeleteDecodesRuleIDs(t *testing.T) {
	c := responseServer(t, func(req rpcRequest) string {
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"ruleids":["12"]}}`, req.ID)
	})
	ids, err := c.DiscoveryRule().Delete(context.Background(), "12")
	if err != nil || len(ids) != 1 || ids[0] != "12" {
		t.Fatalf("discovery rule IDs = %v, error = %v", ids, err)
	}
}

func TestHostInterfaceMassAddAcceptsDocumentedAndServerResult(t *testing.T) {
	for _, result := range []string{`{"interfaceids":["12"]}`, `{"interfaceids":{"interfaceids":["12"]}}`} {
		t.Run(result, func(t *testing.T) {
			c := responseServer(t, func(req rpcRequest) string {
				return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":%s}`, req.ID, result)
			})
			ids, err := c.HostInterface().MassAdd(context.Background(), HostInterfaceMassAddParams{
				Hosts:      []HostRef{{HostID: "1"}},
				Interfaces: []HostInterfaceInput{{Type: 1, Main: 1, UseIP: 1, IP: "127.0.0.1", Port: "10050"}},
			})
			if err != nil || len(ids) != 1 || ids[0] != "12" {
				t.Fatalf("interface IDs = %v, error = %v", ids, err)
			}
		})
	}
}
