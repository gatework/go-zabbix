package zabbix

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"testing"
)

func TestMutationResponseCorrespondsToRequest(t *testing.T) {
	for _, operation := range []string{"create", "update", "delete"} {
		for _, tc := range []struct {
			name, result string
			valid        bool
		}{
			{"missing", `["1"]`, false},
			{"extra", `["1","2","3"]`, false},
			{"duplicate", `["1","01"]`, false},
			{"unrelated", `["8","9"]`, operation == "create"},
			{"reordered canonical ids", `["2","001"]`, true},
		} {
			t.Run(operation+"/"+tc.name, func(t *testing.T) {
				c := responseServer(t, func(req rpcRequest) string {
					return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"groupids":%s}}`, req.ID, tc.result)
				})
				var ids []string
				var err error
				switch operation {
				case "create":
					ids, err = c.HostGroup().Create(t.Context(), HostGroupCreateParams{Name: "a"}, HostGroupCreateParams{Name: "b"})
				case "update":
					ids, err = c.HostGroup().Update(t.Context(), HostGroupUpdateParams{GroupID: "1"}, HostGroupUpdateParams{GroupID: "2"})
				case "delete":
					ids, err = c.HostGroup().Delete(t.Context(), "1", "2")
				}
				if tc.valid {
					if err != nil || len(ids) != 2 {
						t.Fatalf("valid response rejected: ids=%v err=%v", ids, err)
					}
					return
				}
				var requestErr *RequestError
				if ids != nil || !errors.As(err, &requestErr) || requestErr.Kind != ErrorProtocol {
					t.Fatalf("mismatched response accepted: ids=%v err=%v", ids, err)
				}
			})
		}
	}
}

func TestMutationPreservesLegitimateRepeatedIDs(t *testing.T) {
	for _, method := range []string{"task.create", "hostinterface.update", "hostinterface.delete", "event.acknowledge"} {
		t.Run(method, func(t *testing.T) {
			key := "interfaceids"
			if method == "task.create" {
				key = "taskids"
			}
			if method == "event.acknowledge" {
				key = "eventids"
			}
			c := responseServer(t, func(req rpcRequest) string {
				return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"%s":["001","1"]}}`, req.ID, key)
			})
			var ids []string
			var err error
			switch method {
			case "task.create":
				p := TaskCreateParams{Type: 6, Request: map[string]string{"itemid": "42"}}
				ids, err = c.Task().Create(t.Context(), p, p)
			case "hostinterface.update":
				p := HostInterfaceUpdateParams{InterfaceID: "1"}
				ids, err = c.HostInterface().Update(t.Context(), p, p)
			case "hostinterface.delete":
				ids, err = c.HostInterface().Delete(t.Context(), "1", "1")
			case "event.acknowledge":
				ids, err = c.Event().Acknowledge(t.Context(), EventAcknowledgeParams{EventIDs: []string{"001", "1", "1"}})
			}
			if err != nil || !slices.Equal(ids, []string{"001", "1"}) {
				t.Fatalf("legitimate repeated IDs = %v, error = %v", ids, err)
			}
		})
	}
}

