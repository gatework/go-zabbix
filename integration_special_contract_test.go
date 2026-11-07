//go:build integration

package zabbix

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestIntegrationRawConfigurationExport(t *testing.T) {
	c, ctx := loginClient(t)
	_, _, groupID := integrationHost(t, c, ctx)
	raw, err := c.Configuration().ExportRaw(ctx, map[string]any{"host_groups": []string{groupID}})
	integrationCheck(t, err)
	var exported struct {
		Configuration struct {
			Version string `json:"version"`
			Groups  []struct {
				Name string `json:"name"`
			} `json:"host_groups"`
		} `json:"zabbix_export"`
	}
	integrationCheck(t, json.Unmarshal(raw, &exported))
	if exported.Configuration.Version != "7.4" || len(exported.Configuration.Groups) != 1 || exported.Configuration.Groups[0].Name == "" {
		t.Fatalf("incomplete raw export: %s", raw)
	}
}

func TestIntegrationSLIEmptySelectionAndEpoch(t *testing.T) {
	c, ctx := loginClient(t)
	name := fmt.Sprintf("go-zabbix-sli-%d", time.Now().UnixNano())
	ids, err := c.Service().Create(ctx, ServiceCreateParams{
		Name: name, Algorithm: 0, SortOrder: 0,
		Tags: []TagInput{{Tag: "sdk-sli", Value: &name}},
	})
	serviceID := integrationID(t, ids, err)
	integrationDelete(t, c.Service().Delete, serviceID)
	ids, err = c.SLA().Create(ctx, SLACreateParams{
		Name: name, Period: 0, SLO: 99.9, EffectiveDate: integrationValue(0), Timezone: "UTC",
		ServiceTags: []SLATagInput{{Tag: "sdk-sli", Value: &name}},
	})
	slaID := integrationID(t, ids, err)
	integrationDelete(t, c.SLA().Delete, slaID)
	for _, tc := range []struct {
		name       string
		serviceIDs []string
		want       int
	}{
		{"omitted selects matching services", nil, 1},
		{"explicit empty selects no services", []string{}, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := c.SLA().GetSLI(ctx, SLAGetSLIParams{
				SLAID: slaID, ServiceIDs: tc.serviceIDs, Periods: 1, PeriodFrom: integrationValue(int64(0)),
			})
			integrationCheck(t, err)
			if len(result.ServiceIDs) != tc.want || (tc.want == 1 && result.ServiceIDs[0] != serviceID) {
				t.Fatalf("SLI service selection = %+v", result)
			}
			if len(result.Periods) != 1 || result.Periods[0].PeriodFrom != 0 {
				t.Fatalf("explicit epoch was omitted: %+v", result)
			}
		})
	}
}

func TestIntegrationTokenGeneration(t *testing.T) {
	c, ctx := loginClient(t)
	ids, err := c.APIToken().Create(ctx, TokenCreateParams{
		Name: fmt.Sprintf("go-zabbix-token-%d", time.Now().UnixNano()),
	})
	tokenID := integrationID(t, ids, err)
	integrationDelete(t, c.APIToken().Delete, tokenID)
	generated, err := c.APIToken().Generate(ctx, tokenID)
	integrationCheck(t, err)
	if len(generated) != 1 || generated[0].TokenID != tokenID || generated[0].Token == "" {
		t.Fatal("generated token identity or secret is missing")
	}
	apiClient, err := New(*zbxURL, WithAPIToken(generated[0].Token))
	integrationCheck(t, err)
	defer apiClient.CloseIdleConnections()
	_, err = apiClient.HostGroup().Count(ctx, HostGroupCountParams{})
	integrationCheck(t, err)
}
