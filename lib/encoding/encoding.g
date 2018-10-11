is encoding

include "str"
include "mem"
include "io"




string hex_charset = "0123456789abcdef"
# hex converts type T (unknown) to a string containing
# the hex representation of val
func hex(T? val) string {
	# get the info of the type T for use later.
	tInfo = info(T)
	buffer = mem:zero(tInfo.size * 2 + 1)
	byte* offset = &val

	for i = tInfo.size - 1; i >= 0; i -= 1 {
		b = offset[i]
		o = (tInfo.size - i - 1) * 2
		buffer[o] = hex_charset[b >> 4 && 0xf]
		buffer[o+1] = hex_charset[b && 0xf]
	}
	return buffer;
}


# binary converts type T (unknown) to a string containing
# the binary representation of val
func binary(T? val) string {
	# the actual string that will be changed to contain the binary string
	buffer = "";
	bin_buffer = "00000000";
	byte* offset = &val;
	for int i = info(T).size - 1; i >= 0; i -= 1 {
		byte b = offset[i];
		for int o = 7; o >= 0; o -= 1 {
			byte bit = (b >> o) && 1;
			bin_buffer[7 - o] = bit + '0';
		}
		buffer = str:concat(buffer, bin_buffer);
	}
	return buffer;
}

