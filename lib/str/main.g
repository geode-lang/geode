is str


func len(string str) int {
	int len;
	for len <- 0; str[len] != 0; len <- len + 1 {}
	return len;
}