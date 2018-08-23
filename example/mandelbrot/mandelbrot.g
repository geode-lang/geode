is mandelbrot

include "std:io"
include "std:str"
include "std:color"
include "std:math"


func printdensity(int d, int iter) {
	Color c;
	if d > iter {
		c <- color:new_rgb(0,0,0);
	} else {
		# c <- color:new_rgb(d % 255, d % 255, d % 255);
		c <- color:hsv_to_rgb(d * 10.0, 1.0, 1.0);
	}

	io:print("\x1b[48;2;%d;%d;%dm", c.r as int, c.g as int, c.b as int);
	io:print(" "); # Each "pixel" of the mandelbrot is simply a space with a background color
	io:print("\x1b[0m"); # Reset the background color
	return;
}


func mandelconverger(float real, float imag, float iters, float creal, float cimag, int iter) float {
	if iters > iter || (real * real + imag * imag > 4) {
		return iters;
	} else {
		return mandelconverger(real * real - imag * imag + creal, 2.0 * real * imag + cimag, iters + 1.0, creal, cimag, iter);
	}	
}



func mandelconverge(float real, float imag, int iter) float {
	return mandelconverger(real, imag, 0.0, real, imag, iter);
}

func printMandel(float realstart, float imagstart, float zoom, int iter, float width, float height) {
	io:print("%.40f\n", zoom);
	float xmin := realstart - zoom * (width / 2.0);
	float xmax := realstart + zoom * (width / 2.0);
	
	float ymin := imagstart - zoom * (height / 2.0);
	float ymax := imagstart + zoom * (height / 2.0);
	mandelhelp(xmin, xmax, zoom, ymin, ymax, zoom, iter);	
}


func mandelhelp(float xmin, float xmax, float xstep, float ymin, float ymax, float ystep, int iter) {

	io:print("x: %.40f\n", xmin);
	io:print("y: %.40f\n", ymin);
	io:print("zoom: %f\n", xstep);
	io:print("iter: %d\n", iter);
	
	float max := 0;
	
	for float y := ymin; y < ymax; y += ystep {
		for float x := xmin; x < xmax; x += xstep {
			float cov := mandelconverge(x,y, iter);
			printdensity(cov, iter);
		}
		io:print("\n");
	}
	
}

func clear void {
	io:print("\x1b[H\x1b[2J\n");
}



func main int {
	float z := 0.00014;
	float x := -0.9250001355432285;
	float y :=  0.2660002226258663;
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

		io:print("\n");
		
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
		
		if iter >= 500000 {
			iter <- 500000;
		}
		if iter = 0 {
			iter <- 1;
		}
	
		if input = '\x03' || input = 'q' {
			clear();
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
