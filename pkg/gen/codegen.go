package gen

import (
	"fmt"
	"os"

	"github.com/go-llvm/llvm"
)

func init() {

	llvm.InitializeNativeTarget()

	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()

}

// Scope trees represent block scoping by having a root scope
// and children scopes that point back to their parent scope.
type Scope struct {
	Parent   *Scope
	Children []*Scope
	Vals     map[string]llvm.Value
}

// Set a value in this specific scope
func (s *Scope) Set(name string, val llvm.Value) {
	s.Vals[name] = val
}

// Find will traverse the scope tree to find some definition of a symbol
func (s *Scope) Find(name string) (llvm.Value, bool) {
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
	n.Vals = make(map[string]llvm.Value)
	n.Children = make([]*Scope, 0)
	s.Children = append(s.Children, n)
	return n
}

// NewScope creates a scope (for use when generating root scopes)
func NewScope() *Scope {
	n := &Scope{}
	n.Parent = nil
	n.Vals = make(map[string]llvm.Value)
	return n
}

func error(err string) llvm.Value {
	fmt.Println(err)
	os.Exit(-1)
	return llvm.Value{}
}

// Optimize the llvm target
// func Optimize() {
// 	rootFuncPassMgr.Add(execEngine.TargetData())
// 	rootFuncPassMgr.AddPromoteMemoryToRegisterPass()
// 	rootFuncPassMgr.AddInstructionCombiningPass()
// 	rootFuncPassMgr.AddReassociatePass()
// 	rootFuncPassMgr.AddGVNPass()
// 	rootFuncPassMgr.AddCFGSimplificationPass()
// 	rootFuncPassMgr.InitializeFunc()
// }

// func createAlignedAlloca(typ llvm.Type, name string) llvm.Value {
// 	funcEntry := v.currentLLVMFunction().EntryBasicBlock()

// 	// use this builder() for the variable alloca
// 	// this means all allocas go at the start of the function
// 	// so each variable is only allocated once
// 	allocBuilder := llvm.NewBuilder()
// 	defer allocBuilder.Dispose()

// 	allocBuilder.SetInsertPoint(funcEntry, funcEntry.FirstInstruction())

// 	align := targetData.ABITypeAlignment(typ)
// 	alloc := allocBuilder.CreateAlloca(typ, name)
// 	alloc.SetAlignment(align)
// 	return alloc
// }

func (n ifNode) Codegen(scope *Scope, c *Compiler) llvm.Value     { return llvm.Value{} }
func (n forNode) Codegen(scope *Scope, c *Compiler) llvm.Value    { return llvm.Value{} }
func (n unaryNode) Codegen(scope *Scope, c *Compiler) llvm.Value  { return llvm.Value{} }
func (n binaryNode) Codegen(scope *Scope, c *Compiler) llvm.Value { return llvm.Value{} }
func (n fnCallNode) Codegen(scope *Scope, c *Compiler) llvm.Value { return llvm.Value{} }
func (n whileNode) Codegen(scope *Scope, c *Compiler) llvm.Value  { return llvm.Value{} }

// Function Call statement Code Generator
func (n functionCallNode) Codegen(scope *Scope, c *Compiler) llvm.Value {
	callee := c.RootModule.NamedFunction(n.Name)

	if callee.IsNil() {
		return codegenError(fmt.Sprintf("Unknown function %q referenced", n.Name))
	}

	args := []llvm.Value{}
	for _, arg := range n.Args {
		args = append(args, arg.Codegen(scope, c))
		if args[len(args)-1].IsNil() {
			return codegenError(fmt.Sprintf("Argument to function %q failed to generate code", n.Name))
		}
	}

	return c.Builder.CreateCall(callee, args, n.Name+"-call-response")
}

// Return statement Code Generator
func (n returnNode) Codegen(scope *Scope, c *Compiler) llvm.Value {
	retVal := n.Value.Codegen(scope, c)
	c.Builder.CreateRet(retVal)
	return llvm.Value{}
}

// Int Code Generator
func (n intNode) Codegen(scope *Scope, c *Compiler) llvm.Value {
	return llvm.ConstInt(llvm.Int64Type(), uint64(n.Value), true)
}

// Char Code Generator
func (n charNode) Codegen(scope *Scope, c *Compiler) llvm.Value {
	return llvm.ConstInt(llvm.Int8Type(), uint64(n.Value), true)
}

