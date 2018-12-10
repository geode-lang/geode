package ast

import (
	"fmt"
	"os"

	"github.com/geode-lang/geode/pkg/arg"
	"github.com/geode-lang/geode/pkg/gtypes"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/metadata"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// A global number to indicate which `name index` we are on. This way,
// the mangler will never output the same name twice as this number is monotonic
var nameNumber int

func mangleName(name string) string {
	nameNumber++
	return fmt.Sprintf("%s_%d", name, nameNumber)
}

// func branchIfNoTerminator(blk *ir.BasicBlock, to *ir.BasicBlock) {
// 	if blk.Term == nil {
// 		blk.NewBr(to)
// 	}
// }

// Codegen returns some NamespaceNode's arguments
func (n NamespaceNode) Codegen(prog *Program) (value.Value, error) { return nil, nil }

// Handle will do ast-level handling for a dependency node
func (n DependencyNode) Handle(prog *Program) value.Value {

	return nil
}

// Codegen implements Node.Codegen for DependencyNode
func (n DependencyNode) Codegen(prog *Program) (value.Value, error) { return nil, nil }

// Codegen implements Node.Codegen for IfNode
func (n IfNode) Codegen(prog *Program) (value.Value, error) {

	predicate, err := n.If.Codegen(prog)
	if err != nil {
		return nil, err
	}
	zero := constant.NewInt(types.I32, 0)
	// The name of the blocks is prefixed because
	namePrefix := fmt.Sprintf("if.%d.", n.Index)
	parentBlock := prog.Compiler.CurrentBlock()
	c, err := createTypeCast(prog, predicate, types.I32)
	if err != nil {
		return nil, err
	}
	predicate = parentBlock.NewICmp(enum.IPredNE, zero, c)
	parentFunc := parentBlock.Parent

	var thenGenBlk *ir.BasicBlock
	var endBlk *ir.BasicBlock

	thenBlk := parentFunc.NewBlock(mangleName(namePrefix + "then"))

	prog.Compiler.genInBlock(thenBlk, func() error {
		gen, gerr := n.Then.Codegen(prog)
		if gerr != nil {
			return gerr
		}
		thenGenBlk = gen.(*ir.BasicBlock)
		return nil
	})

	elseBlk := parentFunc.NewBlock(mangleName(namePrefix + "else"))
	var elseGenBlk *ir.BasicBlock

	prog.Compiler.genInBlock(elseBlk, func() error {
		// We only want to construct the else block if there is one.
		if n.Else != nil {
			gen, gerr := n.Else.Codegen(prog)
			if gerr != nil {
				return gerr
			}
			elseGenBlk, _ = gen.(*ir.BasicBlock)
		}
		return nil
	})

	endBlk = parentFunc.NewBlock(mangleName(namePrefix + "end"))
	prog.Compiler.PushBlock(endBlk)
	// We need to make sure these blocks have terminators.
	// in order to do that, we branch to the end block

	BranchIfNoTerminator(thenBlk, endBlk)
	BranchIfNoTerminator(thenGenBlk, endBlk)
	BranchIfNoTerminator(elseBlk, endBlk)

	if elseGenBlk != nil {
		BranchIfNoTerminator(elseGenBlk, endBlk)
	}

	parentBlock.NewCondBr(predicate, thenBlk, elseBlk)

	return endBlk, nil
}

// Codegen implements Node.Codegen for CharNode
func (n CharNode) Codegen(prog *Program) (value.Value, error) {
	return constant.NewInt(types.I8, int64(n.Value)), nil
}

// GenAccess returns the value from a given CharNode
func (n CharNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}

// Codegen implements Node.Codegen for UnaryNode
func (n UnaryNode) Codegen(prog *Program) (value.Value, error) {

	// handle reference operation
	if n.Operator == "&" {

		node, ok := n.Operand.(Reference)
		if !ok {
			n.SyntaxError()
			return nil, fmt.Errorf("'&' operator called on non-addressable operand")
		}

		return node.Alloca(prog), nil
	}

	operandValue, err := n.Operand.Codegen(prog)
	if err != nil {
		return nil, err
	}
	if operandValue == nil {
		n.Operand.SyntaxError()
		return nil, fmt.Errorf("nil operand")
	}

	if n.Operator == "-" {

		if types.IsFloat(operandValue.Type()) {
			return prog.Compiler.CurrentBlock().NewFSub(constant.NewFloat(types.Double, 0), operandValue), nil
		} else if types.IsInt(operandValue.Type()) {
			return prog.Compiler.CurrentBlock().NewSub(constant.NewInt(types.I64, 0), operandValue), nil
		}
		return nil, fmt.Errorf("Unable to make a non integer/float into a negative")

	}

	// the not operation is interesting as there is no intrinsic llvm "not" instruction
	// so what I must do is check if the value is 0 with an icmp. Then I xor that with true
	// and sign extend it to an i32 - i32 is a safe value, idk
	if n.Operator == "!" {
		if !types.IsInt(operandValue.Type()) {
			return nil, fmt.Errorf("unable to '!' (not) type %q", operandValue.Type())
		}

		opVal, _ := createTypeCast(prog, operandValue, types.I1)

		eq := prog.Compiler.CurrentBlock().NewICmp(enum.IPredNE, opVal, constant.False)
		inv := prog.Compiler.CurrentBlock().NewXor(eq, constant.True)
		ext := prog.Compiler.CurrentBlock().NewZExt(inv, types.I32)

		return ext, nil

	}

	// handle dereference operation
	if n.Operator == "*" {

		// fmt.Println(prog.Compiler.CurrentFunc())
		if types.IsPointer(operandValue.Type()) {
			return prog.Compiler.CurrentBlock().NewLoad(operandValue), nil
		}
		n.SyntaxError()
		return nil, fmt.Errorf("attempt to dereference a non-pointer variable")
	}

	return operandValue, nil
}

// GenAccess implements Accessable.GenAccess
func (n UnaryNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}

// Codegen implements Node.Codegen for WhileNode
func (n WhileNode) Codegen(prog *Program) (value.Value, error) {

	// The name of the blocks is prefixed because
	namePrefix := fmt.Sprintf("while_%d_", n.Index)
	parentBlock := prog.Compiler.CurrentBlock()

	parentFunc := parentBlock.Parent
	startblock := parentFunc.NewBlock(mangleName(namePrefix + "start"))
	prog.Compiler.PushBlock(startblock)
	predicate, err := n.If.Codegen(prog)
	if err != nil {
		return nil, err
	}
	one := constant.NewInt(types.I1, 1)
	prog.Compiler.PopBlock()
	BranchIfNoTerminator(parentBlock, startblock)
	c, err := createTypeCast(prog, predicate, types.I1)
	if err != nil {
		return nil, err
	}
	predicate = startblock.NewICmp(enum.IPredEQ, one, c)

	var endBlk *ir.BasicBlock

	bodyBlk := parentFunc.NewBlock(mangleName(namePrefix + "body"))
	prog.Compiler.PushBlock(bodyBlk)

	v, err := n.Body.Codegen(prog)
	if err != nil {
		return nil, err
	}
	bodyGenBlk := v.(*ir.BasicBlock)

	// If there is no terminator for the block, IE: no return
	// branch to the merge block

	endBlk = parentFunc.NewBlock(mangleName(namePrefix + "merge"))
	prog.Compiler.PushBlock(endBlk)

	BranchIfNoTerminator(bodyBlk, startblock)
	BranchIfNoTerminator(bodyGenBlk, startblock)

	startblock.NewCondBr(predicate, bodyBlk, endBlk)

	// branchIfNoTerminator(c.CurrentBlock(), endBlk)

	return endBlk, nil
}

func typeSize(t types.Type) int {
	switch t := t.(type) {
	case *types.IntType:
		return int(t.BitSize)
	case *types.FloatType:
		return gtypes.FloatBitSize(t)
	}

	return -1
}

func typesAreLooselyEqual(a, b types.Type) bool {
	return gtypes.IsNumber(a) && gtypes.IsNumber(b)
}

// createTypeCast is where most, if not all, type casting happens in the language.
func createTypeCast(prog *Program, in value.Value, to types.Type) (value.Value, error) {

	inType := in.Type()
	fromInt := types.IsInt(inType)
	fromFloat := types.IsFloat(inType)

	toInt := types.IsInt(to)
	toFloat := types.IsFloat(to)

	inSize := typeSize(inType)
	outSize := typeSize(to)

	// If the cast would not change the type, just return the in value, nil
	if types.Equal(inType, to) {
		return in, nil
	}

	if c, ok := in.(*constant.Int); ok && types.IsInt(to) {
		c.Typ = to.(*types.IntType)
		return c, nil
	}

	if c, ok := in.(*constant.Float); ok && types.IsFloat(to) {
		c.Typ = to.(*types.FloatType)
		return c, nil
	}

	if types.Equal(to, types.Void) {
		return nil, nil
	}

	if types.IsPointer(inType) && types.IsPointer(to) {
		return prog.Compiler.CurrentBlock().NewBitCast(in, to), nil
	}

	if fromFloat && toInt {
		return prog.Compiler.CurrentBlock().NewFPToSI(in, to), nil
	}

	if fromInt && toFloat {
		return prog.Compiler.CurrentBlock().NewSIToFP(in, to), nil
	}

	if fromInt && toInt {
		if inSize < outSize {
			return prog.Compiler.CurrentBlock().NewSExt(in, to), nil
		}
		if inSize == outSize {
			return in, nil
		}
		return prog.Compiler.CurrentBlock().NewTrunc(in, to), nil
	}

	if fromFloat && toFloat {
		if inSize < outSize {
			return prog.Compiler.CurrentBlock().NewFPExt(in, to), nil
		}
		if inSize == outSize {
			return in, nil
		}
		return prog.Compiler.CurrentBlock().NewFPTrunc(in, to), nil
	}

	// If the cast would not change the type, just return the in value, nil
	if types.Equal(inType, to) {
		return in, nil
	}

	if types.IsPointer(inType) && types.IsInt(to) {
		return prog.Compiler.CurrentBlock().NewPtrToInt(in, to), nil
	}

	if types.IsInt(inType) && types.IsPointer(to) {
		return prog.Compiler.CurrentBlock().NewIntToPtr(in, to), nil
	}

	return nil, fmt.Errorf("Failed to typecast type %s to %s", inType.String(), to)
}

// Codegen implements Node.Codegen for ReturnNode
func (n ReturnNode) Codegen(prog *Program) (value.Value, error) {

	var retVal value.Value
	var err error

	if !prog.Compiler.CurrentFunc().Sig.RetType.Equal(types.Void) {
		if n.Value != nil {
			retVal, err = n.Value.Codegen(prog)
			if err != nil {

				return nil, err
			}
			given := retVal.Type()
			expected := prog.Compiler.CurrentFunc().Sig.RetType
			if !types.Equal(given, expected) {
				if !(types.IsInt(given) && types.IsInt(expected)) {
					n.SyntaxError()
					fnName, err := UnmangleFunctionName(prog.Compiler.CurrentFunc().Name())
					if err != nil {

						return nil, err
					}
					expectedName, err := prog.Scope.FindTypeName(expected)
					if err != nil {

						return nil, err
					}
					givenName, err := prog.Scope.FindTypeName(given)
					if err != nil {

						return nil, err
					}

					return nil, fmt.Errorf("incorrect return value for function %s. expected: %s (%s). given: %s (%s)", fnName, expectedName, expected, givenName, given)
				}
				retVal, err = createTypeCast(prog, retVal, prog.Compiler.CurrentFunc().Sig.RetType)
				if err != nil {

					return nil, err
				}
			}
		} else {

			retVal = nil
		}
	}

	ret := prog.Compiler.CurrentBlock().NewRet(retVal)

	if *arg.EnableDebug {
		mdNode := n.Token.DILocation(prog.Scope.DebugInfo)
		md := &metadata.Attachment{
			Name: "dbg",
			Node: mdNode,
		}
		ret.Metadata = append(ret.Metadata, md)
	}

	return retVal, nil
}

// newCharArray returns a character array constant based on the given string.
func newCharArray(s string) *constant.CharArray {
	return constant.NewCharArray([]byte(s))
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

// BranchIfNoTerminator checks if the block has a terminator, and if it doesn't,
// set it to a branch instead.
func BranchIfNoTerminator(block, target *ir.BasicBlock) {
	if block.Term == nil {
		block.NewBr(target)
	}
}

// gep returns a new getelementptr instruction based on the given source address
// and element indices. It handles Geode specific types of which
// ir.NewGetElementPtr is unaware.
func gep(src value.Value, indices ...value.Value) *ir.InstGetElementPtr {
	// Locate element type of src.
	var (
		// Pointer to element type of src.
		elemPtr *types.Type
		// Geode element type.
		elem gtypes.Type
	)
	switch typ := src.Type().(type) {
	case *types.PointerType:
		if e, ok := typ.ElemType.(gtypes.Type); ok {
			elemPtr = &typ.ElemType
			elem = e
		}
	case *types.VectorType:
		t, ok := typ.ElemType.(*types.PointerType)
		if !ok {
			panic(fmt.Errorf("invalid vector element type; expected *types.Pointer, got %T", typ.ElemType))
		}
		if e, ok := t.ElemType.(gtypes.Type); ok {
			elemPtr = &typ.ElemType
			elem = e
		}
	default:
		panic(fmt.Errorf("support for souce type %T not yet implemented", typ))
	}
	// Set element type of src to underlying LLVM IR type of Geode type, prior to
	// calling ir.NewGetElementPtr since it is not aware of Geode types.
	if elemPtr != nil {
		*elemPtr = elem.Underlying()
	}
	inst := ir.NewGetElementPtr(src, indices...)
	// Restore original Geode type.
	if elemPtr != nil {
		*elemPtr = elem
	}
	return inst
}
