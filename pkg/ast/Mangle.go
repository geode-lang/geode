package ast

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/geode-lang/llvm/ir/types"
)

//go:generate stringer -type=ManglePartType

// ManglePartType -
type ManglePartType int

// Parts of a mangled name
const (
	NamespaceMangle ManglePartType = iota
	NameMangle      ManglePartType = iota
	GenericMangle   ManglePartType = iota
)

// ManglePart -
type ManglePart struct {
	value    string
	partType ManglePartType
}

// MangleName is a method container around a list of MangleParts
type MangleName []ManglePart

func (m MangleName) String() string {
	return ""
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

// This is the prefix that will prefix all function names.
// c++ uses _Z, I feel like an equally random value works fine
// so I will use V
const functionNamePrefix = "_V"

// MangleFunctionName will mangle a function name
func MangleFunctionName(origName string, argTypes []types.Type, generics []*GenericSymbol) string {

	name := &bytes.Buffer{}

	fmt.Fprintf(name, "%s", functionNamePrefix)

	parts := splitMany(origName, ":.")
	for i, p := range parts {
		prefix := "N"
		if i == 0 {
			prefix = "M"
		}
		fmt.Fprintf(name, "%s%d%s", prefix, len(p), p)
	}

	for _, gen := range generics {
		fmt.Fprintf(name, "G%d%s", len(gen.Name), gen.Name)
	}

	return name.String()
}

// MangleMatches returns true if the two mangled names are:
//    a) the same namespace
//    b) the same name
//    c) the same generic structure.
func MangleMatches(a, b string) bool {

	aParts := GetMangleParts(a)
	bParts := GetMangleParts(b)

	if len(aParts) != len(bParts) {
		return false
	}

	for i := range aParts {
		ap := aParts[i]
		bp := bParts[i]

		if ap.partType != bp.partType {
			return false
		}

		// Only check if the names are the same if the current
		// char type is not a generic
		if ap.partType == GenericMangle {
			continue
		}

		if ap.value != bp.value {
			return false
		}
	}

	return true

}

// GetMangleParts parses and returns mangleparts for some mangled name
func GetMangleParts(mangled string) []ManglePart {

	typeCharRefs := map[byte]ManglePartType{
		'M': NamespaceMangle,
		'N': NameMangle,
		'G': GenericMangle,
	}

	parts := make([]ManglePart, 0)
	for i := 0; i < len(mangled); i++ {
		c := mangled[i:]
		val := readNumber(c)
		if val != 0 {

			typeChar := mangled[i-1]

			typ, ok := typeCharRefs[typeChar]
			if !ok {
				log.Fatal("Invalid typechar in mangled name %s: %c\n", mangled, typeChar)
			}

			i += len(strconv.Itoa(val)) - 1
			start := i + 1
			end := start + val
			i += val

			part := ManglePart{}

			part.partType = typ
			part.value = mangled[start:end]

			parts = append(parts, part)

		}
	}
	return parts
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

	// GetMangleParts(mangled)

	name := fmt.Sprintf("%s:%s", parsedParts[0], parsedParts[1])

	return name
}

// func
