package ast

import (
	"fmt"

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

// Codegen implements Node.Codegen for AssignmentNode
func (n AssignmentNode) Codegen(s *Scope, c *Compiler) value.Value {
	c.CurrentBlock().AppendInst(NewLLVMComment(n.String()))
	val := n.Value.GenAccess(s, c)
	n.Assignee.GenAssign(s, c, val)
	return val
}
