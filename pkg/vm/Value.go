package vm

import "fmt"

// Value is an interface used to represent a value in the
// virtual machine
type Value interface {
	fmt.Stringer
}
