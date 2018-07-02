#include <stdio.h>
#include <stdarg.h>

// the print function wrapper.
void print(char *fmt, ...) {
	va_list args;
	va_start(args, fmt);
	vprintf(fmt, args);
	va_end(args);
}