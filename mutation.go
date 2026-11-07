package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type writeParams interface {
	validate() error
}

type identifiedWriteParams interface {
	writeParams
	writeID() string
}

type objectWriteParams interface {
	writeParams
	validateResult([]string) error
}

func mutateIDs[T writeParams](
	ctx context.Context, client *Client, method, key string, params []T,
) ([]string, error) {
	if err := validateWriteList(method, params); err != nil {
		return nil, err
	}
	ids, err := mutationIDs(ctx, client, method, key, params, false)
	if err != nil {
		return nil, err
	}
	if len(ids) != len(params) {
		return nil, requestError(method, ErrorProtocol, errors.New("created id count does not match object count"))
	}
	// Task creation may reuse pending tasks, including several entries referring
	// to the same item or to dependent items sharing a master item.
	if method != "task.create" {
		if err := uniqueIDs(ids); err != nil {
			return nil, requestError(method, ErrorProtocol, err)
		}
	}
	return ids, nil
}

func updateIDs[T identifiedWriteParams](
	ctx context.Context, client *Client, method, key string, params []T,
) ([]string, error) {
	if err := validateWriteList(method, params); err != nil {
		return nil, err
	}
	expected := make([]string, len(params))
	for i, param := range params {
		expected[i] = param.writeID()
	}
	ids, err := mutationIDs(ctx, client, method, key, params, false)
	if err != nil {
		return nil, err
	}
	if err := correspondingIDs(ids, expected); err != nil {
		return nil, requestError(method, ErrorProtocol, err)
	}
	return ids, nil
}

func validateWriteList[T writeParams](method string, params []T) error {
	if len(params) == 0 {
		return requestError(method, ErrorInvalidParams, errors.New("at least one object is required"))
	}
	for i, param := range params {
		if err := param.validate(); err != nil {
			return requestError(method, ErrorInvalidParams, fmt.Errorf("object %d: %w", i, err))
		}
	}
	return nil
}

func deleteIDs(ctx context.Context, client *Client, method, key string, ids []string) ([]string, error) {
	if err := checkIDs(ids); err != nil {
		return nil, requestError(method, ErrorInvalidParams, err)
	}
	result, err := mutationIDs(ctx, client, method, key, ids, false)
	if err != nil {
		return nil, err
	}
	if err := correspondingIDs(result, ids); err != nil {
		return nil, requestError(method, ErrorProtocol, err)
	}
	return result, nil
}

func mutateObjectIDs[T objectWriteParams](
	ctx context.Context, client *Client, method, key string, params T, allowEmpty bool,
) ([]string, error) {
	if err := params.validate(); err != nil {
		return nil, requestError(method, ErrorInvalidParams, err)
	}
	ids, err := mutationIDs(ctx, client, method, key, params, allowEmpty)
	if err != nil {
		return nil, err
	}
	if err := params.validateResult(ids); err != nil {
		return nil, requestError(method, ErrorProtocol, err)
	}
	return ids, nil
}

// Compare decimal identities without changing their wire representation or
// passing through a floating-point or machine-sized integer representation.
func canonicalID(id string) string { return strings.TrimLeft(id, "0") }

func correspondingIDs(actual, expected []string) error {
	if len(actual) != len(expected) {
		return errors.New("affected id count does not match request")
	}
	remaining := make(map[string]int, len(expected))
	for _, id := range expected {
		remaining[canonicalID(id)]++
	}
	for _, id := range actual {
		key := canonicalID(id)
		if remaining[key] == 0 {
			return errors.New("affected ids do not match request")
		}
		remaining[key]--
	}
	return nil
}

func uniqueIDs(ids []string) error {
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		key := canonicalID(id)
		if _, exists := seen[key]; exists {
			return errors.New("result contains duplicate ids")
		}
		seen[key] = struct{}{}
	}
	return nil
}

func checkIDs(ids []string) error {
	if len(ids) == 0 {
		return errors.New("at least one id is required")
	}
	for _, id := range ids {
		if err := requireWriteID(id, "id"); err != nil {
			return err
		}
	}
	return nil
}

