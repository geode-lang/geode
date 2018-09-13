package types

import (
	"bytes"
	"fmt"

	"github.com/geode-lang/geode/llvm/enc"
)

// SliceType represents the structure of a struct in geode
type SliceType struct {
	StructType
	// Type name of identified struct type; or empty if struct type literal.
	Name string
	// Struct fields.
	Elem Type
}

// NewSlice returns a new struct type based on the given struct fields.
func NewSlice(elem Type) *SliceType {
	return &SliceType{
		Elem: elem,
	}
}

// String returns the LLVM syntax representation of the type.
func (t *SliceType) String() string {
	if t.Identified() {
		return enc.Local(t.Name)
	}
	return t.Def()
}

// Def returns the LLVM syntax representation of the definition of the type.
func (t *SliceType) Def() string {
	buf := &bytes.Buffer{}
	buf.WriteString("{ ")

	fmt.Fprintf(buf, "%s, %s", NewPointer(t.Elem), I64)
	buf.WriteString(" }")
	return buf.String()
}

// Equal reports whether t and u are of equal type.
func (t *SliceType) Equal(u Type) bool {
	if u, ok := u.(*SliceType); ok {
		// Identified struct types are uniqued by type names, not by structural
		// identity.
		if t.Identified() || u.Identified() {
			return t.Name == u.Name
		}
		// Literal struct types are uniqued by structural identity.
		if len(t.Fields) != len(u.Fields) {
			return false
		}
		for i, tf := range t.Fields {
			uf := u.Fields[i]
			if !tf.Equal(uf) {
				return false
			}
		}
		return true
	}
	return false
}

// GetName returns the name of the type.
func (t *SliceType) GetName() string {
	return t.Name
}

// SetName sets the name of the type.
func (t *SliceType) SetName(name string) {
	t.Name = name
}

// Identified reports whether t is an identified struct type.
func (t *SliceType) Identified() bool {
	return len(t.Name) > 0
}

// ByteCount returns the byte size of the type.
func (t *SliceType) ByteCount() int {
	var size int
	size += t.Elem.ByteCount()
	size += I64.ByteCount()
	return size
}
