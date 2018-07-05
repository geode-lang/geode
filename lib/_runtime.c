#include <stdio.h>

#define USE_GC

#ifdef USE_GC
#include "./tgc/tgc.h"
static tgc_t gc;
#else
#include <stdlib.h>
#endif



void
__GEODE__INIT__GC__(void* stk) {
	#ifdef USE_GC
	// Initialize the garbage collector using argc as the base of the stack
	// This is so the GC can find where to look in it's sweeps
	tgc_start(&gc, &stk);
	#endif
}



extern long __GEODE__main(int argc, char** argv);

int
main(int argc, char** argv) {
	__GEODE__INIT__GC__(&argc);
	
  int prog_ret_val = __GEODE__main(argc, argv);

	return prog_ret_val;
}



char*
__GEODE__alloca(int size) {
	#ifdef USE_GC
	return tgc_alloc(&gc, size);
	#else
	return (char*)malloc(size);
	#endif
}

void
__GEODE_free(char* ptr) {
	#ifdef USE_GC
	return tgc_free(&gc, ptr);
	#else
	return free(ptr);
	#endif
}