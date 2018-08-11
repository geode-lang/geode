is main

# include "std:io"
include "std:str"

# class Person {
# 	string name;
# 	int age;
# 	int friendCount;
# 	int something;
# }

# # func echo<T>(T a) T -> a;

func main int {
	# Person* p := mem:get(sizeof(Person)) as Person*;
	# p.name <- "Nick";
	# p.age <- 19;
	# p.friendCount <- 100000;
	# io:println(format("%d", 30));
	
	byte* data := mem:get(30);
	
	# if str:eq("hello", "hallo") {
	# 	io:println("Equal");
	# } else {
	# 	io:println("Not Equal");
	# }

	return 0;
}


# func foo bool -> false;