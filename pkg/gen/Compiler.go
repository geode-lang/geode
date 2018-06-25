package gen

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"gitlab.com/nickwanninger/geode/pkg/util/log"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

// Compiler contains all information to compile a program into a .o file.
type Compiler struct {
	OutputName string
	// Target             llvm.Target
	// TargetMachine      llvm.TargetMachine
	// TargetData         llvm.TargetData
	RootScope          *Scope
	RootModule         *ir.Module
	blocks             []*ir.BasicBlock
	FN                 *ir.Function // current funciton being compiled
	objectFilesEmitted []string
	Functions          map[string]*ir.Function
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
	return c.RootModule.String()
}

// EmitModuleObject takes an llvm module and emits the object code
func (c *Compiler) EmitModuleObject() string {
	filename := "out.ll"

	llvmir := c.GetLLVMIR()

	writeErr := ioutil.WriteFile(filename, []byte(llvmir), 0666)
	if writeErr != nil {
		panic(writeErr)
	}

	c.objectFilesEmitted = append(c.objectFilesEmitted, filename)

	return filename
}

// CompileTarget is a target to build a binary for
type CompileTarget int

// Supported compile targets to use
const (
	ASMTarget CompileTarget = iota
	BinaryTarget
)

// Compile the llvm files a Compiler instance has emitted
func (c *Compiler) Compile(target CompileTarget) bool {
	linker := "clang"
	linkArgs := make([]string, 0)

	linkArgs = append(linkArgs, "-O3")

	filename := c.OutputName

	if target == ASMTarget {
		linkArgs = append(linkArgs, "-S", "-masm=intel")
		ext := path.Ext(filename)
		filename = filename[0:len(filename)-len(ext)] + ".s"
	}

	// set the output filename
	linkArgs = append(linkArgs, "-o", filename)

	// Append input files to the end of the command
	for _, objFile := range c.objectFilesEmitted {
		linkArgs = append(linkArgs, objFile)
	}

	cmd := exec.Command(linker, linkArgs...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal("failed to compile with clang: `%s`\n\n%s", err.Error(), string(out))
	}

	// log.Printf(string(out))
	// Clean up all the object files that wwere built in the process.
	c.cleanUpObjectFiles()

	return true
}

func (c *Compiler) cleanUpObjectFiles() {
	for _, objFile := range c.objectFilesEmitted {
		os.Remove(objFile)
	}
}

func (c *Compiler) injectTypes() {
	// Does nothing for now
}

// NewCompiler returns a pointer to a new Compiler object.
func NewCompiler(moduleName string, outputName string) *Compiler {
	comp := &Compiler{}
	comp.OutputName = outputName
	// Initialize the module for this compiler.
	comp.RootModule = ir.NewModule()
	comp.RootScope = NewScope()
	comp.blocks = make([]*ir.BasicBlock, 0)
	comp.Functions = make(map[string]*ir.Function)
	i8 := types.I8
	i8ptr := types.NewPointer(i8)

	printf := comp.RootModule.NewFunction("printf", types.I64, ir.NewParam("format", i8ptr))
	printf.Sig.Variadic = true
	comp.Functions["printf"] = printf

	getchar := comp.RootModule.NewFunction("getchar", types.I8)
	comp.Functions["getchar"] = getchar

	comp.injectTypes()

	return comp
}
