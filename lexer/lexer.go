package lexer

import (
	"fmt"
	"strconv"
	"strings"
)

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
			valueStr := string(runes[start:pos])
			valueInt, err := strconv.Atoi(valueStr)
			if err != nil {
				return nil, fmt.Errorf("failed to convert string to int: %s", valueStr)
			}
			cur.Next = &Token{Kind: NUM, Str: valueStr, Val: valueInt}
			cur = cur.Next
			continue
		}

		// if it's a symbol, create a RESERVED token
		if isSymbol(ch) {
			cur.Next = &Token{Kind: RESERVED, Str: string(ch)}
			cur = cur.Next
			pos++
			continue
		}

		// if it's an unknown character, return an error
		return nil, fmt.Errorf("unknown character: %c", ch)
	}

	cur.Next = &Token{Kind: EOF, Str: ""}
	return head.Next, nil
}
