is main
include "std:io"
include "std:mem"
include "std:math"


func main(int argc) int {
	int a := 1;
	if argc >= 2 {
		a <- 2;
		# then 0
		io:print("2 ");
		if argc >= 3 {
			a <- 3;
			io:print("3 ");
			if argc >= 4 {
				a <- 4;
				io:print("4 ");
				if argc >= 5 {
					a <- 5;
					io:print("5 ");
					if argc >= 6 {
						a <- 6;
						io:print("6 ");
					} else {io:print("5 end ");}
				} else {io:print("4 end ");}
			} else {io:print("3 end ");}
		} else {io:print("2 end ");}
		# merge 1
	} else {io:print("1 end ");}
	# merge 0
	io:print("%d\n", a);
	return 0;
}
