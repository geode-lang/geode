package ast

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

// =========================== FunctionDiscoveryWorker ===========================

// FunctionDiscoveryWorker is a struct that
// stores state for function discovery
type FunctionDiscoveryWorker struct {
	Prog *Program
}

// NewFunctionDiscoveryWorker returns an instance of a discovery worker attached to a program
func NewFunctionDiscoveryWorker(prog *Program) *FunctionDiscoveryWorker {
	w := &FunctionDiscoveryWorker{}
	w.Prog = prog
	return w
}

// Discover a function from a worker
// This function can possibly return a newly generated function from a variant
// ie: for generics, unknown type arguments or overloading
func (w *FunctionDiscoveryWorker) Discover(opt *FunctionDiscoveryOptions) (res *FunctionDiscoveryResult, err error) {

	// searchpaths := opt.GetSearchpaths()

	// fmt.Println(searchpaths)

	return res, err
}

// =========================== FunctionDiscoveryOptions ===========================

// FunctionDiscoveryOptions are the options passed into the function discovery
type FunctionDiscoveryOptions struct {
	searchName   string
	searchingPkg *Package
	args         []types.Type
}

// NewFunctionDiscoveryOptions returns a FunctionDiscoveryOptions instance with
// the name and the source package
func NewFunctionDiscoveryOptions(searchName string, searchingPkg *Package) *FunctionDiscoveryOptions {
	opt := &FunctionDiscoveryOptions{}
	opt.searchName = searchName
	opt.searchingPkg = searchingPkg
	return opt
}

// GetSearchpaths returns the searchpaths for the options
func (opt *FunctionDiscoveryOptions) GetSearchpaths() []string {
	searchPaths := make([]string, 0)
	searchPaths = append(searchPaths, opt.searchName)
	searchPaths = append(searchPaths, fmt.Sprintf("%s:%s", opt.searchingPkg.Name, opt.searchName))

	return searchPaths
}

// AddArgs appends variadic arguments to the discovery options
func (opt *FunctionDiscoveryOptions) AddArgs(args ...types.Type) {
	opt.args = append(opt.args, args...)
}

// =========================== FunctionDiscoveryResult ===========================

// FunctionDiscoveryResult is what is returned from
//    FunctionDiscoveryWorker.Discover()
type FunctionDiscoveryResult struct {
	variadic bool
	name     string
	pkg      *Package
	prog     *Program
	Func     *ir.Func
}
