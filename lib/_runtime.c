#include <stdio.h>

#define USE_GC

#ifdef USE_GC
#include "./gc/tgc.h"
static tgc_t gc;
#else
#include <stdlib.h>
#endif




extern long __GEODE__main(int argc, char** argv);

int
main(int argc, char** argv) {
	#ifdef USE_GC
	// Initialize the garbage collector using argc as the base of the stack
	// This is so the GC can find where to look in it's sweeps
	tgc_start(&gc, &argc);
	#endif
	
  int prog_ret_val = __GEODE__main(argc, argv);

	return prog_ret_val;
}


char*
__GEODE__alloca(int size) {
	#ifdef USE_GC
	// tgc_run(&gc);
	return tgc_alloc(&gc, size);
	#else
	return (char*)malloc(size);
	#endif
}

void
__GEODE_free(char* ptr) {
	#ifdef USE_GC
	// tgc_run(&gc);
	return tgc_free(&gc, ptr);
	#else
	return free(ptr);
	#endif
}