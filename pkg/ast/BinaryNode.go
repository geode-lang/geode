package ast

import (
	"bytes"
	"fmt"

	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

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

// CreateBinaryOp produces a geode binary op (just a wrapper around geode-lang/llvm's binary instructions)
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

// InferType implements Node.InferType
func (n BinaryNode) InferType(scope *Scope) string { return n.Left.InferType(scope) }

// GenAccess implements Accessable.GenAccess
func (n BinaryNode) GenAccess(prog *Program) value.Value {
	return n.Codegen(prog)
}

// Codegen implements Node.Codegen for BinaryNode
func (n BinaryNode) Codegen(prog *Program) value.Value {

	// Generate the left and right nodes
	l := n.Left.Codegen(prog)
	r := n.Right.Codegen(prog)

	// Attempt to cast them with casting precidence
	// This means the operation `int + float` will cast the int to a float.
	l, r, t := binaryCast(prog, l, r)

	if l == nil || r == nil {
		n.SyntaxError()
		log.Fatal("An operand to a binary operation `%s` was nil and failed to generate\n", n.OP)
	}

	blk := prog.Compiler.CurrentBlock()

	if op, valid := binaryOperatorTypeMap[n.OP]; valid {
		return CreateBinaryOp(op.I, op.F, blk, t, l, r)
	}

	if op, valid := booleanComparisonOperatorMap[n.OP]; valid {
		return createCmp(blk, op.I, op.F, t, l, r)
	}

	log.Fatal("Invalid Binary Operator")
	return nil
}
