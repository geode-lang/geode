func fib(int n) int {
	if n < 2 {return n;}
	return fib(n - 1) + fib(n - 2);
}

func main(int argc) byte {
	printf("%d\n", fib(30));
	return 0;
}