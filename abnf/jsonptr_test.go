package abnf_test

import (
	"github.com/0x51-dev/jsonptr/abnf"
	"github.com/0x51-dev/jsonptr/abnf/ir"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"testing"
)

func TestRelativeJsonPointer(t *testing.T) {
	for _, test := range []struct {
		ptr      string
		expected []string
	}{
		{"/foo", []string{"foo"}},
		{"/foo/0", []string{"foo", "0"}},
		{"/", []string{""}},
		{"/a~1b", []string{"a/b"}},
		{"/c%d", []string{"c%d"}},
		{"/e^f", []string{"e^f"}},
		{"/g|h", []string{"g|h"}},
		{"/i\\j", []string{"i\\j"}},
		{"/k\"l", []string{"k\"l"}},
		{"/ ", []string{" "}},
		{"/m~0n", []string{"m~n"}},
	} {
		p, err := parser.New([]rune(test.ptr))
		if err != nil {
			t.Fatal(err)
		}
		n, err := p.Parse(op.And{abnf.JsonPointer, op.EOF{}})
		if err != nil {
			t.Fatal(err)
		}
		ptr, err := ir.ParseJsonPointer(n)
		if err != nil {
			t.Fatal(err)
		}
		if len(ptr) != len(test.expected) {
			t.Fatalf("expected %d tokens, got %d", len(test.expected), len(ptr))
		}
		for i, v := range test.expected {
			if ptr[i] != v {
				t.Fatalf("expected %s, got %s", v, ptr[i])
			}
		}
	}
}
