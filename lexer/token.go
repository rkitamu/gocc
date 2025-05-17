package lexer

type TokenKind int

const (
	RESERVED TokenKind = iota
	RETURN
	IDENT
	NUM
	EOF
)

type Token struct {
	Kind TokenKind
	Next *Token
	Str  string
	Val  int
	Pos  int // Position in the input string
}
