

#include "_runtime.h"

#define USE_GC


void
___geodegcinit(char* stack_pointer) {
	tgc_start(&_G_GC, (void*)stack_pointer);
}
