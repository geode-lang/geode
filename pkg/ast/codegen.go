package ast

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	"github.com/nickwanninger/geode/pkg/typesystem"
	"github.com/nickwanninger/geode/pkg/util/log"
)

// A global number to indicate which `name index` we are on. This way,
// the mangler will never output the same name twice as this number is monotonic
var nameNumber int

func mangleName(name string) string {
	nameNumber++
	return fmt.Sprintf("%s_%d", name, nameNumber)
}

func branchIfNoTerminator(blk *ir.BasicBlock, to *ir.BasicBlock) {
	if blk.Term == nil {
		blk.NewBr(to)
	}
}

// Handle will do ast-level handling for a dependency node
func (n DependencyNode) Handle(c *Compiler) value.Value {

	abs, _ := filepath.Abs(c.Package.Source.Path)
	dir := path.Dir(abs)

	for _, dp := range n.Paths {

		depPath := path.Join(dir, dp)
		if n.CLinkage {
			c.Package.AddClinkage(depPath)
		} else {
			c.Package.LoadDep(depPath)
		}

	}

	return nil
}

// Codegen implements Node.Codegen for DependencyNode
func (n DependencyNode) Codegen(scope *Scope, c *Compiler) value.Value { return nil }

// Codegen implements Node.Codegen for IfNode
func (n IfNode) Codegen(scope *Scope, c *Compiler) value.Value {
	predicate := n.If.Codegen(scope, c)
	one := constant.NewInt(1, types.I1)
	// The name of the blocks is prefixed because
	namePrefix := fmt.Sprintf("if.%d.", n.Index)
	parentBlock := c.CurrentBlock()
	predicate = parentBlock.NewICmp(ir.IntEQ, one, createTypeCast(c, predicate, types.I1))
	parentFunc := parentBlock.Parent

	var thenGenBlk *ir.BasicBlock
	var endBlk *ir.BasicBlock

	thenBlk := parentFunc.NewBlock(mangleName(namePrefix + "then"))

	c.genInBlock(thenBlk, func() {
		thenGenBlk = n.Then.Codegen(scope, c).(*ir.BasicBlock)
	})

	elseBlk := parentFunc.NewBlock(mangleName(namePrefix + "else"))
	var elseGenBlk *ir.BasicBlock

	c.genInBlock(elseBlk, func() {
		// We only want to construct the else block if there is one.
		if n.Else != nil {
			elseGenBlk = n.Else.Codegen(scope, c).(*ir.BasicBlock)
		}
	})

	endBlk = parentFunc.NewBlock(mangleName(namePrefix + "end"))
	c.PushBlock(endBlk)
	// We need to make sure these blocks have terminators.
	// in order to do that, we branch to the end block
	branchIfNoTerminator(thenBlk, endBlk)
	branchIfNoTerminator(thenGenBlk, endBlk)
	branchIfNoTerminator(elseBlk, endBlk)

	if elseGenBlk != nil {
		branchIfNoTerminator(elseGenBlk, endBlk)
	}

	parentBlock.NewCondBr(predicate, thenBlk, elseBlk)

	// branchIfNoTerminator(c.CurrentBlock(), endBlk)

	return endBlk
}

// Codegen implements Node.Codegen for ForNode
func (n ForNode) Codegen(scope *Scope, c *Compiler) value.Value {
	// The name of the blocks is prefixed so we can determine which for loop a block is for.
	namePrefix := fmt.Sprintf("for.%X.", n.Index)
	parentBlock := c.CurrentBlock()
	var predicate value.Value
	var condBlk *ir.BasicBlock
	var bodyBlk *ir.BasicBlock
	var bodyGenBlk *ir.BasicBlock
	var endBlk *ir.BasicBlock
	parentFunc := parentBlock.Parent

	condBlk = parentFunc.NewBlock(namePrefix + "cond")

	n.Init.Codegen(scope, c)

	parentBlock.NewBr(condBlk)

	c.genInBlock(condBlk, func() {
		predicate = n.Cond.Codegen(scope, c)
		one := constant.NewInt(1, types.I1)
		predicate = condBlk.NewICmp(ir.IntEQ, one, createTypeCast(c, predicate, types.I1))
	})
	bodyBlk = parentFunc.NewBlock(namePrefix + "body")
	c.genInBlock(bodyBlk, func() {
		bodyGenBlk = n.Body.Codegen(scope, c).(*ir.BasicBlock)

		c.genInBlock(bodyGenBlk, func() {
			n.Step.Codegen(scope, c)
		})
		branchIfNoTerminator(bodyBlk, condBlk)
		branchIfNoTerminator(bodyGenBlk, condBlk)
	})
	endBlk = parentFunc.NewBlock(namePrefix + "end")
	c.PushBlock(endBlk)
	condBlk.NewCondBr(predicate, bodyBlk, endBlk)
	return endBlk
}

