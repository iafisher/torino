package lexer

import "strings"

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
	"fn":     TOKEN_FN,
	"let":    TOKEN_LET,
	"return": TOKEN_RETURN,
}

func (l *Lexer) NextToken() *Token {
	if l.position >= len(l.program) {
		return l.makeToken(TOKEN_EOF, "")
	}

	for isWhitespace(l.program[l.position]) {
		l.advance()
	}

	ch := l.program[l.position]

	// Single character tokens
	switch ch {
	case '+':
		return l.makeTokenAndAdvance(TOKEN_PLUS, "+")
	case '-':
		return l.makeTokenAndAdvance(TOKEN_MINUS, "-")
	case '*':
		return l.makeTokenAndAdvance(TOKEN_ASTERISK, "*")
	case '/':
		if l.position+1 < len(l.program) && l.program[l.position+1] == '/' {
			l.advance()
			l.advance()
			return l.makeToken(TOKEN_DOUBLE_SLASH, "//")
		} else {
			return l.makeTokenAndAdvance(TOKEN_SLASH, "/")
		}
	case '=':
		return l.makeTokenAndAdvance(TOKEN_EQ, "=")
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
		l.advance()
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

	for {
		if l.position >= len(l.program) {
			return str.String(), false
		}

		if l.program[l.position] == '"' {
			break
		} else if l.program[l.position] == '\\' {
			l.advance()
			if l.position >= len(l.program) {
				return str.String(), false
			}
			str.WriteString(decodeEscape(l.program[l.position]))
			l.advance()
		} else {
			str.WriteByte(l.program[l.position])
			l.advance()
		}
	}

	l.advance()

	return str.String(), true
}

func (l *Lexer) makeToken(typ string, value string) *Token {
	return &Token{typ, value, &Location{l.line, l.column}}
}

func (l *Lexer) makeTokenAndAdvance(typ string, value string) *Token {
	tok := l.makeToken(typ, value)
	l.advance()
	return tok
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
