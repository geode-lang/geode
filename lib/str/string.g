is str

include "mem" # needed for allocating memory for strings

class String {
	byte* data;
}



# str:len
# A basic string length function that loops over the bytes
# incrementing a counter till it encounters a null byte
func len(string str) int {
	int len := 0;
	while str[len] != 0 {
		len += 1;
	}
	return len;
}

# str:eq
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
		# Check if this character isn't equal
		if a[i] != b[i] {
			return false;
		}
	}	
	return true;
}



# str:concat
# concatinate two strings into a single string.
func concat(string a, string b) string {
	# cache the size of both strings for later use
	int alen := len(a);
	int blen := len(b);
	# The final number of usable chars in the new string
	int finalLen := alen + blen;
	# allocate finalLen+1 bytes as zero so it is zero
	# delimited
	string combined := mem:zero(finalLen + 1);
	# loop over all usable bytes of the new data and
	# copy over the old strings into the new buffer.
	for int i := 0; i < finalLen; i += 1 {
		if i < alen {
			combined[i] <- a[i];
		} else {
			combined[i] <- b[i - alen];
		}
	}
	return combined;
}