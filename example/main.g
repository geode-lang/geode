is main

include "std:io"
include "std:str"
include "std:mem"

func main(int argc, string* argv) int {
	# argv[0] <- "program";
	# byte v := 65;
	# argv[0][2] <- v;
	
	# for int i := 0; i < argc; i <- i + 1 {
	# 	io:print("strlen('%s') -> %d\n", argv[i], strlen(argv[i]));
	# }
	
	
	string msg := "for int i := 0; i < argc; i <- i + 1 {";
	int length := str:len(msg);
	
	
	io:print("%d\n", length);
	
	int i := 0;
	while 1 {
		
		length <- str:len(msg);
		int index := i % length;
		msg[index] <- (i % 65) + 65;
		
		io:print("%s\n", msg);
		i <- i + 1;
	}
	return 0;
}

