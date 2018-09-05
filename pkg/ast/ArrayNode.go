package ast

import (
	"bytes"
	"fmt"

	"github.com/geode-lang/geode/llvm/ir/constant"
	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
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

// GenAccess -
func (n ArrayNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}

// Codegen implements Node.Codegen for ArrayNode
func (n ArrayNode) Codegen(prog *Program) (value.Value, error) {

	block := prog.Compiler.CurrentBlock()

	var elementType types.Type
	values := make([]value.Value, 0)
	for _, el := range n.Elements {
		val, err := el.Codegen(prog)
		if err != nil {
			return nil, err
		}
		if elementType == nil {
			elementType = val.Type()
		}

		if !types.Equal(val.Type(), elementType) {
			el.SyntaxError()
			return nil, fmt.Errorf("invalid type in array. Array should be of a single type '%s', got '%s'", elementType, val.Type())
		}
		values = append(values, val)
	}
	typ := prog.Compiler.PopType()

	if typ == nil {
		typ = types.NewPointer(values[0].Type())
	}

	itemType := typ.(*types.PointerType).Elem

	arrayType := types.NewArray(itemType, int64(n.Length))

	// arrayLength := int64(itemType.ByteCount() * n.Length)
	var alloca value.Value
	alloca = block.NewAlloca(arrayType)

	// alloca = createTypeCast(prog, alloca, arrayType)

	zero := constant.NewInt(int64(0), types.I64)
	one := constant.NewInt(int64(1), types.I64)
	arrayStart := block.NewGetElementPtr(alloca, zero, zero)
	offset := arrayStart

	for i, val := range values {
		block.AppendInst(NewLLVMComment("[%d] <- %s", i, val.Ident()))
		if i > 0 {
			offset = block.NewGetElementPtr(offset, one)
		}

		c, err := createTypeCast(prog, val, itemType)
		if err != nil {
			return nil, err
		}
		block.NewStore(c, offset)
	}

	return arrayStart, nil
}

func (n ArrayNode) String() string {
	buff := &bytes.Buffer{}
	fmt.Fprintf(buff, "ArrayNode")
	return buff.String()
}
