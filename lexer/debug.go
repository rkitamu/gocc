package lexer

import (
	"fmt"
)

func (l *Lexer) DebugPrintTokens(tok *Token) {
	for i := 0; tok != nil; tok = tok.Next {
		fmt.Printf("[%d] %s", i, tokenKindToString(tok.Kind))
		if tok.Str != "" {
			fmt.Printf("(%q)", tok.Str)
		}
		if tok.Kind == NUM {
			fmt.Printf(" val=%d", tok.Val)
		}
		fmt.Println(" ->")
		i++
	}
	fmt.Println("END")
}

func tokenKindToString(kind TokenKind) string {
	switch kind {
	case RESERVED:
		return "RESERVED"
	case NUM:
		return "NUM"
	case IDENT:
		return "IDENT"
	case EOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}
