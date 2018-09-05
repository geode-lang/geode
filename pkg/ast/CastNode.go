package ast

import (
	"fmt"

	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
)

// CastNode is a structure around a typecast expression
type CastNode struct {
	NodeType
	TokenReference

	Source Node
	Type   GeodeTypeRef
}

// NameString implements Node.NameString
func (n CastNode) NameString() string { return "CastNode" }

// GenAccess implements Accessable.Access for CastNode
func (n CastNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}

// Codegen implements Node.Codegen for CastNode
func (n CastNode) Codegen(prog *Program) (value.Value, error) {

	src, err := n.Source.Codegen(prog)
	if err != nil {
		return nil, err
	}
	t, err := prog.FindType(n.Type.Name)
	if err != nil {
		return nil, err
	}
	for i := 0; i < n.Type.PointerLevel; i++ {
		t = types.NewPointer(t)
	}
	return createTypeCast(prog, src, t)
}

func (n CastNode) String() string {
	return fmt.Sprintf("%s as %s", n.Source, n.Type)
}
