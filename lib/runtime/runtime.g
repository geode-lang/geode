is runtime

link "runtime.c"
link "xmalloc.c"

# safer, gc friendly memory functions.
func xmalloc(int size) byte* ...
func xrealloc(byte* ptr, int size) byte* ...
func memcpy(byte* dest, byte* src, int length) ...
func xmalloc_size(byte* ptr) long ...
func __init_c_runtime() ...
func exit(int status) ...
func kill(int pid, int status) ...



func read(int fd, byte* buf, long nbytes) long ...

# binding to the write syscall
func write(int fd, byte* buf, long nbytes) long ...
func write'(int fd, byte* msg) long {
	len = 0
	while msg[len] != 0 { len += 1 }
	return write(1, msg, len)
}


func wout(byte* msg) long = write'(1, msg)
func werr(byte* msg) int = write'(2, msg)


func raw_copy(byte* source, int len) byte* {
	dest = xmalloc(len);
	memcpy(dest, source, len);
	return dest;
}


func __init_runtime() {
	# this function doesn't do anything right now, but it does
	# get populated with things later on in the compiler phase
	# init the runtime via the c call
	__init_c_runtime();
}

# typeinfo is what is returned from the info(T) call.
# The instance contains information about the type T
class TypeInfo {
	# The size in bytes of the type
	int size

	# the name of the type
	string name
}