// Codegen implements Node.Codegen for CharNode
func (n CharNode) Codegen(scope *Scope, c *Compiler) value.Value { return nil }

// Codegen implements Node.Codegen for UnaryNode
func (n UnaryNode) Codegen(scope *Scope, c *Compiler) value.Value {

	operandValue := n.Operand.Codegen(scope, c)
	if operandValue == nil {
		log.Fatal("nil operand")
	}

	// handle reference operation
	if n.Operator == "&" {
		return operandValue
	}
	// handle dereference operation
	if n.Operator == "*" {
		return c.CurrentBlock().NewLoad(operandValue)
	}

	// fmt.Println(n.Name, operandValue)
	return operandValue
}

// Codegen implements Node.Codegen for WhileNode
func (n WhileNode) Codegen(scope *Scope, c *Compiler) value.Value {

	// The name of the blocks is prefixed because
	namePrefix := fmt.Sprintf("while_%d_", n.Index)
	parentBlock := c.CurrentBlock()

	parentFunc := parentBlock.Parent
	startblock := parentFunc.NewBlock(mangleName(namePrefix + "start"))
	c.PushBlock(startblock)
	predicate := n.If.Codegen(scope, c)
	one := constant.NewInt(1, types.I1)
	c.PopBlock()
	branchIfNoTerminator(parentBlock, startblock)
	predicate = startblock.NewICmp(ir.IntEQ, one, createTypeCast(c, predicate, types.I1))

	var endBlk *ir.BasicBlock

	bodyBlk := parentFunc.NewBlock(mangleName(namePrefix + "body"))
	c.PushBlock(bodyBlk)
	bodyGenBlk := n.Body.Codegen(scope, c).(*ir.BasicBlock)

	// If there is no terminator for the block, IE: no return
	// branch to the merge block

	endBlk = parentFunc.NewBlock(mangleName(namePrefix + "merge"))
	c.PushBlock(endBlk)

	branchIfNoTerminator(bodyBlk, startblock)
	branchIfNoTerminator(bodyGenBlk, startblock)

	startblock.NewCondBr(predicate, bodyBlk, endBlk)

	// branchIfNoTerminator(c.CurrentBlock(), endBlk)

	return endBlk
}

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

	if types.Equal(to, types.Void) {
		return nil
	}

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

	log.Fatal("Failed to typecast type %s to %s\n", inType.String(), to)
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

// CreateBinaryOp produces a geode binary op (just a wrapper around llir/llvm's binary instructions)
func CreateBinaryOp(intstr, fltstr string, blk *ir.BasicBlock, t types.Type, left, right value.Value) value.Value {
	var inst *GeodeBinaryInstr
	if types.IsInt(t) {
		inst = NewGeodeBinaryInstr(intstr, left, right)
	} else {
		inst = NewGeodeBinaryInstr(fltstr, left, right)
	}
	blk.AppendInst(inst)
	return inst
}

