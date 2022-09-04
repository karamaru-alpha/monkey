package parser

import (
	"fmt"
	"testing"

	"github.com/karamaru-alpha/monkey/ast"
	"github.com/karamaru-alpha/monkey/lexer"
)

func TestParser_LetStatement(t *testing.T) {
	input := `
let x = 5;
let y = 10;
let foobar = 838383;
`
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParseError(t, p)

	fmt.Println(program.Statements[0].(*ast.LetStatement).Name.TokenLiteral())
}

func checkParseError(t *testing.T, p *Parser) {
	for _, err := range p.Errors() {
		t.Error(err)
	}
}
