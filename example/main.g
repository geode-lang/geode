include "std::io"
# include "std::mem"


func main int {
	for int i := 0; i >= 0; i <- i + 1 {
		byte* foo := malloc(300);
		print("%p\n", foo);
	}
	
	return 42;
}
