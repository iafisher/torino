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

expr  := infix | call | pexpr | list | map | INT | STRING | SYMBOL
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
	"github.com/iafisher/torino/lexer"
	"strconv"
)

type Parser struct {
	lexer *lexer.Lexer
}

func New(l *lexer.Lexer) *Parser {
	return &Parser{l}
}

func (p *Parser) Parse() *BlockNode {
	return p.parseBlock()
}

func (p *Parser) parseBlock() *BlockNode {
	statements := []Node{}
	for {
		stmt := p.parseStatement()
		statements = append(statements, stmt)
		tkn := p.lexer.NextToken()
		if tkn.Type == lexer.TOKEN_NEWLINE {
			continue
		} else if tkn.Type == lexer.TOKEN_EOF {
			break
		} else {
			panic("parseBlock - unexpected token")
		}
	}
	return &BlockNode{statements}
}

func (p *Parser) parseStatement() Statement {
	return p.parseExpression()
}

func (p *Parser) parseExpression() Expression {
	tkn := p.lexer.NextToken()
	if tkn.Type == lexer.TOKEN_INT {
		v, err := strconv.ParseInt(tkn.Value, 10, 64)
		if err != nil {
			panic("parseExpression - could not parse integer token")
		}
		return &IntegerNode{v}
	} else if tkn.Type == lexer.TOKEN_STRING {
		return &StringNode{tkn.Value}
	} else {
		panic("parseExpression - unexpected token")
	}
}
