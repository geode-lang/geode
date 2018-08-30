is main

include "std:io"

include "std:mem"


func foo int -> 30

func greet -> io:print("hello\n")

func main(int argc, byte** argv) int {
	greet()
	return 1;
}
