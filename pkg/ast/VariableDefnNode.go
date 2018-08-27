package ast

import (
	"bytes"
	"fmt"

	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/value"
)

// VariableDefnNode -
type VariableDefnNode struct {
	NodeType
	TokenReference

	Type      GeodeTypeRef
	HasValue  bool
	Name      NamedReference
	Body      Node
	MustInfer bool

	Package *Package
}

// NameString implements Node.NameString
func (n VariableDefnNode) NameString() string { return "VariableDefnNode" }

// InferType implements Node.InferType
func (n VariableDefnNode) InferType(scope *Scope) string {
	return n.Type.Name
}

// Codegen implements Node.Codegen for VariableDefnNode
func (n VariableDefnNode) Codegen(prog *Program) value.Value {
	scope := prog.Scope

	block := prog.Compiler.CurrentBlock()

	n.Package = prog.Package

	f := block.Parent

	name := n.Name

	var alloc *ir.InstAlloca
	var val value.Value

	found := scope.FindType(n.Type.Name)
	if found == nil {
		log.Fatal("Unable to find type named %q\n", n.Type.Name)
	}
	ty := found.Type
	ty = n.Type.BuildPointerType(ty)
	// block.AppendInst(NewLLVMComment("%s %s", ty, name))
	alloc = createBlockAlloca(f, ty, name.String())

	prog.Compiler.typeCache = alloc.Elem
	scItem := NewVariableScopeItem(name.String(), alloc, PrivateVisibility)
	scope.Add(scItem)

	if n.HasValue {

		// Construct the body
		if n.Body != nil {
			val = n.Body.Codegen(prog)
			if val == nil {
				return val
			}
		}

		val = createTypeCast(prog, val, alloc.Elem)

	} else {
		defa := DefaultValue(alloc.Elem)
		if defa != nil {
			block.NewStore(defa, alloc)
		}
		return nil
	}

	block.NewStore(val, alloc)

	return nil
}

func (n VariableDefnNode) String() string {
	buff := &bytes.Buffer{}
	fmt.Fprintf(buff, "%s %s", n.Type, n.Name)

	if n.HasValue {
		fmt.Fprintf(buff, " := %s", n.Body)
	}
	return buff.String()
}
