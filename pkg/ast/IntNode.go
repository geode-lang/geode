package ast

import (
	"strconv"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
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
	return constant.NewInt(types.I64, n.Value), nil
}

func (n IntNode) String() string {
	return strconv.FormatInt(n.Value, 10)
}

// GenAccess implements Accessable.GenAccess
func (n IntNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}
