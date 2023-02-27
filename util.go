// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/TingYunGo/goagent/libs/tystring"
)

func round(x float64) int {
	if x < 0.0 {
		return -round(-x)
	}
	r := int(x)
	if x-float64(r) < 0.5 {
		return r
	}
	return r + 1
}
func binarySearch(p []float64, value float64) int {
	Begin, Len := 0, len(p)
	for Len > 0 {
		middle := Len / 2
		if value < p[Begin+middle] {
			Len = middle
		} else if p[Begin+middle] < value {
			Begin, Len = Begin+middle+1, Len-middle-1
		} else {
			return Begin + middle
		}
	}
	return -1 - Begin
}
func patchSize(src string, maxSize int) string {
	if len(src) > maxSize {
		src = src[0:maxSize]
	}
	return src
}
func toString(e interface{}) string {
	return fmt.Sprint(e)
}
func callStack(skip int) []string {
	var slice []string
	slice = make([]string, 0, 15)
	opc := uintptr(0)
	for i := skip + 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if opc == pc {
			continue
		}
		fname := getnameByAddr(pc)
		index := strings.Index(fname, "/tingyun/")
		if index > 0 {
			continue
		}
		opc = pc
		//截断源文件名
		index = strings.Index(file, "/src/")
		if index > 0 {
			file = file[index+5 : len(file)]
		}
		slice = append(slice, fmt.Sprintf("%s(%s:%d)", fname, file, line))
	}
	return slice
}
func validCallStack(skip int, removed func(string) bool) []string {
	var slice []string
	slice = make([]string, 0, 15)
	opc := uintptr(0)
	lineRemoved := true
	for i := skip + 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if opc == pc {
			continue
		}
		fname := getnameByAddr(pc)
		if lineRemoved {
			if lineRemoved = removed(fname); lineRemoved {
				continue
			}
		}
		index := strings.Index(fname, "/tingyun/")
		if index > 0 {
			continue
		}
		opc = pc
		//截断源文件名
		index = strings.Index(file, "/src/")
		if index > 0 {
			file = file[index+5 : len(file)]
		}
		slice = append(slice, fmt.Sprintf("%s(%s:%d)", fname, file, line))
	}
	return slice
}
func getnameByAddr(p interface{}) string {
	ptr, _ := strconv.ParseInt(fmt.Sprintf("%x", p), 16, 64)
	return runtime.FuncForPC(uintptr(ptr)).Name()
}
func unicID(t time.Time, p interface{}) string {
	return strings.Replace(fmt.Sprintf("%x-%p", t.UnixNano(), p), "0x", "", -1)
}

func md5sum(src string) string {
	val := md5.New()
	val.Write([]byte(src))
	return fmt.Sprintf("%x", val.Sum(nil))
}

type timeRange struct {
	begin    time.Time
	duration time.Duration
}

func (t *timeRange) GetDuration(endTime time.Time) time.Duration {
	if t.duration == -1 {
		return endTime.Sub(t.begin)
	}
	return t.duration
}
func (t *timeRange) End() {

	t.duration = time.Now().Sub(t.begin)
}

func (t *timeRange) Init() {
	t.begin = time.Now()
	t.duration = -1
}

// EndTime : 时段结束时间
func (t *timeRange) EndTime() time.Time {
	ret := t.begin
	return ret.Add(t.duration)
}

//Inside : 检测 t时段是否是 r时段的子集
func (t *timeRange) Inside(r *timeRange) bool {
	if t.begin.Before(r.begin) || r.duration < t.duration {
		return false
	}
	return !t.EndTime().After(r.EndTime())
}

func jsonReadString(jsonData map[string]interface{}, name string) (string, error) {
	if r, ok := jsonData[name]; !ok { //验证是否有name
		return "", errors.New("Has no " + name)
	} else if v, ok := r.(string); !ok { //类型验证
		return "", errors.New("json \"" + name + "\" not string.")
	} else {
		return v, nil
	}
}
func jsonReadArray(jsonData map[string]interface{}, name string) ([]interface{}, error) {
	if r, ok := jsonData[name]; !ok { //验证是否有name
		return nil, errors.New("Has no " + name)
	} else if v, ok := r.([]interface{}); !ok { //类型验证
		return nil, errors.New("json \"" + name + "\" not array.")
	} else {
		return v, nil
	}
}

