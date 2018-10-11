package ast

import (
	"bytes"
	"fmt"

	"github.com/geode-lang/geode/pkg/lexer"
)

// ExpressionComponents are nodes that make up the expression
// component system defined in parseExpression.go

// ExpComponent is a representation of complex compound expressions
// like call()[1]().foo().bar[12] for example. It is meant to be
// simple to link expression nodes on to, as it is a linked list
// type structure
type ExpComponent interface {
	fmt.Stringer

	Add(ExpComponent)
	Ident() string
	Next() ExpComponent
	ConstructNode(Node) (Node, error)
}

type componentChainNode struct {
	next  ExpComponent
	token lexer.Token
}

func (n *componentChainNode) Add(comp ExpComponent) {
	if n.next != nil {
		n.next.Add(comp)
	} else {
		n.next = comp
	}
}

func (n *componentChainNode) String() string {
	if n.next == nil {
		return ""
	}
	return fmt.Sprintf("%s%s", n.next.Ident(), n.next)
}

func (n *componentChainNode) Next() ExpComponent {
	return n.next
}

// =========================== BaseComponent ===========================

// BaseComponent is the base of a component linked list
// it is used as the generalized start of any component list
type BaseComponent struct {
	componentChainNode
}

// Ident implements ExpComponent.Ident
func (c *BaseComponent) Ident() string {
	return c.String()
}

// ConstructNode returns the ast node for the expression component
func (c *BaseComponent) ConstructNode(prev Node) (Node, error) {
	d := c.Next()
	var err error
	node := prev
	for d != nil {
		node, err = d.ConstructNode(node)
		if err != nil {
			return nil, err
		}
		if node == nil {
			return nil, fmt.Errorf("ConstructNode on %T returned nil, but no error", d)
		}
		d = d.Next()
	}
	return node, nil
}

// =========================== IdentComponent ===========================

// IdentComponent is a component of an expresison that represents
// an identity access
type IdentComponent struct {
	componentChainNode

	Value string
}

// Ident implements ExpComponent.Ident
func (c *IdentComponent) Ident() string {
	return c.Value
}

// ConstructNode returns the ast node for the expression component
func (c *IdentComponent) ConstructNode(prev Node) (Node, error) {
	n := NewIdentNode(c.Value)
	n.Token = c.token
	return n, nil
}

// =========================== IdentDeclComponent ===========================

// IdentDeclComponent is a component of an expresison that represents
// an identity declaration
type IdentDeclComponent struct {
	componentChainNode

	Type TypeNode
	Name IdentNode
}

// Ident implements ExpComponent.Ident
func (c *IdentDeclComponent) Ident() string {
	return fmt.Sprintf("%s %s", c.Type, c.Name)
}

// ConstructNode returns the ast node for the expression component
func (c *IdentDeclComponent) ConstructNode(prev Node) (Node, error) {
	n := VariableDefnNode{}
	n.NodeType = nodeVariableDecl
	n.Token = c.token
	n.Typ = c.Type
	n.Name = c.Name
	return n, nil
}

// =========================== CallComponent ===========================

// CallComponent is an expression component for function calls
type CallComponent struct {
	componentChainNode

	Args []Node
}

// Ident implements ExpComponent.Ident
func (c *CallComponent) Ident() string {
	buf := &bytes.Buffer{}
	buf.WriteString("(")
	for i, arg := range c.Args {
		if arg == nil {
			buf.WriteString("_")
		} else {
			fmt.Fprintf(buf, "%s", arg)
		}

		if i < len(c.Args)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(")")

	return buf.String()
}

// ConstructNode returns the ast node for the expression component
func (c *CallComponent) ConstructNode(prev Node) (Node, error) {

	switch prev.(type) {
	case StringNode:
		n := StringFormatNode{}
		n.Token = c.token
		n.NodeType = nodeStringFormat
		n.Format = prev.(StringNode)
		for _, argc := range c.Args {
			n.Args = append(n.Args, argc)
		}
		return n, nil

	}

	n := FunctionCallNode{}
	n.Token = c.token
	n.NodeType = nodeFunctionCall

	base, ok := prev.(Callable)
	if !ok {
		return nil, fmt.Errorf("function call requires callable - given %T", prev)
	}
	n.Name = base
	for _, argc := range c.Args {
		n.Args = append(n.Args, argc)
	}
	return n, nil
}

// =========================== NumberComponent ===========================

// NumberComponent is an expression component for numbers
type NumberComponent struct {
	componentChainNode

	Value string
}

// ConstructNode returns the ast node for the expression component
func (c *NumberComponent) ConstructNode(prev Node) (Node, error) {
	n, err := GetNumberNodeFromString(c.Value)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, fmt.Errorf("unable to get number type from number component's value")
	}
	return n, nil
}

// Ident implements ExpComponent.Ident
func (c *NumberComponent) Ident() string {
	return c.Value
}

// =========================== SubscriptComponent ===========================

// SubscriptComponent is an expression component for numbers
type SubscriptComponent struct {
	componentChainNode

	Value Node
}

