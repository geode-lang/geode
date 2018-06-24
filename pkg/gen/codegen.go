package gen

import (
	"fmt"
	"os"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"gitlab.com/nickwanninger/geode/pkg/typesystem"
	"gitlab.com/nickwanninger/geode/pkg/util/log"
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
	// b := make([]byte, 2)
	nameNumber++
	// rand.Read(b)
	return fmt.Sprintf("%s_%d", name, nameNumber)
}

func (n ifNode) Codegen(scope *Scope, c *Compiler) value.Value {
	predicate := n.If.Codegen(scope, c)
	one := constant.NewInt(1, types.I1)
	parentBlock := c.CurrentBlock()
	predicate = parentBlock.NewICmp(ir.IntEQ, one, createTypeCast(c, predicate, types.I1))
	parentFunc := parentBlock.Parent

	thenBlk := parentFunc.NewBlock(mangleName("then"))
	elseBlk := &ir.BasicBlock{}
	elseBlk = parentFunc.NewBlock(mangleName("else"))

	endBlk := parentFunc.NewBlock(mangleName("merge"))
	parentBlock.NewCondBr(predicate, thenBlk, elseBlk)

	c.PushBlock(thenBlk)
	n.Then.Codegen(scope, c)
	// If there is no terminator for the block, IE: no return
	// branch to the merge block
	if thenBlk.Term == nil {
		thenBlk.NewBr(endBlk)
	}
	c.PopBlock()

	// We only want to construct the else block if there is one.
	if n.Else != nil {
		c.PushBlock(elseBlk)
		n.Else.Codegen(scope, c)
		if elseBlk.Term == nil {
			elseBlk.NewBr(endBlk)
		}
		c.PopBlock()
	} else {
		// If there is no else block, just break to the merge block
		elseBlk.NewBr(endBlk)
	}

	c.PushBlock(endBlk)

	return endBlk
}

func (n forNode) Codegen(scope *Scope, c *Compiler) value.Value   { return nil }
func (n unaryNode) Codegen(scope *Scope, c *Compiler) value.Value { return nil }

func (n fnCallNode) Codegen(scope *Scope, c *Compiler) value.Value { return nil }
func (n whileNode) Codegen(scope *Scope, c *Compiler) value.Value  { return nil }

func typeSize(t types.Type) int {
	if types.IsInt(t) {
		return t.(*types.IntType).Size
	}
	if types.IsFloat(t) {
		return int(t.(*types.FloatType).Kind)
	}

	return -1
}

func binaryCast(c *Compiler, left, right value.Value) (value.Value, value.Value, types.Type) {
	// Right and Left types
	lt := left.Type()
	rt := right.Type()

	var casted types.Type

	// Get the cast precidence of both sides
	leftPrec := typesystem.CastPrecidence(lt)
	rightPrec := typesystem.CastPrecidence(rt)

	if leftPrec > rightPrec {
		casted = lt
		right = createTypeCast(c, right, lt)
	} else {
		casted = rt
		left = createTypeCast(c, left, rt)
	}
	return left, right, casted
}

// createTypeCast is where most, if not all, type casting happens in the language.
func createTypeCast(c *Compiler, in value.Value, to types.Type) value.Value {
	inType := in.Type()
	fromInt := types.IsInt(inType)
	fromFloat := types.IsFloat(inType)

	toInt := types.IsInt(to)
	toFloat := types.IsFloat(to)

	inSize := typeSize(inType)
	outSize := typeSize(to)

	if fromFloat && toInt {
		return c.CurrentBlock().NewFPToSI(in, to)
	}

	if fromInt && toFloat {
		return c.CurrentBlock().NewSIToFP(in, to)
	}

	if fromInt && toInt {
		if inSize < outSize {
			return c.CurrentBlock().NewSExt(in, to)
		}
		if inSize == outSize {
			return in
		}
		return c.CurrentBlock().NewTrunc(in, to)
	}

	if fromFloat && toFloat {
		if inSize < outSize {
			return c.CurrentBlock().NewFPExt(in, to)
		}
		if inSize == outSize {
			return in
		}
		return c.CurrentBlock().NewFPTrunc(in, to)
	}

	// If the cast would not change the type, just return the in value
	if types.Equal(inType, to) {
		return in
	}

	return codegenError("Failed to typecast")
}

