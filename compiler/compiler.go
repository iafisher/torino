/* Compile ASTs into bytecode instructions.

Author:  Ian Fisher (iafisher@protonmail.com)
Version: February 2019
*/
package compiler

import (
	"fmt"
	"github.com/iafisher/torino/data"
	"github.com/iafisher/torino/parser"
)

// TODO: I might not even need this struct.
type Compiler struct {
}

func New() *Compiler {
	return &Compiler{}
}

func (cmp *Compiler) Compile(ast *parser.BlockNode) []*Instruction {
	program := []*Instruction{}
	for _, stmt := range ast.Statements {
		program = append(program, cmp.compileStatement(stmt)...)
	}
	return program
}

func (cmp *Compiler) compileStatement(stmt parser.Statement) []*Instruction {
	switch v := stmt.(type) {
	case *parser.ExpressionStatement:
		return cmp.compileExpression(v.Expr)
	case *parser.LetNode:
		return cmp.compileLet(v)
	case *parser.AssignNode:
		return cmp.compileAssign(v)
	case *parser.IfNode:
		return cmp.compileIf(v)
	default:
		panic("compileStatement - unknown statement type")
	}
}

func (cmp *Compiler) compileLet(node *parser.LetNode) []*Instruction {
	insts := cmp.compileExpression(node.Value)
	return append(insts, NewInst("STORE_NAME", &data.TorinoString{node.Destination.Value}))
}

func (cmp *Compiler) compileAssign(node *parser.AssignNode) []*Instruction {
	insts := cmp.compileExpression(node.Value)
	return append(insts, NewInst("ASSIGN_NAME", &data.TorinoString{node.Destination.Value}))
}

func (cmp *Compiler) compileIf(ifNode *parser.IfNode) []*Instruction {
	panic("not implemented!")
}

func (cmp *Compiler) compileExpression(expr parser.Expression) []*Instruction {
	insts := []*Instruction{}
	switch v := expr.(type) {
	case *parser.IntegerNode:
		return append(insts, NewInst("PUSH_CONST", &data.TorinoInt{v.Value}))
	case *parser.SymbolNode:
		return append(insts, NewInst("PUSH_NAME", &data.TorinoString{v.Value}))
	case *parser.BoolNode:
		return append(insts, NewInst("PUSH_CONST", &data.TorinoBool{v.Value}))
	case *parser.StringNode:
		return append(insts, NewInst("PUSH_CONST", &data.TorinoString{v.Value}))
	case *parser.InfixNode:
		return cmp.compileInfix(v)
	case *parser.PrefixNode:
		return cmp.compilePrefix(v)
	case *parser.CallNode:
		return cmp.compileCall(v)
	default:
		panic(fmt.Sprintf("compileExpression - unknown expression type %+v (%T)",
			expr, expr))
	}
}

func (cmp *Compiler) compileInfix(infixNode *parser.InfixNode) []*Instruction {
	insts := cmp.compileExpression(infixNode.Right)
	insts = append(insts, cmp.compileExpression(infixNode.Left)...)
	if infixNode.Op == "+" {
		return append(insts, NewInst("BINARY_ADD"))
	} else if infixNode.Op == "-" {
		return append(insts, NewInst("BINARY_SUB"))
	} else if infixNode.Op == "*" {
		return append(insts, NewInst("BINARY_MUL"))
	} else if infixNode.Op == "/" {
		return append(insts, NewInst("BINARY_DIV"))
	} else if infixNode.Op == "==" {
		return append(insts, NewInst("BINARY_EQ"))
	} else if infixNode.Op == ">" {
		return append(insts, NewInst("BINARY_GT"))
	} else if infixNode.Op == "<" {
		return append(insts, NewInst("BINARY_LT"))
	} else if infixNode.Op == ">=" {
		return append(insts, NewInst("BINARY_GE"))
	} else if infixNode.Op == "<=" {
		return append(insts, NewInst("BINARY_LE"))
	} else if infixNode.Op == "and" {
		return append(insts, NewInst("BINARY_AND"))
	} else if infixNode.Op == "or" {
		return append(insts, NewInst("BINARY_OR"))
	} else {
		panic(fmt.Sprintf("compileExpression - unknown infix operator %s", infixNode.Op))
	}
}

func (cmp *Compiler) compilePrefix(prefixNode *parser.PrefixNode) []*Instruction {
	insts := cmp.compileExpression(prefixNode.Arg)
	if prefixNode.Op == "-" {
		return append(insts, NewInst("UNARY_MINUS"))
	} else {
		panic(fmt.Sprintf("compileExpression - unknown prefix operator %s",
			prefixNode.Op))
	}
}

func (cmp *Compiler) compileCall(callNode *parser.CallNode) []*Instruction {
	insts := []*Instruction{}
	for _, e := range callNode.Arglist {
		insts = append(insts, cmp.compileExpression(e)...)
	}
	insts = append(insts, cmp.compileExpression(callNode.Func)...)
	nargs := int64(len(callNode.Arglist))
	return append(insts, NewInst("CALL_FUNCTION", &data.TorinoInt{nargs}))
}
