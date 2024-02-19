// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64 arm64
// +build cgo

package tingyun3

/*

extern int tingyun_go_init(void *);

*/
import "C"
import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/TingYunGo/goagent/libs/tystring"
)

//go:noinline
func ServerMuxHandle(ptr *http.ServeMux, pattern string, handler http.Handler) {
	idPointer.arg7 = *idPointer.idpointer + idPointer.idindex + idPointer.arg1 + idPointer.arg2 + idPointer.arg3 + idPointer.arg4 + idPointer.arg5 + idPointer.arg6 + idPointer.arg7 +
		idPointer.arg8 + idPointer.arg9 + idPointer.arg10 + idPointer.arg11 + idPointer.arg12 + idPointer.arg13 + idPointer.arg14 + idPointer.arg15 + idPointer.arg16 +
		idPointer.arg17 + idPointer.arg18 + idPointer.arg19 + idPointer.arg20
}

//go:noinline
func HttpClientDo(ptr *http.Client, req *http.Request) (*http.Response, error) {
	idPointer.arg6 = *idPointer.idpointer + idPointer.idindex + idPointer.arg1 + idPointer.arg2 + idPointer.arg3 + idPointer.arg4 + idPointer.arg5 + idPointer.arg6 + idPointer.arg7 +
		idPointer.arg8 + idPointer.arg9 + idPointer.arg10 + idPointer.arg11 + idPointer.arg12 + idPointer.arg13 + idPointer.arg14 + idPointer.arg15 + idPointer.arg16 +
		idPointer.arg17 + idPointer.arg18 + idPointer.arg19 + idPointer.arg20
	return nil, nil
}

//go:noinline
func WrapHttpClientDo(ptr *http.Client, req *http.Request) (*http.Response, error) {
	var component *Component = nil
	if action := getAction(); action != nil {
		_, pc := GetCallerPC(3)
		methodName := runtime.FuncForPC(pc).Name()
		component = action.CreateExternalComponent(req.URL.String(), methodName)
		if component != nil {
			if trackID := component.CreateTrackID(); len(trackID) > 0 {
				req.Header.Add("X-Tingyun", trackID)
			}
		}
	}
	defer func() {
		if exception := recover(); exception != nil {
			if component != nil {
				component.setError(exception, "error", true)
				component.Finish()
			}
			panic(exception)
		}
	}()
	res, err := HttpClientDo(ptr, req)
	if component != nil {
		if err != nil {
			component.setError(err, "httpClient", false)
		} else if res != nil {
			if txdata := res.Header.Get("X-Tingyun-Data"); len(txdata) > 0 {
				component.SetTxData(txdata)
				component.SetStatusCode(res.StatusCode)
			}
		}
		component.FixStackEnd(1, func(funcName string) bool {
			token := "net/http"
			return tystring.SubString(funcName, 0, len(token)) == token
		})
	}
	return res, err
}

//go:noinline
func httputilReverseProxyServeHTTP(p *httputil.ReverseProxy, rw http.ResponseWriter, req *http.Request) {
	idPointer.arg6 = *idPointer.idpointer + idPointer.idindex + idPointer.arg1 + idPointer.arg2 + idPointer.arg3 + idPointer.arg4 + idPointer.arg5 + idPointer.arg6 + idPointer.arg7 +
		idPointer.arg8 + idPointer.arg9 + idPointer.arg10 + idPointer.arg11 + idPointer.arg12 + idPointer.arg13 + idPointer.arg14 + idPointer.arg15 + idPointer.arg16 +
		idPointer.arg17 + idPointer.arg18 + idPointer.arg19 + idPointer.arg20
}

