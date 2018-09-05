package ast

import (
	"fmt"

	"github.com/geode-lang/geode/llvm/ir/constant"
	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
)

// SizeofNode returns the size of a type passed into it as a function
type SizeofNode struct {
	NodeType
	TokenReference

	Type GeodeTypeRef
}

// NameString implements Node.NameString
func (n SizeofNode) NameString() string { return "SizeofNode" }

// Codegen implements Node.Codegen for SizeofNode
func (n SizeofNode) Codegen(prog *Program) (value.Value, error) {
	t, err := prog.FindType(n.Type.Name)
	if err != nil {
		return nil, err
	}
	for i := 0; i < n.Type.PointerLevel; i++ {
		t = types.NewPointer(t)
	}

	// prog.Compiler.CurrentBlock().AppendInst(NewLLVMComment("sizeof(%s)", t))
	// https://stackoverflow.com/a/30830445
	elemptr := constant.NewGetElementPtr(constant.NewNull(types.NewPointer(t)), constant.NewInt(1, types.I32))
	return prog.Compiler.CurrentBlock().NewPtrToInt(elemptr, types.I64), nil
}

// GenAccess implements Accessable.Access for SizeofNode
func (n SizeofNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}

func (n SizeofNode) String() string {
	return fmt.Sprintf("sizeof(%s)", n.Type)
}
