is main

include "io"

class Person {
	string name;
}

func main int {
	Person bob;
	
	bob.name = "Bob Smith";
	
	io:print("%s", bob.name);
	return 0;	
}
