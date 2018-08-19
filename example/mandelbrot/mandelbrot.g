is main


include "std:io"
include "std:str"


func printdensity(int d) {
	
	if d > 15 {
		io:print(" ");
		return;
	}
	
	string charset := ".,*xX#";
	io:print("%c", charset[d % str:len(charset)]);
	return;
}


func mandelconverger(float real, float imag, float iters, float creal, float cimag) float {
	if iters > 255 || (real * real + imag * imag > 4) {
		return iters;
	} else {
		return mandelconverger(real * real - imag * imag + creal, 2.0 * real * imag + cimag, iters + 1.0, creal, cimag);
	}	
}



func mandelconverge(float real, float imag) float {
	return mandelconverger(real, imag, 0.0, real, imag);
}

func printMandel(float realstart, float imagstart, float realmag, float imagmag) {
	mandelhelp(realstart, realstart+realmag*128, realmag, imagstart, imagstart+imagmag*64, imagmag);	
}


func mandelhelp(float xmin, float xmax, float xstep, float ymin, float ymax, float ystep) {
	for float y := ymin; y < ymax; y += ystep {
		for float x := xmin; x < xmax; x += xstep {
			printdensity(mandelconverge(x,y));
		}
		io:print("\n");
	}
}

func main int {
	printMandel(-2.5, -1.3, 0.03, 0.04);
	return 0;
}
