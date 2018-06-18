package ast

import (
	"github.com/go-llvm/llvm"
	"github.com/nickwanninger/act/pkg/types"
)

var foo llvm.Type
var foobar llvm.Value

// NodeType -
type NodeType int

// Node -
type Node interface {
	Kind() NodeType
	NameString() string
	Codegen(*Scope) llvm.Value
}

// Kind -
func (t NodeType) Kind() NodeType {
	return t
}

// Inferer interface allows anything that implements it to assume/infer
// the type of itself and return it as a string
type Inferer interface {
	InferType() string
}

// func (t NodeType) String() string {
// 	return "nothing"
// }

const (
	// literals
	nodeInt NodeType = iota
	nodeFloat
	nodeString

	// expressions
	nodeIf
	nodeFor
	nodeUnary  // unary operator  (!, ...)
	nodeBinary // binary operator (+, -, *, ...)

	nodeFnCall
	nodeVariable
	nodeVariableDecl
	nodeVariableReference

	// non-expression statements
	nodeFnPrototype
	nodeFunction
	nodeFunctionCall

	// Other
	nodeBlock
)

type intNode struct {
	NodeType
	Value int64
}

func (n intNode) NameString() string { return "intNode" }

type floatNode struct {
	NodeType
	Value float64
}

func (n floatNode) NameString() string { return "floatNode" }

type stringNode struct {
	NodeType
	Value string
}

func (n stringNode) NameString() string { return "stringNode" }

type ifNode struct {
	NodeType

	// funky notation because `if` and `else` are go keywords
	If   Node
	Then Node
	Else Node
}

func (n ifNode) NameString() string { return "ifNode" }

type forNode struct {
	NodeType

	Counter string
	Start   Node
	Test    Node
	Step    Node
	Body    Node
}

func (n forNode) NameString() string { return "forNode" }

type unaryNode struct {
	NodeType

	Name    string
	Operand Node
}

func (n unaryNode) NameString() string { return "unaryNode" }

type binaryNode struct {
	NodeType

	OP    string
	Left  Node
	Right Node
}

func (n binaryNode) NameString() string { return "binaryNode" }

type fnCallNode struct {
	NodeType
	Calee string
	Args  []Node
}

func (n fnCallNode) NameString() string { return "fnCallNode" }

type variableReferenceNode struct {
	NodeType
	Name string
}

func (n variableReferenceNode) NameString() string { return "variableReferenceNode" }

type variableNode struct {
	NodeType
	Type     *types.VarType
	HasValue bool
	Name     string
	IsArray  bool
	Body     Node
}

func (n variableNode) NameString() string { return "variableNode" }

type returnNode struct {
	NodeType
	Value Node
}

func (n returnNode) NameString() string { return "returnNode" }

type functionNode struct {
	NodeType

	Name       string
	Args       []variableNode
	Body       blockNode
	ReturnType *types.VarType
}

func (n functionNode) NameString() string { return "functionNode" }

type functionCallNode struct {
	NodeType

	Name string
	Args []Node
}

func (n functionCallNode) NameString() string { return "functionCallNode" }

type blockNode struct {
	NodeType
	Nodes []Node
}

func (n blockNode) NameString() string { return "blockNode" }

type whileNode struct {
	NodeType

	Predicate Node
	Body      Node
}

func (n whileNode) NameString() string { return "whileNode" }