//go:noinline
func WraphttputilReverseProxyServeHTTP(p *httputil.ReverseProxy, rw http.ResponseWriter, req *http.Request) {
	var component *Component = nil
	if action := getAction(); action != nil {
		_, pc := GetCallerPC(3)
		methodName := runtime.FuncForPC(pc).Name()
		url := req.URL.String()
		component = action.CreateExternalComponent(url, methodName)
		if component != nil {
			SetComponent(component)
			if trackID := component.CreateTrackID(); len(trackID) > 0 {
				req.Header.Add("X-Tingyun", trackID)
			}
		}
	}
	defer func() {
		if exception := recover(); exception != nil {
			if component != nil {
				component.setError(exception, "httputil.(*ReverseProxy).ServeHTTP", true)
				component.Finish()
			}
			panic(exception)
		}
	}()
	httputilReverseProxyServeHTTP(p, rw, req)
	if component != nil {
		if GetComponent() == component {
			SetComponent(nil)
		}
		res := req.Response
		if res != nil {
			if txdata := res.Header.Get("X-Tingyun-Data"); len(txdata) > 0 {
				component.SetTxData(txdata)
			}
		}
		component.FixStackEnd(1, func(funcName string) bool {
			token := "net/http"
			return tystring.SubString(funcName, 0, len(token)) == token
		})
	}
}

//go:noinline
func httpTransportroundTrip(t *http.Transport, req *http.Request) (*http.Response, error) {
	idPointer.arg6 = *idPointer.idpointer + idPointer.idindex + idPointer.arg1 + idPointer.arg2 + idPointer.arg3 + idPointer.arg4 + idPointer.arg5 + idPointer.arg6 + idPointer.arg7 +
		idPointer.arg8 + idPointer.arg9 + idPointer.arg10 + idPointer.arg11 + idPointer.arg12 + idPointer.arg13 + idPointer.arg14 + idPointer.arg15 + idPointer.arg16 +
		idPointer.arg17 + idPointer.arg18 + idPointer.arg19 + idPointer.arg20
	return nil, nil
}

//go:noinline
func WraphttpTransportroundTrip(t *http.Transport, req *http.Request) (*http.Response, error) {
	url := req.URL.String()
	component := getComponent()
	if component != nil {
		component.SetURL(url)
	}
	r, e := httpTransportroundTrip(t, req)
	if r != nil {
		if component != nil {
			if txdata := r.Header.Get("X-Tingyun-Data"); len(txdata) > 0 {
				component.SetTxData(txdata)
			}
			component.SetStatusCode(r.StatusCode)
		}
	}
	return r, e
}

type NetConnWrapper struct {
	conn net.Conn
	w    *writeWrapper
	url  string
}

func (c *NetConnWrapper) Url() string {
	if c == nil {
		return ""
	}
	return c.url
}

func (c *NetConnWrapper) Read(b []byte) (n int, err error) {
	return c.conn.Read(b)
}
func (c *NetConnWrapper) Write(b []byte) (n int, err error) {
	if c.w != nil {
		if c.w.handleBinaryWrite(b) {
			c.w = nil
		}
	}
	return c.conn.Write(b)
}
func (c *NetConnWrapper) Close() error {
	return c.conn.Close()
}
func (c *NetConnWrapper) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}
func (c *NetConnWrapper) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
func (c *NetConnWrapper) SetDeadline(t time.Time) error {
	return c.conn.SetDeadline(t)
}
func (c *NetConnWrapper) SetReadDeadline(t time.Time) error {
	return c.conn.SetReadDeadline(t)
}
func (c *NetConnWrapper) SetWriteDeadline(t time.Time) error {
	return c.conn.SetWriteDeadline(t)
}

type writeWrapper struct {
	http.ResponseWriter
	w       http.ResponseWriter
	action  *Action
	rules   *dataItemRules
	gid     int64
	wCache  []byte
	url     string
	answerd bool
}

