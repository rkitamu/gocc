package parser

import (
	"fmt"
)

// PrintTree prints the AST in a tree-like format.
func (p Parser) PrintTree(node *Node) {
	printTreeRec(node, "", true)
}

func (p Parser) PrintTreeForMultiStatement(node []*Node) {
	for _, n := range node {
		printTreeRec(n, "", true)
	}
}

func printTreeRec(node *Node, prefix string, isTail bool) {
	if node == nil {
		return
	}

	connector := "├── "
	nextPrefix := prefix + "│   "
	if isTail {
		connector = "└── "
		nextPrefix = prefix + "    "
	}

	label := ""
	if node.Kind == NUM {
		label = fmt.Sprintf("%d", node.Val)
	} else {
		label = fmt.Sprintf("(%s)", nodeKindToString(node.Kind))
	}

	fmt.Printf("%s%s%s\n", prefix, connector, label)

	if node.Lhs != nil || node.Rhs != nil {
		if node.Rhs != nil {
			printTreeRec(node.Lhs, nextPrefix, false)
			printTreeRec(node.Rhs, nextPrefix, true)
		} else if node.Lhs != nil {
			printTreeRec(node.Lhs, nextPrefix, true)
		}
	}
}

func nodeKindToString(kind NodeKind) string {
	switch kind {
	case ADD:
		return "+"
	case SUB:
		return "-"
	case MUL:
		return "*"
	case DIV:
		return "/"
	case EQ:
		return "=="
	case NEQ:
		return "!="
	case LT:
		return "<"
	case LTE:
		return "<="
	case ASSIGN:
		return "="
	case LVAR:
		return "LVAR"
	default:
		return "?"
	}
}
