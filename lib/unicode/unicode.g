is unicode

link "microutf8.c"

include "std:mem"
include "std:str"

func _utf8_strlen(string s) int ...
func _utf8_get_nth_char(string str, string target, int offset, byte* is_bom_present, int n, int max_length) byte ...


func len(string s) int {
	return _utf8_strlen(s);
}

func idx(string s, int i) string {
	string target := mem:zero(4);
	byte* bom;
	_utf8_get_nth_char(s, target, 0, bom, i, str:len(s));
	return target;
}