func (w *writeWrapper) parseSwitchProtocol() bool {
	if w.wCache == nil {
		return false
	}
	if strings.LastIndex(string(w.wCache), "\r\n\r\n") == -1 {
		if len(w.wCache) == 4000 {
			w.wCache = nil
			return true
		}
		return false
	}
	//解析http应答头
	index := strings.Index(string(w.wCache), "\r\n\r\n")
	header := w.wCache[0:index]
	parts := strings.Split(string(header), "\r\n")
	w.wCache = nil
	answer_code := 0
	if len(parts) > 0 {
		splitMapString(parts[0], isSpace, func(_, value string) {
			splitMapString(value, isSpace, func(answerCode, _ string) {
				if code, err := strconv.Atoi(answerCode); err == nil {
					answer_code = code
				}
			})
		})
		// for _, line := range parts {
		// 	fmt.Println("SWITCH:", line)
		// }
		if answer_code > 0 {
			w.action.SetHTTPStatus(uint16(answer_code), 3)
			w.answerd = true
			w.action.Finish()
			if w.gid != 0 {
				removeRoutineLocal(w.gid)
			}
			w.reset()
		}
	}
	return true
}
func (w *writeWrapper) handleBinaryWrite(b []byte) bool {
	if w == nil {
		return true
	}
	w.wCache = limitAppend(w.wCache, b, 4000)
	//协议切换头处理
	return w.parseSwitchProtocol()
}

func (w *writeWrapper) reset() {
	w.w = nil
	w.action = nil
	w.rules = nil
	w.gid = 0
	w.wCache = nil
	w.url = ""
}

func createWriteWraper(w http.ResponseWriter, req *http.Request, action *Action, rule *dataItemRules) http.ResponseWriter {
	r := &writeWrapper{}
	r.w = w
	r.action = action
	r.rules = rule
	r.answerd = false
	r.gid = GetGID()

	proto := "://"
	if req.TLS != nil {
		proto = "s://"
	}
	r.url = proto + req.Host + req.RequestURI

	return r
}

func (w *writeWrapper) onAnswer(statusCode int) {
	if w.answerd {
		return
	}
	if w.action == nil {
		return
	}
	if w.rules != nil {
		for _, item := range w.rules.responseHeader {
			if value := w.w.Header().Get(item); len(value) > 0 {
				w.action.AddResponseParam(item, value)
			}
		}
	}
	if w.action != nil {
		if len(w.action.trackID) > 0 {
			if txData := w.action.GetTxData(); len(txData) > 0 {
				headers := w.w.Header()
				headers.Set("X-Tingyun-Data", txData)
			}
		}
	}

	w.action.SetHTTPStatus(uint16(statusCode), 3)
	w.answerd = true
	if statusCode == 101 {
		w.action.Finish()
		if w.gid != 0 {
			removeRoutineLocal(w.gid)
		}
		w.reset()
	}
}

