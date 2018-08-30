package ast

import (
	"bytes"
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
	BodyParser     *Parser // the parser that can build the body block on demand
	External       bool
	Variadic       bool
	Nomangle       bool
	ReturnType     GeodeTypeRef
	DeclKeyword    FuncDeclKeywordType
	ImplicitReturn bool
	HasUnknownType bool
	Package        *Package
	// A cache so we can remember the name of the function to codegen
	// This is because between the Program.CompileFunction, where we
	// can compile variants, and the codegen section of the function,
	// we lose state. So instead we can store it in the function node
	// itself and just reach into the Variants map to get the correct
	// value for the function
	NameCache string
	Variants  map[string]*ir.Function // A mapping from mangled names to llvm functions

	Compiled bool
	// CompiledValue  *ir.Function

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
			if n.HasUnknownType {
				funcArgs = append(funcArgs, nil)
				argTypes = append(argTypes, nil)
				continue
			} else {
				log.Fatal("Unable to find type with name %q for function %s\n", arg.Type.Name, n.Name)
			}
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
	checkerr := n.Check(prog.Scope, prog.Compiler)
	if checkerr != nil {
		log.Fatal("Check error: %s\n", checkerr.Error())
	}
	funcArgs, _ := n.Arguments(prog.Scope)

	namestring := n.NameCache

	ty := prog.Scope.FindType(n.ReturnType.Name).Type
	ty = n.ReturnType.BuildPointerType(ty)

	function := prog.Compiler.Module.NewFunction(namestring, ty, funcArgs...)

	previousFunction := prog.Compiler.FN
	prog.Compiler.FN = function

	function.Sig.Variadic = n.Variadic

	keyName := fmt.Sprintf("%s:%s", prog.Scope.PackageName, n.Name)

	scopeItem := NewFunctionScopeItem(keyName, n, function, PublicVisibility)
	scopeItem.SetMangled(!n.Nomangle)
	prog.Scope.Add(scopeItem)

	// function.CallConv = ir.CallConvCold

	prog.Scope = prog.Scope.Parent

	prog.Compiler.FN = previousFunction
	return function
}

// MangledName returns the correctly mangled name for some function
func (n FunctionNode) MangledName(prog *Program, types []types.Type) string {
	if n.Name.Value == "main" || n.Package.Name == "_runtime" {
		return n.Name.Value
	}
	// _, types := n.Arguments(prog.Scope)
	return MangleFunctionName(fmt.Sprintf("%s:%s", n.Package.Name, n.Name.Value), types)
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

// Codegen implements Node.Codegen for FunctionNode
func (n FunctionNode) Codegen(prog *Program) value.Value {
	prog.Scope = prog.Scope.SpawnChild()

	checkerr := n.Check(prog.Scope, prog.Compiler)
	if checkerr != nil {
		n.SyntaxError()
		log.Fatal("Check error: %s\n", checkerr.Error())
	}

	namestring := n.Name.String()

	if namestring != "main" || !n.Nomangle {
		namestring = fmt.Sprintf("%s:%s", prog.Package.Name, n.Name)
	}

	function := n.Variants[n.NameCache] // at this point it should only be compiled
	prog.Compiler.FN = function

	// If the function is external (has ... at the end) we don't build a block
	if !n.External {
		// Create the entrypoint to the function
		entryBlock := ir.NewBlock(n.Name.String() + "-entry")

		prog.Compiler.FN.AppendBlock(entryBlock)
		prog.Compiler.PushBlock(entryBlock)

		// Construct the prelude of this function
		// The prelude contains information about
		// initializing the runtime.
		createPrelude(prog, n)
		if len(function.Params()) > 0 {
			// prog.Compiler.CurrentBlock().AppendInst(NewLLVMComment(n.Name.String() + " arguments:"))
		}
		for _, arg := range function.Params() {
			alloc := prog.Compiler.CurrentBlock().NewAlloca(arg.Type())
			prog.Compiler.CurrentBlock().NewStore(arg, alloc)
			// Set the scope item
			scItem := NewVariableScopeItem(arg.Name, alloc, PrivateVisibility)
			prog.Scope.Add(scItem)
		}
		// prog.Compiler.CurrentBlock().AppendInst(NewLLVMComment(fmt.Sprintf("%s code:", n.Name.String())))
		// Gen the body of the function

		if n.BodyParser != nil {
			n.Body = n.BodyParser.parseBlockStmt()
		}

		block := n.Body.Codegen(prog).(*ir.BasicBlock)

		if block.Term == nil {

			ty := prog.Scope.FindType(n.ReturnType.Name).Type

			// log.Error("Function %s is missing a return statement in the root block. Defaulting to 0\n", n.Name)
			v := createTypeCast(prog, constant.NewInt(0, types.I64), ty)
			block.NewRet(v)
		}
		prog.Compiler.PopBlock()
	}

	prog.Scope = prog.Scope.Parent
	return function
}

func createPrelude(prog *Program, n FunctionNode) {
	if prog.Compiler.FN.Name == "main" {
		prog.Compiler.CurrentBlock().AppendInst(NewLLVMRaw(""))
		prog.Compiler.CurrentBlock().AppendInst(NewLLVMComment("Runtime prelude"))
		prog.Compiler.CurrentBlock().AppendInst(NewLLVMRaw(""))
		// Initialize the garbage collector at the first value allocted to the stack.
		// QuickParseIdentifier("").Codegen(prog)
		QuickParseExpression("GC_init();").Codegen(prog)
		QuickParseExpression("GC_enable_incremental();").Codegen(prog)

		prog.Compiler.CurrentBlock().AppendInst(NewLLVMRaw(""))
		prog.Compiler.CurrentBlock().AppendInst(NewLLVMComment("Global Initializations"))
		prog.Compiler.CurrentBlock().AppendInst(NewLLVMRaw(""))
		for _, init := range prog.Initializations {
			init.Codegen(prog)
		}

		prog.Compiler.CurrentBlock().AppendInst(NewLLVMRaw(""))
		prog.Compiler.CurrentBlock().AppendInst(NewLLVMComment("User Code"))
		prog.Compiler.CurrentBlock().AppendInst(NewLLVMRaw(""))
	}
}

func (n FunctionNode) String() string {
	buff := &bytes.Buffer{}

	fmt.Fprintf(buff, "func %s(", n.Name)

	for i, arg := range n.Args {
		fmt.Fprintf(buff, "%s", arg)
		if i < len(n.Args)-1 || n.Variadic {
			fmt.Fprintf(buff, ", ")
		}
	}

	if n.Variadic {
		fmt.Fprintf(buff, "...")
	}

	fmt.Fprintf(buff, ") %s ", n.ReturnType)

	if n.External {
		fmt.Fprintf(buff, "...")
	} else {
		fmt.Fprintf(buff, "%s", n.Body)
	}

	return buff.String()
}
