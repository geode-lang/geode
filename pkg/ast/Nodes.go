package ast

import (
	"fmt"

	"github.com/geode-lang/geode/pkg/lexer"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// NodeType -
type NodeType string

// Kind -
func (t NodeType) Kind() NodeType {
	return t
}

// TokenReference -
type TokenReference struct {
	Token lexer.Token
}

// SyntaxError -
func (t TokenReference) SyntaxError() {
	t.Token.SyntaxError()
}

// Node -
type Node interface {
	Kind() NodeType
	SyntaxError()
	NameString() string
	Codegen(*Scope, *Compiler) value.Value
	InferType(scope *Scope) string
}

// func MarshalNodeToBinary(n Node) ([]byte, error) {

// }
// func (t NodeType) String() string {
// 	return "nothing"
// }

const (
	// literals
	nodeInt          NodeType = "nodeInt"
	nodeFloat                 = "nodeFloat"
	nodeString                = "nodeString"
	nodeChar                  = "nodeChar"
	nodeIf                    = "nodeIf"
	nodeWhile                 = "nodeWhile"
	nodeFor                   = "nodeFor"
	nodeUnary                 = "nodeUnary"
	nodeBinary                = "nodeBinary"
	nodeFnCall                = "nodeFnCall"
	nodeAssignment            = "nodeAssignment"
	nodeVariable              = "nodeVariable"
	nodeVariableDecl          = "nodeVariableDecl"
	nodeFunction              = "nodeFunction"
	nodeFunctionCall          = "nodeFunctionCall"
	nodeClass                 = "nodeClass"
	nodeDependency            = "nodeDependency"
	nodeNamespace             = "nodeNamespace"
	nodeBlock                 = "nodeBlock"
	nodeSubscript             = "nodeSubscript"
	nodeArray                 = "nodeArray"
	nodeDot                   = "nodeDot"
	nodeSizeof                = "nodeSizeof"
	nodeCast                  = "nodeCast"
)

//
//
// StringNode is a string literal

// CharNode is a char literal
// TODO: get parsing working for this.
type CharNode struct {
	NodeType
	TokenReference

	Value int8
}

// NameString implements Node.NameString
func (n CharNode) NameString() string { return "CharNode" }

// InferType implements Node.InferType
func (n CharNode) InferType(scope *Scope) string { return "byte" }

// IfNode is an if statement representation
type IfNode struct {
	NodeType
	TokenReference

	If    Node
	Then  Node
	Else  Node
	Index int
}

// NameString implements Node.NameString
func (n IfNode) NameString() string { return "IfNode" }

// InferType implements Node.InferType
func (n IfNode) InferType(scope *Scope) string { return "void" }

//
// ForNode is a for loop structure representation
type ForNode struct {
	NodeType
	TokenReference

	Index int
	Init  Node
	Cond  Node
	Step  Node
	Body  Node
}

// NameString implements Node.NameString
func (n ForNode) NameString() string { return "ForNode" }

// InferType implements Node.InferType
func (n ForNode) InferType(scope *Scope) string { return "void" }

//
// UnaryNode is a unary operation representation.
// Example:
//     * !a
//     * &value
//
type UnaryNode struct {
	NodeType
	TokenReference

	Operator string
	Operand  Node
}

// NameString implements Node.NameString
func (n UnaryNode) NameString() string { return "UnaryNode" }

// InferType implements Node.InferType
func (n UnaryNode) InferType(scope *Scope) string { return n.Operand.InferType(scope) }

//
// BinaryNode is a binary operation representation
type BinaryNode struct {
	NodeType
	TokenReference
	Accessable

	OP    string
	Left  Node
	Right Node
}

// NameString implements Node.NameString
func (n BinaryNode) NameString() string { return "BinaryNode" }

// InferType implements Node.InferType
func (n BinaryNode) InferType(scope *Scope) string { return n.Left.InferType(scope) }

// GenAccess implements Accessable.GenAccess
func (n BinaryNode) GenAccess(scope *Scope, c *Compiler) value.Value {
	return n.Codegen(scope, c)
}

func (n BinaryNode) String() string {
	return fmt.Sprintf("%s %s %s", n.Left, n.OP, n.Right)
}

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
	TokenReference

	Paths    []string
	CLinkage bool
}

// NameString implements Node.NameString
func (n DependencyNode) NameString() string { return "DependencyNode" }

// InferType implements Node.InferType
func (n DependencyNode) InferType(scope *Scope) string { return "void" }

// ReferenceType is how we go about accessing a variable. Do we just
// want the value, or do we want to assign to it
type ReferenceType int

// The different ways you can access a VariableNode
const (
	ReferenceDefine ReferenceType = iota
	ReferenceAssign
	ReferenceAccessValue
	ReferenceAccessStackAddress
	ReferenceDereference
)

// ReturnNode is how functions return values from any block
// A return node contains the value (another Node) that will be
// codegenned and used in a `NewRet()` call on the parent function
type ReturnNode struct {
	NodeType
	TokenReference

	Value Node
}

// NameString implements Node.NameString
func (n ReturnNode) NameString() string { return "ReturnNode" }

// InferType implements Node.InferType
func (n ReturnNode) InferType(scope *Scope) string { return n.Value.InferType(scope) }

// FunctionCallNode is a function call, example: `foo(a, b, c)`. This would be:
//    Name = "foo"
//    Args = [a, b, c]    <- these are Node references
type FunctionCallNode struct {
	NodeType
	TokenReference

	Name     Reference
	Args     []Node
	Generics []*GenericSymbol
}

// NameString implements Node.NameString
func (n FunctionCallNode) NameString() string { return "FunctionCallNode" }

// InferType implements Node.InferType
func (n FunctionCallNode) InferType(scope *Scope) string {
	return "PLEASE IMPLEMENT ME :)"
}

// GenAccess implements Accessable.GenAccess
func (n FunctionCallNode) GenAccess(s *Scope, c *Compiler) value.Value {
	return n.Codegen(s, c)
}

// WhileNode is a while loop representationvbnm,bvbnm
type WhileNode struct {
	NodeType
	TokenReference

	If    Node
	Body  Node
	Index int
}

// NameString implements Node.NameString
func (n WhileNode) NameString() string { return "WhileNode" }

// InferType implements Node.InferType
func (n WhileNode) InferType(scope *Scope) string { return "void" }

// NamespaceNode -
type NamespaceNode struct {
	NodeType
	TokenReference

	Name string
}

// NameString implements Node.NameString
func (n NamespaceNode) NameString() string { return "NamespaceNode" }

// InferType implements Node.InferType
func (n NamespaceNode) InferType(scope *Scope) string { return "void" }

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
