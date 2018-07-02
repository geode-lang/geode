# Name = "link"
# ExpectedOutput = "42"

include "std::io"
link "link.c"

func foo int ...

func main(int argc) int {
	print("%d", foo());
	return 0;
}
