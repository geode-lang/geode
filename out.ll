@.str_2 = constant [26 x i8] c"hello there. How are you?\00"

@.str_4 = constant [4 x i8] c"%s\0A\00"

declare i64 @printf(i8* %format, ...)

declare i8 @getchar()

declare i8* @malloc(i64 %size)

define i8* @foo() {
entry_1:
	ret i8* getelementptr ([26 x i8], [26 x i8]* @.str_2, i32 0, i32 0)
}

define i64 @main() {
entry_3:
	%0 = call i8* @foo()
	%1 = call i64 (i8*, ...) @printf(i8* getelementptr ([4 x i8], [4 x i8]* @.str_4, i32 0, i32 0), i8* %0)
	ret i64 0
}
