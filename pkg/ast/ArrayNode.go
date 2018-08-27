package ast

import (
	"bytes"
	"fmt"

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

// GenAccess -
func (n ArrayNode) GenAccess(prog *Program) value.Value {
	return n.Codegen(prog)
}

// Codegen implements Node.Codegen for ArrayNode
func (n ArrayNode) Codegen(prog *Program) value.Value {

	block := prog.Compiler.CurrentBlock()

	var elementType types.Type
	values := make([]value.Value, 0)
	for _, el := range n.Elements {
		val := el.Codegen(prog)
		if elementType == nil {
			elementType = val.Type()
		}

		if !types.Equal(val.Type(), elementType) {
			el.SyntaxError()
			log.Fatal("Invalid type in array. Array should be of a single type '%s', got '%s'\n", elementType, val.Type())
		}
		values = append(values, val)
	}
	typ := prog.Compiler.typeCache

	itemType := typ.(*types.PointerType).Elem

	arrayType := types.NewArray(itemType, int64(n.Length))

	alloca := block.NewAlloca(arrayType)

	zero := constant.NewInt(int64(0), types.I64)
	one := constant.NewInt(int64(1), types.I64)
	arrayStart := block.NewGetElementPtr(alloca, zero, zero)
	offset := arrayStart

	for i, val := range values {
		block.AppendInst(NewLLVMComment("[%d] <- %s", i, val.Ident()))
		if i > 0 {
			offset = block.NewGetElementPtr(offset, one)
		}

		block.NewStore(createTypeCast(prog, val, itemType), offset)
	}

	return arrayStart
}

func (n ArrayNode) String() string {
	buff := &bytes.Buffer{}
	fmt.Fprintf(buff, "ArrayNode")
	return buff.String()
}
