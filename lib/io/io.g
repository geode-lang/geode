is io

link "io.c"




# Printing bindings
func print(string format, ...) ...
func fprintf(FILE* handle, string format, ...) ...


func println(string message) {
	print("%s\n", message);
}

