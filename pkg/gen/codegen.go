package gen

import (
	"crypto/rand"
	"fmt"
	"os"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func init() {
	// llvm.InitializeNativeTarget()
	// llvm.InitializeAllTargetInfos()
	// llvm.InitializeAllTargets()
	// llvm.InitializeAllTargetMCs()
	// llvm.InitializeAllAsmParsers()
	// llvm.InitializeAllAsmPrinters()
}

// Scope trees represent block scoping by having a root scope
// and children scopes that point back to their parent scope.
type Scope struct {
	Parent   *Scope
	Children []*Scope
	Vals     map[string]value.Value
}

// Set a value in this specific scope
func (s *Scope) Set(name string, val value.Value) {
	s.Vals[name] = val
}

// Find will traverse the scope tree to find some definition of a symbol
func (s *Scope) Find(name string) (value.Value, bool) {
	val, found := s.Vals[name]
	if !found && s.Parent != nil {
		return s.Parent.Find(name)
	}
	return val, found
}

// SpawnChild takes a parent scope and creates a new variable scope for scoped variable access.
func (s *Scope) SpawnChild() *Scope {
	n := &Scope{}
	n.Parent = s
	n.Vals = make(map[string]value.Value)
	n.Children = make([]*Scope, 0)
	s.Children = append(s.Children, n)
	return n
}

// NewScope creates a scope (for use when generating root scopes)
func NewScope() *Scope {
	n := &Scope{}
	n.Parent = nil
	n.Vals = make(map[string]value.Value)
	return n
}

func error(err string) value.Value {
	fmt.Println(err)
	os.Exit(-1)
	return nil
}

func mangleName(name string) string {
	b := make([]byte, 2)
	nameNumber++
	rand.Read(b)
	return fmt.Sprintf("%s_%x_%d", name, b, nameNumber)
}

func (n ifNode) Codegen(scope *Scope, c *Compiler) value.Value     { return nil }
func (n forNode) Codegen(scope *Scope, c *Compiler) value.Value    { return nil }
func (n unaryNode) Codegen(scope *Scope, c *Compiler) value.Value  { return nil }
func (n binaryNode) Codegen(scope *Scope, c *Compiler) value.Value { return nil }
func (n fnCallNode) Codegen(scope *Scope, c *Compiler) value.Value { return nil }
func (n whileNode) Codegen(scope *Scope, c *Compiler) value.Value  { return nil }

// Function Call statement Code Generator
func (n functionCallNode) Codegen(scope *Scope, c *Compiler) value.Value {
	funcs := c.RootModule.Funcs
	var callee *ir.Function
	for _, fnc := range funcs {
		if fnc.Name == n.Name {
			callee = fnc
			break
		}
	}

	if callee == nil {
		return codegenError(fmt.Sprintf("Unknown function %q referenced", n.Name))
	}

	args := []value.Value{}
	for _, arg := range n.Args {
		args = append(args, arg.Codegen(scope, c))
		if args[len(args)-1] == nil {
			return codegenError(fmt.Sprintf("Argument to function %q failed to generate code", n.Name))
		}
	}

	return c.CurrentBlock().NewCall(callee, args...) // Entry(callee, args, n.Name+"-call-response")
}

// Return statement Code Generator
func (n returnNode) Codegen(scope *Scope, c *Compiler) value.Value {
	retVal := n.Value.Codegen(scope, c)
	c.CurrentBlock().NewRet(retVal)
	return nil
}

// Int Code Generator
func (n intNode) Codegen(scope *Scope, c *Compiler) value.Value {
	// return llvm.ConstInt(llvm.Int64Type(), , true)
	return constant.NewInt(int64(n.Value), types.I32)
}

// Char Code Generator
func (n charNode) Codegen(scope *Scope, c *Compiler) value.Value { return nil }

func newCharArray(s string) *constant.Array {
	var bs []constant.Constant
	for i := 0; i < len(s); i++ {
		b := constant.NewInt(int64(s[i]), types.I8)
		bs = append(bs, b)
	}
	c := constant.NewArray(bs...)
	c.CharArray = true
	return c
}

// String Constant Code Generator
func (n stringNode) Codegen(scope *Scope, c *Compiler) value.Value {
	str := c.RootModule.NewGlobalDef(mangleName(".str"), newCharArray(n.Value))
	str.IsConst = true
	zero := constant.NewInt(0, types.I32)
	return constant.NewGetElementPtr(str, zero, zero)
}

// Float Constant Code Generator
func (n floatNode) Codegen(scope *Scope, c *Compiler) value.Value {
	return constant.NewFloat(n.Value, types.Double)
}

// Variable Reference Node Code Generator
func (n variableReferenceNode) Codegen(scope *Scope, c *Compiler) value.Value {
	v, found := scope.Find(n.Name)
	if !found {
		fmt.Printf("unknown variable name `%s`\n", n.Name)
		os.Exit(-1)
	}
	return c.CurrentBlock().NewLoad(v)
}

// Variable Node Code Generator
func (n variableNode) Codegen(scope *Scope, c *Compiler) value.Value {
	var oldvars = []value.Value{}

	f := c.CurrentBlock().Parent
	name := n.Name
	body := n.Body
	varType := n.Type

	var val value.Value
	if body != nil {
		val = body.Codegen(scope, c)
		if val == nil {
			return val // nil
		}
	}

	alloca := createEntryBlockAlloca(f, varType, n.Name)
	c.CurrentBlock().NewStore(val, alloca)

	oldVar, _ := scope.Find(name)
	oldvars = append(oldvars, oldVar)
	scope.Set(name, alloca)

	// evaluate body now that vars are in scope
	bodyVal := body.Codegen(scope, c)
	if bodyVal == nil {
		return codegenError("body returns nil") // nil
	}

	return bodyVal
}

var nameNumber int

// Code Block Code Generator
func (n blockNode) Codegen(scope *Scope, c *Compiler) value.Value {
	blockScope := scope.SpawnChild()
	name := mangleName("entry")
	c.PushBlock(c.FN.NewBlock(name))
	for _, node := range n.Nodes {
		node.Codegen(blockScope, c)
	}
	return nil
}

// Function Node Statement Code Generation
func (n functionNode) Codegen(scope *Scope, c *Compiler) value.Value {

	funcArgs := make([]*types.Param, 0)
	for _, arg := range n.Args {
		funcArgs = append(funcArgs, ir.NewParam(arg.Name, arg.Type))
	}

	function := c.RootModule.NewFunction(n.Name, n.ReturnType, funcArgs...)
	c.FN = function

	// funcArgs := []llvm.Type{}
	// for _, arg := range n.Args {
	// 	funcArgs = append(funcArgs, arg.Type.LLVMType)
	// }
	// funcType := llvm.FunctionType(n.ReturnType.LLVMType, funcArgs, false)
	// function := llvm.AddFunction(c.RootModule, n.Name, funcType)
	// function.SetVisibility(llvm.DefaultVisibility)

	// if function.Name() != n.Name {
	// 	function.EraseFromParentAsFunction()
	// 	function = c.RootModule.NamedFunction(n.Name)
	// }

	// if function.BasicBlocksCount() != 0 {
	// 	return error("redefinition of function: " + n.Name)
	// }
	// // Set all the parameter names
	// for i, param := range function.Params() {
	// 	param.SetName(n.Args[i].Name)
	// 	scope.Set(n.Args[i].Name, param)
	// }

	// if !n.IsExternal {
	// 	block := llvm.AddBasicBlock(function, "entry")
	// 	c.Builder.SetInsertPointAtEnd(block)

	// 	args := function.Params()
	// 	for i, arg := range args {
	// 		alloca := createEntryBlockAlloca(function, arg.Type(), n.Args[i].Name)
	// 		c.Builder.CreateStore(arg, alloca)
	// 		scope.Set(n.Args[i].Name, alloca)
	// 	}
	// }

	n.Body.Codegen(scope, c)

	// 	if llvm.VerifyFunction(function, llvm.PrintMessageAction) != nil {
	// 		function.EraseFromParentAsFunction()
	// 		return error("function verifiction failed")
	// 	}
	// } else {
	// 	function.SetLinkage(llvm.ExternalLinkage)
	// }

	// c.RootScope.Set(n.Name, function)

	return function
}

// CreateEntryBlockAlloca - Create an alloca instruction in the entry block of
// the function.  This is used for mutable variables etc.
func createEntryBlockAlloca(f *ir.Function, t types.Type, name string) value.Value {
	return nil
}

// Allow functions to return an error isntead of having to manage closing the program each time.
func codegenError(str string, args ...interface{}) value.Value {
	fmt.Fprintf(os.Stderr, "Error: %s\n", fmt.Sprintf(str, args...))
	return nil
}
