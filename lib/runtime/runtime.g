is runtime

link "runtime.c"
link "xmalloc.c"
link "../gc/*.o"


func xmalloc(int size) byte* ...
func memcpy(byte* dest, byte* src, int length) ...
func xmalloc_size(byte* ptr) long ...
func __initruntime() ...
# func raw_copy(byte* target, int length) byte* ...
func exit(int code) ...



func raw_copy(byte* source, int len) byte* {
	byte* dest := xmalloc(len);
	memcpy(dest, source, len);
	return dest;
}
