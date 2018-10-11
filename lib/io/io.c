#include <stdarg.h>
#include <stdio.h>
#include <unistd.h>

#include "../include/mem.h"
#include "../include/xmalloc.h"
#include "io.h"
#include <gc/gc.h>

// the print function wrapper.
void print(char *fmt, ...) {
  va_list args;
  va_start(args, fmt);
  vprintf(fmt, args);
  va_end(args);
}

void sleepms(double ms) { usleep(ms * 1000); }

FILE *get_default_file_descriptor(int index) {
  FILE *fds[] = {stdin, stdout, stderr};
  return fds[index];
}
