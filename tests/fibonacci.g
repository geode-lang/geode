# Name = "fibonacci"
# ExpectedOutput = "832040"

func fib(int n) int {
	if n < 2 {
		return n;
	}
	return fib(n - 1) + fib(n - 2);
}

func main(int argc) int {
	int a;
	a <- 30;
	print("%d", fib(a));
	return 0;
}
