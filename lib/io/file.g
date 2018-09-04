is io

link "io.c"
include "mem"



# File io Functions

# the File class is a standin for the C stdlib.h file handle struct
# it doesn't need to contain any fields because size wil be
# handled later on in the clang phase
class File {}

File* stdout := io:fopen("/dev/stdout", "w+");
File* stderr := io:fopen("/dev/stderr", "w+");
File* stdin := io:fopen("/dev/stdin", "r+");


# Create external linkages to the C stdlib.h files
func fopen(string path, string mode) File* ...
func fseek(File* handle, int offset, int whence) int ...
func ftell(File* handle) int ...
func rewind(File* handle) ...
func fread(string where, int size, int nmemb, File* handle) ...
func fwrite(string what, int size, int nmemb, File* handle) ...
func fclose(File* handle) ...
func getenv(string what) string ...
func fgetc(File* handle) string ...
func fflush(File* handle) int ...
func fprintf(File* handle, byte* format, ...) int ...

func fputs(string str, File* handle) ...

# read_file returns the full content of a file
func read_file(string path) string {
	File* f := fopen(path, "r");
	
	fseek(f, 0, 2);
	int fsize := ftell(f);
	rewind(f);
	
	string data := mem:get(fsize + 1);
	fread(data, fsize, 1, f);
	# ensure the string is null terminated
	data[fsize] <- 0;
	fclose(f);
	return data;
}