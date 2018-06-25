func main(int argc) byte {
	int a := 1;
	if argc >= 2 {
		a <- 2;
		# then 0
		printf("2 ");
		if argc >= 3 {
			a <- 3;
			printf("3 ");
			if argc >= 4 {
				a <- 4;
				printf("4 ");
				if argc >= 5 {
					a <- 5;
					printf("5 ");
					if argc >= 6 {
						a <- 6;
						printf("6 ");
					} else {printf("5 end ");}
				} else {printf("4 end ");}
			} else {printf("3 end ");}
		} else {printf("2 end ");}
		# merge 1
	} else {printf("1 end ");}
	# merge 0
	printf("%d\n", a);
	return 0;
}
