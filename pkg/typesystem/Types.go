package typesystem

import (
	"github.com/geode-lang/geode/pkg/util/log"
	"github.com/geode-lang/geode/llvm/ir/types"
)

var (
	// GlobalTypeMap is the default typemap for geode
	// GlobalTypeMap *TypeMap

	defaultCastPrecs map[types.Type]int
)

func initializeMaps() {
	// GlobalTypeMap = &TypeMap{}
	// GlobalTypeMap.Types = make(map[string]*VarType)
	defaultCastPrecs = make(map[types.Type]int)
}

// CastPrecidence gives the precidence of some type as an integer
func CastPrecidence(t types.Type) int {

	if defaultCastPrecs == nil {
		defaultCastPrecs = make(map[types.Type]int)

		defaultCastPrecs[types.I32] = 32
	}

	p, ok := defaultCastPrecs[t]

	if !ok {
		log.Error("Invalid type precidence request\n")
		return -1
	}

	return p
}
