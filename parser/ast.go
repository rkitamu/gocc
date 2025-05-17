package parser

// NodeKind represents the type of a node in the AST.
type NodeKind int

const (
	ADD    NodeKind = iota // +
	SUB                    // -
	MUL                    // *
	DIV                    // /
	EQ                     // ==
	NEQ                    // !=
	LT                     // <
	LTE                    // <=
	NUM                    // number literal
	ASSIGN                 // =
	LVAR                   // variable
	RETURN                 // return statement
	EOF                    // end of file (optional, not usually needed in AST)
)

// Node represents a node in the abstract syntax tree (AST).
type Node struct {
	Kind   NodeKind // The kind of node (operator, number, etc.)
	Lhs    *Node    // Left-hand side expression
	Rhs    *Node    // Right-hand side expression
	Val    int      // Literal value (only used if Kind == NUM)
	Offset int      // Offset for local variables (only used if Kind == LVAR)
}

type LVar struct {
	Next   *LVar
	Name   string
	Offset int
}
