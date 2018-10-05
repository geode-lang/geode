is mem

link "mem.c"



func bytes_used long ...
func blocks_used long ...

# garbage collected malloc
func GC_gcollect() ...

func heap_size() long ...


func size(byte* ptr) long {
	return xmalloc_size(ptr);
}

func get(long size) byte* {
	return xmalloc(size);
}

func resize(byte* ptr, long size) byte* {
	if size(ptr) < size {
		return xrealloc(ptr, size)
	}
	return ptr
}

func set(byte* ptr, int size, byte val) {
	for i = 0; i < size; i += 1 {
		ptr[i] = val;
	}
}

func zero(int size) byte* {
	data = get(size);
	set(data, size, 0);
	return data;
}


# mem:collect forces a garbage collector collection and returns
# the number of bytes that were freed
func collect() long {
	before = mem:heap_size()
	GC_gcollect()
	after = mem:heap_size()
	return before - after
}
