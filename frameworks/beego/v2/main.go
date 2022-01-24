// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64
// +build cgo

package beegoframe

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/beego/beego/v2/server/web/context"

	"github.com/TingYunGo/goagent"
	beego "github.com/beego/beego/v2/server/web"
	param "github.com/beego/beego/v2/server/web/context/param"
)

const (
	StorageIndexBeego = tingyun3.StorageIndexBeego
)

type handlerInfo struct {
	name   string
	method string
	isFunc bool
}

//go:noinline
func beegoAddMethod(p *beego.ControllerRegister, method, pattern string, f beego.FilterFunc) {
	trampoline.arg1 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
}

var methodMap = map[string]string{
	"GET":     "Get",
	"POST":    "Post",
	"PUT":     "Put",
	"DELETE":  "Delete",
	"PATCH":   "Patch",
	"OPTIONS": "Options",
	"HEAD":    "Head",
	"TRACE":   "Trace",
}

//go:noinline
func WrapbeegoAddMethod(p *beego.ControllerRegister, method, pattern string, f beego.FilterFunc) {
	handlerPC := reflect.ValueOf(f).Pointer()
	methodName := runtime.FuncForPC(handlerPC).Name()
	info := handlerInfo{
		name:   methodName,
		isFunc: true,
	}
	pre := tingyun3.LocalGet(StorageIndexBeego)
	tingyun3.LocalSet(StorageIndexBeego, info)

	defer func() {
		if pre == nil {
			tingyun3.LocalDelete(StorageIndexBeego)
		} else {
			tingyun3.LocalSet(StorageIndexBeego, pre)
		}
	}()
	beegoAddMethod(p, method, pattern, f)
}

//go:noinline
func beegoHandler(p *beego.ControllerRegister, pattern string, h http.Handler, options ...interface{}) {
	trampoline.arg1 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
}

//go:noinline
func WrapbeegoHandler(p *beego.ControllerRegister, pattern string, h http.Handler, options ...interface{}) {
	handlerPC := reflect.ValueOf(h).Pointer()
	methodName := runtime.FuncForPC(handlerPC).Name()
	info := handlerInfo{
		name:   methodName,
		isFunc: true,
	}
	if len(methodName) == 0 {
		reflectVal := reflect.ValueOf(h)
		info.name = reflect.Indirect(reflectVal).Type().String()
		info.method = "ServeHTTP"
		info.isFunc = false
	}
	pre := tingyun3.LocalGet(StorageIndexBeego)
	tingyun3.LocalSet(StorageIndexBeego, info)
	defer func() {
		if pre == nil {
			tingyun3.LocalDelete(StorageIndexBeego)
		} else {
			tingyun3.LocalSet(StorageIndexBeego, pre)
		}
	}()
	beegoHandler(p, pattern, h, options...)
}

//go:noinline
func beegoaddWithMethodParams(p *beego.ControllerRegister, pattern string, c beego.ControllerInterface, methodParams []*param.MethodParam, mappingMethods ...string) {
	trampoline.arg1 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
}

//go:noinline
func WrapbeegoaddWithMethodParams(p *beego.ControllerRegister, pattern string, c beego.ControllerInterface, methodParams []*param.MethodParam, mappingMethods ...string) {
	reflectVal := reflect.ValueOf(c)
	info := handlerInfo{
		name:   reflect.Indirect(reflectVal).Type().String(),
		isFunc: false,
	}
	pre := tingyun3.LocalGet(StorageIndexBeego)
	tingyun3.LocalSet(StorageIndexBeego, info)
	defer func() {
		if pre == nil {
			tingyun3.LocalDelete(StorageIndexBeego)
		} else {
			tingyun3.LocalSet(StorageIndexBeego, pre)
		}
	}()

	beegoaddWithMethodParams(p, pattern, c, methodParams, mappingMethods...)
}

//go:noinline
func beegoAddAutoPrefix(p *beego.ControllerRegister, prefix string, c beego.ControllerInterface) {
	trampoline.arg1 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
}

//go:noinline
func WrapbeegoAddAutoPrefix(p *beego.ControllerRegister, prefix string, c beego.ControllerInterface) {
	reflectVal := reflect.ValueOf(c)
	info := handlerInfo{
		name:   reflect.Indirect(reflectVal).Type().String(),
		isFunc: false,
	}
	pre := tingyun3.LocalGet(StorageIndexBeego)
	defer func() {
		if pre == nil {
			tingyun3.LocalDelete(StorageIndexBeego)
		} else {
			tingyun3.LocalSet(StorageIndexBeego, pre)
		}
	}()
	tingyun3.LocalSet(StorageIndexBeego, info)
	beegoAddAutoPrefix(p, prefix, c)
}

var routeMap map[string]bool = nil

//go:noinline
func beegoaddToRouter(p *beego.ControllerRegister, method, pattern string, r *beego.ControllerInfo) {
	trampoline.arg1 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
}

//go:noinline
func WrapbeegoaddToRouter(p *beego.ControllerRegister, method, pattern string, r *beego.ControllerInfo) {
	if routeMap == nil {
		routeMap = make(map[string]bool)
	}
	info := tingyun3.LocalGet(StorageIndexBeego)
	if _, found := routeMap[pattern]; !found {
		routeMap[pattern] = true
		handler := handlerInfo{}
		if info != nil {
			handler = info.(handlerInfo)
		}
		p.InsertFilter(pattern, beego.BeforeExec, func(ctx *context.Context) {
			action := tingyun3.GetAction()
			if action != nil {
				if len(handler.name) > 0 {
					if handler.isFunc {
						action.SetName("ROUTE", handler.name)
					} else {
						if len(handler.method) > 0 {
							action.SetName("Method", handler.method)
						} else if method, found := methodMap[action.GetMethod()]; found {
							action.SetName("Method", method)
						}
						action.SetName("CLASS", handler.name)
					}
				}
				action.SetName("URI", pattern)
			}
		})
	}
	beegoaddToRouter(p, method, pattern, r)
}

func init() {
	tingyun3.Register(reflect.ValueOf(WrapbeegoaddToRouter).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapbeegoAddMethod).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapbeegoHandler).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapbeegoAddAutoPrefix).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapbeegoaddWithMethodParams).Pointer())
}
