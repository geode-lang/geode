package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// Callable is for the left side of a function call. It has functions for getting the function that it points to, etc...
type Callable interface {
	GetFunc(*Program, []types.Type) (*ir.Function, []value.Value, error)
}
