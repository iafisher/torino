package parser

import (
	"github.com/iafisher/torino/lexer"
	"testing"
)

func TestParseInteger(t *testing.T) {
	p := New(lexer.New("10"))

	tree := p.Parse()
	node, ok := tree.(*IntegerNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected integer, got %T", tree)
	}

	if node.Value != 10 {
		t.Fatalf("Wrong integer value: expected 10, got %d", node.Value)
	}
}
