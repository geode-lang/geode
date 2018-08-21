is main

include "mandelbrot/mandelbrot.g"
include "std:io"
include "std:math"


func clear void {
	io:print("\x1b[H\x1b[2J\n");
}



func main int {
	float z := 0.00014;
	float x := -0.235125;
	float y :=  0.827215;
	int iter := 512;
	
	
	float width := 128;
	float height := 64;
	
	float minzoom := 0.0000000000000001;
	float maxzoom := 0.07;
	while true {
		
		clear();
		mandelbrot:printMandel(x, y, z, iter, width, height);
		
		
		
		io:system("stty raw");
		byte input := io:getchar();
		io:system("stty cooked");
		
		float step := z * 2;
		if input = 'l' {
			x += step;
		}
		
		if input = 'h' {
			x -= step;
		}
		
		if input = 'j' {
			y += step;
		}
		
		if input = 'k' {
			y -= step;
		}
		
		if input = 'x' {
			z *= 2.0;
		}
		
		if input = 'z' {
			z /= 2.0;
		}
		
		if input = 's' {
			iter *= 2;
		}
		
		if input = 'a' {
			iter /= 2;
		
		}
		
		if iter >= 50000 {
			iter <- 50000;
		}
		if iter = 0 {
			iter <- 1;
		}
	
		if input = '\x03' || input = 'q' {
			return 0;
		}
		
		
		
		# Furthest you are allowed to zoom in
		if z < minzoom {
			z <- minzoom;
		}
		
		# Furthest you are allowed to zoom out
		if z > maxzoom {
			z <- maxzoom;
		}
		
		
		
		
	}
	
	
	return 0;
}
