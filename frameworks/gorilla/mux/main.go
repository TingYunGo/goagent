// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64 arm64
// +build cgo

package gorillamux

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/gorilla/mux"

	"github.com/TingYunGo/goagent"
)

func wrapRouteHandlerFuncHandle(name string, f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {

	handler := func(w http.ResponseWriter, r *http.Request) {
		action := tingyun3.GetAction()
		if action != nil {
			if tingyun3.ReadServerConfigBool(tingyun3.ServerConfigBoolAutoActionNaming, true) {
				action.SetName("Route", name)
			}
		}
		f(w, r)
	}
	return handler
}

//go:noinline
func RouterHandleFunc(r *mux.Router, path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapRouterHandleFunc(r *mux.Router, path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {

	handlerPC := reflect.ValueOf(f).Pointer()
	methodName := runtime.FuncForPC(handlerPC).Name()

	return RouterHandleFunc(r, path, wrapRouteHandlerFuncHandle(methodName, f))
}

func getHandlerName(handler http.Handler) string {
	var methodName string
	className := reflect.TypeOf(handler).String()
	if className == "http.HandlerFunc" || className == "HandlerFunc" {
		handlerPC := reflect.ValueOf(handler).Pointer()
		methodName = runtime.FuncForPC(handlerPC).Name()
	} else {
		if len(className) > 0 && className[0] == '*' {
			className = className[1:]
		}
		methodName = className + ".ServeHTTP"
	}
	return methodName
}

func wrapRouterHandleHandle(name string, handler http.Handler) http.Handler {

	newHandler := func(w http.ResponseWriter, r *http.Request) {
		action := tingyun3.GetAction()
		if action != nil {
			if tingyun3.ReadServerConfigBool(tingyun3.ServerConfigBoolAutoActionNaming, true) {
				action.SetName("Route", name)
			}
		}
		handler.ServeHTTP(w, r)
	}
	return http.HandlerFunc(newHandler)

}

//go:noinline
func RouterHandle(r *mux.Router, path string, handler http.Handler) *mux.Route {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapRouterHandle(r *mux.Router, path string, handler http.Handler) *mux.Route {

	methodName := getHandlerName(handler)
	return RouterHandle(r, path, wrapRouterHandleHandle(methodName, handler))
}

func init() {
	tingyun3.Register(reflect.ValueOf(WrapRouterHandleFunc).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapRouterHandle).Pointer())
	tingyun3.Register(reflect.ValueOf(initTrampoline).Pointer())
}
