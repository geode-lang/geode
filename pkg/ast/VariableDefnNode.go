package ast

import (
	"bytes"
	"fmt"

	"github.com/geode-lang/geode/llvm/ir"
	"github.com/geode-lang/geode/llvm/ir/constant"
	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
	"github.com/geode-lang/geode/pkg/util/log"
)

// VariableDefnNode -
type VariableDefnNode struct {
	NodeType
	TokenReference

	Typ            TypeNode
	HasValue       bool
	Name           IdentNode
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

	block := prog.Compiler.CurrentBlock()

	n.Package = prog.Package

	f := block.Parent

	name := n.Name

	prog.Compiler.EmptyTypeStack()

	if !n.NeedsInference {
		found, err := prog.FindType(n.Typ.Name)
		if err != nil {
			return nil, err
		}
		if found == nil {
			n.SyntaxError()
			log.Fatal("Unable to find type named %q for variable declaration\n", n.Typ.Name)
		}
		valType, err = n.Typ.GetType(prog)
		if err != nil {
			return nil, err
		}
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
	prog.Scope.Add(scItem)

	if !n.NeedsInference && val != nil {
		val, err = createTypeCast(prog, val, alloc.Elem)
		if err != nil {
			return nil, err
		}
	}

	// If the value is nil, we need to pull the default value for a given type.
	if val == nil {
		val = constant.NewZeroInitializer(alloc.Elem)
	}

	block.NewStore(val, alloc)

	return alloc, nil
}

// GenAssign implements Assignable.GenAssign
func (n VariableDefnNode) GenAssign(prog *Program, val value.Value, options ...AssignableOption) (value.Value, error) {

	alloc, err := n.Codegen(prog)
	if err != nil {
		return nil, err
	}
	prog.Compiler.CurrentBlock().NewStore(val, alloc)
	return nil, nil
}

// Type implements Assignable.Type
func (n VariableDefnNode) Type(prog *Program) (types.Type, error) {
	return n.Typ.GetType(prog)
}

func (n VariableDefnNode) String() string {
	buff := &bytes.Buffer{}

	if n.NeedsInference {
		fmt.Fprintf(buff, "let %s", n.Name)
	} else {
		fmt.Fprintf(buff, "%s %s", n.Typ, n.Name)
	}

	if n.HasValue {
		fmt.Fprintf(buff, " := %s", n.Body)
	}
	return buff.String()
}
