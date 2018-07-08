package ast

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nickwanninger/geode/pkg/util/log"
)

// FunctionNode is the representation of some function. It has methods
// on it to declare the function as well as codegen. A function has
// a list of VariableNodes for arguments and a single block for a body,
// all of which are codegenned.
type FunctionNode struct {
	NodeType

	Name       string
	Args       []VariableNode
	Body       BlockNode
	External   bool
	Variadic   bool
	Nomangle   bool
	ReturnType GeodeTypeRef
}

// NameString implements Node.NameString
func (n FunctionNode) NameString() string { return "FunctionNode" }

// InferType implements Node.InferType
func (n FunctionNode) InferType(scope *Scope) types.Type { return types.Void }

// Arguments returns some FunctionNode's arguments
func (n FunctionNode) Arguments(scope *Scope) ([]*types.Param, []types.Type) {
	funcArgs := make([]*types.Param, 0)
	argTypes := make([]types.Type, 0)
	for _, arg := range n.Args {
		ty := scope.FindType(arg.Type.Name).Type
		ty = arg.Type.BuildPointerType(ty)
		p := ir.NewParam(arg.Name, ty)
		funcArgs = append(funcArgs, p)
		argTypes = append(argTypes, p.Type())
	}
	return funcArgs, argTypes
}

// Declare declares some FunctionNode's sig
func (n FunctionNode) Declare(scope *Scope, c *Compiler) *ir.Function {

	funcArgs, _ := n.Arguments(scope)

	name := n.Name
	// We need to do some special checks if the function is main. It's special.
	// For instance, it must return int type.
	if name == "main" {
		if n.ReturnType.Name != "int" {
			log.Fatal("Main function must return type int. Called for type '%s'\n", n.ReturnType)
		}
	} else {
		name = n.MangledName(scope, c)
	}

	ty := scope.FindType(n.ReturnType.Name).Type
	ty = n.ReturnType.BuildPointerType(ty)
	function := c.Module.NewFunction(name, ty, funcArgs...)

	c.FN = function

	// if n.Variadic && !n.External {
	// 	log.Fatal("Function '%s' is variadic and has a body. This only allowed for external functions.\n", n.Name)
	// }

	function.Sig.Variadic = n.Variadic

	keyName := fmt.Sprintf("%s:%s", c.Package.NamespaceName, n.Name)
	// fmt.Println(function.Name, function.Sig.Variadic)
	scopeItem := NewFunctionScopeItem(keyName, function, PublicVisibility)
	scopeItem.SetMangled(!n.Nomangle)
	c.Scope.Add(scopeItem)

	return function
}

// MangledName will return the mangled name for a function node
func (n FunctionNode) MangledName(scope *Scope, c *Compiler) string {
	if n.Nomangle {
		return n.Name
	}
	_, argTypes := n.Arguments(scope)
	// Parse the namespace and name from the funciton name
	namespace, name := parseName(n.Name)
	if namespace == "" {
		namespace = c.Package.NamespaceName
	}

	n.Name = fmt.Sprintf("%s:%s", namespace, name)
	name = MangleFunctionName(n.Name, argTypes...)
	return name
}

// Codegen implements Node.Codegen for FunctionNode
func (n FunctionNode) Codegen(scope *Scope, c *Compiler) value.Value {

	name := n.Name

	if name != "main" {
		name = n.MangledName(scope, c)
	}

	declared := c.Scope.FindFunctions(name)
	if len(declared) != 1 {
		log.Fatal("Unable to find function declaration for '%s'\n", name)
	}
	function := declared[0].Value().(*ir.Function)
	c.FN = function

	// If the function is external (has ... at the end) we don't build a block
	if !n.External {
		// Create the entrypoint to the function
		entryBlock := ir.NewBlock("entry")
		c.FN.AppendBlock(entryBlock)
		c.PushBlock(entryBlock)

		if name == "main" {
			createPrelude(scope, c)
		}

		for _, arg := range function.Params() {
			alloc := c.CurrentBlock().NewAlloca(arg.Type())
			c.CurrentBlock().NewStore(arg, alloc)
			// Set the scope iteme
			scItem := NewVariableScopeItem(arg.Name, alloc, PrivateVisibility)
			scope.Add(scItem)
		}
		// Gen the body of the function
		n.Body.Codegen(scope, c)
		if c.CurrentBlock().Term == nil {
			ty := scope.FindType(n.ReturnType.Name).Type
			// log.Warn("Function %s is missing a return statement in the root block. Defaulting to 0\n", n.Name)
			v := createTypeCast(c, constant.NewInt(0, types.I64), ty)
			c.CurrentBlock().NewRet(v)
		}
		c.PopBlock()
	}

	return function
}

func createPrelude(scope *Scope, c *Compiler) {
	// Create a FunctionCallNode that will be the runtime init caller
	n := FunctionCallNode{}
	n.NodeType = nodeFunctionCall
	// Create a call to runtime init.
	n.Name = "_runtime:init"
	n.Codegen(scope, c)
}
