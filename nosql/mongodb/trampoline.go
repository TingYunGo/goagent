package mongodb

import (
	"fmt"
)

type trampolineStruct struct {
	idpointer *int64
	idindex   int64
	arg1      int64
	arg2      int64
	arg3      int64
	arg4      int64
	arg5      int64
	arg6      int64
	arg7      int64
	arg8      int64
	arg9      int64
	arg10     int64
	arg11     int64
	arg12     int64
	arg13     int64
	arg14     int64
	arg15     int64
	arg16     int64
	arg17     int64
	arg18     int64
	arg19     int64
	arg20     int64
}

var trampoline *trampolineStruct = &trampolineStruct{}

//go:noinline
func initTrampoline(p *trampolineStruct) int64 {
	trampoline = p
	fmt.Println(1, p)
	return 0
}
