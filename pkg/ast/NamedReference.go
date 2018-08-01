package ast

import (
	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/value"
)

// NamedReference is a reference to an item on the scope via some name
// this can be used to access variable allocs, function defns, or types
type NamedReference struct {
	Value string
}

// NewNamedReference returns a new name reference with a string as it's name
func NewNamedReference(name string) *NamedReference {
	n := &NamedReference{}

	n.Value = name
	return n
}

func (n *NamedReference) String() string {
	return n.Value
}

// Alloca returns the nearest alloca instruction in this scope with the given name
func (n *NamedReference) Alloca(s *Scope) *ir.InstAlloca {
	scopeitem, found := s.Find(n.Value)
	if !found {
		log.Fatal("Unable to find named reference %s\n", n)
	}

	alloc, cast := scopeitem.Value().(*ir.InstAlloca)
	if !cast {
		log.Fatal("Cast from scope item to alloca with named reference %s failed\n", n)
	}

	return alloc
}

// Load returns a load instruction on a named reference with the given name
func (n *NamedReference) Load(block *ir.BasicBlock, s *Scope) *ir.InstLoad {
	return block.NewLoad(n.Alloca(s))
}

// GenAssign implements Assignable.GenAssign
func (n *NamedReference) GenAssign(s *Scope, c *Compiler, assignment value.Value) value.Value {
	c.CurrentBlock().NewStore(assignment, n.Alloca(s))
	return assignment
}

// GenAccess implements Accessable.GenAccess
func (n *NamedReference) GenAccess(s *Scope, c *Compiler) value.Value {
	return n.Load(c.CurrentBlock(), s)
}
