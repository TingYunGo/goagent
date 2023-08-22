// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// +build linux
// +build amd64 arm64
// +build cgo

package zaplib

import (
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/TingYunGo/goagent"
)

//func (ce *zapcore.CheckedEntry) Write(fields ...zapcore.Field) {

//go:noinline
func zapcoreCheckedEntryWrite(ce *zapcore.CheckedEntry, fields ...zapcore.Field) {
	trampoline.arg1 = *trampoline.idpointer + trampoline.idindex + trampoline.arg1 + trampoline.arg2 + trampoline.arg3 + trampoline.arg4 + trampoline.arg5 + trampoline.arg6 + trampoline.arg7 +
		trampoline.arg8 + trampoline.arg9 + trampoline.arg10 + trampoline.arg11 + trampoline.arg12 + trampoline.arg13 + trampoline.arg14 + trampoline.arg15 + trampoline.arg16 +
		trampoline.arg17 + trampoline.arg18 + trampoline.arg19 + trampoline.arg20
}

//go:noinline
func WrapzapcoreCheckedEntryWrite(ce *zapcore.CheckedEntry, fields ...zapcore.Field) {
	tokenName, tokenValue := tingyun3.GetTrackToken()
	if len(tokenName) > 0 {
		fields = append(fields, zap.String(tokenName, tokenValue))
	}
	zapcoreCheckedEntryWrite(ce, fields...)
}

func init() {
	tingyun3.Register(reflect.ValueOf(WrapzapcoreCheckedEntryWrite).Pointer())
}
