package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nickwanninger/geode/pkg/util/log"
)

func init() {
	// gob.Register(Scope{})
	// gob.Register(FunctionScopeItem{})
	// gob.Register(VariableScopeItem{})
}

// Scope trees represent block scoping by having a root scope
// and children scopes that point back to their parent scope.
type Scope struct {
	Parent *Scope
	Vals   map[string]ScopeItem
	Types  *map[string]*TypeDef
}

// Add a value to this specific scope
func (s *Scope) Add(val ScopeItem) {
	s.Vals[val.Name()] = val
}

// Find will traverse the scope tree to find some definition of a symbol
func (s *Scope) Find(name string) (ScopeItem, bool) {
	for _, v := range s.Vals {
		u := v.Name()
		if u == name {
			return v, true
		}
	}
	if s.Parent != nil {
		return s.Parent.Find(name)
	}
	return nil, false
}

// FindFunctions returns a list of functions that might match the name provided
// The needle can be any of the following: bare name, mangled name
func (s *Scope) FindFunctions(needle string) []FunctionScopeItem {

	// fmt.Println("===========================")
	funcs := make([]FunctionScopeItem, 0)

	unMangled := UnmangleFunctionName(needle)

	_, name := parseName(unMangled)
	// fmt.Println(needle, " -> ", unMangled)
	// fmt.Println(needle, " -> ", name)

	// possibleNames := []string{needle, unMangled, name}
	// fmt.Println(needle)
	for _, v := range s.Vals {
		// if the function is not mangled, check specially
		if name == unMangled {
			if v.Name() == name {
				// fmt.Println("IS MAIN")
				return append(funcs, v.(FunctionScopeItem))
			}
			continue

		}
		// fmt.Println(v.Name(), needle, name, MangleMatches(needle, v.Name()))
		if v.Name() == needle || v.Name() == name || MangleMatches(needle, v.Name()) {
			// fmt.Println(" +", v.Name())
			funcs = append(funcs, v.(FunctionScopeItem))
		} else {
			// fmt.Println(" -", v.Name())
		}

	}

	// fmt.Println("===========================")

	// for _, nm := range possibleNames {
	// 	fnc, found := s.Vals[nm]
	// 	if found {
	// 		return append(funcs, fnc.(FunctionScopeItem))
	// 	}
	// }

	// for _, v := range s.Vals {
	// 	if v.Type() == ScopeItemFunctionType {
	// 		fn := v.(FunctionScopeItem)
	// 		if fn.Name() == needle {
	// 			funcs = append(funcs, fn)
	// 		}
	// 		if UnmangleFunctionName(fn.Name()) == unMangled {
	// 			funcs = append(funcs, fn)

	// 		}
	// 	}
	// }

	if s.Parent != nil {
		funcs = append(funcs, s.Parent.FindFunctions(needle)...)
	}

	return funcs
}

// FindType returns the type stored with a name in this scope
func (s *Scope) FindType(name string) *TypeDef {

	v, ok := (*s.Types)[name]
	if !ok {
		if s.Parent == nil {
			log.Fatal("Unable to find type with name '%s' in scope\n", name)
		}
		return s.Parent.FindType(name)
	}
	return v
}

// InjectPrimitives injects primitve types like int, byte, etc
func (s *Scope) InjectPrimitives() {
	NewTypeDef("byte", types.I8, 1).InjectInto(s)
	NewTypeDef("i16", types.I16, 2).InjectInto(s)
	NewTypeDef("i32", types.I32, 3).InjectInto(s)
	NewTypeDef("int", types.I64, 4).InjectInto(s)
	NewTypeDef("big", types.NewInt(255), 100).InjectInto(s)
	NewTypeDef("float", types.Double, 11).InjectInto(s)
	NewTypeDef("string", types.NewPointer(types.I8), 0).InjectInto(s)
	NewTypeDef("void", types.Void, 0).InjectInto(s)
}

// SpawnChild takes a parent scope and creates a new variable scope for scoped variable access.
func (s *Scope) SpawnChild() *Scope {
	n := &Scope{}
	n.Parent = s
	n.Vals = make(map[string]ScopeItem)
	n.Types = s.Types
	return n
}

