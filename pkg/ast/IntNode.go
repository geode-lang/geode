package ast

import (
	"strconv"

	"github.com/geode-lang/geode/llvm/ir/constant"
	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
)

// IntNode is an integer literal
type IntNode struct {
	NodeType
	TokenReference
	Accessable

	Value int64
}

// NameString implements Node.NameString
func (n IntNode) NameString() string { return "IntNode" }

// Codegen implements Node.Codegen for IntNode
func (n IntNode) Codegen(prog *Program) (value.Value, error) {
	// return llvm.ConstInt(llvm.Int64Type(), , true)
	return constant.NewInt(n.Value, types.I64), nil
}

func (n IntNode) String() string {
	return strconv.FormatInt(n.Value, 10)
}

// GenAccess implements Accessable.GenAccess
func (n IntNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}
