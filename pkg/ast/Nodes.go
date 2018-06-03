package ast

import (
// "github.com/nickwanninger/llvm"
)

// Node -
type Node interface {
	Kind() nodeType
	// codegen() llvm.Value
}

type nodeType int

func (t nodeType) Kind() nodeType {
	return t
}

const (
	// literals
	nodeNumber nodeType = iota

	// expressions
	nodeIf
	nodeFor
	nodeUnary // unary operator  (!, ...)
	nodeBiary // binary operator (+, -, *, ...)

	nodeFnCall
	nodeVariable
	nodeVarialbeExpr

	// non-expression statements
	nodeFnPrototype
	nodeFunction

	// Other
	nodeList
)

type numberNode struct {
	nodeType
	val float64
}

type ifNode struct {
	nodeType

	// funky notation because `if` and `else` are go keywords
	ifN   Node
	thenN Node
	elseN Node
}

type forNode struct {
	nodeType

	counter string
	start   Node
	test    Node
	step    Node
	body    Node
}

type unaryNode struct {
	nodeType

	name    string
	operand Node
}

type binaryNode struct {
	nodeType

	op    string
	left  Node
	right Node
}

type fnCallNode struct {
	nodeType
	callee string
	args   []Node
}

type variableNameNode struct {
	nodeType
	name string
}

type variableNode struct {
	nodeType
	typ  string
	name string
}

type variableExprNode struct {
	nodeType
	typ string
	// vars []struct {
	// 	name string
	// 	node Node
	// }
	body Node
}

type functionNode struct {
	nodeType

	name string
	args []variableNode
	body Node
}

type blockNode struct {
	nodeType

	nodes []Node
}
