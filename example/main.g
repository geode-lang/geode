is main
include "std:io"
include "std:mem"


class Foo {
	
	int value;
	
	func foo() int {
		return this.value;
	}
}


func foo {
	string data := io:format("%d", 300);
	io:print("%s\n", data);
}

func main(int argc, string argv) int {
	while 1 {
		foo();
	}
	return 0;
}
