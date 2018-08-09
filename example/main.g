is main

include "std:io"
include "std:mem"

class Person {
	string name;
	int age;
	int friendCount;
	byte something;
}

func main int {
	Person* p := mem:get(sizeof(Person)) as Person*;
	p.name <- "Nick";
	p.age <- 19;
	p.friendCount <- 100000;
	
	io:print("A single Person takes %d bytes\n", sizeof(Person));
	io:print("%s is %d years old\n", p.name, p.age);
	io:print("%s has %d friends\n", p.name, p.friendCount);
	return 0;
}

