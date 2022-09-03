package lexer

import (
	"github.com/karamaru-alpha/monkey/token"
)

type Lexer struct {
	input        string
	position     int  // 入力における現在の位置
	readPosition int  // これから読み込む位置
	ch           byte // 現在検査中の文字
}

func New(input string) *Lexer {
	return &Lexer{input: input}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {
	l.readChar()
	l.skipWhiteSpace()
	var tok token.Token
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: "=="}
		} else {
			tok = token.New(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
		} else {
			tok = token.New(token.BANG, l.ch)
		}
	case ';':
		tok = token.New(token.SEMICOLUN, l.ch)
	case '(':
		tok = token.New(token.LPAREN, l.ch)
	case ')':
		tok = token.New(token.RPAREN, l.ch)
	case ',':
		tok = token.New(token.COMMA, l.ch)
	case '+':
		tok = token.New(token.PLUS, l.ch)
	case '-':
		tok = token.New(token.MINUS, l.ch)
	case '*':
		tok = token.New(token.ASTERISK, l.ch)
	case '/':
		tok = token.New(token.SLASH, l.ch)
	case '<':
		tok = token.New(token.LT, l.ch)
	case '>':
		tok = token.New(token.GT, l.ch)
	case '{':
		tok = token.New(token.LBRACE, l.ch)
	case '}':
		tok = token.New(token.RBRACE, l.ch)
	case 0:
		tok = token.Token{Type: token.EOF, Literal: ""}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
		} else {
			tok = token.New(token.INLLEGAL, l.ch)
		}
	}
	return tok
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.peekChar()) {
		l.readChar()
	}
	return l.input[position:l.readPosition]
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.peekChar()) {
		l.readChar()
	}
	return l.input[position:l.readPosition]
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
