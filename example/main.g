is main
include "std:io"
include "std:mem"


func main(int argc, string* argv) int {
	int* a := &argc;
	# io:print("Argc: %d\nStack Pointer: %p\n", *a, &_stkptr);
	
	for int i := 0; i < argc; i <- i + 1 {
		io:print("%d: %s\n", i, argv[i]);
	}
	return 0;
}
 