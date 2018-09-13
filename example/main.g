is main

include "io"
include "mem"

func main int {
	h = io:open("LICENSE", "rb")
	content = h.readall()
	h.close()
	io:print("%s\n", content)
	return 0
}