package ast

import (
	"fmt"

	"github.com/geode-lang/llvm/ir/constant"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// SizeofNode returns the size of a type passed into it as a function
type SizeofNode struct {
	NodeType
	TokenReference

	Type GeodeTypeRef
}

// NameString implements Node.NameString
func (n SizeofNode) NameString() string { return "SizeofNode" }

// InferType implements Node.InferType
func (n SizeofNode) InferType(scope *Scope) string {
	return n.Type.Name
}

// Codegen implements Node.Codegen for SizeofNode
func (n SizeofNode) Codegen(prog *Program) value.Value {
	t := prog.Scope.FindType(n.Type.Name).Type
	for i := 0; i < n.Type.PointerLevel; i++ {
		t = types.NewPointer(t)
	}

	prog.Compiler.CurrentBlock().AppendInst(NewLLVMComment("sizeof(%s)", t))
	// https://stackoverflow.com/a/30830445
	elemptr := constant.NewGetElementPtr(constant.NewNull(types.NewPointer(t)), constant.NewInt(1, types.I32))
	return prog.Compiler.CurrentBlock().NewPtrToInt(elemptr, types.I64)
}

// GenAccess implements Accessable.Access for SizeofNode
func (n SizeofNode) GenAccess(prog *Program) value.Value {
	return n.Codegen(prog)
}

func (n SizeofNode) String() string {
	return fmt.Sprintf("sizeof(%S)", n.Type)
}
