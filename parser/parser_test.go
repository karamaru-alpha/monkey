package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/karamaru-alpha/monkey/ast"
	"github.com/karamaru-alpha/monkey/lexer"
	"github.com/karamaru-alpha/monkey/token"
)

func TestParser_LetStatement(t *testing.T) {
	input := "let x = 5;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseError(t, p)

	tests := []*ast.LetStatement{
		{
			Token: token.Token{Type: token.LET, Literal: "let"},
			Name: &ast.Identifier{
				Token: token.Token{Type: token.IDENT, Literal: "x"},
				Value: "x",
			},
			Value: &ast.IntegerLiteral{
				Token: token.Token{Type: token.INT, Literal: "5"},
				Value: 5,
			},
		},
	}

	for i, tt := range tests {
		fmt.Println()
		assert.Equal(t, tt.Value, program.Statements[i].(*ast.LetStatement).Value)
	}
}

func TestParser_ReturnStatement(t *testing.T) {
	input := "return 5"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseError(t, p)

	var tests = []*ast.ReturnStatement{
		{
			Token: token.Token{Type: token.RETURN, Literal: "return"},
			ReturnValue: &ast.IntegerLiteral{
				Token: token.Token{Type: token.INT, Literal: "5"},
				Value: 5,
			},
		},
	}

	for i, tt := range tests {
		assert.Equal(t, tt, program.Statements[i].(*ast.ReturnStatement))
	}
}

func TestOperator_PrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "-a * b",
			expected: "((-a) * b)",
		},
		{
			input:    "1 < 2 == 4 > 3",
			expected: "((1 < 2) == (4 > 3))",
		},
		{
			input:    "1 / 2; 3 * 4",
			expected: "(1 / 2)(3 * 4)",
		},
		{
			input:    "1 + 2 * 3",
			expected: "(1 + (2 * 3))",
		},
		{
			input:    "false != true",
			expected: "(false != true)",
		},
		{
			input:    "1 + 2 + 3",
			expected: "((1 + 2) + 3)",
		},
		{
			input:    "1 / (2 + 3)",
			expected: "(1 / (2 + 3))",
		},
		{
			input:    "1 + (2 + 3) + 4",
			expected: "((1 + (2 + 3)) + 4)",
		},
		{
			input:    "add(1 + 2, 3)",
			expected: "add((1 + 2), 3)",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParseError(t, p)

		actual := program.String()
		assert.Equal(t, tt.expected, actual)

	}
}

func checkParseError(t *testing.T, p *Parser) {
	for _, err := range p.errors {
		t.Error(err)
	}
}
