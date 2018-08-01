package ast

import "github.com/geode-lang/llvm/ir/value"

// Accessable is an interface implementable by
// a node that allows the ability to read the value
// from the node.
type Accessable interface {
	GenAccess(*Scope, *Compiler) value.Value
}
