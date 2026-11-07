package zabbix

import (
	"context"
	"fmt"
	"reflect"
	"testing"
)

func TestItemHTTPReadFields(t *testing.T) {
	const record = `{"itemid":"42","type":"19","url":"https://example.test/metrics","query_fields":[{"name":"tag","value":"first"},{"name":"tag","value":""}],"headers":[{"name":"X-Tag","value":"first"},{"name":"X-Tag","value":""}],"request_method":"1","post_type":"2","posts":"{}","timeout":"5s","status_codes":"200-299","follow_redirects":"0","retrieve_mode":"2","output_format":"1","http_proxy":"http://proxy.test:3128","authtype":"1","username":"fixture-user","password":"fixture-password","verify_peer":"1","verify_host":"1","ssl_cert_file":"fixture.crt","ssl_key_file":"fixture.key","ssl_key_password":"fixture-key-password","allow_traps":"1","trapper_hosts":"192.0.2.1","preprocessing":[{"type":"12","params":"$.value","error_handler":"2","error_handler_params":"fallback"}]}`
	c := queryFixtureClient(t, "item.get", `{"itemids":["42"],"output":"extend","selectPreprocessing":"extend"}`, "["+record+"]")
	items, err := c.Item().Get(context.Background(), ItemGetParams{
		ItemIDs: []string{"42"}, Output: "extend", SelectPreprocessing: "extend",
	})
	if err != nil || len(items) != 1 {
		t.Fatalf("items = %+v, error = %v", items, err)
	}
	assertQueryJSON(t, mustWriteJSON(t, items[0]), []byte(record))
}

func TestItemPreprocessingSelection(t *testing.T) {
	for _, tc := range []struct {
		name       string
		selection  any
		wantParams string
		result     string
		want       []ItemPreprocessing
	}{
		{
			name: "omitted", wantParams: `{}`, result: `[{"itemid":"42"}]`,
		},
		{
			name: "selected empty", selection: "extend",
			wantParams: `{"selectPreprocessing":"extend"}`,
			result:     `[{"itemid":"42","preprocessing":[]}]`, want: []ItemPreprocessing{},
		},
		{
			name: "selected fields", selection: []string{"type", "params"},
			wantParams: `{"selectPreprocessing":["type","params"]}`,
			result:     `[{"itemid":"42","preprocessing":[{"type":"12","params":"$.value"}]}]`,
			want:       []ItemPreprocessing{{Type: "12", Params: "$.value"}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := queryFixtureClient(t, "item.get", tc.wantParams, tc.result)
			items, err := c.Item().Get(context.Background(), ItemGetParams{SelectPreprocessing: tc.selection})
			if err != nil || len(items) != 1 {
				t.Fatalf("items = %+v, error = %v", items, err)
			}
			if !reflect.DeepEqual(items[0].Preprocessing, tc.want) {
				t.Fatalf("preprocessing = %#v, want %#v", items[0].Preprocessing, tc.want)
			}
		})
	}
}

