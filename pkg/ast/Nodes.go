package ast

import (
	"bytes"
	"fmt"

	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
	"github.com/geode-lang/geode/pkg/lexer"
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
	fmt.Stringer
	Kind() NodeType
	SyntaxError()
	NameString() string
	Codegen(*Program) (value.Value, error)
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
	nodeTypeInfo              = "nodeTypeInfo"
	nodeCast                  = "nodeCast"
	nodeBool                  = "nodeBool"
	nodeGlobalDecl            = "nodeGlobalDecl"
	nodeNil                   = "nodeNil"
	nodeIdent                 = "nodeIdent"
	nodeStringFormat          = "nodeStringFormat"
)

//
//
// StringNode is a string literal

// CharNode is a char literal
// TODO: get parsing working for this.
type CharNode struct {
	NodeType
	TokenReference

	Value rune
}

func (n CharNode) String() string {
	buff := &bytes.Buffer{}
	fmt.Fprintf(buff, "CharNode")
	return buff.String()
}

// NameString implements Node.NameString
func (n CharNode) NameString() string { return "CharNode" }

// IfNode is an if statement representation
type IfNode struct {
	NodeType
	TokenReference

	If    Node
	Then  Node
	Else  Node
	Index int
}

func (n IfNode) String() string {
	buff := &bytes.Buffer{}
	fmt.Fprintf(buff, "if %s %s", n.If, n.Then)

	if n.Else != nil {
		fmt.Fprintf(buff, " else %s", n.Else)
	}
	return buff.String()
}

// NameString implements Node.NameString
func (n IfNode) NameString() string { return "IfNode" }

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

func (n UnaryNode) String() string {
	buff := &bytes.Buffer{}
	fmt.Fprintf(buff, "UnaryNode")
	return buff.String()
}

// NameString implements Node.NameString
func (n UnaryNode) NameString() string { return "UnaryNode" }

// DependencyNode is a way of representing the need to include
// a dependency or multiple dependencies. It also works to link
// a c program as well. Paths contains a list of paths to the dependencies
// that the user entered into the statement. These paths are not resolved
// and may not contain a geode source file.
//
// Example:
//    Paths = ["io"]
//    CLinkage = false
///
type DependencyNode struct {
	NodeType
	TokenReference

	Paths    []string
	CLinkage bool
}

func (n DependencyNode) String() string {
	buff := &bytes.Buffer{}

	if n.CLinkage {
		fmt.Fprintf(buff, "link ")
	} else {
		fmt.Fprintf(buff, "include ")
	}

	for i, path := range n.Paths {
		fmt.Fprintf(buff, "\"%s\"", path)

		if i < len(n.Paths)-1 {
			fmt.Fprintf(buff, ", ")
		}
	}
	return buff.String()
}

// NameString implements Node.NameString
func (n DependencyNode) NameString() string { return "DependencyNode" }

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

func (n ReturnNode) String() string {
	buff := &bytes.Buffer{}
	fmt.Fprintf(buff, "return %s", n.Value)
	return buff.String()
}

// NameString implements Node.NameString
func (n ReturnNode) NameString() string { return "ReturnNode" }

// WhileNode is a while loop representationvbnm,bvbnm
type WhileNode struct {
	NodeType
	TokenReference

	If    Node
	Body  Node
	Index int
}

func (n WhileNode) String() string {
	buff := &bytes.Buffer{}
	fmt.Fprintf(buff, "while %s %s", n.If, n.Body)
	return buff.String()
}

// NameString implements Node.NameString
func (n WhileNode) NameString() string { return "WhileNode" }

// NamespaceNode -
type NamespaceNode struct {
	NodeType
	TokenReference

	Name string
}

func (n NamespaceNode) String() string {
	buff := &bytes.Buffer{}
	fmt.Fprintf(buff, "is %s", n.Name)
	return buff.String()
}

// NameString implements Node.NameString
func (n NamespaceNode) NameString() string { return "NamespaceNode" }

// TypeModifier is something like *, [], ?, etc...
type TypeModifier byte

// setup the constants to use globally
const (
	ModifierPointer TypeModifier = iota
	ModifierSlice
	ModifierUnknown
)

// TypeNode -
type TypeNode struct {
	PointerLevel int
	Unknown      bool
	Name         string

	Modifiers []TypeModifier
}

func (n TypeNode) String() string {

	buff := &bytes.Buffer{}

	fmt.Fprintf(buff, "%s", n.Name)

	for _, mod := range n.Modifiers {
		switch mod {
		case ModifierPointer:
			fmt.Fprintf(buff, "*")
		case ModifierSlice:
			fmt.Fprintf(buff, "[]")
		}
	}

	return buff.String()
}

// GetType returns the llvm type representation of the TypeNode
func (n TypeNode) GetType(prog *Program) (types.Type, error) {
	var ty types.Type
	var err error
	ty, err = prog.FindType(n.Name)
	if err != nil {
		return nil, err
	}
	if len(n.Modifiers) > 0 {
		for _, mod := range n.Modifiers {
			switch mod {
			case ModifierPointer:
				ty = types.NewPointer(ty)
			case ModifierSlice:
				ty = types.NewSlice(ty)
			case ModifierUnknown:
				//
			default:
				return nil, fmt.Errorf("unknown type modifier %d on type %q", mod, n)
			}
		}
	}

	return ty, nil
}
