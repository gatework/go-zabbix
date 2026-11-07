package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

// queryArguments keeps protocol validation at the typed namespace boundary.
// Call deliberately remains available for untyped or newer server operations.
func queryArguments(method string, params any) (map[string]json.RawMessage, error) {
	encoded, err := json.Marshal(params)
	if err != nil {
		return nil, requestError(method, ErrorInvalidParams, err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(encoded, &fields); err != nil || fields == nil {
		return nil, requestError(method, ErrorInvalidParams, errors.New("query parameters must be an object"))
	}
	for field, raw := range fields {
		allowCount, selector := querySelectors[method][field]
		if field == "output" || selector {
			if err := validateSelection(raw, allowCount); err != nil {
				return nil, requestError(method, ErrorInvalidParams, fmt.Errorf("%s: %w", field, err))
			}
		}
		if field == "limit" {
			var limit int64
			if err := json.Unmarshal(raw, &limit); err != nil || limit <= 0 || limit > 2147483647 {
				return nil, requestError(method, ErrorInvalidParams, errors.New("limit must be between 1 and 2147483647"))
			}
		}
	}
	return fields, nil
}

func validateSelection(raw []byte, allowCount bool) error {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 {
		return errors.New("selection is empty")
	}
	if raw[0] == '"' {
		var selection string
		if err := json.Unmarshal(raw, &selection); err != nil {
			return err
		}
		if selection == "extend" || (allowCount && selection == "count") {
			return nil
		}
		return errors.New("selection must be extend, an array of fields, or a supported count")
	}
	if raw[0] != '[' {
		return errors.New("selection must be extend or an array of fields")
	}
	var fields []string
	if err := json.Unmarshal(raw, &fields); err != nil {
		return err
	}
	for _, field := range fields {
		if field == "" {
			return errors.New("selected field must not be empty")
		}
	}
	return nil
}

func getList[T any](ctx context.Context, client *Client, method string, params any) ([]T, error) {
	args, err := queryArguments(method, params)
	if err != nil {
		return nil, err
	}
	var raw json.RawMessage
	if err := client.Call(ctx, method, args, &raw); err != nil {
		return nil, err
	}
	items, err := decodeObjectList[T](raw)
	if err != nil {
		return nil, requestError(method, ErrorDecode, err)
	}
	return items, nil
}

func getMap[T any](ctx context.Context, client *Client, method string, params any) (map[string]T, error) {
	args, err := queryArguments(method, params)
	if err != nil {
		return nil, err
	}
	args["preservekeys"] = json.RawMessage("true")
	var raw json.RawMessage
	if err := client.Call(ctx, method, args, &raw); err != nil {
		return nil, err
	}
	items, err := decodeObjectMap[T](raw)
	if err != nil {
		return nil, requestError(method, ErrorDecode, err)
	}
	return items, nil
}

func getCount(ctx context.Context, client *Client, method string, params any) (uint64, error) {
	args, err := queryArguments(method, params)
	if err != nil {
		return 0, err
	}
	args["countOutput"] = json.RawMessage("true")
	var raw json.RawMessage
	if err := client.Call(ctx, method, args, &raw); err != nil {
		return 0, err
	}
	count, err := decodeCount(raw)
	if err != nil {
		return 0, requestError(method, ErrorDecode, err)
	}
	return count, nil
}

func getSingleton[T any](ctx context.Context, client *Client, method string) (*T, error) {
	var raw json.RawMessage
	if err := client.Call(ctx, method, map[string]string{"output": "extend"}, &raw); err != nil {
		return nil, err
	}
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || raw[0] != '{' {
		return nil, requestError(method, ErrorDecode, errors.New("expected an object"))
	}
	var result T
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, requestError(method, ErrorDecode, err)
	}
	return &result, nil
}

// querySelectors records the count capability of each exposed selector.
// Evidence and result fields are documented in docs/api-7.4-contracts.md.
var querySelectors = map[string]map[string]bool{
	"action.get": {
		"selectOperations":         false,
		"selectRecoveryOperations": false,
		"selectUpdateOperations":   false,
		"selectFilter":             false,
	},
	"alert.get": {
		"selectHosts":      false,
		"selectMediatypes": false,
		"selectUsers":      false,
	},
	"connector.get": {
		"selectTags": true,
	},
	"correlation.get": {
		"selectFilter":     false,
		"selectOperations": false,
	},
	"dashboard.get": {
		"selectPages":      false,
		"selectUsers":      false,
		"selectUserGroups": false,
	},
	"dhost.get": {
		"selectDServices": true,
	},
	"discoveryrule.get": {
		"selectHosts":         false,
		"selectFilter":        false,
		"selectLLDMacroPaths": false,
		"selectPreprocessing": false,
	},
	"drule.get": {
		"selectDChecks": true,
		"selectDHosts":  true,
	},
	"dservice.get": {
		"selectDRules": false,
		"selectDHosts": false,
		"selectHosts":  true,
	},
	"event.get": {
		"selectHosts":        false,
		"selectTags":         false,
		"selectAcknowledges": true,
	},
	"graph.get": {
		"selectGraphItems": false,
		"selectHosts":      false,
		"selectItems":      false,
		"selectTemplates":  false,
	},
	"graphitem.get": {
		"selectGraphs": false,
	},
	"graphprototype.get": {
		"selectGraphItems":    false,
		"selectHosts":         false,
		"selectItems":         false,
		"selectDiscoveryRule": false,
	},
	"host.get": {
		"selectHostGroups":      false,
		"selectInterfaces":      true,
		"selectItems":           true,
		"selectTriggers":        true,
		"selectParentTemplates": true,
		"selectTags":            false,
		"selectMacros":          false,
		"selectInventory":       false,
	},
	"hostgroup.get": {
		"selectHosts": true,
	},
	"hostinterface.get": {
		"selectItems": true,
		"selectHosts": false,
	},
	"hostprototype.get": {
		"selectGroupLinks":      false,
		"selectGroupPrototypes": false,
		"selectDiscoveryRule":   false,
		"selectParentHost":      false,
		"selectInterfaces":      false,
		"selectTemplates":       true,
		"selectTags":            false,
		"selectMacros":          false,
	},
	"iconmap.get": {
		"selectMappings": false,
	},
	"item.get": {
		"selectHosts":    false,
		"selectTriggers": true,
		"selectTags":     false,
	},
	"itemprototype.get": {
		"selectHosts":         false,
		"selectTags":          false,
		"selectPreprocessing": false,
		"selectDiscoveryRule": false,
	},
	"maintenance.get": {
		"selectTimeperiods": false,
		"selectTags":        false,
		"selectHostGroups":  false,
		"selectHosts":       false,
	},
	"map.get": {
		"selectSelements":  false,
		"selectLinks":      false,
		"selectUrls":       false,
		"selectUsers":      false,
		"selectUserGroups": false,
		"selectShapes":     false,
		"selectLines":      false,
	},
	"mediatype.get": {
		"selectMessageTemplates": false,
		"selectUsers":            false,
	},
	"problem.get": {
		"selectAcknowledges":    true,
		"selectTags":            false,
		"selectSuppressionData": false,
	},
	"proxy.get": {
		"selectHosts":      true,
		"selectProxyGroup": false,
	},
	"proxygroup.get": {
		"selectProxies": true,
	},
	"regexp.get": {
		"selectExpressions": false,
	},
	"report.get": {
		"selectUsers":      false,
		"selectUserGroups": false,
	},
	"role.get": {
		"selectRules": false,
	},
	"script.get": {
		"selectHosts":      false,
		"selectHostGroups": false,
	},
	"service.get": {
		"selectChildren":    true,
		"selectParents":     true,
		"selectTags":        true,
		"selectProblemTags": true,
		"selectStatusRules": true,
	},
	"sla.get": {
		"selectSchedule":          true,
		"selectExcludedDowntimes": true,
		"selectServiceTags":       true,
	},
	"template.get": {
		"selectTemplateGroups": false,
		"selectHosts":          true,
		"selectItems":          true,
		"selectTriggers":       true,
		"selectTags":           false,
		"selectMacros":         false,
	},
	"templatedashboard.get": {
		"selectPages": false,
	},
	"templategroup.get": {
		"selectTemplates": true,
	},
	"trigger.get": {
		"selectHosts":        false,
		"selectItems":        false,
		"selectTags":         false,
		"selectDependencies": false,
	},
	"triggerprototype.get": {
		"selectHosts":        false,
		"selectItems":        false,
		"selectTags":         false,
		"selectDependencies": false,
	},
	"user.get": {
		"selectMedias":  false,
		"selectUsrgrps": false,
		"selectRole":    false,
	},
	"usergroup.get": {
		"selectUsers":               false,
		"selectTagFilters":          false,
		"selectHostGroupRights":     false,
		"selectTemplateGroupRights": false,
	},
	"valuemap.get": {
		"selectMappings": true,
	},
}
