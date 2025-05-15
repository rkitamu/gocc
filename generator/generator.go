package generator

import (
	"fmt"
	"strings"
	"rkitamu/gocc/parser"
)

type Generator struct {
	sb *strings.Builder
}

func NewGenerator() *Generator {
	return &Generator{
		sb: &strings.Builder{},
	}
}

func (g *Generator) Generate(node *parser.Node) string {
	g.emit(".intel_syntax noprefix")
	g.emit(".global main")
	g.emit("main:")

	// TODO: now we only support a single expression
	g.emitExpr(node)
	g.emit("  pop rax")
	g.emit("  ret")

	g.emit("ret")
	g.emit(".section .note.GNU-stack,\"\",@progbits")

	return g.sb.String()
}

func (g *Generator) emit(line string) {
	fmt.Fprintln(g.sb, line)
}

func (g *Generator) emitExpr(node *parser.Node) {
	if node.Kind == parser.NUM {
		g.emit(fmt.Sprintf("  push %d", node.Val))
		return
	}

	g.emitExpr(node.Lhs)
	g.emitExpr(node.Rhs)

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
}
