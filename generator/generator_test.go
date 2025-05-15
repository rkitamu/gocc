package generator_test

import (
	"strings"
	"testing"

	"rkitamu/gocc/generator"
	"rkitamu/gocc/parser"
)

func TestGenerator_Addition(t *testing.T) {
	node := &parser.Node{
		Kind: parser.ADD,
		Lhs:  &parser.Node{Kind: parser.NUM, Val: 1},
		Rhs:  &parser.Node{Kind: parser.NUM, Val: 2},
	}

	gen := generator.NewGenerator()
	asm := gen.Generate(node)

	checks := []string{
		"push 1",
		"push 2",
		"pop rdi",
		"pop rax",
		"add rax, rdi",
		"push rax",
		"pop rax",
		"ret",
	}

	for _, line := range checks {
		if !strings.Contains(asm, line) {
			t.Errorf("expected line '%s' not found in generated assembly:\n%s", line, asm)
		}
	}
}

func TestGenerator_Equality(t *testing.T) {
	node := &parser.Node{
		Kind: parser.EQ,
		Lhs:  &parser.Node{Kind: parser.NUM, Val: 42},
		Rhs:  &parser.Node{Kind: parser.NUM, Val: 42},
	}

	gen := generator.NewGenerator()
	asm := gen.Generate(node)

	expected := []string{
		"cmp rax, rdi",
		"sete al",
		"movzb rax, al",
	}

	for _, e := range expected {
		if !strings.Contains(asm, e) {
			t.Errorf("expected instruction '%s' not found in:\n%s", e, asm)
		}
	}
}

func TestGenerator_Subtraction(t *testing.T) {
	node := &parser.Node{
		Kind: parser.SUB,
		Lhs:  &parser.Node{Kind: parser.NUM, Val: 5},
		Rhs:  &parser.Node{Kind: parser.NUM, Val: 3},
	}

	gen := generator.NewGenerator()
	asm := gen.Generate(node)

	if !strings.Contains(asm, "sub rax, rdi") {
		t.Errorf("sub instruction missing in:\n%s", asm)
	}
}

func TestGenerator_NestedExpr(t *testing.T) {
	node := &parser.Node{
		Kind: parser.ADD,
		Lhs:  &parser.Node{Kind: parser.NUM, Val: 1},
		Rhs: &parser.Node{
			Kind: parser.MUL,
			Lhs:  &parser.Node{Kind: parser.NUM, Val: 2},
			Rhs:  &parser.Node{Kind: parser.NUM, Val: 3},
		},
	}

	gen := generator.NewGenerator()
	asm := gen.Generate(node)

	if !strings.Contains(asm, "imul rax, rdi") {
		t.Errorf("imul instruction missing in nested expression:\n%s", asm)
	}
}

func TestGenerator_LessThan(t *testing.T) {
	node := &parser.Node{
		Kind: parser.LT,
		Lhs:  &parser.Node{Kind: parser.NUM, Val: 1},
		Rhs:  &parser.Node{Kind: parser.NUM, Val: 2},
	}

	gen := generator.NewGenerator()
	asm := gen.Generate(node)

	expected := []string{
		"cmp rax, rdi",
		"setl al",
		"movzb rax, al",
	}

	for _, line := range expected {
		if !strings.Contains(asm, line) {
			t.Errorf("expected '%s' in:\n%s", line, asm)
		}
	}
}
