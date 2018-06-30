# Name = "Conditional nesting"
# RunArgs = ["1", "2", "3"]
# ExpectedOutput = "2 3 4 4 end 4\n"

func main(int argc) int {
	int a := 1;
	if argc >= 2 {
		a <- 2;
		# then 0
		print("2 ");
		if argc >= 3 {
			a <- 3;
			print("3 ");
			if argc >= 4 {
				a <- 4;
				print("4 ");
				if argc >= 5 {
					a <- 5;
					print("5 ");
					if argc >= 6 {
						a <- 6;
						print("6 ");
					} else {print("5 end ");}
				} else {print("4 end ");}
			} else {print("3 end ");}
		} else {print("2 end ");}
		# merge 1
	} else {print("1 end ");}
	# merge 0
	print("%d\n", a);
	return 0;
}
