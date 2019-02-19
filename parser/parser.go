/* Grammar of Torino.

	start := block

	block := (stmt NEWLINE)*
	stmt  := let | fn | for | while | if | break | continue | return | expr

	let      := LET SYMBOL ASSIGN expr
	fn       := FN SYMBOL LPAREN params? RPAREN brace-block
	for      := FOR SYMBOL IN expr brace-block
	while    := WHILE expr brace-block
	if       := IF expr brace-block elif* else?
	elif     := ELIF expr brace-block
	else     := ELSE brace-block
	break    := BREAK
	continue := CONTINUE
	return   := RETURN expr?

	brace-block := LBRACE NEWLINE block RBRACE

	expr  := infix | call | pexpr | list | map | INT | STRING | SYMBOL | TRUE | FALSE
	pexpr := LPAREN expr RPAREN
	infix := expr OP expr
	call  := SYMBOL LPAREN args? RPAREN
	list  := LBRACKET args? RBRACKET
	map   := LBRACKET mapargs? RBRACKET

	params  := (SYMBOL COMMA)* SYMBOL
	args    := (expr COMMA)* expr
	mapargs := (maparg COMMA)* maparg
	maparg  := expr COLON expr

Infix operators have the usual precedence.

Author:  Ian Fisher (iafisher@protonmail.com)
Version: February 2019
*/
package parser

import (
	"fmt"
	"github.com/iafisher/torino/lexer"
	"strconv"
)

type Parser struct {
	lexer *lexer.Lexer

	curToken *lexer.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l, nil}
	p.nextToken()
	return p
}

func (p *Parser) Parse() *BlockNode {
	p.skipNewlines()
	return p.parseBlock(true)
}

func (p *Parser) parseBlock(topLevel bool) *BlockNode {
	statements := []Statement{}
	for {
		stmt := p.parseStatement(topLevel)
		statements = append(statements, stmt)

		if p.checkCurToken(lexer.TOKEN_NEWLINE) {
			p.skipNewlines()
		} else if p.checkCurToken(lexer.TOKEN_EOF) || p.checkCurToken(lexer.TOKEN_RBRACE) {
			break
		} else {
			panic(fmt.Sprintf("parseBlock - unexpected token %s", p.curToken.Type))
		}

		if p.checkCurToken(lexer.TOKEN_EOF) || p.checkCurToken(lexer.TOKEN_RBRACE) {
			break
		}
	}
	return &BlockNode{statements}
}

func (p *Parser) parseStatement(topLevel bool) Statement {
	if p.checkCurToken(lexer.TOKEN_LET) {
		return p.parseLetStatement()
	} else if p.checkCurToken(lexer.TOKEN_FOR) {
		return p.parseForStatement()
	} else if p.checkCurToken(lexer.TOKEN_WHILE) {
		return p.parseWhileStatement()
	} else if p.checkCurToken(lexer.TOKEN_IF) {
		return p.parseIfStatement()
	} else if p.checkCurToken(lexer.TOKEN_RETURN) {
		return p.parseReturnStatement()
	} else if p.checkCurToken(lexer.TOKEN_FN) {
		if !topLevel {
			panic("parseStatement - function declarations must be at top level")
		}
		return p.parseFnStatement()
	} else if p.checkCurToken(lexer.TOKEN_BREAK) {
		p.nextToken()
		return &BreakNode{}
	} else if p.checkCurToken(lexer.TOKEN_CONTINUE) {
		p.nextToken()
		return &ContinueNode{}
	} else {
		expr := p.parseExpression(PREC_LOWEST)
		if p.checkCurToken(lexer.TOKEN_ASSIGN) {
			sym, ok := expr.(*SymbolNode)
			if !ok {
				panic("parseStatement - cannot assign to non-symbol")
			}
			p.nextToken()
			lhs := p.parseExpression(PREC_LOWEST)
			return &AssignNode{sym, lhs}
		} else {
			return &ExpressionStatement{expr}
		}
	}
}

func (p *Parser) parseLetStatement() Statement {
	p.nextToken()
	if p.checkCurToken(lexer.TOKEN_SYMBOL) {
		dest := &SymbolNode{p.curToken.Value}
		p.nextToken()
		if !p.checkCurToken(lexer.TOKEN_ASSIGN) {
			panic("parseLetStatement - expected =")
		}
		p.nextToken()
		v := p.parseExpression(PREC_LOWEST)
		return &LetNode{dest, v}
	} else {
		panic("parseLetStatement - expected symbol")
	}
}

