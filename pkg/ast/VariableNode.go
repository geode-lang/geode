package ast

import (
	"fmt"
	"os"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// VariableNode is a generic variable statement representation
// this contains a reference type inside it that tellst the
// code generator what kind of variable statement to build
type VariableNode struct {
	NodeType
	TokenReference

	Type         GeodeTypeRef
	HasValue     bool
	Name         *NamedReference
	IsPointer    bool
	RefType      ReferenceType
	IndexExpr    Node
	IsArray      bool
	Reassignment bool
	Body         Node
}

// NameString implements Node.NameString
func (n VariableNode) NameString() string { return "VariableNode" }

// InferType implements Node.InferType
func (n VariableNode) InferType(scope *Scope) string {

	if n.RefType == ReferenceDefine {
		found := scope.FindType(n.Type.Name)
		if found == nil {
			return "void"
		}

		return found.Name
	}

	return "void"

}

// Codegen implements Node.Codegen for VariableNode
func (n VariableNode) Codegen(scope *Scope, c *Compiler) value.Value {

	block := c.CurrentBlock()
	// f := block.Parent

	name := n.Name
	var alloc *ir.InstAlloca
	var val value.Value

	// fmt.Printf("VARIABLE NODE %s\n", name)

	// fmt.Printf("%s -> %s\n", n.Name, n.InferType(scope))

	if n.RefType == ReferenceAccessValue || n.RefType == ReferenceAccessStackAddress {
		v, found := scope.Find(name.String())
		if !found {
			n.SyntaxError()

			// spew.Dump(n)
			fmt.Printf("unknown variable name `%s`\n", name)
			os.Exit(-1)
		}

		alloc = v.Value().(*ir.InstAlloca)

		if n.RefType == ReferenceAccessStackAddress {
			return alloc
		}

		if n.RefType == ReferenceDereference {
			return alloc
		}

		val = block.NewLoad(v.Value())

		if n.IndexExpr != nil {
			if types.IsPointer(val.Type()) {
				// zero := constant.NewInt(0, types.I32)
				index := n.IndexExpr.Codegen(scope, c)
				ptr := block.NewGetElementPtr(val, index)
				val = block.NewLoad(ptr)
			}

		}
		return val
	}

	return nil
}
