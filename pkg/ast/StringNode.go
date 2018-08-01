package ast

import (
	"github.com/geode-lang/llvm/ir/constant"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// StringNode -
type StringNode struct {
	NodeType
	TokenReference

	Value string
}

// NameString implements Node.NameString
func (n StringNode) NameString() string { return "StringNode" }

// InferType implements Node.InferType
func (n StringNode) InferType(scope *Scope) string { return "string" }

// Codegen implements Node.Codegen for StringNode
func (n StringNode) Codegen(scope *Scope, c *Compiler) value.Value {
	str := c.Module.NewGlobalDef(mangleName(".str"), newCharArray(n.Value))
	// str.IsConst = true
	zero := constant.NewInt(0, types.I32)
	return constant.NewGetElementPtr(str, zero, zero)
}

// GenAccess implements Accessable.GenAccess
func (n StringNode) GenAccess(s *Scope, c *Compiler) value.Value {
	return n.Codegen(s, c)
}
