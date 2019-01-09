package ast

import (
	"sync"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

// Compiler contains all information to
// compile a program from nodes to llvm ir
type Compiler struct {
	Name string
	// A reference to the scope in the package for easier access
	Package *Package
	Module  *ir.Module
	blocks  []*ir.Block
	// FN            *ir.Func // current funciton being compiled
	typeStack     []types.Type
	typestacklock sync.RWMutex

	fnStack     []*ir.Func
	fnstacklock sync.RWMutex
}

// CurrentBlock -
func (c *Compiler) CurrentBlock() *ir.Block {
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
	n.fnStack = c.fnStack
	n.typeStack = c.typeStack
	return n
}

// PushBlock -
func (c *Compiler) PushBlock(blk *ir.Block) {

	c.blocks = append(c.blocks, blk)
}

// PopBlock -
func (c *Compiler) PopBlock() *ir.Block {
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
func (c *Compiler) FunctionDefined(fn *ir.Func) bool {
	for _, defined := range c.Module.Funcs {
		if defined == fn {
			return true
		}
	}
	return false
}

// NewComment appends a comment string to the current block
func (c *Compiler) NewComment(comment string) {
	curBlock := c.CurrentBlock()
	curBlock.Insts = append(curBlock.Insts, NewLLVMComment(comment))

}

func (c *Compiler) genInBlock(blk *ir.Block, fn func() error) error {
	c.PushBlock(blk)
	err := fn()
	c.PopBlock()
	return err
}

// PushType appends a type to the compiler's type stack
func (c *Compiler) PushType(t types.Type) {
	c.typestacklock.Lock()
	c.typeStack = append(c.typeStack, t)
	c.typestacklock.Unlock()
}

// PopType removes an Item from the top of the stack
func (c *Compiler) PopType() (item types.Type) {
	c.typestacklock.Lock()
	item = c.typeStack[len(c.typeStack)-1]
	c.typeStack = c.typeStack[0 : len(c.typeStack)-1]
	c.typestacklock.Unlock()
	return item
}

// EmptyTypeStack does exactly what it seems
func (c *Compiler) EmptyTypeStack() {
	c.typeStack = make([]types.Type, 0)
}

// NewCompiler returns a pointer to a new Compiler object.
func NewCompiler(prog *Program) *Compiler {
	comp := &Compiler{}

	// Initialize the module for this compiler.
	comp.Module = prog.Module

	comp.blocks = make([]*ir.Block, 0)
	comp.typeStack = make([]types.Type, 0)
	return comp
}

// PushFunc appends a Func to the compiler's Func stack
func (c *Compiler) PushFunc(fn *ir.Func) {
	// fmt.Println("pushing", fn.Name)
	c.fnstacklock.Lock()
	c.fnStack = append(c.fnStack, fn)
	c.fnstacklock.Unlock()
}

// PopFunc removes an Item from the top of the stack of functions
func (c *Compiler) PopFunc() (fn *ir.Func) {
	c.fnstacklock.Lock()
	if len(c.fnStack) >= 1 {
		fn = c.fnStack[len(c.fnStack)-1]
		c.fnStack = c.fnStack[0 : len(c.fnStack)-1]
	}
	c.fnstacklock.Unlock()
	// fmt.Println("popping", fn.Name)
	return fn
}

// CurrentFunc returns the top of the function stack
func (c *Compiler) CurrentFunc() *ir.Func {
	return c.fnStack[len(c.fnStack)-1]
}
