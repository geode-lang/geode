is io

link "io.c"

include "std:mem"
include "std:str"



# Printing bindings
func print(string format, ...) ...
func format(string format, ...) string ...

func println(string message) {
	print("%s\n", message);
}


# File io Functions

# FILE is a standin for the C stdlib.h file handle struct
# it doesn't need to contain any fields because size wil be
# handled later on in the clang phase
class FILE {}

# Create external linkages to the C stdlib.h files
func fopen(string path, string mode) FILE* ...
func fseek(FILE* handle, int offset, int whence) int ...
func ftell(FILE* handle) int ...
func rewind(FILE* handle) ...
func fread(string where, int size, int nmemb, FILE* handle) ...
func fclose(FILE* handle) ...
func getenv(string what) string ...
func fgetc(FILE* handle) string ...


# readFile returns the full content of a file
func readFile(string path) string {
	FILE* f := fopen(path, "r");
	
	fseek(f, 0, 2);
	int fsize := ftell(f);
	rewind(f);
	
	string data := mem:get(fsize + 1);
	fread(data, fsize, 1, f);
	# ensure the string is null terminated
	data[fsize] <- 0;
	return data;
}