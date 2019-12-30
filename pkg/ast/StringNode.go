package ast

import (
	"fmt"

	"github.com/geode-lang/geode/pkg/arg"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
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
		str.Immutable = true
		prog.StringDefs[n.Value] = str
	}

	var val value.Value
	zero := constant.NewInt(types.I32, 0)
	val = constant.NewGetElementPtr(str.ContentType, str, zero, zero)

	if !*arg.DisableStringDataCopy {
		length := constant.NewInt(types.I32, int64(len([]byte(n.Value))+1))
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
