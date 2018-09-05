package ast

import (
	"fmt"

	"github.com/geode-lang/geode/pkg/arg"
	"github.com/geode-lang/geode/llvm/ir"
	"github.com/geode-lang/geode/llvm/ir/constant"
	"github.com/geode-lang/geode/llvm/ir/types"
	"github.com/geode-lang/geode/llvm/ir/value"
)

// StringNode -
type StringNode struct {
	NodeType
	TokenReference

	Value string
}

// NameString implements Node.NameString
func (n StringNode) NameString() string { return "StringNode" }

var strIndex = 0

// Codegen implements Node.Codegen for StringNode
func (n StringNode) Codegen(prog *Program) (value.Value, error) {

	var str *ir.Global

	if found, exists := prog.StringDefs[n.Value]; exists {
		str = found
	} else {
		name := fmt.Sprintf(".str.%X", strIndex)
		strIndex++
		str = prog.Compiler.Module.NewGlobalDef(name, newCharArray(n.Value))
		str.IsConst = true
		str.Immutable()
		prog.StringDefs[n.Value] = str
	}

	var val value.Value
	zero := constant.NewInt(0, types.I32)
	val = constant.NewGetElementPtr(str, zero, zero)

	if !*arg.DisableStringDataCopy {
		length := constant.NewInt(int64(len([]byte(n.Value))+1), types.I32)
		v, err := prog.NewRuntimeFunctionCall("raw_copy", val, length)
		if err != nil {
			return nil, err
		}
		val = v
	}

	return val, nil
}

// GenAccess implements Accessable.GenAccess
func (n StringNode) GenAccess(prog *Program) (value.Value, error) {
	return n.Codegen(prog)
}

func (n StringNode) String() string {
	return fmt.Sprintf("%q", n.Value)
}
