package info

import (
	"encoding/json"
	"fmt"
)

// Item is an interface that has methods used to display information
type Item interface {
}

type context struct {
	tokens []Item
	nodes  []Item
}

// global info context
var gic *context

func init() {
	gic = &context{}
}

// AddToken adds a token to the info context
func AddToken(t Item) {
	gic.tokens = append(gic.tokens, t)
}

// AddNode adds a node to the info context
func AddNode(n Item) {
	gic.nodes = append(gic.nodes, n)
}

// Dump info to the console
func Dump() {
	fmt.Printf("Info:\n")
	fmt.Printf("Token Count: %d\n", len(gic.tokens))
}

// DumpJSON dumps the context to the console
func DumpJSON() {
	data := make(map[string]interface{})

	// data["tokens"] = gic.tokens

	data["nodes"] = gic.nodes

	j, _ := json.MarshalIndent(data, "", "    ")
	fmt.Printf("%s", string(j))
}
