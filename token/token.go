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

var keywords = map[string]Type{
	"fn":  FUNCTION,
	"let": LET,
}

func New(typ Type, ch byte) Token {
	return Token{Type: typ, Literal: string(ch)}
}

func LookupIdent(ident string) Type {
	if typ, ok := keywords[ident]; ok {
		return typ
	}
	return IDENT
}
