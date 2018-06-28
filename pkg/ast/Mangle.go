package ast

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir/types"
	"github.com/nickwanninger/geode/pkg/util/log"
)

func init() {
	// fmt.Println("hello")
}

const functionmangleprefix = "__geode_"

// MangleFunctionName will mangle a function name like how cpp does :)
func MangleFunctionName(origName string, argTypes ...types.Type) string {
	name := fmt.Sprintf("%s", origName)

	if len(argTypes) > 0 {
		name += functionmangleprefix
	}
	for _, arg := range argTypes {
		name += fmt.Sprintf("%s", arg)
	}

	return name
}

// UnmangleFunctionName takes some mangled name and returns the unmangled one
func UnmangleFunctionName(mangled string) string {
	pos := strings.Index(mangled, functionmangleprefix)
	if pos == -1 {
		log.Fatal("Unable to unmangle function name '%s'\n", mangled)
	}
	return mangled[0:pos]
}
