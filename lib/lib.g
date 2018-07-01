func print(string format, ...) ...

func print(int a) -> print("%d\n", a);
func print(float a) -> print("%f\n", a);

# Read a file's contents entirely
func readfile(string path) byte* ...

func exp(int x, int n) int {
	if n = 0 {
		return 1;
	}
	return x * exp(x, n - 1);
}