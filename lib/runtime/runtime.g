is runtime

link "runtime.c"
link "xmalloc.c"
# link "../gc/*.o"

# safer, gc friendly memory functions.
func xmalloc(int size) byte* ...
func xrealloc(byte* ptr, int size) byte* ...
func memcpy(byte* dest, byte* src, int length) ...
func xmalloc_size(byte* ptr) long ...
func __initruntime() ...
func exit(int status) ...
func kill(int pid, int status) ...



# binding to the write syscall
func write(int fd, byte* buf, long nbytes) long ...

func write'(int fd, byte* msg) long {
	len = 0
	while msg[len] != 0 { len += 1 }
	return write(1, msg, len)
}
 
func out(byte* msg) long -> write'(1, msg)
func werr(byte* msg) int -> write'(2, msg)


func raw_copy(byte* source, int len) byte* {
	dest = xmalloc(len);
	memcpy(dest, source, len);
	return dest;
}

class TypeInfo {
	int size
	string name
}

