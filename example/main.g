is main

include "io"
include "encoding"

byte* start;
byte* end;

func dump_stack() {
	# Set end to the current end of the stack
	byte a
	end <- &a + 1
	
	io:print("Address        Binary   0x char\n")
	for let i := start; i >= end; i -= 1 {
		byte* addr := i as byte*
		byte c := *addr
		io:print("%p ", addr)
		io:print("%s ", encoding:binary(*addr))
		io:print("%s ", encoding:hex(*addr))
		if (c >= 0x20) || (c >= 0x7e) {
			io:print("'%c'", c)
		}
		io:print("\n")
	}
}

func main(int argc, byte** argv) int {
	start <- &argc
	dump_stack()
	return 0
}
