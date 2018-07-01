func main(int argc, byte** argv) int {
	for int i := 0; i < argc; i <- i + 1 {
		print("%d: %s\n", i + 1, argv[i]);
	}
	return argc;
}