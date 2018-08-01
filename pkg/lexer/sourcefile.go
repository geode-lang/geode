package lexer

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/geode-lang/geode/pkg/util/log"
)

// Sourcefile is a wrapper around a rune array
// from a file, along with some loading functions
type Sourcefile struct {
	Path     string
	Name     string
	contents []rune
}

// NewSourcefile resolves a filename and creates a source file
func NewSourcefile(name string) (*Sourcefile, error) {
	s := &Sourcefile{}
	s.Name = name
	return s, nil
}

// Hash - Get the hash of the sourcefile. This is a truncated sha1.
func (s *Sourcefile) Hash() []byte {
	h := sha1.New()
	io.WriteString(h, s.String())
	return h.Sum(nil)
}

// HashName returns the name and the hash.
func (s *Sourcefile) HashName() string {
	return fmt.Sprintf("%s_%x", s.Name, s.Hash()[:2])
}

// LoadFile -
func (s *Sourcefile) LoadFile(src string) error {
	s.Path = src
	bytes, err := ioutil.ReadFile(src)
	if err != nil {
		log.Fatal("Unable to read file at path '%s'\n", src)
	}
	s.Name = src
	log.Debug("Reading %s\n", src)
	s.LoadBytes(bytes)
	return nil
}

// ResolveFile resolves a filename and loads it
func (s *Sourcefile) ResolveFile(path string) error {
	log.Debug("Resolving filename %q\n", path)
	p, e := ResolveFileName(path, ".g")
	if e != nil {
		log.Fatal("Unable to resolve path '%s'\n", path)
	}
	s.Name = p
	return s.LoadFile(p)
}

// LoadString takes a string and loads it
func (s *Sourcefile) LoadString(source string) {
	runes := []rune(source)
	s.contents = runes
}

// LoadBytes takes an array of bytes and loads it into the source
func (s *Sourcefile) LoadBytes(bytes []byte) {
	s.LoadString(string(bytes))
}

func (s *Sourcefile) String() string {
	return string(s.contents)
}

// Bytes returns the source as a byte array
func (s *Sourcefile) Bytes() []byte {
	return []byte(string(s.contents))
}

// ResolveFileName - if the filename passed in is a folder, look in that folder for a main.g
// if the filename is not, look for a file matching that filename, but with a .g extension
func ResolveFileName(filename, suffix string) (string, error) {
	// Grab the stats of the file
	stats, err := os.Stat(filename)

	// If there was an error (file doesnt exist)
	if err != nil {
		// Try resolving the filename with .g extension
		if !strings.HasSuffix(filename, suffix) {
			return ResolveFileName(filename+suffix, suffix)
		}
		// There was no file by that name, so we fail
		return "", fmt.Errorf("fatal error: No such file or directory %s", filename)
	}
	if stats.IsDir() {
		log.Debug("Resolving: Looking in %s for main.g\n", filename)
		return ResolveFileName(filename+"/main.g", suffix)
	}

	return filename, nil
}
