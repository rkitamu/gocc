package lexer

type TokenKind int

const (
	RESERVED TokenKind = iota
	NUM
	EOF
)

type Token struct {
	Kind TokenKind
	Next *Token
	Str  string
	Val  int
}
