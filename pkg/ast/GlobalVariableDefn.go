package ast

import (
	"bytes"
	"fmt"

	"github.com/geode-lang/geode/llvm/ir"
	"github.com/geode-lang/geode/llvm/ir/constant"
	"github.com/geode-lang/geode/llvm/ir/value"
)

// GlobalVariableDeclNode -
type GlobalVariableDeclNode struct {
	NodeType
	TokenReference

	Type     TypeNode
	External bool
	Name     IdentNode
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

	varType, err := n.Type.GetType(prog)
	if err != nil {
		return nil, err
	}

	init := constant.NewZeroInitializer(varType)

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

	var val value.Value
	var err error

	pkgCache := prog.Package
	prog.Package = n.Package
	if n.Body != nil {

		assign := AssignmentNode{}

		assign.NodeType = nodeAssignment

		assign.Token = n.Token

		assign.Assignee = n.Name

		assign.Value = n.Body.(Accessable)

		val, err = assign.Codegen(prog)
		if err != nil {
			prog.Package = pkgCache
			return nil, err
		}
	}

	prog.Package = pkgCache
	return val, nil
}

func (n GlobalVariableDeclNode) String() string {
	buff := &bytes.Buffer{}
	fmt.Fprintf(buff, "%s %s", n.Type, n.Name)

	if !n.External {
		fmt.Fprintf(buff, " = %s", n.Body)
	}
	return buff.String()
}
