package ast

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
)

// Scope trees represent block scoping by having a root scope
// and children scopes that point back to their parent scope.
type Scope struct {
	Parent   *Scope
	Children []*Scope
	Vals     map[string]ScopeItem
}

// Add a value to this specific scope
func (s *Scope) Add(val ScopeItem) {
	s.Vals[val.Name()] = val
}

// Find will traverse the scope tree to find some definition of a symbol
func (s *Scope) Find(name string) (ScopeItem, bool) {
	val, found := s.Vals[name]
	if !found && s.Parent != nil {
		return s.Parent.Find(name)
	}
	return val, found
}

// SpawnChild takes a parent scope and creates a new variable scope for scoped variable access.
func (s *Scope) SpawnChild() *Scope {
	n := &Scope{}
	n.Parent = s
	n.Vals = make(map[string]ScopeItem)
	n.Children = make([]*Scope, 0)
	s.Children = append(s.Children, n)
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
	value value.Value
	vis   Visibility
	name  string
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

// NewVariableScopeItem constructs a function scope item
func NewVariableScopeItem(name string, value value.Value, vis Visibility) VariableScopeItem {
	item := VariableScopeItem{}
	item.name = name
	item.value = value
	item.vis = vis
	return item
}
