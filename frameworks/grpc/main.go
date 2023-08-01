// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64 arm64
// +build cgo

package grpcframe

import (
	"context"
	"io"
	"net"
	"reflect"
	"runtime"

	"github.com/TingYunGo/goagent"
	"github.com/TingYunGo/goagent/libs/tystring"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

//go:noinline
func ServerRegisterService(s *grpc.Server, sd *grpc.ServiceDesc, ss interface{}) {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
}

type methodHandler func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error)
type StreamHandler func(srv interface{}, stream grpc.ServerStream) error

func getMetadata(md metadata.MD, name string) string {
	if md != nil {
		if values := md.Get(name); len(values) > 0 {
			return values[0]
		}
	}
	return ""
}
func parseIP(addr string) string {
	for id := len(addr); id > 0; id-- {
		if addr[id-1] == ':' {
			addr = tystring.SubString(addr, 0, id-1)
			break
		}
	}
	return addr
}
func wrapMethodsHandler(serviceName, className, methodName string, handler func(interface{}, context.Context, func(interface{}) error, grpc.UnaryServerInterceptor) (interface{}, error)) func(interface{}, context.Context, func(interface{}) error, grpc.UnaryServerInterceptor) (interface{}, error) {
	wrapper := func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {

		md, _ := metadata.FromIncomingContext(ctx)

		action, _ := tingyun3.CreateAction(className, methodName)

		if action != nil {

			if trackID := getMetadata(md, "X-Tingyun"); len(trackID) > 0 {

				action.SetTrackID(trackID)
				action.SetBackEnabled(false)
			}

			if p, ok := peer.FromContext(ctx); ok {
				action.SetName("CLIENTIP", parseIP(p.Addr.String()))
			}

			catchRequestHeaders := tingyun3.CatchRequestHeaders()
			for _, item := range catchRequestHeaders {
				if value := getMetadata(md, item); len(value) > 0 {
					action.AddRequestParam(item, value)
				}
			}
			reqUrl := "grpc://" + serviceName + "/" + className + "/" + methodName
			action.SetURL(reqUrl)

			tingyun3.SetAction(action)

			ctx = context.WithValue(ctx, "TingYunWebAction", action)
		}
		defer func() {
			a := tingyun3.GetAction()
			if a != nil && a != action {
				a.Finish()
			}
			tingyun3.LocalClear()
			exception := recover()
			if exception != nil {
				action.SetError(exception)
			}
			action.Finish()
			//re throw
			if exception != nil {
				panic(exception)
			}
		}()

		return handler(srv, ctx, dec, interceptor)
	}
	return wrapper
}

func wrapStreamsHandler(serviceName, className, methodName string, handler grpc.StreamHandler) grpc.StreamHandler {
	wrapper := func(srv interface{}, stream grpc.ServerStream) error {

		ctx := stream.Context()
		md, _ := metadata.FromIncomingContext(ctx)

		action, _ := tingyun3.CreateAction(className, methodName)

		if action != nil {

			if trackID := getMetadata(md, "X-Tingyun"); len(trackID) > 0 {

				action.SetTrackID(trackID)
				action.SetBackEnabled(false)
			}

			if p, ok := peer.FromContext(ctx); ok {
				action.SetName("CLIENTIP", parseIP(p.Addr.String()))
			}

			catchRequestHeaders := tingyun3.CatchRequestHeaders()
			for _, item := range catchRequestHeaders {
				if value := getMetadata(md, item); len(value) > 0 {
					action.AddRequestParam(item, value)
				}
			}
			reqUrl := "grpc://" + serviceName + "/" + className + "/" + methodName
			action.SetURL(reqUrl)
			tingyun3.SetAction(action)
		}
		defer func() {
			a := tingyun3.GetAction()
			if a != nil && a != action {
				a.Finish()
			}
			tingyun3.LocalClear()
			exception := recover()
			if exception != nil {
				action.SetError(exception)
			}
			action.Finish()
			//re throw
			if exception != nil {
				panic(exception)
			}
		}()

		return handler(srv, stream)
	}
	return wrapper
}

//go:noinline
func WrapServerRegisterService(s *grpc.Server, sd *grpc.ServiceDesc, ss interface{}) {

	className := getHandlerName(ss)

	for id, _ := range sd.Methods {
		sd.Methods[id].Handler = wrapMethodsHandler(sd.ServiceName, className, sd.Methods[id].MethodName, sd.Methods[id].Handler)
	}

	for id, _ := range sd.Streams {
		sd.Streams[id].Handler = wrapStreamsHandler(sd.ServiceName, className, sd.Streams[id].StreamName, sd.Streams[id].Handler)
	}
	ServerRegisterService(s, sd, ss)
}

