#include <math.h>
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "../include/xmalloc.h"
#include <gc/gc.h>

// #define DEBUG_XMALLOC

static pthread_mutex_t mutex = PTHREAD_MUTEX_INITIALIZER;

static void xmalloc_lock() { pthread_mutex_lock(&mutex); }

static void xmalloc_unlock() { pthread_mutex_unlock(&mutex); }

// eXtended memory things, like xmalloc, xcalloc, and xfree
// These functions give a little bit more information in a prelude

static long allocated_before_collect = 0;

long __memoryused = 0;
int __blocksallocated = 0;
long __alloc_index = 0;

long bytes_used() { return __memoryused; }
long blocks_used() { return __blocksallocated; }

static xmalloc_prelude_t *xmalloc_getprelude(void *ptr) {
  return (xmalloc_prelude_t *)(ptr - PRELUDE_SIZE);
}

long xmalloc_size(void *ptr) {
  xmalloc_lock();
  xmalloc_prelude_t *prelude = xmalloc_getprelude(ptr);
  xmalloc_unlock();
  return prelude->size;
}

long heap_size() { return allocated_before_collect; }

void xfree(void *ptr) {
  xmalloc_lock();
  // Don't free a null pointer
  if (ptr == NULL) {
    return;
  }

  void *new_ptr = ptr - PRELUDE_SIZE;

  xmalloc_prelude_t *prelude = xmalloc_getprelude(ptr);
  __memoryused -= prelude->size;

  __blocksallocated--;
  GC_FREE(new_ptr);
#ifdef DEBUG_XMALLOC
  printf("[DEBUG] xfree(%p) -> %u bytes\n", ptr, prelude->size);
#endif

  xmalloc_unlock();
}

static void xfinalizer(GC_PTR obj, GC_PTR x) {
  xmalloc_prelude_t *prelude = (xmalloc_prelude_t *)obj;
#ifdef DEBUG_XMALLOC
  printf("[DEBUG] gc_xfree(%p) -> %u bytes\n", obj, prelude->size);
#endif
  allocated_before_collect -= prelude->size;
}

void *xmalloc(size_t size) {

  void *realptr = (void *)GC_MALLOC(size + PRELUDE_SIZE);
  GC_register_finalizer(realptr, xfinalizer, 0, 0, 0);

  xmalloc_lock();
  // GC_gcollect();
  allocated_before_collect += size;
  if (realptr == NULL) {
    fprintf(stderr, "Fatal: memory exhausted (xmalloc of %zu bytes).\n", size);
    exit(EXIT_FAILURE);
  }

  __memoryused += size;
  __blocksallocated++;
#ifdef DEBUG_XMALLOC
  printf("[DEBUG] xmalloc(%u) -> %p\n", size, realptr);
#endif

  xmalloc_prelude_t *prelude = realptr;
  prelude->size = size;
  prelude->alloc_count = 1;
  prelude->alloc_index = __alloc_index;

  __alloc_index++;

  xmalloc_unlock();

  return (void *)(realptr + PRELUDE_SIZE);
}

void *xrealloc(void *ptr, size_t newsize) {
  // Give them a new block of memory if

  // there isnt anything to reallocate
  if (ptr == NULL)
    return xmalloc(newsize);

  // The real pointer is offset by PRELUDE_SIZE
  void *real_ptr = ptr - PRELUDE_SIZE;

  // Pull the prelude data out of the pointer
  xmalloc_prelude_t *prelude = xmalloc_getprelude(ptr);
  size_t oldsize = prelude->size;

  void *newptr = GC_REALLOC(real_ptr, newsize + PRELUDE_SIZE);
  if (newptr == NULL) {
    fprintf(stderr,
            "Fatal: Memory reallocation of %p to %zu bytes from %zu bytes "
            "failed.\n",
            ptr, newsize, oldsize);
    exit(EXIT_FAILURE);
  }

  xmalloc_lock();
  xmalloc_prelude_t *new_prelude = newptr;
  new_prelude->size = newsize;
  new_prelude->alloc_count++;

  // Update the "memory used" value

  __memoryused += newsize - oldsize;
  xmalloc_unlock();

#ifdef DEBUG_XMALLOC
  printf("[DEBUG] xrealloc(%p, %lu) -> %p\n", ptr, newsize, newptr);
#endif

  return (void *)(newptr + PRELUDE_SIZE);
}

void *xcalloc(unsigned count, unsigned size) {
  unsigned int n = count * size;
  // Errors should be handled in the xmalloc function
  void *new_mem = (void *)xmalloc(n);
  memset(new_mem, '\0', n);
  return new_mem;
}
