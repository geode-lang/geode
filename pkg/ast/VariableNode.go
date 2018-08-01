package ast

import (
	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/value"
)

// VariableNode is a generic variable statement representation
// this contains a reference type inside it that tellst the
// code generator what kind of variable statement to build
type VariableNode struct {
	NodeType
	TokenReference

	Type         GeodeTypeRef
	HasValue     bool
	Name         *NamedReference
	IsPointer    bool
	RefType      ReferenceType
	IsArray      bool
	Reassignment bool
	Body         Accessable
}

// NameString implements Node.NameString
func (n VariableNode) NameString() string { return "VariableNode" }

// InferType implements Node.InferType
func (n VariableNode) InferType(scope *Scope) string {

	if n.RefType == ReferenceDefine {
		found := scope.FindType(n.Type.Name)
		if found == nil {
			return "void"
		}

		return found.Name
	}

	return "void"

}

func (n VariableNode) String() string {
	return n.Name.String()
}

// Codegen implements Node.Codegen for VariableNode
func (n VariableNode) Codegen(scope *Scope, c *Compiler) value.Value {

	block := c.CurrentBlock()

	switch n.RefType {
	case ReferenceDereference, ReferenceAccessStackAddress:
		alloc := n.Name.Alloca(scope)
		return alloc
	case ReferenceAccessValue:
		val := n.Name.Load(block, scope)
		return val
	}

	return nil
}

// GenAddress returns the instruction allocation
func (n VariableNode) GenAddress(s *Scope, c *Compiler) *ir.InstAlloca {
	return n.Name.Alloca(s)
}

func (n VariableNode) GenAccess(s *Scope, c *Compiler) value.Value {
	return n.Codegen(s, c)
}
