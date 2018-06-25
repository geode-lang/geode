func main(int argc) byte {
	int a := 1;
	if argc >= 2 {
		a <- 2;
		# then 0
		printf("at least 2 args\n");
		if argc >= 3 {
			a <- 3;
			printf("at least 3 args\n");
			if argc >= 4 {
				a <- 4;
				printf("at least 4 args\n");
				if argc >= 5 {
					a <- 5;
					printf("at least 5 args\n");
					if argc >= 6 {
						a <- 6;
						printf("at least 6 args\n");
					} else {printf("5's it\n");}
				} else {printf("4's it\n");}
			} else {printf("3's it\n");}
		} else {printf("2's it\n");}
		# merge 1
	} else {printf("1 arg :(\n");}
	# merge 0
	printf("%d\n", a);
	return 0;
}
