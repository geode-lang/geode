package ast

import (
	"fmt"

	"github.com/geode-lang/geode/pkg/arg"
	"github.com/geode-lang/llvm/ir"
	"github.com/geode-lang/llvm/ir/constant"
	"github.com/geode-lang/llvm/ir/types"
	"github.com/geode-lang/llvm/ir/value"
)

// StringNode -
type StringNode struct {
	NodeType
	TokenReference

	Value string
}

// NameString implements Node.NameString
func (n StringNode) NameString() string { return "StringNode" }

// InferType implements Node.InferType
func (n StringNode) InferType(scope *Scope) string { return "string" }

var strIndex = 0

// Codegen implements Node.Codegen for StringNode
func (n StringNode) Codegen(prog *Program) value.Value {

	var str *ir.Global

	if found, exists := prog.StringDefs[n.Value]; exists {
		str = found
	} else {
		name := fmt.Sprintf(".str.%X", strIndex)
		strIndex++
		str = prog.Compiler.Module.NewGlobalDef(name, newCharArray(n.Value))
		str.IsConst = true
		prog.StringDefs[n.Value] = str
	}

	var val value.Value
	zero := constant.NewInt(0, types.I32)
	val = constant.NewGetElementPtr(str, zero, zero)

	if !*arg.DisableStringDataCopy {
		length := constant.NewInt(int64(len([]byte(n.Value))+1), types.I32)
		val = prog.NewRuntimeFunctionCall("raw_copy", val, length)
	}

	return val
}

// GenAccess implements Accessable.GenAccess
func (n StringNode) GenAccess(prog *Program) value.Value {
	return n.Codegen(prog)
}

func (n StringNode) String() string {
	return fmt.Sprintf("%q", n.Value)
}
