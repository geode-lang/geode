# global variables 3
is main

include "std:io"

int a := test();

func test int {
	io:print("called");
	return 3;
}

func main int {
	# Test code here.
	return a;
}
