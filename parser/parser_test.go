package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/karamaru-alpha/monkey/ast"
	"github.com/karamaru-alpha/monkey/lexer"
	"github.com/karamaru-alpha/monkey/token"
)

func TestParser_LetStatement(t *testing.T) {
	input := `let x = 5;
let y = 10;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseError(t, p)

	var tests = []*ast.LetStatement{
		{
			Token: token.Token{Type: token.LET, Literal: "let"},
			Name: &ast.Identifier{
				Token: token.Token{Type: token.IDENT, Literal: "x"},
				Value: "x",
			},
			Value: nil,
		},
		{
			Token: token.Token{Type: token.LET, Literal: "let"},
			Name: &ast.Identifier{
				Token: token.Token{Type: token.IDENT, Literal: "y"},
				Value: "y",
			},
			Value: nil,
		},
	}

	for i, tt := range tests {
		assert.Equal(t, tt, program.Statements[i].(*ast.LetStatement))
	}
}

func TestParser_ReturnStatement(t *testing.T) {
	input := `return 5;
return add(1, 2);
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseError(t, p)

	var tests = []*ast.ReturnStatement{
		{
			Token:       token.Token{Type: token.RETURN, Literal: "return"},
			ReturnValue: nil,
		},
		{
			Token:       token.Token{Type: token.RETURN, Literal: "return"},
			ReturnValue: nil,
		},
	}

	for i, tt := range tests {
		assert.Equal(t, tt, program.Statements[i].(*ast.ReturnStatement))
	}
}

func TestParser_IdentifierExpression(t *testing.T) {
	//input := "foobar;"
	//
	//l := lexer.New(input)
	//p := New(l)
	//program := p.ParseProgram()
	//checkParseError(t, p)
	//
	//var tests = []*ast.ReturnStatement{
	//	{
	//		Token:       token.Token{Type: token.RETURN, Literal: "return"},
	//		ReturnValue: nil,
	//	},
	//	{
	//		Token:       token.Token{Type: token.RETURN, Literal: "return"},
	//		ReturnValue: nil,
	//	},
	//}
	//
	//for i, tt := range tests {
	//	assert.Equal(t, tt, program.Statements[i].(*ast.ReturnStatement))
	//}
}

func checkParseError(t *testing.T, p *Parser) {
	for _, err := range p.errors {
		t.Error(err)
	}
}
