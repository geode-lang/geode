package ast

import (
	"fmt"

	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/constant"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// FuncDeclKeywordType lets the compiler keep track of
// what keyword was used to declare a function
type FuncDeclKeywordType int

// Type of function declaration keywords
const (
	DeclKeywordFunc FuncDeclKeywordType = iota
	DeclKeywordPure
)

// FunctionNode is the representation of some function. It has methods
// on it to declare the function as well as codegen. A function has
// a list of VariableNodes for arguments and a single block for a body,
// all of which are codegenned.
type FunctionNode struct {
	NodeType
	TokenReference

	Name           NamedReference
	Args           []VariableDefnNode
	Body           BlockNode
	External       bool
	Variadic       bool
	Nomangle       bool
	ReturnType     GeodeTypeRef
	Generics       []*GenericSymbol
	DeclKeyword    FuncDeclKeywordType
	ImplicitReturn bool

	line   int
	column int
}

// NameString implements Node.NameString
func (n FunctionNode) NameString() string { return "FunctionNode" }

// InferType implements Node.InferType
func (n FunctionNode) InferType(scope *Scope) string {
	return "function" //scope.FindType(n.ReturnType.Name).Type
}

// Arguments returns some FunctionNode's arguments
func (n FunctionNode) Arguments(scope *Scope) ([]*types.Param, []types.Type) {
	funcArgs := make([]*types.Param, 0)
	argTypes := make([]types.Type, 0)
	for _, arg := range n.Args {
		found := scope.FindType(arg.Type.Name)
		if found == nil {
			log.Fatal("Unable to find type with name %q\n", arg.Type.Name)
		}
		ty := found.Type
		ty = arg.Type.BuildPointerType(ty)
		p := ir.NewParam(arg.Name.String(), ty)
		funcArgs = append(funcArgs, p)
		argTypes = append(argTypes, p.Type())
	}
	return funcArgs, argTypes
}

// Declare declares some FunctionNode's sig
func (n FunctionNode) Declare(prog *Program) *ir.Function {

	prog.Scope = prog.Scope.SpawnChild()
	c := prog.Compiler
	checkerr := n.Check(prog.Scope, c)
	if checkerr != nil {
		log.Fatal("Check error: %s\n", checkerr.Error())
	}
	funcArgs, _ := n.Arguments(prog.Scope)

	namestring := n.Name.String()
	// We need to do some special checks if the function is main. It's special.
	// For instance, it must return int type.
	if namestring == "main" {
		if n.ReturnType.Name != "int" {
			log.Fatal("Main function must return type int. Called for type '%s'\n", n.ReturnType)
		}
	} else {
		namestring = n.MangledName(prog.Scope, c, nil)
	}

	ty := prog.Scope.FindType(n.ReturnType.Name).Type
	ty = n.ReturnType.BuildPointerType(ty)

	function := c.Module.NewFunction(namestring, ty, funcArgs...)

	c.FN = function

	// if n.Variadic && !n.External {
	// 	log.Fatal("Function '%s' is variadic and has a body. This only allowed for external functions.\n", n.Name)
	// }

	function.Sig.Variadic = n.Variadic

	keyName := fmt.Sprintf("%s:%s", c.Scope.PackageName, n.Name)

	scopeItem := NewFunctionScopeItem(keyName, n, function, PublicVisibility)
	scopeItem.SetMangled(!n.Nomangle)
	c.Scope.Add(scopeItem)

	// c.Module.NewGlobalDecl(fmt.Sprintf("_ret_%s", function.Name), function.Sig.Ret)
	c.FN = nil

	prog.Scope = prog.Scope.Parent
	return function
}

// MangledName will return the mangled name for a function node
func (n FunctionNode) MangledName(scope *Scope, c *Compiler, generics []*GenericSymbol) string {
	var ns string
	if n.Nomangle {
		return n.Name.String()
	}
	_, argTypes := n.Arguments(scope)
	// Parse the namespace and name from the funciton name
	namespace, name := parseName(n.Name.String())
	if namespace == "" {
		namespace = c.Scope.PackageName
	}

	ns = fmt.Sprintf("%s:%s", namespace, name)
	name = MangleFunctionName(ns, argTypes, n.Generics)
	return name
}

// Check makes sure a function follows the correct limitations set by the language
// ex:
//    when the function is pure, it cannot accept pointer or have a block as a body.
func (n FunctionNode) Check(scope *Scope, c *Compiler) error {
	if n.DeclKeyword == DeclKeywordPure {
		_, argtypes := n.Arguments(scope)
		for _, arg := range argtypes {
			if types.IsPointer(arg) {
				return fmt.Errorf("pure function '%s' is not allowed to accept pointers as arguments", n.Name)
			}
		}
		if n.ReturnType.PointerLevel != 0 {
			return fmt.Errorf("pure function '%s' is not allowed to return a pointer", n.Name)
		}

		if !n.ImplicitReturn {
			return fmt.Errorf("pure function '%s' must have an implcit return, not a block", n.Name)
		}
	}
	return nil
}

// CodegenGeneric takes some generic type symbols, checks if they could work, and generates
// a new function using those as types.
func (n FunctionNode) CodegenGeneric(prog *Program, g []*GenericSymbol) value.Value {
	if len(n.Generics) != len(g) {
		n.SyntaxError()
		log.Fatal("Generics used in function call on '%s' are not of the correct length. Passed: %d, Expected: %d", n.Name, len(g), len(n.Generics))
	}

	return n.Codegen(prog)

}

// Codegen implements Node.Codegen for FunctionNode
func (n FunctionNode) Codegen(prog *Program) value.Value {
	scope := prog.Scope
	c := prog.Compiler

	checkerr := n.Check(scope, c)
	if checkerr != nil {
		n.SyntaxError()
		log.Fatal("Check error: %s\n", checkerr.Error())
	}

	namestring := n.Name.String()

	if namestring != "main" || !n.Nomangle {
		namestring = n.MangledName(scope, c, n.Generics)
	}

	declared, _ := c.Scope.FindFunctions(namestring)

	if len(declared) != 1 {
		n.SyntaxError()
		log.Fatal("Unable to find function declaration for '%s'\n", namestring)
	}
	function := declared[0].Value().(*ir.Function)
	c.FN = function

	// If the function is external (has ... at the end) we don't build a block
	if !n.External {
		// Create the entrypoint to the function
		entryBlock := ir.NewBlock(n.Name.String() + "-entry")
		c.FN.AppendBlock(entryBlock)
		c.PushBlock(entryBlock)

		// Construct the prelude of this function
		// The prelude contains information about
		// initializing the runtime.
		createPrelude(prog, n)
		if len(function.Params()) > 0 {
			c.CurrentBlock().AppendInst(NewLLVMComment(n.Name.String() + " arguments:"))
		}
		for _, arg := range function.Params() {
			alloc := c.CurrentBlock().NewAlloca(arg.Type())
			c.CurrentBlock().NewStore(arg, alloc)
			// Set the scope item
			scItem := NewVariableScopeItem(arg.Name, alloc, PrivateVisibility)
			scope.Add(scItem)
		}
		// c.CurrentBlock().AppendInst(NewLLVMComment(fmt.Sprintf("%s code:", n.Name.String())))
		// Gen the body of the function
		n.Body.Codegen(prog)
		if c.CurrentBlock().Term == nil {
			ty := scope.FindType(n.ReturnType.Name).Type
			// log.Error("Function %s is missing a return statement in the root block. Defaulting to 0\n", n.Name)
			v := createTypeCast(prog, constant.NewInt(0, types.I64), ty)
			c.CurrentBlock().NewRet(v)
		}
		c.PopBlock()
	}
	return function
}

func createPrelude(prog *Program, n FunctionNode) {
	if prog.Compiler.FN.Name == "main" {
		prog.Compiler.CurrentBlock().AppendInst(NewLLVMComment("runtime prelude:"))
		// Initialize the garbage collector at the first value allocted to the stack.
		QuickParseIdentifier("byte __GC_BASE_POINTER;").Codegen(prog)
		QuickParseExpression("___geodegcinit(&__GC_BASE_POINTER);").Codegen(prog)

		// fields := map[string]value.Value{}

		// fields["foo"] = constant.NewInt(12, types.I64)
		// fields["bar"] = constant.NewInt(42, types.I64)
		// NewClassInstance(prog, prog.Scope.FindType("Foo").Type.(*types.StructType), fields)

	}

}
