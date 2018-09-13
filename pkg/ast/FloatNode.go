package ast

import (
	"fmt"

	"github.com/geode-lang/geode/llvm/ir/constant"
	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
)

// FloatNode is a float literla
type FloatNode struct {
	NodeType
	TokenReference

	Value float64
}

// NameString implements Node.NameString
func (n FloatNode) NameString() string { return "FloatNode" }

// Codegen implements Node.Codegen for FloatNode
func (n FloatNode) Codegen(prog *Program) (value.Value, error) {
	return constant.NewFloat(n.Value, types.Double), nil
}

// GenAccess implements Accessable.GenAccess
func (n FloatNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}

func (n FloatNode) String() string {
	return fmt.Sprintf("%f", n.Value)
}
