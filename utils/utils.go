package utils

import (
	"fmt"
	"runtime"

	"github.com/TingYunGo/goagent"
)

func getnameByAddr(p uintptr) string {
	if r := runtime.FuncForPC(p); r == nil {
		return ""
	} else {
		file, line := r.FileLine(p)
		return fmt.Sprintf("%x:%s(%s:%d)", p, r.Name(), file, line)
	}
}

func PrintCaller(skip, stacks int) {
	stackList := make([]uintptr, stacks)
	count := runtime.Callers(skip+1, stackList)
	fmt.Println("Routine:", tingyun3.GetGID())
	for i := 0; i < count; i++ {
		name := getnameByAddr(stackList[i] - 1)
		fmt.Println(i, name)
	}
}