// String Constant Code Generator
func (n stringNode) Codegen(scope *Scope, c *Compiler) llvm.Value {
	var backingArrayPointer llvm.Value
	return backingArrayPointer
}

// Float Constant Code Generator
func (n floatNode) Codegen(scope *Scope, c *Compiler) llvm.Value {
	return llvm.ConstFloat(llvm.DoubleType(), n.Value)
}

// Variable Reference Node Code Generator
func (n variableReferenceNode) Codegen(scope *Scope, c *Compiler) llvm.Value {
	v, found := scope.Find(n.Name)
	if !found {
		fmt.Printf("unknown variable name `%s`\n", n.Name)
		os.Exit(-1)
	}
	return c.Builder.CreateLoad(v, n.Name)
}

// Variable Node Code Generator
func (n variableNode) Codegen(scope *Scope, c *Compiler) llvm.Value {
	var oldvars = []llvm.Value{}

	f := c.Builder.GetInsertBlock().Parent()
	name := n.Name
	body := n.Body
	varType := n.Type.LLVMType

	var val llvm.Value
	if body != nil {
		val = body.Codegen(scope, c)
		if val.IsNil() {
			return val // nil
		}
	}

	alloca := createEntryBlockAlloca(f, varType, n.Name)
	c.Builder.CreateStore(val, alloca)

	oldVar, _ := scope.Find(name)
	oldvars = append(oldvars, oldVar)
	scope.Set(name, alloca)

	// evaluate body now that vars are in scope
	bodyVal := body.Codegen(scope, c)
	if bodyVal.IsNil() {
		return codegenError("body returns nil") // nil
	}

	// pop old values
	// for i := range n.Vars {
	// 	Scope.Set(name, oldvars[1])
	// 	// namedVals[n.vars[i].name] = oldvars[i]
	// }

	return bodyVal
}

// Code Block Code Generator
func (n blockNode) Codegen(scope *Scope, c *Compiler) llvm.Value {
	blockScope := scope.SpawnChild()
	for _, node := range n.Nodes {
		node.Codegen(blockScope, c)
	}
	return llvm.Value{}
}

// Function Node Statement Code Generation
func (n functionNode) Codegen(scope *Scope, c *Compiler) llvm.Value {

	funcArgs := []llvm.Type{}
	for _, arg := range n.Args {
		funcArgs = append(funcArgs, arg.Type.LLVMType)
	}
	funcType := llvm.FunctionType(n.ReturnType.LLVMType, funcArgs, false)
	function := llvm.AddFunction(c.RootModule, n.Name, funcType)
	function.SetVisibility(llvm.DefaultVisibility)

	if function.Name() != n.Name {
		function.EraseFromParentAsFunction()
		function = c.RootModule.NamedFunction(n.Name)
	}

	if function.BasicBlocksCount() != 0 {
		return error("redefinition of function: " + n.Name)
	}
	// Set all the parameter names
	for i, param := range function.Params() {
		param.SetName(n.Args[i].Name)
		scope.Set(n.Args[i].Name, param)
	}

	if !n.IsExternal {
		block := llvm.AddBasicBlock(function, "entry")
		c.Builder.SetInsertPointAtEnd(block)

		// createArgAlloca(n, function, scope)

		args := function.Params()
		for i, arg := range args {
			alloca := createEntryBlockAlloca(function, arg.Type(), n.Args[i].Name)
			c.Builder.CreateStore(arg, alloca)
			scope.Set(n.Args[i].Name, alloca)
		}

		n.Body.Codegen(scope, c)

		if llvm.VerifyFunction(function, llvm.PrintMessageAction) != nil {
			function.EraseFromParentAsFunction()
			return error("function verifiction failed")
		}
	} else {
		function.SetLinkage(llvm.ExternalLinkage)
	}

	c.RootScope.Set(n.Name, function)

	return function
}

// CreateEntryBlockAlloca - Create an alloca instruction in the entry block of
// the function.  This is used for mutable variables etc.
func createEntryBlockAlloca(f llvm.Value, t llvm.Type, name string) llvm.Value {
	tmpB := llvm.NewBuilder()
	tmpB.SetInsertPoint(f.EntryBasicBlock(), f.EntryBasicBlock().FirstInstruction())
	return tmpB.CreateAlloca(t, name)
}

// Allow functions to return an error isntead of having to manage closing the program each time.
func codegenError(str string, args ...interface{}) llvm.Value {
	fmt.Fprintf(os.Stderr, "Error: %s\n", fmt.Sprintf(str, args...))
	return llvm.Value{}
}
