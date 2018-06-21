package lexer

type Sourcefile struct {
	Path     string
	Name     string
	Contents []rune
}

// NewSourcefile resolves a filename and creates a source file
func NewSourcefile(name string) (*Sourcefile, error) {
	source := &Sourcefile{}

	return source, nil
}
