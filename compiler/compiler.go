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

func (cmp *Compiler) compileExpression(expr parser.Expression) []*Instruction {
	insts := []*Instruction{}
	switch v := expr.(type) {
	case *parser.IntegerNode:
		return append(insts, NewInst("PUSH_CONST", &data.TorinoInt{v.Value}))
	case *parser.SymbolNode:
		return append(insts, NewInst("PUSH_NAME", &data.TorinoString{v.Value}))
	case *parser.BoolNode:
		return append(insts, NewInst("PUSH_CONST", &data.TorinoBool{v.Value}))
	case *parser.InfixNode:
		insts = append(insts, cmp.compileExpression(v.Right)...)
		insts = append(insts, cmp.compileExpression(v.Left)...)
		if v.Op == "+" {
			return append(insts, NewInst("BINARY_ADD"))
		} else if v.Op == "-" {
			return append(insts, NewInst("BINARY_SUB"))
		} else if v.Op == "*" {
			return append(insts, NewInst("BINARY_MUL"))
		} else if v.Op == "/" {
			return append(insts, NewInst("BINARY_DIV"))
		} else if v.Op == "==" {
			return append(insts, NewInst("BINARY_EQ"))
		} else if v.Op == ">" {
			return append(insts, NewInst("BINARY_GT"))
		} else if v.Op == "<" {
			return append(insts, NewInst("BINARY_LT"))
		} else if v.Op == ">=" {
			return append(insts, NewInst("BINARY_GE"))
		} else if v.Op == "<=" {
			return append(insts, NewInst("BINARY_LE"))
		} else if v.Op == "and" {
			return append(insts, NewInst("BINARY_AND"))
		} else if v.Op == "or" {
			return append(insts, NewInst("BINARY_OR"))
		} else {
			panic(fmt.Sprintf("compileExpression - unknown infix operator %s", v.Op))
		}
	case *parser.PrefixNode:
		insts = append(insts, cmp.compileExpression(v.Arg)...)
		if v.Op == "-" {
			return append(insts, NewInst("UNARY_MINUS"))
		} else {
			panic(fmt.Sprintf("compileExpression - unknown prefix operator %s", v.Op))
		}
	case *parser.CallNode:
		for _, e := range v.Arglist {
			insts = append(insts, cmp.compileExpression(e)...)
		}
		insts = append(insts, cmp.compileExpression(v.Func)...)
		nargs := int64(len(v.Arglist))
		return append(insts, NewInst("CALL_FUNCTION", &data.TorinoInt{nargs}))
	default:
		panic(fmt.Sprintf("compileExpression - unknown expression type %+v (%T)",
			expr, expr))
	}
}
