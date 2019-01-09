package ast

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
)

// LLVMComment is a Geode pseud-comment instruction. It implements the
// ir.Instruction interface by embedding a NOP add instruction.
type LLVMComment struct {
	// Comment; may contain multiple lines.
	data string

	// LLVMComment implements ir.Instruction by embedding a NOP add instruction.
	*ir.InstAdd
}

// NewLLVMComment returns a new Geode comment pseudo-instruction.
func NewLLVMComment(format string, args ...interface{}) *LLVMComment {
	zero := constant.NewInt(types.I64, 0)
	nop := ir.NewAdd(zero, zero)
	nop.SetName("nop")
	return &LLVMComment{
		data:    fmt.Sprintf(format, args...),
		InstAdd: nop,
	}
}

// LLString returns the LLVM syntax representation of the instruction.
func (inst *LLVMComment) LLString() string {
	// Handle multi-line comments.
	data := strings.Replace(inst.data, "\n", "; ", -1)
	return fmt.Sprintf("; %s", data)
}
