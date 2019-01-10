package ast

import (
	"bytes"
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func createCmp(blk *ir.Block, i enum.IPred, f enum.FPred, t types.Type, left, right value.Value) value.Value {

	var val value.Value

	if types.IsInt(t) {
		val = blk.NewICmp(i, left, right)
	}
	if types.IsFloat(t) {
		val = blk.NewFCmp(f, left, right)
	}

	return val
}

// CreateBinaryOp produces a geode binary op (just a wrapper around geode-lang/geode/llvm's binary instructions)
func CreateBinaryOp(intstr, fltstr string, blk *ir.Block, t types.Type, left, right value.Value) value.Value {

	var val BinaryInstruction
	if types.IsInt(t) {
		val = NewGeodeBinaryInstr(intstr, left, right)
	} else {
		val = NewGeodeBinaryInstr(fltstr, left, right)
	}

	blk.Insts = append(blk.Insts, val)

	return val
}

type numericalBinaryOperator struct {
	I string
	F string
}

type comparisonOperation struct {
	I enum.IPred
	F enum.FPred
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
	"==": {enum.IPredEQ, enum.FPredOEQ},
	"!=": {enum.IPredNE, enum.FPredONE},
	">":  {enum.IPredSGT, enum.FPredOGT},
	">=": {enum.IPredSGE, enum.FPredOGE},
	"<":  {enum.IPredSLT, enum.FPredOLT},
	"<=": {enum.IPredSLE, enum.FPredOLE},
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

// CodegenCompoundOperator generates a compound operator expression
func CodegenCompoundOperator(prog *Program, left, right Node, compop string) (value.Value, error) {
	var op string
	var ok bool

	switch compop {
	case "+=":
		op = "+"
	case "-=":
		op = "-"
	case "*=":
		op = "*"
	case "/=":
		op = "/"
	default:
		return nil, fmt.Errorf("unknown compound assignment %q", compop)

	}

	n := AssignmentNode{}
	n.Assignee, ok = left.(Assignable)
	if !ok {
		return nil, fmt.Errorf("left hand side of compound assignment %q is not assignable", compop)
	}

	binary := BinaryNode{}

	binary.Left = left
	binary.Right = right
	binary.OP = op

	n.Value = binary

	return n.Codegen(prog)

}

// Codegen implements Node.Codegen for BinaryNode
func (n BinaryNode) Codegen(prog *Program) (value.Value, error) {

	if n.OP == "=" {
		lhs, ok := n.Left.(Assignable)
		if !ok {
			return nil, fmt.Errorf("attempt to assign to a non assignable value '%s' (%T)", n.Left, n.Left)
		}

		rhs, ok := n.Right.(Accessable)
		if !ok {
			return nil, fmt.Errorf("attempt to assign with a non accessable value '%s'", n.Right)
		}

		a := AssignmentNode{}
		a.Assignee = lhs
		a.Value = rhs
		a.NodeType = nodeAssignment
		return a.Codegen(prog)
	}

	switch n.OP {
	case "+=", "-=", "*=", "/=":
		return CodegenCompoundOperator(prog, n.Left, n.Right, n.OP)
	case "+", "-":
		add := AddSubNode{}
		add.Left = n.Left
		add.Right = n.Right
		add.Sub = n.OP == "-"
		add.TokenReference = n.TokenReference
		add.NodeType = nodeBinary
		return add.Codegen(prog)
	}

	if n.Left == nil || n.Right == nil {
		n.SyntaxError()
		return nil, fmt.Errorf("invalid binary expression")
	}
	// Generate the left and right nodes
	l, err := n.Left.Codegen(prog)
	if err != nil {
		return nil, err
	}
	r, err := n.Right.Codegen(prog)
	if err != nil {
		return nil, err
	}

	mustCastToPtr := false
	var finalPointerType types.Type

	if types.IsPointer(l.Type()) {
		mustCastToPtr = true
		finalPointerType = l.Type()
	}

	if types.IsPointer(r.Type()) {
		mustCastToPtr = true
		finalPointerType = r.Type()
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
		return nil, fmt.Errorf("invalid binary operator %s", n.OP)
	}

	if resultcast != nil {
		value, _ = createTypeCast(prog, value, resultcast)
	}

	if mustCastToPtr {
		value, _ = createTypeCast(prog, value, finalPointerType)
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

// ----------------------- Add Node -----------------------

// AddSubNode represents an addition or subtraction operation
type AddSubNode struct {
	NodeType
	TokenReference
	Accessable

	Sub bool

	Left  Node
	Right Node
}

func (n AddSubNode) String() string {
	buff := &bytes.Buffer{}
	op := "+"
	if n.Sub {
		op = "-"
	}
	fmt.Fprintf(buff, "%s %s %s", n.Left, op, n.Right)
	return buff.String()
}

// NameString implements Node.NameString
func (n AddSubNode) NameString() string { return "AddSubNode" }

// GenAccess implements Accessable.GenAccess
func (n AddSubNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}

// Codegen implements Node.Codegen for AddSubNode
func (n AddSubNode) Codegen(prog *Program) (value.Value, error) {

	opname := "add"
	if n.Sub {
		opname = "sub"
	}

	var result value.Value
	var err error
	// Generate the left and right nodes
	l, err := n.Left.Codegen(prog)
	if err != nil {
		return nil, err
	}
	r, err := n.Right.Codegen(prog)
	if err != nil {
		return nil, err
	}
	// TODO: handle unsigned numbers... (maybe)
	left, right, t, resultcast := binaryCast(prog, l, r)

	// float add/sub operations on numeric types are prefixed with 'f'
	if types.IsFloat(t) {
		opname = "f" + opname
	}

	inst := NewGeodeBinaryInstr(opname, left, right)
	result = inst
	curBlock := prog.Compiler.CurrentBlock()
	curBlock.Insts = append(curBlock.Insts, inst)

	if resultcast != nil {
		result, err = createTypeCast(prog, result, resultcast)
		if err != nil {
			return nil, err
		}
	}

	// finalType is the resultant type to cast to at the end
	// var finalType types.Type
	return result, nil
}
