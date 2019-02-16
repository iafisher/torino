package parser

type Node interface{}
type Expression interface {
	Node
	expressionNode()
}
type Statement interface {
	Node
	statementNode()
}

type BlockNode struct {
	Statements []Node
}

type ExpressionStatement struct {
	Expr Expression
}

func (n *ExpressionStatement) statementNode() {}

type LetNode struct {
	Destination *SymbolNode
	Value       Expression
}

func (n *LetNode) statementNode() {}

type IntegerNode struct {
	Value int64
}

func (n *IntegerNode) expressionNode() {}

type BoolNode struct {
	Value bool
}

func (n *BoolNode) expressionNode() {}

type StringNode struct {
	Value string
}

func (n *StringNode) expressionNode() {}

type SymbolNode struct {
	Value string
}

func (n *SymbolNode) expressionNode() {}
