is io

link "io.c"

include "std:mem"


# Printing bindings
func print(string format, ...) ...
func format(string format, ...) string ...

func println(string message) {
	print("%s\n", message);
}