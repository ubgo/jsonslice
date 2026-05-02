package jsonslice

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestJsonSlice_JSONRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		in   JsonSlice
		want string
	}{
		{"nil", nil, "null"},
		{"empty", JsonSlice{}, "[]"},
		{"mixed", JsonSlice{"a", float64(1), true, nil}, `["a",1,true,null]`},
		{"nested", JsonSlice{[]any{"x", float64(2)}, map[string]any{"k": "v"}}, `[["x",2],{"k":"v"}]`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.in)
			if err != nil {
				t.Fatalf("Marshal: %v", err)
			}
			if string(data) != tc.want {
				t.Fatalf("Marshal: got %s want %s", data, tc.want)
			}

			var got JsonSlice
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatalf("Unmarshal: %v", err)
			}
			data2, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("re-Marshal: %v", err)
			}
			if string(data2) != tc.want {
				t.Fatalf("round trip: got %s want %s", data2, tc.want)
			}
		})
	}
}

func TestJsonSlice_UnmarshalJSON_Invalid(t *testing.T) {
	var j JsonSlice
	if err := j.UnmarshalJSON([]byte(`{"k":"v"}`)); err == nil {
		t.Fatal("expected error for object input")
	}
	if err := j.UnmarshalJSON([]byte(`not json`)); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestJsonSlice_SQLValue(t *testing.T) {
	var nilSlice JsonSlice
	v, err := nilSlice.Value()
	if err != nil {
		t.Fatalf("Value(nil): %v", err)
	}
	if v != nil {
		t.Fatalf("nil slice Value() = %v, want nil", v)
	}

	s := JsonSlice{"a", float64(1)}
	v, err = s.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	got, ok := v.([]byte)
	if !ok {
		t.Fatalf("Value type = %T, want []byte", v)
	}
	if string(got) != `["a",1]` {
		t.Fatalf("Value = %s", got)
	}
}

func TestJsonSlice_SQLScan(t *testing.T) {
	cases := []struct {
		name    string
		src     any
		wantNil bool
		want    string
		wantErr bool
	}{
		{"nil src", nil, true, "", false},
		{"bytes", []byte(`["a",1]`), false, `["a",1]`, false},
		{"string", `["b",2]`, false, `["b",2]`, false},
		{"empty bytes", []byte{}, true, "", false},
		{"int unsupported", 42, false, "", true},
		{"bad json", []byte(`{`), false, "", true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var j JsonSlice
			err := j.Scan(tc.src)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("Scan: %v", err)
			}
			if tc.wantNil {
				if j != nil {
					t.Fatalf("want nil, got %v", j)
				}
				return
			}
			data, _ := json.Marshal(j)
			if string(data) != tc.want {
				t.Fatalf("Scan got %s want %s", data, tc.want)
			}
		})
	}
}

func TestJsonSlice_MarshalGQL(t *testing.T) {
	var buf bytes.Buffer
	JsonSlice{"x", float64(1)}.MarshalGQL(&buf)
	if buf.String() != `["x",1]` {
		t.Fatalf("MarshalGQL = %s", buf.String())
	}
}

func TestJsonSlice_UnmarshalGQL(t *testing.T) {
	cases := []struct {
		name    string
		in      any
		want    string
		wantErr bool
	}{
		{"nil", nil, "null", false},
		{"slice", []any{"a", float64(1)}, `["a",1]`, false},
		{"raw json array", json.RawMessage(`["b",2]`), `["b",2]`, false},
		{"raw json object", json.RawMessage(`{"k":"v"}`), "", true},
		{"unsupported", make(chan int), "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var j JsonSlice
			err := j.UnmarshalGQL(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalGQL: %v", err)
			}
			data, _ := json.Marshal(j)
			if string(data) != tc.want {
				t.Fatalf("got %s want %s", data, tc.want)
			}
		})
	}
}
