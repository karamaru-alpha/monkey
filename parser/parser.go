package parser

import (
	"fmt"
	"strconv"

	"github.com/karamaru-alpha/monkey/ast"
	"github.com/karamaru-alpha/monkey/lexer"
	"github.com/karamaru-alpha/monkey/token"
)

type prefixParseFn func() ast.Expression

type infixParseFn func(left ast.Expression) ast.Expression

type Parser struct {
	lex    *lexer.Lexer
	errors []error

	currentToken token.Token
	peekToken    token.Token

	prefixParseFns map[token.Type]prefixParseFn
	infixParseFns  map[token.Type]infixParseFn
}

func New(lex *lexer.Lexer) *Parser {
	p := &Parser{lex: lex}

	p.prefixParseFns = make(map[token.Type]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	p.infixParseFns = make(map[token.Type]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) registerPrefix(typ token.Type, fn prefixParseFn) {
	p.prefixParseFns[typ] = fn
}

func (p *Parser) registerInfix(typ token.Type, fn infixParseFn) {
	p.infixParseFns[typ] = fn
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
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

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
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()
	for p.currentToken.Type != token.SEMICOLUN {
		p.nextToken()
	}
	return stmt
}

const (
	LOWEST      = iota + 1
	EQUALS      // =
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x or !x
	CALL        // func(x)
)

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	stmt.Expression = p.parseExpression(LOWEST)
	if p.peekToken.Type == token.SEMICOLUN {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Errorf("preFn is not registered. TokenType: %s", p.currentToken.Type))
		return nil
	}
	leftExp := prefix()

	for p.peekToken.Type != token.SEMICOLUN && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, err)
		return nil
	}
	lit.Value = value

	return lit
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currentToken, Value: p.currentToken.Type == token.TRUE}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if p.peekToken.Type != token.RPAREN {
		return nil
	}
	p.nextToken()
	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

var precedences = map[token.Type]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
}

func (p *Parser) currentPrecedence() int {
	if precedence, ok := precedences[p.currentToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() int {
	if precedence, ok := precedences[p.peekToken.Type]; ok {
		return precedence
	}
	return LOWEST
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}
	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.currentToken}

	if p.peekToken.Type != token.LPAREN {
		return nil
	}
	p.nextToken()

	p.nextToken()
	expression.Condition = p.parseExpression(LOWEST)

	if p.peekToken.Type != token.RPAREN {
		return nil
	}
	p.nextToken()

	if p.peekToken.Type != token.LBRACE {
		return nil
	}
	p.nextToken()

	expression.Consequence = p.parseBlockStatement()

	if p.peekToken.Type == token.ELSE {
		p.nextToken()

		if p.peekToken.Type != token.LBRACE {
			return nil
		}
		p.nextToken()

		expression.Alternative = p.parseBlockStatement()
	}

	return expression
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = make([]ast.Statement, 0)
	p.nextToken()

	for p.currentToken.Type != token.RBRACE && p.currentToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	return block
}
