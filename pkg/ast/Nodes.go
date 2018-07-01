package ast

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nickwanninger/geode/pkg/typesystem"
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
	nodeWhile
	nodeFor
	nodeUnary  // unary operator  (!, ...)
	nodeBinary // binary operator (+, -, *, ...)

	nodeFnCall
	nodeCast
	nodeVariable
	nodeVariableDecl
	nodeVariableReference

	// non-expression statements
	nodeFnPrototype
	nodeFunction
	nodeFunctionCall
	nodeClass

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

type castNode struct {
	NodeType
	From Node
	To   string
}

func (n castNode) NameString() string                { return "castNode" }
func (n castNode) InferType(scope *Scope) types.Type { return typesystem.GlobalTypeMap.GetType(n.To) }

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
	Index int
	Init  Node
	Cond  Node
	Step  Node
	Body  Node
}

func (n forNode) NameString() string                { return "forNode" }
func (n forNode) InferType(scope *Scope) types.Type { return types.Void }

type unaryNode struct {
	NodeType

	Operator string
	Operand  Node
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

// type variableReferenceNode struct {
// 	NodeType
// 	Index     int
// 	IndexExpr Node
// 	Name      string
// }

// func (n variableReferenceNode) NameString() string                { return "variableReferenceNode" }
// func (n variableReferenceNode) InferType(scope *Scope) types.Type { return types.Void }

// ReferenceType is how we go about accessing a variable. Do we just
// want the value, or do we want to assign to it
type ReferenceType int

// The different ways you can access a variableNode
const (
	ReferenceDefine ReferenceType = iota
	ReferenceAssign
	ReferenceAccess
)

type variableNode struct {
	NodeType
	Type         GeodeTypeRef
	HasValue     bool
	Name         string
	IsPointer    bool
	RefType      ReferenceType
	IndexExpr    Node
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
	External   bool
	Variadic   bool
	Nomangle   bool
	ReturnType GeodeTypeRef
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

	If    Node
	Body  Node
	Index int
}

func (n whileNode) NameString() string                { return "whileNode" }
func (n whileNode) InferType(scope *Scope) types.Type { return types.Void }

type classNode struct {
	NodeType

	Name      string
	Methods   []functionNode
	Variables []variableNode
}

func (n classNode) NameString() string                { return "classNode" }
func (n classNode) InferType(scope *Scope) types.Type { return types.Void }

// GeodeTypeRef -
type GeodeTypeRef struct {
	Array        bool
	PointerLevel int
	Name         string
}

// BuildPointerType will take some type and apply a level of nested pointers
func (r GeodeTypeRef) BuildPointerType(t types.Type) types.Type {
	for i := 0; i < r.PointerLevel; i++ {
		t = types.NewPointer(t)
	}
	return t
}
