is mem

link "mem.c"

# garbage collected malloc
func nomangle gcmalloc(int size) byte* ...

func get(int size) byte* {
	return gcmalloc(size);
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
