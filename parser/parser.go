package parser

import (
	"fmt"
	"rkitamu/gocc/errors"
	"rkitamu/gocc/lexer"
)

type Parser struct {
	current *lexer.Token
	Code    []*Node
	input   string
}

func NewParser(token *lexer.Token, input string) *Parser {
	return &Parser{
		current: token,
		Code:    make([]*Node, 0),
		input:   input,
	}
}

// Parse parses the input tokens and returns the root node of the parse tree.
// supports the following grammar:
// program = stmt*
// stmt = expr ";"
// expr = assign
// assign = equality ("=" assign)?
// equality = relational ("==" relational | "!=" relational)*
// relational = add ("<" add | "<=" add | ">" add | ">=" add)*
// add = mul ("+" mul | "-" mul)*
// mul = unary ("*" unary | "/" unary)*
// unary = ("+" | "-")? unary | primary
// primary = num | ident | "(" expr ")"
func (p *Parser) Parse() error {
	return p.program()
}

// program = stmt*
func (p *Parser) program() error {
	for i := 0; !p.atEnd(); i++ {
		node, err := p.stmt()
		if err != nil {
			return err
		}
		p.Code = append(p.Code, node)
	}
	return nil
}

// stmt = expr ";"
func (p *Parser) stmt() (*Node, error) {
	node, err := p.expr()
	if err != nil {
		return nil, err
	}
	if err := p.expect(";"); err != nil {
		return nil, err
	}
	return node, nil
}

// expr = assign
func (p *Parser) expr() (*Node, error) {
	return p.assign()
}

// assign = equality ("=" assign)?
func (p *Parser) assign() (*Node, error) {
	node, err := p.equality()
	if err != nil {
		return nil, err
	}
	if p.match("=") {
		p.advance()
		rhs, err := p.assign()
		if err != nil {
			return nil, err
		}
		node = &Node{Kind: ASSIGN, Lhs: node, Rhs: rhs}
	}
	return node, nil
}

// equality = relational ("==" relational | "!=" relational)*
func (p *Parser) equality() (*Node, error) {
	node, err := p.relational()
	if err != nil {
		return nil, err
	}

	for {
		switch {
		case p.match("=="):
			p.advance()
			rhs, err := p.relational()
			if err != nil {
				return nil, err
			}
			node = &Node{Kind: EQ, Lhs: node, Rhs: rhs}
		case p.match("!="):
			p.advance()
			rhs, err := p.relational()
			if err != nil {
				return nil, err
			}
			node = &Node{Kind: NEQ, Lhs: node, Rhs: rhs}
		default:
			return node, nil
		}
	}
}

// relational = add ("<" add | "<=" add | ">" add | ">=" add)*
func (p *Parser) relational() (*Node, error) {
	node, err := p.add()
	if err != nil {
		return nil, err
	}
	for {
		switch {
		case p.match("<"):
			p.advance()
			rhs, err := p.add()
			if err != nil {
				return nil, err
			}
			node = &Node{Kind: LT, Lhs: node, Rhs: rhs}
		case p.match("<="):
			p.advance()
			rhs, err := p.add()
			if err != nil {
				return nil, err
			}
			node = &Node{Kind: LTE, Lhs: node, Rhs: rhs}
		case p.match(">"):
			p.advance()
			lhs, err := p.add()
			if err != nil {
				return nil, err
			}
			// ">" is equivalent to "<" in reverse
			node = &Node{Kind: LT, Lhs: lhs, Rhs: node}
		case p.match(">="):
			p.advance()
			lhs, err := p.add()
			if err != nil {
				return nil, err
			}
			// ">=" is equivalent to "<=" in reverse
			node = &Node{Kind: LTE, Lhs: lhs, Rhs: node}
		default:
			return node, nil
		}
	}
}

// add = mul ("+" mul | "-" mul)*
func (p *Parser) add() (*Node, error) {
	node, err := p.mul()
	if err != nil {
		return nil, err
	}
	for {
		switch {
		case p.match("+"):
			p.advance()
			rhs, err := p.mul()
			if err != nil {
				return nil, err
			}
			node = &Node{Kind: ADD, Lhs: node, Rhs: rhs}
		case p.match("-"):
			p.advance()
			rhs, err := p.mul()
			if err != nil {
				return nil, err
			}
			node = &Node{Kind: SUB, Lhs: node, Rhs: rhs}
		default:
			return node, nil
		}
	}
}

// mul = unary ("*" unary | "/" unary)*
func (p *Parser) mul() (*Node, error) {
	node, err := p.unary()
	if err != nil {
		return nil, err
	}
	for {
		switch {
		case p.match("*"):
			p.advance()
			rhs, err := p.unary()
			if err != nil {
				return nil, err
			}
			node = &Node{Kind: MUL, Lhs: node, Rhs: rhs}
		case p.match("/"):
			p.advance()
			rhs, err := p.unary()
			if err != nil {
				return nil, err
			}
			node = &Node{Kind: DIV, Lhs: node, Rhs: rhs}
		default:
			return node, nil
		}
	}
}

// unary = ("+" | "-")? unary | primary
func (p *Parser) unary() (*Node, error) {
	if p.match("+") {
		p.advance()
		node, err := p.unary()
		if err != nil {
			return nil, err
		}
		return node, nil
	}

	if p.match("-") {
		p.advance()
		node, err := p.unary()
		if err != nil {
			return nil, err
		}
		// "-" unary is equivalent to 0 - val
		return &Node{Kind: SUB, Lhs: &Node{Kind: NUM, Val: 0}, Rhs: node}, nil
	}

	return p.primary()
}

// primary = num | ident | "(" expr ")"
func (p *Parser) primary() (*Node, error) {
	if p.match("(") {
		p.advance()
		node, err := p.expr()
		if err != nil {
			return nil, err
		}
		if err := p.expect(")"); err != nil {
			return nil, err
		}
		return node, nil
	}

	if p.current.Kind == lexer.NUM {
		val := p.current.Val
		p.advance()
		return &Node{Kind: NUM, Val: val}, nil
	} else if p.current.Kind == lexer.IDENT {
		offset := (rune(p.current.Str[0]) - 'a' + 1) * 8
		p.advance()
		return &Node{Kind: LVAR, Offset: int(offset)}, nil
	} else {
		return nil, errors.NewPosError(
			fmt.Sprintf("expected number or identifier, but got %s", p.current.Str),
			p.input,
			p.current.Pos,
		)
	}
}

func (p *Parser) atEnd() bool {
	return p.current == nil || p.current.Kind == lexer.EOF
}

func (p *Parser) advance() {
	if p.current != nil {
		p.current = p.current.Next
	}
}

func (p *Parser) match(op string) bool {
	return p.current != nil && p.current.Str == op
}

func (p *Parser) expect(op string) error {
	if p.current == nil {
		return errors.NewPosError(
			fmt.Sprintf("expected %s, but got EOF", op),
			p.input,
			len(p.input), // assuming EOF is at the end of the input
		)
	}
	if p.current.Str != op {
		return errors.NewPosError(
			fmt.Sprintf("expected %s, but got %s", op, p.current.Str),
			p.input,
			p.current.Pos,
		)
	}
	p.advance()
	return nil
}
