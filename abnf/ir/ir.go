package ir

import (
	"fmt"
	"github.com/0x51-dev/upeg/parser"
	"strconv"
	"strings"
)

func ParseJsonPointer(n *parser.Node) ([]string, error) {
	if n.Name != "JsonPointer" {
		return nil, fmt.Errorf("expected JsonPointer, got %s", n.Name)
	}
	var tokens []string
	for _, n := range n.Children() {
		switch n.Name {
		case "ReferenceToken":
			v := n.Value()
			v = strings.ReplaceAll(v, "~1", "/")
			v = strings.ReplaceAll(v, "~0", "~") // Order important!
			tokens = append(tokens, v)
		default:
			return nil, fmt.Errorf("expected ReferenceToken, got %s", n.Name)
		}
	}
	return tokens, nil
}

type RelativeJsonPointer struct {
	NonNegativeInteger int
	IndexManipulation  *int
	JsonPointer        *[]string
}

func ParseRelativeJsonPointer(n *parser.Node) (*RelativeJsonPointer, error) {
	if n.Name != "RelativeJsonPointer" {
		return nil, fmt.Errorf("expected RelativeJsonPointer, got %s", n.Name)
	}
	var ptr RelativeJsonPointer
	for _, n := range n.Children() {
		switch n.Name {
		case "JsonPointer":
			v, err := ParseJsonPointer(n)
			if err != nil {
				return nil, err
			}
			ptr.JsonPointer = &v
		case "OriginSpecification":
			for _, n := range n.Children() {
				switch n.Name {
				case "NonNegativeInteger":
					v, err := strconv.Atoi(n.Value())
					if err != nil {
						return nil, err
					}
					ptr.NonNegativeInteger = v
				case "IndexManipulation":
					str := n.Value()
					v, err := strconv.Atoi(str[1:])
					if err != nil {
						return nil, err
					}
					switch str[0] {
					case '+':
					case '-':
						v = -v
					default:
						return nil, fmt.Errorf("expected + or -, got %c", str[0])
					}
					ptr.IndexManipulation = &v
				default:
					return nil, fmt.Errorf("expected NonNegativeInteger or IndexManipulation, got %s", n.Name)
				}
			}
		default:
			return nil, fmt.Errorf("expected JsonPointer or OriginSpecification, got %s", n.Name)
		}
	}
	return &ptr, nil
}