func TestNumericMutationIDContracts(t *testing.T) {
	for _, method := range []string{"mfa.delete", "hostinterface.massremove"} {
		for _, result := range []string{`[9007199254740993]`, `["9007199254740993"]`, `[null]`, `[1.5]`, `[1e2]`, `[0]`} {
			t.Run(method+"/"+result, func(t *testing.T) {
				key := "mfaids"
				if method == "hostinterface.massremove" {
					key = "interfaceids"
				}
				c := responseServer(t, func(req rpcRequest) string {
					return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"%s":%s}}`, req.ID, key, result)
				})
				var ids []string
				var err error
				if method == "mfa.delete" {
					ids, err = c.MFA().Delete(t.Context(), "9007199254740993")
				} else {
					ids, err = c.HostInterface().MassRemove(t.Context(), HostInterfaceMassRemoveParams{
						HostIDs: []string{"1"}, Interfaces: []HostInterfaceMatchInput{{IP: "192.0.2.1", Port: "10050"}},
					})
				}
				valid := result == `[9007199254740993]` || result == `["9007199254740993"]`
				if valid && (err != nil || !slices.Equal(ids, []string{"9007199254740993"})) {
					t.Fatalf("valid numeric ID lost: ids=%v err=%v", ids, err)
				}
				if !valid && (err == nil || ids != nil) {
					t.Fatalf("invalid numeric ID accepted: ids=%v err=%v", ids, err)
				}
			})
		}
	}
}

func TestSingletonResponseCorrespondsToSubmittedFields(t *testing.T) {
	for _, result := range []string{`["http_auth_enabled","ldap_auth_enabled"]`, `["ldap_auth_enabled"]`, `["http_auth_enabled","http_auth_enabled"]`} {
		t.Run(result, func(t *testing.T) {
			c := responseServer(t, func(req rpcRequest) string {
				return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":%s}`, req.ID, result)
			})
			fields, err := c.Authentication().Update(t.Context(), AuthenticationUpdateParams{HTTPAuthEnabled: new(0)})
			if err == nil || fields != nil {
				t.Fatalf("incorrect updated fields accepted: %v, %v", fields, err)
			}
		})
	}
	for _, provider := range []string{"", "OpenStreetMap.Mapnik"} {
		t.Run("geomaps/"+provider, func(t *testing.T) {
			c := responseServer(t, func(req rpcRequest) string {
				fields := `["geomaps_tile_provider"]`
				if provider != "" {
					fields = `["geomaps_tile_provider","geomaps_tile_url","geomaps_max_zoom","geomaps_attribution"]`
				}
				return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":%s}`, req.ID, fields)
			})
			if _, err := c.Settings().Update(t.Context(), SettingsUpdateParams{GeomapsTileProvider: &provider}); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAutoregistrationUpdateAcceptsBooleanSuccess(t *testing.T) {
	c := responseServer(t, func(req rpcRequest) string {
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":true}`, req.ID)
	})
	err := c.Autoregistration().Update(context.Background(), AutoregistrationUpdateParams{TLSAccept: new(1)})
	if err != nil {
		t.Fatalf("valid boolean response rejected: %v", err)
	}
}

func TestAutoregistrationUpdateRejectsInvalidReceipt(t *testing.T) {
	for _, result := range []string{`false`, `null`, `[]`, `["tls_accept"]`, `{}`, `"true"`} {
		t.Run(result, func(t *testing.T) {
			c := responseServer(t, func(req rpcRequest) string {
				return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":%s}`, req.ID, result)
			})
			if err := c.Autoregistration().Update(t.Context(), AutoregistrationUpdateParams{TLSAccept: new(1)}); err == nil {
				t.Fatal("invalid registration receipt accepted")
			}
		})
	}
}

func TestSpecialMutationResponseCorrespondence(t *testing.T) {
	for _, tc := range []struct {
		name, key, result string
		call              func(*Client) ([]string, error)
		valid             bool
	}{
		{"massadd product", "interfaceids", `["10","11","12","13"]`, addFixtureInterfaces, true},
		{"massadd missing", "interfaceids", `["10","11","12"]`, addFixtureInterfaces, false},
		{"massadd duplicate", "interfaceids", `["10","11","12","010"]`, addFixtureInterfaces, false},
		{"massremove multiple matches", "interfaceids", `[10,11,12]`, removeFixtureInterfaces, true},
		{"massremove duplicate", "interfaceids", `[10,"010"]`, removeFixtureInterfaces, false},
		{"replace ordered identities", "interfaceids", `["041","50","41","51"]`, replaceFixtureInterfaces, true},
		{"replace wrong position", "interfaceids", `["50","41","41","51"]`, replaceFixtureInterfaces, false},
		{"replace reuses existing", "interfaceids", `["41","041","41","51"]`, replaceFixtureInterfaces, false},
		{"replace repeats new", "interfaceids", `["41","50","41","050"]`, replaceFixtureInterfaces, false},
		{"replace missing", "interfaceids", `["41","50","41"]`, replaceFixtureInterfaces, false},
		{"task missing", "taskids", `["42"]`, createFixtureTasks, false},
		{"event missing", "eventids", `[42]`, acknowledgeFixtureEvents, false},
		{"event extra", "eventids", `[42,43,44]`, acknowledgeFixtureEvents, false},
		{"event raw deduplication", "eventids", `[43,42]`, acknowledgeFixtureEvents, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := responseServer(t, func(req rpcRequest) string {
				return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"%s":%s}}`, req.ID, tc.key, tc.result)
			})
			ids, err := tc.call(c)
			if tc.valid {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			var requestErr *RequestError
			if ids != nil || !errors.As(err, &requestErr) || requestErr.Kind != ErrorProtocol {
				t.Fatalf("incorrect receipt accepted: ids=%v err=%v", ids, err)
			}
		})
	}
}

func addFixtureInterfaces(c *Client) ([]string, error) {
	item := HostInterfaceInput{Type: 1, Main: 0, UseIP: 1, IP: "192.0.2.1", Port: "10050"}
	return c.HostInterface().MassAdd(context.Background(), HostInterfaceMassAddParams{
		Hosts: []HostRef{{HostID: "1"}, {HostID: "2"}}, Interfaces: []HostInterfaceInput{item, item},
	})
}

func removeFixtureInterfaces(c *Client) ([]string, error) {
	return c.HostInterface().MassRemove(context.Background(), HostInterfaceMassRemoveParams{
		HostIDs: []string{"1"}, Interfaces: []HostInterfaceMatchInput{{IP: "192.0.2.1", Port: "10050"}},
	})
}

func replaceFixtureInterfaces(c *Client) ([]string, error) {
	return c.HostInterface().ReplaceHostInterfaces(context.Background(), HostInterfaceReplaceParams{
		HostID: "1", Interfaces: []HostInterfaceUpdateInput{
			{InterfaceID: new("41")},
			{Type: new(1), Main: new(0), UseIP: new(1), IP: new("192.0.2.1"), DNS: new(""), Port: new("10050")},
			{InterfaceID: new("41")},
			{Type: new(1), Main: new(0), UseIP: new(1), IP: new("192.0.2.2"), DNS: new(""), Port: new("10050")},
		},
	})
}

func createFixtureTasks(c *Client) ([]string, error) {
	p := TaskCreateParams{Type: 6, Request: map[string]string{"itemid": "42"}}
	return c.Task().Create(context.Background(), p, p)
}

func acknowledgeFixtureEvents(c *Client) ([]string, error) {
	return c.Event().Acknowledge(context.Background(), EventAcknowledgeParams{EventIDs: []string{"42", "43", "42"}})
}