// Codegen implements Node.Codegen for BinaryNode
func (n BinaryNode) Codegen(scope *Scope, c *Compiler) value.Value {
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
		return CreateBinaryOp("add", "fadd", blk, t, l, r)
	case "-":
		return CreateBinaryOp("sub", "fsub", blk, t, l, r)
	case "*":
		return CreateBinaryOp("mul", "fmul", blk, t, l, r)
	case "/":
		return CreateBinaryOp("sdiv", "fdiv", blk, t, l, r)
	case "%":
		return CreateBinaryOp("srem", "frem", blk, t, l, r)
	case ">>":
		return CreateBinaryOp("lshr", "lshr", blk, t, l, r)
	case "<<":
		return CreateBinaryOp("shl", "shl", blk, t, l, r)
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

// Codegen implements Node.Codegen for CastNode
func (n CastNode) Codegen(scope *Scope, c *Compiler) value.Value {
	from := n.From.Codegen(scope, c)
	return createTypeCast(c, from, n.InferType(scope))
}

// Codegen implements Node.Codegen for FunctionCallNode
func (n FunctionCallNode) Codegen(scope *Scope, c *Compiler) value.Value {

	// scopeItem, found := c.Scope.Find(n.Name)

	args := []value.Value{}
	argTypes := []types.Type{}
	for _, arg := range n.Args {
		a := arg.Codegen(scope, c)
		// fmt.Println(a.Type())
		args = append(args, a)
		argTypes = append(argTypes, a.Type())
		if args[len(args)-1] == nil {
			return codegenError(fmt.Sprintf("Argument to function %q failed to generate code", n.Name))
		}
	}
	name := MangleFunctionName(n.Name, argTypes...)
	functionOptions := c.Scope.FindFunctions(name)
	funcCount := len(functionOptions)

	if funcCount > 1 {
		log.Fatal("Too many options for function call '%s'\n", name)
	} else if funcCount == 0 {
		log.Fatal("Unable to find function '%s' in scope of module '%s'\n", name, c.Name)
	}

	fnScopeItem := functionOptions[0]
	if !c.FunctionDefined(fnScopeItem.function) {
		c.Module.AppendFunction(fnScopeItem.function)
	}

	callee := fnScopeItem.Value().(*ir.Function)
	if callee == nil {
		return codegenError(fmt.Sprintf("Unknown function %q referenced", name))
	}

	// fmt.Println(n.Name, callee.Type())

	return c.CurrentBlock().NewCall(callee, args...)
}

// Codegen implements Node.Codegen for ReturnNode
func (n ReturnNode) Codegen(scope *Scope, c *Compiler) value.Value {
	var retVal value.Value

	if c.FN.Sig.Ret != types.Void {
		if n.Value != nil {
			retVal = n.Value.Codegen(scope, c)
			// retVal = createTypeCast(c, retVal, c.FN.Sig.Ret)
		} else {
			retVal = nil
		}
	}

	c.CurrentBlock().NewRet(retVal)

	return retVal
}

// Codegen implements Node.Codegen for IntNode
func (n IntNode) Codegen(scope *Scope, c *Compiler) value.Value {
	// return llvm.ConstInt(llvm.Int64Type(), , true)
	return constant.NewInt(n.Value, types.I64)
}

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

func canBeIndexed(val value.Value) bool {
	return types.IsArray(val.Type()) || types.IsPointer(val.Type())
}

// Codegen implements Node.Codegen for StringNode
func (n StringNode) Codegen(scope *Scope, c *Compiler) value.Value {
	str := c.Module.NewGlobalDef(mangleName(".str"), newCharArray(n.Value))
	str.IsConst = true
	zero := constant.NewInt(0, types.I32)
	return constant.NewGetElementPtr(str, zero, zero)
}

// Codegen implements Node.Codegen for FloatNode
func (n FloatNode) Codegen(scope *Scope, c *Compiler) value.Value {
	return constant.NewFloat(n.Value, types.Double)
}

// Codegen implements Node.Codegen for VariableNode
func (n VariableNode) Codegen(scope *Scope, c *Compiler) value.Value {
	block := c.CurrentBlock()
	f := block.Parent

	name := n.Name
	var alloc *ir.InstAlloca
	var val value.Value

	if n.RefType == ReferenceAccess || n.RefType == ReferenceDereference {
		v, found := scope.Find(name)
		if !found {
			fmt.Printf("unknown variable name `%s`\n", name)
			os.Exit(-1)
		}

		alloc = v.Value().(*ir.InstAlloca)

		if n.RefType == ReferenceDereference {
			return alloc
		}

		val = block.NewLoad(v.Value())

		if n.IndexExpr != nil {
			if types.IsPointer(val.Type()) {
				// zero := constant.NewInt(0, types.I32)
				index := n.IndexExpr.Codegen(scope, c)
				ptr := block.NewGetElementPtr(val, index)
				val = block.NewLoad(ptr)
			}

		}
		return val
	}

	if n.RefType == ReferenceAssign {
		v, found := scope.Find(name)
		if !found {
			fmt.Println(v, "Not found")
		}
		alloc = v.Value().(*ir.InstAlloca)
	} else if n.RefType == ReferenceDefine {
		ty := scope.FindType(n.Type.Name).Type
		ty = n.Type.BuildPointerType(ty)
		alloc = createBlockAlloca(f, ty, name)
		scItem := NewVariableScopeItem(n.Name, alloc, PrivateVisibility)
		scope.Add(scItem)
	}

	if n.HasValue {
		// Construct the body
		if n.Body != nil {
			val = n.Body.Codegen(scope, c)
			if val == nil {
				return val // nil
			}
		}
		val = createTypeCast(c, val, alloc.Elem)

	} else {
		// Default to 0 from issue:
		// https://github.com/nickwanninger/geode/issues/5
		// val = createTypeCast(c, constant.NewInt(0, types.I64), alloc.Elem)
		// val = nil
		return nil
	}
	block.NewStore(val, alloc)

	return nil
}

// Codegen implements Node.Codegen for BlockNode
func (n BlockNode) Codegen(scope *Scope, c *Compiler) value.Value {
	blockScope := scope.SpawnChild()
	// spew.Dump(n.Nodes)
	for _, node := range n.Nodes {
		node.Codegen(blockScope, c)
	}
	return c.CurrentBlock()
}

// Codegen implements Node.Codegen for ClassNode
func (n ClassNode) Codegen(scope *Scope, c *Compiler) value.Value {

	fields := make([]types.Type, 0)

	for _, f := range n.Variables {
		t := f.Type.Name
		ty := scope.FindType(t).Type
		fields = append(fields, ty)
	}

	structDefn := types.NewStruct(fields...)

	NewTypeDef(n.Name, structDefn, -1).InjectInto(scope)
	// structDefn.Opaque = true
	structDefn.SetName(n.Name)
	c.Module.NewType(n.Name, structDefn)
	// fmt.Println(t, structDefn)
	return nil
}

// Arguments returns some FunctionNode's arguments
func (n FunctionNode) Arguments(scope *Scope) ([]*types.Param, []types.Type) {
	funcArgs := make([]*types.Param, 0)
	argTypes := make([]types.Type, 0)
	for _, arg := range n.Args {
		ty := scope.FindType(arg.Type.Name).Type
		ty = arg.Type.BuildPointerType(ty)
		p := ir.NewParam(arg.Name, ty)
		funcArgs = append(funcArgs, p)
		argTypes = append(argTypes, p.Type())
	}
	return funcArgs, argTypes
}

// Declare declares some FunctionNode's sig
func (n FunctionNode) Declare(scope *Scope, c *Compiler) *ir.Function {

	funcArgs, argTypes := n.Arguments(scope)
	// We need to do some special checks if the function is main. It's special.
	// For instance, it must return int type.
	if n.Name == "main" {
		if n.ReturnType.Name != "int" {
			log.Fatal("Main function must return type int. Called for type '%s'\n", n.ReturnType)
		}
	}

	name := n.Name
	if !n.Nomangle {
		name = MangleFunctionName(n.Name, argTypes...)
	}
	ty := scope.FindType(n.ReturnType.Name).Type
	ty = n.ReturnType.BuildPointerType(ty)
	function := c.Module.NewFunction(name, ty, funcArgs...)

	c.FN = function

	if n.Variadic && !n.External {
		log.Fatal("Function '%s' is variadic and has a body. This only allowed for external functions.\n", n.Name)
	}

	function.Sig.Variadic = n.Variadic

	// fmt.Println(function.Name, function.Sig.Variadic)
	scopeItem := NewFunctionScopeItem(name, function, PublicVisibility)
	scopeItem.SetMangled(!n.Nomangle)
	c.Scope.Add(scopeItem)

	return function
}

// Codegen implements Node.Codegen for FunctionNode
func (n FunctionNode) Codegen(scope *Scope, c *Compiler) value.Value {

	_, argTypes := n.Arguments(scope)

	name := n.Name
	if !n.Nomangle {
		name = MangleFunctionName(n.Name, argTypes...)
	}

	declared := c.Scope.FindFunctions(name)
	if len(declared) != 1 {
		log.Fatal("Unable to find function declaration for '%s'\n", name)
	}

	// fmt.Println(c.Scope.Vals)

	function := declared[0].Value().(*ir.Function)
	c.FN = function

	// If the function is external (has ... at the end) we don't build a block
	if !n.External {
		name := mangleName("entry")
		c.PushBlock(c.FN.NewBlock(name))

		for _, arg := range function.Params() {
			alloc := c.CurrentBlock().NewAlloca(arg.Type())
			c.CurrentBlock().NewStore(arg, alloc)
			// Set the scope item
			scItem := NewVariableScopeItem(arg.Name, alloc, PrivateVisibility)
			scope.Add(scItem)
		}
		// Gen the body of the function
		n.Body.Codegen(scope, c)
		if c.CurrentBlock().Term == nil {
			ty := scope.FindType(n.ReturnType.Name).Type
			// log.Warn("Function %s is missing a return statement in the root block. Defaulting to 0\n", n.Name)
			v := createTypeCast(c, constant.NewInt(0, types.I64), ty)
			c.CurrentBlock().NewRet(v)
		}
		c.PopBlock()
	}

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