func jsonReadObjects(jsonData map[string]interface{}, name string) (map[string]interface{}, error) {
	if r, ok := jsonData[name]; !ok { //验证是否有name
		return nil, errors.New("Has no " + name)
	} else if v, ok := r.(map[string]interface{}); !ok { //类型验证
		return nil, errors.New("json \"" + name + "\" not objects.")
	} else {
		return v, nil
	}
}
func jsonReadBool(jsonData map[string]interface{}, name string) (bool, error) {
	if r, ok := jsonData[name]; !ok { //验证是否有name
		return false, errors.New("Has no " + name)
	} else if v, ok := r.(bool); !ok { //类型验证
		return false, errors.New("json \"" + name + "\" not bool.")
	} else {
		return v, nil
	}
}
func readInt(v interface{}) (int, error) {
	switch r := v.(type) {
	case float64:
		return int(r), nil
	case float32:
		return int(r), nil
	case int:
		return r, nil
	case int32:
		return int(r), nil
	case int64:
		return int(r), nil
	case uint32:
		return int(r), nil
	case uint64:
		return int(r), nil
	default:
		return 0, errors.New(fmt.Sprint(v, ":  not int value."))
	}
}
func readInt64(v interface{}) (int64, error) {
	switch r := v.(type) {
	case float64:
		return int64(r), nil
	case float32:
		return int64(r), nil
	case int:
		return int64(r), nil
	case int32:
		return int64(r), nil
	case int64:
		return r, nil
	case uint32:
		return int64(r), nil
	case uint64:
		return int64(r), nil
	default:
		return 0, errors.New(fmt.Sprint(v, ":  not int value."))
	}
}
func readFloat(v interface{}) (float64, error) {
	switch r := v.(type) {
	case float64:
		return r, nil
	case float32:
		return float64(r), nil
	case int:
		return float64(r), nil
	case int32:
		return float64(r), nil
	case int64:
		return float64(r), nil
	case uint32:
		return float64(r), nil
	case uint64:
		return float64(r), nil
	default:
		return 0.0, errors.New(fmt.Sprint(v, ":  not float value."))
	}
}
func jsonReadInt(jsonData map[string]interface{}, name string) (int, error) {
	t, ok := jsonData[name]
	if !ok {
		return 0, errors.New("Has no " + name)
	}
	return readInt(t)
}
func jsonReadInt64(jsonData map[string]interface{}, name string) (int64, error) {
	t, ok := jsonData[name]
	if !ok {
		return 0, errors.New("Has no " + name)
	}
	return readInt64(t)
}
func jsonReadFloat(jsonData map[string]interface{}, name string) (float64, error) {
	t, ok := jsonData[name]
	if !ok {
		return 0.0, errors.New("Has no " + name)
	}
	switch r := t.(type) {
	case float64:
		return t.(float64), nil
	case float32:
		return float64(r), nil
	case int:
		return float64(r), nil
	case int32:
		return float64(r), nil
	case int64:
		return float64(r), nil
	case uint32:
		return float64(r), nil
	case uint64:
		return float64(r), nil
	default:
		return 0.0, errors.New(fmt.Sprint(name+": ", t, " not float value."))
	}
}
func jsonToString(jsonData map[string]interface{}, name string) (string, error) {
	if r, ok := jsonData[name]; !ok {
		return "", errors.New("Has no " + name)
	} else if v, ok := r.(string); ok {
		return v, nil
	} else {
		switch t := r.(type) {
		case float64:
			return fmt.Sprintf("%d", int64(t)), nil
		case float32:
			return fmt.Sprintf("%d", int64(t)), nil
		case int:
			return fmt.Sprintf("%d", int64(t)), nil
		case int32:
			return fmt.Sprintf("%d", int64(t)), nil
		case int64:
			return fmt.Sprintf("%d", t), nil
		case uint32:
			return fmt.Sprintf("%d", int64(t)), nil
		case uint64:
			return fmt.Sprintf("%d", t), nil
		default:
			return "", errors.New(fmt.Sprint(name+": ", t, " not string or int value."))
		}
	}
}
func parseMethod(method string) (string, string) {
	array := strings.Split(method, "::")
	arrayLen := len(array)
	if arrayLen > 1 {
		classRet := array[0]
		for i := 1; i < arrayLen-1; i++ {
			classRet = classRet + "::" + array[i]
		}
		return classRet, array[arrayLen-1]
	}
	array = strings.Split(method, ".")
	arrayLen = len(array)
	if arrayLen == 1 {
		return "", method
	}
	classRet := array[0]
	for i := 1; i < arrayLen-1; i++ {
		classRet = classRet + "." + array[i]
	}
	return classRet, array[arrayLen-1]
}

//go:nosplit
func getExt(name string, token byte) string {
	for i := len(name); i > 0; i-- {
		if name[i-1] == token {
			return name[i:]
		}
	}
	return name
}

//go:nosplit
func fileExtName(name string) string {
	for i := len(name); i > 0; i-- {
		if name[i-1] == '.' {
			return name[i:]
		}
	}
	return ""
}
func isBool(value string) bool {
	return tystring.CaseCMP(value, "true") == 0 || tystring.CaseCMP(value, "false") == 0
}
func isNumber(value string) bool {
	for i := 0; i < len(value); i++ {
		v := value[i] - '0'
		if v < 0 || v > 9 {
			return false
		}
	}
	return true
}
func toBool(value string) (bool, error) {
	if !isBool(value) {
		return false, errors.New("not boolean")
	}
	return tystring.CaseCMP(value, "true") == 0, nil
}

