is io

link "io.c"
include "mem"



# FILE_DESCRIPTOR io Functions

# the FILE_DESCRIPTOR class is a standin for the C stdlib.h file handle struct
# it doesn't need to contain any fields because size wil be
# handled later on in the clang phase
class FILE_DESCRIPTOR {}


# Some default file descriptors that are used pretty commonly
FILE_DESCRIPTOR* stdout = fopen("/dev/stdout", "w+");
FILE_DESCRIPTOR* stderr = fopen("/dev/stderr", "w+");
FILE_DESCRIPTOR* stdin = fopen("/dev/stdin", "r+");



func fgets(byte* buf, int len, FILE_DESCRIPTOR* fd) ...



# File is a class wrapper around the FILE_DESCRIPTOR that has more methods
# that allow you to act on `this` instead of passing around a FILE_DESCRIPTOR
class File {

	# handle is a pointer to an OS Specific file descriptor
	# returned by the c bindings for fopen
	FILE_DESCRIPTOR* handle

	# write some string to the file handle
	func puts(string msg) {
		io:fputs(msg, this.handle)
		this.flush()
	}

	# wrap around the stdlib io:fflush
	func flush {
		io:fflush(this.handle)
	}

	# wrap around the stdlib io:ftell
	func pos long {
		return io:ftell(this.handle)
	}

	# wrap around the stdlib io:fseek
	func seek(int off, int loc) long {
		return io:fseek(this.handle, off, loc)
	}

	# wrap around the stdlib io:fclose
	func close {
		io:fclose(this.handle)
	}

	# wrap around the stdlib io:fread
	func read(string where, int size, int nmemb) {
		io:fread(where, size, nmemb, this.handle)
	}

	# Return the filesize in bytes
	func size long {
		this.seek(0, 2);
		fsize = this.pos();
		io:rewind(this.handle)
		return fsize
	}

	func readall byte* {
		# get the size of the file
		fsize = this.size()
		# allocate enough space for the buffer
		buffer = mem:get(fsize + 1);
		# read the file into the buffer entirely
		res = io:fread(buffer, fsize, 1, this.handle);
		# ensure the string is null terminated
		buffer[fsize] = 0;
		# return the buffer we read
		return buffer
	}
}


func open(string path, string mode) File* {
	File* f = mem:get(info(File).size)
	f.handle = fopen(path, mode)
	return f
}



# Create external linkages to the C stdlib.h files
func fopen(string path, string mode) FILE_DESCRIPTOR* ...
func fseek(FILE_DESCRIPTOR* handle, int offset, int whence) int ...
func ftell(FILE_DESCRIPTOR* handle) long ...
func rewind(FILE_DESCRIPTOR* handle) ...
func fread(string where, int size, int nmemb, FILE_DESCRIPTOR* handle) long ...
func fwrite(string what, int size, int nmemb, FILE_DESCRIPTOR* handle) ...
func fclose(FILE_DESCRIPTOR* handle) ...
func getenv(string what) string ...
func fgetc(FILE_DESCRIPTOR* handle) string ...
func fflush(FILE_DESCRIPTOR* handle) int ...
func fprintf(FILE_DESCRIPTOR* handle, byte* format, ...) int ...
func ferror(FILE_DESCRIPTOR* handle) int ...
func feof(FILE_DESCRIPTOR* handle) int ...
func fputs(string str, FILE_DESCRIPTOR* handle) ...

# read_file returns the full content of a file
func read_file(string path) string {
	FILE_DESCRIPTOR* f = fopen(path, "r");
	
	fseek(f, 0, 2);
	int fsize = ftell(f);
	rewind(f);
	
	string data = mem:get(fsize + 1);
	fread(data, fsize, 1, f);
	# ensure the string is null terminated
	data[fsize] = 0;
	fclose(f);
	return data;
}
