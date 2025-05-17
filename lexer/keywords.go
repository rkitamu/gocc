package lexer

// Keywords maps keyword strings to their corresponding TokenKind.
var Keywords = map[string]TokenKind{
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	//"while":  WHILE,
	//"for":    FOR,
}
