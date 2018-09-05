package ast

import (
	"fmt"

	"github.com/geode-lang/geode/llvm/ir/constant"
	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
)

// NilNode -
type NilNode struct {
	NodeType
	TokenReference
}

// NameString implements Node.NameString
func (n NilNode) NameString() string { return "NilNode" }

// Codegen implements Node.Codegen for NilNode
func (n NilNode) Codegen(prog *Program) (value.Value, error) {
	return constant.NewNull(types.NewPointer(types.I8)), nil
}

// GenAccess implements Accessable.GenAccess
func (n NilNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}

func (n NilNode) String() string {
	return fmt.Sprintf("nil")
}
