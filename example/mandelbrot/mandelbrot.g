is mandelbrot


include "std:io"
include "std:str"
include "std:color"



func printdensity(int d, int iter) {
	Color c;
	if d > iter {
		c <- color:NewColorRGB(0,0,0);
	} else {
		c <- color:NewColorHSV(d, 0.68, 1);
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
		
	io:print("x: %f\n", xmin);
	io:print("y: %f\n", ymin);
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

