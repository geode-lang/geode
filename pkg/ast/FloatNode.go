package ast

import (
	"github.com/geode-lang/llvm/ir/constant"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// FloatNode is a float literla
type FloatNode struct {
	NodeType
	TokenReference

	Value float64
}

// NameString implements Node.NameString
func (n FloatNode) NameString() string { return "FloatNode" }

// InferType implements Node.InferType
func (n FloatNode) InferType(scope *Scope) string { return "float" }

// Codegen implements Node.Codegen for FloatNode
func (n FloatNode) Codegen(prog *Program) value.Value {
	return constant.NewFloat(n.Value, types.Double)
}

// GenAccess implements Accessable.GenAccess
func (n FloatNode) GenAccess(prog *Program) value.Value {
	return n.Codegen(prog)
}
