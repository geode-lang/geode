@.str_fa90_2 = constant [14 x i8] c"hello, world\0A\00"

declare i64 @printf(i8* %format, ...)

declare i8 @getchar()

define i32 @main(i32 %argc) {
entry_b741_1:
	%0 = alloca i32
	store i32 %argc, i32* %0
	%1 = call i64 (i8*, ...) @printf(i8* getelementptr ([14 x i8], [14 x i8]* @.str_fa90_2, i32 0, i32 0))
	%2 = trunc i64 0 to i32
	ret i32 %2
}