func createAdd(blk *ir.BasicBlock, t types.Type, left, right value.Value) value.Value {
	if types.IsInt(t) {
		return blk.NewAdd(left, right)
	}
	if types.IsFloat(t) {
		return blk.NewFAdd(left, right)
	}
	log.Fatal("Creation of add instruction failed. `%s + %s`\n", left.Type(), right.Type())
	return nil
}

func createSub(blk *ir.BasicBlock, t types.Type, left, right value.Value) value.Value {
	if types.IsInt(t) {
		return blk.NewSub(left, right)
	}
	if types.IsFloat(t) {
		return blk.NewFSub(left, right)
	}
	log.Fatal("Creation of sub instruction failed. `%s - %s`\n", left.Type(), right.Type())
	return nil
}

func createMul(blk *ir.BasicBlock, t types.Type, left, right value.Value) value.Value {
	if types.IsInt(t) {
		return blk.NewMul(left, right)
	}
	if types.IsFloat(t) {
		return blk.NewFMul(left, right)
	}
	log.Fatal("Creation of mul instruction failed. `%s * %s`\n", left.Type(), right.Type())
	return nil
}

func createDiv(blk *ir.BasicBlock, t types.Type, left, right value.Value) value.Value {
	if types.IsInt(t) {
		return blk.NewSDiv(left, right)
	}
	if types.IsFloat(t) {
		return blk.NewFDiv(left, right)
	}
	log.Fatal("Creation of div instruction failed. `%s ÷ %s`\n", left.Type(), right.Type())
	return nil
}

func createRem(blk *ir.BasicBlock, t types.Type, left, right value.Value) value.Value {
	if types.IsInt(t) {
		return blk.NewSRem(left, right)
	}
	if types.IsFloat(t) {
		return blk.NewFRem(left, right)
	}
	log.Fatal("Creation of rem instruction failed. `%s % %s`\n", left.Type(), right.Type())
	return nil
}

func createCmp(blk *ir.BasicBlock, i ir.IntPred, f ir.FloatPred, t types.Type, left, right value.Value) value.Value {
	if types.IsInt(t) {
		return blk.NewICmp(i, left, right)
	}
	if types.IsFloat(t) {
		return blk.NewFCmp(f, left, right)
	}
	log.Fatal("Creation of rem instruction failed. `%s % %s`\n", left.Type(), right.Type())
	return nil
}
func (n binaryNode) Codegen(scope *Scope, c *Compiler) value.Value {
	// Generate the left and right nodes
	l := n.Left.Codegen(scope, c)
	r := n.Right.Codegen(scope, c)

	// Attempt to cast them with casting precidence
	// This means the operation `int + float` will cast the int to a float.
	l, r, t := binaryCast(c, l, r)

	if l == nil || r == nil {
		log.Fatal("An operand to a binart operation `%s` was nil and failed to generate\n", n.OP)
	}

	blk := c.CurrentBlock()

	switch n.OP {
	case "+":
		return createAdd(blk, t, l, r)
	case "-":
		return createSub(blk, t, l, r)
	case "*":
		return createMul(blk, t, l, r)
	case "/":
		return createDiv(blk, t, l, r)
	case "%":
		return createRem(blk, t, l, r)
	case "=":
		return createCmp(blk, ir.IntEQ, ir.FloatOEQ, t, l, r)
	case "!=":
		return createCmp(blk, ir.IntNE, ir.FloatONE, t, l, r)
	case ">":
		return createCmp(blk, ir.IntSGT, ir.FloatOGT, t, l, r)
	case ">=":
		return createCmp(blk, ir.IntSGE, ir.FloatOGE, t, l, r)
	case "<":
		return createCmp(blk, ir.IntSLT, ir.FloatOLT, t, l, r)
	case "<=":
		return createCmp(blk, ir.IntSLE, ir.FloatOLE, t, l, r)
	default:
		return codegenError("invalid binary operator")
	}
}

