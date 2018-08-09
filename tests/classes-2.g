# Name = "classes 2 (Assigning and Accessing)"
# ExpectedOutput = "Bob Smith"

is main

include "std:io"

class Person {
	string name;
}

func main int {
	Person bob;
	
	bob.name <- "Bob Smith";
	
	io:print("%s", bob.name);
	return 0;	
}
