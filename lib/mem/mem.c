#include <stdio.h>
#include <stdarg.h>
#include <stdlib.h>
 
#include "../include/gc/gc.h"
// #include "../include/_runtime.h"


// extern long __memoryused;

long used() {
	return 0;
}

int malloc_count = 0;
int current_mallocs = 0;

void
GC_FREE_HANDLE(void* ptr) {
	current_mallocs--;
}


void
rawfree(char* ptr) {
	free((void*)ptr);
}