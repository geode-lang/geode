package types

import (
	"github.com/go-llvm/llvm"
)

var (
	// GlobalTypeMap is a global mapping for this build
	GlobalTypeMap *TypeMap
)

func init() {
	GlobalTypeMap = NewTypeMap("global")
}

// VarType represents a type in the language
// ie: int, float, string, etc.
type VarType struct {
	Name      string
	IsStruct  bool
	IsPointer bool
	Size      int // How many bytes does one item take up?
	Fields    map[string]*VarType
	LLVMType  llvm.Type
}

func (v VarType) String() string {
	return v.Name
}

// NewType takes some configuration and returns a VarType with those configs.
func NewType(name string, isStruct bool, size int, isPointer bool, llvmType llvm.Type) *VarType {
	t := &VarType{}
	t.Name = name
	t.IsStruct = isStruct
	t.IsPointer = isPointer
	t.Size = size
	t.LLVMType = llvmType
	return t
}

// Declare the base types in the language
var (
	DefaultIntType    = NewType("int", false, 8, false, llvm.Int64Type())
	DefaultFloatType  = NewType("float", false, 8, false, llvm.DoubleType())
	DefaultCharType   = NewType("char", false, 8, false, llvm.Int8Type())
	DefaultStringType = NewType("string", false, 8, false, llvm.PointerType(llvm.Int8Type(), 0))
	DefaultVoidType   = NewType("void", false, 8, false, llvm.VoidType())
	// DefaultStringType = NewType("string", false, 1, true, llvm.Type{}) // strings are a pointer to a byte, with a size of 1
	// DefaultCharType = NewType("char", false, 1, false, llvm.Int8Type())
)

// TypeMap is a mapping from type name to type object.
// This will be used inside the parser to store custom
// type definitions
type TypeMap struct {
	Name  string
	Types map[string]*VarType // Types is a map of strings to instances (pointers) of VarType
}

// PopulateDefaults fills the TypeMap with default Types for the langauge
func (m *TypeMap) PopulateDefaults() {
	m.Types["int"] = DefaultIntType
	m.Types["float"] = DefaultFloatType
	m.Types["char"] = DefaultCharType
	m.Types["string"] = DefaultStringType
	m.Types["void"] = DefaultVoidType
	// m.Types["string"] = DefaultStringType
	// m.Types["char"] = DefaultStringType
}

// GetType takes a type name and returns a reference to that type as a pointer
func (m *TypeMap) GetType(name string) *VarType {
	return m.Types[name]
}

// // TypeReference is an internal link to a certain type inside the program
// type TypeReference struct {
// 	Map  *TypeMap
// 	Type *VarType
// }

// NewTypeMap returns an empty typemap containing
// no default type representations
func NewTypeMap(name string) *TypeMap {
	t := &TypeMap{}
	t.Name = name
	t.Types = make(map[string]*VarType)
	t.PopulateDefaults()
	return t
}
