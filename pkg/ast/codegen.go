package ast

import (
	"fmt"
	"github.com/go-llvm/llvm"
	"os"
)

var (
	rootModule             = llvm.NewModule("act")
	rootFuncPassMgr        = llvm.NewFunctionPassManagerForModule(rootModule)
	nativeInitErr          = llvm.InitializeNativeTarget()
	builder                = llvm.NewBuilder()
	rootScope              = &Scope{nil, make([]*Scope, 0), make(map[string]llvm.Value)}
	execEngine, jitInitErr = llvm.NewMCJITCompiler(rootModule, llvm.MCJITCompilerOptions{})
)

// Scope -
type Scope struct {
	Parent   *Scope
	Children []*Scope
	Vals     map[string]llvm.Value
}

// Set -
func (s *Scope) Set(name string, val llvm.Value) {
	s.Vals[name] = val
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

func init() {
	if nativeInitErr != nil {
		fmt.Fprintln(os.Stderr, nativeInitErr)
		os.Exit(-1)
	}
	if jitInitErr != nil {
		fmt.Fprintln(os.Stderr, jitInitErr)
		os.Exit(-1)
	}

}

func error(err string) llvm.Value {
	fmt.Println(err)
	os.Exit(-1)
	return llvm.Value{}
}

// Optimize the llvm target
func Optimize() {
	rootFuncPassMgr.Add(execEngine.TargetData())
	rootFuncPassMgr.AddPromoteMemoryToRegisterPass()
	rootFuncPassMgr.AddInstructionCombiningPass()
	rootFuncPassMgr.AddReassociatePass()
	rootFuncPassMgr.AddGVNPass()
	rootFuncPassMgr.AddCFGSimplificationPass()
	rootFuncPassMgr.InitializeFunc()
}

func (n intNode) Codegen() llvm.Value               { return llvm.Value{} }
func (n stringNode) Codegen() llvm.Value            { return llvm.Value{} }
func (n ifNode) Codegen() llvm.Value                { return llvm.Value{} }
func (n forNode) Codegen() llvm.Value               { return llvm.Value{} }
func (n unaryNode) Codegen() llvm.Value             { return llvm.Value{} }
func (n binaryNode) Codegen() llvm.Value            { return llvm.Value{} }
func (n fnCallNode) Codegen() llvm.Value            { return llvm.Value{} }
func (n variableReferenceNode) Codegen() llvm.Value { return llvm.Value{} }
func (n variableNode) Codegen() llvm.Value          { return llvm.Value{} }
func (n returnNode) Codegen() llvm.Value            { return llvm.Value{} }
func (n functionCallNode) Codegen() llvm.Value      { return llvm.Value{} }
func (n blockNode) Codegen() llvm.Value             { return llvm.Value{} }
func (n whileNode) Codegen() llvm.Value             { return llvm.Value{} }

func (n floatNode) Codegen() llvm.Value {
	return llvm.ConstFloat(llvm.DoubleType(), n.Value)
}

func createEntryBlockAlloca(f llvm.Value, t llvm.Type, name string) llvm.Value {
	tmpB := llvm.NewBuilder()
	tmpB.SetInsertPoint(f.EntryBasicBlock(), f.EntryBasicBlock().FirstInstruction())
	return tmpB.CreateAlloca(t, name)
}

func (n functionNode) Codegen() llvm.Value {
	scope := GetRootScope().SpawnChild()

	funcArgs := []llvm.Type{}
	for _, arg := range n.Args {
		funcArgs = append(funcArgs, arg.Type.LLVMType)
	}
	funcType := llvm.FunctionType(n.ReturnType.LLVMType, funcArgs, false)
	function := llvm.AddFunction(rootModule, n.Name, funcType)
	function.SetVisibility(llvm.DefaultVisibility)

	if function.Name() != n.Name {
		function.EraseFromParentAsFunction()
		function = rootModule.NamedFunction(n.Name)
	}

	if function.BasicBlocksCount() != 0 {
		return error("redefinition of function: " + n.Name)
	}
	// Set all the parameter names
	for i, param := range function.Params() {
		param.SetName(n.Args[i].Name)
		scope.Set(n.Args[i].Name, param)
	}

	block := llvm.AddBasicBlock(function, "entry")
	builder.SetInsertPointAtEnd(block)

	// createArgAlloca(n, function, scope)

	args := function.Params()
	for i, arg := range args {
		alloca := createEntryBlockAlloca(function, arg.Type(), n.Args[i].Name)
		builder.CreateStore(arg, alloca)
		scope.Set(n.Args[i].Name, alloca)
	}

	n.Body.Codegen()

	if llvm.VerifyFunction(function, llvm.PrintMessageAction) != nil {
		function.EraseFromParentAsFunction()
		return error("function verifiction failed")
	}

	rootScope.Set(n.Name, function)

	return function
}

// GetRootScope returns rootScope
func GetRootScope() *Scope {
	return rootScope
}

// GetLLVMIR returns the llvm IR of the compiled program
func GetLLVMIR() string {
	return rootModule.String()
}
