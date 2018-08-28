is main

include "std:io"
include "std:mem"

func foo() int {
	let data := mem:get(10);
	
	for int i := 0; i < 10; i += 1 {
		data[i] <- i;
	}
	return data[4];
}

func main int {	
	let i := 1000;
	while true {
		io:print("%d\n", foo());
	}
	return 0;
}