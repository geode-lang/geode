// Package gtypes declares the data types of Geode.
package gtypes

import "github.com/llir/llvm/ir/types"

// Type is a Geode type.
type Type interface {
	types.Type
	// Underlying returns the underlying LLVM IR type of the Geode type.
	Underlying() types.Type
}

// IsNumber reports whether the given type is an integer or floating-point type.
func IsNumber(t types.Type) bool {
	return types.IsInt(t) || types.IsFloat(t)
}

// IsSlice reports whether the given type is a Geode slice type.
func IsSlice(t types.Type) bool {
	_, ok := t.(*SliceType)
	return ok
}

// IsStruct reports whether the given type is a struct type.
func IsStruct(t types.Type) bool {
	if u, ok := t.(Type); ok {
		t = u.Underlying()
	}
	return types.IsStruct(t)
}
