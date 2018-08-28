is primes

include "std:io"
func main int {
	let i := 0;
	while true {
		if check_prime(i) {
			io:print("%d\n", i);
		}
		i+=1;
	}
	return 0;
}


func check_prime(int a) bool {
	if a % 2 = 0 {
		return false;
	}
	int c;
  for c <- 2; c <= a - 1; c += 1 { 
		if a % c = 0 {
			return false;
		}
	}
	return true;
}