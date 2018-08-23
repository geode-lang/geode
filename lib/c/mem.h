#ifndef __MEM__
#define __MEM__

char* gcmalloc(long size);
char* rawmalloc(long size);
void rawfree(char* ptr);

#endif