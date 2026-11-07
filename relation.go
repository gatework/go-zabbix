package zabbix

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

// SelectedObject represents a selector that returns one object or an empty array.
// Present is false when the selector was omitted from the response. When Present
// is true, a nil Value denotes the server's empty-array result.
type SelectedObject[T any] struct {
	Value   *T
	Present bool
}

// UnmarshalJSON leaves the receiver unchanged if the result cannot be decoded.
func (o *SelectedObject[T]) UnmarshalJSON(data []byte) error {
	raw := bytes.TrimSpace(data)
	var value *T
	switch {
	case len(raw) > 0 && raw[0] == '{':
		value = new(T)
		if err := json.Unmarshal(raw, value); err != nil {
			return err
		}
	case isEmptyArray(raw):
	default:
		return errors.New("selected object must be an object or an empty array")
	}
	*o = SelectedObject[T]{Value: value, Present: true}
	return nil
}

func (o SelectedObject[T]) MarshalJSON() ([]byte, error) {
	if o.Value != nil {
		return json.Marshal(o.Value)
	}
	if o.Present {
		return []byte("[]"), nil
	}
	return []byte("null"), nil
}

func isEmptyArray(raw []byte) bool {
	if len(raw) == 0 || raw[0] != '[' {
		return false
	}
	var values []json.RawMessage
	return json.Unmarshal(raw, &values) == nil && len(values) == 0
}

// Relation holds a selected collection or its count. An omitted selection has
// both fields nil; a selected empty collection has a non-nil empty Items slice;
// a count of zero has a non-nil Count pointing to zero. At most one field is set.
// Only selectors documented to accept "count" use this type.
type Relation[T any] struct {
	Items []T
	Count *uint64
}

// UnmarshalJSON accepts an object array or a non-negative decimal count and
// leaves the receiver unchanged if the complete value cannot be decoded.
func (r *Relation[T]) UnmarshalJSON(data []byte) error {
	raw := bytes.TrimSpace(data)
	var next Relation[T]
	if len(raw) > 0 && raw[0] == '[' {
		items, err := decodeObjectList[T](raw)
		if err != nil {
			return err
		}
		next.Items = items
	} else {
		count, err := decodeCount(raw)
		if err != nil {
			return err
		}
		next.Count = &count
	}
	*r = next
	return nil
}

// MarshalJSON preserves the Zabbix array/count representation. A zero Relation
// marshals as null; response model fields use omitzero to omit unselected data.
func (r Relation[T]) MarshalJSON() ([]byte, error) {
	if r.Count != nil {
		if r.Items != nil {
			return nil, errors.New("relation contains both items and count")
		}
		return json.Marshal(strconv.FormatUint(*r.Count, 10))
	}
	return json.Marshal(r.Items)
}

func decodeCount(raw []byte) (uint64, error) {
	raw = bytes.TrimSpace(raw)
	if !json.Valid(raw) {
		return 0, errors.New("invalid count JSON")
	}
	var digits string
	if raw[0] == '"' {
		if err := json.Unmarshal(raw, &digits); err != nil {
			return 0, err
		}
	} else {
		digits = string(raw)
	}
	if digits == "" {
		return 0, errors.New("empty count")
	}
	for _, digit := range digits {
		if digit < '0' || digit > '9' {
			return 0, errors.New("count must be a non-negative decimal integer")
		}
	}
	return strconv.ParseUint(digits, 10, 64)
}

func decodeObjectList[T any](raw []byte) ([]T, error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || raw[0] != '[' {
		return nil, errors.New("expected an object array")
	}
	var values []json.RawMessage
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil, err
	}
	items := make([]T, len(values))
	for i, value := range values {
		value = bytes.TrimSpace(value)
		if len(value) == 0 || value[0] != '{' {
			return nil, fmt.Errorf("element %d must be an object", i)
		}
		if err := json.Unmarshal(value, &items[i]); err != nil {
			return nil, fmt.Errorf("element %d: %w", i, err)
		}
	}
	return items, nil
}

func decodeObjectMap[T any](raw []byte) (map[string]T, error) {
	raw = bytes.TrimSpace(raw)
	if len(raw) > 0 && raw[0] == '[' {
		var values []json.RawMessage
		if err := json.Unmarshal(raw, &values); err != nil {
			return nil, err
		}
		if len(values) != 0 {
			return nil, errors.New("expected an object map or an empty array")
		}
		return make(map[string]T), nil
	}
	if len(raw) == 0 || raw[0] != '{' {
		return nil, errors.New("expected an object map or an empty array")
	}
	var values map[string]json.RawMessage
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil, err
	}
	items := make(map[string]T, len(values))
	for key, value := range values {
		value = bytes.TrimSpace(value)
		if key == "" || len(value) == 0 || value[0] != '{' {
			return nil, errors.New("map entries must have a non-empty key and object value")
		}
		var item T
		if err := json.Unmarshal(value, &item); err != nil {
			return nil, fmt.Errorf("map entry: %w", err)
		}
		items[key] = item
	}
	return items, nil
}

func decodeObjectGroups[T any](raw []byte) (map[string][]T, error) {
	raw = bytes.TrimSpace(raw)
	if isEmptyArray(raw) {
		return make(map[string][]T), nil
	}
	if len(raw) == 0 || raw[0] != '{' {
		return nil, errors.New("expected an object of arrays or an empty array")
	}
	var values map[string]json.RawMessage
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil, err
	}
	groups := make(map[string][]T, len(values))
	for key, value := range values {
		if key == "" {
			return nil, errors.New("group key must not be empty")
		}
		items, err := decodeObjectList[T](value)
		if err != nil {
			return nil, fmt.Errorf("group: %w", err)
		}
		groups[key] = items
	}
	return groups, nil
}
