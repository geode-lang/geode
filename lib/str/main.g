is str

class String {
	byte* data;
}


func len(string str) int {
	int len := 0;
	while str[len] != 0 {
		len <- len + 1;
	}
	return len;
}