// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

//与server通信,login, upload

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	postRequest "github.com/TingYunGo/goagent/utils/httprequest"
)

//登陆中..1.inredirect
//       2.inInit
type serviceDC struct {
	locked       int32 //login由redirect状态切换到init状态时提供保护，防止这时候Release
	configs      *configurations
	request      *postRequest.Request
	aliveRequest *postRequest.Request
	uploadHost   string
	lastAlive    time.Time
	lastValidDC  int
	uploadSet    map[*postRequest.Request]int
}

func (s *serviceDC) keepAlive(callback func(error, map[string]interface{})) {
	//getCmd?version=3.2.0&sessionKey=4112274
	currTime := time.Now()
	if currTime.Sub(s.lastAlive) < 30*time.Second {
		return
	}
	if s.aliveRequest != nil {
		return
	}
	requrl, err := s.makeTraceURL("getCmd")
	if err != nil {
		return
	}
	s.lastAlive = currTime
	mtime := s.configs.serverExt.CIntegers.Read(configServerConfigIntegerMTime, 0)
	b, _ := json.Marshal(map[string]interface{}{
		"mTime": mtime,
	})
	Log().Println(LevelInfo, "getCmd:", requrl)
	Log().Println(LevelInfo|Audit, "Request Data:", string(b))
	s.aliveRequest, err = postRequest.New(requrl, map[string]string{}, b, time.Second*10, func(data []byte, statusCode int, err error) {
		if err == nil {
			Log().Println(LevelInfo, "getCmd Status Code:", statusCode)
			if len(data) > 0 {
				Log().Println(LevelInfo|Audit, "getCmd Response Data:", string(data))
			}
		}
		jsonData, er := parseJSON(data, statusCode, err)
		s.aliveRequest = nil
		callback(er, jsonData)
	})
}

//Login --启动登陆过程,如果已经在login中,返回error
func (s *serviceDC) Login(callback func(error, map[string]interface{})) error {
	if s.request != nil {
		return errors.New("Login already Startd")
	}
	if host := s.configs.local.CStrings.Read(configLocalStringNbsHost, ""); len(host) == 0 {
		return errors.New("No collector address in configuration file")
	}
	appName := s.configs.local.CStrings.Read(configLocalStringNbsAppName, "GO_LANG")
	license := s.configs.local.CStrings.Read(configLocalStringNbsLicenseKey, "_")

	requrl := fmt.Sprintf("%s/redirect?app=%s&license=%s&request=entry&version=%s", getRedirectHost(s, s.getConfigProtocol()), url.QueryEscape(appName), license, "3.2.0")
	params := make(map[string]string)
	var err error = nil
	Log().Println(LevelInfo, "Redirect:", requrl)
	s.request, err = postRequest.New(requrl, params, []byte("{}"), time.Second*10, func(data []byte, statusCode int, err error) {
		s.request = nil
		//完成回调,在另一个routine中触发
		use := atomic.AddInt32(&s.locked, 1)
		defer atomic.AddInt32(&s.locked, -1)
		if use != 1 {
			return
		}
		if err == nil {
			Log().Println(LevelInfo, "Redirect Status Code:", statusCode)
			if len(data) > 0 {
				Log().Println(LevelInfo|Audit, "Redirect Response Data:", string(data))
			}
		}
		var e error = nil
		for {
			var jsonData map[string]interface{}
			if jsonData, e = parseJSON(data, statusCode, err); e != nil {
				break
			}
			if s.uploadHost, e = parseRedirectResult(jsonData); e != nil {
				break
			}
			var b []byte
			b, e = makeLoginRequest()
			if e != nil {
				break
			}
			requrl := fmt.Sprintf("%s://%s/init?app=%s&license=%s&request=login&version=%s", s.getConfigProtocol(), s.uploadHost, url.QueryEscape(appName), url.QueryEscape(license), "3.2.0")
			Log().Println(LevelInfo, "Login:", requrl)
			Log().Println(LevelInfo|Audit, "Login Request: ", string(b))
			s.request, e = postRequest.New(requrl, map[string]string{ /*"Content-Encoding": "deflate"*/ }, b, time.Second*10, func(data []byte, statusCode int, err error) {
				use := atomic.AddInt32(&s.locked, 1)
				defer atomic.AddInt32(&s.locked, -1)
				if use != 1 {
					return
				}
				if err == nil {
					Log().Println(LevelInfo, "Login Status Code:", statusCode)
					Log().Println(LevelInfo, "Login Response Data:", string(data))
				}
				r, er := parseJSON(data, statusCode, err)
				callback(er, r)
			})
			break
		}
		if e != nil {
			s.lastValidDC++
			callback(e, nil)
		}
	})
	if err != nil {
		s.lastValidDC++
	}
	return err
}
func (s *serviceDC) getConfigProtocol() string {
	protocol := "https"
	if !s.configs.local.CBools.Read(configLocalBoolSSL, false) {
		protocol = "http"
	}
	return protocol
}
func (s *serviceDC) makeTraceURL(request string) (string, error) {
	sessionKey := s.configs.server.CStrings.Read(configServerStringAppSessionKey, "")
	if sessionKey == "" {
		return "", errors.New("makeTraceUrl: " + request + " server session key not found.")
	}
	appName := s.configs.local.CStrings.Read(configLocalStringNbsAppName, "GO_LANG")
	license := s.configs.local.CStrings.Read(configLocalStringNbsLicenseKey, "_")
	requrl := fmt.Sprintf("%s://%s/%s?app=%s&license=%s&sessionKey=%s&version=%s", s.getConfigProtocol(), s.uploadHost, request, url.QueryEscape(appName), url.QueryEscape(license), url.QueryEscape(sessionKey), "3.2.0")
	return requrl, nil
}

