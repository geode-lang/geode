package ast

import (
	"strings"

	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/metadata"
	"github.com/geode-lang/llvm/ir/types"
)

// LLVMIdent implements ir.Instruction
// This will be easier than doing individual work
// per operation
type LLVMIdent struct {
	// Parent basic block.
	Parent *ir.BasicBlock

	data string

	typ types.Type

	// Map from metadata identifier (e.g. !dbg) to metadata associated with the
	// instruction.
	Metadata map[string]*metadata.Metadata
}

// NewLLVMIdent returns a new geode binary instruction
func NewLLVMIdent(typ types.Type, lines ...string) *LLVMIdent {
	return &LLVMIdent{
		data:     strings.Join(lines, "\n"),
		typ:      typ,
		Metadata: make(map[string]*metadata.Metadata),
	}
}

// Type returns the type of the instruction.
func (inst *LLVMIdent) Type() types.Type {
	return inst.typ
}

// Ident returns the identifier associated with the instruction.
func (inst *LLVMIdent) Ident() string {
	return inst.data
}

// GetName returns the name of the local variable associated with the
// instruction.
func (inst *LLVMIdent) GetName() string {
	return inst.data
}

// SetName sets the name of the local variable associated with the instruction.
func (inst *LLVMIdent) SetName(name string) {

}

// String returns the LLVM syntax representation of the instruction.
func (inst *LLVMIdent) String() string {
	return inst.data
}

// GetParent returns the parent basic block of the instruction.
func (inst *LLVMIdent) GetParent() *ir.BasicBlock {
	return inst.Parent
}

// SetParent sets the parent basic block of the instruction.
func (inst *LLVMIdent) SetParent(parent *ir.BasicBlock) {
	inst.Parent = parent
}
