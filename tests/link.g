# Name = "link"
# ExpectedOutput = "42"
is main

include "std::io"
link "link.c"

func foo int ...

func main(int argc) int {
	io:print("%d", foo());
	return 0;
}
