is main

include "std:io"

func main(int argc, string* argv) int {
	argv[0] <- "program";
	for int i := 0; i < argc; i <- i + 1 {
		io:print("strlen('%s') -> %d\n", argv[i], strlen(argv[i]));
	}
	return 0;
}

func strlen(string str) int {
	int len;
	for len <- 0; str[len] != 0; len <- len + 1 {}
	return len;
}