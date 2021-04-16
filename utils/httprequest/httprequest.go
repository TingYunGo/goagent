// Copyright 2016-2021 冯立强 fenglq@tingyun.com.  All rights reserved.

//Post请求异步封装
package postRequest

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/TingYunGo/goagent/utils/zip"
)

type Request struct {
	lock     sync.Mutex
	callback func(data []byte, statusCode int, err error)
}

func (r *Request) answer(data []byte, statusCode int, err error) {
	r.lock.Lock()
	callback := r.callback
	r.callback = nil
	r.lock.Unlock()
	if callback != nil {
		callback(data, statusCode, err)
	}
}

//释放请求对象，不管返回结果
func (r *Request) Release() {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.callback = nil
}

type ioReader struct {
	data  []byte
	begin int
}

//go:nosplit
func min(a, b int) int {

	if a < b {
		return a
	}
	return b
}
func (b *ioReader) Read(p []byte) (n int, err error) {
	if b.data == nil {
		return 0, io.EOF
	}
	if len(b.data) == b.begin {
		b.data = nil
		return 0, io.EOF
	}
	res := copy(p, b.data[b.begin:])
	b.begin += res
	if b.begin == len(b.data) {
		b.data = nil
	}
	return res, nil
}

//发起一个post请求,返回请求对象
func New(url string, params map[string]string, data []byte, duration time.Duration, callback func(data []byte, statusCode int, err error)) (*Request, error) {
	var body []byte
	var err error = nil
	if v, ok := params["Content-Encoding"]; ok && v == "deflate" {
		body, err = zip.Deflate(data)
	} else {
		body = data
	}
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", url, &ioReader{data: body})
	if nil != err {
		return nil, err
	}
	useParams := make(map[string]string)
	useParams["Accept-Encoding"] = "identity, deflate"
	useParams["Content-Type"] = "Application/json;charset=UTF-8"
	useParams["User-Agent"] = "TingYun-Agent/GoLang"
	for k, v := range params {
		useParams[k] = v
	}
	res := &Request{callback: callback}
	for k, v := range useParams {
		request.Header.Add(k, v)
	}
	go func(request *http.Request) {
		client := &http.Client{Timeout: duration}
		defer func() {
			if exception := recover(); exception != nil {
				fmt.Println(exception)
			}
			request.Body.Close()
		}()
		response, err := client.Do(request)
		if err != nil {
			res.answer(nil, -1, err)
			return
		}
		defer response.Body.Close()
		if response.StatusCode == 200 {
			if b, err := ioutil.ReadAll(response.Body); err != nil { //server返回200，然后读数据失败....
				res.answer(nil, 200, err)
			} else {
				encoding := response.Header.Get("Content-Encoding")
				if encoding == "gzip" || encoding == "deflate" {
					d, err := zip.Inflate(b)
					if err == nil {
						res.answer(d, 200, nil)
						return
					}
				}
				res.answer(b, 200, nil)
			}
		} else {
			res.answer(nil, response.StatusCode, nil)
		}

	}(request)
	return res, nil
}
