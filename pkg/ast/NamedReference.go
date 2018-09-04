package ast

import (
	"fmt"
	"strings"

	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// NameType is a type to notate what kind of name a NamedReference is
type NameType int

// Various NameType constants
const (
	ClassMethodNameType NameType = iota
)

// NamedReference is a reference to an item on the scope via some name
// this can be used to access variable allocs, function defns, or types
type NamedReference struct {
	Value    string
	NameType NameType
}

// NewNamedReference returns a new name reference with a string as it's name
func NewNamedReference(name string) NamedReference {
	n := NamedReference{}

	n.Value = name
	return n
}

// GetFunc implements Callable.GetFunc
func (n NamedReference) GetFunc(prog *Program, argTypes []types.Type) (*ir.Function, error) {

	ns, nm := ParseName(n.String())
	if ns == "" {
		ns = prog.Scope.PackageName
	} else if !prog.Package.HasAccessToPackage(ns) {
		return nil, fmt.Errorf("package %s doesn't load package %s but attempts to call %s:%s", prog.Scope.PackageName, ns, ns, nm)
	}
	searchNames := []string{
		fmt.Sprintf("%s:%s", ns, nm),
		fmt.Sprintf("%s:%s", prog.Scope.PackageName, nm),
		nm,
	}
	return prog.FindFunction(searchNames, argTypes)
}

func (n NamedReference) String() string {
	return n.Value
}

// Alloca returns the nearest alloca instruction in this scope with the given name
func (n NamedReference) Alloca(prog *Program) value.Value {

	searchPaths := make([]string, 0)
	searchPaths = append(searchPaths, n.Value)
	searchPaths = append(searchPaths, fmt.Sprintf("%s:%s", prog.Package.Name, n.Value))

	scopeitem, found := prog.Scope.Find(searchPaths)

	// fmt.Println(prog.Scope.AllNames())
	if !found {
		log.Fatal("Unable to find named reference %s, search paths: [%s]\n", n, strings.Join(searchPaths, ", "))
	}

	var alloc value.Value
	success := false

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
func (n NamedReference) Load(block *ir.BasicBlock, prog *Program) *ir.InstLoad {
	return block.NewLoad(n.Alloca(prog))
}

// GenAssign implements Assignable.GenAssign
func (n NamedReference) GenAssign(prog *Program, assignment value.Value) (value.Value, error) {
	alloca := n.Alloca(prog)
	prog.Compiler.CurrentBlock().NewStore(assignment, alloca)
	return assignment, nil
}

// GenAccess implements Accessable.GenAccess
func (n NamedReference) GenAccess(prog *Program) (value.Value, error) {
	return n.Load(prog.Compiler.CurrentBlock(), prog), nil
}

// Type implements Assignable.Type
func (n NamedReference) Type(prog *Program) (types.Type, error) {
	ref := n.Alloca(prog)

	if alloca, success := ref.(*ir.InstAlloca); success {
		return alloca.Elem, nil
	}

	if global, success := ref.(*ir.Global); success {
		return global.Type().(*types.PointerType).Elem, nil
	}
	return types.Void, nil
}
