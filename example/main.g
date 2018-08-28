is main

include "std:io"
include "std:mem"
include "std:math"


func test {
	io:print("%p. Heap Size: %d\n", mem:get(math:rand() % 1000), mem:GC_get_heap_size());
}

func main int {	
	let i := 0;
	while i <= 500 {
		test();
		i+=1;
	}
	io:print("final\n");
	io:print("%p. Heap Size: %d\n", mem:get(30), mem:GC_get_heap_size());
	return 0;
}