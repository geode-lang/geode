package ast

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/geode-lang/geode/llvm/ir/types"
)

//go:generate stringer -type=ManglePartType

// ManglePartType -
type ManglePartType int

// Parts of a mangled name
const (
	NamespaceMangle ManglePartType = iota
	NameMangle
	TypeMangle
	ReturnTypeMangle
	GenericMangle
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
// so I will use X
const functionNamePrefix = "_X"
const globalVariableNamePrefix = "_V"
const separator = `:`

// MangleFunctionName will mangle a function name
func MangleFunctionName(origName string, types []types.Type, ret types.Type) string {

	buff := &bytes.Buffer{}

	fmt.Fprintf(buff, "%s", functionNamePrefix)

	parts := splitMany(origName, ":.")
	for i, p := range parts {
		prefix := "N"
		if i == 0 {
			prefix = "M"
		}
		fmt.Fprintf(buff, separator+"%s%s", prefix, p)
	}

	for _, t := range types {
		fmt.Fprintf(buff, separator+"T%s", t)
	}

	fmt.Fprintf(buff, separator+"R%s", ret)

	return buff.String()
}

// MangleVariableName will mangle a Variable name
func MangleVariableName(origName string) string {

	buff := &bytes.Buffer{}

	fmt.Fprintf(buff, "%s", globalVariableNamePrefix)

	parts := splitMany(origName, ":.")
	for i, p := range parts {
		prefix := "N"
		if i == 0 {
			prefix = "M"
		}
		fmt.Fprintf(buff, separator+"%s%s", prefix, p)
	}
	return buff.String()
}

// MangleMatches returns true if the two mangled names are:
//    a) the same namespace
//    b) the same name
//    c) the same generic structure.
func MangleMatches(a, b string) bool {

	aParts, aerr := GetMangleParts(a)
	bParts, berr := GetMangleParts(b)
	if aerr != nil || berr != nil {
		return false
	}

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
func GetMangleParts(mangled string) ([]ManglePart, error) {

	typeCharRefs := map[byte]ManglePartType{
		'M': NamespaceMangle,
		'N': NameMangle,
		'T': TypeMangle,
		'R': ReturnTypeMangle,
		'G': GenericMangle,
	}

	rawParts := strings.Split(mangled, separator)
	parts := make([]ManglePart, 0, len(rawParts))

	if rawParts[0] == functionNamePrefix {
		rawParts = rawParts[1:]
	} else {
		for _, p := range rawParts {
			part := ManglePart{}
			part.partType = NameMangle
			part.value = p
			parts = append(parts, part)
		}

		return parts, nil
	}

	for _, rawPart := range rawParts {
		typeChar := rawPart[0]

		typ, ok := typeCharRefs[typeChar]
		if !ok {
			return nil, fmt.Errorf("invalid typechar in mangled name %s: %c", mangled, typeChar)
		}

		part := ManglePart{}
		part.partType = typ
		part.value = rawPart[1:]
		parts = append(parts, part)
	}

	return parts, nil
}

// UnmangleFunctionName takes some mangled name and returns the unmangled one
func UnmangleFunctionName(mangled string) (string, error) {
	if mangled == "main" || !strings.HasPrefix(mangled, functionNamePrefix) {
		return mangled, nil
	}

	parsedParts, err := GetMangleParts(mangled)
	if err != nil {
		return "nil", err
	}

	if len(parsedParts) == 1 {
		return parsedParts[0].value, nil
	}

	buff := &bytes.Buffer{}
	for _, part := range parsedParts {
		if part.partType == NamespaceMangle {
			fmt.Fprintf(buff, "%s"+separator, part.value)
		}

		if part.partType == NameMangle {
			fmt.Fprintf(buff, "%s", part.value)
		}
	}

	return buff.String(), nil
}

// ParseName returns the namespace and the name of a string
func ParseName(combined string) (string, string) {
	var namespace, name string
	parts := strings.Split(combined, separator)
	name = parts[len(parts)-1]
	if len(parts) > 1 {
		namespace = parts[0]
	}

	return namespace, name
}
