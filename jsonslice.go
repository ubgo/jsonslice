// Package jsonslice is a JSON-backed slice type that round-trips through
// SQL columns and GraphQL scalars without dragging gqlgen or any database
// driver into the dependency graph.
//
// The core type is [JsonSlice], an alias for []any that:
//
//   - Marshals/unmarshals as a JSON array via the standard encoding/json
//     interfaces.
//   - Implements sql.Scanner and driver.Valuer so it can be stored in a
//     JSON / JSONB column directly.
//   - Implements MarshalGQL/UnmarshalGQL so gqlgen will pick it up as a
//     custom scalar via duck typing — no gqlgen import is required, and
//     consumers that don't use gqlgen pay nothing.
//
// The zero value is an empty slice. nil is treated as the JSON value null
// when scanning; on marshal a nil [JsonSlice] is rendered as null too.
package jsonslice

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// JsonSlice is a slice of arbitrary JSON values.
type JsonSlice []any

// MarshalJSON implements [json.Marshaler]. A nil slice marshals to the JSON
// value null; an empty non-nil slice marshals to [].
func (j JsonSlice) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return json.Marshal([]any(j))
}

// UnmarshalJSON implements [json.Unmarshaler]. The JSON value null
// deserialises to a nil slice.
func (j *JsonSlice) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*j = nil
		return nil
	}
	var raw []any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("jsonslice: unmarshal: %w", err)
	}
	*j = raw
	return nil
}

// Value implements [driver.Valuer]. The zero / nil slice persists as SQL
// NULL; otherwise the slice is encoded as a JSON array.
func (j JsonSlice) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal([]any(j))
}

// Scan implements [sql.Scanner]. Accepts []byte, string, or nil. NULL
// scans into a nil slice.
func (j *JsonSlice) Scan(src any) error {
	if src == nil {
		*j = nil
		return nil
	}
	var data []byte
	switch v := src.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("jsonslice: cannot scan %T", src)
	}
	if len(data) == 0 {
		*j = nil
		return nil
	}
	return j.UnmarshalJSON(data)
}

// MarshalGQL is the gqlgen-compatible marshal hook. gqlgen detects this
// method by name; we therefore need not import gqlgen. Errors are rare
// (the slice already round-tripped through JSON to reach a resolver) and
// would indicate a programmer error encoding an unsupported type into the
// slice — those are surfaced by writing an empty array, never panicking
// inside a request goroutine.
func (j JsonSlice) MarshalGQL(w io.Writer) {
	data, err := json.Marshal([]any(j))
	if err != nil {
		_, _ = w.Write([]byte("[]"))
		return
	}
	_, _ = w.Write(data)
}

// UnmarshalGQL is the gqlgen-compatible unmarshal hook. It accepts the
// already-parsed JSON value gqlgen passes (a []any for a JSON array, or a
// scalar that the receiver chooses to wrap), re-encodes it, and decodes
// into the slice.
func (j *JsonSlice) UnmarshalGQL(v any) error {
	if v == nil {
		*j = nil
		return nil
	}
	if arr, ok := v.([]any); ok {
		*j = arr
		return nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("jsonslice: %T is not a JSON array", v)
	}
	var raw []any
	if err := json.Unmarshal(data, &raw); err != nil {
		return errors.New("jsonslice: value is not a JSON array")
	}
	*j = raw
	return nil
}
