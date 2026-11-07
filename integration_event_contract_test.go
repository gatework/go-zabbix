//go:build integration

package zabbix

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

// integrationProblemEvent creates only disposable fixture data. The server's
// configuration cache may take a minute to see a new trapper item.
func integrationProblemEvent(t *testing.T, c *Client, ctx context.Context) (string, string) {
	t.Helper()
	hostID, hostName, _ := integrationHost(t, c, ctx)
	ids, err := c.Item().Create(ctx, ItemCreateParams{
		HostID: hostID, Name: "SDK trapper", Key: "sdk.problem", Type: 2, ValueType: 3,
	})
	itemID := integrationID(t, ids, err)
	ids, err = c.Trigger().Create(ctx, TriggerCreateParams{
		Description: "SDK problem", Expression: fmt.Sprintf("last(/%s/sdk.problem)>0", hostName),
		Priority: integrationValue(3), ManualClose: integrationValue(1),
	})
	triggerID := integrationID(t, ids, err)
	poll, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	var lastItemError string
	for {
		var pushed struct {
			Response string `json:"response"`
			Data     []struct {
				ItemID json.Number `json:"itemid"`
				Error  string      `json:"error"`
			} `json:"data"`
		}
		// history.push is outside the typed namespace surface. Call is the
		// supported escape hatch for this fixture-producing operation.
		err := c.Call(poll, "history.push", []map[string]any{{"itemid": itemID, "value": 1}}, &pushed)
		integrationCheck(t, err)
		if pushed.Response != "success" || len(pushed.Data) != 1 {
			t.Fatalf("unexpected fixture history response: %+v", pushed)
		}
		lastItemError = pushed.Data[0].Error
		if lastItemError == "" {
			events, err := c.Event().Get(poll, EventGetParams{
				ObjectIDs: []string{triggerID}, Source: integrationValue(0), Object: integrationValue(0),
				Value: []int{1}, Output: []string{"eventid"},
			})
			integrationCheck(t, err)
			if len(events) == 1 && events[0].EventID != "" {
				return hostID, events[0].EventID
			}
		}
		select {
		case <-poll.Done():
			t.Fatalf("fixture problem event was not created: %v (last item error %q)", poll.Err(), lastItemError)
		case <-ticker.C:
		}
	}
}

func TestIntegrationEventZeroValues(t *testing.T) {
	c, ctx := loginClient(t)
	_, eventID := integrationProblemEvent(t, c, ctx)
	ids, err := c.Event().Acknowledge(ctx, EventAcknowledgeParams{
		EventIDs: []string{eventID}, Action: AckActionChangeSeverity, Severity: integrationValue(0),
	})
	if got := integrationID(t, ids, err); got != eventID {
		t.Fatalf("acknowledged event = %s, want %s", got, eventID)
	}
	_, err = c.Event().Acknowledge(ctx, EventAcknowledgeParams{
		EventIDs: []string{eventID}, Action: AckActionSuppress, SuppressUntil: integrationValue(int64(0)),
	})
	integrationCheck(t, err)
	events, err := c.Event().Get(ctx, EventGetParams{
		EventIDs: []string{eventID}, Output: []string{"eventid", "severity"}, SelectAcknowledges: "extend",
	})
	integrationCheck(t, err)
	if len(events) != 1 || events[0].Severity != "0" {
		t.Fatalf("severity zero did not persist: %+v", events)
	}
	var changed, suppressed bool
	for _, update := range events[0].Acknowledges.Items {
		changed = changed || (update.Action == "8" && update.NewSeverity == "0")
		suppressed = suppressed || (update.Action == "32" && update.SuppressUntil == "0")
	}
	if !changed || !suppressed {
		t.Fatalf("explicit zero update records missing: %+v", events[0].Acknowledges.Items)
	}
}
