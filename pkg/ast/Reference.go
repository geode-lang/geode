package ast

import (
	"fmt"

	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/value"
)

// Reference is an interface that other things can implement that
// allows assignment, accessing, and gathering type information
type Reference interface {
	fmt.Stringer
	Accessable
	Assignable
	Alloca(s *Scope, c *Compiler) value.Value
	Load(block *ir.BasicBlock, s *Scope, c *Compiler) *ir.InstLoad
}
