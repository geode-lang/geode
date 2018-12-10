package ast

import (
	"fmt"

	"github.com/geode-lang/geode/pkg/gtypes"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// SubscriptNode is a recursive subscript operation
type SubscriptNode struct {
	NodeType
	TokenReference
	Reference

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

	if gtypes.IsSlice(src.Type()) {
		fmt.Println(src.Type())
		zero := constant.NewInt(types.I64, 0)
		curBlock := prog.Compiler.CurrentBlock()
		inst := gep(src, zero)
		src = inst
		curBlock.Insts = append(curBlock.Insts, inst)
	}
	curBlock := prog.Compiler.CurrentBlock()
	inst := gep(src, idx)
	curBlock.Insts = append(curBlock.Insts, inst)
	return inst, nil
}

// Codegen implements Node.Codegen for SubscriptNode
func (n SubscriptNode) Codegen(prog *Program) (value.Value, error) {
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
func (n SubscriptNode) GenAssign(prog *Program, val value.Value, options ...AssignableOption) (value.Value, error) {
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

// Alloca implements Reference.Alloca
func (n SubscriptNode) Alloca(prog *Program) value.Value {
	ptr, _ := n.GenElementPtr(prog)
	return ptr
}

// Load implements Reference.Load
func (n SubscriptNode) Load(blk *ir.BasicBlock, prog *Program) *ir.InstLoad {
	ld, _ := n.Codegen(prog)
	return ld.(*ir.InstLoad)
}
