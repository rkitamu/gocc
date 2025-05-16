package parser_test

import (
	"github.com/google/go-cmp/cmp"
	"rkitamu/gocc/lexer"
	"rkitamu/gocc/parser"
	"testing"
)

func EqualAST(a, b *parser.Node) bool {
	if a == nil || b == nil {
		return a == b
	}
	if a.Kind != b.Kind || a.Val != b.Val {
		return false
	}
	return EqualAST(a.Lhs, b.Lhs) && EqualAST(a.Rhs, b.Rhs)
}

func TestParse_GrammarCoverage(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		tokens *lexer.Token
		want   *parser.Node
	}{
		{
			name:  "equality ==: 1 + 2 == 3",
			input: "1 + 2 == 3",
			tokens: &lexer.Token{
				Kind: lexer.NUM, Val: 1, Str: "1", Next: &lexer.Token{
					Kind: lexer.RESERVED, Str: "+", Next: &lexer.Token{
						Kind: lexer.NUM, Val: 2, Str: "2", Next: &lexer.Token{
							Kind: lexer.RESERVED, Str: "==", Next: &lexer.Token{
								Kind: lexer.NUM, Val: 3, Str: "3", Next: &lexer.Token{
									Kind: lexer.EOF,
								},
							},
						},
					},
				},
			},
			want: &parser.Node{Kind: parser.EQ,
				Lhs: &parser.Node{Kind: parser.ADD,
					Lhs: &parser.Node{Kind: parser.NUM, Val: 1},
					Rhs: &parser.Node{Kind: parser.NUM, Val: 2},
				},
				Rhs: &parser.Node{Kind: parser.NUM, Val: 3},
			},
		},
		{
			name:  "add +: 1 + 2 + 3",
			input: "1 + 2 + 3",
			tokens: &lexer.Token{
				Kind: lexer.NUM, Val: 1, Str: "1", Next: &lexer.Token{
					Kind: lexer.RESERVED, Str: "+", Next: &lexer.Token{
						Kind: lexer.NUM, Val: 2, Str: "2", Next: &lexer.Token{
							Kind: lexer.RESERVED, Str: "+", Next: &lexer.Token{
								Kind: lexer.NUM, Val: 3, Str: "3", Next: &lexer.Token{
									Kind: lexer.EOF,
								},
							},
						},
					},
				},
			},
			want: &parser.Node{Kind: parser.ADD,
				Lhs: &parser.Node{Kind: parser.ADD,
					Lhs: &parser.Node{Kind: parser.NUM, Val: 1},
					Rhs: &parser.Node{Kind: parser.NUM, Val: 2},
				},
				Rhs: &parser.Node{Kind: parser.NUM, Val: 3},
			},
		},
		{
			name:   "unary -: -1 + 2",
			input:  "-1 + 2",
			tokens: &lexer.Token{Kind: lexer.RESERVED, Str: "-", Next: &lexer.Token{Kind: lexer.NUM, Val: 1, Str: "1", Next: &lexer.Token{Kind: lexer.RESERVED, Str: "+", Next: &lexer.Token{Kind: lexer.NUM, Val: 2, Str: "2", Next: &lexer.Token{Kind: lexer.EOF}}}}},
			want: &parser.Node{Kind: parser.ADD,
				Lhs: &parser.Node{Kind: parser.SUB,
					Lhs: &parser.Node{Kind: parser.NUM, Val: 0},
					Rhs: &parser.Node{Kind: parser.NUM, Val: 1},
				},
				Rhs: &parser.Node{Kind: parser.NUM, Val: 2},
			},
		},
		{
			name:   "relational >=: 5 >= 1 + 2",
			input:  "5 >= 1 + 2",
			tokens: &lexer.Token{Kind: lexer.NUM, Val: 5, Str: "5", Next: &lexer.Token{Kind: lexer.RESERVED, Str: ">=", Next: &lexer.Token{Kind: lexer.NUM, Val: 1, Str: "1", Next: &lexer.Token{Kind: lexer.RESERVED, Str: "+", Next: &lexer.Token{Kind: lexer.NUM, Val: 2, Str: "2", Next: &lexer.Token{Kind: lexer.EOF}}}}}},
			want: &parser.Node{Kind: parser.LTE,
				Lhs: &parser.Node{Kind: parser.ADD,
					Lhs: &parser.Node{Kind: parser.NUM, Val: 1},
					Rhs: &parser.Node{Kind: parser.NUM, Val: 2},
				},
				Rhs: &parser.Node{Kind: parser.NUM, Val: 5},
			},
		},
		{
			name:   "grouping (1 + 2) * 3: (1 + 2) * 3",
			input:  "(1 + 2) * 3",
			tokens: &lexer.Token{Kind: lexer.RESERVED, Str: "(", Next: &lexer.Token{Kind: lexer.NUM, Val: 1, Str: "1", Next: &lexer.Token{Kind: lexer.RESERVED, Str: "+", Next: &lexer.Token{Kind: lexer.NUM, Val: 2, Str: "2", Next: &lexer.Token{Kind: lexer.RESERVED, Str: ")", Next: &lexer.Token{Kind: lexer.RESERVED, Str: "*", Next: &lexer.Token{Kind: lexer.NUM, Val: 3, Str: "3", Next: &lexer.Token{Kind: lexer.EOF}}}}}}}},
			want: &parser.Node{Kind: parser.MUL,
				Lhs: &parser.Node{Kind: parser.ADD,
					Lhs: &parser.Node{Kind: parser.NUM, Val: 1},
					Rhs: &parser.Node{Kind: parser.NUM, Val: 2},
				},
				Rhs: &parser.Node{Kind: parser.NUM, Val: 3},
			},
		},
		{
			name:   "not equal !=: 1 != 2",
			input:  "1 != 2",
			tokens: &lexer.Token{Kind: lexer.NUM, Val: 1, Str: "1", Next: &lexer.Token{Kind: lexer.RESERVED, Str: "!=", Next: &lexer.Token{Kind: lexer.NUM, Val: 2, Str: "2", Next: &lexer.Token{Kind: lexer.EOF}}}},
			want: &parser.Node{Kind: parser.NEQ,
				Lhs: &parser.Node{Kind: parser.NUM, Val: 1},
				Rhs: &parser.Node{Kind: parser.NUM, Val: 2},
			},
		},
		{
			name:   "less than <: 1 < 2",
			input:  "1 < 2",
			tokens: &lexer.Token{Kind: lexer.NUM, Val: 1, Str: "1", Next: &lexer.Token{Kind: lexer.RESERVED, Str: "<", Next: &lexer.Token{Kind: lexer.NUM, Val: 2, Str: "2", Next: &lexer.Token{Kind: lexer.EOF}}}},
			want: &parser.Node{Kind: parser.LT,
				Lhs: &parser.Node{Kind: parser.NUM, Val: 1},
				Rhs: &parser.Node{Kind: parser.NUM, Val: 2},
			},
		},
		{
			name:   "multiply *: 2 * 3",
			input:  "2 * 3",
			tokens: &lexer.Token{Kind: lexer.NUM, Val: 2, Str: "2", Next: &lexer.Token{Kind: lexer.RESERVED, Str: "*", Next: &lexer.Token{Kind: lexer.NUM, Val: 3, Str: "3", Next: &lexer.Token{Kind: lexer.EOF}}}},
			want: &parser.Node{Kind: parser.MUL,
				Lhs: &parser.Node{Kind: parser.NUM, Val: 2},
				Rhs: &parser.Node{Kind: parser.NUM, Val: 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.tokens, tt.input)
			got, err := p.Parse()
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("AST mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
