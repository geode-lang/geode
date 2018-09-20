package ast

import (
	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
)

// Assignable is an interface that a node can implement that
// allows the definition of what happens when you attempt to
// assign a value to it.
type Assignable interface {
	GenAssign(*Program, value.Value, ...AssignableOption) (value.Value, error)
	Type(*Program) (types.Type, error)
}

// AssignableOption is a type that will be passed into the GenAssign function for options
type AssignableOption int

// Declare the various options for a GenAssign call
const (
	AssignableDeclare AssignableOption = iota
)
