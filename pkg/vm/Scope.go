package vm

import "fmt"

// Scope is a way to store values in the vm in a nested key/value pair
// most of the time the nesting will only be per function call or function
// scoping as LLVM is a flat assembler with SSA stuff
type Scope struct {
	Parent *Scope
	index  int
	values map[string]Value
}

var scopeindex = 0

// NewScope returns a newly heap allocated scope
func NewScope() *Scope {
	s := &Scope{}
	s.index = scopeindex
	scopeindex++
	s.values = make(map[string]Value)
	return s
}

// SpawnChild spawns a scope child, storing the parent reference in the
// new child reference
func (s *Scope) SpawnChild() *Scope {
	child := NewScope()
	child.Parent = s
	return child
}

// Get an item in the scope's values by traveling upward in the parent tree
func (s *Scope) Get(key string) (Value, error) {
	if val, exists := s.values[key]; exists {
		return val, nil
	}
	if s.Parent == nil {
		return nil, fmt.Errorf("unable to Get scope item with name %q", key)
	}
	return s.Parent.Get(key)
}

// Set a value in the current scope
func (s *Scope) Set(key string, val Value) error {
	s.values[key] = val
	return nil
}
