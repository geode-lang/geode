package ast

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
	nodeWhile
	nodeFor
	nodeUnary  // unary operator  (!, ...)
	nodeBinary // binary operator (+, -, *, ...)

	nodeFnCall
	nodeCast
	nodeVariable
	nodeVariableDecl
	nodeVariableReference

	nodeCompoundAssignment

	// non-expression statements
	nodeFnPrototype
	nodeFunction
	nodeFunctionCall
	nodeClass
	nodeDependency
	nodeNamespace

	// Other
	nodeBlock
)

// IntNode is an integer literal
type IntNode struct {
	NodeType
	Value int64
}

// NameString implements Node.NameString
func (n IntNode) NameString() string { return "IntNode" }

// InferType implements Node.InferType
func (n IntNode) InferType(scope *Scope) types.Type { return types.I64 }

//
//
// FloatNode is a float literla
type FloatNode struct {
	NodeType
	Value float64
}

// NameString implements Node.NameString
func (n FloatNode) NameString() string { return "FloatNode" }

// InferType implements Node.InferType
func (n FloatNode) InferType(scope *Scope) types.Type { return types.Double }

//
//
// StringNode is a string literal
type StringNode struct {
	NodeType
	Value string
}

// NameString implements Node.NameString
func (n StringNode) NameString() string { return "StringNode" }

// InferType implements Node.InferType
func (n StringNode) InferType(scope *Scope) types.Type { return types.NewPointer(types.I8) }

// CharNode is a char literal
// TODO: get parsing working for this.
type CharNode struct {
	NodeType
	Value int8
}

// NameString implements Node.NameString
func (n CharNode) NameString() string { return "CharNode" }

// InferType implements Node.InferType
func (n CharNode) InferType(scope *Scope) types.Type { return types.I8 }

// CastNode is a type cast "function" call. TODO: Replace this with normal function calls and check
// in the function call codegen function
type CastNode struct {
	NodeType
	From Node
	To   string
}

// NameString implements Node.NameString
func (n CastNode) NameString() string { return "CastNode" }

// InferType implements Node.InferType
func (n CastNode) InferType(scope *Scope) types.Type { return scope.FindType(n.To).Type }

// IfNode is an if statement representation
type IfNode struct {
	NodeType
	If    Node
	Then  Node
	Else  Node
	Index int
}

// NameString implements Node.NameString
func (n IfNode) NameString() string { return "IfNode" }

// InferType implements Node.InferType
func (n IfNode) InferType(scope *Scope) types.Type { return types.Void }

//
// ForNode is a for loop structure representation
type ForNode struct {
	NodeType
	Index int
	Init  Node
	Cond  Node
	Step  Node
	Body  Node
}

// NameString implements Node.NameString
func (n ForNode) NameString() string { return "ForNode" }

// InferType implements Node.InferType
func (n ForNode) InferType(scope *Scope) types.Type { return types.Void }

//
// UnaryNode is a unary operation representation.
// Example:
//     * !a
//     * &value
//
type UnaryNode struct {
	NodeType

	Operator string
	Operand  Node
}

// NameString implements Node.NameString
func (n UnaryNode) NameString() string { return "UnaryNode" }

// InferType implements Node.InferType
func (n UnaryNode) InferType(scope *Scope) types.Type { return n.Operand.InferType(scope) }

//
// BinaryNode is a binary operation representation
type BinaryNode struct {
	NodeType

	OP    string
	Left  Node
	Right Node
}

// NameString implements Node.NameString
func (n BinaryNode) NameString() string { return "BinaryNode" }

// InferType implements Node.InferType
func (n BinaryNode) InferType(scope *Scope) types.Type { return types.Void }

// DependencyNode is a way of representing the need to include
// a dependency or multiple dependencies. It also works to link
// a c program as well. Paths contains a list of paths to the dependencies
// that the user entered into the statement. These paths are not resolved
// and may not contain a geode source file.
//
// Example:
//    Paths = ["std:io"]
//    CLinkage = false
///
type DependencyNode struct {
	NodeType
	Paths    []string
	CLinkage bool
}

// NameString implements Node.NameString
func (n DependencyNode) NameString() string { return "DependencyNode" }

// InferType implements Node.InferType
func (n DependencyNode) InferType(scope *Scope) types.Type { return types.Void }

// ReferenceType is how we go about accessing a variable. Do we just
// want the value, or do we want to assign to it
type ReferenceType int

// The different ways you can access a VariableNode
const (
	ReferenceDefine ReferenceType = iota
	ReferenceAssign
	ReferenceAccess
	ReferenceDereference
)

// ReturnNode is how functions return values from any block
// A return node contains the value (another Node) that will be
// codegenned and used in a `NewRet()` call on the parent function
type ReturnNode struct {
	NodeType
	Value Node
}

// NameString implements Node.NameString
func (n ReturnNode) NameString() string { return "ReturnNode" }

// InferType implements Node.InferType
func (n ReturnNode) InferType(scope *Scope) types.Type { return types.Void }

// FunctionCallNode is a function call, example: `foo(a, b, c)`. This would be:
//    Name = "foo"
//    Args = [a, b, c]    <- these are Node references
type FunctionCallNode struct {
	NodeType

	Name string
	Args []Node
}

// NameString implements Node.NameString
func (n FunctionCallNode) NameString() string { return "FunctionCallNode" }

// InferType implements Node.InferType
func (n FunctionCallNode) InferType(scope *Scope) types.Type { return types.Void }

// WhileNode is a while loop representationvbnm,bvbnm
type WhileNode struct {
	NodeType

	If    Node
	Body  Node
	Index int
}

// NameString implements Node.NameString
func (n WhileNode) NameString() string { return "WhileNode" }

// InferType implements Node.InferType
func (n WhileNode) InferType(scope *Scope) types.Type { return types.Void }

// NamespaceNode -
type NamespaceNode struct {
	NodeType

	Name string
}

// NameString implements Node.NameString
func (n NamespaceNode) NameString() string { return "NamespaceNode" }

// InferType implements Node.InferType
func (n NamespaceNode) InferType(scope *Scope) types.Type { return types.Void }

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
