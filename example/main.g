is main

include "std:io"

func foo int -> 42;

func main int {
	string* a := ["hello", 1, "how", "are", "you"];
	io:print("%s\n", a[2]);
	return 0;
}

