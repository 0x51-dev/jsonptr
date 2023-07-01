package ir

import (
	"fmt"
	"github.com/0x51-dev/upeg/parser"
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
