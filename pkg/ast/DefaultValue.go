package ast

import (
	"github.com/geode-lang/geode/llvm/ir/constant"
	"github.com/geode-lang/geode/llvm/ir/types"
)

// DefaultValue takes some type and returns a default value for it.
func DefaultValue(t types.Type) constant.Constant {

	var value constant.Constant

	if types.IsInt(t) {
		value = constant.NewInt(0, t)
	}

	if types.IsFloat(t) {
		value = constant.NewFloat(0.0, t)
	}

	if types.IsStruct(t) {

		structType := t.(*types.StructType)

		fields := make([]constant.Constant, 0, len(structType.Fields))
		for _, fieldType := range structType.Fields {
			fields = append(fields, DefaultValue(fieldType))
		}
		stct := constant.NewStruct(fields...)
		stct.Typ = structType
		value = stct
	}

	if types.IsPointer(t) {
		value = constant.NewNull(t)
	}

	return value
}
