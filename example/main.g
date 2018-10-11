is main

# github.com/geode-lang/geode

func main int {
	log("hello, world")
	return 0
}

# func parseint(string str, int base) long {
# 	set = "0123456789abcdef"
# 	setlen = str:len(set)
# 	long res = 0
# 	digit = 0
# 	for i = 0; str[i] != 0; i += 1 {
# 		for c = 0; c <= setlen; c += 1 {
# 			if set[c] == str[i] {
# 				digit = c
# 			}
# 		}
# 		res = res * base + digit
# 	}
# 	return res
# }
