# Name = "FizzBuzz"
# ExpectedOutput = "FizzBuzz 1 2 Fizz 4 Buzz "
include "std::io"

func main int {
	for int i := 0; i <= 5; i <- i + 1 {
		int written := 0;
		if i % 3 = 0 {
			written <- 1;
			print("Fizz");
		}
		if i % 5 = 0 {
			written <- 1;
			print("Buzz");
		}
		if written = 0 {
			print("%d ", i);
		} else {
			print(" ");	
		}
	}
	return 0;
}