package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/karamaru-alpha/monkey/token"
)

func TestLexer_NextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);

if (5 < 10) {
	return true;
} else {
	return false;
}

10 == 10;
10 != 9;

!-/*
[1, 2, "hoge", {"key": "val"}];
`
	expected := []token.Token{
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "five"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "5"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "ten"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.INT, Literal: "10"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "add"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.FUNCTION, Literal: "fn"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENT, Literal: "x"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.IDENT, Literal: "x"},
		{Type: token.PLUS, Literal: "+"},
		{Type: token.IDENT, Literal: "y"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.LET, Literal: "let"},
		{Type: token.IDENT, Literal: "result"},
		{Type: token.ASSIGN, Literal: "="},
		{Type: token.IDENT, Literal: "add"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.IDENT, Literal: "five"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.IDENT, Literal: "ten"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.IF, Literal: "if"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.INT, Literal: "5"},
		{Type: token.LT, Literal: "<"},
		{Type: token.INT, Literal: "10"},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.RETURN, Literal: "return"},
		{Type: token.TRUE, Literal: "true"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.ELSE, Literal: "else"},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.RETURN, Literal: "return"},
		{Type: token.FALSE, Literal: "false"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.INT, Literal: "10"},
		{Type: token.EQ, Literal: "=="},
		{Type: token.INT, Literal: "10"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.INT, Literal: "10"},
		{Type: token.NOT_EQ, Literal: "!="},
		{Type: token.INT, Literal: "9"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.BANG, Literal: "!"},
		{Type: token.MINUS, Literal: "-"},
		{Type: token.SLASH, Literal: "/"},
		{Type: token.ASTERISK, Literal: "*"},
		{Type: token.LBRACKET, Literal: "["},
		{Type: token.INT, Literal: "1"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.INT, Literal: "2"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.STRING, Literal: "hoge"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.LBRACE, Literal: "{"},
		{Type: token.STRING, Literal: "key"},
		{Type: token.COLON, Literal: ":"},
		{Type: token.STRING, Literal: "val"},
		{Type: token.RBRACE, Literal: "}"},
		{Type: token.RBRACKET, Literal: "]"},
		{Type: token.SEMICOLON, Literal: ";"},
		{Type: token.EOF, Literal: ""},
	}

	l := New(input)
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], l.NextToken())
	}
}
