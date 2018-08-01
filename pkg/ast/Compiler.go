package ast

import (
	"github.com/geode-lang/llvm/ir"
)

// Compiler contains all information to
// compile a program from nodes to llvm ir
type Compiler struct {
	Name string
	// A reference to the scope in the package for easier access
	Scope   *Scope
	Package *Package
	Module  *ir.Module
	blocks  []*ir.BasicBlock
	FN      *ir.Function // current funciton being compiled
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
func NewCompiler(module *ir.Module, moduleName string, pkg *Package) *Compiler {
	comp := &Compiler{}
	comp.Package = pkg
	// Initialize the module for this compiler.
	comp.Module = module
	comp.Name = moduleName

	comp.Scope = pkg.Scope
	comp.blocks = make([]*ir.BasicBlock, 0)
	return comp
}
