is main

include "std:io"
include "std:encoding"

func main int {
	float a := 1;
	while true {
		io:print("%8f: %s\n", a, encoding:base64(a));
		a += (a / 1000.0);
	}
	return 0;
}
