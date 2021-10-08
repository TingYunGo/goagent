// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64

package irisframe

import (
	"fmt"
	"net"
	"reflect"
	"runtime"

	"github.com/TingYunGo/goagent/libs/tystring"

	"github.com/TingYunGo/goagent"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/router"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/websocket"
	"github.com/kataras/neffos"
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
	return func(ctx *context.Context) {
		if tingyun3.LocalGet(irisRoutineLocalIndex) == nil {
			action := tingyun3.GetAction()
			if action != nil {
				if len(info.name) == 0 && len(info.method) > 0 {
					action.SetName("ROUTE", info.method)
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

//go:noinline
func irisHandleDir(api *router.APIBuilder, requestPath string, fsOrDir interface{}, opts ...router.DirOptions) (routes []*router.Route) {
	fmt.Println(api, requestPath, fsOrDir, opts)
	return nil
}

//go:noinline
func WrapirisHandleDir(api *router.APIBuilder, requestPath string, fsOrDir interface{}, opts ...router.DirOptions) (routes []*router.Route) {
	return irisHandleDir(api, requestPath, fsOrDir, opts...)
}

//github.com/kataras/neffos.(*Conn).handleMessage
//go:noinline
func neffosConnhandleMessage(c *neffos.Conn, msg neffos.Message) error {
	fmt.Println(c, msg)
	return nil
}

func readBoolean(name string, defaultValue bool) bool {
	retvalue := defaultValue
	if value, found := tingyun3.ConfigRead(name); found {
		if v, ok := value.(bool); ok {
			retvalue = v
		}
	}
	return retvalue
}

//go:noinline
func WrapneffosConnhandleMessage(c *neffos.Conn, msg neffos.Message) error {
	action := tingyun3.GetAction()
	preaction := action
	if preaction == nil {
		r := c.Socket().Request()
		if readBoolean("websocket_enabled", false) {
			if action, _ = tingyun3.CreateAction("URI", r.URL.Path); action != nil {
				tingyun3.SetAction(action)
				if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
					action.SetName("CLIENTIP", host)
				}
			}
		}
	}

	defer func() {
		if preaction == nil && action != nil {
			action.Finish()
			tingyun3.LocalClear()
		}
	}()
	return neffosConnhandleMessage(c, msg)
}

//go:noinline
func websocketUpgrade(ctx *context.Context, idGen websocket.IDGenerator, s *neffos.Server) *neffos.Conn {
	fmt.Println(ctx, idGen, s)
	return nil
}

//go:noinline
func WrapwebsocketUpgrade(ctx *context.Context, idGen websocket.IDGenerator, s *neffos.Server) *neffos.Conn {

	action := tingyun3.GetAction()
	if action != nil {
		action.SetName("websocket", "Upgrade")
	}
	r := websocketUpgrade(ctx, idGen, s)
	return r
}

//go:noinline
func neffosfireEvent(e neffos.Events, c *neffos.NSConn, msg neffos.Message) error {
	fmt.Println(e, c, msg)
	return nil
}

//go:noinline
func WrapneffosfireEvent(e neffos.Events, c *neffos.NSConn, msg neffos.Message) error {
	action := tingyun3.GetAction()
	if action != nil {
		if h, found := e[msg.Event]; found {
			handlerPC := reflect.ValueOf(h).Pointer()
			methodName := runtime.FuncForPC(handlerPC).Name()
			if len(methodName) > 0 {
				action.SetName("Method", methodName)
				action.SetName("Websocket", methodName)
			}
		}
	}
	return neffosfireEvent(e, c, msg)
}

//go:noinline
func neffosmakeEventFromMethod(v reflect.Value, method reflect.Method, eventMatcher neffos.EventMatcherFunc) (eventName string, cb neffos.MessageHandlerFunc) {
	fmt.Println(v, method, eventMatcher)
	return "", nil
}

//go:noinline
func WrapneffosmakeEventFromMethod(v reflect.Value, method reflect.Method, eventMatcher neffos.EventMatcherFunc) (eventName string, cb neffos.MessageHandlerFunc) {
	name, handler := neffosmakeEventFromMethod(v, method, eventMatcher)
	if handler == nil {
		return name, handler
	}
	className := v.Type().String()
	methodName := method.Name
	if tystring.SubString(className, 0, 1) == "*" {
		className = tystring.SubString(className, 1, len(className))
	}
	eventName = name
	cb = func(nsconn *neffos.NSConn, msg neffos.Message) error {
		action := tingyun3.GetAction()
		if action != nil {
			if len(name) > 0 {
				action.SetName("Method", methodName)
			}
			if len(className) > 0 {
				action.SetName("CLASS", className)
			}
		}
		return handler(nsconn, msg)
	}
	return
}

func init() {
	tingyun3.Register(reflect.ValueOf(WrapirisCreateRoutes).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapwebsocketUpgrade).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapneffosfireEvent).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapirishandleMany).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapneffosConnhandleMessage).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapneffosmakeEventFromMethod).Pointer())
}
