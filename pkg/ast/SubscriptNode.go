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
func (n SubscriptNode) GenElementPtr(s *Scope, c *Compiler) *ir.InstGetElementPtr {
	src := n.Source.GenAccess(s, c)
	idx := n.Index.GenAccess(s, c)
	return c.CurrentBlock().NewGetElementPtr(src, idx)
}

// Codegen implements Node.Codegen for SubscriptNode
func (n SubscriptNode) Codegen(scope *Scope, c *Compiler) value.Value {
	// c.CurrentBlock().AppendInst(NewLLVMComment("%s", n))
	return c.CurrentBlock().NewLoad(n.GenElementPtr(scope, c))
}

// GenAccess implements Accessable.GenAccess
func (n SubscriptNode) GenAccess(s *Scope, c *Compiler) value.Value {
	return n.Codegen(s, c)
}

// GenAssign generates an assignment at the address
func (n SubscriptNode) GenAssign(s *Scope, c *Compiler, val value.Value) value.Value {
	c.CurrentBlock().NewStore(val, n.GenElementPtr(s, c))
	return val
}

// Type returns the type of the node.
func (n SubscriptNode) Type(s *Scope, c *Compiler) types.Type {

	tmpBlock := ir.NewBlock("")

	tmpC := NewCompiler(ir.NewModule(), "tmp", NewPackage("", nil, s, nil), nil)
	tmpC.PushBlock(tmpBlock)

	load := tmpBlock.NewLoad(n.GenElementPtr(s, c))
	return load.Type()
}
