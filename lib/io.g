link "io.c"

func print(string format, ...) ...
func print(int a) -> print("%d\n", a);
func print(float a) -> print("%f\n", a);
func sprintf(byte* buf, string format, ...) ...


# Some File IO functions.
# The point of these functions in the stdlib is to make
# it easy for other parts of the stdlib to read files and
# streams like stdio/stderr. 
func __openfile(string path, string mode) byte* ...
func __readchar(byte* fp) byte ...
func __fileeof(byte* fp) int ...
func __filewritestring(byte*fp, string data) int ...


# class File {
# 	byte* fp;
# 	func readc() byte {
# 		return __readchar(this.fp);
# 	}
# 	func readall() string {
# 		# ...
# 	}
# }

# func openFile(string path, string mode) File {
	
# }