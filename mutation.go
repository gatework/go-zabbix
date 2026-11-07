package zabbix

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type writeParams interface {
	validate() error
}

func mutateIDs[T writeParams](
	ctx context.Context, client *Client, method, key string, params []T,
) ([]string, error) {
	if len(params) == 0 {
		return nil, requestError(method, ErrorInvalidParams, errors.New("at least one object is required"))
	}
	for i, param := range params {
		if err := param.validate(); err != nil {
			return nil, requestError(method, ErrorInvalidParams, fmt.Errorf("object %d: %w", i, err))
		}
	}
	return mutationIDs(ctx, client, method, key, params, false)
}

func deleteIDs(ctx context.Context, client *Client, method, key string, ids []string) ([]string, error) {
	if err := checkIDs(ids); err != nil {
		return nil, requestError(method, ErrorInvalidParams, err)
	}
	return mutationIDs(ctx, client, method, key, ids, false)
}

func mutateObjectIDs[T writeParams](
	ctx context.Context, client *Client, method, key string, params T, allowEmpty bool,
) ([]string, error) {
	if err := params.validate(); err != nil {
		return nil, requestError(method, ErrorInvalidParams, err)
	}
	return mutationIDs(ctx, client, method, key, params, allowEmpty)
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
	if err := json.Unmarshal(raw, &ids); err != nil {
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

func mutateFields[T writeParams](
	ctx context.Context, client *Client, method string, params T,
) ([]string, error) {
	if err := params.validate(); err != nil {
		return nil, requestError(method, ErrorInvalidParams, err)
	}
	body, err := json.Marshal(params)
	if err != nil {
		return nil, requestError(method, ErrorInvalidParams, err)
	}
	if string(body) == "{}" {
		return nil, requestError(method, ErrorInvalidParams, errors.New("at least one setting is required"))
	}
	var fields []string
	if err := client.Call(ctx, method, params, &fields); err != nil {
		return nil, err
	}
	if len(fields) == 0 {
		return nil, requestError(method, ErrorProtocol, errors.New("missing updated field names"))
	}
	for _, field := range fields {
		if field == "" {
			return nil, requestError(method, ErrorProtocol, errors.New("empty updated field name"))
		}
	}
	return fields, nil
}
