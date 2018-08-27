is io

link "io.c"
include "std:mem"



# File io Functions

# the FILE class is a standin for the C stdlib.h file handle struct
# it doesn't need to contain any fields because size wil be
# handled later on in the clang phase
class FILE {}

FILE* stdout := io:fopen("/dev/stdout", "w+");
FILE* stderr := io:fopen("/dev/stderr", "w+");


# Create external linkages to the C stdlib.h files
func fopen(string path, string mode) FILE* ...
func fseek(FILE* handle, int offset, int whence) int ...
func ftell(FILE* handle) int ...
func rewind(FILE* handle) ...
func fread(string where, int size, int nmemb, FILE* handle) ...
func fwrite(string what, int size, int nmemb, FILE* handle) ...
func fclose(FILE* handle) ...
func getenv(string what) string ...
func fgetc(FILE* handle) string ...
func fflush(FILE* handle) int ...

func fputs(string str, FILE* handle) ...


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