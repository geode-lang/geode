#include "../include/_runtime.h"

#define USE_GC



void
__runtimeinit() {
	GC_INIT()
	GC_enable_incremental();
}
