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
	int i := 0;
	while i < 255 {
		i + 1;
		io:print("%d\n", i);
	}
	return 0;
}
