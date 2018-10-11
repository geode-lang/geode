package preprocessor

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// State -
type State struct {
	Reader *bufio.Reader
	Result io.Writer
	Macros map[string]*Macro
}

// New creates a preprocessor State
func New() *State {
	pp := &State{}
	return pp
}

// Run the preprocessor on a file
func (pp *State) Run(source string) (string, error) {
	val := pp.Process(source)
	return val, nil
}

// Process a source string
func (pp *State) Process(source string) string {

	w := newWalker(source)

	for true {
		ln := w.next()
		ln = strings.TrimLeft(ln, " \t")
		if w.eof() {
			break
		}
		if len(ln) == 0 {
			continue
		}

		if ln[0] == '@' {
			w.drop()
			continue
		}
		fmt.Println(ln)
	}

	return strings.Join(w.lines, "\n")
}

func (pp *State) parsemacro(w *walker) *Macro {

	return nil
}

// NewMacro creates a new macro and adds it to the state
func (pp *State) NewMacro(name string, args []string, body string) *Macro {
	m := &Macro{}
	m.Name = name
	m.Args = args
	m.Body = body
	m.state = pp
	pp.Macros[name] = m
	return m
}

type walker struct {
	lines []string
	index int
}

func newWalker(s string) *walker {
	w := &walker{}
	w.lines = strings.Split(s, "\n")
	w.index = -1
	return w
}

func (w *walker) next() string {
	w.index++
	if w.index < len(w.lines) {
		return w.lines[w.index]
	}
	return ""
}

// drop the currently indexed line
func (w *walker) drop() string {
	ln := w.lines[w.index]
	w.lines = append(w.lines[:w.index], w.lines[w.index+1:]...)
	w.index--
	return ln
}

func (w *walker) eof() bool {
	return w.index >= len(w.lines)
}
