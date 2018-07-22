#include <stdio.h>
#include <stdarg.h>
#include <stdlib.h>

#include "../c/_runtime.h"

char*
gcmalloc(long size) {
	return tgc_alloc(&_G_GC, size);
}

void
rawfree(char* ptr) {
	free((void*)ptr);
}