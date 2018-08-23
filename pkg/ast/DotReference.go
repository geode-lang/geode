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
func (n DotReference) BaseType(prog *Program) types.Type {
	base := n.Base.Alloca(prog)
	baseType := base.(*ir.InstAlloca).Elem
	for types.IsPointer(baseType) {
		baseType = baseType.(*types.PointerType).Elem
	}
	return baseType
}

// Alloca returns the nearest alloca instruction in this scope with the given name
func (n DotReference) Alloca(prog *Program) value.Value {
	base := n.Base.Alloca(prog)
	index := 0
	baseType := n.BaseType(prog)

	// An allocation is always a pointer, so we need to figure out what it is pointing to
	// here, I coerce base's type into a *PointerType and pull the Elem type out of it.
	// This finds this part of the InstAlloca:
	//      %_1 = alloca i8
	//                   ^^
	ptr := base.Type().(*types.PointerType)
	elemType := ptr.Elem

	// If the type that the alloca points to is a pointer, we need to load from the pointer
	if types.IsPointer(elemType) {
		base = prog.Compiler.CurrentBlock().NewLoad(base)

	}

	structType := baseType.(*types.StructType)
	index = structType.FieldIndex(n.Field.String())

	zero := constant.NewInt(0, types.I32)
	fieldOffset := constant.NewInt(int64(index), types.I32)
	gen := prog.Compiler.CurrentBlock().NewGetElementPtr(base, zero, fieldOffset)

	return gen
}

// Load returns a load instruction on a named reference with the given name
func (n DotReference) Load(block *ir.BasicBlock, prog *Program) *ir.InstLoad {
	target := n.Alloca(prog).(*ir.InstGetElementPtr)
	target.Typ = types.NewPointer(n.Type(prog))
	return block.NewLoad(target)
}

// GenAssign implements Assignable.GenAssign
func (n DotReference) GenAssign(prog *Program, assignment value.Value) value.Value {
	target := n.Alloca(prog)
	prog.Compiler.CurrentBlock().NewStore(assignment, target)
	return assignment
}

// GenAccess implements Accessable.GenAccess
func (n DotReference) GenAccess(prog *Program) value.Value {
	return n.Load(prog.Compiler.CurrentBlock(), prog)
}

// Type implements Assignable.Type
func (n DotReference) Type(prog *Program) types.Type {
	baseType := n.BaseType(prog).(*types.StructType)
	index := baseType.FieldIndex(n.Field.String())
	return baseType.Fields[index]
}
