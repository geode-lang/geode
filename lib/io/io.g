is io

link "io.c"


# Printing bindings
func print(string format, ...) ...
func format(string format, ...) string ...

func println(string message) {
	print("%s\n", message);
}



func getchar() byte ...
func system(string command) int ...
func sleepms(float ms) ...