func (w *writeWrapper) Header() http.Header {
	if w.w == nil {
		return nil
	}
	return w.w.Header()
}
func (w *writeWrapper) Write(b []byte) (int, error) {
	if w.w == nil {
		return -1, errors.New("null writer")
	}
	w.onAnswer(int(200))
	return w.w.Write(b)
}
func (w *writeWrapper) ReadFrom(src io.Reader) (n int64, err error) {
	return w.w.(io.ReaderFrom).ReadFrom(src)
}
func (w *writeWrapper) Flush() {
	w.w.(http.Flusher).Flush()
}
func (w *writeWrapper) WriteHeader(statusCode int) {
	if w.w != nil {
		w.onAnswer(statusCode)
		w.w.WriteHeader(statusCode)
	}
}
func (w *writeWrapper) CloseNotify() <-chan bool {
	return w.w.(http.CloseNotifier).CloseNotify()
}
func (w *writeWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	c, rw, e := w.w.(http.Hijacker).Hijack()
	if c != nil && e == nil && readLocalConfigBool(configLocalBoolGorillaWebsocket, false) {
		conn := &NetConnWrapper{c, w, w.url}
		c = conn
	}
	return c, rw, e
}
func CatchRequestHeaders() []string {
	if app == nil {
		return nil
	}
	rule := app.configs.dataItemRules.Get()
	if rule == nil {
		return nil
	}
	return rule.requestHeader
}
func wrapHandler(pattern string, handler http.Handler) http.Handler {
	h := handler
	var methodName string
	// isRouteMode := true
	className := reflect.TypeOf(handler).String()
	if className == "http.HandlerFunc" || className == "HandlerFunc" {
		handlerPC := reflect.ValueOf(handler).Pointer()
		methodName = runtime.FuncForPC(handlerPC).Name()
	} else {
		// isRouteMode = false
		if len(className) > 0 && className[0] == '*' {
			className = className[1:]
		}
		methodName = className + ".ServeHTTP"
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var component *Component = nil
		action := getAction()
		preAction := false
		if action != nil {
			preAction = true
			component = action.CreateComponent(methodName)
		}
		if !preAction {
			if readServerConfigBool(ServerConfigBoolAutoActionNaming, true) {
				action, _ = CreateAction("ROUTER", methodName)
			} else {
				action, _ = CreateAction("URI", r.URL.Path)
				if action != nil {
					action.method = methodName
					action.root.method = methodName
				}
			}
			if action != nil {
				if len(action.name) == 0 {
					action.name = namingCustomizeName(r)
				}
				r = r.WithContext(context.WithValue(r.Context(), "TingYunWebAction", action))
			}
		}
		action.SetName("CLIENTIP", parseIP(r.RemoteAddr))
		resWriter := w

		if action != nil && !preAction {
			action.SetHTTPMethod(strings.ToUpper(r.Method))
			if trackID := r.Header.Get("X-Tingyun"); len(trackID) > 0 {
				action.SetTrackID(trackID)
			}
			rule := app.configs.dataItemRules.Get()
			for _, item := range rule.requestHeader {
				if value := r.Header.Get(item); len(value) > 0 {
					action.AddRequestParam(item, value)
				}
			}
			if readServerConfigBool(configServerConfigBoolCaptureParams, false) {
				protocol := "http"
				if r.TLS != nil {
					protocol = "https"
				}
				action.SetURL(protocol + "://" + r.Host + r.RequestURI)
			}
			resWriter = createWriteWraper(w, r, action, rule)
			setAction(action)
		}
		defer func() {
			a := getAction()
			if a != nil && a != action {
				a.Finish()
			}
			routineLocalRemove()
			exception := recover()
			if exception != nil && !preAction {
				action.setError(exception, "error", 2, true)
			}
			if component != nil {
				component.Finish()
			} else if action != nil {
				resWriter.(*writeWrapper).reset()
			}
			action.Finish()
			//re throw
			if exception != nil {
				panic(exception)
			}
		}()
		h.ServeHTTP(resWriter, r)
	})
}

//go:noinline
func WrapServerMuxHandle(ptr *http.ServeMux, pattern string, handler http.Handler) {
	// fmt.Println("Wrap: ", pattern, ", By: ", reflect.TypeOf(handler).String())
	ServerMuxHandle(ptr, pattern, wrapHandler("ServeMuxHandle: pattern="+pattern, handler))
}

//go:noinline
func ServerMux121handle(ptr *http.ServeMux, pattern string, handler http.Handler) {
	idPointer.arg7 = *idPointer.idpointer + idPointer.idindex + idPointer.arg1 + idPointer.arg2 + idPointer.arg3 + idPointer.arg4 + idPointer.arg5 + idPointer.arg6 + idPointer.arg7 +
		idPointer.arg8 + idPointer.arg9 + idPointer.arg10 + idPointer.arg11 + idPointer.arg12 + idPointer.arg13 + idPointer.arg14 + idPointer.arg15 + idPointer.arg16 +
		idPointer.arg17 + idPointer.arg18 + idPointer.arg19 + idPointer.arg20
}

//go:noinline
func WrapServerMux121handle(ptr *http.ServeMux, pattern string, handler http.Handler) {
	ServerMux121handle(ptr, pattern, wrapHandler("ServeMux121.handle: pattern="+pattern, handler))
}

//func (mux *ServeMux) registerErr(patstr string, handler Handler) error

