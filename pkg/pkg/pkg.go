// Package pkg is where the package manager is defined
package pkg

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/geode-lang/geode/pkg/arg"
	"github.com/geode-lang/geode/pkg/util"
	"github.com/geode-lang/geode/pkg/util/color"
	input "github.com/tcnksm/go-input"
)

// PackageRule is a definition of rules for a single dependency
type PackageRule struct {
	Repo       string
	CommitLock string
}

// NewPackageRule constructs a packagerule pointer
func NewPackageRule(repo string, commitlock string) *PackageRule {
	pkg := &PackageRule{
		Repo:       repo,
		CommitLock: commitlock,
	}
	return pkg
}

// PackageManagerEnv is the structural representation of the geodepkg.toml
type PackageManagerEnv struct {
	Name     string
	Repo     string
	Packages []*PackageRule
	Test     int
}

// HandleCommand pulls from the global args package and handles `geode pkg ...
func HandleCommand() {
	var err error
	_, err = util.BashCmd("git --version")
	if err != nil {
		fmt.Println(err)
		fmt.Println("Git not found. Please install git to use 'geode pkg'")
		os.Exit(1)
	}

	if *arg.PkgInit {
		Init()
		os.Exit(0)
	}

	_, err = Config()
	if err != nil {
		fmt.Printf("Missing a geodepkg.toml config.\nRun %s to get started\n", color.Green("geode pkg --init"))
		return
	}

	EditConfig(func(env *PackageManagerEnv) {
		env.Test = env.Test + 1
	})
	// pretty.Print(env)
}

// InitCommand is what is called when the user enters `geode pkg init`
func InitCommand() {
	Init()
	os.Exit(0)
}

// Init creates geode.toml if it isn't already setup
func Init() {

	if _, err := Config(); err == nil {
		fmt.Println("geodepkg.toml already exists")
		return
	}
	fmt.Println("This utility will walk you through creating a geodepkg.toml file.")

	dir, _ := os.Getwd()

	ui := &input.UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}

	var repo string

	resp, err := util.BashCmd("git config remote.origin.url")
	if err == nil {
		repo = strings.TrimSpace(resp)
	}

	name, _ := ui.Ask("Project Name", &input.Options{
		Default: filepath.Base(filepath.Dir(dir)),
	})
	repo, _ = ui.Ask("Git Repo", &input.Options{
		Default: repo,
	})

	env := &PackageManagerEnv{}

	env.Name = name
	env.Repo = repo

	env.Packages = make([]*PackageRule, 0)
	WriteConfig(env)
	fmt.Println("Created geodepkg.toml")
}

// Config reads the config from the current directory's geodepkg.toml
func Config() (*PackageManagerEnv, error) {
	env := &PackageManagerEnv{}
	_, err := toml.DecodeFile("geodepkg.toml", env)
	if err != nil {
		return nil, err
	}
	return env, nil
}

// EditConfig takes an editor function and runs an edit on the config and
// re-writes it to the disk
func EditConfig(editor func(*PackageManagerEnv)) {
	env, _ := Config()
	editor(env)
	WriteConfig(env)
}

// WriteConfig takes a config and writes it to the toml file
func WriteConfig(env *PackageManagerEnv) {
	buff := &bytes.Buffer{}
	toml.NewEncoder(buff).Encode(env)
	ioutil.WriteFile("geodepkg.toml", buff.Bytes(), os.ModePerm)
}
