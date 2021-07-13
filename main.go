// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64

package tingyun3

/*
#cgo LDFLAGS: -L${SRCDIR} -ltingyungosdk

extern void tingyun_go_init(void *);
*/
import "C"
import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

//go:noinline
func ServerMuxHandle(ptr uintptr, pattern string, handler http.Handler) {
	fmt.Println(pattern, handler)
}

//go:noinline
func HttpClientDo(ptr uintptr, req *http.Request) (*http.Response, error) {
	fmt.Println(ptr, req)
	return nil, nil
}

//go:noinline
func WrapHttpClientDo(ptr uintptr, req *http.Request) (*http.Response, error) {
	var component *Component = nil
	if action := getAction(); action != nil {
		_, pc := GetCallerPC(3)
		methodName := runtime.FuncForPC(pc).Name()
		component = action.CreateExternalComponent(req.URL.String(), methodName)
		if trackID := component.CreateTrackID(); len(trackID) > 0 {
			req.Header.Add("X-Tingyun", trackID)
		}
	}
	defer func() {
		if exception := recover(); exception != nil {
			component.setError(exception, "error")
			component.Finish()
			panic(exception)
		}
	}()
	res, err := HttpClientDo(ptr, req)
	if component != nil {
		if err != nil {
			component.setError(err, "httpClient")
		} else if res != nil {
			if txdata := res.Header.Get("X-Tingyun-Data"); len(txdata) > 0 {
				component.SetTxData(txdata)
			}
		}
		component.Finish()
	}
	return res, err
}

type writeWrapper struct {
	http.ResponseWriter
	w       http.ResponseWriter
	action  *Action
	rules   *dataItemRules
	answerd bool
}

func (w *writeWrapper) reset() {
	w.w = nil
	w.action = nil
	w.rules = nil
}

func createWriteWraper(w http.ResponseWriter, action *Action, rule *dataItemRules) http.ResponseWriter {
	r := &writeWrapper{}
	r.w = w
	r.action = action
	r.rules = rule
	r.answerd = false
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
	return w.w.(http.Hijacker).Hijack()
}
func wrapHandler(pattern string, handler http.Handler) http.Handler {
	h := handler
	var methodName string
	isRouteMode := true
	className := reflect.TypeOf(handler).String()
	if className == "http.HandlerFunc" || className == "HandlerFunc" {
		handlerPC := reflect.ValueOf(handler).Pointer()
		methodName = runtime.FuncForPC(handlerPC).Name()
	} else {
		isRouteMode = false
		if len(className) > 0 && className[0] == '*' {
			className = className[1:]
		}
		methodName = className + ".ServeHTTP"
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("on: ", pattern, ", by: ", methodName)
		var component *Component = nil
		action := getAction()
		if action != nil {
			component = action.CreateComponent(methodName)
		}
		if component == nil {
			if isRouteMode {
				action, _ = CreateAction("ROUTER", methodName)
			} else {
				action, _ = CreateAction("URI", r.URL.Path)
				if action != nil {
					action.method = methodName
					action.root.method = methodName
				}
			}
		}
		action.SetName("CLIENTIP", parseIP(r.RemoteAddr))
		resWriter := w

		if action != nil && component == nil {
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
			resWriter = createWriteWraper(w, action, rule)
			setAction(action)
		}
		defer func() {
			exception := recover()
			if exception != nil && component == nil {
				action.setError(exception, "error", 2)
			}
			if component != nil {
				component.Finish()
			} else {
				action.Finish()
				routineLocalRemove()
				if action != nil {
					resWriter.(*writeWrapper).reset()
				}
			}
			//re throw
			if exception != nil {
				panic(exception)
			}
		}()
		h.ServeHTTP(resWriter, r)
	})
}

type httpListenAddress struct {
	Addr string
	tls  bool
}

var httpListenAddr httpListenAddress

//go:noinline
func WrapServerMuxHandle(ptr uintptr, pattern string, handler http.Handler) {
	// fmt.Println("Wrap: ", pattern, ", By: ", reflect.TypeOf(handler).String())
	ServerMuxHandle(ptr, pattern, wrapHandler(pattern, handler))
}

//go:noinline
func ServerServe(srv *http.Server, l net.Listener) error {
	fmt.Println(srv, l)
	return nil
}

//go:noinline
func WrapServerServe(srv *http.Server, l net.Listener) error {

	pre := httpListenAddr
	httpListenAddr = httpListenAddress{
		Addr: srv.Addr,
		tls:  srv.TLSConfig != nil,
	}
	if app != nil {
		app.logger.Println(LevelDebug, "http.Server.Serve:", httpListenAddr.Addr)
	}

	if srv.Handler != nil {
		srv.Handler = wrapHandler("", srv.Handler)
	}
	e := ServerServe(srv, l)
	httpListenAddr = pre
	return e
}

//go:noinline
func httpNotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, r)
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
			resWriter = createWriteWraper(w, action, rule)
		}
	}
	httpNotFound(resWriter, r)
	if action != nil && !found {
		action.Finish()
	}
}

// GetGID return goroutine id
//go:noinline
func GetGID() int64 {
	fmt.Println(1)
	return 0
}

// GetCallerPC return caller pc
//go:noinline
func GetCallerPC(layer int) (l int, pc uintptr) {
	if pc, _, _, success := runtime.Caller(layer); success {
		return layer, pc
	}
	return 0, 0
}

//GetCallerName : 取layer层调用栈函数名
//go:noinline
func GetCallerName(layer int) string {
	if _, pc := GetCallerPC(layer + 1); pc != 0 {
		return runtime.FuncForPC(pc).Name()
	}
	return ""
}

// Register : native method
func Register(p uintptr) {
	C.tingyun_go_init(unsafe.Pointer(p))
}
func init() {
	C.tingyun_go_init(unsafe.Pointer(reflect.ValueOf(WraphttpNotFound).Pointer()))
	C.tingyun_go_init(unsafe.Pointer(reflect.ValueOf(WrapServerServe).Pointer()))
	C.tingyun_go_init(unsafe.Pointer(reflect.ValueOf(WrapServerMuxHandle).Pointer()))
	C.tingyun_go_init(unsafe.Pointer(reflect.ValueOf(WrapHttpClientDo).Pointer()))
	C.tingyun_go_init(unsafe.Pointer(reflect.ValueOf(GetGID).Pointer()))
}
