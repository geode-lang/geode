package ast

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/constant"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"

	"github.com/geode-lang/geode/pkg/util"
	"github.com/geode-lang/geode/pkg/util/log"
)

func parseName(combined string) (string, string) {
	var namespace, name string
	parts := strings.Split(combined, ":")
	name = parts[len(parts)-1]
	if len(parts) > 1 {
		namespace = parts[0]
	}

	return namespace, name
}

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

// Codegen returns some NamespaceNode's arguments
func (n NamespaceNode) Codegen(prog *Program) value.Value { return nil }

// Handle will do ast-level handling for a dependency node
func (n DependencyNode) Handle(prog *Program) value.Value {

	return nil
}

// Codegen implements Node.Codegen for DependencyNode
func (n DependencyNode) Codegen(prog *Program) value.Value { return nil }

// Codegen implements Node.Codegen for IfNode
func (n IfNode) Codegen(prog *Program) value.Value {

	predicate := n.If.Codegen(prog)
	zero := constant.NewInt(0, types.I32)
	// The name of the blocks is prefixed because
	namePrefix := fmt.Sprintf("if.%d.", n.Index)
	parentBlock := prog.Compiler.CurrentBlock()
	predicate = parentBlock.NewICmp(ir.IntNE, zero, createTypeCast(prog, predicate, types.I32))
	parentFunc := parentBlock.Parent

	var thenGenBlk *ir.BasicBlock
	var endBlk *ir.BasicBlock

	thenBlk := parentFunc.NewBlock(mangleName(namePrefix + "then"))

	prog.Compiler.genInBlock(thenBlk, func() {
		thenGenBlk = n.Then.Codegen(prog).(*ir.BasicBlock)
	})

	elseBlk := parentFunc.NewBlock(mangleName(namePrefix + "else"))
	var elseGenBlk *ir.BasicBlock

	prog.Compiler.genInBlock(elseBlk, func() {
		// We only want to construct the else block if there is one.
		if n.Else != nil {
			elseGenBlk = n.Else.Codegen(prog).(*ir.BasicBlock)
		}
	})

	endBlk = parentFunc.NewBlock(mangleName(namePrefix + "end"))
	prog.Compiler.PushBlock(endBlk)
	// We need to make sure these blocks have terminators.
	// in order to do that, we branch to the end block
	branchIfNoTerminator(thenBlk, endBlk)
	branchIfNoTerminator(thenGenBlk, endBlk)
	branchIfNoTerminator(elseBlk, endBlk)

	if elseGenBlk != nil {
		branchIfNoTerminator(elseGenBlk, endBlk)
	}

	parentBlock.NewCondBr(predicate, thenBlk, elseBlk)

	return endBlk
}

// Codegen implements Node.Codegen for ForNode
func (n ForNode) Codegen(prog *Program) value.Value {

	// The name of the blocks is prefixed so we can determine which for loop a block is for.
	namePrefix := fmt.Sprintf("for.%X.", n.Index)
	parentBlock := prog.Compiler.CurrentBlock()
	prog.Scope = prog.Scope.SpawnChild()

	var predicate value.Value
	var condBlk *ir.BasicBlock
	var bodyBlk *ir.BasicBlock
	var bodyGenBlk *ir.BasicBlock
	var endBlk *ir.BasicBlock
	parentFunc := parentBlock.Parent

	condBlk = parentFunc.NewBlock(namePrefix + "cond")

	n.Init.Codegen(prog)

	parentBlock.NewBr(condBlk)

	prog.Compiler.genInBlock(condBlk, func() {
		predicate = n.Cond.Codegen(prog)
		one := constant.NewInt(1, types.I1)
		predicate = condBlk.NewICmp(ir.IntEQ, one, createTypeCast(prog, predicate, types.I1))
	})
	bodyBlk = parentFunc.NewBlock(namePrefix + "body")
	prog.Compiler.genInBlock(bodyBlk, func() {
		bodyGenBlk = n.Body.Codegen(prog).(*ir.BasicBlock)

		prog.Compiler.genInBlock(bodyGenBlk, func() {
			n.Step.Codegen(prog)
		})
		branchIfNoTerminator(bodyBlk, condBlk)
		branchIfNoTerminator(bodyGenBlk, condBlk)
	})
	endBlk = parentFunc.NewBlock(namePrefix + "end")
	prog.Compiler.PushBlock(endBlk)
	condBlk.NewCondBr(predicate, bodyBlk, endBlk)

	prog.Scope = prog.Scope.Parent
	return endBlk
}

