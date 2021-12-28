// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
package tingyun3

import (
	"runtime"
)

type httpListenAddress struct {
	Addr string
	tls  bool
}

var httpListenAddr httpListenAddress

//GetCallerName : 取layer层调用栈函数名
//go:noinline
func GetCallerName(layer int) string {
	if _, pc := GetCallerPC(layer + 1); pc != 0 {
		return runtime.FuncForPC(pc).Name()
	}
	return ""
}

// GetCallerPC return caller pc
//go:noinline
func GetCallerPC(layer int) (l int, pc uintptr) {
	if pc, _, _, success := runtime.Caller(layer); success {
		return layer, pc
	}
	return 0, 0
}
