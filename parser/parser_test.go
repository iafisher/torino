package parser

import (
	"github.com/iafisher/torino/lexer"
	"testing"
)

func TestParseInteger(t *testing.T) {
	p := New(lexer.New("10"))

	tree := p.parseExpression()
	node, ok := tree.(*IntegerNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *IntegerNode, got %T", tree)
	}

	if node.Value != 10 {
		t.Fatalf("Wrong integer value: expected 10, got %d", node.Value)
	}
}

func TestParseString(t *testing.T) {
	p := New(lexer.New("\"hello\\n\""))

	tree := p.parseExpression()
	node, ok := tree.(*StringNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *StringNode, got %T", tree)
	}

	if node.Value != "hello\n" {
		t.Fatalf("Wrong string value: expected \"hello\\n\", got %q", node.Value)
	}
}

func TestParseSymbol(t *testing.T) {
	p := New(lexer.New("foo"))

	tree := p.parseExpression()
	node, ok := tree.(*SymbolNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *SymbolNode, got %T", tree)
	}

	if node.Value != "foo" {
		t.Fatalf("Wrong symbol value: expected foo, got %s", node.Value)
	}
}
