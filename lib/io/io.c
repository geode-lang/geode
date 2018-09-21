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

// the format function wrapper.
char *format(char *fmt, ...) {
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

void sleepms(double ms) { usleep(ms * 1000); }