type requestContext struct {
	request *postRequest.Request
}

//上传数据,如果inLogin, 返回false,否则创建request,
func (s *serviceDC) Upload(data []byte, callback func(err error, rCode int, httpStatus int)) (*postRequest.Request, error) {
	requrl, err := s.makeTraceURL("trace")
	if err != nil {
		return nil, err
	}
	Log().Println(LevelInfo, "Upload", len(data), "bytes:", requrl)
	Log().Println(LevelInfo|Audit, "Upload Request Data:", len(data))
	return postRequest.New(requrl, map[string]string{ /*"Content-Encoding": "deflate"*/ }, data, time.Second*10, func(data []byte, statusCode int, err error) {
		use := atomic.AddInt32(&s.locked, 1)
		defer atomic.AddInt32(&s.locked, -1)
		if use != 1 {
			return
		}
		if err == nil {
			if len(data) > 0 {
				Log().Println(LevelInfo, "Upload Status Code:", statusCode, ", Data:", string(data))
			} else {
				Log().Println(LevelInfo, "Upload Status Code:", statusCode)
			}
		}
		r, er := parseJSON(data, statusCode, err)
		if er != nil {
			Log().Println(LevelError, "Upload Error:", er, r)
			callback(er, -2, statusCode)
		} else if status, er := jsonReadString(r, "status"); er == nil && status == "success" {
			callback(nil, -1, statusCode)
		} else {
			Log().Println(LevelError, "Upload Result:", string(data))
			callback(errors.New(string(data)), -1, statusCode)
		}
	})
}
func (s *serviceDC) Release() {
	for {
		if use := atomic.AddInt32(&s.locked, 1); use == 1 {
			break
		}
		atomic.AddInt32(&s.locked, -1)
		time.Sleep(1 * time.Millisecond)
	}
	s.request = nil
	s.configs = nil
}
func (s *serviceDC) init(config *configurations) {
	s.configs = config
	s.request = nil
	s.locked = 0
	s.lastValidDC = 0
	s.uploadSet = map[*postRequest.Request]int{}
}

func parseJSON(data []byte, statusCode int, err error) (map[string]interface{}, error) {

	if err != nil { //http过程有错误
		return nil, err
	}
	if statusCode != 200 {
		return nil, fmt.Errorf("server response status %d", statusCode)
	}
	jsonData := make(map[string]interface{})
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, err
	}
	return jsonData, nil
}

//解析redirect服务器返回值，成功则返回loginhost,失败则返回error
func parseRedirectResult(jsonData map[string]interface{}) (string, error) {
	result, resok := jsonData["result"]
	var resString string = "has no result"
	if resok {
		resString = fmt.Sprint(result)
	}
	if status, ok := jsonData["status"]; !ok { //验证是否有status
		return "", errors.New("server result have no status")
	} else if v, ok := status.(string); !ok { //类型验证
		return "", errors.New("server result status format error")
	} else if v != "success" { //值验证
		firstRun = false
		return "", errors.New("server result not success: " + resString)
	}
	if !resok {
		return "", errors.New("Redirect server status is success, no result")
	}
	return resString, nil
}

func getRedirectHost(s *serviceDC, protocol string) string {
	hosts := strings.Split(s.configs.local.CStrings.Read(configLocalStringNbsHost, ""), ",")
	if len(hosts) == 0 {
		return ""
	}
	host := hosts[s.lastValidDC%len(hosts)]

	array := strings.Split(host, "://")
	if len(array) > 1 {
		host = array[1]
	}
	array = strings.Split(host, "/")
	if len(array) > 1 {
		host = array[0]
	}
	array = strings.Split(host, ":")
	if len(array) > 1 {
		host = array[0]
	}
	port := 80
	if len(array) > 1 {
		if p, e := strconv.Atoi(array[1]); e == nil {
			port = p
		}
	}
	if protocol != "http" {
		port = 443
	}
	port = int(s.configs.local.CIntegers.Read(configLocalIntegerNbsPort, int64(port)))
	return fmt.Sprintf("%s://%s:%d", protocol, host, port)
}

type RequestHandler struct {
	request *postRequest.Request
}
