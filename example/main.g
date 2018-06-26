# Implicit returns
func foo string -> "hello there. How are you?";

func main int {
	printf("%s\n", foo());
	return 0;
}