/* Compile ASTs into bytecode instructions.

Author:  Ian Fisher (iafisher@protonmail.com)
Version: February 2019
*/
package compiler

import (
	"github.com/iafisher/torino/data"
	"github.com/iafisher/torino/parser"
)

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
	default:
		panic("compileStatement - unknown statement type")
	}
}

func (cmp *Compiler) compileLet(node *parser.LetNode) []*Instruction {
	insts := cmp.compileExpression(node.Value)
	return append(insts, NewInst("STORE_NAME", &data.TorinoString{node.Destination.Value}))
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
			panic("compileExpression - unknown operator")
		}
	default:
		panic("compileExpression - unknown expression type")
	}
}
