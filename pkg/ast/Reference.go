package ast

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

// Reference is an interface that other things can implement that
// allows assignment, accessing, and gathering type information
type Reference interface {
	fmt.Stringer
	Accessable
	Assignable

	Alloca(*Program) value.Value
	Load(*ir.Block, *Program) *ir.InstLoad
}
