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

//go:noinline
func GetRootCallerName(layer int) string {
	var pc uintptr = 0
	for {
		if addr, _, _, success := runtime.Caller(layer); success {
			pc = addr
			layer++
		} else {
			break
		}
	}
	if pc != 0 {
		return runtime.FuncForPC(pc).Name()
	}
	return ""
}

//go:noinline
func MatchCallerName(layer int, funcname string) bool {
	for {
		if addr, _, _, success := runtime.Caller(layer); success {
			if runtime.FuncForPC(addr).Name() == funcname {
				return true
			}
			layer++
		} else {
			break
		}
	}
	return false
}

// GetCallerPC return caller pc
//go:noinline
func GetCallerPC(layer int) (l int, pc uintptr) {
	if pc, _, _, success := runtime.Caller(layer); success {
		return layer, pc
	}
	return 0, 0
}

const (
	StorageIndexDatabase = 1 + 8*0
	StorageIndexRedis    = 1 + 8*1
	StorageIndexEcho     = 1 + 8*2
	StorageIndexMongo    = 1 + 8*3
	StorageIndexBeego    = 1 + 8*4
	StorageIndexIris     = 1 + 8*5
	StorageIndexGin      = 1 + 8*6
	StorageGorillaWS     = 1 + 8*7
)
