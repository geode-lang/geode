is main

include "mandelbrot/mandelbrot.g"
include "std:io"
include "std:math"



func main int {
	# io:print("\u0e27\u0e23\u0e0d\u0e32");
	# string clear := "\x1b[H\x1b[2J\n";
	float i := 0.02;
	while true {
		i -= 0.0001;
		io:sleepms(60);
		io:print("\x1b[H\x1b[2J\n");
		mandelbrot:printMandel(-0.925, 0.266, i, i);
	}
	
	
	return 0;
}