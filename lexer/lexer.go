package lexer

import (
	"fmt"
	"strings"
)

type Lexer struct {
	program  string
	position int
	line     int
	column   int
}

func New(program string) *Lexer {
	return &Lexer{program: program, position: 0, line: 1, column: 1}
}

var keywords = map[string]string{
	"and":    TOKEN_AND,
	"fn":     TOKEN_FN,
	"for":    TOKEN_FOR,
	"if":     TOKEN_IF,
	"in":     TOKEN_IN,
	"let":    TOKEN_LET,
	"or":     TOKEN_OR,
	"return": TOKEN_RETURN,
	"while":  TOKEN_WHILE,
}

func (l *Lexer) NextToken() *Token {
	err := l.skipWhitespaceAndComments()
	if err {
		fmt.Println("ouch!")
		return l.makeToken(TOKEN_UNKNOWN, "")
	}

	if l.position >= len(l.program) {
		return l.makeToken(TOKEN_EOF, "")
	}

	ch := l.program[l.position]

	// Single and double character tokens
	switch ch {
	case ',':
		return l.makeTokenAndAdvance(TOKEN_COMMA, ",")
	case '+':
		return l.makeTokenAndAdvance(TOKEN_PLUS, "+")
	case '-':
		return l.makeTokenAndAdvance(TOKEN_MINUS, "-")
	case '*':
		return l.makeTokenAndAdvance(TOKEN_ASTERISK, "*")
	case '/':
		if l.peek('/') {
			tok := l.makeToken(TOKEN_DOUBLE_SLASH, "//")
			l.advance()
			l.advance()
			return tok
		} else {
			return l.makeTokenAndAdvance(TOKEN_SLASH, "/")
		}
	case '=':
		if l.peek('=') {
			tok := l.makeToken(TOKEN_EQ, "==")
			l.advance()
			l.advance()
			return tok
		} else {
			return l.makeTokenAndAdvance(TOKEN_ASSIGN, "=")
		}
	case '<':
		if l.peek('=') {
			tok := l.makeToken(TOKEN_LE, "<=")
			l.advance()
			l.advance()
			return tok
		} else {
			return l.makeTokenAndAdvance(TOKEN_LT, "<")
		}
	case '>':
		if l.peek('=') {
			tok := l.makeToken(TOKEN_GE, ">=")
			l.advance()
			l.advance()
			return tok
		} else {
			return l.makeTokenAndAdvance(TOKEN_GT, ">")
		}
	case '(':
		return l.makeTokenAndAdvance(TOKEN_LPAREN, "(")
	case ')':
		return l.makeTokenAndAdvance(TOKEN_RPAREN, ")")
	case '{':
		return l.makeTokenAndAdvance(TOKEN_LBRACE, "{")
	case '}':
		return l.makeTokenAndAdvance(TOKEN_RBRACE, "}")
	case '\n':
		return l.makeTokenAndAdvance(TOKEN_NEWLINE, "\n")
	}

	// Multi character tokens
	switch {
	case ch == '"':
		value, ok := l.readString()
		if ok {
			return l.makeToken(TOKEN_STRING, value)
		} else {
			return l.makeToken(TOKEN_UNKNOWN, value)
		}
	case canStartIdentifier(ch):
		value := l.readIdentifier()
		keywordType, ok := keywords[value]
		if ok {
			return l.makeToken(keywordType, value)
		} else {
			return l.makeToken(TOKEN_SYMBOL, value)
		}
	case isDigit(ch):
		value := l.readInteger()
		return l.makeToken(TOKEN_INT, value)
	default:
		return l.makeTokenAndAdvance(TOKEN_UNKNOWN, string(ch))
	}
}

func (l *Lexer) advance() {
	if l.position < len(l.program) {
		if l.program[l.position] == '\n' {
			l.line += 1
			l.column = 1
		}
		l.position += 1
	}
}

func (l *Lexer) peek(ch byte) bool {
	return l.position+1 < len(l.program) && l.program[l.position+1] == ch
}

func (l *Lexer) skipWhitespaceAndComments() bool {
	for l.onCommentOrWhitespace() {
		if isWhitespace(l.program[l.position]) {
			for l.position < len(l.program) && isWhitespace(l.program[l.position]) {
				l.advance()
			}
		} else {
			err := l.readComment()
			if err {
				return true
			}
		}
	}
	return false
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for l.position < len(l.program) && isIdentifierChar(l.program[l.position]) {
		l.advance()
	}
	return l.program[start:l.position]
}

func (l *Lexer) readInteger() string {
	start := l.position
	for l.position < len(l.program) && isDigit(l.program[l.position]) {
		l.advance()
	}
	return l.program[start:l.position]
}

func (l *Lexer) readString() (string, bool) {
	var str strings.Builder

	// Skip the opening quote.
	l.advance()

	for {
		if l.position >= len(l.program) {
			return str.String(), false
		}

		ch := l.program[l.position]
		if ch == '"' {
			break
		} else if ch == '\\' {
			l.advance()
			if l.position >= len(l.program) {
				return str.String(), false
			}
			str.WriteString(decodeEscape(l.program[l.position]))
			l.advance()
		} else if ch == '\n' {
			// TODO: Better error for this?
			l.advance()
			return str.String(), false
		} else {
			str.WriteByte(l.program[l.position])
			l.advance()
		}
	}

	l.advance()

	return str.String(), true
}

func (l *Lexer) readComment() bool {
	// Skip the initial slash and asterisk.
	l.advance()
	l.advance()

	for l.position < len(l.program) && !strings.HasPrefix(l.program[l.position:], "*/") {
		l.advance()
	}

	if l.position == len(l.program) {
		return true
	} else {
		l.advance()
		l.advance()
		return false
	}
}

func (l *Lexer) makeToken(typ string, value string) *Token {
	return &Token{typ, value, &Location{l.line, l.column}}
}

func (l *Lexer) makeTokenAndAdvance(typ string, value string) *Token {
	tok := l.makeToken(typ, value)
	l.advance()
	return tok
}

func (l *Lexer) onCommentOrWhitespace() bool {
	return l.position < len(l.program) && isWhitespace(l.program[l.position]) ||
		strings.HasPrefix(l.program[l.position:], "/*")
}

func canStartIdentifier(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ch == '_'
}

func isIdentifierChar(ch byte) bool {
	return canStartIdentifier(ch) || isDigit(ch)
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isWhitespace(ch byte) bool {
	// Note that newline is not whitespace as it can be syntactically significant.
	return ch == ' ' || ch == '\t' || ch == '\v' || ch == '\f'
}

func decodeEscape(ch byte) string {
	// Same escape sequences as Go
	// (see https://golang.org/ref/spec, section "Rune literals")
	switch ch {
	case 'a':
		return "\a"
	case 'b':
		return "\b"
	case 'f':
		return "\f"
	case 'n':
		return "\n"
	case 'r':
		return "\r"
	case 't':
		return "\t"
	case 'v':
		return "\v"
	case '\\':
		return "\\"
	case '"':
		return "\""
	default:
		return "\\" + string(ch)
	}
}
