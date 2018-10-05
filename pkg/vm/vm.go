package vm

import (
	"fmt"

	"github.com/geode-lang/geode/llvm/ir"
	"github.com/geode-lang/geode/llvm/ir/value"
)

// VirtualMachine is a structure that can run a *ir.Module
// in the context of the geode programming language
type VirtualMachine struct {
	Module *ir.Module
	Scope  *Scope
}

// New constructs a new VM with the module passed
func New(mod *ir.Module) *VirtualMachine {
	vm := &VirtualMachine{}
	vm.Module = mod
	return vm
}

func (v *VirtualMachine) String() string {
	return "TODO: IMPLEMENT A DECENT STRINGER FOR A VM"
}

// RunFunctionName runs a function in the virtual machine with arguments
func (v *VirtualMachine) RunFunctionName(fnName string, args ...Value) (Value, error) {

	var function *ir.Function

	for _, fn := range v.Module.Funcs {
		if fn.Name == fnName {
			function = fn
		}
	}

	if function == nil {
		return nil, fmt.Errorf("unable to find function %q", fnName)
	}
	return v.RunFunction(function, args...)
}

// Eval a value
func (v *VirtualMachine) Eval(val value.Value) {
	fmt.Printf("%s\n", val)
}

// EvalInst evaluates an instruction
func (v *VirtualMachine) EvalInst(i ir.Instruction) {
	fmt.Printf("%s\n", i)
}

// RunFunction runs a single function in the virtual machine's context
func (v *VirtualMachine) RunFunction(fn *ir.Function, args ...Value) (Value, error) {

	// go into a child scope for this function
	v.Scope = v.Scope.SpawnChild()

	for _, i := range fn.Blocks[0].Insts {
		v.EvalInst(i)
	}

	// Pop the function scope
	v.Scope = v.Scope.Parent
	return nil, nil
}
