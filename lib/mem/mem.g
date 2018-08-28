is mem

link "mem.c"

# garbage collected malloc
func GC_malloc(int size) byte* ...
func GC_get_heap_size() long ...

func malloc(long size) byte* ...

# func size(T? ptr) long {
# 	return xmalloc_size(ptr as byte*);
# }



func raw(long size) byte* {
	return malloc(size);
}


func get(long size) byte* {
	return GC_malloc(size);
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


func used() long ...