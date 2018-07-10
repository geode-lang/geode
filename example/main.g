is main
include "std:io"
include "std:mem"

func main int {
	int i := 0;
	byte* first := mem:get(3000);
	io:print("First: %p\n", first);	
	while i < 20 {
		i <- i + 1;
		byte* ptr := mem:get(300);
		io:print("%d: %p\n", i, ptr);
	}
	
	io:print("this is a nice meme. %p\n", mem:get(300));
	return 0;
}
