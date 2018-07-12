package ast

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/metadata"
	"github.com/llir/llvm/ir/types"
)

// LLVMComment implements ir.Instruction
// This will be easier than doing individual work
// per operation
type LLVMComment struct {
	// Parent basic block.
	Parent *ir.BasicBlock

	data string

	typ types.Type

	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewLLVMComment returns a new geode binary instruction
func NewLLVMComment(typ types.Type, lines ...string) *LLVMComment {
	return &LLVMComment{
		data:     strings.Join(lines, "\n"),
		typ:      typ,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

// Type returns the type of the instruction.
func (inst *LLVMComment) Type() types.Type {
	return inst.typ
}

// Ident returns the identifier associated with the instruction.
func (inst *LLVMComment) Ident() string {
	return "COMMENT"
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *LLVMComment) GetName() string {
	return "COMMENT"
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *LLVMComment) SetName(name string) {

}

// String returns the LLVM syntax representation of the instruction.
func (inst *LLVMComment) String() string {
	return fmt.Sprintf("; %s", inst.data)
}

// GetParent returns the parent basic block of the instruction.
func (inst *LLVMComment) GetParent() *ir.BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *LLVMComment) SetParent(parent *ir.BasicBlock) {
	inst.Parent = parent
}
