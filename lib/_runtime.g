link "./tgc/tgc.c"
link "_runtime.c"

func __GEODE__alloca(int size) byte* ...
func __GEODE__free(int size) byte* ...

func malloc(int size) byte* {
	return __GEODE__alloca(size);
}

func free(byte* ptr) {
	return __GEODE__free(ptr);
}