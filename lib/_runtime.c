#include "./gc/tgc.h"

static tgc_t gc;

extern long __GEODE__main(int argc, char** argv);

int
main(int argc, char** argv) {
	
	tgc_start(&gc, &argc);
  int prog_ret_val = __GEODE__main(argc, argv);
  tgc_stop(&gc);
	return prog_ret_val;
}


char*
__GEODE__alloca(int size) {
	// tgc_run(&gc);
	return tgc_alloc(&gc, size);
}