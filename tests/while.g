# Name = "while"
# ExpectedOutput = "255"

func main int {
	int i := 0;
	while i < 255 {
		i <- i + 1;
	}
	print("%d", i);
	return 0;
}