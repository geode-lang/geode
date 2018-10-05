package ast

import (
	"fmt"

	"github.com/geode-lang/geode/pkg/arg"

	"github.com/geode-lang/geode/llvm/ir"
	"github.com/geode-lang/geode/llvm/ir/metadata"
	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
	"github.com/geode-lang/geode/pkg/util/log"
)

// NameType is a type to notate what kind of name a IdentNode is
type NameType int

// Various NameType constants
const (
	ClassMethodNameType NameType = iota
)

// IdentNode is a reference to an item on the scope via some name
// this can be used to access variable allocs, function defns, or types
type IdentNode struct {
	NodeType
	TokenReference
	Accessable
	Assignable
	Callable

	Value    string
	NameType NameType
}

// NewIdentNode returns a new name reference with a string as it's name
func NewIdentNode(name string) IdentNode {
	n := IdentNode{}

	n.Value = name
	return n
}

// NameString implements Node.NameString
func (n IdentNode) NameString() string { return "IdentNode" }

// GetFunc implements Callable.GetFunc
func (n IdentNode) GetFunc(prog *Program, argTypes []types.Type) (*ir.Function, []value.Value, error) {

	ns, nm := ParseName(n.String())
	if ns == "" {
		ns = prog.Scope.PackageName
	} else if !prog.Package.HasAccessToPackage(ns) {
		return nil, nil, fmt.Errorf("package %s doesn't load package %s but attempts to call %s:%s", prog.Scope.PackageName, ns, ns, nm)
	}
	searchNames := []string{
		fmt.Sprintf("%s:%s", ns, nm),
		fmt.Sprintf("%s:%s", prog.Package.Name, nm),
		nm,
	}
	f, err := prog.FindFunction(searchNames, argTypes)
	return f, nil, err
}

func (n IdentNode) String() string {
	return n.Value
}

// Alloca returns the nearest alloca instruction in this scope with the given name
func (n IdentNode) Alloca(prog *Program) value.Value {

	searchPaths := make([]string, 0)
	searchPaths = append(searchPaths, n.Value)
	searchPaths = append(searchPaths, fmt.Sprintf("%s:%s", prog.Package.Name, n.Value))

	if prog.Scope == nil {
		n.SyntaxError()
		fmt.Println(n)
		return nil
	}
	scopeitem, found := prog.Scope.Find(searchPaths)

	var alloc value.Value
	success := false

	// fmt.Println(prog.Scope.AllNames())
	if !found {
		// log.Fatal("Unable to find named reference %s, search paths: [%s]\n", n, strings.Join(searchPaths, ", "))

		// If it is not found, I need to create a new node. Assignment will never fail when assigning to
		return nil
	}

	if alloc, success = scopeitem.(VariableScopeItem).Value().(*ir.InstAlloca); success {
		return alloc
	}

	if alloc, success = scopeitem.(VariableScopeItem).Value().(*ir.Global); success {
		return alloc
	}

	log.Fatal("Unknown Type in VariableScopeItem in search paths: [%s]\n", searchPaths)
	return nil
}

// Load returns a load instruction on a named reference with the given name
func (n IdentNode) Load(block *ir.BasicBlock, prog *Program) *ir.InstLoad {
	alloc := n.Alloca(prog)
	if alloc == nil {
		return nil
	}
	return block.NewLoad(alloc)
}

// GenAssign implements Assignable.GenAssign
func (n IdentNode) GenAssign(prog *Program, assignment value.Value, options ...AssignableOption) (value.Value, error) {
	alloca := n.Alloca(prog)

	if alloca == nil {
		alloca = prog.Compiler.CurrentBlock().NewAlloca(assignment.Type())
		prog.Scope.Add(NewVariableScopeItem(n.Value, alloca, PublicVisibility))
	}
	store := prog.Compiler.CurrentBlock().NewStore(assignment, alloca)

	if *arg.EnableDebug {
		md := &metadata.Metadata{}
		md.Add(metadata.NewRaw(n.Token.DILocation(prog.Scope.DebugInfo)))
		store.Metadata["dbg"] = md
	}

	return assignment, nil
}

// Codegen implements Node.Codegen
func (n IdentNode) Codegen(prog *Program) (value.Value, error) {
	return n.GenAccess(prog)
}

// GenAccess implements Accessable.GenAccess
func (n IdentNode) GenAccess(prog *Program) (value.Value, error) {
	load := n.Load(prog.Compiler.CurrentBlock(), prog)
	if load == nil {
		return nil, fmt.Errorf("unable to load/access value for identifier %q", n)
	}
	return load, nil
}

// Type implements Assignable.Type
func (n IdentNode) Type(prog *Program) (types.Type, error) {
	ref := n.Alloca(prog)

	if alloca, success := ref.(*ir.InstAlloca); success {
		return alloca.Elem, nil
	}

	if global, success := ref.(*ir.Global); success {
		return global.Type().(*types.PointerType).Elem, nil
	}
	return nil, nil
}
