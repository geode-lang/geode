package ast

import (
	"bytes"
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// FunctionCallNode is a function call, example: `foo(a, b, c)`. This would be:
//    Name = "foo"
//    Args = [a, b, c]    <- these are Node references
type FunctionCallNode struct {
	NodeType
	TokenReference

	Name Callable
	Args []Node
}

// NewRuntimeFunctionCall returns a new function call value

// NameString implements Node.NameString
func (n FunctionCallNode) NameString() string { return "FunctionCallNode" }

func (n FunctionCallNode) String() string {
	buff := &bytes.Buffer{}

	fmt.Fprintf(buff, "%s(", n.Name)
	for i, arg := range n.Args {
		fmt.Fprintf(buff, "%s", arg)
		if i < len(n.Args)-1 {
			fmt.Fprintf(buff, ", ")
		}
	}

	fmt.Fprintf(buff, ")")
	return buff.String()
}

// GenAccess implements Accessable.GenAccess
func (n FunctionCallNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}

// Codegen implements Node.Codegen for FunctionCallNode
func (n FunctionCallNode) Codegen(prog *Program) (value.Value, error) {

	// var name string
	var err error

	args := []value.Value{}
	argTypes := []types.Type{}

	for _, arg := range n.Args {

		if ac, isAccessable := arg.(Accessable); isAccessable {
			val, err := ac.GenAccess(prog)
			if err != nil {
				return nil, err
			}

			args = append(args, val)
			argTypes = append(argTypes, val.Type())
			if args[len(args)-1] == nil {
				return nil, fmt.Errorf("argument to function %q failed to generate code", n.Name)
			}
		} else {
			arg.SyntaxError()
			return nil, fmt.Errorf("argument to function call to '%s' is not accessable (has no readable value). Node type %s", n.Name, arg.Kind())
		}
	}

	callee, prependingArgs, err := n.Name.GetFunc(prog, argTypes)
	if err != nil {
		return nil, err
	}
	if prependingArgs != nil {
		args = append(prependingArgs, args...)

		prependingTypes := []types.Type{}
		for _, arg := range prependingArgs {
			prependingTypes = append(prependingTypes, arg.Type())
		}
		argTypes = append(prependingTypes, argTypes...)
	}

	if callee == nil {
		return nil, fmt.Errorf("unknown function %q referenced at %s", n.Name, n.Token.FileInfo())
	}

	// Attempt to typecast all the args into the correct type
	for i, paramType := range callee.Sig.Params {
		args[i], _ = createTypeCast(prog, args[i], paramType)
	}

	// Varargs require type conversion to a standardized type
	// So we will use the same type promotion c uses
	//  if int && type != i32 -> type = i32
	//  if fnn && type != f64 -> type = f64
	arguments := make([]value.Value, 0, len(args))

	for i, arg := range args {

		if callee.Sig.Variadic && i >= len(callee.Params) {
			if types.IsInt(arg.Type()) {
				if !types.Equal(arg.Type(), types.I32) {
					c, err := createTypeCast(prog, arg, types.I32)
					if err != nil {
						return nil, err
					}
					arguments = append(arguments, c)
					continue
				}

			}

			if types.IsFloat(arg.Type()) {
				if !types.Equal(arg.Type(), types.Double) {
					c, err := createTypeCast(prog, arg, types.Double)
					if err != nil {
						return nil, err
					}
					arguments = append(arguments, c)
					continue
				}
			}
		}
		arguments = append(arguments, arg)
	}

	return prog.Compiler.CurrentBlock().NewCall(callee, arguments...), nil
}

// Alloca implements Reference.Alloca
func (n FunctionCallNode) Alloca(prog *Program) value.Value {
	val, err := n.Codegen(prog)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	alloc := prog.Compiler.CurrentBlock().NewAlloca(val.Type())
	prog.Compiler.CurrentBlock().NewStore(val, alloc)
	return alloc
}

// Load implements Reference.Load
func (n FunctionCallNode) Load(blk *ir.BasicBlock, prog *Program) *ir.InstLoad {
	ld, _ := n.Codegen(prog)
	return ld.(*ir.InstLoad)
}

// GenAssign implement Assignable.GenAssign
func (n FunctionCallNode) GenAssign(prog *Program, _ value.Value, options ...AssignableOption) (value.Value, error) {
	return nil, fmt.Errorf("unable to assign to a function call")
}

// Type implement Assignable.Type
func (n FunctionCallNode) Type(prog *Program) (types.Type, error) {
	val, err := n.Codegen(prog)

	if val == nil {
		return nil, err
	}
	return val.Type(), err
}
