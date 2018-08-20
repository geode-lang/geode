# Name = "FizzBuzz"
# ExpectedOutput = "FizzBuzz 1 2 Fizz 4 Buzz "
is main
include "std:io"


func main int {
	for int i := 0; i <= 5; i <- i + 1 {
		int written := 0;
		if i % 3 = 0 {
			written <- 1;
			io:print("Fizz");
		}
		if i % 5 = 0 {
			written <- 1;
			io:print("Buzz");
		}
		if written = 0 {
			io:print("%d ", i);
		} else {
			io:print(" ");	
		}
	}
	return 0;
}