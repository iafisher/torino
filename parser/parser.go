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
	lexer    *lexer.Lexer
	curToken *lexer.Token
	errors   []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l, nil, nil}
	p.nextToken()
	return p
}

func (p *Parser) Parse() (*BlockNode, bool) {
	p.skipNewlines()
	return p.parseBlock(true)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) recordError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseBlock(topLevel bool) (*BlockNode, bool) {
	statements := []Statement{}
	for {
		stmt, ok := p.parseStatement(topLevel)
		if !ok {
			return nil, false
		}
		statements = append(statements, stmt)

		if p.checkCurToken(lexer.TOKEN_NEWLINE) {
			p.skipNewlines()
		} else if p.checkCurToken(lexer.TOKEN_EOF) || p.checkCurToken(lexer.TOKEN_RBRACE) {
			break
		} else {
			err := fmt.Sprintf("unexpected token %s while parsing block", p.curToken.Type)
			p.recordError(err)
			return nil, false
		}

		if p.checkCurToken(lexer.TOKEN_EOF) || p.checkCurToken(lexer.TOKEN_RBRACE) {
			break
		}
	}

	return &BlockNode{statements}, true
}

func (p *Parser) parseStatement(topLevel bool) (Statement, bool) {
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
			p.recordError("function declarations must be at top level")
			return nil, false
		}
		return p.parseFnStatement()
	} else if p.checkCurToken(lexer.TOKEN_BREAK) {
		p.nextToken()
		return &BreakNode{}, true
	} else if p.checkCurToken(lexer.TOKEN_CONTINUE) {
		p.nextToken()
		return &ContinueNode{}, true
	} else {
		expr, ok := p.parseExpression(PREC_LOWEST)
		if !ok {
			return nil, false
		}

		if p.checkCurToken(lexer.TOKEN_ASSIGN) {
			sym, ok := expr.(*SymbolNode)
			if !ok {
				p.recordError("cannot assign to non-symbol")
				return nil, false
			}
			p.nextToken()
			lhs, ok := p.parseExpression(PREC_LOWEST)
			if !ok {
				return nil, false
			}
			return &AssignNode{sym, lhs}, true
		} else {
			return &ExpressionStatement{expr}, true
		}
	}
}

func (p *Parser) parseLetStatement() (Statement, bool) {
	p.nextToken()
	if p.checkCurToken(lexer.TOKEN_SYMBOL) {
		dest := &SymbolNode{p.curToken.Value}
		p.nextToken()
		if !p.checkCurToken(lexer.TOKEN_ASSIGN) {
			p.recordError("expected = while parsing let statement")
			return nil, false
		}
		p.nextToken()
		v, ok := p.parseExpression(PREC_LOWEST)
		if !ok {
			return nil, false
		}
		return &LetNode{dest, v}, true
	} else {
		p.recordError("expected symbol while parsing let statement")
		return nil, false
	}
}

func (p *Parser) parseForStatement() (Statement, bool) {
	p.nextToken()
	if p.checkCurToken(lexer.TOKEN_SYMBOL) {
		sym := &SymbolNode{p.curToken.Value}
		p.nextToken()
		if !p.checkCurToken(lexer.TOKEN_IN) {
			p.recordError("expected in while parsing for loop")
			return nil, false
		}
		p.nextToken()
		iter, ok := p.parseExpression(PREC_LOWEST)
		if !ok {
			return nil, false
		}
		body, ok := p.parseBracedBlock()
		if !ok {
			return nil, false
		}
		return &ForNode{sym, iter, body}, true
	} else {
		p.recordError("expected symbol while parsing for loop")
		return nil, false
	}
}

func (p *Parser) parseWhileStatement() (Statement, bool) {
	p.nextToken()
	cond, ok := p.parseExpression(PREC_LOWEST)
	if !ok {
		return nil, false
	}
	body, ok := p.parseBracedBlock()
	if !ok {
		return nil, false
	}
	return &WhileNode{cond, body}, true
}

func (p *Parser) parseIfStatement() (Statement, bool) {
	p.nextToken()
	clauses := []*IfClause{}

	// Parse the if block.
	cond, ok := p.parseExpression(PREC_LOWEST)
	if !ok {
		return nil, false
	}
	body, ok := p.parseBracedBlock()
	if !ok {
		return nil, false
	}
	clauses = append(clauses, &IfClause{cond, body})

	// Parse zero or more elif blocks.
	for p.checkCurToken(lexer.TOKEN_ELIF) {
		p.nextToken()
		elifCond, ok := p.parseExpression(PREC_LOWEST)
		if !ok {
			return nil, false
		}
		elifBody, ok := p.parseBracedBlock()
		if !ok {
			return nil, false
		}
		clauses = append(clauses, &IfClause{elifCond, elifBody})
	}

	// Parse an optional else block.
	var elseBody *BlockNode = nil
	if p.checkCurToken(lexer.TOKEN_ELSE) {
		p.nextToken()
		elseBody, ok = p.parseBracedBlock()
		if !ok {
			return nil, false
		}
	}

	return &IfNode{clauses, elseBody}, true
}

func (p *Parser) parseReturnStatement() (Statement, bool) {
	p.nextToken()
	if p.checkCurToken(lexer.TOKEN_NEWLINE) || p.checkCurToken(lexer.TOKEN_EOF) {
		return &ReturnNode{nil}, true
	} else {
		expr, ok := p.parseExpression(PREC_LOWEST)
		if !ok {
			return nil, false
		}

		return &ReturnNode{expr}, true
	}
}