// Codegen implements Node.Codegen for CharNode
func (n CharNode) Codegen(prog *Program) value.Value {
	return constant.NewInt(int64(n.Value), types.I8)
}

// GenAccess returns the value from a given CharNode
func (n CharNode) GenAccess(prog *Program) value.Value {
	return n.Codegen(prog)
}

// Codegen implements Node.Codegen for UnaryNode
func (n UnaryNode) Codegen(prog *Program) value.Value {

	operandValue := n.Operand.Codegen(prog)
	if operandValue == nil {
		n.Operand.SyntaxError()
		log.Fatal("nil operand")
	}

	if n.Operator == "-" {

		if types.IsFloat(operandValue.Type()) {
			return prog.Compiler.CurrentBlock().NewFSub(constant.NewFloat(0, types.Double), operandValue)
		} else if types.IsInt(operandValue.Type()) {
			return prog.Compiler.CurrentBlock().NewSub(constant.NewInt(0, types.I64), operandValue)
		}
		n.SyntaxError()
		log.Fatal("Unable to make a non integer/float into a negative\n")

	}

	// handle reference operation
	if n.Operator == "&" {
		return operandValue
	}
	// handle dereference operation
	if n.Operator == "*" {
		if types.IsPointer(operandValue.Type()) {
			return prog.Compiler.CurrentBlock().NewLoad(operandValue)
		}
		n.SyntaxError()
		log.Fatal("attempt to dereference a non-pointer variable\n")
	}

	return operandValue
}

// GenAccess implements Accessable.GenAccess
func (n UnaryNode) GenAccess(prog *Program) value.Value {
	return n.Codegen(prog)
}

