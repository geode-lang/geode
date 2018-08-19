package typesystem

import (
	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/geode-lang/llvm/ir/types"
)

var (
	// GlobalTypeMap is the default typemap for geode
	GlobalTypeMap *TypeMap

	castPrecidences map[types.Type]int
)

// Declare the base types in the language
var (
	GeodeBool     = NewType("bool", types.I1, 1)
	GeodeI8       = NewType("byte", types.I8, 2)
	GeodeI32      = NewType("i32", types.I32, 3)
	GeodeU64      = NewType("uint", types.I64, 4)
	GeodeI64      = NewType("int", types.I64, 4)
	GeodeF64      = NewType("float", types.Double, 11) // All floats are doubles
	GeodeString   = NewType("string", types.NewPointer(GeodeI8.LLVMType), 0)
	GeodeVoidType = NewType("void", types.Void, 0)
)

func initializeMaps() {
	GlobalTypeMap = &TypeMap{}
	GlobalTypeMap.Types = make(map[string]*VarType)
	castPrecidences = make(map[types.Type]int)
}

// TypeMap is a mapping from type name to type object.
// This will be used inside the lexer to store custom
// type definitions
type TypeMap struct {
	Name  string
	Types map[string]*VarType // Types is a map of strings to instances (pointers) of VarType
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

// CastPrecidence gives the precidence of some type as an integer
func CastPrecidence(t types.Type) int {
	p, ok := castPrecidences[t]

	if !ok {
		log.Error("Invalid type precidence request\n")
		return -1
	}

	return p
}

// NewType takes some configuration and returns a VarType with those configs.
func NewType(name string, llvmType types.Type, prec int) *VarType {

	if GlobalTypeMap == nil || castPrecidences == nil {
		initializeMaps()
	}
	t := &VarType{}
	t.Name = name
	t.LLVMType = llvmType
	castPrecidences[llvmType] = prec
	GlobalTypeMap.Types[t.Name] = t
	return t
}
