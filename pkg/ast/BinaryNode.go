package ast

import (
	"bytes"
	"fmt"

	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

func createCmp(blk *ir.BasicBlock, i ir.IntPred, f ir.FloatPred, t types.Type, left, right value.Value) value.Value {

	var val value.Value

	if types.IsInt(t) {
		val = blk.NewICmp(i, left, right)
	}
	if types.IsFloat(t) {
		val = blk.NewFCmp(f, left, right)
	}

	return val
}

// CreateBinaryOp produces a geode binary op (just a wrapper around geode-lang/llvm's binary instructions)
func CreateBinaryOp(intstr, fltstr string, blk *ir.BasicBlock, t types.Type, left, right value.Value) value.Value {

	var val *GeodeBinaryInstr
	if types.IsInt(t) {
		val = NewGeodeBinaryInstr(intstr, left, right)
	} else {
		val = NewGeodeBinaryInstr(fltstr, left, right)
	}

	blk.AppendInst(val)

	return val
}

type numericalBinaryOperator struct {
	I string
	F string
}

type comparisonOperation struct {
	I ir.IntPred
	F ir.FloatPred
}

var binaryOperatorTypeMap = map[string]numericalBinaryOperator{
	"+":  {"add", "fadd"},
	"-":  {"sub", "fsub"},
	"*":  {"mul", "fmul"},
	"/":  {"sdiv", "fdiv"},
	"%":  {"srem", "frem"},
	">>": {"lshr", "lshr"},
	"<<": {"shl", "shl"},
	"||": {"or", "or"},
	"&&": {"and", "and"},
	"^":  {"xor", "xor"},
}

var booleanComparisonOperatorMap = map[string]comparisonOperation{
	"=":  {ir.IntEQ, ir.FloatOEQ},
	"!=": {ir.IntNE, ir.FloatONE},
	">":  {ir.IntSGT, ir.FloatOGT},
	">=": {ir.IntSGE, ir.FloatOGE},
	"<":  {ir.IntSLT, ir.FloatOLT},
	"<=": {ir.IntSLE, ir.FloatOLE},
}

// BinaryNode is a representation of a binary operation
type BinaryNode struct {
	NodeType
	TokenReference
	Accessable

	OP    string
	Left  Node
	Right Node
}

func (n BinaryNode) String() string {
	buff := &bytes.Buffer{}
	fmt.Fprintf(buff, "%s %s %s", n.Left, n.OP, n.Right)
	return buff.String()
}

// NameString implements Node.NameString
func (n BinaryNode) NameString() string { return "BinaryNode" }

// GenAccess implements Accessable.GenAccess
func (n BinaryNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}

// Codegen implements Node.Codegen for BinaryNode
func (n BinaryNode) Codegen(prog *Program) (value.Value, error) {

	// Generate the left and right nodes
	l, err := n.Left.Codegen(prog)
	if err != nil {
		return nil, err
	}
	r, err := n.Right.Codegen(prog)
	if err != nil {
		return nil, err
	}

	// Attempt to cast them with casting precidence
	// This means the operation `int + float` will cast the int to a float.
	l, r, t, resultcast := binaryCast(prog, l, r)

	if l == nil || r == nil {
		n.SyntaxError()
		return nil, fmt.Errorf("an operand to a binary operation `%s` was nil and failed to generate", n.OP)
	}

	blk := prog.Compiler.CurrentBlock()

	var value value.Value

	if op, valid := binaryOperatorTypeMap[n.OP]; valid {
		value = CreateBinaryOp(op.I, op.F, blk, t, l, r)
	}

	if op, valid := booleanComparisonOperatorMap[n.OP]; valid {
		value = createCmp(blk, op.I, op.F, t, l, r)
	}

	if value == nil {
		return nil, fmt.Errorf("invalid binary operator")
	}

	if resultcast != nil {
		value, _ = createTypeCast(prog, value, resultcast)
	}

	return value, nil

}

func binaryCast(prog *Program, left, right value.Value) (value.Value, value.Value, types.Type, types.Type) {

	var resultcast types.Type
	if types.IsPointer(left.Type()) {
		left = prog.Compiler.CurrentBlock().NewPtrToInt(left, types.I64)
		resultcast = left.Type()
	}

	if types.IsPointer(right.Type()) {
		right = prog.Compiler.CurrentBlock().NewPtrToInt(right, types.I64)
		resultcast = right.Type()
	}

	// Right and Left types
	lt := left.Type()
	rt := right.Type()

	var casted types.Type

	// Get the cast precidence of both sides
	leftPrec := prog.CastPrecidence(lt)
	rightPrec := prog.CastPrecidence(rt)

	if leftPrec > rightPrec {
		casted = lt
		right, _ = createTypeCast(prog, right, lt)
	} else {
		casted = rt
		left, _ = createTypeCast(prog, left, rt)
	}
	return left, right, casted, resultcast
}
