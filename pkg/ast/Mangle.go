package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/llir/llvm/ir/types"
)

func init() {
	// fmt.Println("hello")
}

func readNumber(str string) int {
	var value int
	if _, err := fmt.Sscanf(str, "%d", &value); err != nil {
		value = 0
	}
	return value
}

func splitMany(s string, splits string) []string {
	m := make(map[rune]int)
	for _, r := range splits {
		m[r] = 1
	}

	splitter := func(r rune) bool {
		return m[r] == 1
	}

	return strings.FieldsFunc(s, splitter)
}

const functionNamePrefix = "_G"

// MangleFunctionName will mangle a function name like how cpp does :)
func MangleFunctionName(origName string, argTypes ...types.Type) string {

	name := functionNamePrefix
	parts := splitMany(origName, ":.")
	for i, p := range parts {
		prefix := ""
		if i == 0 {
			prefix = "N"
		}
		name += fmt.Sprintf("%s%d%s", prefix, len(p), p)
	}

	return name
}

// UnmangleFunctionName takes some mangled name and returns the unmangled one
func UnmangleFunctionName(mangled string) string {
	if mangled == "main" || !strings.HasPrefix(mangled, functionNamePrefix) {
		return mangled
	}

	parsedParts := make([]string, 0)
	for i := 0; i < len(mangled); i++ {
		c := mangled[i:]
		val := readNumber(c)
		if val != 0 {
			i += len(strconv.Itoa(val)) - 1
			start := i + 1
			end := start + val
			i += val
			parsedParts = append(parsedParts, mangled[start:end])
		}
	}

	// fmt.Println(mangled)

	name := fmt.Sprintf("%s:%s", parsedParts[0], parsedParts[1])

	return name
}
