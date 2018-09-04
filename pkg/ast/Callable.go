package ast

import (
	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/types"
)

// Callable is for the left side of a function call. It has functions for getting the function that it points to, etc...
type Callable interface {
	GetFunc(*Program, []types.Type) (*ir.Function, error)
}
