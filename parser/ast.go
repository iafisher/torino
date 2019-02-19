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
	Statements []Statement
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

type FnNode struct {
	Symbol *SymbolNode
	Params []*SymbolNode
	Body   *BlockNode
}

func (n *FnNode) statementNode() {}

type AssignNode struct {
	Destination *SymbolNode
	Value       Expression
}

func (n *AssignNode) statementNode() {}

type IfNode struct {
	Clauses []*IfClause
	Else    *BlockNode
}

type IfClause struct {
	Cond Expression
	Body *BlockNode
}

func (n *IfNode) statementNode() {}

type InfixNode struct {
	Op    string
	Left  Expression
	Right Expression
}

func (n *InfixNode) expressionNode() {}

type PrefixNode struct {
	Op  string
	Arg Expression
}

func (n *PrefixNode) expressionNode() {}

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

type WhileNode struct {
	Cond  Expression
	Block *BlockNode
}

func (n *WhileNode) statementNode() {}

type BreakNode struct{}

func (n *BreakNode) statementNode() {}

type ContinueNode struct{}

func (n *ContinueNode) statementNode() {}

type ReturnNode struct {
	Value Expression
}

func (n *ReturnNode) statementNode() {}

type IntegerNode struct {
	Value int
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
