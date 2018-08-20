package ast

// ParseResult is what is returned from a file parsing.
type ParseResult struct {
	Pkg   string
	Nodes []Node
	Paths []string
}

// FlattenParseResults takes an array of ParseResults and flattens their PkgNames together
func FlattenParseResults(results []ParseResult) (finals []ParseResult) {

	for _, r := range results {
		found := false
		for _, f := range finals {
			if f.Pkg == r.Pkg {
				f.Nodes = append(f.Nodes, r.Nodes...)
				f.Paths = append(f.Paths, r.Paths...)
				found = true
				break
			}
		}

		if found == false {
			finals = append(finals, r)
		}

	}
	return results
}
