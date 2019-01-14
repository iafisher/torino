package lexer

const (
	// Keywords
	TOKEN_FN     = "TOKEN_FN"
	TOKEN_LET    = "TOKEN_LET"
	TOKEN_RETURN = "TOKEN_RETURN"

	// Operators
	TOKEN_PLUS         = "TOKEN_PLUS"
	TOKEN_MINUS        = "TOKEN_MINUS"
	TOKEN_ASTERISK     = "TOKEN_ASTERISK"
	TOKEN_SLASH        = "TOKEN_SLASH"
	TOKEN_DOUBLE_SLASH = "TOKEN_DOUBLE_SLASH"
	TOKEN_EQ           = "TOKEN_EQ"

	// Value literals
	TOKEN_SYMBOL = "TOKEN_SYMBOL"
	TOKEN_INT    = "TOKEN_INT"
	TOKEN_STRING = "TOKEN_STRING"

	TOKEN_LPAREN = "TOKEN_LPAREN"
	TOKEN_RPAREN = "TOKEN_RPAREN"
	TOKEN_LBRACE = "TOKEN_LBRACE"
	TOKEN_RBRACE = "TOKEN_RBRACE"

	TOKEN_NEWLINE = "TOKEN_NEWLINE"
	TOKEN_EOF     = "TOKEN_EOF"
	TOKEN_UNKNOWN = "TOKEN_UNKNOWN"
)

type Token struct {
	Type  string
	Value string
	Loc   *Location
}

type Location struct {
	Line   int
	Column int
}
