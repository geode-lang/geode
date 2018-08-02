package ast

import (
	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/geode-lang/llvm/ir/constant"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// ArrayNode -
type ArrayNode struct {
	NodeType
	TokenReference

	Elements []Node
	Length   int
}

// NameString implements Node.NameString
func (n ArrayNode) NameString() string { return "ArrayNode" }

// InferType implements Node.InferType
func (n ArrayNode) InferType(scope *Scope) string { return "void" }

// Codegen implements Node.Codegen for ArrayNode
func (n ArrayNode) Codegen(scope *Scope, c *Compiler) value.Value {
	block := c.CurrentBlock()

	var elementType types.Type
	values := make([]value.Value, 0)
	for _, el := range n.Elements {
		val := el.Codegen(scope, c)
		if elementType == nil {
			elementType = val.Type()
		}

		if !types.Equal(val.Type(), elementType) {
			el.SyntaxError()
			log.Fatal("Invalid type in array. Array should be of a single type (%s)\n", elementType)
		}
		values = append(values, val)
	}

	typ := types.NewArray(values[0].Type(), int64(n.Length))

	alloca := block.NewAlloca(typ)

	zero := constant.NewInt(int64(0), types.I64)
	one := constant.NewInt(int64(1), types.I64)
	arrayStart := block.NewGetElementPtr(alloca, zero, zero)
	offset := arrayStart

	for i, val := range values {
		if i > 0 {
			offset = block.NewGetElementPtr(offset, one)
		}

		block.NewStore(val, offset)
	}

	return arrayStart
}
