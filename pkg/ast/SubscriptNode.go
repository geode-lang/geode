package ast

import (
	"fmt"

	"github.com/geode-lang/geode/llvm/ir"
	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
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

func (n SubscriptNode) String() string {
	return fmt.Sprintf("%s[%s]", n.Source, n.Index)
}

// GenElementPtr returns a generated GetElementPtr for this subscript operation
func (n SubscriptNode) GenElementPtr(prog *Program) (*ir.InstGetElementPtr, error) {
	src, err := n.Source.GenAccess(prog)
	if err != nil {
		return nil, err
	}
	idx, err := n.Index.GenAccess(prog)
	if err != nil {
		return nil, err
	}
	return prog.Compiler.CurrentBlock().NewGetElementPtr(src, idx), nil
}

// Codegen implements Node.Codegen for SubscriptNode
func (n SubscriptNode) Codegen(prog *Program) (value.Value, error) {
	// c.CurrentBlock().AppendInst(NewLLVMComment("%s", n))
	ptr, err := n.GenElementPtr(prog)
	if err != nil {
		return nil, err
	}
	return prog.Compiler.CurrentBlock().NewLoad(ptr), nil
}

// GenAccess implements Accessable.GenAccess
func (n SubscriptNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}

// GenAssign generates an assignment at the address
func (n SubscriptNode) GenAssign(prog *Program, val value.Value) (value.Value, error) {
	ptr, err := n.GenElementPtr(prog)
	if err != nil {
		return nil, err
	}
	prog.Compiler.CurrentBlock().NewStore(val, ptr)
	return val, nil
}

// Type returns the type of the node.
func (n SubscriptNode) Type(prog *Program) (types.Type, error) {

	tmpBlock := ir.NewBlock("")

	ptr, err := n.GenElementPtr(prog)
	if err != nil {
		return nil, err
	}

	load := tmpBlock.NewLoad(ptr)
	return load.Type(), nil
}
