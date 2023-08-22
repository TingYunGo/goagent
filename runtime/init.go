// Copyright 2023 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64 arm64
// +build cgo

package wrapruntime

import (
	"reflect"

	"github.com/TingYunGo/goagent"
)

func init() {
	tingyun3.Register(reflect.ValueOf(Wrapchansend).Pointer())
	tingyun3.Register(reflect.ValueOf(Wrapchanrecv).Pointer())
	tingyun3.Register(reflect.ValueOf(Wrapselectgo).Pointer())
	tingyun3.Register(reflect.ValueOf(WrapTypeString).Pointer())
}
