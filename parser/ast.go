package parser

type Node interface{}

type BlockNode struct {
	Statements []Node
}

type IntegerNode struct {
	Value int
}
