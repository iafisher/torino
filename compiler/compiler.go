/* Compile ASTs into bytecode instructions.

Author:  Ian Fisher (iafisher@protonmail.com)
Version: February 2019
*/
package compiler

import (
	"errors"
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

func (cmp *Compiler) Compile(ast *parser.BlockNode) ([]*Instruction, error) {
	program := []*Instruction{}
	for _, stmt := range ast.Statements {
		stmtCode, err := cmp.compileStatement(stmt)
		if err != nil {
			return nil, err
		}
		program = append(program, stmtCode...)
	}
	return program, nil
}

func (cmp *Compiler) compileStatement(stmt parser.Statement) ([]*Instruction, error) {
	switch v := stmt.(type) {
	case *parser.ExpressionStatement:
		return cmp.compileExpression(v.Expr)
	case *parser.LetNode:
		return cmp.compileLet(v)
	case *parser.AssignNode:
		return cmp.compileAssign(v)
	case *parser.IfNode:
		return cmp.compileIf(v)
	case *parser.FnNode:
		return cmp.compileFn(v)
	case *parser.ReturnNode:
		return cmp.compileReturn(v)
	case *parser.WhileNode:
		return cmp.compileWhile(v)
	case *parser.ForNode:
		return cmp.compileFor(v)
	default:
		return nil, errors.New(fmt.Sprintf("unknown statement type %T", stmt))
	}
}

func (cmp *Compiler) compileExpression(expr parser.Expression) ([]*Instruction, error) {
	insts := []*Instruction{}
	switch v := expr.(type) {
	case *parser.IntegerNode:
		return append(insts, NewInst("PUSH_CONST", &data.TorinoInt{v.Value})), nil
	case *parser.SymbolNode:
		return append(insts, NewInst("PUSH_NAME", &data.TorinoString{v.Value})), nil
	case *parser.BoolNode:
		return append(insts, NewInst("PUSH_CONST", &data.TorinoBool{v.Value})), nil
	case *parser.StringNode:
		return append(insts, NewInst("PUSH_CONST", &data.TorinoString{v.Value})), nil
	case *parser.ListNode:
		return cmp.compileList(v)
	case *parser.MapNode:
		return cmp.compileMap(v)
	case *parser.InfixNode:
		return cmp.compileInfix(v)
	case *parser.PrefixNode:
		return cmp.compilePrefix(v)
	case *parser.CallNode:
		return cmp.compileCall(v)
	case *parser.IndexNode:
		return cmp.compileIndex(v)
	default:
		return nil, errors.New(fmt.Sprintf("unknown expression type %+v (%T)", expr, expr))
	}
}

func (cmp *Compiler) compileLet(node *parser.LetNode) ([]*Instruction, error) {
	insts, err := cmp.compileExpression(node.Value)
	if err != nil {
		return nil, err
	}
	return append(insts, NewInst("STORE_NAME", &data.TorinoString{node.Destination.Value})), nil
}

func (cmp *Compiler) compileAssign(node *parser.AssignNode) ([]*Instruction, error) {
	insts, err := cmp.compileExpression(node.Value)
	if err != nil {
		return nil, err
	}

	insts = append(insts, NewInst("ASSIGN_NAME", &data.TorinoString{node.Destination.Value}))
	return insts, nil
}

func (cmp *Compiler) compileIf(ifNode *parser.IfNode) ([]*Instruction, error) {
	compiledBodies := make([][]*Instruction, 0, len(ifNode.Clauses))
	compiledConds := make([][]*Instruction, 0, len(ifNode.Clauses))
	endJump := 0
	for _, clause := range ifNode.Clauses {
		cond, err := cmp.compileExpression(clause.Cond)
		if err != nil {
			return nil, err
		}

		body, err := cmp.Compile(clause.Body)
		if err != nil {
			return nil, err
		}

		compiledConds = append(compiledConds, cond)
		compiledBodies = append(compiledBodies, body)

		endJump += len(cond) + len(body) + 1
	}

	var elseCode []*Instruction
	var err error
	if ifNode.Else != nil {
		elseCode, err = cmp.Compile(ifNode.Else)
		if err != nil {
			return nil, err
		}

		endJump += len(elseCode)
	}

	insts := []*Instruction{}
	for i, code := range compiledBodies {
		endJump -= (len(code) + len(compiledConds[i]) + 1)

		insts = append(insts, compiledConds[i]...)
		jump := &data.TorinoInt{len(code) + 2}

		insts = append(insts, NewInst("REL_JUMP_IF_FALSE", jump))

		insts = append(insts, code...)
		insts = append(insts, NewInst("REL_JUMP", &data.TorinoInt{endJump + 1}))
	}

	if elseCode != nil {
		insts = append(insts, elseCode...)
	}

	return insts, nil
}

func (cmp *Compiler) compileFn(fnNode *parser.FnNode) ([]*Instruction, error) {
	insts := []*Instruction{}

	body, err := cmp.Compile(fnNode.Body)
	if err != nil {
		return nil, err
	}

	insts = append(insts, NewInst("PUSH_CONST", &TorinoFunction{fnNode.Params, body}))
	return append(insts, NewInst("STORE_NAME", &data.TorinoString{fnNode.Symbol.Value})), nil
}

func (cmp *Compiler) compileReturn(returnNode *parser.ReturnNode) ([]*Instruction, error) {
	insts, err := cmp.compileExpression(returnNode.Value)
	if err != nil {
		return nil, err
	}
	return append(insts, NewInst("RETURN_VALUE")), nil
}

func (cmp *Compiler) compileWhile(whileNode *parser.WhileNode) ([]*Instruction, error) {
	cond, err := cmp.compileExpression(whileNode.Cond)
	if err != nil {
		return nil, err
	}

	body, err := cmp.Compile(whileNode.Block)
	if err != nil {
		return nil, err
	}

	insts := cond

	endJump := &data.TorinoInt{len(body) + 2}
	insts = append(insts, NewInst("REL_JUMP_IF_FALSE", endJump))

	insts = append(insts, body...)
	startJump := &data.TorinoInt{-(len(cond) + len(body) + 1)}
	return append(insts, NewInst("REL_JUMP", startJump)), nil
}

func (cmp *Compiler) compileFor(forNode *parser.ForNode) ([]*Instruction, error) {
	// Compile the iterator value and the body.
	iterCode, err := cmp.compileExpression(forNode.Iter)
	if err != nil {
		return nil, err
	}

	bodyCode, err := cmp.Compile(forNode.Block)
	if err != nil {
		return nil, err
	}

	insts := []*Instruction{
		// Just need to initialize the loop variable with some throwaway value,
		// so we can subsequently use ASSIGN_NAME without worrying about an
		// undefined symbol error. Hacky but it works.
		NewInst("PUSH_CONST", &data.TorinoInt{0}),
		NewInst("STORE_NAME", &data.TorinoString{forNode.Symbol.Value}),
	}

	insts = append(insts, iterCode...)

	endJump := len(bodyCode) + 3
	insts = append(insts, NewInst("LIST_NEXT", &data.TorinoInt{endJump}))
	insts = append(insts, NewInst("ASSIGN_NAME", &data.TorinoString{forNode.Symbol.Value}))
	insts = append(insts, bodyCode...)
	startJump := -(len(bodyCode) + 2)
	insts = append(insts, NewInst("REL_JUMP", &data.TorinoInt{startJump}))
	return insts, nil
}

func (cmp *Compiler) compileList(listNode *parser.ListNode) ([]*Instruction, error) {
	insts := []*Instruction{}
	for i := len(listNode.Values) - 1; i >= 0; i-- {
		exprCode, err := cmp.compileExpression(listNode.Values[i])
		if err != nil {
			return nil, err
		}

		insts = append(insts, exprCode...)
	}

	insts = append(insts, NewInst("MAKE_LIST", &data.TorinoInt{len(listNode.Values)}))
	return insts, nil
}

func (cmp *Compiler) compileMap(mapNode *parser.MapNode) ([]*Instruction, error) {
	insts := []*Instruction{}
	for _, item := range mapNode.Values {
		keyCode, err := cmp.compileExpression(item.Key)
		if err != nil {
			return nil, err
		}

		valCode, err := cmp.compileExpression(item.Value)
		if err != nil {
			return nil, err
		}

		insts = append(insts, keyCode...)
		insts = append(insts, valCode...)
	}
	insts = append(insts, NewInst("MAKE_MAP", &data.TorinoInt{len(mapNode.Values)}))
	return insts, nil
}

func (cmp *Compiler) compileInfix(infixNode *parser.InfixNode) ([]*Instruction, error) {
	insts, err := cmp.compileExpression(infixNode.Right)
	if err != nil {
		return nil, err
	}

	leftCode, err := cmp.compileExpression(infixNode.Left)
	if err != nil {
		return nil, err
	}

	insts = append(insts, leftCode...)
	if infixNode.Op == "+" {
		return append(insts, NewInst("BINARY_ADD")), nil
	} else if infixNode.Op == "-" {
		return append(insts, NewInst("BINARY_SUB")), nil
	} else if infixNode.Op == "*" {
		return append(insts, NewInst("BINARY_MUL")), nil
	} else if infixNode.Op == "/" {
		return append(insts, NewInst("BINARY_DIV")), nil
	} else if infixNode.Op == "==" {
		return append(insts, NewInst("BINARY_EQ")), nil
	} else if infixNode.Op == ">" {
		return append(insts, NewInst("BINARY_GT")), nil
	} else if infixNode.Op == "<" {
		return append(insts, NewInst("BINARY_LT")), nil
	} else if infixNode.Op == ">=" {
		return append(insts, NewInst("BINARY_GE")), nil
	} else if infixNode.Op == "<=" {
		return append(insts, NewInst("BINARY_LE")), nil
	} else if infixNode.Op == "and" {
		return append(insts, NewInst("BINARY_AND")), nil
	} else if infixNode.Op == "or" {
		return append(insts, NewInst("BINARY_OR")), nil
	} else {
		return nil, errors.New(fmt.Sprintf("unknown infix operator %s", infixNode.Op))
	}
}

func (cmp *Compiler) compilePrefix(prefixNode *parser.PrefixNode) ([]*Instruction, error) {
	insts, err := cmp.compileExpression(prefixNode.Arg)
	if err != nil {
		return nil, err
	}

	if prefixNode.Op == "-" {
		return append(insts, NewInst("UNARY_MINUS")), nil
	} else {
		return nil, errors.New(fmt.Sprintf("unknown prefix operator %s", prefixNode.Op))
	}
}

func (cmp *Compiler) compileCall(callNode *parser.CallNode) ([]*Instruction, error) {
	insts := []*Instruction{}
	for _, e := range callNode.Arglist {
		exprCode, err := cmp.compileExpression(e)
		if err != nil {
			return nil, err
		}

		insts = append(insts, exprCode...)
	}
	fCode, err := cmp.compileExpression(callNode.Func)
	if err != nil {
		return nil, err
	}

	insts = append(insts, fCode...)
	nargs := len(callNode.Arglist)
	insts = append(insts, NewInst("CALL_FUNCTION", &data.TorinoInt{nargs}))
	return insts, nil
}

func (cmp *Compiler) compileIndex(indexNode *parser.IndexNode) ([]*Instruction, error) {
	insts, err := cmp.compileExpression(indexNode.Index)
	if err != nil {
		return nil, err
	}

	indexedCode, err := cmp.compileExpression(indexNode.Indexed)
	if err != nil {
		return nil, err
	}

	insts = append(insts, indexedCode...)
	return append(insts, NewInst("BINARY_INDEX")), nil
}

// Some data types, defined here because they use the compiler.Instruction object,
// which would create a circular import path if they were defined in the data
// package.

type TorinoCode struct {
	Code []*Instruction
}

func (t *TorinoCode) Torino() {}

func (t *TorinoCode) String() string {
	return "<code object>"
}

func (t *TorinoCode) Repr() string {
	return t.String()
}

type TorinoFunction struct {
	Params []*parser.SymbolNode
	Body   []*Instruction
}

func (t *TorinoFunction) Torino() {}

func (t *TorinoFunction) String() string {
	return "<function object>"
}

func (t *TorinoFunction) Repr() string {
	return t.String()
}
