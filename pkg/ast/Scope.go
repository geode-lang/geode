package ast

import (
	"encoding/json"
	"fmt"

	"github.com/geode-lang/geode/pkg/util"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/metadata"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func init() {
}

// Scope trees represent block scoping by having a root scope
// and children scopes that point back to their parent scope.
type Scope struct {
	Parent      *Scope                `json:"-"`
	Index       int                   `json:"index"`
	Children    []*Scope              `json:"children"`
	Vals        map[string]ScopeItem  `json:"values"`
	Types       map[string]*ScopeType `json:"types"`
	PackageName string                `json:"package_name"`
	DebugInfo   *metadata.DISubprogram
}

// Add a value to this specific scope
func (s *Scope) Add(val ScopeItem) {
	s.Vals[val.Name()] = val
}

// Find will traverse the scope tree to find some definition of a symbol
func (s *Scope) Find(searchPaths []string) (ScopeItem, bool) {
	for _, v := range s.Vals {
		for _, name := range searchPaths {
			if v.Name() == name {
				return v, true
			}
		}
	}
	if s.Parent != nil {
		return s.Parent.Find(searchPaths)
	}
	return nil, false
}

// GetSimilarName returns the most similar name in the parent scopes
func (s *Scope) GetSimilarName(name string) (string, float64) {
	names := s.GetNames()
	tmp := s.Parent
	for tmp != nil {
		names = append(names, tmp.GetNames()...)
		tmp = tmp.Parent
	}
	type info struct {
		name string
		dist float64
	}
	closest := info{
		name: names[0],
		dist: util.Jaro(names[0], name),
	}
	for _, n := range names[1:] {
		d := util.Jaro(n, name)
		if d > closest.dist {
			closest.dist = d
			closest.name = n
		}
	}
	return closest.name, closest.dist
}

// AllNames returns a recursive lookup of all names in a scope tree
func (s *Scope) AllNames() []string {
	return s.GetRoot().GetNames()
}

// GetNames returns just the names for this scope
func (s *Scope) GetNames() []string {
	names := make([]string, 0)

	for k := range s.Vals {
		names = append(names, k)
	}
	for _, c := range s.Children {
		names = append(names, c.GetNames()...)
	}

	return names
}

// FindFunctions returns a list of functions that might match the name provided
// The needle can be any of the following: bare name, mangled name
func (s *Scope) FindFunctions(needle string) ([]FunctionScopeItem, []GenericTemplateScopeItem, error) {

	funcs := make([]FunctionScopeItem, 0)
	generics := make([]GenericTemplateScopeItem, 0)

	unMangled, err := UnmangleFunctionName(needle)
	if err != nil {
		return nil, nil, err
	}

	_, name := ParseName(unMangled)

	for _, v := range s.Vals {
		// if the function is not mangled, check specially
		if name == unMangled {
			if v.Name() == name {
				return append(funcs, v.(FunctionScopeItem)), generics, nil
			}
			continue
		}
		if v.Name() == needle || v.Name() == name || MangleMatches(needle, v.Name()) {
			funcs = append(funcs, v.(FunctionScopeItem))
		}

	}

	if s.Parent != nil {
		fn, gn, err := s.Parent.FindFunctions(needle)
		if err != nil {
			return nil, nil, err
		}
		funcs = append(funcs, fn...)
		generics = append(generics, gn...)
	}

	return funcs, generics, nil
}

// FindType returns the type stored with a name in this scope
func (s *Scope) FindType(names ...string) *ScopeType {
	var v *ScopeType
	var ok bool
	for _, name := range names {

		v, ok = s.Types[name]
		if ok {
			return v
		}
	}
	if s.Parent == nil {
		return nil
	}
	return s.Parent.FindType(names...)
}

// FindTypeName returns the geode defined type name
// for an llvm type representation
func (s *Scope) FindTypeName(t types.Type) (string, error) {
	for _, val := range s.Types {
		if types.Equal(val.Type, t) {
			return val.Name, nil
		}
	}
	if s.Parent == nil {
		return "", fmt.Errorf("unable to find type %s in any scope", t)
	}
	return s.Parent.FindTypeName(t)
}

func (s *Scope) String() string {
	j, _ := json.MarshalIndent(s, "", "    ")
	return string(j)
}

// GetTypeName takes a type and returns the human name
// that the compiler and lexer understands
func (s *Scope) GetTypeName(t types.Type) string {
	return "void"
}

// InjectPrimitives injects primitve types like int, byte, etc
func (s *Scope) InjectPrimitives() {
	s.RegisterType("bool", types.I1, 1)
	s.RegisterType("byte", types.I8, 2)
	s.RegisterType("short", types.I16, 3)
	s.RegisterType("int", types.I32, 4)
	s.RegisterType("long", types.I64, 5)

	s.RegisterType("big", types.NewInt(128), 128)
	s.RegisterType("large", types.NewInt(256), 256)
	s.RegisterType("huge", types.NewInt(512), 512)

	s.RegisterType("float", types.Double, 11)
	s.RegisterType("string", types.NewPointer(types.I8), 0)
	s.RegisterType("void", types.Void, 0)
}

