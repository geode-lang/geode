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
___geodegcinit(void* stk) {
	#ifdef USE_GC
	// Initialize the garbage collector using argc as the base of the stack
	// This is so the GC can find where to look in it's sweeps
	tgc_start(&gc, &stk);
	#endif
}



char*
rawmalloc(int size) {
	#ifdef USE_GC
	return tgc_alloc(&gc, size);
	#else
	return (char*)malloc(size);
	#endif
}
