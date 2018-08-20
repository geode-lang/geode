package ast

import (
	"fmt"

	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// SubscriptNode is a recursive subscript operation
type SubscriptNode struct {
	NodeType
	TokenReference

	Source Accessable // Source is the value being indexed
	Index  Accessable // Index is the offset to index by (must be full integer)
}

// NameString implements Node.NameString
func (n SubscriptNode) NameString() string { return "SubscriptNode" }

// InferType implements Node.InferType
func (n SubscriptNode) InferType(scope *Scope) string { return "void" }

func (n SubscriptNode) String() string {
	return fmt.Sprintf("%s[%s]", n.Source, n.Index)
}

// GenElementPtr returns a generated GetElementPtr for this subscript operation
func (n SubscriptNode) GenElementPtr(prog *Program) *ir.InstGetElementPtr {
	src := n.Source.GenAccess(prog)
	idx := n.Index.GenAccess(prog)
	return prog.Compiler.CurrentBlock().NewGetElementPtr(src, idx)
}

// Codegen implements Node.Codegen for SubscriptNode
func (n SubscriptNode) Codegen(prog *Program) value.Value {
	// c.CurrentBlock().AppendInst(NewLLVMComment("%s", n))
	return prog.Compiler.CurrentBlock().NewLoad(n.GenElementPtr(prog))
}

// GenAccess implements Accessable.GenAccess
func (n SubscriptNode) GenAccess(prog *Program) value.Value {
	return n.Codegen(prog)
}

// GenAssign generates an assignment at the address
func (n SubscriptNode) GenAssign(prog *Program, val value.Value) value.Value {
	prog.Compiler.CurrentBlock().NewStore(val, n.GenElementPtr(prog))
	return val
}

// Type returns the type of the node.
func (n SubscriptNode) Type(prog *Program) types.Type {

	tmpBlock := ir.NewBlock("")

	load := tmpBlock.NewLoad(n.GenElementPtr(prog))
	return load.Type()
}
