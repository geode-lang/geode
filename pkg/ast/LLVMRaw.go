package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/metadata"
	"github.com/llir/llvm/ir/types"
)

// LLVMRaw implements ir.Instruction
// This will be easier than doing individual work
// per operation
type LLVMRaw struct {
	// Parent basic block.
	Parent *ir.BasicBlock

	data string

	typ types.Type

	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewLLVMRaw returns a new geode binary instruction
func NewLLVMRaw(content string) *LLVMRaw {
	return &LLVMRaw{
		data:     content,
		typ:      types.Void,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

// Type returns the type of the instruction.
func (inst *LLVMRaw) Type() types.Type {
	return inst.typ
}

// Ident returns the identifier associated with the instruction.
func (inst *LLVMRaw) Ident() string {
	return inst.data
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *LLVMRaw) GetName() string {
	return inst.data
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *LLVMRaw) SetName(name string) {

}

// String returns the LLVM syntax representation of the instruction.
func (inst *LLVMRaw) String() string {
	return inst.data
}

// GetParent returns the parent basic block of the instruction.
func (inst *LLVMRaw) GetParent() *ir.BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *LLVMRaw) SetParent(parent *ir.BasicBlock) {
	inst.Parent = parent
}
