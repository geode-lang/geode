#include "../include/runtime.h"

#include <stdarg.h>
#include <stdio.h>

void __init_c_runtime() {
  GC_init();
  GC_enable_incremental();
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
