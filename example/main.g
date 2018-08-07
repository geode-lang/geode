is main

include "std:io"
include "std:str"

func main int {
	
	string test := "\x01f44c";
	io:print(test);
	
	
	# io:print("[");
	# for int i := 0; i < str:len(test); i += 1 {
	# 	io:print("%02x ", test[i]);
	# }
	# io:print("\b]\n");
	
	return 0;
}