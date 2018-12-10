// Package gtypes declares the data types of Geode.
package gtypes

import "github.com/llir/llvm/ir/types"

// IsNumber reports whether the given type is an integer or floating-point type.
func IsNumber(t types.Type) bool {
	return types.IsInt(t) || types.IsFloat(t)
}
