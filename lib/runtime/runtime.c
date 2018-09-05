#include "../include/runtime.h"

#include "stdio.h"



void __initruntime() {
	GC_init();
	GC_enable_incremental();
}