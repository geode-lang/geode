package ast

import (
	"io/ioutil"
	"strings"

	"github.com/llir/llvm/ir"
)

// Compiler contains all information to compile a program into a .o file.
type Compiler struct {
	Name string
	// Target             llvm.Target
	// TargetMachine      llvm.TargetMachine
	// TargetData         llvm.TargetData
	Scope              *Scope
	Module             *ir.Module
	blocks             []*ir.BasicBlock
	FN                 *ir.Function // current funciton being compiled
	objectFilesEmitted []string
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

// GetLLVMIR returns the llvm repr of the module
func (c *Compiler) GetLLVMIR() string {
	return c.Module.String()
}

// EmitModuleObject takes an llvm module and emits the object code
func (c *Compiler) EmitModuleObject() string {
	filename := strings.Replace(c.Name, ".g", "", -1) + ".ll"

	llvmir := c.GetLLVMIR()

	writeErr := ioutil.WriteFile(filename, []byte(llvmir), 0666)
	if writeErr != nil {
		panic(writeErr)
	}

	c.objectFilesEmitted = append(c.objectFilesEmitted, filename)

	return filename
}

func (c *Compiler) runInBlock(blk *ir.BasicBlock, fn func()) {
	c.PushBlock(blk)
	fn()
	c.PopBlock()
}

// NewCompiler returns a pointer to a new Compiler object.
func NewCompiler(moduleName string) *Compiler {
	comp := &Compiler{}
	// Initialize the module for this compiler.
	comp.Module = ir.NewModule()
	comp.Name = moduleName

	comp.Scope = NewScope()
	comp.blocks = make([]*ir.BasicBlock, 0)
	return comp
}
