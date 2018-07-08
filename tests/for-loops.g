# Name = "For Loops"
# ExpectedOutput = "1057046400"
is main
include "std::io"


func main int {
	int i := 0;
	for int x := 0; x <= 255; x <- x + 1 {
		for int y := 0; y < 255; y <- y + 1 {
			i <- i + x * y;
		}
	}
	io:print("%d", i);
	return 0;
}