// +build linux
// +build amd64
// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

/*
#cgo LDFLAGS: -L${SRCDIR} -ltingyungosdk

extern void init(void *);
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
func WrapHttpClientDo(ptr uintptr, req *http.Request) (*http.Response, error) {
	var component *Component = nil
	if action := getAction(); action != nil {
		_, pc := GetCallerPC(3)
		method_name := runtime.FuncForPC(pc).Name()
		component = action.CreateExternalComponent(req.URL.String(), method_name)
		if trackId := component.CreateTrackID(); len(trackId) > 0 {
			req.Header.Add("X-Tingyun", trackId)
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
	if len(w.action.trackID) > 0 {
		//写跨应用追踪应答
		if txData := w.action.GetTxData(); len(txData) > 0 {
			headers := w.w.Header()
			headers.Set("X-Tingyun-Data", txData)
			// fmt.Println("Set txData:", txData)
		} else {
			// fmt.Println("Get txData failed\n")
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
	var method_name string
	class_name := reflect.TypeOf(handler).String()
	// fmt.Println("handle type: ", class_name)
	if class_name == "http.HandlerFunc" || class_name == "HandlerFunc" {
		handler_pc := reflect.ValueOf(handler).Pointer()
		method_name = runtime.FuncForPC(handler_pc).Name()
	} else {
		if len(class_name) > 0 && class_name[0] == '*' {
			class_name = class_name[1:]
		}
		method_name = class_name + ".ServeHTTP"
	}
	// fmt.Println("handle :", pattern, ", =>Method: ", method_name)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		action, _ := CreateAction("ROUTER", method_name)
		var rule *dataItemRules = nil
		resWriter := w
		if action != nil {
			action.SetHTTPMethod(strings.ToUpper(r.Method))
			if trackId := r.Header.Get("X-Tingyun"); len(trackId) > 0 {
				action.SetTrackID(trackId)
				// fmt.Println("Set track id: ", trackId)
			}
			rule = app.configs.dataItemRules.Get()
			for _, item := range rule.requestHeader {
				if value := r.Header.Get(item); len(value) > 0 {
					action.AddRequestParam(item, value)
				}
			}
			if readServerConfigBool(configServerConfigBoolCaptureParams, false) {
				action.url = r.URL.RequestURI()
			}
			resWriter = createWriteWraper(w, action, rule)
			setAction(action)
		}
		defer func() {
			if exception := recover(); exception != nil {
				//异常处理
				// fmt.Println("Exception:", exception)
				action.setError(exception, "error", 2)
				action.Finish()
				routineLocalRemove()
				if action != nil {
					resWriter.(*writeWrapper).reset()
				}
				//重新抛出异常
				panic(exception)
			} else {
				action.Finish()
				routineLocalRemove()
				if action != nil {
					resWriter.(*writeWrapper).reset()
				}
			}
		}()
		h.ServeHTTP(resWriter, r)
	})
}

//go:noinline
func WrapServerMuxHandle(ptr uintptr, pattern string, handler http.Handler) {
	ServerMuxHandle(ptr, pattern, wrapHandler(pattern, handler))
}

//go:noinline
func ServerServe(srv *http.Server, l net.Listener) error {
	fmt.Println(srv, l)
	return nil
}

//go:noinline
func WrapServerServe(srv *http.Server, l net.Listener) error {
	if srv.Handler != nil {
		srv.Handler = wrapHandler("", srv.Handler)
	}
	return ServerServe(srv, l)
}

//net/http.(*Server).Serve
//net/http.(*Server).ServeTLS
//net/http.(*Server).ListenAndServe
//net/http.(*Server).ListenAndServeTLS

// GetGID is Return the goroutine id
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
	C.init(unsafe.Pointer(p))
}
func init() {
	C.init(unsafe.Pointer(reflect.ValueOf(WrapServerServe).Pointer()))
	C.init(unsafe.Pointer(reflect.ValueOf(WrapServerMuxHandle).Pointer()))
	C.init(unsafe.Pointer(reflect.ValueOf(WrapHttpClientDo).Pointer()))
	C.init(unsafe.Pointer(reflect.ValueOf(GetGID).Pointer()))
}
