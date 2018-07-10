is main
include "std:io"
include "std:mem"

func main int {
	float pi := 3.1415926535897932384626433832795028841971693993751058209749445923078164;
	string data := io:format("%.50g", pi);
	io:print("%s\n", data);
	return 0;
}
