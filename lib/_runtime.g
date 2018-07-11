is _runtime

link "./tgc/tgc.c"
link "c/_runtime.c"

func ___geodegcinit(byte* stk) ...

func init(byte* stk) {
	___geodegcinit(stk);
}


