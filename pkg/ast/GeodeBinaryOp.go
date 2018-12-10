package ast

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

// BinaryInstruction is a Geode binary instruction.
type BinaryInstruction interface {
	ir.Instruction
	value.Named
}

// NewGeodeBinaryInstr returns a new binary instruction based on the given
// operation.
func NewGeodeBinaryInstr(op string, x, y value.Value) BinaryInstruction {
	switch op {
	// Binary instruction.
	case "add":
		return ir.NewAdd(x, y)
	case "fadd":
		return ir.NewFAdd(x, y)
	case "sub":
		return ir.NewSub(x, y)
	case "fsub":
		return ir.NewFSub(x, y)
	case "mul":
		return ir.NewMul(x, y)
	case "fmul":
		return ir.NewFMul(x, y)
	case "udiv":
		return ir.NewUDiv(x, y)
	case "sdiv":
		return ir.NewSDiv(x, y)
	case "fdiv":
		return ir.NewFDiv(x, y)
	case "urem":
		return ir.NewURem(x, y)
	case "srem":
		return ir.NewSRem(x, y)
	case "frem":
		return ir.NewFRem(x, y)
	// Bitwise instructions.
	case "shl":
		return ir.NewShl(x, y)
	case "lshr":
		return ir.NewLShr(x, y)
	case "ashr":
		return ir.NewAShr(x, y)
	case "and":
		return ir.NewAnd(x, y)
	case "or":
		return ir.NewOr(x, y)
	case "xor":
		return ir.NewXor(x, y)
	default:
		panic(fmt.Errorf("support for binary instruction op %q not yet implemented", op))
	}
}
