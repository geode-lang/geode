package ast

import (
	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/value"
)

// VariableDefnNode -
type VariableDefnNode struct {
	NodeType
	TokenReference

	Type     GeodeTypeRef
	HasValue bool
	Name     *NamedReference
	Body     Node
}

// NameString implements Node.NameString
func (n VariableDefnNode) NameString() string { return "VariableDefnNode" }

// InferType implements Node.InferType
func (n VariableDefnNode) InferType(scope *Scope) string {
	return n.Type.Name
}

// Codegen implements Node.Codegen for VariableDefnNode
func (n VariableDefnNode) Codegen(scope *Scope, c *Compiler) value.Value {

	block := c.CurrentBlock()

	f := block.Parent

	name := n.Name
	var alloc *ir.InstAlloca
	var val value.Value

	ty := scope.FindType(n.Type.Name).Type
	ty = n.Type.BuildPointerType(ty)
	block.AppendInst(NewLLVMComment("%s %s", ty, name))
	alloc = createBlockAlloca(f, ty, name.String())

	c.typeCache = alloc.Elem
	scItem := NewVariableScopeItem(name.String(), alloc, PrivateVisibility)
	scope.Add(scItem)

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
		return nil
	}

	block.AppendInst(NewLLVMComment("%s <- %s", name, val.Type()))

	block.NewStore(val, alloc)

	return nil
}
