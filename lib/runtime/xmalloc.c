#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>

#include "../include/xmalloc.h"

// #define DEBUG_XMALLOC

// eXtended memory things, like xmalloc, xcalloc, and xfree
// These functions give a little bit more information in a prelude

long __memoryused = 0;
int __blocksallocated = 0;

static xmalloc_prelude_t*
xmalloc_getprelude(void* ptr) {
	return (xmalloc_prelude_t*)(ptr - PRELUDE_SIZE);
}

long
xmalloc_size(void* ptr) {
	xmalloc_prelude_t* prelude = xmalloc_getprelude(ptr);
	return prelude->size;
}


void
xfree(void* ptr) {

	// Don't free a null pointer
	if (ptr == NULL) {
		return;
	}

	void* new_ptr = ptr - PRELUDE_SIZE;

	xmalloc_prelude_t* prelude = xmalloc_getprelude(ptr);
	__memoryused -= prelude->size;
	__blocksallocated--;
	free(new_ptr);
	#ifdef DEBUG_XMALLOC
	printf("[DEBUG] Freed %zu bytes\n", prelude->size);
	#endif
}


void*
xmalloc(size_t size) {
	
	void* realptr = (void*) malloc(size + PRELUDE_SIZE);
	
	if (realptr == NULL) {
		fprintf(stderr, "Fatal: memory exhausted (xmalloc of %zu bytes).\n", size);
		exit(EXIT_FAILURE);
	}
	__memoryused += size;
	__blocksallocated++;
	#ifdef DEBUG_XMALLOC
	printf("[DEBUG] Allocated %u bytes\n", size);
	#endif

	xmalloc_prelude_t* prelude = realptr;
	prelude->size = size;
	prelude->alloc_count = 1;
	return (void*)(realptr + PRELUDE_SIZE);
}


void*
xrealloc(void* ptr, size_t newsize) {
	// Give them a new block of memory if
	// there isnt anything to reallocate
	if (ptr == NULL) return xmalloc(newsize);
	
	// The real pointer is offset by PRELUDE_SIZE
	void* real_ptr = ptr - PRELUDE_SIZE;
	
	// Pull the prelude data out of the pointer
	xmalloc_prelude_t* prelude = xmalloc_getprelude(ptr);
	size_t oldsize = prelude->size;

	void* newptr = realloc(real_ptr, newsize + PRELUDE_SIZE);
	if (newptr == NULL) {
		fprintf(stderr, "Fatal: Memory reallocation of %p to %zu bytes from %zu bytes failed.\n", ptr, newsize, oldsize);
		exit(EXIT_FAILURE);
	}
	xmalloc_prelude_t* new_prelude = newptr;
	new_prelude->size = newsize;
	new_prelude->alloc_count++;

	// Update the "memory used" value
	__memoryused += newsize - oldsize;
	return (void*)(newptr + PRELUDE_SIZE);
}


void*
xcalloc(unsigned count, unsigned size) {
	unsigned int n = count * size;
	// Errors should be handled in the xmalloc function
	void* new_mem = (void*) xmalloc(n);
	memset(new_mem, '\0', n);
	return new_mem;
}


