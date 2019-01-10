package ast

import (
	"bytes"
	"fmt"

	"github.com/geode-lang/geode/pkg/arg"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// FuncDeclKeywordType lets the compiler keep track of
// what keyword was used to declare a function
type FuncDeclKeywordType int

// Type of function declaration keywords
const (
	DeclKeywordFunc FuncDeclKeywordType = iota
	DeclKeywordPure
)

// FunctionArg represents a single argument to a function
type FunctionArg struct {
	Type TypeNode
	Name string
}

// FunctionNode is the representation of some function. It has methods
// on it to declare the function as well as codegen. A function has
// a list of VariableNodes for arguments and a single block for a body,
// all of which are codegenned.
type FunctionNode struct {
	NodeType
	TokenReference

	Name           IdentNode
	Args           []FunctionArg
	Body           BlockNode
	BodyParser     *Parser // the parser that can build the body block on demand
	External       bool
	Variadic       bool
	Nomangle       bool
	ReturnType     TypeNode
	DeclKeyword    FuncDeclKeywordType
	ImplicitReturn bool
	HasUnknownType bool
	Package        *Package
	IsMethod       bool

	// A cache so we can remember the name of the function to codegen
	// This is because between the Program.GetFunction, where we
	// can compile variants, and the codegen section of the function,
	// we lose state. So instead we can store it in the function node
	// itself and just reach into the Variants map to get the correct
	// value for the function
	NameCache string
	Variants  map[string]*ir.Func // A mapping from mangled names to llvm functions

	Compiled bool
	// CompiledValue  *ir.Func

	line   int
	column int
}

// NameString implements Node.NameString
func (n FunctionNode) NameString() string { return "FunctionNode" }

// Arguments returns some FunctionNode's arguments
func (n FunctionNode) Arguments(prog *Program) ([]*ir.Param, []types.Type, error) {
	funcArgs := make([]*ir.Param, 0)
	argTypes := make([]types.Type, 0)
	for _, arg := range n.Args {
		found, _ := prog.FindType(arg.Type.Name)
		if found == nil {
			if n.HasUnknownType {
				funcArgs = append(funcArgs, nil)
				argTypes = append(argTypes, nil)
				continue
			} else {
				return nil, nil, fmt.Errorf("unable to find type with name %q for function %s (%s)", arg.Type.Name, n.Name, n.Token.FileInfo())
			}
		}
		ty, err := arg.Type.GetType(prog)
		if err != nil {
			return nil, nil, err
		}
		p := ir.NewParam(arg.Name, ty)
		funcArgs = append(funcArgs, p)
		argTypes = append(argTypes, p.Type())
	}
	return funcArgs, argTypes, nil
}

// Declare a function in the module of the program for future use. This allows recursive calls
// to the function in the codegen step. This is also the last step for a function that is external
// as external functions only need a declaration for their signature.
// This function will generate the function object and push it to the map of other functions
// in the program for use in calls
func (n FunctionNode) Declare(prog *Program) (*ir.Func, error) {
	// Spawn a new scope
	prog.ScopeDown(n.Token)

	// Check that the funciton is allowed or follows rules set in place
	// by geode
	checkerr := n.Check(prog)
	if checkerr != nil {
		return nil, fmt.Errorf("check error: %s", checkerr.Error())
	}
	funcArgs, _, err := n.Arguments(prog)
	if err != nil {
		return nil, err
	}

	namestring := n.NameCache

	ty, err := n.ReturnType.GetType(prog)
	if err != nil {
		return nil, err
	}

	function := prog.Compiler.Module.NewFunc(namestring, ty, funcArgs...)

	prog.Compiler.PushFunc(function)
	defer prog.Compiler.PopFunc()

	function.Sig.Variadic = n.Variadic
	keyName := fmt.Sprintf("%s:%s", prog.Scope.PackageName, n.Name)

	scopeItem := NewFunctionScopeItem(keyName, n, function, PublicVisibility)
	scopeItem.SetMangled(!n.Nomangle)
	prog.Scope.GetRoot().Add(scopeItem)

	if err := prog.ScopeUp(); err != nil {
		return nil, err
	}
	return function, nil
}

// MangledName returns the correctly mangled name for some function
func (n FunctionNode) MangledName(prog *Program, types []types.Type) string {

	ret, _ := n.ReturnType.GetType(prog)
	if n.IsMethod {
		return MangleFunctionName(n.Name.Value, types, ret)
	}

	if n.Name.Value == "main" || (n.Package != nil && n.Package.Name == "runtime") {
		return n.Name.Value
	}
	// _, types := n.Arguments(prog.Scope)
	return MangleFunctionName(fmt.Sprintf("%s:%s", n.Package.Name, n.Name.Value), types, ret)
}

// Check makes sure a function follows the correct limitations set by the language
// ex:
//    when the function is pure, it cannot accept pointer or have a block as a body.
func (n FunctionNode) Check(prog *Program) error {
	if n.DeclKeyword == DeclKeywordPure {
		_, argtypes, err := n.Arguments(prog)
		if err != nil {
			return err
		}
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
func (n FunctionNode) Codegen(prog *Program) (value.Value, error) {
	prog.ScopeDown(n.Token)

	checkerr := n.Check(prog)
	if checkerr != nil {
		n.SyntaxError()
		return nil, fmt.Errorf("check error: %s", checkerr.Error())
	}

	namestring := n.Name.String()

	if (namestring != "main" || !n.Nomangle) && !n.IsMethod {
		namestring = fmt.Sprintf("%s:%s", prog.Package.Name, n.Name)
	}

	function := n.Variants[n.NameCache] // at this point it should only be compiled

	prog.Compiler.PushFunc(function)
	defer prog.Compiler.PopFunc()

	// If the function is external (has ... at the end) we don't build a block
	if !n.External {
		// Create the entrypoint to the function
		curFunc := prog.Compiler.CurrentFunc()
		entryBlock := curFunc.NewBlock(n.Name.String() + "_entry")
		prog.Compiler.PushBlock(entryBlock)

		// Construct the prelude of this function
		// The prelude contains information about
		// initializing the runtime.
		createInitializationPrelude(prog, n)
		if len(function.Params) > 0 {
			// prog.Compiler.CurrentBlock().AppendInst(NewLLVMComment(n.Name.String() + " arguments:"))
		}
		for _, arg := range function.Params {
			alloc := prog.Compiler.CurrentBlock().NewAlloca(arg.Type())
			prog.Compiler.CurrentBlock().NewStore(arg, alloc)
			// Set the scope item
			scItem := NewVariableScopeItem(arg.Name(), alloc, PrivateVisibility)
			prog.Scope.Add(scItem)
		}
		// Gen the body of the function
		if n.BodyParser != nil {
			n.Body = n.BodyParser.parseBlockStmt()
		}
		var block *ir.Block
		var ok bool
		gen, err := n.Body.Codegen(prog)
		if err != nil {
			return nil, err
		}

		if block, ok = gen.(*ir.Block); !ok {
			return nil, fmt.Errorf("type assertion to block in function node failed")
		}

		// if the block we ended on does not return, we need to either error or return a new void
		if block.Term == nil {

			retType, err := prog.FindType(n.ReturnType.Name)
			if err != nil {
				return nil, err
			}
			if retType.Equal(types.Void) {
				// Automatically return void from the function
				// new ret interpets a nil value as returning void
				block.NewRet(nil)
			} else {
				return nil, fmt.Errorf("Function %s does not end in a return statement", namestring)
			}

		}
		prog.Compiler.PopBlock()
	}

	if err := prog.ScopeUp(); err != nil {
		return nil, err
	}
	return function, nil
}

func createInitializationPrelude(prog *Program, n FunctionNode) {

	// if the user disabled the runtime, we should just not do anything special
	// with preludes or whatnot.
	if *arg.DisableRuntime {
		return
	}
	if prog.Compiler.CurrentFunc().Name() == "main" {

		prog.NewRuntimeFunctionCall("__init_runtime")

		prog.Compiler.NewComment("User Code:")
	}

	if prog.Compiler.CurrentFunc().Name() == "__init_runtime" {
		prog.Compiler.NewComment("Runtime Prelude:")

		if len(prog.Initializations) > 0 {
			prog.Compiler.NewComment("Global Initializations:")
			for _, init := range prog.Initializations {
				init.Codegen(prog)
			}
		}
	}

	// if prog.Compiler.CurrentFunc().Name == "init"
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
