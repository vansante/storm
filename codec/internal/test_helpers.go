package internal

import (
	"encoding/gob"
	"reflect"
	"testing"

	"github.com/vansante/storm/v3/codec"
)

type testStruct struct {
	Name string
}

// RoundtripTester is a test helper to test a MarshalUnmarshaler
func RoundtripTester(t *testing.T, c codec.MarshalUnmarshaler, vals ...any) {
	var val, to any
	if len(vals) > 0 {
		if len(vals) != 2 {
			panic("Wrong number of vals, expected 2")
		}
		val = vals[0]
		to = vals[1]
	} else {
		val = &testStruct{Name: "test"}
		to = &testStruct{}
	}

	encoded, err := c.Marshal(val)
	if err != nil {
		t.Fatal("Encode error:", err)
	}
	err = c.Unmarshal(encoded, to)
	if err != nil {
		t.Fatal("Decode error:", err)
	}
	if !equalExportedFields(val, to) {
		t.Fatalf("Roundtrip codec mismatch, expected\n%#v\ngot\n%#v", val, to)
	}
}

// equalExportedFields compares two values like reflect.DeepEqual, but when
// traversing structs it only considers exported fields. This allows roundtrip
// comparisons of protobuf-generated messages, which carry unexported internal
// fields (state, sizeCache, unknownFields) that differ between instances.
func equalExportedFields(a, b any) bool {
	return equalExportedValues(reflect.ValueOf(a), reflect.ValueOf(b))
}

func equalExportedValues(a, b reflect.Value) bool {
	if !a.IsValid() || !b.IsValid() {
		return a.IsValid() == b.IsValid()
	}
	if a.Type() != b.Type() {
		return false
	}
	switch a.Kind() {
	case reflect.Ptr, reflect.Interface:
		if a.IsNil() || b.IsNil() {
			return a.IsNil() == b.IsNil()
		}
		return equalExportedValues(a.Elem(), b.Elem())
	case reflect.Struct:
		for i := 0; i < a.NumField(); i++ {
			if !a.Type().Field(i).IsExported() {
				continue
			}
			if !equalExportedValues(a.Field(i), b.Field(i)) {
				return false
			}
		}
		return true
	case reflect.Slice, reflect.Array:
		if a.Len() != b.Len() {
			return false
		}
		for i := 0; i < a.Len(); i++ {
			if !equalExportedValues(a.Index(i), b.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Map:
		if a.Len() != b.Len() {
			return false
		}
		iter := a.MapRange()
		for iter.Next() {
			bv := b.MapIndex(iter.Key())
			if !bv.IsValid() {
				return false
			}
			if !equalExportedValues(iter.Value(), bv) {
				return false
			}
		}
		return true
	default:
		return reflect.DeepEqual(a.Interface(), b.Interface())
	}
}

func init() {
	gob.Register(&testStruct{})
}
