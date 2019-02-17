/* Grammar of Torino.

start := block

block := (stmt NEWLINE)*
stmt  := let | fn | for | while | if | break | return | expr

let    := LET SYMBOL ASSIGN expr
fn     := FN SYMBOL LPAREN params? RPAREN brace-block
for    := FOR SYMBOL IN expr brace-block
while  := WHILE expr brace-block
if     := IF expr brace-block elif* else?
elif   := ELIF expr brace-block
else   := ELSE brace-block
break  := BREAK
return := RETURN expr?

brace-block := LBRACE NEWLINE block RBRACE

expr  := infix | call | pexpr | list | map
       | INT | STRING | SYMBOL | TRUE | FALSE
pexpr := LPAREN expr RPAREN
infix := expr OP expr
call  := SYMBOL LPAREN args? RPAREN
list  := LBRACKET args? RBRACKET
map   := LBRACKET mapargs? RBRACKET

params  := (SYMBOL COMMA)* SYMBOL
args    := (expr COMMA)* expr
mapargs := (maparg COMMA)* maparg
maparg  := expr COLON expr
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
	return p.parseBlock()
}

func (p *Parser) parseBlock() *BlockNode {
	statements := []Node{}
	for {
		stmt := p.parseStatement()
		statements = append(statements, stmt)
		if p.checkCurToken(lexer.TOKEN_NEWLINE) {
			for p.checkCurToken(lexer.TOKEN_NEWLINE) {
				p.nextToken()
			}
		} else if p.checkCurToken(lexer.TOKEN_EOF) {
			break
		} else {
			panic(fmt.Sprintf("parseBlock - unexpected token %s", p.curToken.Type))
		}
	}
	return &BlockNode{statements}
}

func (p *Parser) parseStatement() Statement {
	if p.checkCurToken(lexer.TOKEN_LET) {
		return p.parseLetStatement()
	} else {
		return &ExpressionStatement{p.parseExpression(PREC_LOWEST)}
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

func (p *Parser) parseExpression(precedence int) Expression {
	left := p.parsePrefix()
	p.nextToken()

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
	if p.checkCurToken(lexer.TOKEN_INT) {
		v, err := strconv.ParseInt(p.curToken.Value, 10, 64)
		if err != nil {
			panic("parseExpression - could not parse integer token")
		}
		return &IntegerNode{v}
	} else if p.checkCurToken(lexer.TOKEN_STRING) {
		return &StringNode{p.curToken.Value}
	} else if p.checkCurToken(lexer.TOKEN_SYMBOL) {
		return &SymbolNode{p.curToken.Value}
	} else if p.checkCurToken(lexer.TOKEN_TRUE) {
		return &BoolNode{true}
	} else if p.checkCurToken(lexer.TOKEN_FALSE) {
		return &BoolNode{false}
	} else if p.checkCurToken(lexer.TOKEN_LPAREN) {
		p.nextToken()
		expr := p.parseExpression(PREC_LOWEST)
		if !p.checkCurToken(lexer.TOKEN_RPAREN) {
			panic("parseExpression - expected )")
		}
		return expr
	} else {
		panic(fmt.Sprintf("parseExpression - unexpected token %s", p.curToken.Type))
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

func (p *Parser) checkCurToken(expectedType string) bool {
	return p.curToken.Type == expectedType
}

func (p *Parser) nextToken() *lexer.Token {
	p.curToken = p.lexer.NextToken()
	return p.curToken
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
	PREC_ADD_SUB
	PREC_MUL_DIV
	PREC_PREFIX
	PREC_CALL
)

var precedenceMap = map[string]int{
	lexer.TOKEN_PLUS:     PREC_ADD_SUB,
	lexer.TOKEN_MINUS:    PREC_ADD_SUB,
	lexer.TOKEN_ASTERISK: PREC_MUL_DIV,
	lexer.TOKEN_SLASH:    PREC_MUL_DIV,
	lexer.TOKEN_LPAREN:   PREC_CALL,
}
