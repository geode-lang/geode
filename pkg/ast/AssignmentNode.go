package ast

import (
	"fmt"

	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
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

func (n AssignmentNode) String() string {
	return fmt.Sprintf("%s <- %s", n.Assignee, n.Value)
}

// GenAccess returns the value of the assignment
func (n AssignmentNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}

// Codegen implements Node.Codegen for AssignmentNode
func (n AssignmentNode) Codegen(prog *Program) (value.Value, error) {
	var err error
	targetType, _ := n.Assignee.Type(prog)
	prog.Compiler.PushType(targetType)

	val, err := n.Value.GenAccess(prog)
	if err != nil {
		return nil, err
	}

	if targetType != nil && !types.Equal(val.Type(), targetType) {
		val, err = createTypeCast(prog, val, targetType)
		if err != nil {
			return nil, err
		}
	}

	// fmt.Println(val)
	n.Assignee.GenAssign(prog, val)
	return val, nil
}
