package ast

import (
	"github.com/llir/llvm/ir/value"
)

// BlockNode is a block statement. A block statement is just an array of Nodes
// that run in sequence.
type BlockNode struct {
	NodeType
	Nodes []Node
}

// NameString implements Node.NameString
func (n BlockNode) NameString() string { return "BlockNode" }

// InferType implements Node.InferType
func (n BlockNode) InferType(scope *Scope) string { return "void" }

// Codegen implements Node.Codegen for BlockNode
func (n BlockNode) Codegen(scope *Scope, c *Compiler) value.Value {
	blockScope := scope.SpawnChild()
	// spew.Dump(n.Nodes)
	for _, node := range n.Nodes {
		node.Codegen(blockScope, c)
	}
	return c.CurrentBlock()
}