func fileExist(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}
func readLink(path string) string {
	if link, err := os.Readlink(path); err != nil {
		return path
	} else {
		if strings.Index(link, "/") == 0 {
			return link
		}
		if pathIndex := strings.LastIndex(path, "/"); pathIndex >= 0 {
			return string(path[0:pathIndex]) + "/" + link
		}
		return link
	}
}

//go:nosplit
func isSpace(t byte) bool {
	return (t <= 0x20) && (t > 0)
}

//go:nosplit
func isHex(ch byte) bool {
	return ch-'0' < 10 || ch-'a' < 6 || ch-'A' < 6
}
func trimString(value string, isSep func(byte) bool) string {
	return tystring.TrimString(value, isSep)
}
func trimSpace(value string) string {
	return trimString(value, isSpace)
}
func fileCat(filename string) string {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(bytes)
}
func getContainerID() string {
	bytes, _ := ioutil.ReadFile("/proc/self/cgroup")
	lines := strings.Split(string(bytes), "\n")
	for _, line := range lines {
		parts := strings.Split(trimSpace(line), ":")
		groupPath := ""
		if len(parts) > 0 {
			groupPath = parts[len(parts)-1]
		}
		parts = strings.Split(groupPath, "/")
		groupName := ""
		for _, part := range parts {
			if v := trimSpace(part); len(v) > 0 {
				groupName = v
			}
		}
		if fileExtName(groupName) == "scope" {
			groupName = groupName[0 : len(groupName)-6]
		}
		if offset := len(groupName); offset >= 64 {
			for offset > 0 && isHex(groupName[offset-1]) {
				offset--
			}
			if len(groupName)-offset == 64 {
				return groupName[offset:]
			}
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
func parseEnv(envline string) (string, string) {
	for id := 0; id < len(envline); id++ {
		if envline[id] == '=' {
			return tystring.SubString(envline, 0, id), tystring.SubString(envline, id+1, len(envline))
		}
	}
	return envline, ""
}
func parsePort(addr string) (uint16, error) {
	for id := len(addr); id > 0; id-- {
		if addr[id-1] == ':' {
			addr = tystring.SubString(addr, id, len(addr))
			break
		}
	}
	if len(addr) == 0 {
		return 0, errors.New("no port listen")
	}
	v, e := strconv.Atoi(addr)
	if e != nil {
		return 0, e
	}
	if v > 65535 || v < 0 {
		return 0, errors.New("port out of bound")
	}
	return uint16(v), nil
}
func caseCMP(a, b string) int {
	return tystring.CaseCMP(a, b)
}
func subString(str string, begin, size int) string {
	return tystring.SubString(str, begin, size)
}
func parseHost(url string) string {
	for id := 0; id < len(url); id++ {
		if !tystring.IsAlpha(url[id]) {
			if tystring.SubString(url, id, 3) == "://" {
				url = tystring.SubString(url, id+3, len(url))
			}
			break
		}
	}
	for id := 0; id < len(url); id++ {
		if url[id] == '/' {
			return tystring.SubString(url, 0, id)
		}
	}
	return url
}
func parseURI(url string) string {
	for id := 0; id < len(url); id++ {
		if !tystring.IsAlpha(url[id]) {
			if tystring.SubString(url, id, 3) == "://" {
				url = tystring.SubString(url, id+3, len(url))
			}
			break
		}
	}
	for id := 0; id < len(url); id++ {
		if url[id] == '/' || url[id] == '?' {
			return tystring.SubString(url, id, len(url))
		}
	}
	return ""
}
func parseUriRequest(uri string) string {
	for id := 0; id < len(uri); id++ {
		if uri[id] == '?' {
			return tystring.SubString(uri, 0, id)
		}
	}
	return uri
}
func parseQueryString(uri string) string {
	for id := 0; id < len(uri); id++ {
		if uri[id] == '?' {
			return tystring.SubString(uri, id+1, len(uri))
		}
	}
	return ""
}
func splitMapString(source string, isSep func(byte) bool, handler func(string, string)) {
	tystring.SplitMapString(source, isSep, handler)
}

func splitStrings(source string, handler func(string) bool, isSep func(byte) bool) {
	tystring.SplitStrings(source, isSep, handler)
}
func limitAppend(a, b []byte, size int) []byte {
	if a == nil {
		if b == nil {
			return nil
		}
		if len(b) <= size {
			return b
		}
		return b[0:size]
	}
	leftSize := len(a)
	if leftSize >= size {
		return a[0:size]
	}
	if b == nil {
		return a
	}
	size = size - leftSize
	return append(a, b[0:size]...)
}
