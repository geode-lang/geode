is encoding

link "b64_encode.c"

func b64_encode(string src, int len) string ...

include "std:str"
include "std:mem"


# hex converts type T (unknown) to a string containing
# the hex representation of val
func hex(T? val) string {
	string charset := "0123456789abcdef";
	string buffer := mem:zero(sizeof(T) * 2 + 1);
	byte* offset := &val;
	for int i := sizeof(T) - 1; i >= 0; i -= 1 {
		byte b := offset[i];
		int o := (sizeof(T) - i - 1) * 2;
		buffer[o] <- charset[b >> 4 && 0xf];
		buffer[o+1] <- charset[b && 0xf];
	}
	return buffer;
}


# binary converts type T (unknown) to a string containing
# the binary representation of val
func binary(T? val) string {
	# the actual string that will be changed to contain the binary string
	string buffer := "";
	string bin_buffer := "00000000";
	byte* offset := &val;
	for int i := sizeof(T) - 1; i >= 0; i -= 1 {
		byte b := offset[i];
		for int o := 7; o >= 0; o -= 1 {
			byte bit := (b >> o) && 1;
			bin_buffer[7 - o] <- bit + '0';
		}
		buffer <- str:concat(buffer, bin_buffer);
	}
	return buffer;
}

func base64(T? val) string {
	byte* source := &val;
	
	return b64_encode(source, sizeof(T));
}