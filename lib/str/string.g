is str

include "mem" # needed for allocating memory for strings

class String {
	byte* data;
}



# str:len
# A basic string length function that loops over the bytes
# incrementing a counter till it encounters a null byte
func len(string str) int {
	int len = 0;
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
	alen = len(a);
	blen = len(b);

	# if they arent the same length, they must not be equal
	if alen != blen {
		return false;
	}

	# loop over each char and test if they are equal
	for i = 0; i < alen; i += 1 {
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
	alen = len(a);
	blen = len(b);
	# The final number of usable chars in the new string
	finalLen = alen + blen;
	# allocate finalLen+1 bytes as zero so it is zero
	# delimited
	string combined = mem:zero(finalLen + 1);
	# loop over all usable bytes of the new data and
	# copy over the old strings into the new buffer.
	for i = 0; i < finalLen; i += 1 {
		if i < alen {
			combined[i] = a[i];
		} else {
			combined[i] = b[i - alen];
		}
	}
	return combined;
}


# str:hash
#
#    djb2 hash function implementation for strings
#
# this algorithm (k=33) was first reported by dan bernstein many years ago
# in comp.lang.c. another version of this algorithm (now favored by bernstein)
# uses xor: hash(i) = hash(i - 1) * 33 ^ str[i]; the magic of number 33
# (why it works better than many other constants, prime or not) has
# never been adequately explained.
func hash(string str) long {
	long hash = 5381;
	size = len(str);

	for c = 0; c < size; c += 1 {
		hash = ((hash << 5) + hash) + c # hash * 33 + c
	}
	return hash
}




# split str by all characters in sset and return
# a NULL terminated string buffer
func split(string str, string sset) string* {
	strSize = info(string).size
	count = 0
	splits = mem:get(strSize * count)


	string ptr = str;
	while str[i] {
		char = str[i]

		# check if the char is in the subset
		for c = 0; c < len(sset); c += 1 {
			if 
		}
	}


	return splits
}