// RegisterType takes information about some type and binds it to this scope
func (s *Scope) RegisterType(name string, t types.Type, prec int) {
	s.Types[name] = NewScopeType(name, t, prec)
}

// SpawnChild takes a parent scope and creates a new variable scope for scoped variable access.
func (s *Scope) SpawnChild() *Scope {
	child := NewScope()
	child.Parent = s
	child.Vals = make(map[string]ScopeItem)
	child.Types = make(map[string]*ScopeType)
	child.PackageName = s.PackageName
	s.Children = append(s.Children, child)
	return child
}

// GetRoot returns the root of the scope
func (s *Scope) GetRoot() *Scope {
	scope := s

	for scope.Parent != nil {
		scope = scope.Parent
	}
	return scope
}

var scopeIndex = 0

// NewScope creates a scope (for use when generating root scopes)
func NewScope() *Scope {
	n := &Scope{}
	n.Index = scopeIndex
	scopeIndex++
	n.Parent = nil
	n.Vals = make(map[string]ScopeItem)
	n.Types = make(map[string]*ScopeType)
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
// GenericTemplateScopeItem implements ScopeItem.
// This is used so we can store functions in the scope (mainly in the root scope)
type GenericTemplateScopeItem struct {
	function *ir.Func
	vis      Visibility
	name     string
	types    []TypeNode
	node     Node
	mangled  bool
}

// Value implements ScopeItem.Value()
func (item GenericTemplateScopeItem) Value() value.Value {
	return item.function
}

// Type implements ScopeItem.Type()
func (item GenericTemplateScopeItem) Type() ScopeItemType {
	return ScopeItemFunctionType
}

// Visibility implements ScopeItem.Visibility()
func (item GenericTemplateScopeItem) Visibility() Visibility {
	return item.vis
}

// Name implements ScopeItem.Name()
func (item GenericTemplateScopeItem) Name() string {
	return item.name
}

// Mangled implements ScopeItem.Mangled()
func (item GenericTemplateScopeItem) Mangled() bool {
	return item.mangled
}

// SetMangled implements ScopeItem.SetMangled()
func (item GenericTemplateScopeItem) SetMangled(m bool) {
	item.mangled = m
}

// Node implements ScopeItem.Node()
func (item GenericTemplateScopeItem) Node() Node {
	return item.node
}

// NewGenericTemplateScopeItem constructs a function scope item
func NewGenericTemplateScopeItem(name string) GenericTemplateScopeItem {
	item := GenericTemplateScopeItem{}
	return item
}

//
//
// FunctionScopeItem implements ScopeItem.
// This is used so we can store functions in the scope (mainly in the root scope)
type FunctionScopeItem struct {
	function *ir.Func
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
func NewFunctionScopeItem(name string, node FunctionNode, function *ir.Func, vis Visibility) FunctionScopeItem {
	item := FunctionScopeItem{}
	item.name = function.Name()
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
	value    value.Value
	vis      Visibility
	name     string
	mangled  bool
	node     VariableNode
	varIndex int
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

// MangledName returns the instance unique name for this variable (to fix the "variable already defined" bug)
func (item VariableScopeItem) MangledName() string {
	return fmt.Sprintf("%s_%d", item.name, item.varIndex)
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

var varIndex = 0

// NewVariableScopeItem constructs a function scope item
func NewVariableScopeItem(name string, value value.Value, vis Visibility) VariableScopeItem {
	item := VariableScopeItem{}
	item.name = name
	item.value = value

	item.vis = vis
	item.varIndex = varIndex

	// Here we need to do something special. This is in order to fix the bug where you cannot define
	// a variable if it has already been defined in another block in the same function
	// Example of the bug:
	//      for int i := 0; i < 200; i <- i + 1 {}
	//      for int i := 0; i < 200; i <- i + 1 {}
	// LLVM would complain in the second loop because `i` has already been defined in this "function"
	// even if the scopes are different.
	// if v, is := value.(*ir.InstAlloca); is {
	// 	v.Name = fmt.Sprintf("_%s%d", item.name, varIndex)
	// }

	varIndex++
	return item
}

// ScopeType is a storage for types in the scope. They are stored seperately from variables.
type ScopeType struct {
	Type types.Type
	Name string
	Prec int
}

// NewScopeType constructs a function scope item
func NewScopeType(name string, t types.Type, prec int) *ScopeType {
	item := &ScopeType{}
	item.Name = name
	item.Type = t
	item.Prec = prec
	return item
}
