package ast

import (
	"bytes"
	"fmt"

	"github.com/geode-lang/geode/llvm/ir"
	"github.com/geode-lang/geode/llvm/ir/constant"
	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
)

//
// ForNode is a for loop structure representation
type ForNode struct {
	NodeType
	TokenReference

	Index int
	Init  Node
	Cond  Node
	Step  Node
	Body  Node
}

func (n ForNode) String() string {
	buff := &bytes.Buffer{}
	fmt.Fprintf(buff, "for %s; %s; %s %s", n.Init, n.Cond, n.Step, n.Body)
	return buff.String()
}

// NameString implements Node.NameString
func (n ForNode) NameString() string { return "ForNode" }

// Codegen implements Node.Codegen for ForNode
func (n ForNode) Codegen(prog *Program) (value.Value, error) {

	// The name of the blocks is prefixed so we can determine which for loop a block is for.
	namePrefix := fmt.Sprintf("F%X_", n.Index)
	parentBlock := prog.Compiler.CurrentBlock()

	prog.ScopeDown(n.Token)
	var err error
	var predicate value.Value
	var condBlk *ir.BasicBlock
	var bodyBlk *ir.BasicBlock
	var bodyGenBlk *ir.BasicBlock
	var endBlk *ir.BasicBlock
	parentFunc := parentBlock.Parent

	condBlk = parentFunc.NewBlock(namePrefix + "cond")

	n.Init.Codegen(prog)

	parentBlock.NewBr(condBlk)

	err = prog.Compiler.genInBlock(condBlk, func() error {
		predicate, _ = n.Cond.Codegen(prog)
		one := constant.NewInt(1, types.I1)

		c, err := createTypeCast(prog, predicate, types.I1)
		if err != nil {
			return err
		}
		predicate = condBlk.NewICmp(ir.IntEQ, one, c)
		return nil
	})

	if err != nil {
		return nil, err
	}

	bodyBlk = parentFunc.NewBlock(namePrefix + "body")

	stepBlk := parentFunc.NewBlock(namePrefix + "step")

	err = prog.Compiler.genInBlock(bodyBlk, func() error {

		scp := prog.Scope
		gen, err := n.Body.Codegen(prog)
		if err != nil {
			return err
		}
		prog.Scope = scp

		// prog.Compiler.CurrentBlock().AppendInst(NewLLVMComment("branch to the step"))
		bodyGenBlk = gen.(*ir.BasicBlock)

		if err != nil {
			return err
		}

		bodyGenBlk.BranchIfNoTerminator(stepBlk)
		bodyBlk.BranchIfNoTerminator(stepBlk)

		return nil
	})

	err = prog.Compiler.genInBlock(stepBlk, func() error {
		scp := prog.Scope
		_, err := n.Step.Codegen(prog)
		prog.Scope = scp
		return err
	})

	stepBlk.BranchIfNoTerminator(condBlk)
	endBlk = parentFunc.NewBlock(namePrefix + "end")
	prog.Compiler.PushBlock(endBlk)
	condBlk.NewCondBr(predicate, bodyBlk, endBlk)

	if err := prog.ScopeUp(); err != nil {
		return nil, err
	}
	return endBlk, nil
}
