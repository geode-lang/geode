package ast

import (
	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// NameType is a type to notate what kind of name a NamedReference is
type NameType int

// Various NameType constants
const (
	ClassMethodNameType NameType = iota
)

// NamedReference is a reference to an item on the scope via some name
// this can be used to access variable allocs, function defns, or types
type NamedReference struct {
	Value    string
	NameType NameType
}

// NewNamedReference returns a new name reference with a string as it's name
func NewNamedReference(name string) NamedReference {
	n := NamedReference{}

	n.Value = name
	return n
}

func (n NamedReference) String() string {
	return n.Value
}

// Alloca returns the nearest alloca instruction in this scope with the given name
func (n NamedReference) Alloca(prog *Program) value.Value {

	scopeitem, found := prog.Scope.Find(n.Value)
	if !found {
		log.Fatal("Unable to find named reference %s\n", n)
	}

	alloc, cast := scopeitem.(VariableScopeItem).Value().(*ir.InstAlloca)
	if !cast {
		log.Fatal("Cast from scope item to alloca with named reference %s failed\n", n)
	}

	return alloc
}

// Load returns a load instruction on a named reference with the given name
func (n NamedReference) Load(block *ir.BasicBlock, prog *Program) *ir.InstLoad {
	return block.NewLoad(n.Alloca(prog))
}

// GenAssign implements Assignable.GenAssign
func (n NamedReference) GenAssign(prog *Program, assignment value.Value) value.Value {
	prog.Compiler.CurrentBlock().NewStore(assignment, n.Alloca(prog))
	return assignment
}

// GenAccess implements Accessable.GenAccess
func (n NamedReference) GenAccess(prog *Program) value.Value {
	return n.Load(prog.Compiler.CurrentBlock(), prog)
}

// Type implements Assignable.Type
func (n NamedReference) Type(prog *Program) types.Type {
	return n.Alloca(prog).(*ir.InstAlloca).Elem
}
