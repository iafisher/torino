package lexer

const (
	// Keywords
	TOKEN_FN     = "TOKEN_FN"
	TOKEN_FOR    = "TOKEN_FOR"
	TOKEN_IF     = "TOKEN_IF"
	TOKEN_LET    = "TOKEN_LET"
	TOKEN_RETURN = "TOKEN_RETURN"
	TOKEN_WHILE  = "TOKEN_WHILE"

	// Operators
	TOKEN_PLUS         = "TOKEN_PLUS"
	TOKEN_MINUS        = "TOKEN_MINUS"
	TOKEN_ASTERISK     = "TOKEN_ASTERISK"
	TOKEN_SLASH        = "TOKEN_SLASH"
	TOKEN_DOUBLE_SLASH = "TOKEN_DOUBLE_SLASH"
	TOKEN_EQ           = "TOKEN_EQ"
	TOKEN_GT           = "TOKEN_GT"
	TOKEN_LT           = "TOKEN_LT"
	TOKEN_GE           = "TOKEN_GE"
	TOKEN_LE           = "TOKEN_LE"
	TOKEN_AND          = "TOKEN_AND"
	TOKEN_OR           = "TOKEN_OR"
	TOKEN_IN           = "TOKEN_IN"

	// Value literals
	TOKEN_SYMBOL = "TOKEN_SYMBOL"
	TOKEN_INT    = "TOKEN_INT"
	TOKEN_STRING = "TOKEN_STRING"
	TOKEN_TRUE   = "TOKEN_TRUE"
	TOKEN_FALSE  = "TOKEN_FALSE"

	TOKEN_ASSIGN = "TOKEN_ASSIGN"
	TOKEN_COMMA  = "TOKEN_COMMA"

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
