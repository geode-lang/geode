package ast

import (
	"fmt"

	"github.com/geode-lang/geode/llvm/ir"
	"github.com/geode-lang/geode/llvm/ir/constant"
	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
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
	baseType := base.Type()
	for types.IsPointer(baseType) {
		baseType = baseType.(*types.PointerType).Elem
	}
	return baseType
}

// BaseAddr returns the true address of the base, be it through loads, etc...
func (n DotReference) BaseAddr(prog *Program) value.Value {
	var val value.Value
	val = n.Base.Alloca(prog)
	for {
		load := ir.NewLoad(val)
		if types.IsPointer(load.Type()) {
			prog.Compiler.CurrentBlock().AppendInst(load)
			val = load
		} else {
			break
		}
	}
	return val
}

// GetFunc implemnets Callable.GetFunc
func (n DotReference) GetFunc(prog *Program, argTypes []types.Type) (*ir.Function, []value.Value, error) {

	class := n.BaseType(prog)

	name, err := prog.Scope.FindTypeName(class)
	if err != nil {
		return nil, nil, err
	}

	args := make([]value.Value, 0)

	args = append(args, n.BaseAddr(prog))

	fieldName := n.Field.String()

	// fmt.Println(funcName)

	argTypes = append([]types.Type{types.NewPointer(class)}, argTypes...)

	searchNames := []string{
		fmt.Sprintf("%s.%s", name, fieldName),
		fmt.Sprintf("runtime:%s.%s", name, fieldName),
	}

	// fmt.Println(searchNames)
	fn, err := prog.FindFunction(searchNames, argTypes)
	// fmt.Println(fn, err)

	// for k := range prog.Functions {
	// 	fmt.Println(k)
	// }

	return fn, args, err
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

// NameString implements Node.NameString
func (n DotReference) NameString() string { return "DotReference" }

// Codegen implements Node.Codegen
func (n DotReference) Codegen(prog *Program) (value.Value, error) {
	return n.GenAccess(prog)
}

// Load returns a load instruction on a named reference with the given name
func (n DotReference) Load(block *ir.BasicBlock, prog *Program) *ir.InstLoad {
	target := n.Alloca(prog).(*ir.InstGetElementPtr)
	t, _ := n.Type(prog)
	target.Typ = types.NewPointer(t)
	return block.NewLoad(target)
}

// GenAssign implements Assignable.GenAssign
func (n DotReference) GenAssign(prog *Program, assignment value.Value, options ...AssignableOption) (value.Value, error) {
	target := n.Alloca(prog)
	prog.Compiler.CurrentBlock().NewStore(assignment, target)
	return assignment, nil
}

// GenAccess implements Accessable.GenAccess
func (n DotReference) GenAccess(prog *Program) (value.Value, error) {
	return n.Load(prog.Compiler.CurrentBlock(), prog), nil
}

// Type implements Assignable.Type
func (n DotReference) Type(prog *Program) (types.Type, error) {
	baseType := n.BaseType(prog).(*types.StructType)
	index := baseType.FieldIndex(n.Field.String())
	return baseType.Fields[index], nil
}
