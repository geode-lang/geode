is mem

link "mem.c"

# garbage collected malloc
func GC_gcollect() ...

func heap_size() long ...


func size(byte* ptr) long {
	return xmalloc_size(ptr);
}

func get(long size) byte* {
	return xmalloc(size);
}

func set(byte* ptr, int size, byte val) {
	for int i := 0; i < size; i += 1 {
		ptr[i] <- val;
	}
}

func zero(int size) byte* {
	byte* data := get(size);
	set(data, size, 0 as byte);
	return data;
}


# mem:collect forces a garbage collector collection and returns
# the number of bytes that were freed
func collect() long {
	let before := mem:heap_size()
	GC_gcollect()
	let after := mem:heap_size()
	return before - after
}