package gtypes

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
)

// ByteCount returns the byte size of the type.
func ByteCount(t types.Type) int {
	switch t := t.(type) {
	case *types.ArrayType:
		return ArrayByteCount(t)
	case *types.StructType:
		return StructByteCount(t)
	case *types.VoidType:
		return VoidByteCount(t)
	case *types.FuncType:
		return FuncByteCount(t)
	case *types.LabelType:
		return LabelByteCount(t)
	case *types.MetadataType:
		return MetadataByteCount(t)
	case *types.IntType:
		return IntByteCount(t)
	case *types.FloatType:
		return FloatByteCount(t)
	case *types.PointerType:
		return PointerByteCount(t)
	case *types.VectorType:
		return VectorByteCount(t)
	case *SliceType:
		return SliceByteCount(t)
	default:
		panic(fmt.Errorf("support for type %T not yet implemented", t))
	}
}

// ArrayByteCount returns the byte size of the type.
func ArrayByteCount(t *types.ArrayType) int {
	return 8
}

// StructByteCount returns the byte size of the type.
func StructByteCount(t *types.StructType) int {
	var size int
	for _, ty := range t.Fields {
		size += ByteCount(ty)
	}
	return size
}

// VoidByteCount returns the byte size of the type.
func VoidByteCount(t *types.VoidType) int {
	return 8
}

// FuncByteCount returns the byte size of the type.
func FuncByteCount(t *types.FuncType) int {
	return 8
}

// LabelByteCount returns the byte size of the type.
func LabelByteCount(t *types.LabelType) int {
	return 8
}

// MetadataByteCount returns the byte size of the type.
func MetadataByteCount(t *types.MetadataType) int {
	return 0
}

// IntByteCount returns the byte size of the type.
func IntByteCount(t *types.IntType) int {
	return int(t.BitSize) / 8
}

// FloatByteCount returns the byte size of the type.
func FloatByteCount(t *types.FloatType) int {
	switch t.Kind {
	case types.FloatKindHalf:
		return 2
	case types.FloatKindFloat:
		return 4
	case types.FloatKindDouble:
		return 8
	case types.FloatKindFP128:
		return 16
	case types.FloatKindX86_FP80:
		return 10
	case types.FloatKindPPC_FP128:
		return 16
	default:
		panic(fmt.Errorf("support for floating-point kind %q not yet implemented", t.Kind))
	}
}

// PointerByteCount returns the byte size of the type.
func PointerByteCount(t *types.PointerType) int {
	return 8
}

// VectorByteCount returns the byte size of the type.
func VectorByteCount(t *types.VectorType) int {
	return ByteCount(t.ElemType) * int(t.Len)
}

// SliceByteCount returns the byte size of the type.
func SliceByteCount(t *SliceType) int {
	size := ByteCount(t.ElemType)
	size += ByteCount(types.I64)
	return size
}

// FloatBitSize returns the bit size of the given floating-point type.
func FloatBitSize(t *types.FloatType) int {
	switch t.Kind {
	case types.FloatKindHalf:
		return 16
	case types.FloatKindFloat:
		return 32
	case types.FloatKindDouble:
		return 64
	case types.FloatKindFP128:
		return 128
	case types.FloatKindX86_FP80:
		return 80
	case types.FloatKindPPC_FP128:
		return 128
	default:
		panic(fmt.Errorf("support for floating-point kind %q not yet implemented", t.Kind))
	}
}