//go:noinline
func ServeMuxregisterErr(ptr *http.ServeMux, pattern string, handler http.Handler) error {
	idPointer.arg7 = *idPointer.idpointer + idPointer.idindex + idPointer.arg1 + idPointer.arg2 + idPointer.arg3 + idPointer.arg4 + idPointer.arg5 + idPointer.arg6 + idPointer.arg7 +
		idPointer.arg8 + idPointer.arg9 + idPointer.arg10 + idPointer.arg11 + idPointer.arg12 + idPointer.arg13 + idPointer.arg14 + idPointer.arg15 + idPointer.arg16 +
		idPointer.arg17 + idPointer.arg18 + idPointer.arg19 + idPointer.arg20
	return nil
}

//go:noinline
func WrapServeMuxregisterErr(ptr *http.ServeMux, pattern string, handler http.Handler) error {
	return ServeMuxregisterErr(ptr, pattern, wrapHandler("ServeMux.registerErr: pattern="+pattern, handler))
}

//go:noinline
func ServerMuxHandler(ptr *http.ServeMux, r *http.Request) (h http.Handler, pattern string) {
	idPointer.arg5 = *idPointer.idpointer + idPointer.idindex + idPointer.arg1 + idPointer.arg2 + idPointer.arg3 + idPointer.arg4 + idPointer.arg5 + idPointer.arg6 + idPointer.arg7 +
		idPointer.arg8 + idPointer.arg9 + idPointer.arg10 + idPointer.arg11 + idPointer.arg12 + idPointer.arg13 + idPointer.arg14 + idPointer.arg15 + idPointer.arg16 +
		idPointer.arg17 + idPointer.arg18 + idPointer.arg19 + idPointer.arg20
	return nil, ""
}

//go:noinline
func WrapServerMuxHandler(ptr *http.ServeMux, r *http.Request) (h http.Handler, pattern string) {

	hres, pattern := ServerMuxHandler(ptr, r)
	className := reflect.TypeOf(hres).String()
	if className == "http.HandlerFunc" || className == "HandlerFunc" {
		handlerPC := reflect.ValueOf(hres).Pointer()
		if runtime.FuncForPC(handlerPC).Name() == "net/http.NotFound" {
			return http.HandlerFunc(WraphttpNotFound), pattern
		}
	}
	return hres, pattern
}

//go:noinline
func HttpServerServe(srv *http.Server, l net.Listener) error {
	idPointer.arg3 = *idPointer.idpointer + idPointer.idindex + idPointer.arg1 + idPointer.arg2 + idPointer.arg3 + idPointer.arg4 + idPointer.arg5 + idPointer.arg6 + idPointer.arg7 +
		idPointer.arg8 + idPointer.arg9 + idPointer.arg10 + idPointer.arg11 + idPointer.arg12 + idPointer.arg13 + idPointer.arg14 + idPointer.arg15 + idPointer.arg16 +
		idPointer.arg17 + idPointer.arg18 + idPointer.arg19 + idPointer.arg20
	return nil
}

func AppendListenAddress(address string) {
	listens.Append(address)
}

//go:noinline
func WrapHttpServerServe(srv *http.Server, l net.Listener) error {
	addr := srv.Addr
	network := ""
	if len(addr) == 0 && l != nil {
		a := l.Addr()
		network = a.Network()
		addr = a.String()
	}
	if len(addr) > 0 {
		listens.Append(addr)
	}
	if app != nil {
		app.logger.Println(LevelDebug, "http.Server.Serve:", network, addr)
	}

	if srv.Handler != nil {
		srv.Handler = wrapHandler("HttpServerServe", srv.Handler)
	}

	e := HttpServerServe(srv, l)
	return e

}

//go:noinline
func httpNotFound(w http.ResponseWriter, r *http.Request) {
	idPointer.arg1 = *idPointer.idpointer + idPointer.idindex + idPointer.arg1 + idPointer.arg2 + idPointer.arg3 + idPointer.arg4 + idPointer.arg5 + idPointer.arg6 + idPointer.arg7 +
		idPointer.arg8 + idPointer.arg9 + idPointer.arg10 + idPointer.arg11 + idPointer.arg12 + idPointer.arg13 + idPointer.arg14 + idPointer.arg15 + idPointer.arg16 +
		idPointer.arg17 + idPointer.arg18 + idPointer.arg19 + idPointer.arg20
}

