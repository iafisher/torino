package parser

type Node interface{}
type Expression interface {
	Node
}
type Statement interface {
	Node
}

type BlockNode struct {
	Statements []Node
}

type IntegerNode struct {
	Value int64
}

type BoolNode struct {
	Value bool
}

type StringNode struct {
	Value string
}

type SymbolNode struct {
	Value string
}
