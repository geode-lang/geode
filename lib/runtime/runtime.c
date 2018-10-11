#include "../include/runtime.h"

#include "../include/xmalloc.h"
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h>

void exit_handle(void) { GC_gcollect(); }

void __init_c_runtime() {
  atexit(exit_handle);
  GC_init();
  // GC_enable_incremental();
}

void fatalf(int err, char *fmt, ...) {
  fputs("Error: ", stderr);
  va_list vargs;
  va_start(vargs, fmt);
  vfprintf(stderr, fmt, vargs);
  printf("\n");
  va_end(vargs);
  exit(err);
}

char *__runtime_str_format(char *fmt, ...) {
  va_list checkArgs;
  va_start(checkArgs, fmt);
  long size = vsnprintf(NULL, 0, fmt, checkArgs);
  va_end(checkArgs);
  // Allocate memory for the string
  char *buffer = xmalloc(size + 1);
  // Reparse the args... There is no way around this, sadly
  va_list args;
  va_start(args, fmt);
  vsnprintf(buffer, size + 1, fmt, args);
  va_end(args);
  return buffer;
}