//go:noinline
func WraphttpNotFound(w http.ResponseWriter, r *http.Request) {
	action := getAction()
	found := action != nil
	resWriter := w
	if action == nil {
		if action, _ = CreateAction("URI", r.URL.Path); action != nil {

			if trackID := r.Header.Get("X-Tingyun"); len(trackID) > 0 {
				action.SetTrackID(trackID)
			}

			rule := app.configs.dataItemRules.Get()
			for _, item := range rule.requestHeader {
				if value := r.Header.Get(item); len(value) > 0 {
					action.AddRequestParam(item, value)
				}
			}
			if readServerConfigBool(configServerConfigBoolCaptureParams, false) {
				protocol := "http"
				if r.TLS != nil {
					protocol = "https"
				}
				action.SetName("CLIENTIP", parseIP(r.RemoteAddr))
				action.SetURL(protocol + "://" + r.Host + r.RequestURI)
			}
			resWriter = createWriteWraper(w, r, action, rule)
		}
	}
	httpNotFound(resWriter, r)
	if action != nil && !found {
		action.Finish()
	}
}

var idPointer *pidStruct = &pidStruct{}

// GetGID return goroutine id
//go:noinline
func GetGID() int64 {
	return *idPointer.idpointer + idPointer.idindex + idPointer.arg1 + idPointer.arg2 + idPointer.arg3 + idPointer.arg4 + idPointer.arg5 + idPointer.arg6 + idPointer.arg7 +
		idPointer.arg8 + idPointer.arg9 + idPointer.arg10 + idPointer.arg11 + idPointer.arg12 + idPointer.arg13 + idPointer.arg14 + idPointer.arg15
}

//go:noinline
func setGID(p *pidStruct) int64 {
	idPointer = p
	fmt.Println(1, p)
	return 0
}

// Register : native method
func Register(p uintptr) {
	C.tingyun_go_init(unsafe.Pointer(p))
}

type ClientDoFunc func(*http.Client, *http.Request) (*http.Response, error)

type pidStruct struct {
	idpointer *int64
	idindex   int64
	arg1      int64
	arg2      int64
	arg3      int64
	arg4      int64
	arg5      int64
	arg6      int64
	arg7      int64
	arg8      int64
	arg9      int64
	arg10     int64
	arg11     int64
	arg12     int64
	arg13     int64
	arg14     int64
	arg15     int64
	arg16     int64
	arg17     int64
	arg18     int64
	arg19     int64
	arg20     int64
}

//go:noinline
func setRoutineID(p *pidStruct) {
	idPointer = p
}

func init() {
	C.tingyun_go_init(unsafe.Pointer(reflect.ValueOf(WrapHttpServerServe).Pointer()))
	C.tingyun_go_init(unsafe.Pointer(reflect.ValueOf(WrapServerMuxHandler).Pointer()))
	C.tingyun_go_init(unsafe.Pointer(reflect.ValueOf(WraphttpNotFound).Pointer()))
	C.tingyun_go_init(unsafe.Pointer(reflect.ValueOf(WrapServerMuxHandle).Pointer()))
	C.tingyun_go_init(unsafe.Pointer(reflect.ValueOf(WrapServerMux121handle).Pointer()))
	C.tingyun_go_init(unsafe.Pointer(reflect.ValueOf(WrapServeMuxregisterErr).Pointer()))
	C.tingyun_go_init(unsafe.Pointer(reflect.ValueOf(WrapHttpClientDo).Pointer()))
	C.tingyun_go_init(unsafe.Pointer(reflect.ValueOf(WraphttputilReverseProxyServeHTTP).Pointer()))
	C.tingyun_go_init(unsafe.Pointer(reflect.ValueOf(WraphttpTransportroundTrip).Pointer()))
	C.tingyun_go_init(unsafe.Pointer(reflect.ValueOf(GetGID).Pointer()))
	C.tingyun_go_init(unsafe.Pointer(reflect.ValueOf(setGID).Pointer()))
}
