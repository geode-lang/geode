package ast

import (
	"fmt"

	"github.com/geode-lang/geode/llvm/ir"
	"github.com/geode-lang/geode/llvm/ir/value"
)

// Reference is an interface that other things can implement that
// allows assignment, accessing, and gathering type information
type Reference interface {
	fmt.Stringer
	Accessable
	Assignable
	Alloca(*Program) value.Value
	Load(*ir.BasicBlock, *Program) *ir.InstLoad
}