// Codegen implements Node.Codegen for WhileNode
func (n WhileNode) Codegen(prog *Program) value.Value {

	// The name of the blocks is prefixed because
	namePrefix := fmt.Sprintf("while_%d_", n.Index)
	parentBlock := prog.Compiler.CurrentBlock()

	parentFunc := parentBlock.Parent
	startblock := parentFunc.NewBlock(mangleName(namePrefix + "start"))
	prog.Compiler.PushBlock(startblock)
	predicate := n.If.Codegen(prog)
	one := constant.NewInt(1, types.I1)
	prog.Compiler.PopBlock()
	branchIfNoTerminator(parentBlock, startblock)
	predicate = startblock.NewICmp(ir.IntEQ, one, createTypeCast(prog, predicate, types.I1))

	var endBlk *ir.BasicBlock

	bodyBlk := parentFunc.NewBlock(mangleName(namePrefix + "body"))
	prog.Compiler.PushBlock(bodyBlk)
	bodyGenBlk := n.Body.Codegen(prog).(*ir.BasicBlock)

	// If there is no terminator for the block, IE: no return
	// branch to the merge block

	endBlk = parentFunc.NewBlock(mangleName(namePrefix + "merge"))
	prog.Compiler.PushBlock(endBlk)

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

func binaryCast(prog *Program, left, right value.Value) (value.Value, value.Value, types.Type) {
	// Right and Left types
	lt := left.Type()
	rt := right.Type()

	var casted types.Type

	// Get the cast precidence of both sides
	leftPrec := prog.CastPrecidence(lt)
	rightPrec := prog.CastPrecidence(rt)

	if leftPrec > rightPrec {
		casted = lt
		right = createTypeCast(prog, right, lt)
	} else {
		casted = rt
		left = createTypeCast(prog, left, rt)
	}
	return left, right, casted
}

func typesAreLooselyEqual(a, b types.Type) bool {
	return types.IsNumber(a) && types.IsNumber(b)
}

// createTypeCast is where most, if not all, type casting happens in the language.
func createTypeCast(prog *Program, in value.Value, to types.Type) value.Value {

	inType := in.Type()
	fromInt := types.IsInt(inType)
	fromFloat := types.IsFloat(inType)

	toInt := types.IsInt(to)
	toFloat := types.IsFloat(to)

	inSize := typeSize(inType)
	outSize := typeSize(to)

	if c, ok := in.(*constant.Int); ok && types.IsInt(to) {
		c.Typ = to.(*types.IntType)
		return c
	}

	if c, ok := in.(*constant.Float); ok && types.IsFloat(to) {
		c.Typ = to.(*types.FloatType)
		return c
	}

	if types.Equal(to, types.Void) {
		return nil
	}

	if types.IsPointer(inType) && types.IsPointer(to) {
		return prog.Compiler.CurrentBlock().NewBitCast(in, to)
	}

	if fromFloat && toInt {
		return prog.Compiler.CurrentBlock().NewFPToSI(in, to)
	}

	if fromInt && toFloat {
		return prog.Compiler.CurrentBlock().NewSIToFP(in, to)
	}

	if fromInt && toInt {
		if inSize < outSize {
			return prog.Compiler.CurrentBlock().NewSExt(in, to)
		}
		if inSize == outSize {
			return in
		}
		return prog.Compiler.CurrentBlock().NewTrunc(in, to)
	}

	if fromFloat && toFloat {
		if inSize < outSize {
			return prog.Compiler.CurrentBlock().NewFPExt(in, to)
		}
		if inSize == outSize {
			return in
		}
		return prog.Compiler.CurrentBlock().NewFPTrunc(in, to)
	}

	// If the cast would not change the type, just return the in value
	if types.Equal(inType, to) {
		return in
	}

	if types.IsPointer(inType) && types.IsInt(to) {
		return prog.Compiler.CurrentBlock().NewPtrToInt(in, to)
	}

	log.Fatal("Failed to typecast type %s to %s\n", inType.String(), to)
	return nil
}

// func (n FunctionCallNode) GenAccess(prog *Program) value.Value {
// 	return n.Codegen(prog)
// }

// Codegen implements Node.Codegen for FunctionCallNode
func (n FunctionCallNode) Codegen(prog *Program) value.Value {

	// scopeItem, found := c.Scope.Find(n.Name)

	args := []value.Value{}
	argTypes := []types.Type{}
	// argStrings := []string{}
	for _, arg := range n.Args {
		if ac, isAccessable := arg.(Accessable); isAccessable {
			val := ac.GenAccess(prog)

			args = append(args, val)
			argTypes = append(argTypes, val.Type())
			if args[len(args)-1] == nil {
				return codegenError(fmt.Sprintf("Argument to function %q failed to generate code", n.Name))
			}
		} else {
			arg.SyntaxError()
			log.Fatal("Argument to function call to '%s' is not accessable (has no readable value). Node type %s\n", n.Name, arg.Kind())
		}
	}

	// First we need to check if the function call is actually a call to a class's constructor.
	// Because in geode, calling a class name results in the constructor being called for said class.
	// class := c.Scope.FindType(n.Name.String())
	// if class != nil {
	// 	return GenerateClassConstruction(n.Name.String(), class.Type, c.Scope, c, args)
	// }

	// name := n.Name

	// fmt.Println(n.Name, reflect.TypeOf(n.Name).Name())

	var searchNames []string

	// originalName := n.Name.String()

	// dotParts := strings.Split(originalName, ".")

	var name string

	switch v := n.Name.(type) {
	case DotReference:
		// t := v.BaseType(prog)
		// fmt.Println(v.Base.Alloca(prog))
		for name, val := range prog.Scope.Types {
			fmt.Println(name, val.Type)
		}

		os.Exit(1)
		break
	case NamedReference:

		ns, nm := parseName(v.String())
		name = nm
		if ns == "" {
			ns = prog.Scope.PackageName
		} else if !prog.Package.HasAccessToPackage(ns) {
			n.SyntaxError()
			log.Fatal("Package %s doesn't load package %s but attempts to call %s:%s.\n", prog.Scope.PackageName, ns, ns, nm)
		}
		searchNames = []string{
			nm,
			fmt.Sprintf("%s:%s", ns, nm),
			fmt.Sprintf("%s:%s", prog.Scope.PackageName, nm),
		}
	default:
		log.Fatal("Unknown type node passed to a function call %q\n", reflect.TypeOf(n.Name).Name())
	}

	var callee *ir.Function
	for _, name := range searchNames {
		compOpts := FunctionCompilationOptions{}
		compOpts.ArgTypes = argTypes
		callee = prog.CompileFunction(name, compOpts)
		if callee != nil {
			break
		}
	}

	if callee == nil {
		log.Fatal("Unknown function %q referenced\n", name)
	}

	// Attempt to typecast all the args into the correct type
	// This is skipped with variadic functions
	if !callee.Sig.Variadic {
		for i := range args {
			// fmt.Println(callee.Sig.Params[i])
			args[i] = createTypeCast(prog, args[i], callee.Sig.Params[i].Type())
		}
	}

	// Varargs require type conversion to a standardized type
	// So we will use the same type promotion c uses
	//  if int && type != i32 -> type = i32
	//  if fnn && type != f64 -> type = f64
	arguments := make([]value.Value, 0, len(args))

	for i, arg := range args {

		if callee.Sig.Variadic && i >= len(callee.Params()) {
			if types.IsInt(arg.Type()) {
				if !types.Equal(arg.Type(), types.I32) {
					arguments = append(arguments, createTypeCast(prog, arg, types.I32))
					continue
				}

			}

			if types.IsFloat(arg.Type()) {
				if !types.Equal(arg.Type(), types.Double) {
					arguments = append(arguments, createTypeCast(prog, arg, types.Double))
					continue
				}
			}
		}
		arguments = append(arguments, arg)
	}

	// prog.Compiler.CurrentBlock().AppendInst(NewLLVMComment("%s", n.String()))
	return prog.Compiler.CurrentBlock().NewCall(callee, arguments...)
}

// Codegen implements Node.Codegen for ReturnNode
func (n ReturnNode) Codegen(prog *Program) value.Value {

	var retVal value.Value

	if prog.Compiler.FN.Sig.Ret != types.Void {
		if n.Value != nil {
			retVal = n.Value.Codegen(prog)
			given := retVal.Type()
			expected := prog.Compiler.FN.Sig.Ret
			if !types.Equal(given, expected) {
				if !(types.IsInt(given) && types.IsInt(expected)) {
					n.SyntaxError()
					fnName := UnmangleFunctionName(prog.Compiler.FN.Name)
					expectedName, err := prog.Scope.FindTypeName(expected)
					util.EatError(err)
					givenName, err := prog.Scope.FindTypeName(given)
					util.EatError(err)
					n.SyntaxError()
					log.Fatal("Incorrect return value for function %s. Expected: %s (%s). Given: %s (%s)\n", fnName, expectedName, expected, givenName, given)
				} else {
					retVal = createTypeCast(prog, retVal, prog.Compiler.FN.Sig.Ret)
				}
			}
		} else {
			retVal = nil
		}
	}

	prog.Compiler.CurrentBlock().NewRet(retVal)

	return retVal
}

func newCharArray(s string) *constant.Array {
	var bs []constant.Constant
	for i := 0; i < len(s); i++ {
		b := constant.NewInt(int64(s[i]), types.I8)
		bs = append(bs, b)
	}
	bs = append(bs, constant.NewInt(0, types.I8))
	c := constant.NewArray(bs...)
	c.CharArray = true
	return c
}

// CreateEntryBlockAlloca - Create an alloca instruction in the entry block of
// the function.  This is used for mutable variables etc.
func createBlockAlloca(f *ir.Function, elemType types.Type, name string) *ir.InstAlloca {
	// Create a new allocation in the root of the function
	alloca := f.Blocks[0].NewAlloca(elemType)
	// Set the name of the allocation (the variable name)
	// alloca.SetName(name)
	return alloca
}

// Allow functions to return an error isntead of having to manage closing the program each time.
func codegenError(str string, args ...interface{}) value.Value {
	fmt.Fprintf(os.Stderr, "Error: %s\n", fmt.Sprintf(str, args...))
	return nil
}
