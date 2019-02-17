package parser

import (
	"github.com/iafisher/torino/lexer"
	"testing"
)

func TestParseInteger(t *testing.T) {
	p := New(lexer.New("10"))

	tree := p.parseExpression(PREC_LOWEST)

	checkInteger(t, tree, 10)
}

func TestParseString(t *testing.T) {
	p := New(lexer.New("\"hello\\n\""))

	tree := p.parseExpression(PREC_LOWEST)
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

	tree := p.parseExpression(PREC_LOWEST)

	checkSymbol(t, tree, "foo")
}

func TestParseBool(t *testing.T) {
	p := New(lexer.New("true"))

	tree := p.parseExpression(PREC_LOWEST)
	node, ok := tree.(*BoolNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *BoolNode, got %T", tree)
	}

	if !node.Value {
		t.Fatalf("Wrong boolean value: expected true, got false")
	}
}

func TestParseLet(t *testing.T) {
	p := New(lexer.New("let x = 10"))

	tree := p.parseStatement()
	node, ok := tree.(*LetNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *LetNode, got %T", tree)
	}

	if node.Destination.Value != "x" {
		t.Fatalf("Wrong destination value: expected x, got %s", node.Destination.Value)
	}

	checkInteger(t, node.Value, 10)
}

func TestParseTwoLets(t *testing.T) {
	p := New(lexer.New("let x = 10\nlet y = 20"))

	tree := p.Parse()
	if len(tree.Statements) != 2 {
		t.Fatalf("Expected 2 statements, got %d", len(tree.Statements))
	}

	stmt1, ok := tree.Statements[0].(*LetNode)
	if !ok {
		t.Fatalf("Wrong AST type for first statement: expected *LetNode, got %T",
			tree.Statements[0])
	}

	if stmt1.Destination.Value != "x" {
		t.Fatalf("Wrong destination value: expected x, got %s", stmt1.Destination.Value)
	}

	checkInteger(t, stmt1.Value, 10)

	stmt2, ok := tree.Statements[1].(*LetNode)
	if !ok {
		t.Fatalf("Wrong AST type for second statement: expected *LetNode, got %T",
			tree.Statements[1])
	}

	if stmt2.Destination.Value != "y" {
		t.Fatalf("Wrong destination value: expected x, got %s", stmt2.Destination.Value)
	}

	checkInteger(t, stmt2.Value, 20)
}

func TestParseSimpleArithmetic(t *testing.T) {
	p := New(lexer.New("5 + x"))

	tree := p.parseExpression(PREC_LOWEST)
	node, ok := tree.(*InfixNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *InfixNode, got %T", tree)
	}

	if node.Op != "+" {
		t.Fatalf("Wrong operator: expected +, got %s", node.Op)
	}

	checkInteger(t, node.Left, 5)
	checkSymbol(t, node.Right, "x")
}

func TestParseArithmeticPrecedence(t *testing.T) {
	p := New(lexer.New("5 * 2 + 4"))

	tree := p.parseExpression(PREC_LOWEST)

	node := checkInfix(t, tree, "+")
	checkInteger(t, node.Right, 4)

	left := checkInfix(t, node.Left, "*")
	checkInteger(t, left.Left, 5)
	checkInteger(t, left.Right, 2)
}

func TestParseParentheses(t *testing.T) {
	p := New(lexer.New("5 * (2 + 4)"))

	tree := p.parseExpression(PREC_LOWEST)

	node := checkInfix(t, tree, "*")
	checkInteger(t, node.Left, 5)

	right := checkInfix(t, node.Right, "+")
	checkInteger(t, right.Left, 2)
	checkInteger(t, right.Right, 4)
}

func TestParseCallExpression(t *testing.T) {
	p := New(lexer.New("f(x)"))

	tree := p.parseExpression(PREC_LOWEST)

	callNode := checkCall(t, tree, "f", 1)
	checkSymbol(t, callNode.Arglist[0], "x")
}

func TestParseCallExpressionWithNoArgs(t *testing.T) {
	p := New(lexer.New("f()"))

	tree := p.parseExpression(PREC_LOWEST)

	checkCall(t, tree, "f", 0)
}

func TestParseComplexCallExpression(t *testing.T) {
	p := New(lexer.New("7 * add(4 + 2, x - 1) / 10"))

	tree := p.parseExpression(PREC_LOWEST)

	divNode := checkInfix(t, tree, "/")
	checkInteger(t, divNode.Right, 10)

	mulNode := checkInfix(t, divNode.Left, "*")
	checkInteger(t, mulNode.Left, 7)

	callNode := checkCall(t, mulNode.Right, "add", 2)

	addNode := checkInfix(t, callNode.Arglist[0], "+")
	checkInteger(t, addNode.Left, 4)
	checkInteger(t, addNode.Right, 2)

	subNode := checkInfix(t, callNode.Arglist[1], "-")
	checkSymbol(t, subNode.Left, "x")
	checkInteger(t, subNode.Right, 1)
}

func checkInteger(t *testing.T, n Node, v int64) {
	intNode, ok := n.(*IntegerNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *IntegerNode, got %T", n)
	}

	if intNode.Value != v {
		t.Fatalf("Wrong value for integer: expected %d, got %d", v, intNode.Value)
	}
}

func checkSymbol(t *testing.T, n Node, v string) {
	symNode, ok := n.(*SymbolNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *SymbolNode, got %T", n)
	}

	if symNode.Value != v {
		t.Fatalf("Wrong value for symbol: expected %s, got %s", v, symNode.Value)
	}
}

func checkInfix(t *testing.T, n Node, op string) *InfixNode {
	infixNode, ok := n.(*InfixNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *InfixNode, got %T", n)
	}

	if infixNode.Op != op {
		t.Fatalf("Wrong operator: expected %s, got %s", op, infixNode.Op)
	}

	return infixNode
}

func checkCall(t *testing.T, n Node, sym string, nargs int) *CallNode {
	callNode, ok := n.(*CallNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *CallNode, got %T", n)
	}

	checkSymbol(t, callNode.Func, sym)

	if len(callNode.Arglist) != nargs {
		t.Fatalf("Wrong number of args: expected %d, got %d", nargs, len(callNode.Arglist))
	}

	return callNode
}
