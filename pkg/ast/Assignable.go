package ast

import (
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// Assignable is an interface that a node can implement that
// allows the definition of what happens when you attempt to
// assign a value to it.
type Assignable interface {
	GenAssign(*Scope, *Compiler, value.Value) value.Value
	Type(*Scope) types.Type
}