// NewScope creates a scope (for use when generating root scopes)
func NewScope() *Scope {
	n := &Scope{}
	n.Parent = nil
	n.Vals = make(map[string]ScopeItem)
	typemap := make(map[string]*TypeDef)
	n.Types = &typemap
	return n
}

// ScopeItem is what the scope contains.ScopeItem
// When you set a value in the scope, you must wrap
// it in a struct implementing ScopeItem
type ScopeItem interface {
	Type() ScopeItemType
	Value() value.Value // an llvm value
	Visibility() Visibility
	Name() string
	Mangled() bool
	SetMangled(m bool)
	Node() Node
}

// ScopeItemType -
type ScopeItemType int

// The ScopeItemTypes available
const (
	ScopeItemFunctionType ScopeItemType = iota
	ScopeItemVariableType
	ScopeItemTypeType
)

// Visibility is the access modifier of some scope variable.
type Visibility int

// Visibility options for scope items
const (
	PublicVisibility Visibility = iota
	PrivateVisibility
)

//
//
// FunctionScopeItem implements ScopeItem.
// This is used so we can store functions in the scope (mainly in the root scope)
type FunctionScopeItem struct {
	function *ir.Function
	vis      Visibility
	name     string
	mangled  bool
	node     FunctionNode
}

// Value implements ScopeItem.Value()
func (item FunctionScopeItem) Value() value.Value {
	return item.function
}

// Type implements ScopeItem.Type()
func (item FunctionScopeItem) Type() ScopeItemType {
	return ScopeItemFunctionType
}

// Visibility implements ScopeItem.Visibility()
func (item FunctionScopeItem) Visibility() Visibility {
	return item.vis
}

// Name implements ScopeItem.Name()
func (item FunctionScopeItem) Name() string {
	return item.name
}

// Mangled implements ScopeItem.Mangled()
func (item FunctionScopeItem) Mangled() bool {
	return item.mangled
}

// SetMangled implements ScopeItem.SetMangled()
func (item FunctionScopeItem) SetMangled(m bool) {
	item.mangled = m
}

// Node implements ScopeItem.Node()
func (item FunctionScopeItem) Node() Node {
	return item.node
}

// NewFunctionScopeItem constructs a function scope item
func NewFunctionScopeItem(name string, node FunctionNode, function *ir.Function, vis Visibility) FunctionScopeItem {
	item := FunctionScopeItem{}
	item.name = function.Name
	item.function = function
	item.vis = vis
	item.node = node
	return item
}

//
//
// VariableScopeItem implements ScopeItem.
// This is used so we can store functions in the scope (mainly in the root scope)
type VariableScopeItem struct {
	value   value.Value
	vis     Visibility
	name    string
	mangled bool
	node    VariableNode
}

// Value implements ScopeItem.Value()
func (item VariableScopeItem) Value() value.Value {
	return item.value
}

// Type implements ScopeItem.Type()
func (item VariableScopeItem) Type() ScopeItemType {
	return ScopeItemVariableType
}

// Visibility implements ScopeItem.Visibility()
func (item VariableScopeItem) Visibility() Visibility {
	return item.vis
}

// Name implements ScopeItem.Name()
func (item VariableScopeItem) Name() string {
	return item.name
}

// Mangled implements ScopeItem.Mangled()
func (item VariableScopeItem) Mangled() bool {
	return item.mangled
}

// SetMangled implements ScopeItem.SetMangled()
func (item VariableScopeItem) SetMangled(m bool) {
	item.mangled = m
}

// Node implements ScopeItem.Node()
func (item VariableScopeItem) Node() Node {
	return item.node
}

// NewVariableScopeItem constructs a function scope item
func NewVariableScopeItem(name string, value value.Value, vis Visibility) VariableScopeItem {
	item := VariableScopeItem{}
	item.name = name
	item.value = value
	item.vis = vis
	return item
}

// TypeDef is a storage for types in the scope. They are stored seperately from variables.
type TypeDef struct {
	Type types.Type
	Name string
	Prec int
}

// NewTypeDef constructs a function scope item
func NewTypeDef(name string, t types.Type, prec int) *TypeDef {
	item := &TypeDef{}
	item.Name = name
	item.Type = t
	item.Prec = prec
	return item
}

// InjectInto will inject the type into a given scope
func (t *TypeDef) InjectInto(s *Scope) {
	(*s.Types)[t.Name] = t
}
