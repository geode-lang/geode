package ast

import (
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// CastNode is a structure around a typecast expression
type CastNode struct {
	NodeType
	TokenReference

	Source Node
	Type   GeodeTypeRef
}

// NameString implements Node.NameString
func (n CastNode) NameString() string { return "CastNode" }

// InferType implements Node.InferType
func (n CastNode) InferType(scope *Scope) string {
	return n.Type.Name
}

// GenAccess implements Accessable.Access for CastNode
func (n CastNode) GenAccess(prog *Program) value.Value {
	return n.Codegen(prog)
}

// Codegen implements Node.Codegen for CastNode
func (n CastNode) Codegen(prog *Program) value.Value {

	scope := prog.Scope
	src := n.Source.Codegen(prog)
	t := scope.FindType(n.Type.Name).Type
	for i := 0; i < n.Type.PointerLevel; i++ {
		t = types.NewPointer(t)
	}
	return createTypeCast(prog, src, t)
}
