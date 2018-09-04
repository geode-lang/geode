is __ideas__

include "io"
include "encoding"

protocol Addable {
	func add(Addable a) int;
	func value() int;
}

class Foo(a) {
	int a;
	func add(Addable other) int -> this.value() + other.value();
	func value int -> this.a;
}

func add(Addable a, Addable b) int {
	return a.add(b);
}

func main int {
	let a := Foo(12);
	let b := Foo(3);
	return add(a, b)
}