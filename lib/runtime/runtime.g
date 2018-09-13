is runtime

link "runtime.c"
link "xmalloc.c"
link "../gc/*.o"

# safer, gc friendly memory functions.
func xmalloc(int size) byte* ...
func xrealloc(byte* ptr, int size) byte* ...

func memcpy(byte* dest, byte* src, int length) ...
func xmalloc_size(byte* ptr) long ...
func __initruntime() ...
# func raw_copy(byte* target, int length) byte* ...
func exit(int code) ...




func raw_copy(byte* source, int len) byte* {
	dest = xmalloc(len);
	memcpy(dest, source, len);
	return dest;
}


class TypeInfo {
	int size
	string name
}



# geode bindings to the unsafe c malloc/free.
# these allow use when you disable the GC with the --no-runtime flag
func malloc(int size) byte* ...
func free(byte* ptr) ...