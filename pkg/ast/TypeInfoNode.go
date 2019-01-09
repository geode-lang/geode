package ast

import (
	"fmt"

	"github.com/geode-lang/geode/pkg/gtypes"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// TypeInfoDeclaration is what will be used in the program to keep track of globals
type TypeInfoDeclaration struct {
	Global  *ir.Global
	Defined bool
}

// TypeInfoNode returns the size of a type passed into it as a function
type TypeInfoNode struct {
	NodeType
	TokenReference

	T TypeNode
}

// NameString implements Node.NameString
func (n TypeInfoNode) NameString() string { return "TypeInfoNode" }

// Codegen implements Node.Codegen for TypeInfoNode
func (n TypeInfoNode) Codegen(prog *Program) (value.Value, error) {
	found, ok := prog.TypeInfoDefs[n.T.String()]
	if ok && found.Defined {
		return found.Global, nil
	}

	analyzeType, err := n.T.GetType(prog)
	if err != nil {
		return nil, err
	}

	// allocation was not found, so we make a new global one.
	typ, _ := n.Type(prog)

	sct := typ.(*gtypes.StructType)
	globl := prog.Module.NewGlobal(fmt.Sprintf("type_info_%s", n.T), sct)

	globl.Init = constant.NewZeroInitializer(sct)

	prog.TypeInfoDefs[n.T.String()] = &TypeInfoDeclaration{
		Global:  globl,
		Defined: false,
	}

	// https://stackoverflow.com/a/30830445
	elemptr := constant.NewGetElementPtr(constant.NewNull(types.NewPointer(analyzeType)), constant.NewInt(types.I32, 1))

	size := prog.Compiler.CurrentBlock().NewPtrToInt(elemptr, types.I64)

	nameNode := StringNode{}
	nameNode.Value = n.T.Name
	name, _ := nameNode.Codegen(prog)

	inst := NewClassInstance(prog, sct, map[string]value.Value{
		"size": size,
		"name": name,
	})

	prog.Compiler.CurrentBlock().NewStore(inst, globl)

	prog.TypeInfoDefs[n.T.String()] = &TypeInfoDeclaration{
		Global:  globl,
		Defined: true,
	}

	return globl, nil
}

// GenAccess implements Accessable.Access for TypeInfoNode
func (n TypeInfoNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}

func (n TypeInfoNode) String() string {
	return fmt.Sprintf("info(%s)", n.T)
}

// Alloca implements Reference.Alloca
func (n TypeInfoNode) Alloca(prog *Program) value.Value {

	if found, ok := prog.TypeInfoDefs[n.T.String()]; ok {
		return found.Global
	}

	v, _ := n.Codegen(prog)
	return v
}

// Load implements Reference.Load
func (n TypeInfoNode) Load(blk *ir.Block, prog *Program) *ir.InstLoad {
	return blk.NewLoad(n.Alloca(prog))
}

// GenAssign implements Assignable.GenAssign
func (n TypeInfoNode) GenAssign(prog *Program, val value.Value, options ...AssignableOption) (value.Value, error) {
	return nil, fmt.Errorf("unable to assign to any part of a type info call")
}

// Type implements Assignable.Type
func (n TypeInfoNode) Type(prog *Program) (types.Type, error) {
	return prog.Scope.FindType("TypeInfo").Type, nil
}
