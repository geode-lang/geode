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
	Type         GeodeTypeRef
	HasValue     bool
	Name         string
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
func (n VariableNode) InferType(scope *Scope) types.Type { return types.Void }

// Codegen implements Node.Codegen for VariableNode
func (n VariableNode) Codegen(scope *Scope, c *Compiler) value.Value {
	block := c.CurrentBlock()
	f := block.Parent

	name := n.Name
	var alloc *ir.InstAlloca
	var val value.Value

	if n.RefType == ReferenceAccessValue || n.RefType == ReferenceAccessStackAddress {
		v, found := scope.Find(name)
		if !found {
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

	if n.RefType == ReferenceAssign {
		v, found := scope.Find(name)
		if !found {
			fmt.Println(v, "Not found")
		}
		alloc = v.Value().(*ir.InstAlloca)
	} else if n.RefType == ReferenceDefine {
		ty := scope.FindType(n.Type.Name).Type
		ty = n.Type.BuildPointerType(ty)
		alloc = createBlockAlloca(f, ty, name)
		scItem := NewVariableScopeItem(n.Name, alloc, PrivateVisibility)
		scope.Add(scItem)
	}

	if n.HasValue {
		// Construct the body
		if n.Body != nil {
			val = n.Body.Codegen(scope, c)
			if val == nil {
				return val // nil
			}
		}
		val = createTypeCast(c, val, alloc.Elem)

	} else {
		// Default to 0 from issue:
		// https://github.com/nickwanninger/geode/issues/5
		// val = createTypeCast(c, constant.NewInt(0, types.I64), alloc.Elem)
		// val = nil
		return nil
	}
	block.NewStore(val, alloc)

	return nil
}