func TestItemHTTPWriteWire(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		const want = `[{"hostid":"1","name":"HTTP fixture","key_":"http.fixture","type":19,"value_type":4,"delay":"{$INTERVAL}","url":"https://{HOST.CONN}/metrics","query_fields":[{"name":"tag","value":"first"},{"name":"tag","value":""}],"headers":[{"name":"X-Tag","value":"first"},{"name":"X-Tag","value":""}],"request_method":1,"post_type":2,"posts":"{\"value\":{$BODY}}","timeout":"{$TIMEOUT}","status_codes":"200,{$STATUS}","follow_redirects":0,"retrieve_mode":2,"output_format":1,"http_proxy":"http://proxy.test:3128","authtype":1,"username":"fixture-user","password":"fixture-password","verify_peer":1,"verify_host":1,"ssl_cert_file":"fixture.crt","ssl_key_file":"fixture.key","ssl_key_password":"fixture-key-password","allow_traps":1,"trapper_hosts":"192.0.2.1","preprocessing":[{"type":12,"params":"$.value","error_handler":2,"error_handler_params":""}]}]`
		c := queryFixtureClient(t, "item.create", want, `{"itemids":["42"]}`)
		ids, err := c.Item().Create(context.Background(), ItemCreateParams{
			HostID: "1", Name: "HTTP fixture", Key: "http.fixture", Type: 19, ValueType: 4,
			Delay: writeTestPointer("{$INTERVAL}"), URL: writeTestPointer("https://{HOST.CONN}/metrics"),
			QueryFields:   []ItemHTTPFieldInput{{Name: "tag", Value: "first"}, {Name: "tag", Value: ""}},
			Headers:       []ItemHTTPFieldInput{{Name: "X-Tag", Value: "first"}, {Name: "X-Tag", Value: ""}},
			RequestMethod: writeTestPointer(1), PostType: writeTestPointer(2), Posts: writeTestPointer(`{"value":{$BODY}}`),
			Timeout: writeTestPointer("{$TIMEOUT}"), StatusCodes: writeTestPointer("200,{$STATUS}"),
			FollowRedirects: writeTestPointer(0), RetrieveMode: writeTestPointer(2), OutputFormat: writeTestPointer(1),
			HTTPProxy: writeTestPointer("http://proxy.test:3128"), AuthType: writeTestPointer(1),
			Username: writeTestPointer("fixture-user"), Password: writeTestPointer("fixture-password"),
			VerifyPeer: writeTestPointer(1), VerifyHost: writeTestPointer(1),
			SSLCertFile: writeTestPointer("fixture.crt"), SSLKeyFile: writeTestPointer("fixture.key"),
			SSLKeyPassword: writeTestPointer("fixture-key-password"), AllowTraps: writeTestPointer(1),
			TrapperHosts: writeTestPointer("192.0.2.1"),
			Preprocessing: []ItemPreprocessingInput{{
				Type: 12, Params: writeTestPointer("$.value"), ErrorHandler: writeTestPointer(2), ErrorHandlerParams: writeTestPointer(""),
			}},
		})
		if err != nil || !reflect.DeepEqual(ids, []string{"42"}) {
			t.Fatalf("create ids = %v, error = %v", ids, err)
		}
	})
	t.Run("update clears collections and sends explicit zero", func(t *testing.T) {
		const want = `[{"itemid":"42","url":"https://example.test/metrics","query_fields":[],"headers":[],"request_method":0,"post_type":0,"posts":"","timeout":"","status_codes":"","follow_redirects":0,"retrieve_mode":0,"output_format":0,"http_proxy":"","authtype":0,"username":"","password":"","verify_peer":0,"verify_host":0,"ssl_cert_file":"","ssl_key_file":"","ssl_key_password":"","allow_traps":0,"trapper_hosts":"","preprocessing":[]}]`
		c := queryFixtureClient(t, "item.update", want, `{"itemids":["42"]}`)
		ids, err := c.Item().Update(context.Background(), ItemUpdateParams{
			ItemID: "42", URL: writeTestPointer("https://example.test/metrics"),
			QueryFields: []ItemHTTPFieldInput{}, Headers: []ItemHTTPFieldInput{},
			RequestMethod: writeTestPointer(0), PostType: writeTestPointer(0), Posts: writeTestPointer(""),
			Timeout: writeTestPointer(""), StatusCodes: writeTestPointer(""),
			FollowRedirects: writeTestPointer(0), RetrieveMode: writeTestPointer(0), OutputFormat: writeTestPointer(0),
			HTTPProxy: writeTestPointer(""), AuthType: writeTestPointer(0),
			Username: writeTestPointer(""), Password: writeTestPointer(""),
			VerifyPeer: writeTestPointer(0), VerifyHost: writeTestPointer(0),
			SSLCertFile: writeTestPointer(""), SSLKeyFile: writeTestPointer(""), SSLKeyPassword: writeTestPointer(""),
			AllowTraps: writeTestPointer(0), TrapperHosts: writeTestPointer(""), Preprocessing: []ItemPreprocessingInput{},
		})
		if err != nil || !reflect.DeepEqual(ids, []string{"42"}) {
			t.Fatalf("update ids = %v, error = %v", ids, err)
		}
	})
}

