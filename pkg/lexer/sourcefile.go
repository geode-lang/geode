package lexer

import (
	"crypto/sha1"
	"io"
	"io/ioutil"

	"github.com/nickwanninger/geode/pkg/util/log"
)

// Sourcefile is a
type Sourcefile struct {
	Path     string
	Name     string
	Contents []rune
}

// NewSourcefile resolves a filename and creates a source file
func NewSourcefile(name string) (*Sourcefile, error) {
	s := &Sourcefile{}
	s.Name = name
	return s, nil
}

// Hash - Get the has of the sourcefile
func (s *Sourcefile) Hash() []byte {
	h := sha1.New()
	io.WriteString(h, s.String())
	return h.Sum(nil)
}

// LoadFile -
func (s *Sourcefile) LoadFile(path string) error {
	s.Path = path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Unable to read file at path '%s'\n", path)
	}
	s.LoadBytes(bytes)
	return nil
}

// LoadString takes a string and loads it
func (s *Sourcefile) LoadString(source string) {
	runes := []rune(source)
	s.Contents = runes
}

// LoadBytes takes an array of bytes and loads it into the source
func (s *Sourcefile) LoadBytes(bytes []byte) {
	s.LoadString(string(bytes))
}

func (s *Sourcefile) String() string {
	return string(s.Contents)
}

// Bytes returns the source as a byte array
func (s *Sourcefile) Bytes() []byte {
	return []byte(string(s.Contents))
}
