package typesystem

import (
	"github.com/llir/llvm/ir/types"
)

var (
	// GlobalTypeMap is the default typemap for geode
	GlobalTypeMap *TypeMap
)

// TypeMap is a mapping from type name to type object.
// This will be used inside the lexer to store custom
// type definitions
type TypeMap struct {
	Name  string
	Types map[string]*VarType // Types is a map of strings to instances (pointers) of VarType
}

func initGlobalTypeMap() {
	GlobalTypeMap = &TypeMap{}
	GlobalTypeMap.Types = make(map[string]*VarType)
}

// VarType represents a type in the language
// ie: int, float, string, etc.
type VarType struct {
	Name     string
	LLVMType types.Type
}

func (v VarType) String() string {
	return v.Name
}

// GetType returns the llvm type of some string name
func (t *TypeMap) GetType(name string) types.Type {
	v, ok := t.Types[name]
	if !ok {
		return nil
	}
	return v.LLVMType
}

// NewType takes some configuration and returns a VarType with those configs.
func NewType(name string, llvmType types.Type) *VarType {
	if GlobalTypeMap == nil {
		initGlobalTypeMap()
	}
	t := &VarType{}
	t.Name = name
	t.LLVMType = llvmType
	GlobalTypeMap.Types[t.Name] = t
	return t
}

// Declare the base types in the language
var (
	GeodeI8       = NewType("i8", types.I8)
	GeodeI32      = NewType("i32", types.I32)
	GeodeI64      = NewType("i64", types.I64)
	GeodeF32      = NewType("float", types.Float)
	GeodeF64      = NewType("double", types.Double)
	GeodeBool     = NewType("bool", types.I8)
	GeodeString   = NewType("string", types.NewPointer(GeodeI8.LLVMType))
	GeodeVoidType = NewType("void", types.Void)
	// GeodeI8 = NewType("char", false, 1, false, llvm.Int8Type())
)
