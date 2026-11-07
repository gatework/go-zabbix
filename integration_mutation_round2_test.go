//go:build integration

package zabbix

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestIntegrationAutoregistrationUpdateReceipt(t *testing.T) {
	c, ctx := loginClient(t)
	original, err := c.Autoregistration().Get(ctx)
	integrationCheck(t, err)
	if original == nil {
		t.Fatal("autoregistration settings missing")
	}
	tlsAccept, err := strconv.Atoi(original.TLSAccept)
	integrationCheck(t, err)
	params := AutoregistrationUpdateParams{TLSAccept: &tlsAccept}
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := c.Autoregistration().Update(cleanup, params); err != nil {
			t.Errorf("restore autoregistration settings: %v", err)
			return
		}
		restored, err := c.Autoregistration().Get(cleanup)
		if err != nil {
			t.Errorf("read restored autoregistration settings: %v", err)
			return
		}
		if restored == nil || restored.TLSAccept != original.TLSAccept {
			t.Error("autoregistration TLS mode was not restored")
		}
	})

	// A same-value update returns true too. Changing this mode could require or
	// erase write-only PSKs, so exercise the receipt without modifying secrets.
	integrationCheck(t, c.Autoregistration().Update(ctx, params))
	settings, err := c.Autoregistration().Get(ctx)
	integrationCheck(t, err)
	if settings == nil || settings.TLSAccept != original.TLSAccept {
		t.Fatal("autoregistration TLS mode did not survive the update")
	}
}

func integrationMFAFixture(t *testing.T, c *Client, ctx context.Context, name string) string {
	t.Helper()
	ids, err := c.MFA().Create(ctx, MFACreateParams{
		Name: name, Type: integrationValue(1),
		HashFunction: integrationValue(1), CodeLength: integrationValue(6),
	})
	mfaID := integrationID(t, ids, err)
	t.Cleanup(func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		methods, err := c.MFA().Get(cleanup, MFAGetParams{
			MFAIDs: []string{mfaID}, Output: []string{"mfaid"},
		})
		if err != nil {
			t.Errorf("read MFA fixture for cleanup: %v", err)
			return
		}
		// A failed receipt decode can follow a committed deletion.
		if len(methods) == 0 {
			return
		}
		if _, err := c.MFA().Delete(cleanup, mfaID); err != nil {
			t.Errorf("MFA fixture cleanup: %v", err)
		}
	})
	return mfaID
}

func TestIntegrationMFADeleteNumericIDs(t *testing.T) {
	c, ctx := loginClient(t)
	name := fmt.Sprintf("go-zabbix-mfa-%d", time.Now().UnixNano())
	existing, err := c.MFA().Get(ctx, MFAGetParams{Output: []string{"mfaid"}, Limit: 1})
	integrationCheck(t, err)
	if len(existing) == 0 {
		// Zabbix makes the first method its default. Keep the regression target
		// unassigned, and delete this temporary default last during cleanup.
		integrationMFAFixture(t, c, ctx, name+"-keeper")
	}
	mfaID := integrationMFAFixture(t, c, ctx, name)
	params := MFAGetParams{
		MFAIDs: []string{mfaID}, Output: []string{"mfaid", "name", "type", "hash_function", "code_length"},
	}
	methods, err := c.MFA().Get(ctx, params)
	integrationCheck(t, err)
	if len(methods) != 1 || methods[0].MFAID != mfaID || methods[0].Name != name ||
		methods[0].Type != "1" || methods[0].HashFunction != "1" || methods[0].CodeLength != "6" {
		t.Fatalf("TOTP fixture not persisted: %+v", methods)
	}

	// Zabbix 7.4 returns numeric IDs for mfa.delete.
	ids, err := c.MFA().Delete(ctx, mfaID)
	if got := integrationID(t, ids, err); got != mfaID {
		t.Fatalf("deleted MFA id = %s, want %s", got, mfaID)
	}
	methods, err = c.MFA().Get(ctx, params)
	integrationCheck(t, err)
	if len(methods) != 0 {
		t.Fatalf("deleted MFA method still persisted: %+v", methods)
	}
}

func TestIntegrationHostInterfaceMassRemoveNumericIDs(t *testing.T) {
	c, ctx := loginClient(t)
	hostID, _, _ := integrationHost(t, c, ctx)
	ids, err := c.HostInterface().Create(ctx, HostInterfaceCreateParams{
		HostID: hostID, Type: 1, Main: 1, UseIP: 1, IP: "127.0.0.1", Port: "10050",
	})
	mainID := integrationID(t, ids, err)
	ids, err = c.HostInterface().Create(ctx, HostInterfaceCreateParams{
		HostID: hostID, Type: 1, Main: 0, UseIP: 1, IP: "127.0.0.2", Port: "10051",
	})
	removeID := integrationID(t, ids, err)
	interfaces, err := c.HostInterface().Get(ctx, HostInterfaceGetParams{
		InterfaceIDs: []string{removeID}, Output: "extend",
	})
	integrationCheck(t, err)
	if len(interfaces) != 1 || interfaces[0].InterfaceID != removeID || interfaces[0].HostID != hostID ||
		interfaces[0].Main != "0" || interfaces[0].IP != "127.0.0.2" || interfaces[0].DNS != "" || interfaces[0].Port != "10051" {
		t.Fatalf("removable interface not persisted: %+v", interfaces)
	}

	// A populated massremove result contains numeric IDs, unlike the empty
	// result exercised by TestIntegrationHostContractRoundTrip.
	ids, err = c.HostInterface().MassRemove(ctx, HostInterfaceMassRemoveParams{
		HostIDs:    []string{hostID},
		Interfaces: []HostInterfaceMatchInput{{IP: "127.0.0.2", DNS: "", Port: "10051"}},
	})
	if got := integrationID(t, ids, err); got != removeID {
		t.Fatalf("removed interface id = %s, want %s", got, removeID)
	}
	interfaces, err = c.HostInterface().Get(ctx, HostInterfaceGetParams{HostIDs: []string{hostID}, Output: "extend"})
	integrationCheck(t, err)
	if len(interfaces) != 1 || interfaces[0].InterfaceID != mainID || interfaces[0].HostID != hostID ||
		interfaces[0].Main != "1" || interfaces[0].IP != "127.0.0.1" || interfaces[0].Port != "10050" {
		t.Fatalf("massremove did not leave only the unmatched main interface: %+v", interfaces)
	}
}
