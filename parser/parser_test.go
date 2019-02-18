package parser

import (
	"github.com/iafisher/torino/lexer"
	"testing"
)

func TestParseInteger(t *testing.T) {
	tree := parseExpressionHelper(t, "10")

	checkInteger(t, tree, 10)
}

func TestParseString(t *testing.T) {
	tree := parseExpressionHelper(t, "\"hello\\n\"")

	node, ok := tree.(*StringNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *StringNode, got %T", tree)
	}

	if node.Value != "hello\n" {
		t.Fatalf("Wrong string value: expected \"hello\\n\", got %q", node.Value)
	}
}

func TestParseSymbol(t *testing.T) {
	tree := parseExpressionHelper(t, "foo")

	checkSymbol(t, tree, "foo")
}

func TestParseBool(t *testing.T) {
	tree := parseExpressionHelper(t, "true")

	node, ok := tree.(*BoolNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *BoolNode, got %T", tree)
	}

	if !node.Value {
		t.Fatalf("Wrong boolean value: expected true, got false")
	}
}

func TestParseLet(t *testing.T) {
	tree := parseStatementHelper(t, "let x = 10")

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
	tree := parseExpressionHelper(t, "5 + x")

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
	tree := parseExpressionHelper(t, "5 * 2 + 4")

	node := checkInfix(t, tree, "+")
	checkInteger(t, node.Right, 4)

	left := checkInfix(t, node.Left, "*")
	checkInteger(t, left.Left, 5)
	checkInteger(t, left.Right, 2)
}

func TestParseParentheses(t *testing.T) {
	tree := parseExpressionHelper(t, "5 * (2 + 4)")

	node := checkInfix(t, tree, "*")
	checkInteger(t, node.Left, 5)

	right := checkInfix(t, node.Right, "+")
	checkInteger(t, right.Left, 2)
	checkInteger(t, right.Right, 4)
}

func TestParseParentheses2(t *testing.T) {
	tree := parseExpressionHelper(t, "(2 + 4) * 5")

	node := checkInfix(t, tree, "*")
	checkInteger(t, node.Right, 5)

	left := checkInfix(t, node.Left, "+")
	checkInteger(t, left.Left, 2)
	checkInteger(t, left.Right, 4)
}

func TestParseCallExpression(t *testing.T) {
	tree := parseExpressionHelper(t, "f(x)")

	callNode := checkCall(t, tree, "f", 1)
	checkSymbol(t, callNode.Arglist[0], "x")
}

func TestParseCallExpressionWithNoArgs(t *testing.T) {
	tree := parseExpressionHelper(t, "f()")

	checkCall(t, tree, "f", 0)
}

func TestParseComplexCallExpression(t *testing.T) {
	tree := parseExpressionHelper(t, "7 * add(4 + 2, x - 1) / 10")

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

func TestParseCallExpressionWithNonSymbol(t *testing.T) {
	tree := parseExpressionHelper(t, "(x + y)()")

	callNode := checkCall(t, tree, "", 0)
	addNode := checkInfix(t, callNode.Func, "+")
	checkSymbol(t, addNode.Left, "x")
	checkSymbol(t, addNode.Right, "y")
}

func TestParseAssignNode(t *testing.T) {
	tree := parseStatementHelper(t, "x = x + 1")

	assignNode, ok := tree.(*AssignNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *AssignNode, got %T", tree)
	}

	checkSymbol(t, assignNode.Destination, "x")
	addNode := checkInfix(t, assignNode.Value, "+")
	checkSymbol(t, addNode.Left, "x")
	checkInteger(t, addNode.Right, 1)
}

func TestParseForLoop(t *testing.T) {
	tree := parseStatementHelper(t, "for c in string {\nprint(c)\n}")

	node, ok := tree.(*ForNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *ForNode, got %T", tree)
	}

	checkSymbol(t, node.Symbol, "c")
	checkSymbol(t, node.Iter, "string")

	if len(node.Block.Statements) != 1 {
		t.Fatalf("Wrong number of statements in block: expected 1, got %d",
			len(node.Block.Statements))
	}

	eStmt, ok := node.Block.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("Wrong AST type: expected *ExpressionStatement, got %T",
			node.Block.Statements[0])
	}

	callNode := checkCall(t, eStmt.Expr, "print", 1)
	checkSymbol(t, callNode.Arglist[0], "c")
}

func TestParseWhileLoop(t *testing.T) {
	tree := parseStatementHelper(t, "while x > 0 {\nprint(x)\nx = x - 1\n}")

	node, ok := tree.(*WhileNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *WhileNode, got %T", tree)
	}

	cmpNode := checkInfix(t, node.Cond, ">")
	checkSymbol(t, cmpNode.Left, "x")
	checkInteger(t, cmpNode.Right, 0)

	if len(node.Block.Statements) != 2 {
		t.Fatalf("Wrong number of statements in block: expected 2, got %d",
			len(node.Block.Statements))
	}

	eStmt, ok := node.Block.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("Wrong AST type: expected *ExpressionStatement, got %T",
			node.Block.Statements[0])
	}

	callNode := checkCall(t, eStmt.Expr, "print", 1)
	checkSymbol(t, callNode.Arglist[0], "x")

	assignNode, ok := node.Block.Statements[1].(*AssignNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *AssignNode, got %T",
			node.Block.Statements[1])
	}

	checkSymbol(t, assignNode.Destination, "x")
	subNode := checkInfix(t, assignNode.Value, "-")
	checkSymbol(t, subNode.Left, "x")
	checkInteger(t, subNode.Right, 1)
}

func TestParseIf(t *testing.T) {
	input := `
if x > 0 {
} elif x == 0 {
} else {
}
`
	tree := parseStatementHelper(t, input)
	node, ok := tree.(*IfNode)
	if !ok {
		t.Fatalf("Wrong AST type: expected *IfNode, got %T", tree)
	}

	if len(node.Clauses) != 2 {
		t.Fatalf("Wrong number of if-elif clauses: expected 2, got %d",
			len(node.Clauses))
	}

	cmpNode := checkInfix(t, node.Clauses[0].Cond, ">")
	checkSymbol(t, cmpNode.Left, "x")
	checkInteger(t, cmpNode.Right, 0)
}

// Helper functions

func parseHelper(input string) *BlockNode {
	p := New(lexer.New(input))
	return p.Parse()
}

func parseExpressionHelper(t *testing.T, input string) Expression {
	return extractExpression(t, parseHelper(input))
}

func parseStatementHelper(t *testing.T, input string) Statement {
	return extractStatement(t, parseHelper(input))
}

func extractExpression(t *testing.T, bn *BlockNode) Expression {
	if len(bn.Statements) != 1 {
		t.Fatalf("Wrong number of statements: expected 1, got %d", len(bn.Statements))
	}

	eStmt, ok := bn.Statements[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("Wrong AST type: expected *ExpressionStatement, got %T",
			bn.Statements[0])
	}

	return eStmt.Expr
}

func extractStatement(t *testing.T, bn *BlockNode) Statement {
	if len(bn.Statements) != 1 {
		t.Fatalf("Wrong number of statements: expected 1, got %d", len(bn.Statements))
	}
	return bn.Statements[0]
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

	if sym != "" {
		checkSymbol(t, callNode.Func, sym)
	}

	if len(callNode.Arglist) != nargs {
		t.Fatalf("Wrong number of args: expected %d, got %d", nargs, len(callNode.Arglist))
	}

	return callNode
}
