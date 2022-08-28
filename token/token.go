package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	// その他
	INLLEGAL = "ILLEGAL" // 未知のトークン
	EOF      = "EOF"     // ファイル終端

	// 識別子 リテラル
	IDENT = "IDENT"
	INT   = "INT"

	// 演算子
	ASSIGN = "="
	PLUS   = "+"

	// デリミタ
	COMMA     = ","
	SEMICOLUN = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// キーワード
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

func New(typ Type, ch byte) Token {
	return Token{Type: typ, Literal: string(ch)}
}
