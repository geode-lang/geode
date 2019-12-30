package gtypes

import (
	"github.com/llir/llvm/ir/types"
)

// SliceType type is a Geode slice type.
type SliceType struct {
	// Element type.
	ElemType types.Type

	// A Geode slice type is implemented as an LLVM struct type.
	//    { elem*, length }
	*types.StructType
}

// NewSlice returns a new Geode slice type based on the given element type.
func NewSlice(elem types.Type) *SliceType {
	length := types.I64
	typ := types.NewStruct(types.NewPointer(elem), length)
	return &SliceType{
		ElemType:   elem,
		StructType: typ,
	}
}

// Underlying returns the underlying LLVM IR type of the Geode slice type.
func (t *SliceType) Underlying() types.Type {
	return t.StructType
}

// Equal reports whether t and u are of equal type.
func (t *SliceType) Equal(u types.Type) bool {
	if u, ok := u.(*SliceType); ok {
		return t.StructType.Equal(u.StructType)
	}
	return u.Equal(t.StructType)
}
