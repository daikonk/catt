package lexer

import (
	"go_interpreter/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
  while (false) { 10 % 10; }
10 != 9;
  fn
"foobar"
"foo bar"
  for (var i = 0; i < 10; var i = i + 1) { 10 / 10 };
  if (true && false) { 10 / 10 }
`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.WHILE, "while"},
		{token.LPAR, "("},
		{token.FALSE, "false"},
		{token.RPAR, ")"},
		{token.LBRAC, "{"},
		{token.INT, "10"},
		{token.MODULO, "%"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.RBRAC, "}"},
		{token.INT, "10"},
		{token.NOT_EQ, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.FUNCTION, "fn"},
		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},
		{token.FOR, "for"},
		{token.LPAR, "("},
		{token.VAR, "var"},
		{token.IDENT, "i"},
		{token.ASSIGN, "="},
		{token.INT, "0"},
		{token.SEMICOLON, ";"},
		{token.IDENT, "i"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"},
		{token.IDENT, "i"},
		{token.ASSIGN, "="},
		{token.IDENT, "i"},
		{token.PLUS, "+"},
		{token.INT, "1"},
		{token.RPAR, ")"},
		{token.LBRAC, "{"},
		{token.INT, "10"},
		{token.SLASH, "/"},
		{token.INT, "10"},
		{token.RBRAC, "}"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAR, "("},
		{token.TRUE, "true"},
		{token.AND, "&&"},
		{token.FALSE, "false"},
		{token.RPAR, ")"},
		{token.LBRAC, "{"},
		{token.INT, "10"},
		{token.SLASH, "/"},
		{token.INT, "10"},
		{token.RBRAC, "}"},
		{token.EOF, ""},
	}
	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
