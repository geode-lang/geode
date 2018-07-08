is mem


func rawmalloc(int size) byte* ...

func get(int size) byte* {
	return rawmalloc(size);
}