package gtypes

import (
	"github.com/llir/llvm/ir/types"
)

// StructType type is a Geode struct type.
type StructType struct {
	// Field names.
	Names []string

	// A Geode struct type is implemented as an LLVM struct type.
	*types.StructType
}

// NewStruct returns a new Geode struct type based on the given field types. The
// associated field names may be specified through t.Names.
func NewStruct(fields ...types.Type) *StructType {
	return &StructType{
		StructType: types.NewStruct(fields...),
	}
}

// FieldIndex returns the index of some field in the struct, or -1 if not
// present.
func (t *StructType) FieldIndex(name string) int {
	for i, n := range t.Names {
		if n == name {
			return i
		}
	}
	return -1
}