//go:noinline
func ServerServe(s *grpc.Server, lis net.Listener) error {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapServerServe(s *grpc.Server, lis net.Listener) error {
	addr := lis.Addr().String()
	if len(addr) > 0 {
		tingyun3.AppendListenAddress(addr)
	}
	return ServerServe(s, lis)
}

func getHandlerName(ss interface{}) string {
	className := reflect.TypeOf(ss).String()
	if len(className) > 0 && className[0] == '*' {
		className = className[1:]
	}
	return className
}

func matchMethod(method, matcher string) bool {
	return tystring.SubString(method, 0, len(matcher)) == matcher
}
func isNativeMethod(method string) bool {

	if matchMethod(method, "github.com/TingYunGo/goagent") {
		return true
	}
	return false
}

//go:noinline
func getCallName(skip int) (callerName string) {
	stackCount := skip + 3
	skip++
	stackList := make([]uintptr, stackCount)
	count := runtime.Callers(skip, stackList)

	for i := 0; i < count; i++ {

		callerName = runtime.FuncForPC(stackList[i]).Name()
		if !isNativeMethod(callerName) {
			break
		}
	}
	return
}

//go:noinline
func ClientConnInvoke(cc *grpc.ClientConn, ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil
}

//go:noinline
func WrapClientConnInvoke(cc *grpc.ClientConn, ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error {

	var component *tingyun3.Component = (*tingyun3.Component)(nil)
	callerName := ""
	url := ""
	target := ""

	if action, _ := tingyun3.FindAction(ctx); action != nil {
		callerName = getCallName(2)
		target = cc.Target()
		url = "grpc://" + target
		if url[len(url)-1] != '/' {
			url += "/"
		}
		if len(method) > 0 {
			uri := method
			if uri[0] == '/' {
				uri = uri[1:]
			}
			url += uri
		}
		component = action.CreateExternalComponent(url, callerName)
	}

	if trackID := component.CreateTrackID(); len(trackID) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("X-Tingyun", trackID))
	}

	if component == nil {
		return ClientConnInvoke(cc, ctx, method, args, reply, opts...)
	}
	if len(target) > 0 {
		component.SetURL(url)
	}

	defer func() {
		if exception := recover(); exception != nil {
			component.SetError(exception, callerName+"("+method+")", 3)
			component.Finish()
			panic(exception)
		}
	}()

	ret := ClientConnInvoke(cc, ctx, method, args, reply, opts...)

	if ret != nil {
		component.SetException(ret, callerName+"("+method+")", 3)
	}
	component.Finish()

	return ret
}

//go:noinline
func ClientConnNewStream(cc *grpc.ClientConn, ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
	return nil, nil
}

//go:noinline
func WrapClientConnNewStream(cc *grpc.ClientConn, ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {

	var component *tingyun3.Component = (*tingyun3.Component)(nil)
	callerName := ""
	url := ""
	target := ""

	if action, _ := tingyun3.FindAction(ctx); action != nil {
		callerName = getCallName(2)
		target = cc.Target()
		url = "grpc://" + target
		if url[len(url)-1] != '/' {
			url += "/"
		}
		if len(method) > 0 {
			uri := method
			if uri[0] == '/' {
				uri = uri[1:]
			}
			url += uri
		}
		component = action.CreateExternalComponent(url, callerName)
	}

	if component == nil {
		return ClientConnNewStream(cc, ctx, desc, method, opts...)
	}

	if trackID := component.CreateTrackID(); len(trackID) > 0 {
		ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("X-Tingyun", trackID))
	}

	if len(target) > 0 {
		component.SetURL(url)
	}

	defer func() {
		if exception := recover(); exception != nil {
			component.SetError(exception, callerName+"("+method+")", 3)
			component.Finish()
			panic(exception)
		}
	}()

	cs, ret := ClientConnNewStream(cc, ctx, desc, method, opts...)

	if ret != nil {
		component.SetException(ret, callerName+"("+method+")", 3)
	}
	if cs == nil {
		component.Finish()
	} else {
		cs = &ClientStreamWrapper{
			component,
			cs,
		}
	}

	return cs, ret
}

type ClientStreamWrapper struct {
	component *tingyun3.Component
	upstream  grpc.ClientStream
}

func (w *ClientStreamWrapper) Header() (metadata.MD, error) {
	return w.upstream.Header()
}
func (w *ClientStreamWrapper) Trailer() metadata.MD {
	return w.upstream.Trailer()
}
func (w *ClientStreamWrapper) CloseSend() error {
	return w.upstream.CloseSend()
}
func (w *ClientStreamWrapper) Context() context.Context {
	return w.upstream.Context()
}
func (w *ClientStreamWrapper) SendMsg(m interface{}) error {
	return w.upstream.SendMsg(m)
}
func (w *ClientStreamWrapper) RecvMsg(m interface{}) error {
	ret := w.upstream.RecvMsg(m)
	if ret != nil {
		if w.component != nil {
			if ret != io.EOF {
				w.component.SetException(ret, "ClientStream.RecvMsg", 2)
			}
			w.component.Finish()
		}
	} else {
		w.component.Finish()
	}
	return ret
}

func init() {
	tingyun3.Register(reflect.ValueOf(WrapServerRegisterService).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapServerServe).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapClientConnInvoke).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapClientConnNewStream).Pointer())
	tingyun3.Register(reflect.ValueOf(initTrampoline).Pointer())
}