// Function Call statement Code Generator
func (n functionCallNode) Codegen(scope *Scope, c *Compiler) value.Value {
	callee := c.Functions[n.Name]

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
	retValCoerced := createTypeCast(c, retVal, c.FN.Sig.Ret)
	c.CurrentBlock().NewRet(retValCoerced)
	return nil
}

// Int Code Generator
func (n intNode) Codegen(scope *Scope, c *Compiler) value.Value {
	// return llvm.ConstInt(llvm.Int64Type(), , true)
	return constant.NewInt(n.Value, types.I64)
}

// Char Code Generator
func (n charNode) Codegen(scope *Scope, c *Compiler) value.Value { return nil }

func newCharArray(s string) *constant.Array {
	var bs []constant.Constant
	for i := 0; i < len(s); i++ {
		b := constant.NewInt(int64(s[i]), types.I8)
		bs = append(bs, b)
	}
	bs = append(bs, constant.NewInt(int64(0), types.I8))
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
	f := c.CurrentBlock().Parent
	name := n.Name
	body := n.Body

	var val value.Value
	if body != nil {
		val = body.Codegen(scope, c)
		if val == nil {
			return val // nil
		}
	}

	var alloc *ir.InstAlloca
	if n.Reassignment {
		v, found := scope.Find(name)
		if !found {
			fmt.Println(v, "Not found")
		}
		alloc = v.(*ir.InstAlloca)
	} else {
		alloc = createBlockAlloca(f, n.Type, name)
		scope.Set(name, alloc)
	}

	val = createTypeCast(c, val, alloc.Elem)

	c.CurrentBlock().NewStore(val, alloc)
	return nil
}

var nameNumber int

// Code Block Code Generator
func (n blockNode) Codegen(scope *Scope, c *Compiler) value.Value {
	blockScope := scope.SpawnChild()
	// spew.Dump(n.Nodes)
	for _, node := range n.Nodes {
		node.Codegen(blockScope, c)
	}

	// spew.Dump(c.CurrentBlock())

	// if c.CurrentBlock().Term == nil {
	// 	c.CurrentBlock().NewRet(constant.NewInt(0, types.Void))
	// }
	return nil
}

// Function Node Statement Code Generation
func (n functionNode) Codegen(scope *Scope, c *Compiler) value.Value {

	funcArgs := make([]*types.Param, 0)
	for _, arg := range n.Args {
		p := ir.NewParam(arg.Name, arg.Type)
		funcArgs = append(funcArgs, p)

	}

	function := c.RootModule.NewFunction(n.Name, n.ReturnType, funcArgs...)

	c.FN = function
	// Set the function name map to the function call
	c.Functions[n.Name] = function
	name := mangleName("entry")
	c.PushBlock(c.FN.NewBlock(name))

	for _, arg := range function.Params() {
		alloc := c.CurrentBlock().NewAlloca(arg.Type())
		c.CurrentBlock().NewStore(arg, alloc)
		scope.Set(arg.Name, alloc)
	}

	n.Body.Codegen(scope, c)

	// Given that we are at the end of the function's codegen phase
	// we can assume the current block is the last block. So we can
	// make sure it has a terminator, and if there isn't one we need
	// to create one only if the return type of the function is void
	// otherwise we need to Error fatally

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
	// 	block := llvm.AddBasicBlock(fgiunction, "entry")
	// 	c.Builder.SetInsertPointAtEnd(block)

	// 	args := function.Params()
	// 	for i, arg := range args {
	// 		alloca := createBlockAlloca(function, arg.Type(), n.Args[i].Name)
	// 		c.Builder.CreateStore(arg, alloca)
	// 		scope.Set(n.Args[i].Name, alloca)
	// 	}
	// }

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
func createBlockAlloca(f *ir.Function, elemType types.Type, name string) *ir.InstAlloca {
	// Create a new allocation in the root of the function
	alloca := f.Blocks[0].NewAlloca(elemType)
	// Set the name of the allocation (the variable name)
	alloca.SetName(name)
	return alloca
}

// Allow functions to return an error isntead of having to manage closing the program each time.
func codegenError(str string, args ...interface{}) value.Value {
	fmt.Fprintf(os.Stderr, "Error: %s\n", fmt.Sprintf(str, args...))
	return nil
}
