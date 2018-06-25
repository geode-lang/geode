package gen

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// NodeType -
type NodeType int

// Node -
type Node interface {
	Kind() NodeType
	NameString() string
	Codegen(*Scope, *Compiler) value.Value
	InferType(scope *Scope) types.Type
}

// Kind -
func (t NodeType) Kind() NodeType {
	return t
}

// func (t NodeType) String() string {
// 	return "nothing"
// }

const (
	// literals
	nodeInt NodeType = iota
	nodeFloat
	nodeString
	nodeChar

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

func (n intNode) NameString() string                { return "intNode" }
func (n intNode) InferType(scope *Scope) types.Type { return types.I64 }

type floatNode struct {
	NodeType
	Value float64
}

func (n floatNode) NameString() string                { return "floatNode" }
func (n floatNode) InferType(scope *Scope) types.Type { return types.Double }

type stringNode struct {
	NodeType
	Value string
}

func (n stringNode) NameString() string                { return "stringNode" }
func (n stringNode) InferType(scope *Scope) types.Type { return types.NewPointer(types.I8) }

type charNode struct {
	NodeType
	Value int8
}

func (n charNode) NameString() string                { return "charNode" }
func (n charNode) InferType(scope *Scope) types.Type { return types.I8 }

type ifNode struct {
	NodeType
	If    Node
	Then  Node
	Else  Node
	Index int
}

func (n ifNode) NameString() string                { return "ifNode" }
func (n ifNode) InferType(scope *Scope) types.Type { return types.Void }

type forNode struct {
	NodeType

	Counter string
	Start   Node
	Test    Node
	Step    Node
	Body    Node
}

func (n forNode) NameString() string                { return "forNode" }
func (n forNode) InferType(scope *Scope) types.Type { return types.Void }

type unaryNode struct {
	NodeType

	Name    string
	Operand Node
}

func (n unaryNode) NameString() string                { return "unaryNode" }
func (n unaryNode) InferType(scope *Scope) types.Type { return n.Operand.InferType(scope) }

type binaryNode struct {
	NodeType

	OP    string
	Left  Node
	Right Node
}

func (n binaryNode) NameString() string                { return "binaryNode" }
func (n binaryNode) InferType(scope *Scope) types.Type { return types.Void }

type fnCallNode struct {
	NodeType
	Calee string
	Args  []Node
}

func (n fnCallNode) NameString() string                { return "fnCallNode" }
func (n fnCallNode) InferType(scope *Scope) types.Type { return types.Void }

type variableReferenceNode struct {
	NodeType
	Name string
}

func (n variableReferenceNode) NameString() string                { return "variableReferenceNode" }
func (n variableReferenceNode) InferType(scope *Scope) types.Type { return types.Void }

type variableNode struct {
	NodeType
	Type         types.Type
	HasValue     bool
	Name         string
	IsPointer    bool
	IsArray      bool
	Reassignment bool
	Body         Node
}

func (n variableNode) NameString() string                { return "variableNode" }
func (n variableNode) InferType(scope *Scope) types.Type { return types.Void }

type returnNode struct {
	NodeType
	Value Node
}

func (n returnNode) NameString() string                { return "returnNode" }
func (n returnNode) InferType(scope *Scope) types.Type { return types.Void }

type functionNode struct {
	NodeType

	Name       string
	Args       []variableNode
	Body       blockNode
	IsExternal bool
	ReturnType types.Type
}

func (n functionNode) NameString() string                { return "functionNode" }
func (n functionNode) InferType(scope *Scope) types.Type { return types.Void }

type functionCallNode struct {
	NodeType

	Name string
	Args []Node
}

func (n functionCallNode) NameString() string                { return "functionCallNode" }
func (n functionCallNode) InferType(scope *Scope) types.Type { return types.Void }

type blockNode struct {
	NodeType
	Nodes []Node
}

func (n blockNode) NameString() string                { return "blockNode" }
func (n blockNode) InferType(scope *Scope) types.Type { return types.Void }

type whileNode struct {
	NodeType

	Predicate Node
	Body      Node
}

func (n whileNode) NameString() string                { return "whileNode" }
func (n whileNode) InferType(scope *Scope) types.Type { return types.Void }
