// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64

package beegoframe

import (
	"net/http"
	"reflect"
	"runtime"

	"github.com/astaxie/beego/context"

	"github.com/TingYunGo/goagent"
	beego "github.com/astaxie/beego"
	param "github.com/astaxie/beego/context/param"
)

const (
	beegoRoutineLocalIndex = 9 + 8 + 8 + 8
)

type handlerInfo struct {
	name   string
	method string
	isFunc bool
}

var tempVar = 0x1234567890

//go:noinline
func beegoAddMethod(p *beego.ControllerRegister, method, pattern string, f beego.FilterFunc) {
	tempVar += 10
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
	tingyun3.LocalSet(beegoRoutineLocalIndex, info)
	defer tingyun3.LocalDelete(beegoRoutineLocalIndex)
	beegoAddMethod(p, method, pattern, f)
}

//go:noinline
func beegoHandler(p *beego.ControllerRegister, pattern string, h http.Handler, options ...interface{}) {
	tempVar += 10
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
	tingyun3.LocalSet(beegoRoutineLocalIndex, info)
	defer tingyun3.LocalDelete(beegoRoutineLocalIndex)
	beegoHandler(p, pattern, h, options...)
}

//go:noinline
func beegoaddWithMethodParams(p *beego.ControllerRegister, pattern string, c beego.ControllerInterface, methodParams []*param.MethodParam, mappingMethods ...string) {
	tempVar += 10
}

//go:noinline
func WrapbeegoaddWithMethodParams(p *beego.ControllerRegister, pattern string, c beego.ControllerInterface, methodParams []*param.MethodParam, mappingMethods ...string) {
	reflectVal := reflect.ValueOf(c)
	info := handlerInfo{
		name:   reflect.Indirect(reflectVal).Type().String(),
		isFunc: false,
	}
	tingyun3.LocalSet(beegoRoutineLocalIndex, info)
	defer tingyun3.LocalDelete(beegoRoutineLocalIndex)
	beegoaddWithMethodParams(p, pattern, c, methodParams, mappingMethods...)
}

//go:noinline
func beegoAddAutoPrefix(p *beego.ControllerRegister, prefix string, c beego.ControllerInterface) {
	tempVar += 10
}

//go:noinline
func WrapbeegoAddAutoPrefix(p *beego.ControllerRegister, prefix string, c beego.ControllerInterface) {
	reflectVal := reflect.ValueOf(c)
	info := handlerInfo{
		name:   reflect.Indirect(reflectVal).Type().String(),
		isFunc: false,
	}
	tingyun3.LocalSet(beegoRoutineLocalIndex, info)
	defer tingyun3.LocalDelete(beegoRoutineLocalIndex)
	beegoAddAutoPrefix(p, prefix, c)
}

var routeMap map[string]bool = nil

//go:noinline
func beegoaddToRouter(p *beego.ControllerRegister, method, pattern string, r *beego.ControllerInfo) {
	tempVar += 10
}

//go:noinline
func WrapbeegoaddToRouter(p *beego.ControllerRegister, method, pattern string, r *beego.ControllerInfo) {
	if routeMap == nil {
		routeMap = make(map[string]bool)
	}
	info := tingyun3.LocalGet(beegoRoutineLocalIndex)
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