func (p *Parser) parseFnStatement() (Statement, bool) {
	p.nextToken()
	if !p.checkCurToken(lexer.TOKEN_SYMBOL) {
		p.recordError("expected symbol while parsing function declaration")
		return nil, false
	}
	sym := &SymbolNode{p.curToken.Value}

	p.nextToken()
	if !p.checkCurToken(lexer.TOKEN_LPAREN) {
		p.recordError("expected ( while parsing function declaration")
		return nil, false
	}
	p.nextToken()
	params, ok := p.parseParamList()
	if !ok {
		return nil, false
	}

	body, ok := p.parseBracedBlock()
	if !ok {
		return nil, false
	}
	return &FnNode{sym, params, body}, true
}

func (p *Parser) parseExpression(precedence int) (Expression, bool) {
	left, ok := p.parsePrefix()
	if !ok {
		return nil, false
	}

	for {
		// Keep consuming infix operators until we hit either a non-infix token or an
		// infix operator with a lower precedence.
		if infixPrecedence, ok := precedenceMap[p.curToken.Type]; ok {
			if precedence < infixPrecedence {
				if p.curToken.Type == lexer.TOKEN_LPAREN {
					p.nextToken()
					arglist, ok := p.parseArglist(lexer.TOKEN_RPAREN)
					if !ok {
						return nil, false
					}
					left = &CallNode{left, arglist}
				} else {
					left, ok = p.parseInfix(left, getPrecedence(p.curToken.Type))
					if !ok {
						return nil, false
					}
				}
			} else {
				break
			}
		} else {
			break
		}
	}

	return left, true
}

func (p *Parser) parsePrefix() (Expression, bool) {
	typ := p.curToken.Type
	val := p.curToken.Value
	p.nextToken()
	if typ == lexer.TOKEN_INT {
		v, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			p.recordError("could not parse integer token")
			return nil, false
		}
		return &IntegerNode{int(v)}, true
	} else if typ == lexer.TOKEN_STRING {
		return &StringNode{val}, true
	} else if typ == lexer.TOKEN_SYMBOL {
		return &SymbolNode{val}, true
	} else if typ == lexer.TOKEN_TRUE {
		return &BoolNode{true}, true
	} else if typ == lexer.TOKEN_FALSE {
		return &BoolNode{false}, true
	} else if typ == lexer.TOKEN_LPAREN {
		expr, ok := p.parseExpression(PREC_LOWEST)
		if !ok {
			return nil, false
		}

		if !p.checkCurToken(lexer.TOKEN_RPAREN) {
			p.recordError("expected )")
			return nil, false
		}

		p.nextToken()
		return expr, true
	} else if typ == lexer.TOKEN_MINUS {
		expr, ok := p.parseExpression(PREC_PREFIX)
		return &PrefixNode{val, expr}, ok
	} else if typ == lexer.TOKEN_LBRACKET {
		values, ok := p.parseArglist(lexer.TOKEN_RBRACKET)
		return &ListNode{values}, ok
	} else {
		p.recordError(fmt.Sprintf("unexpected token %s", p.curToken.Type))
		return nil, false
	}
}

func (p *Parser) parseInfix(left Expression, precedence int) (Expression, bool) {
	operator := p.curToken.Value
	p.nextToken()
	right, ok := p.parseExpression(precedence)
	if !ok {
		return nil, false
	}
	return &InfixNode{operator, left, right}, true
}

func (p *Parser) parseArglist(terminator string) ([]Expression, bool) {
	arglist := []Expression{}
	// Special case for empty arglist
	if p.checkCurToken(terminator) {
		p.nextToken()
		return arglist, true
	}

	for {
		expr, ok := p.parseExpression(PREC_LOWEST)
		if !ok {
			return nil, false
		}
		arglist = append(arglist, expr)

		if p.checkCurToken(lexer.TOKEN_COMMA) {
			p.nextToken()
		} else if p.checkCurToken(terminator) {
			p.nextToken()
			break
		} else {
			p.recordError(fmt.Sprintf("unexpected token %s while parsing argument list",
				p.curToken.Type))
			return nil, false
		}
	}
	return arglist, true
}

func (p *Parser) parseParamList() ([]*SymbolNode, bool) {
	paramlist := []*SymbolNode{}
	// Special case for empty list
	if p.checkCurToken(lexer.TOKEN_RPAREN) {
		p.nextToken()
		return paramlist, true
	}

	for {
		if !p.checkCurToken(lexer.TOKEN_SYMBOL) {
			p.recordError("expected symbol while parsing parameter list")
			return nil, false
		}
		paramlist = append(paramlist, &SymbolNode{p.curToken.Value})

		p.nextToken()
		if p.checkCurToken(lexer.TOKEN_COMMA) {
			p.nextToken()
		} else if p.checkCurToken(lexer.TOKEN_RPAREN) {
			p.nextToken()
			break
		} else {
			p.recordError(fmt.Sprintf("unexpected token %s while parsing parameter list",
				p.curToken.Type))
			return nil, false
		}
	}
	return paramlist, true
}

func (p *Parser) parseBracedBlock() (*BlockNode, bool) {
	if !p.checkCurToken(lexer.TOKEN_LBRACE) {
		p.recordError("expected { while parsing block")
		return nil, false
	}
	p.nextToken()

	p.skipNewlines()

	if p.checkCurToken(lexer.TOKEN_RBRACE) {
		p.nextToken()
		return &BlockNode{[]Statement{}}, true
	}

	block, ok := p.parseBlock(false)
	if !ok {
		return nil, false
	}

	if !p.checkCurToken(lexer.TOKEN_RBRACE) {
		p.recordError("expected } while parsing block")
		return nil, false
	}
	p.nextToken()

	return block, true
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
