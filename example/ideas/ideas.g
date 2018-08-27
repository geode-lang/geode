is __ideas__

include "std:io"
include "std:encoding"

# The new class system would basically be the go interface system.
# it is a "class" of types, all having methods attached to them
# that the class, in this case Addable, contains.
class Addable {
	func add() int;
	func value() int;
}

# Data would replace the current class system, basically representing
# a higher level struct than C has.
# all methods would be promoted to the global scope at compile time
# and would require an implicit "this" pointer to the instance to be
# passed in. This would be handled at compile time.
data Foo(a, int c) {
	int a; # a has whatever value was passed in the "constructor"
	int b := c * 30;
	func add(Addable other) int -> this.value() + other.value();
	func value int -> this.a;
}


func add(Addable a, Addable b) int {
	return a.add(b);
}

