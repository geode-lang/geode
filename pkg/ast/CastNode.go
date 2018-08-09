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
func (n CastNode) GenAccess(s *Scope, c *Compiler) value.Value {
	return n.Codegen(s, c)
}

// Codegen implements Node.Codegen for CastNode
func (n CastNode) Codegen(scope *Scope, c *Compiler) value.Value {
	src := n.Source.Codegen(scope, c)
	t := scope.FindType(n.Type.Name).Type
	for i := 0; i < n.Type.PointerLevel; i++ {
		t = types.NewPointer(t)
	}
	return createTypeCast(c, src, t)
}
