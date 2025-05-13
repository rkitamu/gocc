package lexer

import (
	"fmt"
	"strings"
)

type TokenKind int
const (
	RESERVED TokenKind = iota
	NUM
	EOF
)

type Token struct {
	kind TokenKind;
	str string
	next *Token;
}

func isSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isSymbol(ch rune) bool {
	return strings.ContainsRune("+-*/=()", ch)
}

// Lex takes an input string and returns a linked list of tokens.
func Lex(input string) (*Token, error) {
	runes := []rune(input)
	pos := 0

	head := &Token{}
	cur := head

	for pos < len(runes) {
		ch := runes[pos]

		// skip whitespace
		if isSpace(ch) {
			pos++
			continue
		}

		// if it's a digit, create a NUM token
		if isDigit(ch) {
			start := pos
			for pos < len(runes) && isDigit(runes[pos]) {
				pos++
			}
			value := string(runes[start:pos])
			cur.next = &Token{kind: NUM, str: value}
			cur = cur.next
			continue
		}

		// if it's a symbol, create a RESERVED token
		if isSymbol(ch) {
			cur.next = &Token{kind: RESERVED, str: string(ch)}
			cur = cur.next
			pos++
			continue
		}

		// if it's an unknown character, return an error
		return nil, fmt.Errorf("unknown character: %c", ch)
	}

	cur.next = &Token{kind: EOF, str: ""}
	return head.next, nil
}
