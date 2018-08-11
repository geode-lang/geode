is str

include "std:mem" # needed for allocating memory for strings

class String {
	byte* data;
}



# A basic string length function that loops over the bytes
# incrementing a counter till it encounters a null byte
func len(string str) int {
	int len := 0;
	while str[len] != 0 {
		len <- len + 1;
	}
	return len;
}


# The equal function goes through several tests
# before finally returning true. If any of these
# tests fail along the way, the function will 
# return false.
# These tests include:
#   - same length?
#   - looping over testing each byte for equality
#     across both strings
func eq(string a, string b) bool {
	
	# cache the size of both strings for later use
	int alen := len(a);
	int blen := len(b);
	
	# if they arent the same length, they must not be equal
	if alen != blen {
		return false;
	}
	
	# loop over each char and test if they are equal
	for int i := 0; i < alen; i += 1 {
		if a[i] != b[i] {
			return false;
		}
	}	
	return true;
}



func concat(string a, string b) string {
	# cache the size of both strings for later use
	int alen := len(a);
	int blen := len(b);
	return "hello";
}