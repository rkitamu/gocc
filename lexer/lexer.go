package lexer

import (
	"fmt"
	"strconv"
	"strings"

	"rkitamu/gocc/errors"
)

type Lexer struct {
	input string
}

func NewLexer(input string) *Lexer {
	return &Lexer{input: input}
}

func isSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func isSymbol(ch rune) bool {
	return strings.ContainsRune("+-*/=()<>", ch)
}

// Lex takes an input string and returns a linked list of tokens.
func (p *Lexer) Lex() (*Token, error) {
	runes := []rune(p.input)
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

		// if it's a letter, create an IDENT token
		if 'a' <= ch && ch <= 'z' {
			cur.Next = &Token{Kind: IDENT, Str: string(ch), Pos: pos}
			cur = cur.Next
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
				return nil, errors.NewPosError(
					fmt.Sprintf("invalid numeric literal: %s", valueStr),
					p.input,
					start,
				)
			}
			cur.Next = &Token{Kind: NUM, Str: valueStr, Val: valueInt, Pos: start}
			cur = cur.Next
			continue
		}

		// if it's a symbol, check for multi-character operators
		if pos+1 < len(runes) {
			two := string(runes[pos : pos+2])
			switch two {
			case "==", "!=", "<=", ">=":
				cur.Next = &Token{Kind: RESERVED, Str: two, Pos: pos}
				cur = cur.Next
				pos += 2
				continue
			}
		}

		// if it's a symbol, create a RESERVED token
		if isSymbol(ch) {
			cur.Next = &Token{Kind: RESERVED, Str: string(ch), Pos: pos}
			cur = cur.Next
			pos++
			continue
		}

		// if it's an unknown character, return an error
		return nil, errors.NewPosError(
			fmt.Sprintf("unexpected character: %c", ch),
			p.input,
			pos,
		)
	}

	cur.Next = &Token{Kind: EOF, Str: "EOF", Pos: pos}
	return head.Next, nil
}
