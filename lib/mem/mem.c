#include <stdio.h>
#include <stdarg.h>
#include <stdlib.h>

#include "../c/_runtime.h"

int malloc_count = 0;
int current_mallocs = 0;

void
GC_FREE_HANDLE(void* ptr) {
	current_mallocs--;
}

char*
gcmalloc(long size) {
	void* ptr = tgc_alloc(&_G_GC, size);
	tgc_set_dtor(&_G_GC, ptr, GC_FREE_HANDLE);
	current_mallocs++;
	malloc_count++;
	return ptr;
}

void
rawfree(char* ptr) {
	free((void*)ptr);
}