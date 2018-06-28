func foo (float a, int b) float {
	return float(printf("%d\n", int(a)));
}

func foo (int a) int {
	return printf("%d\n", a);
}

func main int {
	return int(foo(1.1, 1));
}