func TestItemHTTPPartialUpdates(t *testing.T) {
	for _, tc := range []struct {
		name   string
		params ItemUpdateParams
		want   string
	}{
		{"unrelated field omits HTTP configuration", ItemUpdateParams{ItemID: "42", Name: writeTestPointer("renamed")}, `[{"itemid":"42","name":"renamed"}]`},
		{"type may retain stored URL and interval", ItemUpdateParams{ItemID: "42", Type: writeTestPointer(19)}, `[{"itemid":"42","type":19}]`},
		{"post type may retain stored body", ItemUpdateParams{ItemID: "42", PostType: writeTestPointer(2)}, `[{"itemid":"42","post_type":2}]`},
		{"body may retain stored post type", ItemUpdateParams{ItemID: "42", Posts: writeTestPointer(`{"value":{$BODY}}`)}, `[{"itemid":"42","posts":"{\"value\":{$BODY}}"}]`},
		{"retrieve mode may retain stored request method", ItemUpdateParams{ItemID: "42", RetrieveMode: writeTestPointer(2)}, `[{"itemid":"42","retrieve_mode":2}]`},
		{"allowed hosts may retain stored allow traps", ItemUpdateParams{ItemID: "42", TrapperHosts: writeTestPointer("{$ALLOWED_HOSTS}")}, `[{"itemid":"42","trapper_hosts":"{$ALLOWED_HOSTS}"}]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := queryFixtureClient(t, "item.update", tc.want, `{"itemids":["42"]}`)
			ids, err := c.Item().Update(context.Background(), tc.params)
			if err != nil || !reflect.DeepEqual(ids, []string{"42"}) {
				t.Fatalf("update ids = %v, error = %v", ids, err)
			}
		})
	}
}

func TestItemHTTPCreateRequiresURL(t *testing.T) {
	c := responseServer(t, func(req rpcRequest) string {
		t.Error("HTTP item without a URL reached upstream")
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"itemids":["42"]}}`, req.ID)
	})
	_, err := c.Item().Create(context.Background(), ItemCreateParams{
		HostID: "1", Name: "HTTP fixture", Key: "http.fixture", Type: 19,
		ValueType: 4, Delay: writeTestPointer("1m"),
	})
	assertQueryError(t, err, "item.create", ErrorInvalidParams)
}

func TestItemHTTPCreateRejectsInvalidShape(t *testing.T) {
	for _, tc := range []struct {
		name   string
		change func(*ItemCreateParams)
	}{
		{"empty URL", func(p *ItemCreateParams) { p.URL = writeTestPointer("") }},
		{"missing interval", func(p *ItemCreateParams) { p.Delay = nil }},
		{"empty interval", func(p *ItemCreateParams) { p.Delay = writeTestPointer("") }},
		{"JSON without body", func(p *ItemCreateParams) { p.PostType = writeTestPointer(2) }},
		{"XML with empty body", func(p *ItemCreateParams) { p.PostType, p.Posts = writeTestPointer(3), writeTestPointer("") }},
		{"query field without name", func(p *ItemCreateParams) { p.QueryFields = []ItemHTTPFieldInput{{Value: "value"}} }},
		{"header without name", func(p *ItemCreateParams) { p.Headers = []ItemHTTPFieldInput{{Value: "value"}} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := responseServer(t, func(req rpcRequest) string {
				t.Error("invalid HTTP item reached upstream")
				return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"itemids":["42"]}}`, req.ID)
			})
			params := ItemCreateParams{
				HostID: "1", Name: "HTTP fixture", Key: "http.fixture", Type: 19, ValueType: 4,
				URL: writeTestPointer("https://example.test/metrics"), Delay: writeTestPointer("1m"),
			}
			tc.change(&params)
			_, err := c.Item().Create(context.Background(), params)
			assertQueryError(t, err, "item.create", ErrorInvalidParams)
		})
	}
}

func TestItemHTTPUpdateRejectsInvalidFieldNames(t *testing.T) {
	for _, tc := range []struct {
		name   string
		params ItemUpdateParams
	}{
		{"query field", ItemUpdateParams{ItemID: "42", QueryFields: []ItemHTTPFieldInput{{Value: "value"}}}},
		{"header", ItemUpdateParams{ItemID: "42", Headers: []ItemHTTPFieldInput{{Value: "value"}}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := responseServer(t, func(req rpcRequest) string {
				t.Error("invalid HTTP item reached upstream")
				return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"itemids":["42"]}}`, req.ID)
			})
			_, err := c.Item().Update(context.Background(), tc.params)
			assertQueryError(t, err, "item.update", ErrorInvalidParams)
		})
	}
}
