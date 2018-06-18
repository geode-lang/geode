package ast

import (
	"fmt"
	"github.com/go-llvm/llvm"
	"io/ioutil"
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

// Scope trees represent block scoping by having a root scope
// and children scopes that point back to their parent scope.
type Scope struct {
	Parent   *Scope
	Children []*Scope
	Vals     map[string]llvm.Value
}

// Set a value in this specific scope
func (s *Scope) Set(name string, val llvm.Value) {
	// _, alreadyDefined := s.Vals[name]
	// if alreadyDefined {
	// 	error(fmt.Sprintf("variable `%s` already defined in scope", name))
	// }
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

func init() {
	if nativeInitErr != nil {
		fmt.Fprintln(os.Stderr, nativeInitErr)
		os.Exit(-1)
	}
	if jitInitErr != nil {
		fmt.Fprintln(os.Stderr, jitInitErr)
		os.Exit(-1)
	}

	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()

	CPU := "generic"
	features := ""
	targetTripple := llvm.DefaultTargetTriple()
	opt := llvm.CodeGenLevelNone
	reloc := llvm.RelocDefault
	model := llvm.CodeModelDefault
	target, err := llvm.GetTargetFromTriple(targetTripple)
	if err != nil {
		panic(err)
	}
	target.CreateTargetMachine(targetTripple, CPU, features, opt, reloc, model)
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

func (n ifNode) Codegen(scope *Scope) llvm.Value           { return llvm.Value{} }
func (n forNode) Codegen(scope *Scope) llvm.Value          { return llvm.Value{} }
func (n unaryNode) Codegen(scope *Scope) llvm.Value        { return llvm.Value{} }
func (n binaryNode) Codegen(scope *Scope) llvm.Value       { return llvm.Value{} }
func (n fnCallNode) Codegen(scope *Scope) llvm.Value       { return llvm.Value{} }
func (n functionCallNode) Codegen(scope *Scope) llvm.Value { return llvm.Value{} }

func (n whileNode) Codegen(scope *Scope) llvm.Value { return llvm.Value{} }

func (n returnNode) Codegen(scope *Scope) llvm.Value {
	retVal := n.Value.Codegen(scope)
	builder.CreateRet(retVal)
	return llvm.Value{}
}

func (n intNode) Codegen(scope *Scope) llvm.Value {
	return llvm.ConstInt(llvm.Int64Type(), uint64(n.Value), true)
}

var stringIndex int

func (n stringNode) Codegen(scope *Scope) llvm.Value {
	// length := len(n.Value)

	var backingArrayPointer llvm.Value

	globString := builder.CreateGlobalStringPtr(n.Value, ".str")
	fmt.Println(globString)
	// backingArray := v.Create

	return backingArrayPointer
}

func (n floatNode) Codegen(scope *Scope) llvm.Value {
	return llvm.ConstFloat(llvm.DoubleType(), n.Value)
}

func (n variableReferenceNode) Codegen(scope *Scope) llvm.Value {
	v, found := scope.Find(n.Name)
	if !found {
		fmt.Printf("unknown variable name `%s`\n", n.Name)
		os.Exit(-1)
	}
	return builder.CreateLoad(v, n.Name)
}

func (n variableNode) Codegen(scope *Scope) llvm.Value {
	var oldvars = []llvm.Value{}

	f := builder.GetInsertBlock().Parent()
	name := n.Name
	body := n.Body
	varType := n.Type.LLVMType

	var val llvm.Value
	if body != nil {
		val = body.Codegen(scope)
		if val.IsNil() {
			return val // nil
		}
	}

	alloca := createEntryBlockAlloca(f, varType, n.Name)
	builder.CreateStore(val, alloca)

	oldVar, _ := scope.Find(name)
	oldvars = append(oldvars, oldVar)
	scope.Set(name, alloca)

	// evaluate body now that vars are in scope
	bodyVal := body.Codegen(scope)
	if bodyVal.IsNil() {
		return ErrorV("body returns nil") // nil
	}

	// pop old values
	// for i := range n.Vars {
	// 	Scope.Set(name, oldvars[1])
	// 	// namedVals[n.vars[i].name] = oldvars[i]
	// }

	return bodyVal
}

func (n blockNode) Codegen(scope *Scope) llvm.Value {
	blockScope := scope.SpawnChild()

	for _, node := range n.Nodes {
		node.Codegen(blockScope)
	}

	return llvm.Value{}

}

func createEntryBlockAlloca(f llvm.Value, t llvm.Type, name string) llvm.Value {
	tmpB := llvm.NewBuilder()
	tmpB.SetInsertPoint(f.EntryBasicBlock(), f.EntryBasicBlock().FirstInstruction())
	return tmpB.CreateAlloca(t, name)
}

func (n functionNode) Codegen(scope *Scope) llvm.Value {

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

	n.Body.Codegen(scope)

	// if llvm.VerifyFunction(function, llvm.PrintMessageAction) != nil {
	// 	function.EraseFromParentAsFunction()
	// 	return error("function verifiction failed")
	// }

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

// ErrorV -
func ErrorV(str string) llvm.Value {
	fmt.Fprintf(os.Stderr, "Error: %v\n", str)
	return llvm.Value{}
}

// EmitModuleObject takes an llvm module and emits the object code
func EmitModuleObject() string {
	filename := "out.o"

	target, targetMachineTrippleError := llvm.GetTargetFromTriple(llvm.DefaultTargetTriple())
	if targetMachineTrippleError != nil {
		panic(targetMachineTrippleError)
	}

	targetMachine := target.CreateTargetMachine(llvm.DefaultTargetTriple(), "", "", llvm.CodeGenLevelNone, llvm.RelocDefault, llvm.CodeModelDefault)
	// targetData := targetMachine.TargetData()

	membuf, emitErr := targetMachine.EmitToMemoryBuffer(rootModule, llvm.ObjectFile)
	if emitErr != nil {
		panic(emitErr)
	}

	writeErr := ioutil.WriteFile(filename, membuf.Bytes(), 0666)
	if writeErr != nil {
		panic(writeErr)
	}

	return filename
}
