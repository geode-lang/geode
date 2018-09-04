package ast

import (
	"bytes"
	"fmt"

	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
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

	var name string
	var err error

	args := []value.Value{}
	argTypes := []types.Type{}

	// argStrings := []string{}
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
	callee, err := n.Name.GetFunc(prog, argTypes)
	if err != nil {
		return nil, err
	}

	if callee == nil {
		return nil, fmt.Errorf("unknown function %q referenced at %s", name, n.Token.FileInfo())
	}

	// Attempt to typecast all the args into the correct type
	// This is skipped with variadic functions
	if !callee.Sig.Variadic {
		for i := range args {
			args[i], _ = createTypeCast(prog, args[i], callee.Sig.Params[i].Type())
		}
	}

	// Varargs require type conversion to a standardized type
	// So we will use the same type promotion c uses
	//  if int && type != i32 -> type = i32
	//  if fnn && type != f64 -> type = f64
	arguments := make([]value.Value, 0, len(args))

	for i, arg := range args {

		if callee.Sig.Variadic && i >= len(callee.Params()) {
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

	// prog.Compiler.CurrentBlock().AppendInst(NewLLVMComment("%s", n.String()))
	return prog.Compiler.CurrentBlock().NewCall(callee, arguments...), nil
}
