#include <stdio.h>

extern int add (int a, int b);
int main() {
	printf("%d", add(1, 2));
	return 1;
}