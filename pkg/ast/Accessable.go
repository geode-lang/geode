package ast

import (
	"github.com/geode-lang/geode/llvm/ir/value"
)

// Accessable is an interface implementable by
// a node that allows the ability to read the value
// from the node.
type Accessable interface {
	GenAccess(*Program) (value.Value, error)
}
