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

type InfixNode struct {
	Op    string
	Left  Expression
	Right Expression
}

func (n *InfixNode) expressionNode() {}

type CallNode struct {
	Func    Expression
	Arglist []Expression
}

func (n *CallNode) expressionNode() {}

type ForNode struct {
	Symbol *SymbolNode
	Iter   Expression
	Block  *BlockNode
}

func (n *ForNode) statementNode() {}

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
