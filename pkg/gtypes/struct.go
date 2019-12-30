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

// Underlying returns the underlying LLVM IR type of the Geode struct type.
func (t *StructType) Underlying() types.Type {
	return t.StructType
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

// Equal reports whether t and u are of equal type.
func (t *StructType) Equal(u types.Type) bool {
	if u, ok := u.(*StructType); ok {
		if len(t.Names) != len(u.Names) {
			return false
		}
		for i := range t.Names {
			if t.Names[i] != u.Names[i] {
				return false
			}
		}
		return t.StructType.Equal(u.StructType)
	}
	return u.Equal(t.StructType)
}
