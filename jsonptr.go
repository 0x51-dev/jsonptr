package jsonptr

import (
	"fmt"
	"github.com/0x51-dev/jsonptr/abnf"
	"github.com/0x51-dev/jsonptr/abnf/ir"
	"github.com/0x51-dev/upeg/parser"
	"github.com/0x51-dev/upeg/parser/op"
	"strconv"
)

// DEPS: go install github.com/0x51-dev/upeg/cmd/abnf
//go:generate abnf --in=abnf/jsonptr.abnf --out=abnf/jsonptr.go --ignore=unescaped,escaped

type JsonPointer []string

func ParseJsonPointer(ptr string) (JsonPointer, error) {
	p, err := parser.New([]rune(ptr))
	if err != nil {
		return nil, err
	}
	n, err := p.Parse(op.And{abnf.JsonPointer, op.EOF{}})
	if err != nil {
		return nil, err
	}
	return ir.ParseJsonPointer(n)
}

func (ptr JsonPointer) Eval(document map[string]any) (any, error) {
	return ptr.evalMap(document)
}

func (ptr JsonPointer) evalAny(v any) (any, error) {
	switch v := v.(type) {
	case map[string]any:
		return ptr.evalMap(v)
	case []any:
		return ptr.evalArray(v)
	default:
		return nil, fmt.Errorf("expected map or array, got %T", v)
	}
}

func (ptr JsonPointer) evalArray(arr []any) (any, error) {
	i, err := strconv.Atoi(ptr[0])
	if err != nil {
		return nil, err
	}
	if i < 0 || i >= len(arr) {
		return nil, fmt.Errorf("index %d out of bounds", i)
	}
	v := arr[i]
	if len(ptr) == 1 {
		return v, nil
	}
	return ptr[1:].evalAny(v)
}

func (ptr JsonPointer) evalMap(document map[string]any) (any, error) {
	if len(ptr) == 0 {
		return nil, nil
	}
	v, ok := document[ptr[0]]
	if !ok {
		return nil, fmt.Errorf("key %s not found", ptr[0])
	}
	if len(ptr) == 1 {
		return v, nil
	}
	return ptr[1:].evalAny(v)
}
