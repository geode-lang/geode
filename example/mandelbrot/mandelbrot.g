is mandel

include "io"
include "str"
include "math"


func printdensity(int d, int iter) {
	r = d % 255
	if d > iter {
		r = 0
	}
	a = io:format("\x1b[48;2;%d;%d;%dm ", r, r, r)
	io:fputs(a, io:stdout)
	io:fputs("\x1b[0m", io:stdout)
}


func mandelconverger(float real, float imag, float iters, float creal, float cimag, int iter) float {
	if iters > iter || (real * real + imag * imag > 4) {
		return iters
	} else {
		return mandelconverger(real * real - imag * imag + creal, 2.0 * real * imag + cimag, iters + 1.0, creal, cimag, iter)
	}
	return 0.0
}




func mandelconverge(float real, float imag, int iter) float {
	return mandelconverger(real, imag, 0.0, real, imag, iter)
}

func printMandel(float realstart, float imagstart, float zoom, int iter, float width, float height) {
	# x-values
	xmin = realstart - zoom * (width / 2.0)
	xmax = realstart + zoom * (width / 2.0)
	# y-values
	ymin = imagstart - zoom * (height / 2.0)
	ymax = imagstart + zoom * (height / 2.0)
	mandelhelp(xmin, xmax, zoom, ymin, ymax, zoom, iter)	
}


func mandelhelp(float xmin, float xmax, float xstep, float ymin, float ymax, float ystep, int iter) {

	
	max = 0.0
	
	for y = ymin; y < ymax; y += ystep {
		for x = xmin; x < xmax; x += xstep {
			cov = mandelconverge(x,y, iter)
			printdensity(cov, iter)
		}
		io:fputs("\n", io:stdout)
	}
	io:fflush(io:stdout)
	
}

func clear {
	io:print("\x1b[H\x1b[2J\n")
}



func main int {
	z = 0.00014
	x = -0.9250001355432285
	y =  0.2660002226258663
	iter = 512
	
	
	width = 128.0
	height = 64.0
	
	minzoom = 0.0000000000000001
	maxzoom = 0.07
	while true {
		
		clear()
		printMandel(x, y, z, iter, width, height)

		io:system("stty raw")
		input = io:getchar()
		io:system("stty cooked")

		io:print("\n")
		
		step = z * 2
		if input == 'l' {
			x += step
		}
		
		if input == 'h' {
			x -= step
		}
		
		if input == 'j' {
			y += step
		}
		
		if input == 'k' {
			y -= step
		}
		
		if input == 'x' {
			z = z * 2.0
		}
		
		if input == 'z' {
			z = z / 2.0
		}
		
		if input == 's' {
			iter = iter * 2
		}
		
		if input == 'a' {
			iter = iter / 2
		
		}
		
		if iter >= 500000 {
			iter = 500000
		}
		if iter == 0 {
			iter = 1
		}

		if input == '\x03' || input == 'q' {
			clear()
			return 0
		}

		# Furthest you are allowed to zoom in
		if z < minzoom {
			z = minzoom
		}

		# Furthest you are allowed to zoom out
		if z > maxzoom {
			z = maxzoom
		}

	}

	return 0
}
