package ast

import (
	"strconv"

	"github.com/geode-lang/llvm/ir/constant"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
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

// InferType implements Node.InferType
func (n IntNode) InferType(scope *Scope) string { return "int" }

// Codegen implements Node.Codegen for IntNode
func (n IntNode) Codegen(prog *Program) value.Value {
	// return llvm.ConstInt(llvm.Int64Type(), , true)
	return constant.NewInt(n.Value, types.I64)
}

func (n IntNode) String() string {
	return strconv.FormatInt(n.Value, 10)
}

// GenAccess implements Accessable.GenAccess
func (n IntNode) GenAccess(prog *Program) value.Value {
	return n.Codegen(prog)
}
