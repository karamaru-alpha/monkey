package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/karamaru-alpha/monkey/token"
)

func TestLexer_NextToken(t *testing.T) {
	input := "+(),;"

	expected := []token.Token{
		{Type: token.PLUS, Literal: "+"},
		{Type: token.LPAREN, Literal: "("},
		{Type: token.RPAREN, Literal: ")"},
		{Type: token.COMMA, Literal: ","},
		{Type: token.SEMICOLUN, Literal: ";"},
	}

	l := New(input)
	res := make([]token.Token, 0, len(input))
	for i := 0; i < len(input); i++ {
		res = append(res, l.NextToken())
	}

	assert.Equal(t, expected, res)
}
