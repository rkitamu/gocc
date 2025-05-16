package lexer

import "testing"

func TestLexer(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		want    []Token
		wantErr bool
	}{
		{
			name:  "simple test",
			input: "1+2",
			want: []Token{
				{Kind: NUM, Str: "1"},
				{Kind: RESERVED, Str: "+"},
				{Kind: NUM, Str: "2"},
				{Kind: EOF, Str: ""},
			},
			wantErr: false,
		},
		{
			name:  "complex test",
			input: "3*(4-5)",
			want: []Token{
				{Kind: NUM, Str: "3"},
				{Kind: RESERVED, Str: "*"},
				{Kind: RESERVED, Str: "("},
				{Kind: NUM, Str: "4"},
				{Kind: RESERVED, Str: "-"},
				{Kind: NUM, Str: "5"},
				{Kind: RESERVED, Str: ")"},
				{Kind: EOF, Str: ""},
			},
			wantErr: false,
		},
		{
			name:  "digit test",
			input: "1234567890",
			want: []Token{
				{Kind: NUM, Str: "1234567890"},
				{Kind: EOF, Str: ""},
			},
			wantErr: false,
		},
		{
			name:  "plus test",
			input: "1+2",
			want: []Token{
				{Kind: NUM, Str: "1"},
				{Kind: RESERVED, Str: "+"},
				{Kind: NUM, Str: "2"},
				{Kind: EOF, Str: ""},
			},
			wantErr: false,
		},
		{
			name:  "minus test",
			input: "1-2",
			want: []Token{
				{Kind: NUM, Str: "1"},
				{Kind: RESERVED, Str: "-"},
				{Kind: NUM, Str: "2"},
				{Kind: EOF, Str: ""},
			},
			wantErr: false,
		},
		{
			name:  "multiply test",
			input: "1*2",
			want: []Token{
				{Kind: NUM, Str: "1"},
				{Kind: RESERVED, Str: "*"},
				{Kind: NUM, Str: "2"},
				{Kind: EOF, Str: ""},
			},
			wantErr: false,
		},
		{
			name:  "divide test",
			input: "1/2",
			want: []Token{
				{Kind: NUM, Str: "1"},
				{Kind: RESERVED, Str: "/"},
				{Kind: NUM, Str: "2"},
				{Kind: EOF, Str: ""},
			},
			wantErr: false,
		},
		{
			name:  "less than test",
			input: "1<2",
			want: []Token{
				{Kind: NUM, Str: "1"},
				{Kind: RESERVED, Str: "<"},
				{Kind: NUM, Str: "2"},
				{Kind: EOF, Str: ""},
			},
			wantErr: false,
		},
		{
			name:  "greater than test",
			input: "1>2",
			want: []Token{
				{Kind: NUM, Str: "1"},
				{Kind: RESERVED, Str: ">"},
				{Kind: NUM, Str: "2"},
				{Kind: EOF, Str: ""},
			},
			wantErr: false,
		},
		{
			name:  "equal test",
			input: "1 == 2",
			want: []Token{
				{Kind: NUM, Str: "1"},
				{Kind: RESERVED, Str: "=="},
				{Kind: NUM, Str: "2"},
				{Kind: EOF, Str: ""},
			},
			wantErr: false,
		},
		{
			name:  "not equal test",
			input: "1 != 2",
			want: []Token{
				{Kind: NUM, Str: "1"},
				{Kind: RESERVED, Str: "!="},
				{Kind: NUM, Str: "2"},
				{Kind: EOF, Str: ""},
			},
			wantErr: false,
		},
		{
			name:  "less than or equal test",
			input: "1 <= 2",
			want: []Token{
				{Kind: NUM, Str: "1"},
				{Kind: RESERVED, Str: "<="},
				{Kind: NUM, Str: "2"},
				{Kind: EOF, Str: ""},
			},
			wantErr: false,
		},
		{
			name:  "greater than or equal test",
			input: "1 >= 2",
			want: []Token{
				{Kind: NUM, Str: "1"},
				{Kind: RESERVED, Str: ">="},
				{Kind: NUM, Str: "2"},
				{Kind: EOF, Str: ""},
			},
			wantErr: false,
		},
		{
			name:    "error test",
			input:   "1+2a",
			want:    nil,
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			lexer := NewLexer(c.input)
			got, err := lexer.Lex()

			if c.wantErr {
				if err == nil {
					t.Errorf("Lex() expected error, but got none")
				}
				return // ✅ エラーが出るのが期待通りなので終了
			}

			if err != nil {
				t.Errorf("Lex() unexpected error: %v", err)
				return
			}

			// ✅ トークンの比較
			curr := got
			for i, token := range c.want {
				if curr == nil {
					t.Errorf("Lex() returned too few tokens, missing token: %v", token)
					break
				}
				if curr.Kind != token.Kind || curr.Str != token.Str {
					t.Errorf("Token %d: got = %v, want = %v", i, curr, token)
				}
				curr = curr.Next
			}
			if curr != nil {
				t.Errorf("Lex() returned too many tokens, extra: %v", curr)
			}
		})
	}
}
