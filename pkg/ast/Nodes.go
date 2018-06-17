package ast

import (
	"github.com/nickwanninger/llvm"
)

var foo llvm.Value

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
	nodeInt nodeType = iota
	nodeFloat

	// expressions
	nodeIf
	nodeFor
	nodeUnary  // unary operator  (!, ...)
	nodeBinary // binary operator (+, -, *, ...)

	nodeFnCall
	nodeVariable
	nodeVarialbeExpr

	// non-expression statements
	nodeFnPrototype
	nodeFunction

	// Other
	nodeBlock
)

type intNode struct {
	nodeType
	Value int64
}

type floatNode struct {
	nodeType
	Value float64
}

type ifNode struct {
	nodeType

	// funky notation because `if` and `else` are go keywords
	If   Node
	Then Node
	Else Node
}

type forNode struct {
	nodeType

	Counter string
	Start   Node
	Test    Node
	Step    Node
	Body    Node
}

type unaryNode struct {
	nodeType

	Name    string
	Operand Node
}

type binaryNode struct {
	nodeType

	OP    string
	Left  Node
	Right Node
}

type fnCallNode struct {
	nodeType
	Calee string
	Args  []Node
}

type variableNameNode struct {
	nodeType
	Name string
}

type variableNode struct {
	nodeType
	Type string
	Name string
}

type variableExprNode struct {
	nodeType
	Type string
	// vars []struct {
	// 	name string
	// 	node Node
	// }
	Body Node
}

type returnNode struct {
	nodeType
	Value Node
}

type functionNode struct {
	nodeType

	Name       string
	Args       []variableNode
	Body       blockNode
	ReturnType string
}

type blockNode struct {
	nodeType
	Nodes []Node
}

type whileNode struct {
	nodeType

	Predicate Node
	Body      Node
}
