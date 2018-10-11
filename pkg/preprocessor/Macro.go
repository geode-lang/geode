package preprocessor

// Macro stores information for a single string replacement macro
type Macro struct {
	Name  string
	Args  []string
	Body  string
	state *State
}

// Expand the macro into it's resulting string
func (m *Macro) Expand(args []string) {

}
