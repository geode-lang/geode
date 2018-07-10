#include <stdio.h>
#include <stdarg.h>

#define USE_GC

#ifdef USE_GC
#include "./tgc/tgc.h"
static tgc_t gc;
#else
#include <stdlib.h>
#endif



void
___geodegcinit(char* stack_pointer) {
	tgc_start(&gc, (void*)stack_pointer);
}



char*
gcmalloc(long size) {
	#ifdef USE_GC
	return tgc_alloc(&gc, size);
	#else
	return (char*)malloc(size);
	#endif
}


char*
rawmalloc(long size) {
	return (char*)malloc(size);
}


void
rawfree(char* ptr) {
	free((void*)ptr);
}