func (p *Parser) parseForStatement() Statement {
	p.nextToken()
	if p.checkCurToken(lexer.TOKEN_SYMBOL) {
		sym := &SymbolNode{p.curToken.Value}
		p.nextToken()
		if !p.checkCurToken(lexer.TOKEN_IN) {
			panic("parseForStatement - expected in")
		}
		p.nextToken()
		iter := p.parseExpression(PREC_LOWEST)
		body := p.parseBracedBlock()
		return &ForNode{sym, iter, body}
	} else {
		panic("parseForStatement - expected symbol")
	}
}

func (p *Parser) parseWhileStatement() Statement {
	p.nextToken()
	cond := p.parseExpression(PREC_LOWEST)
	body := p.parseBracedBlock()
	return &WhileNode{cond, body}
}

func (p *Parser) parseIfStatement() Statement {
	p.nextToken()
	clauses := []*IfClause{}

	// Parse the if block.
	cond := p.parseExpression(PREC_LOWEST)
	body := p.parseBracedBlock()
	clauses = append(clauses, &IfClause{cond, body})

	// Parse zero or more elif blocks.
	for p.checkCurToken(lexer.TOKEN_ELIF) {
		p.nextToken()
		elifCond := p.parseExpression(PREC_LOWEST)
		elifBody := p.parseBracedBlock()
		clauses = append(clauses, &IfClause{elifCond, elifBody})
	}

	// Parse an optional else block.
	var elseBody *BlockNode = nil
	if p.checkCurToken(lexer.TOKEN_ELSE) {
		p.nextToken()
		elseBody = p.parseBracedBlock()
	}

	return &IfNode{clauses, elseBody}
}

func (p *Parser) parseReturnStatement() Statement {
	p.nextToken()
	if p.checkCurToken(lexer.TOKEN_NEWLINE) || p.checkCurToken(lexer.TOKEN_EOF) {
		return &ReturnNode{nil}
	} else {
		return &ReturnNode{p.parseExpression(PREC_LOWEST)}
	}
}

func (p *Parser) parseFnStatement() Statement {
	p.nextToken()
	if !p.checkCurToken(lexer.TOKEN_SYMBOL) {
		panic("parseFnStatement - expected symbol")
	}
	sym := &SymbolNode{p.curToken.Value}

	p.nextToken()
	if !p.checkCurToken(lexer.TOKEN_LPAREN) {
		panic("parseFnStatement - expected (")
	}
	p.nextToken()
	params := p.parseParamList()

	body := p.parseBracedBlock()
	return &FnNode{sym, params, body}
}

func (p *Parser) parseExpression(precedence int) Expression {
	left := p.parsePrefix()

	for {
		// Keep consuming infix operators until we hit either a non-infix token or an
		// infix operator with a lower precedence.
		if infixPrecedence, ok := precedenceMap[p.curToken.Type]; ok {
			if precedence < infixPrecedence {
				if p.curToken.Type == lexer.TOKEN_LPAREN {
					p.nextToken()
					arglist := p.parseArglist()
					left = &CallNode{left, arglist}
				} else {
					left = p.parseInfix(left, getPrecedence(p.curToken.Type))
				}
			} else {
				break
			}
		} else {
			break
		}
	}

	return left
}

func (p *Parser) parsePrefix() Expression {
	typ := p.curToken.Type
	val := p.curToken.Value
	p.nextToken()
	if typ == lexer.TOKEN_INT {
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			panic("parsePrefix - could not parse integer token")
		}
		return &IntegerNode{int(v)}
	} else if typ == lexer.TOKEN_STRING {
		return &StringNode{val}
	} else if typ == lexer.TOKEN_SYMBOL {
		return &SymbolNode{val}
	} else if typ == lexer.TOKEN_TRUE {
		return &BoolNode{true}
	} else if typ == lexer.TOKEN_FALSE {
		return &BoolNode{false}
	} else if typ == lexer.TOKEN_LPAREN {
		expr := p.parseExpression(PREC_LOWEST)
		if !p.checkCurToken(lexer.TOKEN_RPAREN) {
			panic("parsePrefix - expected )")
		}
		p.nextToken()
		return expr
	} else if typ == lexer.TOKEN_MINUS {
		expr := p.parseExpression(PREC_PREFIX)
		return &PrefixNode{val, expr}
	} else {
		panic(fmt.Sprintf("parsePrefix - unexpected token %s", p.curToken.Type))
	}
}

