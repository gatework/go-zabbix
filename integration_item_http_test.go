//go:build integration

package zabbix

import (
	"reflect"
	"testing"
)

func TestIntegrationItemHTTPRoundTrip(t *testing.T) {
	c, ctx := loginClient(t)
	hostID, _, _ := integrationHost(t, c, ctx)
	// The item remains disabled, so these tests need no HTTP poller or target.
	ids, err := c.Item().Create(ctx, ItemCreateParams{
		HostID: hostID, Name: "HTTP fixture", Key: "sdk.http.fixture", Type: 19, ValueType: 4,
		Status: integrationValue(1), Delay: integrationValue("1m"), URL: integrationValue("https://example.invalid/metrics"),
		QueryFields:   []ItemHTTPFieldInput{{Name: "tag", Value: "first"}, {Name: "tag", Value: ""}},
		Headers:       []ItemHTTPFieldInput{{Name: "X-Tag", Value: "first"}, {Name: "X-Tag", Value: ""}},
		RequestMethod: integrationValue(1), PostType: integrationValue(2), Posts: integrationValue(`{"value":"before"}`),
		Timeout: integrationValue("5s"), StatusCodes: integrationValue("200-299"),
		FollowRedirects: integrationValue(0), RetrieveMode: integrationValue(2), OutputFormat: integrationValue(1),
		HTTPProxy: integrationValue("http://proxy.invalid:3128"), AuthType: integrationValue(1),
		Username: integrationValue("fixture-user"), Password: integrationValue("fixture-password"),
		VerifyPeer: integrationValue(1), VerifyHost: integrationValue(1),
		SSLCertFile: integrationValue("fixture.crt"), SSLKeyFile: integrationValue("fixture.key"),
		SSLKeyPassword: integrationValue("fixture-key-password"), AllowTraps: integrationValue(1),
		TrapperHosts: integrationValue("192.0.2.1"),
		Preprocessing: []ItemPreprocessingInput{
			{Type: 12, Params: integrationValue("$.value"), ErrorHandler: integrationValue(2), ErrorHandlerParams: integrationValue("")},
			{Type: 25, Params: integrationValue("a\nb")},
		},
	})
	itemID := integrationID(t, ids, err)
	// integrationHost removes the host and its items if an assertion fails.
	params := ItemGetParams{
		ItemIDs: []string{itemID}, SelectPreprocessing: "extend",
		Output: []string{
			"itemid", "hostid", "name", "key_", "type", "value_type", "status", "delay", "url", "query_fields", "headers",
			"request_method", "post_type", "posts", "timeout", "status_codes", "follow_redirects", "retrieve_mode",
			"output_format", "http_proxy", "authtype", "username", "password", "verify_peer", "verify_host",
			"ssl_cert_file", "ssl_key_file", "ssl_key_password", "allow_traps", "trapper_hosts",
		},
	}
	want := Item{
		ItemID: itemID, HostID: hostID, Name: "HTTP fixture", Key: "sdk.http.fixture", Type: "19", ValueType: "4",
		Status: "1", Delay: "1m", URL: "https://example.invalid/metrics",
		QueryFields:   []ItemHTTPField{{Name: "tag", Value: "first"}, {Name: "tag", Value: ""}},
		Headers:       []ItemHTTPField{{Name: "X-Tag", Value: "first"}, {Name: "X-Tag", Value: ""}},
		RequestMethod: "1", PostType: "2", Posts: `{"value":"before"}`,
		Timeout: "5s", StatusCodes: "200-299", FollowRedirects: "0", RetrieveMode: "2", OutputFormat: "1",
		HTTPProxy: "http://proxy.invalid:3128", AuthType: "1", Username: "fixture-user", Password: "fixture-password",
		VerifyPeer: "1", VerifyHost: "1", SSLCertFile: "fixture.crt", SSLKeyFile: "fixture.key",
		SSLKeyPassword: "fixture-key-password", AllowTraps: "1", TrapperHosts: "192.0.2.1",
		Preprocessing: []ItemPreprocessing{
			{Type: "12", Params: "$.value", ErrorHandler: "2", ErrorHandlerParams: ""},
			{Type: "25", Params: "a\nb", ErrorHandler: "0", ErrorHandlerParams: ""},
		},
	}
	items, err := c.Item().Get(ctx, params)
	integrationCheck(t, err)
	if !reflect.DeepEqual(items, []Item{want}) {
		t.Fatalf("stored HTTP item = %+v, want %+v", items, want)
	}

	ids, err = c.Item().Update(ctx, ItemUpdateParams{ItemID: itemID, Posts: integrationValue(`{"value":"after"}`)})
	if got := integrationID(t, ids, err); got != itemID {
		t.Fatalf("updated item ID = %s, want %s", got, itemID)
	}
	want.Posts = `{"value":"after"}`
	items, err = c.Item().Get(ctx, params)
	integrationCheck(t, err)
	if !reflect.DeepEqual(items, []Item{want}) {
		t.Fatalf("partial update changed omitted settings: got %+v, want %+v", items, want)
	}

	ids, err = c.Item().Update(ctx, ItemUpdateParams{
		ItemID: itemID, Headers: []ItemHTTPFieldInput{}, QueryFields: []ItemHTTPFieldInput{}, Preprocessing: []ItemPreprocessingInput{},
	})
	if got := integrationID(t, ids, err); got != itemID {
		t.Fatalf("updated item ID = %s, want %s", got, itemID)
	}
	want.Headers = []ItemHTTPField{}
	want.QueryFields = []ItemHTTPField{}
	want.Preprocessing = []ItemPreprocessing{}
	items, err = c.Item().Get(ctx, params)
	integrationCheck(t, err)
	if !reflect.DeepEqual(items, []Item{want}) {
		t.Fatalf("cleared HTTP collections = %+v, want %+v", items, want)
	}

	ids, err = c.Item().Delete(ctx, itemID)
	if got := integrationID(t, ids, err); got != itemID {
		t.Fatalf("deleted item ID = %s, want %s", got, itemID)
	}
	items, err = c.Item().Get(ctx, params)
	integrationCheck(t, err)
	if items == nil || len(items) != 0 {
		t.Fatalf("deleted item still returned: %+v", items)
	}
}
