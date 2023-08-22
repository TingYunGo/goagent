// Copyright 2023 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64 arm64
// +build cgo

package wrapruntime

import (
	"unsafe"
)

//func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool

//go:noinline
func chansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return false
}

type ChanSendHandle func(name string, ep unsafe.Pointer, callerpc uintptr) bool

var handlerList []ChanSendHandle = nil

func HandleChanSend(handler ChanSendHandle) {
	handlerList = append(handlerList, handler)
}

//go:noinline
func Wrapchansend(c *hchan, ep unsafe.Pointer, block bool, callerpc uintptr) bool {
	for _, handler := range handlerList {
		if handler(TypeString(c.elemtype), ep, callerpc) {
			break
		}
	}
	return chansend(c, ep, block, callerpc)
}

//func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {

type ChanRecvHandler interface {
	Ret(selected, received bool)
}

type ChanRecvHandle func(name string, ep unsafe.Pointer, block bool) ChanRecvHandler

var recvHandlerList []ChanRecvHandle = nil

func HandleChanRecv(handler ChanRecvHandle) {
	recvHandlerList = append(recvHandlerList, handler)
}

//go:noinline
func chanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return false, false
}

//go:noinline
func Wrapchanrecv(c *hchan, ep unsafe.Pointer, block bool) (selected, received bool) {
	var recvHandler ChanRecvHandler = nil
	for _, handler := range recvHandlerList {
		if recvHandler = handler(TypeString(c.elemtype), ep, block); recvHandler != nil {
			break
		}
	}
	selected, received = chanrecv(c, ep, block)

	if recvHandler != nil {
		recvHandler.Ret(selected, received)
	}
	return
}

type hchan struct {
	qcount   uint           // total data in the queue
	dataqsiz uint           // size of the circular queue
	buf      unsafe.Pointer // points to an array of dataqsiz elements
	elemsize uint16
	closed   uint32
	elemtype uintptr // element type
}
