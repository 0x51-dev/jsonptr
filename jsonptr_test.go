package jsonptr_test

import (
	"encoding/json"
	"github.com/0x51-dev/jsonptr"
	"testing"
)

func TestJsonPointer_Eval(t *testing.T) {
	documentStr := `{
    	"foo": ["bar", "baz"],
    	"": 0,
    	"a/b": 1,
    	"c%d": 2,
    	"e^f": 3,
    	"g|h": 4,
    	"i\\j": 5,
    	"k\"l": 6,
   		" ": 7,
    	"m~n": 8
   	}`
	var document map[string]any
	if err := json.Unmarshal([]byte(documentStr), &document); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		ptr string
		val any
	}{
		{"/foo", []any{"bar", "baz"}},
		{"/foo/0", "bar"},
		{"/", 0},
		{"/a~1b", 1},
		{"/c%d", 2},
		{"/e^f", 3},
		{"/g|h", 4},
		{"/i\\j", 5},
		{"/k\"l", 6},
		{"/ ", 7},
		{"/m~0n", 8},
	} {
		ptr, err := jsonptr.ParseJsonPointer(test.ptr)
		if err != nil {
			t.Fatal(err)
		}
		val, err := ptr.Eval(document)
		if err != nil {
			t.Fatal(err)
		}
		cmp(t, test.val, val)
	}
}

func TestRelativeJsonPointer_Eval(t *testing.T) {
	documentStr := `{
		"foo": ["bar", "baz", "biz"],
		"highly": {
			"nested": {
				"objects": true
			}
		}
	}`
	var document map[string]any
	if err := json.Unmarshal([]byte(documentStr), &document); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		path  string
		tests []struct {
			ptr string
			val any
		}
	}{
		{
			path: "/foo/1",
			tests: []struct {
				ptr string
				val any
			}{
				{"0", "baz"},
				{"1/0", "bar"},
				{"0-1", "bar"},
				{"2/highly/nested/objects", true},
				{"0#", 1},
				{"0+1#", 2},
				{"1#", "foo"},
			},
		},
		{
			path: "/highly/nested",
			tests: []struct {
				ptr string
				val any
			}{
				{"0/objects", true},
				{"1/nested/objects", true},
				{"2/foo/0", "bar"},
				{"0#", "nested"},
				{"1#", "highly"},
			},
		},
	} {
		t.Run(test.path, func(t *testing.T) {
			path := test.path
			for _, test := range test.tests {
				start, err := jsonptr.ParseJsonPointer(path)
				if err != nil {
					t.Fatal(err)
				}
				ptr, err := jsonptr.ParseRelativeJsonPointer(test.ptr)
				if err != nil {
					t.Fatal(err)
				}
				val, err := ptr.Eval(start, document)
				if err != nil {
					t.Fatal(err)
				}
				cmp(t, test.val, val)
			}
		})
	}
}

func cmp(t *testing.T, a, b any) {
	switch a := a.(type) {
	case []any:
		for i, a := range a {
			if a != b.([]any)[i] {
				t.Fatalf("expected %v, got %v", a, b)
			}
		}
	default:
		switch a := a.(type) {
		case string:
			if b != a {
				t.Fatalf("expected %v, got %v", a, b)
			}
		case int:
			if b != float64(a) { // JSON numbers are always float64.
				t.Fatalf("expected %v, got %v", a, b)
			}
		}
	}
}