func (p *Parser) parseInfix(left Expression, precedence int) Expression {
	operator := p.curToken.Value
	p.nextToken()
	right := p.parseExpression(precedence)
	return &InfixNode{operator, left, right}
}

func (p *Parser) parseArglist() []Expression {
	arglist := []Expression{}
	// Special case for empty arglist
	if p.checkCurToken(lexer.TOKEN_RPAREN) {
		p.nextToken()
		return arglist
	}

	for {
		expr := p.parseExpression(PREC_LOWEST)
		arglist = append(arglist, expr)

		if p.checkCurToken(lexer.TOKEN_COMMA) {
			p.nextToken()
		} else if p.checkCurToken(lexer.TOKEN_RPAREN) {
			p.nextToken()
			break
		} else {
			panic(fmt.Sprintf("parseArglist - unexpected token %s", p.curToken.Type))
		}
	}
	return arglist
}

func (p *Parser) parseParamList() []*SymbolNode {
	paramlist := []*SymbolNode{}
	// Special case for empty list
	if p.checkCurToken(lexer.TOKEN_RPAREN) {
		p.nextToken()
		return paramlist
	}

	for {
		if !p.checkCurToken(lexer.TOKEN_SYMBOL) {
			panic(fmt.Sprintf("parseParamList - expected symbol, got %s",
				p.curToken.Type))
		}
		paramlist = append(paramlist, &SymbolNode{p.curToken.Value})

		p.nextToken()
		if p.checkCurToken(lexer.TOKEN_COMMA) {
			p.nextToken()
		} else if p.checkCurToken(lexer.TOKEN_RPAREN) {
			p.nextToken()
			break
		} else {
			panic(fmt.Sprintf("parseParamList - unexpected token %s", p.curToken.Type))
		}
	}
	return paramlist
}

func (p *Parser) parseBracedBlock() *BlockNode {
	if !p.checkCurToken(lexer.TOKEN_LBRACE) {
		panic("parseBracedBlock - expected {")
	}
	p.nextToken()

	p.skipNewlines()

	if p.checkCurToken(lexer.TOKEN_RBRACE) {
		p.nextToken()
		return &BlockNode{[]Statement{}}
	}

	block := p.parseBlock(false)

	if !p.checkCurToken(lexer.TOKEN_RBRACE) {
		panic("parseBracedBlock - expected }")
	}
	p.nextToken()

	return block
}

func (p *Parser) checkCurToken(expectedType string) bool {
	return p.curToken.Type == expectedType
}

func (p *Parser) nextToken() *lexer.Token {
	p.curToken = p.lexer.NextToken()
	return p.curToken
}

func (p *Parser) skipNewlines() {
	for p.checkCurToken(lexer.TOKEN_NEWLINE) {
		p.nextToken()
	}
}

func getPrecedence(tokType string) int {
	if prec, ok := precedenceMap[tokType]; ok {
		return prec
	} else {
		return PREC_LOWEST
	}
}

const (
	_ int = iota
	PREC_LOWEST
	PREC_OR
	PREC_AND
	PREC_CMP
	PREC_ADD_SUB
	PREC_MUL_DIV
	PREC_PREFIX
	PREC_CALL
)

var precedenceMap = map[string]int{
	lexer.TOKEN_EQ:       PREC_CMP,
	lexer.TOKEN_GT:       PREC_CMP,
	lexer.TOKEN_GE:       PREC_CMP,
	lexer.TOKEN_LT:       PREC_CMP,
	lexer.TOKEN_LE:       PREC_CMP,
	lexer.TOKEN_PLUS:     PREC_ADD_SUB,
	lexer.TOKEN_MINUS:    PREC_ADD_SUB,
	lexer.TOKEN_ASTERISK: PREC_MUL_DIV,
	lexer.TOKEN_SLASH:    PREC_MUL_DIV,
	lexer.TOKEN_LPAREN:   PREC_CALL,
	lexer.TOKEN_AND:      PREC_AND,
	lexer.TOKEN_OR:       PREC_OR,
}
