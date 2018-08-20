package ast

import (
	"github.com/geode-lang/llvm/ir"
)

// Addressable is an interface implementable by
// a node that allows the ability to read address
// of a node.
type Addressable interface {
	GenAddress(*Program) *ir.InstAlloca
}
