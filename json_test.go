package pgtype_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgtype/testutil"
)

func TestJSONTranscode(t *testing.T) {
	testutil.TestSuccessfulTranscode(t, "json", []interface{}{
		pgtype.JSON("{}"),
		pgtype.JSON("null"),
		pgtype.JSON("42"),
		pgtype.JSON(`"hello"`),
		pgtype.JSON(nil),
	})
}

func TestJSONSet(t *testing.T) {
	successfulTests := []struct {
		source interface{}
		result pgtype.JSON
	}{
		{source: "{}", result: pgtype.JSON("{}")},
		{source: []byte("{}"), result: pgtype.JSON("{}")},
		{source: ([]byte)(nil), result: nil},
		{source: (*string)(nil), result: nil},
		{source: []int{1, 2, 3}, result: pgtype.JSON("[1,2,3]")},
		{source: map[string]interface{}{"foo": "bar"}, result: pgtype.JSON(`{"foo":"bar"}`)},
	}

	for i, tt := range successfulTests {
		var d pgtype.JSON
		err := d.Set(tt.source)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !reflect.DeepEqual(d, tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, d)
		}
	}
}

func TestJSONAssignTo(t *testing.T) {
	var s string
	var ps *string
	var b []byte

	rawStringTests := []struct {
		src      pgtype.JSON
		dst      *string
		expected string
	}{
		{src: pgtype.JSON("{}"), dst: &s, expected: "{}"},
	}

	for i, tt := range rawStringTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if *tt.dst != tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, *tt.dst)
		}
	}

	rawBytesTests := []struct {
		src      pgtype.JSON
		dst      *[]byte
		expected []byte
	}{
		{src: pgtype.JSON("{}"), dst: &b, expected: []byte("{}")},
		{src: nil, dst: &b, expected: (([]byte)(nil))},
	}

	for i, tt := range rawBytesTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if !bytes.Equal(tt.expected, *tt.dst) {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, *tt.dst)
		}
	}

	var mapDst map[string]interface{}
	type structDst struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var strDst structDst

	unmarshalTests := []struct {
		src      pgtype.JSON
		dst      interface{}
		expected interface{}
	}{
		{src: pgtype.JSON(`{"foo":"bar"}`), dst: &mapDst, expected: map[string]interface{}{"foo": "bar"}},
		{src: pgtype.JSON(`{"name":"John","age":42}`), dst: &strDst, expected: structDst{Name: "John", Age: 42}},
	}
	for i, tt := range unmarshalTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if dst := reflect.ValueOf(tt.dst).Elem().Interface(); !reflect.DeepEqual(dst, tt.expected) {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, dst)
		}
	}

	pointerAllocTests := []struct {
		src      pgtype.JSON
		dst      **string
		expected *string
	}{
		{src: nil, dst: &ps, expected: ((*string)(nil))},
	}

	for i, tt := range pointerAllocTests {
		err := tt.src.AssignTo(tt.dst)
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if *tt.dst != tt.expected {
			t.Errorf("%d: expected %v to assign %v, but result was %v", i, tt.src, tt.expected, *tt.dst)
		}
	}
}

func TestJSONMarshalJSON(t *testing.T) {
	successfulTests := []struct {
		source pgtype.JSON
		result string
	}{
		{source: nil, result: "null"},
		{source: pgtype.JSON("{\"a\": 1}"), result: "{\"a\": 1}"},
	}
	for i, tt := range successfulTests {
		r, err := tt.source.MarshalJSON()
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if string(r) != tt.result {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, string(r))
		}
	}
}

func TestJSONUnmarshalJSON(t *testing.T) {
	successfulTests := []struct {
		source string
		result pgtype.JSON
	}{
		{source: "null", result: nil},
		{source: "{\"a\": 1}", result: pgtype.JSON("{\"a\": 1}")},
	}
	for i, tt := range successfulTests {
		var r pgtype.JSON
		err := r.UnmarshalJSON([]byte(tt.source))
		if err != nil {
			t.Errorf("%d: %v", i, err)
		}

		if string(r) != string(tt.result) {
			t.Errorf("%d: expected %v to convert to %v, but it was %v", i, tt.source, tt.result, r)
		}
	}
}
