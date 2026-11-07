package zabbix

import (
	"fmt"
	"testing"
)

func TestItemPreprocessingRejectsIncompleteSteps(t *testing.T) {
	for _, tc := range []struct {
		name string
		step ItemPreprocessingInput
	}{
		{"unknown type", ItemPreprocessingInput{Type: 31}},
		{"missing type", ItemPreprocessingInput{}},
		{"JSONPath without params", ItemPreprocessingInput{Type: 12, ErrorHandler: writeTestPointer(0)}},
		{"JSONPath without error handler", ItemPreprocessingInput{Type: 12, Params: writeTestPointer("$.value")}},
		{"invalid handler", ItemPreprocessingInput{Type: 6, ErrorHandler: writeTestPointer(4)}},
		{"unsupported check with default handler", ItemPreprocessingInput{Type: 26, Params: writeTestPointer("-1"), ErrorHandler: writeTestPointer(0)}},
		{"set value without value", ItemPreprocessingInput{Type: 6, ErrorHandler: writeTestPointer(2)}},
		{"set error without message", ItemPreprocessingInput{Type: 6, ErrorHandler: writeTestPointer(3)}},
		{"set empty error", ItemPreprocessingInput{Type: 6, ErrorHandler: writeTestPointer(3), ErrorHandlerParams: writeTestPointer("")}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := responseServer(t, func(req rpcRequest) string {
				t.Error("incomplete preprocessing reached upstream")
				return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"itemids":["42"]}}`, req.ID)
			})
			_, err := c.Item().Update(t.Context(), ItemUpdateParams{ItemID: "42", Preprocessing: []ItemPreprocessingInput{tc.step}})
			assertQueryError(t, err, "item.update", ErrorInvalidParams)
		})
	}
}

func TestItemPreprocessingPreservesEmptyFallbackAndMacros(t *testing.T) {
	c := queryFixtureClient(t, "item.update",
		`[{"itemid":"42","preprocessing":[{"type":12,"params":"{$JSONPATH}","error_handler":2,"error_handler_params":""},{"type":21,"params":"return value;"},{"type":19}]}]`,
		`{"itemids":["42"]}`)
	_, err := c.Item().Update(t.Context(), ItemUpdateParams{
		ItemID: "42", Preprocessing: []ItemPreprocessingInput{
			{Type: 12, Params: writeTestPointer("{$JSONPATH}"), ErrorHandler: writeTestPointer(2), ErrorHandlerParams: writeTestPointer("")},
			{Type: 21, Params: writeTestPointer("return value;")},
			{Type: 19},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}
