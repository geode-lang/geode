package gen

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/go-llvm/llvm"
)

// Compiler contains all information to compile a program into a .o file.
type Compiler struct {
	OutputName         string
	Target             llvm.Target
	TargetMachine      llvm.TargetMachine
	TargetData         llvm.TargetData
	RootScope          *Scope
	RootModule         llvm.Module
	RootPassManager    llvm.PassManager
	Builder            llvm.Builder
	CodeModel          llvm.CodeModel
	RelocMode          llvm.RelocMode
	Opt                llvm.CodeGenOptLevel
	objectFilesEmitted []string
}

// GetLLVMIR returns the llvm repr of the module
func (c *Compiler) GetLLVMIR() string {

	return c.RootModule.String()
}

// EmitModuleObject takes an llvm module and emits the object code
func (c *Compiler) EmitModuleObject() string {
	filename := "out.o"

	c.TargetMachine = c.Target.CreateTargetMachine(llvm.DefaultTargetTriple(), "", "", llvm.CodeGenLevelNone, llvm.RelocDefault, llvm.CodeModelDefault)
	// targetData := targetMachine.TargetData()

	membuf, emitErr := c.TargetMachine.EmitToMemoryBuffer(c.RootModule, llvm.ObjectFile)
	if emitErr != nil {
		panic(emitErr)
	}

	writeErr := ioutil.WriteFile(filename, membuf.Bytes(), 0666)
	if writeErr != nil {
		panic(writeErr)
	}

	c.objectFilesEmitted = append(c.objectFilesEmitted, filename)

	return filename
}

// Compile the object files a Compiler instance has emitted
func (c *Compiler) Compile() bool {
	linker := "cc"
	linkArgs := make([]string, 0)

	linkArgs = append(linkArgs, "-fno-PIE", "-nodefaultlibs", "-lc", "-lm")

	linkArgs = append(linkArgs, "-o", c.OutputName)

	for _, objFile := range c.objectFilesEmitted {
		linkArgs = append(linkArgs, objFile)
	}

	cmd := exec.Command(linker, linkArgs...)

	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("failed to link object files: `%s`\n%s", err.Error(), string(out))
		return false
	}

	// Clean up all the object files that wwere built in the process.
	c.cleanUpObjectFiles()

	return true
}

func (c *Compiler) cleanUpObjectFiles() {
	for _, objFile := range c.objectFilesEmitted {
		os.Remove(objFile)
	}
}

// NewCompiler returns a pointer to a new Compiler object.
func NewCompiler(moduleName string, outputName string) *Compiler {
	comp := &Compiler{}
	comp.OutputName = outputName
	// Initialize the module for this compiler.
	comp.RootModule = llvm.NewModule(moduleName)
	comp.RootPassManager = llvm.NewFunctionPassManagerForModule(comp.RootModule)
	comp.Builder = llvm.NewBuilder()
	comp.RootScope = NewScope()
	CPU := "generic"
	features := ""
	comp.CodeModel = llvm.CodeModelDefault
	comp.RelocMode = llvm.RelocDefault

	targetTripple := llvm.DefaultTargetTriple()
	comp.Opt = llvm.CodeGenLevelAggressive
	comp.Target, _ = llvm.GetTargetFromTriple(targetTripple)
	comp.Target.CreateTargetMachine(targetTripple, CPU, features, comp.Opt, comp.RelocMode, comp.CodeModel)
	return comp
}