// mutationIDs validates the operation-specific result key. A successful RPC
// envelope alone does not prove that the server returned the affected IDs.
func mutationIDs(
	ctx context.Context, client *Client, method, key string, params any, allowEmpty bool,
) ([]string, error) {
	var result map[string]json.RawMessage
	if err := client.Call(ctx, method, params, &result); err != nil {
		return nil, err
	}
	raw, ok := result[key]
	if !ok || string(raw) == "null" {
		return nil, requestError(method, ErrorProtocol, fmt.Errorf("missing %s result", key))
	}
	// Zabbix 7.4 massadd wraps create's entire result in interfaceids once
	// more. Accept that observed shape only for this method, as well as the
	// documented flat array. Other malformed mutation results remain errors.
	if method == "hostinterface.massadd" && bytes.HasPrefix(bytes.TrimSpace(raw), []byte("{")) {
		var nested map[string]json.RawMessage
		if err := json.Unmarshal(raw, &nested); err != nil {
			return nil, requestError(method, ErrorDecode, err)
		}
		raw, ok = nested[key]
		if !ok || string(raw) == "null" {
			return nil, requestError(method, ErrorProtocol, fmt.Errorf("missing nested %s result", key))
		}
	}
	var ids []string
	var err error
	switch method {
	case "event.acknowledge", "hostinterface.massremove", "mfa.delete":
		// These operations return PHP array keys, which may be JSON integers.
		ids, err = decimalMutationIDs(raw)
	default:
		err = json.Unmarshal(raw, &ids)
	}
	if err != nil {
		return nil, requestError(method, ErrorDecode, err)
	}
	if len(ids) == 0 && allowEmpty {
		return ids, nil
	}
	if err := checkIDs(ids); err != nil {
		return nil, requestError(method, ErrorProtocol, err)
	}
	return ids, nil
}

func decimalMutationIDs(raw json.RawMessage) ([]string, error) {
	var values []json.RawMessage
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil, err
	}
	ids := make([]string, len(values))
	for i, value := range values {
		value = bytes.TrimSpace(value)
		if len(value) != 0 && value[0] == '"' {
			if err := json.Unmarshal(value, &ids[i]); err != nil {
				return nil, err
			}
		} else {
			// The caller's positive-decimal check rejects null, fractions,
			// exponent notation, and all other non-ID tokens.
			ids[i] = string(value)
		}
	}
	return ids, nil
}

func mutateFields[T writeParams](
	ctx context.Context, client *Client, method string, params T,
) ([]string, error) {
	body, expected, err := settingArguments(method, params)
	if err != nil {
		return nil, err
	}
	if method == "settings.update" {
		var provider string
		if raw, ok := expected["geomaps_tile_provider"]; ok {
			if err := json.Unmarshal(raw, &provider); err != nil {
				return nil, requestError(method, ErrorInvalidParams, err)
			}
		}
		if provider != "" {
			for _, field := range []string{"geomaps_tile_url", "geomaps_max_zoom", "geomaps_attribution"} {
				expected[field] = nil
			}
		}
	}
	var fields []string
	if err := client.Call(ctx, method, body, &fields); err != nil {
		return nil, err
	}
	if len(fields) != len(expected) {
		return nil, requestError(method, ErrorProtocol, errors.New("updated fields do not match request"))
	}
	for _, field := range fields {
		if _, ok := expected[field]; !ok {
			return nil, requestError(method, ErrorProtocol, errors.New("unexpected or duplicate updated field"))
		}
		delete(expected, field)
	}
	return fields, nil
}

func mutateSetting[T writeParams](ctx context.Context, client *Client, method string, params T) error {
	body, _, err := settingArguments(method, params)
	if err != nil {
		return err
	}
	var success bool
	if err := client.Call(ctx, method, body, &success); err != nil {
		return err
	}
	if !success {
		return requestError(method, ErrorProtocol, errors.New("update result is not true"))
	}
	return nil
}

func settingArguments[T writeParams](method string, params T) (json.RawMessage, map[string]json.RawMessage, error) {
	if err := params.validate(); err != nil {
		return nil, nil, requestError(method, ErrorInvalidParams, err)
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, nil, requestError(method, ErrorInvalidParams, err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil {
		return nil, nil, requestError(method, ErrorInvalidParams, err)
	}
	if len(fields) == 0 {
		return nil, nil, requestError(method, ErrorInvalidParams, errors.New("at least one setting is required"))
	}
	return body, fields, nil
}
