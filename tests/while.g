# Name = "while"
# ExpectedOutput = "255"
include "std::io"

func main int {
	int i := 0;
	while i < 255 {
		i <- i + 1;
	}
	print("%d", i);
	return 0;
}