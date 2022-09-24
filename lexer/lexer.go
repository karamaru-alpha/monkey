package lexer

import (
	"github.com/karamaru-alpha/monkey/token"
)

type Lexer struct {
	input        string
	position     int  // 入力における現在の位置
	ch           byte // 現在検査中の文字
	readPosition int  // 次読み込む位置(position+1)
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
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			return token.Token{Type: token.EQ, Literal: "=="}
		} else {
			return token.New(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			return token.Token{Type: token.NOT_EQ, Literal: "!="}
		} else {
			return token.New(token.BANG, l.ch)
		}
	case ';':
		return token.New(token.SEMICOLON, l.ch)
	case '(':
		return token.New(token.LPAREN, l.ch)
	case ')':
		return token.New(token.RPAREN, l.ch)
	case '[':
		return token.New(token.LBRACKET, l.ch)
	case ']':
		return token.New(token.RBRACKET, l.ch)
	case ',':
		return token.New(token.COMMA, l.ch)
	case '+':
		return token.New(token.PLUS, l.ch)
	case '-':
		return token.New(token.MINUS, l.ch)
	case '*':
		return token.New(token.ASTERISK, l.ch)
	case '/':
		return token.New(token.SLASH, l.ch)
	case '<':
		return token.New(token.LT, l.ch)
	case '>':
		return token.New(token.GT, l.ch)
	case '{':
		return token.New(token.LBRACE, l.ch)
	case '}':
		return token.New(token.RBRACE, l.ch)
	case '"':
		var tok token.Token
		tok.Type = token.STRING
		tok.Literal = l.readString()
		return tok
	case 0:
		return token.Token{Type: token.EOF, Literal: ""}
	default:
		var tok token.Token
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		}
		if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		}
		return token.New(token.ILLEGAL, l.ch)
	}
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

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.peekChar()) {
		l.readChar()
	}
	return l.input[position:l.readPosition]
}

func isLetter(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}
