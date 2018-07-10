is mem


# garbage collected malloc
func gcmalloc(int size) byte* ...
func get(int size) byte* {
	return gcmalloc(size);
}


# non garbage collected malloc
func rawmalloc(int size) byte* ...
func unsafe(int size) byte* {
	return rawmalloc(size);
}


func rawfree(byte* ptr) ...
func free(byte* ptr) {
	
}