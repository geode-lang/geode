package main

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/nickwanninger/sheppard/pkg/parser"
	"io/ioutil"
	"os"
)

func main() {
	filename := "example/foo"
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println(err)
	}

	src := string(data)

	p := parser.New(src, filename)
	perr := p.Run()

	if perr != nil {
		fmt.Println(perr)
		os.Exit(1)
	}

	fmt.Printf("\n\n\n\n\n%s", spew.Sdump(p))
}
