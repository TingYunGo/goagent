// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64

package echoframe

import (
	"reflect"
	"runtime"

	"github.com/TingYunGo/goagent"
	"github.com/labstack/echo/v4"
)

func getHandlerName(handler echo.HandlerFunc) string {
	handlerPC := reflect.ValueOf(handler).Pointer()
	return runtime.FuncForPC(handlerPC).Name()
}

func wrapHandler(method, route string, handler echo.HandlerFunc) echo.HandlerFunc {
	methodName := getHandlerName(handler)
	httpMethod := method
	routePath := route
	return func(ctx echo.Context) error {
		action := tingyun3.GetAction()
		if action != nil {
			action.SetName("Route", methodName)
			action.SetName("URI", routePath)
			action.SetHTTPMethod(httpMethod)
		}
		return handler(ctx)
	}
}

//go:noinline
func echoEchoAdd(ptr *echo.Echo, method, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapechoEchoAdd(ptr *echo.Echo, method, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route {
	tingyun3.LocalSet(9+8, "handled")
	wrapper := wrapHandler(method, path, handler)
	r := echoEchoAdd(ptr, method, path, wrapper, middleware...)
	tingyun3.LocalDelete(9 + 8)
	return r
}

//go:noinline
func echoEchoadd(ptr *echo.Echo, host, method, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route {
	trampoline.arg2 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapechoEchoadd(ptr *echo.Echo, host, method, path string, handler echo.HandlerFunc, middleware ...echo.MiddlewareFunc) *echo.Route {
	tingyun3.LocalSet(9+8, "handled")
	wrapper := wrapHandler(method, path, handler)
	r := echoEchoadd(ptr, host, method, path, wrapper, middleware...)
	tingyun3.LocalDelete(9 + 8)
	return r
}

//go:noinline
func echoRouterAdd(ptr *echo.Router, method, path string, h echo.HandlerFunc) {
	trampoline.arg3 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
}

//go:noinline
func WrapechoRouterAdd(ptr *echo.Router, method, path string, h echo.HandlerFunc) {
	if tingyun3.LocalGet(9+8) == nil && h != nil {
		h = wrapHandler(method, path, h)
	}
	echoRouterAdd(ptr, method, path, h)
}

func init() {
	tingyun3.Register(reflect.ValueOf(WrapechoEchoAdd).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapechoEchoadd).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapechoRouterAdd).Pointer())
	tingyun3.Register(reflect.ValueOf(initTrampoline).Pointer())
}
