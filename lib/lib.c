#include <stdio.h>
#include <stdlib.h>
#include <stdarg.h>


// Readfile just takes some path
// and returns bytes containing
// the content of the file.
char* readfile(char* path) {
	FILE *f = fopen(path, "rb");
	fseek(f, 0, SEEK_END);
	long fsize = ftell(f);
	fseek(f, 0, SEEK_SET); // same as rewind(f);
	char *string = malloc(fsize + 1);
	fread(string, fsize, 1, f);
	fclose(f);
	string[fsize] = 0;
	return string;
}

// the print function wrapper.
void print(char *fmt, ...) {
	va_list args;
	va_start(args, fmt);
	vprintf(fmt, args);
	va_end(args);
}