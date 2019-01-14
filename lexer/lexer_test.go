package lexer

import "testing"

func TestNextToken(t *testing.T) {
	input := `
fn f(y) {
	return y / 3 + 7 * -2 // 4
}

let x = f(10)
let s = "\n\c\\\""`
	tests := []struct {
		expectedType  string
		expectedValue string
	}{
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_FN, "fn"},
		{TOKEN_SYMBOL, "f"},
		{TOKEN_LPAREN, "("},
		{TOKEN_SYMBOL, "y"},
		{TOKEN_RPAREN, ")"},
		{TOKEN_LBRACE, "{"},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_RETURN, "return"},
		{TOKEN_SYMBOL, "y"},
		{TOKEN_SLASH, "/"},
		{TOKEN_INT, "3"},
		{TOKEN_PLUS, "+"},
		{TOKEN_INT, "7"},
		{TOKEN_ASTERISK, "*"},
		{TOKEN_MINUS, "-"},
		{TOKEN_INT, "2"},
		{TOKEN_DOUBLE_SLASH, "//"},
		{TOKEN_INT, "4"},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_RBRACE, "}"},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_LET, "let"},
		{TOKEN_SYMBOL, "x"},
		{TOKEN_EQ, "="},
		{TOKEN_SYMBOL, "f"},
		{TOKEN_LPAREN, "("},
		{TOKEN_INT, "10"},
		{TOKEN_RPAREN, ")"},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_LET, "let"},
		{TOKEN_SYMBOL, "s"},
		{TOKEN_EQ, "="},
		{TOKEN_STRING, "\n\\c\\\""},
		{TOKEN_EOF, ""},
	}

	l := New(input)
	for _, tt := range tests {
		got := l.NextToken()
		if got.Type != tt.expectedType {
			t.Fatalf("Wrong token type: got %q, expected %q",
				got.Type, tt.expectedType)
		}

		if got.Value != tt.expectedValue {
			t.Fatalf("Wrong token value: got %q, expected %q (type %q)",
				got.Value, tt.expectedValue, got.Type)
		}
	}
}

func TestUnclosedStringLiterals(t *testing.T) {
	tests := []string{
		`"`,
		`"\`,
	}

	for _, tt := range tests {
		l := New(tt)
		first := l.NextToken()
		if first.Type != TOKEN_UNKNOWN {
			t.Fatalf("Expected unknown token, got %q", first.Type)
		}

		second := l.NextToken()
		if second.Type != TOKEN_EOF {
			t.Fatalf("Expected EOF token, got %q", second.Type)
		}
	}
}
