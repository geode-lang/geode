is mem

link "mem.c"

# garbage collected malloc
func nomangle gcmalloc(int size) byte* ...

func get(int size) byte* {
	return gcmalloc(size);
}
