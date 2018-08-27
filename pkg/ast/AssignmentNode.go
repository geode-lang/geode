package ast

import (
	"fmt"

	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// AssignmentNode is a node that has an assignable and a value.
// the Codegen function on the node will generate an assign and an
// access and update scopes in accordance
type AssignmentNode struct {
	NodeType
	TokenReference

	Assignee Assignable
	Value    Accessable
}

// NameString implements Node.NameString
func (n AssignmentNode) NameString() string { return "AssignmentNode" }

// InferType implements Node.InferType
func (n AssignmentNode) InferType(scope *Scope) string {
	return "void"
}

func (n AssignmentNode) String() string {
	return fmt.Sprintf("%s <- %s", n.Assignee, n.Value)
}

// GenAccess returns the value of the assignment
func (n AssignmentNode) GenAccess(prog *Program) value.Value {
	return n.Codegen(prog)
}

// Codegen implements Node.Codegen for AssignmentNode
func (n AssignmentNode) Codegen(prog *Program) value.Value {

	prog.Compiler.CurrentBlock().AppendInst(NewLLVMComment(n.String()))
	targetType := n.Assignee.Type(prog)
	prog.Compiler.typeCache = targetType

	val := n.Value.GenAccess(prog)

	if !types.Equal(val.Type(), targetType) {
		val = createTypeCast(prog, val, targetType)
	}

	n.Assignee.GenAssign(prog, val)
	return val
}
