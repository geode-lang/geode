package ast

import (
	"bytes"
	"fmt"

	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/value"
)

// GlobalVariableDeclNode -
type GlobalVariableDeclNode struct {
	NodeType
	TokenReference

	Type     GeodeTypeRef
	External bool
	Name     NamedReference
	Body     Node

	GlobalDecl *ir.Global
	Package    *Package
}

// NameString implements Node.NameString
func (n GlobalVariableDeclNode) NameString() string { return "GlobalVariableDeclNode" }

// Declare a global variable declaration
func (n GlobalVariableDeclNode) Declare(prog *Program) (value.Value, error) {
	var name string

	if n.External {
		name = n.Name.Value
	} else {
		name = fmt.Sprintf("%s:%s", prog.Package.Name, n.Name)
	}

	found, err := prog.FindType(n.Type.Name)
	if err != nil {
		return nil, err
	}
	varType := n.Type.BuildPointerType(found)

	init := DefaultValue(varType)

	decl := prog.Module.NewGlobalDef(name, init)

	if !n.External {
		decl.Name = MangleVariableName(name)
	}

	n.GlobalDecl = decl
	n.Package = prog.Package

	scopeName := fmt.Sprintf("%s:%s", prog.Package.Name, n.Name)
	n.Name.Value = scopeName
	prog.Scope.GetRoot().Add(NewVariableScopeItem(scopeName, decl, PublicVisibility))

	prog.RegisterGlobalVariableInitialization(&n)

	return decl, nil
}

// Codegen a global variable declaration
func (n GlobalVariableDeclNode) Codegen(prog *Program) (value.Value, error) {

	if n.Body != nil {
		assign := AssignmentNode{}

		assign.NodeType = nodeAssignment

		assign.Token = n.Token

		assign.Assignee = n.Name

		assign.Value = n.Body.(Accessable)

		return assign.Codegen(prog)
	}
	return nil, nil
}

func (n GlobalVariableDeclNode) String() string {
	buff := &bytes.Buffer{}
	fmt.Fprintf(buff, "%s %s", n.Type, n.Name)

	if !n.External {
		fmt.Fprintf(buff, " := %s", n.Body)
	}
	return buff.String()
}
