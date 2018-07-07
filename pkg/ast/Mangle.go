package ast

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir/types"
)

func init() {
	// fmt.Println("hello")
}

const functionNamePrefix = "_GN"

// MangleFunctionName will mangle a function name like how cpp does :)
func MangleFunctionName(origName string, argTypes ...types.Type) string {

	var name string

	name += functionNamePrefix

	name += strings.Replace(origName, ":", "$", -1)

	if len(argTypes) > 0 {
		name += "-"
	}
	for _, arg := range argTypes {
		name += fmt.Sprintf("%s", arg)
	}

	return name
}

// UnmangleFunctionName takes some mangled name and returns the unmangled one
func UnmangleFunctionName(mangled string) string {
	if mangled == "main" {
		return mangled
	}
	if !strings.HasPrefix(mangled, functionNamePrefix) {
		return mangled
	}

	mangled = mangled[len(functionNamePrefix):]
	// fmt.Println("MANGLED:", mangled)

	typePos := strings.Index(mangled, "-")
	if typePos != -1 {
		mangled = mangled[:typePos]
	}

	mangled = strings.Replace(mangled, "$", ":", -1)
	return mangled
}
