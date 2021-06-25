// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64

package irisframe

import (
	"fmt"
	"reflect"
	"runtime"
	"time"

	"github.com/TingYunGo/goagent/libs/tystring"

	"github.com/TingYunGo/goagent"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/mvc"
)

const (
	irisRoutineLocalIndex = 9 + 8*4
)

type controllerInfo struct {
	name   string
	method string
}

//go:noinline
func irisCreateRoutes(api *router.APIBuilder, methods []string, relativePath string, handlers ...context.Handler) []*router.Route {
	fmt.Println(api, methods, relativePath, handlers)
	return nil
}
func wrapHandler(handler context.Handler, path string) context.Handler {
	info := controllerInfo{}
	if i := tingyun3.LocalGet(irisRoutineLocalIndex); i != nil {
		controller := i.(controllerInfo)
		info.name, info.method = controller.name, controller.method
	}
	if len(info.method) == 0 {
		handlerPC := reflect.ValueOf(handler).Pointer()
		info.method = runtime.FuncForPC(handlerPC).Name()
	}
	token := "git.codemonky.net/TingYunGo/goagent"
	if tystring.SubString(info.method, 0, len(token)) == token {
		return handler
	}
	return func(ctx context.Context) {
		if tingyun3.LocalGet(irisRoutineLocalIndex) == nil {
			action := tingyun3.GetAction()
			if action != nil {
				if len(info.name) == 0 && len(info.method) > 0 {
					action.SetName("ROUTE", info.name)
				} else if len(info.name) > 0 {
					action.SetName("Method", info.method)
					action.SetName("CLASS", info.name)
				}
				if len(info.name) == 0 && len(info.method) == 0 {
					action.SetName("URI", path)
				}
			}
			tingyun3.LocalSet(irisRoutineLocalIndex, 1)
			defer tingyun3.LocalDelete(irisRoutineLocalIndex)
		}
		handler(ctx)
	}
}

//go:noinline
func WrapirisCreateRoutes(api *router.APIBuilder, methods []string, relativePath string, handlers ...context.Handler) []*router.Route {
	for i := 0; i < len(handlers); i++ {
		handlers[i] = wrapHandler(handlers[i], api.GetRelPath()+relativePath)
	}
	return irisCreateRoutes(api, methods, relativePath, handlers...)
}

//go:noinline
func irishandleMany(c *mvc.ControllerActivator, method, path, funcName string, override bool, middleware ...context.Handler) []*router.Route {
	fmt.Println(c, method, path, funcName, override, middleware)
	return nil
}

//go:noinline
func WrapirishandleMany(c *mvc.ControllerActivator, method, path, funcName string, override bool, middleware ...context.Handler) []*router.Route {
	info := controllerInfo{
		name:   c.Name(),
		method: funcName,
	}
	tingyun3.LocalSet(irisRoutineLocalIndex, info)
	defer tingyun3.LocalDelete(irisRoutineLocalIndex)
	return irishandleMany(c, method, path, funcName, override, middleware...)
}

func init() {
	tingyun3.Register(reflect.ValueOf(WrapirisCreateRoutes).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapirishandleMany).Pointer())
}
