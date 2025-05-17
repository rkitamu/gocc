package generator

import (
	"fmt"
	"rkitamu/gocc/parser"

	"strings"
)

type Generator struct {
	sb *strings.Builder
}

func NewGenerator() *Generator {
	return &Generator{
		sb: &strings.Builder{},
	}
}

func (g *Generator) Generate(node *parser.Node) (string, error) {
	g.emit(".intel_syntax noprefix")
	g.emit(".global main")
	g.emit("main:")

	// TODO: now we only support a single expression
	if err := g.emitExpr(node); err != nil {
		return "", err
	}
	g.emit("  pop rax")
	g.emit("  ret")

	g.emit(".section .note.GNU-stack,\"\",@progbits")
	return g.sb.String(), nil
}

func (g *Generator) GenerateForMultiStatement(node []*parser.Node) (string, error) {
	g.emit(".intel_syntax noprefix")
	g.emit(".global main")
	g.emit("main:")

	// 変数26個分の領域を確保
	g.emit("  push rbp")
	g.emit("  mov rbp, rsp")
	g.emit("  sub rsp, 208")

	for _, n := range node {
		err := g.emitExpr(n)
		if err != nil {
			return "", err
		}
		g.emit("  pop rax")
	}

	g.emit(".section .note.GNU-stack,\"\",@progbits")
	return g.sb.String(), nil
}

func (g *Generator) emit(line string) {
	fmt.Fprintln(g.sb, line)
}

func (g *Generator) emitLval(node *parser.Node) error {
	if node.Kind == parser.LVAR {
		g.emit("  mov rax, rbp")
		g.emit(fmt.Sprintf("  sub rax, %d", node.Offset))
		g.emit("  push rax")
	} else {
		return fmt.Errorf("not lval: ")
	}

	return nil
}

func (g *Generator) emitExpr(node *parser.Node) error {
	if node.Kind == parser.NUM {
		g.emit(fmt.Sprintf("  push %d", node.Val))
		return nil
	} else if node.Kind == parser.LVAR {
		err := g.emitLval(node)
		if err != nil {
			return err
		}
		g.emit("  pop rax")
		g.emit("  mov rax, [rax]")
		g.emit("  push rax")
		return nil
	} else if node.Kind == parser.ASSIGN {
		if err := g.emitLval(node.Lhs); err != nil {
			return err
		}
		if err := g.emitExpr(node.Rhs); err != nil {
			return err
		}

		g.emit("  pop rdi")
		g.emit("  pop rax")
		g.emit("  mov [rax], rdi")
		g.emit("  push rdi")
		return nil
	} else if node.Kind == parser.RETURN {
		if err := g.emitExpr(node.Lhs); err != nil {
			return err
		}
		g.emit("  pop rax")
		g.emit("  mov rsp, rbp")
		g.emit("  pop rbp")
		g.emit("  ret")
		return nil
	}

	if err := g.emitExpr(node.Lhs); err != nil {
		return err
	}
	if err := g.emitExpr(node.Rhs); err != nil {
		return err
	}

	g.emit("  pop rdi")
	g.emit("  pop rax")

	switch node.Kind {
	case parser.ADD:
		g.emit("  add rax, rdi")
	case parser.SUB:
		g.emit("  sub rax, rdi")
	case parser.MUL:
		g.emit("  imul rax, rdi")
	case parser.DIV:
		g.emit("  cqo")
		g.emit("  idiv rdi")
	case parser.EQ:
		g.emit("  cmp rax, rdi")
		g.emit("  sete al")
		g.emit("  movzb rax, al")
	case parser.NEQ:
		g.emit("  cmp rax, rdi")
		g.emit("  setne al")
		g.emit("  movzb rax, al")
	case parser.LT:
		g.emit("  cmp rax, rdi")
		g.emit("  setl al")
		g.emit("  movzb rax, al")
	case parser.LTE:
		g.emit("  cmp rax, rdi")
		g.emit("  setle al")
		g.emit("  movzb rax, al")
	}

	g.emit("  push rax")

	return nil
}
