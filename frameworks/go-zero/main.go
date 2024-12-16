//go:build linux && (amd64 || arm64) && cgo
// +build linux
// +build amd64 arm64
// +build cgo

// Copyright 2024 冯立强 fenglq@tingyun.com.  All rights reserved.
package gozeroframe

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/TingYunGo/goagent"
	"github.com/zeromicro/go-zero/rest"
)

//go:noinline
func _rest_Server_AddRoutes(s *rest.Server, rs []rest.Route, opts ...rest.RouteOption) {
	trampoline.arg2 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
}

func wrapRouteHandle(http_method, name string, h http.HandlerFunc) http.HandlerFunc {

	handlerPointer := runtime.FuncForPC(reflect.ValueOf(h).Pointer())
	if handlerPointer != nil {
		name = handlerPointer.Name()
	}
	handler := func(w http.ResponseWriter, request *http.Request) {
		action := tingyun3.GetAction()
		timeoutHander := false
		if action == nil {
			action, _ = tingyun3.FindAction(request.Context())
			if action != nil {
				timeoutHander = true
				tingyun3.SetAction(action)
			}
		}
		if action != nil && tingyun3.ReadServerConfigBool(tingyun3.ServerConfigBoolAutoActionNaming, true) {
			action.SetName("Route", name)
			action.SetHTTPMethod(http_method)
		}
		defer func() {
			if timeoutHander {
				tingyun3.LocalClear()
			}
			exception := recover()
			if exception != nil && action != nil {
				action.SetError(exception)
			}
			if exception != nil {
				action.Finish()
				panic(exception)
			}
		}()
		h(w, request)
	}
	return handler
}

//go:noinline
func Wrap_rest_Server_AddRoutes(s *rest.Server, rs []rest.Route, opts ...rest.RouteOption) {

	for id := 0; id < len(rs); id++ {
		h := rs[id].Handler
		rs[id].Handler = wrapRouteHandle(rs[id].Method, rs[id].Path, h)
	}
	_rest_Server_AddRoutes(s, rs, opts...)
}

func init() {
	tingyun3.Register(reflect.ValueOf(Wrap_rest_Server_AddRoutes).Pointer())
}
