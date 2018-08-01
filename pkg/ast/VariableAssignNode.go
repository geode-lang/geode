package ast

import (
	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/value"
)

// VariableAssignNode is a generic variable statement representation
// this contains a reference type inside it that tellst the
// code generator what kind of variable statement to build
type VariableAssignNode struct {
	NodeType
	TokenReference

	Type     GeodeTypeRef
	HasValue bool
	Target   *NamedReference
	Body     Node
}

// NameString implements Node.NameString
func (n VariableAssignNode) NameString() string { return "VariableAssignNode" }

// InferType implements Node.InferType
func (n VariableAssignNode) InferType(scope *Scope) string {
	return n.Type.Name
}

// Codegen implements Node.Codegen for VariableAssignNode
func (n VariableAssignNode) Codegen(scope *Scope, c *Compiler) value.Value {

	block := c.CurrentBlock()

	name := n.Target
	var alloc *ir.InstAlloca
	var val value.Value

	v, found := scope.Find(name.String())
	if !found {
		log.Fatal("%s not found\n", v)
	}
	alloc = v.Value().(*ir.InstAlloca)

	// Construct the body
	if n.Body != nil {
		val = n.Body.Codegen(scope, c)
		if val == nil {
			return val // nil
		}
	} else {
		log.Fatal("Missing body on variable assignment to name '%s'\n", name)
	}
	val = createTypeCast(c, val, alloc.Elem)
	block.NewStore(val, alloc)
	return nil
}
