package ast

import (
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
	Name         Reference
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
func (n VariableNode) Codegen(prog *Program) value.Value {

	block := prog.Compiler.CurrentBlock()

	switch n.RefType {
	case ReferenceDereference, ReferenceAccessStackAddress:
		alloc := n.Name.Alloca(prog)
		return alloc
	case ReferenceAccessValue:
		val := n.Name.Load(block, prog)
		return val
	}

	return nil
}

// GenAddress returns the instruction allocation
func (n VariableNode) GenAddress(prog *Program) value.Value {
	return n.Name.Alloca(prog)
}

// GenAccess returns the value of a VariableNode
func (n VariableNode) GenAccess(prog *Program) value.Value {
	return n.Codegen(prog)
}
