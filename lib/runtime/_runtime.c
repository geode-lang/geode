#include "../include/_runtime.h"

#include "stdio.h"

char*
raw_copy(char* source, int length) {
	int len = strlen(source);
	char* dest = GC_MALLOC(len + 1);
	memcpy(dest, source, length);
	return dest;
}
