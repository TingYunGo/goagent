// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64 arm64
// +build cgo

package kratosframe

import (
	"context"
	"net/http"
	"reflect"
	"runtime"

	"github.com/TingYunGo/goagent"
	_ "github.com/TingYunGo/goagent/frameworks/grpc"
	"github.com/go-kratos/kratos/v2/middleware"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
)

type responseWriter struct {
	code int
	w    http.ResponseWriter
}
type wrapper struct {
	router *khttp.Router
	req    *http.Request
	res    http.ResponseWriter
	w      responseWriter
}

//go:noinline
func TransportHttpContextMiddleware(c *wrapper, h middleware.Handler) middleware.Handler {
	trampoline.arg2 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapTransportHttpContextMiddleware(c *wrapper, h middleware.Handler) middleware.Handler {

	r := TransportHttpContextMiddleware(c, h)
	if r == nil {

		return r
	}
	//callerName := tingyun3.GetCallerName(2)
	//className := reflect.TypeOf(r).String()
	tyWrapper := func(ctx context.Context, req interface{}) (interface{}, error) {

		action := tingyun3.GetAction()
		if action != nil {
			// fmt.Println("Action ", action.GetName())
			// action.SetName("Route", method)
		}

		v, e := r(ctx, req)

		return v, e
	}
	return tyWrapper
}


//go:noinline
func TransportHttpRouterHandle(r *khttp.Router, method, relativePath string, h khttp.HandlerFunc, filters ...khttp.FilterFunc) {
	trampoline.arg2 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
}


func wrapRouteHandle(name string, h khttp.HandlerFunc) khttp.HandlerFunc {

	handler := func(ctx khttp.Context) error {

		action := tingyun3.GetAction()
		if action != nil {
			if tingyun3.ReadServerConfigBool(tingyun3.ServerConfigBoolAutoActionNaming, true) {
				action.SetName("Route", name)
			}
		}

		e := h(ctx)
		if e != nil && action != nil {
			action.SetException(e)
		}
		return e
	}
	return handler
}

//go:noinline
func WrapTransportHttpRouterHandle(r *khttp.Router, method, relativePath string, h khttp.HandlerFunc, filters ...khttp.FilterFunc) {
	handlerPC := reflect.ValueOf(h).Pointer()
	methodName := runtime.FuncForPC(handlerPC).Name()

	TransportHttpRouterHandle(r, method, relativePath, wrapRouteHandle(methodName, h), filters...)
}

func init() {
	tingyun3.Register(reflect.ValueOf(WrapTransportHttpContextMiddleware).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapTransportHttpRouterHandle).Pointer())
	tingyun3.Register(reflect.ValueOf(initTrampoline).Pointer())
}
