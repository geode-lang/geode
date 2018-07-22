package ast

// NamedReference is a reference to an item on the scope via some name
// this can be used to access variable allocs, function defns, or types
type NamedReference struct {
	Value string
}

// NewNamedReference returns a new name reference with a string as it's name
func NewNamedReference(name string) *NamedReference {
	n := &NamedReference{}

	n.Value = name
	return n
}

func (n *NamedReference) String() string {
	return n.Value
}
