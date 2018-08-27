is main

include "std:io"

int a := 3 + 3;

func main int {
	io:fputs(io:format("%d", a), io:stderr);
	return 0;
}