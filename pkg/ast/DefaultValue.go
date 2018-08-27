package ast

import (
	"github.com/geode-lang/llvm/ir/constant"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// DefaultValue takes some type and returns a default value for it.
func DefaultValue(t types.Type) value.Value {

	if types.IsInt(t) {
		return constant.NewInt(0, t.(*types.IntType))
	}

	if types.IsFloat(t) {
		return constant.NewFloat(0, t.(*types.FloatType))
	}

	if types.IsStruct(t) {
		return nil
	}

	return nil
}