// ConstructNode returns the ast node for the expression component
func (c *SubscriptComponent) ConstructNode(prev Node) (Node, error) {

	n := &SubscriptNode{}
	n.Token = c.token
	n.NodeType = nodeSubscript
	var ok bool
	n.Source, ok = prev.(Accessable)
	if !ok {
		return nil, fmt.Errorf("previous node in SubscriptComponent is not accessable: %T", prev)
	}

	val := c.Value

	n.Index, ok = val.(Accessable)
	if !ok {
		return nil, fmt.Errorf("index node in SubscriptComponent is not accessable: %T, %T", val, c.Value)
	}
	return n, nil
}

// Ident implements ExpComponent.Ident
func (c *SubscriptComponent) Ident() string {
	return "[" + c.Value.String() + "]"
}

// =========================== ArrayComponent ===========================

// ArrayComponent is an expression component for numbers
type ArrayComponent struct {
	componentChainNode

	Values []Node
}

// Ident implements ExpComponent.Ident
func (c *ArrayComponent) Ident() string {
	buf := &bytes.Buffer{}

	buf.WriteString("[")
	for i, arg := range c.Values {
		if arg == nil {
			buf.WriteString("_")
		} else {
			fmt.Fprintf(buf, "%s", arg)
		}

		if i < len(c.Values)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString("]")
	return buf.String()
}

// ConstructNode returns the ast node for the expression component
func (c *ArrayComponent) ConstructNode(prev Node) (Node, error) {
	n := ArrayNode{}
	n.Token = c.token

	n.Length = len(c.Values)
	n.NodeType = nodeArray

	for _, item := range c.Values {
		n.Elements = append(n.Elements, item)
	}
	return n, nil
}

// =========================== DotComponent ===========================

// DotComponent is an expression component for numbers
type DotComponent struct {
	componentChainNode

	Value string
}

// Ident implements ExpComponent.Ident
func (c *DotComponent) Ident() string {
	return fmt.Sprintf(".%s", c.Value)
}

// ConstructNode returns the ast node for the expression component
func (c *DotComponent) ConstructNode(prev Node) (Node, error) {
	n := DotReference{}
	n.Token = c.token
	n.NodeType = nodeDot

	base, ok := prev.(Reference)
	if !ok {
		return nil, fmt.Errorf("dot component requires a reference type on the lhs. instead got %T", prev)
	}
	n.Base = base
	n.Field = NewIdentNode(c.Value)
	return n, nil
}

// =========================== StringComponent ===========================

// StringComponent is an expression component for numbers
type StringComponent struct {
	componentChainNode

	Value string
}

// Ident implements ExpComponent.Ident
func (c *StringComponent) Ident() string {
	return fmt.Sprintf(".%s", c.Value)
}

// ConstructNode returns the ast node for the expression component
func (c *StringComponent) ConstructNode(prev Node) (Node, error) {
	n := StringNode{}
	n.Token = c.token
	n.NodeType = nodeString
	val := c.Value[1 : len(c.Value)-1]
	escaped, _ := UnescapeString(val)
	n.Value = escaped
	return n, nil
}

// =========================== ParenthesisComponent ===========================

// ParenthesisComponent is an expression component for numbers
type ParenthesisComponent struct {
	componentChainNode

	Value Node
}

// Ident implements ExpComponent.Ident
func (c *ParenthesisComponent) Ident() string {
	return fmt.Sprintf(".%s", c.Value)
}

// ConstructNode returns the ast node for the expression component
func (c *ParenthesisComponent) ConstructNode(prev Node) (Node, error) {
	return c.Value, nil
}

// =========================== BooleanComponent ===========================

// BooleanComponent is an expression component for numbers
type BooleanComponent struct {
	componentChainNode

	Value string
}

// Ident implements ExpComponent.Ident
func (c *BooleanComponent) Ident() string {
	return fmt.Sprintf("%s", c.Value)
}

// ConstructNode returns the ast node for the expression component
func (c *BooleanComponent) ConstructNode(prev Node) (Node, error) {

	n := BooleanNode{}
	n.Token = c.token
	n.NodeType = nodeBool

	n.Value = c.Value

	return n, nil
}

// =========================== CharComponent ===========================

// CharComponent is an expression component for numbers
type CharComponent struct {
	componentChainNode

	Value string
}

// Ident implements ExpComponent.Ident
func (c *CharComponent) Ident() string {
	return fmt.Sprintf("%s", c.Value)
}

// ConstructNode returns the ast node for the expression component
func (c *CharComponent) ConstructNode(prev Node) (Node, error) {

	n := CharNode{}
	n.Token = c.token
	n.NodeType = nodeBool

	n.Value = []rune(c.Value)[1]

	return n, nil
}

// =========================== TypeInfoComponent ===========================

// TypeInfoComponent is an expression component for numbers
type TypeInfoComponent struct {
	componentChainNode

	Type TypeNode
}

// Ident implements ExpComponent.Ident
func (c *TypeInfoComponent) Ident() string {
	node, _ := c.ConstructNode(nil)
	return fmt.Sprintf("%s", node)
}

// ConstructNode returns the ast node for the expression component
func (c *TypeInfoComponent) ConstructNode(prev Node) (Node, error) {
	n := TypeInfoNode{}
	n.Token = c.token
	n.NodeType = nodeTypeInfo
	n.T = c.Type
	return n, nil
}
