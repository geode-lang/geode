func fib(int n) int {
	if n < 2 {
		return n;
	}
	return fib(n - 1) + fib(n - 2);
}

func main(int argc) int {
	int a;
	a <- 30;
	printf("%d\n", fib(a));
	return 0;
}
