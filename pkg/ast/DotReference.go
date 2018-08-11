package ast

import (
	"fmt"

	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/constant"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// DotReference -
type DotReference struct {
	NodeType
	TokenReference
	Base  Reference
	Field Reference
}

func (n DotReference) String() string {
	return fmt.Sprintf("%s.%s", n.Base, n.Field)
}

// BaseType returns the type of the base struct to a class
func (n DotReference) BaseType(s *Scope, c *Compiler) types.Type {
	base := n.Base.Alloca(s, c)
	baseType := base.(*ir.InstAlloca).Elem
	for types.IsPointer(baseType) {
		baseType = baseType.(*types.PointerType).Elem
	}
	return baseType
}

// Alloca returns the nearest alloca instruction in this scope with the given name
func (n DotReference) Alloca(s *Scope, c *Compiler) value.Value {
	base := n.Base.Alloca(s, c)
	index := 0
	baseType := n.BaseType(s, c)

	// An allocation is always a pointer, so we need to figure out what it is pointing to
	// here, I coerce base's type into a *PointerType and pull the Elem type out of it.
	// This finds this part of the InstAlloca:
	//      %_1 = alloca i8
	//                   ^^
	ptr := base.Type().(*types.PointerType)
	elemType := ptr.Elem

	// If the type that the alloca points to is a pointer, we need to load from the pointer
	if types.IsPointer(elemType) {
		base = c.CurrentBlock().NewLoad(base)
	}
	structType := baseType.(*types.StructType)
	index = structType.FieldIndex(n.Field.String())

	zero := constant.NewInt(0, types.I32)
	fieldOffset := constant.NewInt(int64(index), types.I32)
	gen := c.CurrentBlock().NewGetElementPtr(base, zero, fieldOffset)

	return gen
}

// Load returns a load instruction on a named reference with the given name
func (n DotReference) Load(block *ir.BasicBlock, s *Scope, c *Compiler) *ir.InstLoad {
	target := n.Alloca(s, c).(*ir.InstGetElementPtr)
	target.Typ = types.NewPointer(n.Type(s, c))
	return block.NewLoad(target)
}

// GenAssign implements Assignable.GenAssign
func (n DotReference) GenAssign(s *Scope, c *Compiler, assignment value.Value) value.Value {
	target := n.Alloca(s, c)
	c.CurrentBlock().NewStore(assignment, target)
	return assignment
}

// GenAccess implements Accessable.GenAccess
func (n DotReference) GenAccess(s *Scope, c *Compiler) value.Value {
	return n.Load(c.CurrentBlock(), s, c)
}

// Type implements Assignable.Type
func (n DotReference) Type(s *Scope, c *Compiler) types.Type {
	baseType := n.BaseType(s, c).(*types.StructType)
	index := baseType.FieldIndex(n.Field.String())
	return baseType.Fields[index]
}
