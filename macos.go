// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build !linux !cgo

package tingyun3

import (
	"runtime"
	"strconv"
	"strings"
)

func GetGID() int64 {
	buffer := make([]byte, 32)
	runtime.Stack(buffer, false)
	if parts := Split(string(buffer), " "); parts[0] == "goroutine" && len(parts) > 1 {
		if goid, err := strconv.ParseInt(parts[1], 10, 0); err == nil {
			return goid
		}
	}
	return -1
}

func Split(s, sep string) []string {
	sep_len := len(sep)
	if sep_len == 0 {
		return []string{s}
	}
	count := 0
	index := 0
	for i := strings.Index(s[index:], sep); index < len(s); i = strings.Index(s[index:], sep) {
		if i != 0 {
			count++
		}
		if i < 0 {
			break
		}
		index += i + sep_len
	}
	if count == 0 {
		return []string{}
	}
	r := make([]string, count)
	index = 0
	for i := 0; i < count; i++ {
		sep_index := 0
		for sep_index = strings.Index(s[index:], sep); sep_index == 0; sep_index = strings.Index(s[index:], sep) {
			index += sep_len
		}
		if sep_index > 0 {
			r[i] = string(s[index : index+sep_index])
		} else {
			r[i] = string(s[index:])
		}
		index += sep_index
	}
	return r
}
