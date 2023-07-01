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

type RelativeJsonPointer ir.RelativeJsonPointer

func ParseRelativeJsonPointer(ptr string) (*RelativeJsonPointer, error) {
	p, err := parser.New([]rune(ptr))
	if err != nil {
		return nil, err
	}
	n, err := p.Parse(op.And{abnf.RelativeJsonPointer, op.EOF{}})
	if err != nil {
		return nil, err
	}
	r, err := ir.ParseRelativeJsonPointer(n)
	if err != nil {
		return nil, err
	}
	return (*RelativeJsonPointer)(r), nil
}

func (ptr RelativeJsonPointer) Eval(start JsonPointer, document map[string]any) (any, error) {
	current := start[:]

	// 1. Processing the non-negative-integer prefix.
	switch ptr.NonNegativeInteger {
	case 0: // Skip!
	default:
		for i := 0; i < ptr.NonNegativeInteger; i++ {
			if len(current) == 0 {
				// If the current referenced value is the root of the document, then evaluation fails.
				return nil, fmt.Errorf("referencing root document")
			}
			// If the referenced value is an item within an array/object, then the new referenced value is that array/object.
			current = current[:len(current)-1]
		}
	}

	// 2. Processing the index-manipulation suffix.
	if ptr.IndexManipulation != nil {
		v, err := strconv.Atoi(current[len(current)-1])
		if err != nil {
			return nil, fmt.Errorf("referencing non-integer key")
		}
		v += *ptr.IndexManipulation
		current[len(current)-1] = strconv.Itoa(v)
	}

	// 3. Processing the JSON Pointer suffix.
	if ptr.JsonPointer != nil {
		current = append(current, *ptr.JsonPointer...)
	} else {
		// The remainder of the Relative JSON Pointer is the character '#'.
		if len(current) == 0 {
			// If the current referenced value is the root of the document, then evaluation fails.
			return nil, fmt.Errorf("referencing root document")
		}
		p, err := current[:len(current)-1].Eval(document)
		if err != nil {
			return nil, err
		}
		switch p.(type) {
		case nil, map[string]any:
			return current[len(current)-1], nil
		case []any:
			v := current[len(current)-1]
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("referencing non-integer key")
			}
			return float64(i), nil
		default:
			return nil, fmt.Errorf("referencing non-object, non-array")
		}
	}

	return current.Eval(document)
}
