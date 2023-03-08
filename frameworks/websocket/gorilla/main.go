// Copyright 2023 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64 arm64
// +build cgo

package gorillaframe

import (
	"io"
	"reflect"
	"strconv"
	"time"

	"github.com/gorilla/websocket"

	"github.com/TingYunGo/goagent"
	"github.com/TingYunGo/goagent/libs/tystring"
)

const (
	StorageGorillaWS = tingyun3.StorageGorillaWS
)

func matchMethod(method, matcher string) bool {
	return tystring.SubString(method, 0, len(matcher)) == matcher
}
func isNativeMethod(method string) bool {

	if matchMethod(method, "io.") {
		return true
	}
	if matchMethod(method, "io/") {
		return true
	}
	if matchMethod(method, "github.com/gorilla/websocket") {
		return true
	}
	return false
}

//go:noinline
func getCallName(skip int) (callerName string) {
	skip++
	callerName = tingyun3.GetCallerName(skip)
	for isNativeMethod(callerName) {
		skip++
		callerName = tingyun3.GetCallerName(skip)
	}
	return
}

func onWebsocketTaskEnd() {
	action := tingyun3.GetAction()
	if action == nil {
		return
	}
	if component := tingyun3.GetComponent(); component != nil {
		component.Finish()
	}
	igonre_duration := tingyun3.ReadLocalConfigInteger(tingyun3.ConfigLocalIntegerWebsocketIgnore, 0)
	if action.Duration() < time.Duration(igonre_duration)*time.Millisecond {
		action.Ignore()
	} else {
		action.SetStatusCode(1)
		action.Finish()
	}
	tingyun3.LocalClear()
}
func getWebsockType(messageType int) string {
	if messageType == 1 {
		return "websocket.text"
	}
	if messageType == 2 {
		return "websocket.bin"
	}
	return "websocket." + strconv.Itoa(messageType)
}
func (r *messageReader) onWebsocketTaskBegin() {
	methodName := getCallName(3)
	action, _ := tingyun3.CreateAction(getWebsockType(r.messageType), methodName)
	if action == nil {
		return
	}
	action.SetName("CLIENTIP", r.addr)
	if len(r.url) > 0 {
		action.SetURL(r.url)
	}
	component := action.CreateComponent(methodName)
	tingyun3.SetAction(action)
	tingyun3.SetComponent(component)
}

type messageReader struct {
	r           io.Reader
	size        int64
	messageType int
	addr        string
	url         string
}

func (r *messageReader) Read(b []byte) (int, error) {
	size, err := r.r.Read(b)
	if size > 0 {
		r.size += int64(size)
	}
	if err == io.EOF {
		r.onWebsocketTaskBegin()
	}
	return size, err
}

func (r *messageReader) Close() error {
	return r.r.(io.Closer).Close()
}

//go:noinline
func ConnNextReader(c *websocket.Conn) (messageType int, r io.Reader, err error) {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return -1, nil, nil
}

//go:noinline
func WrapConnNextReader(c *websocket.Conn) (int, io.Reader, error) {
	onWebsocketTaskEnd()
	msgType, r, err := ConnNextReader(c)

	if r != nil && msgType != 8 {
		remoteIP := ""

		tystring.SplitMapString(c.RemoteAddr().String(), func(t byte) bool {
			return t == ':'
		}, func(ip, _ string) {
			remoteIP = ip
		})
		url := ""

		raw := c.UnderlyingConn()
		if wrpper_conn, ok := raw.(*tingyun3.NetConnWrapper); ok {
			url = wrpper_conn.Url()
			if len(url) > 0 {
				url = "ws" + url
			}
		}
		r = &messageReader{r, 0, msgType, remoteIP, url}
	}
	return msgType, r, err
}

func init() {
	tingyun3.Register(reflect.ValueOf(WrapConnNextReader).Pointer())
	tingyun3.Register(reflect.ValueOf(initTrampoline).Pointer())
}
