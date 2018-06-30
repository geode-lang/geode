package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/nickwanninger/geode/pkg/util/log"
)

// Scope trees represent block scoping by having a root scope
// and children scopes that point back to their parent scope.
type Scope struct {
	Parent *Scope
	Vals   map[string]ScopeItem
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
func (s *Scope) FindFunctions(name string) []FunctionScopeItem {
	funcs := make([]FunctionScopeItem, 0)

	// First check the name without unmangling it
	fnc, found := s.Vals[name]
	if found {
		return append(funcs, fnc.(FunctionScopeItem))
	}

	unMangled := UnmangleFunctionName(name)
	fnc, found = s.Vals[unMangled]
	if found {
		return append(funcs, fnc.(FunctionScopeItem))
	}

	for _, v := range s.Vals {
		if v.Type() == ScopeItemFunctionType {
			fn := v.(FunctionScopeItem)
			if fn.Name() == name {
				funcs = append(funcs, fn)
			}
		}
	}
	if s.Parent != nil {
		funcs = append(funcs, s.Parent.FindFunctions(name)...)
	}

	return funcs

}

// SpawnChild takes a parent scope and creates a new variable scope for scoped variable access.
func (s *Scope) SpawnChild() *Scope {
	n := &Scope{}
	n.Parent = s
	n.Vals = make(map[string]ScopeItem)
	return n
}

// NewScope creates a scope (for use when generating root scopes)
func NewScope() *Scope {
	n := &Scope{}
	n.Parent = nil
	n.Vals = make(map[string]ScopeItem)
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

// NewFunctionScopeItem constructs a function scope item
func NewFunctionScopeItem(name string, function *ir.Function, vis Visibility) FunctionScopeItem {
	item := FunctionScopeItem{}
	item.name = name
	item.function = function
	item.vis = vis
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

// NewVariableScopeItem constructs a function scope item
func NewVariableScopeItem(name string, value value.Value, vis Visibility) VariableScopeItem {
	item := VariableScopeItem{}
	item.name = name
	item.value = value
	item.vis = vis
	return item
}

//
//
// TypeScopeItem implements ScopeItem.
// This is used so we can store functions in the scope (mainly in the root scope)
type TypeScopeItem struct {
	value   value.Value
	vis     Visibility
	name    string
	mangled bool
}

// Value implements ScopeItem.Value()
func (item TypeScopeItem) Value() value.Value {
	log.Error("Unable to get value from scope item type '%s'\n", item.Name())
	return nil
}

// Type implements ScopeItem.Type()
func (item TypeScopeItem) Type() ScopeItemType {
	return ScopeItemTypeType
}

// Visibility implements ScopeItem.Visibility()
func (item TypeScopeItem) Visibility() Visibility {
	return item.vis
}

// Name implements ScopeItem.Name()
func (item TypeScopeItem) Name() string {
	return item.name
}

// Mangled implements ScopeItem.Mangled()
func (item TypeScopeItem) Mangled() bool {
	return item.mangled
}

// SetMangled implements ScopeItem.SetMangled()
func (item TypeScopeItem) SetMangled(m bool) {
	item.mangled = m
}

// NewTypeScopeItem constructs a function scope item
func NewTypeScopeItem(name string, value value.Value, vis Visibility) TypeScopeItem {
	item := TypeScopeItem{}
	item.name = name
	item.value = value
	item.vis = vis
	return item
}
