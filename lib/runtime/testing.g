is runtime

include "io"

# the testing section of runtime includes functions that can be
# used in the testing system of geode

# fatalf takes an exit status, a format, and a variadic list
#        and logs the formatted information to stderr then
#        exits with the status code provided
func fatalf(int err, byte* fmt, ...) ...

# assert takes a message and a boolean case and if the case 
# is false, it logs the message with an "Assertion Failed:" prefix
# and exits the program
func assert(byte* msg, bool case) {
	if !case {
		fatalf(-1, "Assertion Failed: %s", msg) # simply fatally log to stderr
	}
}
