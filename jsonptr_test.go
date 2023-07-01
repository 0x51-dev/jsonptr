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
		switch v := test.val.(type) {
		case []any:
			for i, v := range v {
				if v != val.([]any)[i] {
					t.Fatalf("expected %v, got %v", test.val, val)
				}
			}
		default:
			switch v := test.val.(type) {
			case string:
				if val != v {
					t.Fatalf("expected %v, got %v", test.val, val)
				}
			case int:
				if val != float64(v) { // JSON numbers are always float64.
					t.Fatalf("expected %v, got %v", test.val, val)
				}
			}
		}
	}
}
