package ast

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// ClassNode -
type ClassNode struct {
	NodeType

	Name      string
	Methods   []FunctionNode
	Variables []VariableNode
}

// NameString implements Node.NameString
func (n ClassNode) NameString() string { return "ClassNode" }

// InferType implements Node.InferType
func (n ClassNode) InferType(scope *Scope) types.Type { return types.Void }

// Codegen implements Node.Codegen for ClassNode
func (n ClassNode) Codegen(scope *Scope, c *Compiler) value.Value {

	fields := make([]types.Type, 0)

	for _, f := range n.Variables {
		t := f.Type.Name
		ty := scope.FindType(t).Type
		fields = append(fields, ty)
	}

	structDefn := types.NewStruct(fields...)

	NewTypeDef(n.Name, structDefn, -1).InjectInto(scope)
	// structDefn.Opaque = true
	structDefn.SetName(n.Name)
	c.Module.NewType(n.Name, structDefn)
	// fmt.Println(t, structDefn)
	return nil
}
