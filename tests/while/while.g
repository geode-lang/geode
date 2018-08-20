# Name = "while"
# ExpectedOutput = "255"
is main
include "std:io"

func main int {
	int i := 0;
	while i < 255 {
		i <- i + 1;
	}
	io:print("%d", i);
	return 0;
}