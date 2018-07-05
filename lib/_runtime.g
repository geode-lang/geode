link "./gc/tgc.c"
link "_runtime.c"


func __GEODE__alloca(int size) byte* ...

func malloc(int size) byte* {
	return __GEODE__alloca(size);
}