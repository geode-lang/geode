package ast

import (
	"fmt"

	"github.com/geode-lang/llvm/ir/constant"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// NilNode -
type NilNode struct {
	NodeType
	TokenReference
}

// NameString implements Node.NameString
func (n NilNode) NameString() string { return "NilNode" }

// InferType implements Node.InferType
func (n NilNode) InferType(scope *Scope) string { return "string" }

// Codegen implements Node.Codegen for NilNode
func (n NilNode) Codegen(prog *Program) value.Value {
	return constant.NewNull(types.NewPointer(types.I8))
}

// GenAccess implements Accessable.GenAccess
func (n NilNode) GenAccess(prog *Program) value.Value {
	return n.Codegen(prog)
}

func (n NilNode) String() string {
	return fmt.Sprintf("nil")
}
