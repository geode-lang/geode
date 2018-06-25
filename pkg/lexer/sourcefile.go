package lexer

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"

	"gitlab.com/nickwanninger/geode/pkg/util/log"

	"github.com/timtadh/data-structures/errors"
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
func (s *Sourcefile) Hash() string {
	hasher := sha1.New()
	hasher.Write([]byte(string(s.Contents)))
	sha := fmt.Sprintf("%x", hasher.Sum(nil))
	return sha
}

func (s *Sourcefile) Read(path string) error {
	s.Path = path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Unable to read file at path '%s'\n", path)
	}
	runes := []rune(string(bytes))
	s.Contents = runes
	// for _, rn := range runes {
	// 	fmt.Println(rn, strconv.QuoteRune(rn))
	// }
	return errors.Errorf("Unable to read file")
}
