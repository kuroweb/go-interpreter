package lexer

import (
	"testing"
	"example.com/token"
)

func TestNestToken(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		exprectedType    token.TokenType
		exprectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "{"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}
}

l := New(input)

for i, tt := range tests {
	tok := l.NextToken()

	if tok.Type != tt.exprectedType {
		t.Fatalf("tests[%d] - tokentype wrong. exprected=%q, got=%q",
			i, tt.exprectedType, tok.Type)
	}

	if tok.Literal != tt.exprectedLiteral {
		t.Fatalf("tests[%d] - literal wrong. exprected=%q, got=%q")
			i. tt.exprectedLiteral, tok.Literal)
	}
}
