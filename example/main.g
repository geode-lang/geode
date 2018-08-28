is main

include "std:io"
include "std:mem"

func foo() byte* {
	let data := mem:get(10);
	
	for int i := 0; i < 10; i += 1 {
		data[i] <- i;
	}
	return data;
}

func main int {	
	let i := 1000;
	while true {
		io:print("%p\n", foo());
	}
	return 0;
}