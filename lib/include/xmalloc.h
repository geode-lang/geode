#ifndef __xmalloc__
#define __xmalloc__

#include <stdlib.h>

#include "./gc/gc.h"

typedef struct {
	long size;
	int alloc_count;
} xmalloc_prelude_t;


typedef struct {
	unsigned long memoryused;
	unsigned int blocksallocated;
	char readable[20];
} xmalloc_stat_t;


#define PRELUDE_SIZE (sizeof(xmalloc_prelude_t))





xmalloc_stat_t xmemstat();
long xmalloc_size(void* ptr);
void xfree(void* ptr);
void* xmalloc(size_t size);
void* xcalloc(unsigned count, unsigned size);
void* xrealloc(void* ptr, size_t newsize);


#endif