package ast

import (
	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/types"
)

// Compiler contains all information to
// compile a program from nodes to llvm ir
type Compiler struct {
	Name string
	// A reference to the scope in the package for easier access
	Package   *Package
	Module    *ir.Module
	blocks    []*ir.BasicBlock
	FN        *ir.Function // current funciton being compiled
	typeCache types.Type
}

// CurrentBlock -
func (c *Compiler) CurrentBlock() *ir.BasicBlock {
	l := len(c.blocks)
	if l == 0 {
		return nil
	}
	blk := (c.blocks)[l-1]
	return blk
}

// Copy a compiler's data into a new compiler
func (c *Compiler) Copy() *Compiler {
	n := &Compiler{}
	n.Package = c.Package
	n.Module = c.Module
	n.blocks = c.blocks
	n.FN = c.FN
	n.typeCache = c.typeCache
	return n
}

// PushBlock -
func (c *Compiler) PushBlock(blk *ir.BasicBlock) {

	c.blocks = append(c.blocks, blk)
}

// PopBlock -
func (c *Compiler) PopBlock() *ir.BasicBlock {
	l := len(c.blocks)
	if l == 0 {
		return nil
	}

	blk := (c.blocks)[l-1]
	c.blocks = (c.blocks)[:l-1]
	return blk
}

// FunctionDefined returns whether or not a function
// with a name has been defined in the module's scope
func (c *Compiler) FunctionDefined(fn *ir.Function) bool {
	for _, defined := range c.Module.Funcs {
		if defined == fn {
			return true
		}
	}
	return false
}

func (c *Compiler) genInBlock(blk *ir.BasicBlock, fn func()) {
	c.PushBlock(blk)
	fn()
	c.PopBlock()
}

// NewCompiler returns a pointer to a new Compiler object.
func NewCompiler(prog *Program) *Compiler {
	comp := &Compiler{}

	// Initialize the module for this compiler.
	comp.Module = prog.Module

	comp.blocks = make([]*ir.BasicBlock, 0)
	return comp
}
