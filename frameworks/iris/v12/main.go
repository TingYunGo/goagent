// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64
// +build cgo

package irisframe

import (
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
	StorageIndexIris = tingyun3.StorageIndexIris
)

type controllerInfo struct {
	name   string
	method string
}

//go:noinline
func irisCreateRoutes(api *router.APIBuilder, methods []string, relativePath string, handlers ...context.Handler) []*router.Route {
	trampoline.arg1 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}
func wrapHandler(handler context.Handler, path string) context.Handler {
	info := controllerInfo{}
	if i := tingyun3.LocalGet(StorageIndexIris); i != nil {
		controller := i.(controllerInfo)
		info.name, info.method = controller.name, controller.method
	}
	if len(info.method) == 0 {
		handlerPC := reflect.ValueOf(handler).Pointer()
		info.method = runtime.FuncForPC(handlerPC).Name()
	}
	token := "github.com/TingYunGo/goagent"
	if tystring.SubString(info.method, 0, len(token)) == token {
		return handler
	}
	return func(ctx context.Context) {
		if tingyun3.LocalGet(StorageIndexIris) == nil {
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
			tingyun3.LocalSet(StorageIndexIris, 1)
			defer tingyun3.LocalDelete(StorageIndexIris)
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
	trampoline.arg2 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapirishandleMany(c *mvc.ControllerActivator, method, path, funcName string, override bool, middleware ...context.Handler) []*router.Route {
	info := controllerInfo{
		name:   c.Name(),
		method: funcName,
	}
	tingyun3.LocalSet(StorageIndexIris, info)
	defer tingyun3.LocalDelete(StorageIndexIris)
	return irishandleMany(c, method, path, funcName, override, middleware...)
}

//go:noinline
func routerFileServer(directory string, opts ...router.DirOptions) context.Handler {
	trampoline.arg3 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WraprouterFileServer(directory string, opts ...router.DirOptions) context.Handler {
	handler := routerFileServer(directory, opts...)
	return func(ctx context.Context) {
		if tingyun3.LocalGet(StorageIndexIris) == nil {
			action := tingyun3.GetAction()
			if action != nil {
				action.SetName("URI", ctx.Request().URL.Path)
			}
			tingyun3.LocalSet(StorageIndexIris, 1)
			defer tingyun3.LocalDelete(StorageIndexIris)
		}
		handler(ctx)
	}
}

//github.com/kataras/neffos.(*Conn).handleMessage
//go:noinline
func neffosConnhandleMessage(c *neffos.Conn, msg neffos.Message) error {
	trampoline.arg4 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
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
func websocketUpgrade(ctx context.Context, idGen websocket.IDGenerator, s *neffos.Server) *neffos.Conn {
	trampoline.arg5 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapwebsocketUpgrade(ctx context.Context, idGen websocket.IDGenerator, s *neffos.Server) *neffos.Conn {

	action := tingyun3.GetAction()
	if action != nil {
		action.SetName("websocket", "Upgrade")
	}
	r := websocketUpgrade(ctx, idGen, s)
	return r
}

//go:noinline
func neffosmakeEventFromMethod(v reflect.Value, method reflect.Method, eventMatcher neffos.EventMatcherFunc) (eventName string, cb neffos.MessageHandlerFunc) {
	trampoline.arg6 = trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
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
	tingyun3.Register(reflect.ValueOf(WraprouterFileServer).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapwebsocketUpgrade).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapirishandleMany).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapneffosConnhandleMessage).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapneffosmakeEventFromMethod).Pointer())
	tingyun3.Register(reflect.ValueOf(initTrampoline).Pointer())
}
