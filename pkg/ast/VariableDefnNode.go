package ast

import (
	"bytes"
	"fmt"

	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// VariableDefnNode -
type VariableDefnNode struct {
	NodeType
	TokenReference

	Type           GeodeTypeRef
	HasValue       bool
	Name           NamedReference
	Body           Node
	NeedsInference bool

	Package *Package
}

// NameString implements Node.NameString
func (n VariableDefnNode) NameString() string { return "VariableDefnNode" }

// Codegen implements Node.Codegen for VariableDefnNode
func (n VariableDefnNode) Codegen(prog *Program) (value.Value, error) {
	var err error

	var alloc *ir.InstAlloca
	var val value.Value
	var valType types.Type

	scope := prog.Scope

	block := prog.Compiler.CurrentBlock()

	n.Package = prog.Package

	f := block.Parent

	name := n.Name

	prog.Compiler.EmptyTypeStack()

	if !n.NeedsInference {
		found, err := prog.FindType(n.Type.Name)
		if err != nil {
			return nil, err
		}
		if found == nil {
			n.SyntaxError()
			log.Fatal("Unable to find type named %q for variable declaration\n", n.Type.Name)
		}
		valType = found
		valType = n.Type.BuildPointerType(valType)
	} else {

		if n.HasValue && n.Body != nil {
			v, err := n.Body.Codegen(prog)
			if err != nil {
				return nil, err
			}
			val = v
		}
		valType = val.Type()
	}

	alloc = createBlockAlloca(f, valType, name.String())

	if !n.NeedsInference {
		prog.Compiler.PushType(valType)
		if n.HasValue && n.Body != nil {
			v, err := n.Body.Codegen(prog)
			if err != nil {
				return nil, err
			}
			val = v
		}
	}

	prog.Compiler.PushType(alloc.Elem)
	scItem := NewVariableScopeItem(name.String(), alloc, PrivateVisibility)
	scope.Add(scItem)

	if !n.NeedsInference && val != nil {
		val, err = createTypeCast(prog, val, alloc.Elem)
		if err != nil {
			return nil, err
		}
	}

	// If the value is nil, we need to pull the default value for a given type.
	if val == nil {
		val = DefaultValue(alloc.Elem)
	}

	block.NewStore(val, alloc)

	return nil, nil
}

func (n VariableDefnNode) String() string {
	buff := &bytes.Buffer{}

	if n.NeedsInference {
		fmt.Fprintf(buff, "let %s", n.Name)
	} else {
		fmt.Fprintf(buff, "%s %s", n.Type, n.Name)
	}

	if n.HasValue {
		fmt.Fprintf(buff, " := %s", n.Body)
	}
	return buff.String()
}
