package gen

// RuntimeSource is the source the runtime will use
const RuntimeSource string = `
# The exponent operator function
func exp(int x, int n) int {
	if n = 0 {
		return 1;
	}
	return x * exp(x, n - 1);
}
`
