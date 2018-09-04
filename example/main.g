is main

include "io"
include "encoding"

byte* start;
byte* end;

func write_stack(io:File* target) {
	# Set end to the current end of the stack
	byte a
	end <- &a + 1
	io:fprintf(target, "Address        Binary   0x char\n")
	for let i := start; i >= end; i -= 1 {
		byte* addr := i as byte*
		byte c := *addr
		io:fprintf(target, "%p ", addr)
		io:fprintf(target, "%s ", encoding:binary(*addr))
		io:fprintf(target, "%s ", encoding:hex(*addr))
		if (c >= 0x20) || (c >= 0x7e) {
			io:fprintf(target, "'%c'", c)
		}
		io:fprintf(target, "\n")
	}
}

func main(int argc, byte** argv) int {
	start <- &argc	
	io:File* f := io:fopen("foo.txt", "w+")
	write_stack(io:stdout)
	io:fclose(f)
	return 0
}
