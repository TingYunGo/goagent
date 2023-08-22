// Copyright 2023 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64 arm64
// +build cgo

package wrapruntime

import (
	"unsafe"
)

//go:noinline
func selectgo(cas0 *scase, order0 *uint16, pc0 *uintptr, nsends, nrecvs int, block bool) (int, bool) {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return 0, false
}

type scase struct {
	c    *hchan         // chan
	elem unsafe.Pointer // data element
}

//go:noinline
func Wrapselectgo(cas0 *scase, order0 *uint16, pc0 *uintptr, nsends, nrecvs int, block bool) (int, bool) {

	var selectHandler SelectHandler = nil

	ncases := nsends + nrecvs

	if ncases == 0 {
		return selectgo(cas0, order0, pc0, nsends, nrecvs, block)
	}

	cas1 := (*[1 << 16]scase)(unsafe.Pointer(cas0))
	scases := cas1[:ncases:ncases]

	casi := 0
	var cas *scase
	typeName := ""
	for _, handler := range selectHandlerList {
		for casi = 0; casi < ncases; casi++ {
			cas = &scases[casi]
			if cas.c == nil {
				continue
			}
			typeName = TypeString(cas.c.elemtype)
			if selectHandler = handler(typeName); selectHandler != nil {
				break
			}
		}
		if selectHandler != nil {
			break
		}
	}
	id, ok := selectgo(cas0, order0, pc0, nsends, nrecvs, block)
	if selectHandler != nil && ok && id >= 0 {
		cas = &scases[id]
		typeName = TypeString(cas.c.elemtype)
		selectHandler.Ret(typeName, cas.elem, id, ok)
	}
	return id, ok
}

type SelectHandler interface {
	Ret(tyneName string, elem unsafe.Pointer, retId int, retOk bool)
}

type SelectHandle func(typeName string) SelectHandler

var selectHandlerList []SelectHandle = nil

func HandleSelect(handler SelectHandle) {
	selectHandlerList = append(selectHandlerList, handler)
}
