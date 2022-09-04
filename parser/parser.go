package parser

import (
	"fmt"

	"github.com/karamaru-alpha/monkey/ast"
	"github.com/karamaru-alpha/monkey/lexer"
	"github.com/karamaru-alpha/monkey/token"
)

type Parser struct {
	lex *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token
	errors       []error
}

func New(lex *lexer.Lexer) *Parser {
	p := &Parser{lex: lex}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	var stmt = &ast.LetStatement{Token: p.currentToken}

	if p.peekToken.Type != token.IDENT {
		p.errors = append(p.errors, fmt.Errorf("wrong token. expected: %s, actual: %s", token.IDENT, p.peekToken.Type))
		return nil
	}
	p.nextToken()
	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if p.peekToken.Type != token.ASSIGN {
		p.errors = append(p.errors, fmt.Errorf("wrong token. expected: %s, actual: %s", token.ASSIGN, p.peekToken.Type))
		return nil
	}
	p.nextToken()

	for p.currentToken.Type != token.SEMICOLUN {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	var stmt = &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()
	for p.currentToken.Type != token.SEMICOLUN {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) Errors() []error {
	return p.errors
}
