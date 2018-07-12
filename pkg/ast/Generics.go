package ast

// GenericSymbol is a way to represent a generic in a function/struct/class definition
type GenericSymbol struct {
	Name string
}

func (g *GenericSymbol) String() string {
	return g.Name
}
