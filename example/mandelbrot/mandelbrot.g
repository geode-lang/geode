is mandelbrot


include "std:io"
include "std:str"

func printdensity(int d) {
	
	if d > 255 {
		io:print(" ");
		return;
	}
	
	string charset := ".,-*x#";
	io:print("%c", charset[d % str:len(charset)]);
	return;
}


func mandelconverger(float real, float imag, float iters, float creal, float cimag) float {
	if iters > 250 || (real * real + imag * imag > 4) {
		return iters;
	} else {
		return mandelconverger(real * real - imag * imag + creal, 2.0 * real * imag + cimag, iters + 1.0, creal, cimag);
	}	
}



func mandelconverge(float real, float imag) float {
	return mandelconverger(real, imag, 0.0, real, imag);
}

func printMandel(float realstart, float imagstart, float realmag, float imagmag) {
	
	
	float width := 120;
	float height := 60;
	float xmin := realstart - realmag*(width/2.0);
	float xmax := realstart + realmag*(width/2.0);
	
	float ymin := imagstart - realmag*(height/2.0);
	float ymax := imagstart + realmag*(height/2.0);
	mandelhelp(xmin, xmax, realmag, ymin, ymax, imagmag);	
}


func mandelhelp(float xmin, float xmax, float xstep, float ymin, float ymax, float ystep) {
	
	io:print("(%f-%f, %f-%f)\n", xmin, xmax, ymin, ymax);
	
	io:print("\n");
	
	for float y := ymin; y < ymax; y += ystep {
		for float x := xmin; x < xmax; x += xstep {
			printdensity(mandelconverge(x,y));
		}
		io:print("\n");
	}
}

