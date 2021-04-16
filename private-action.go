// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"strings"
	"sync/atomic"
	"time"

	"github.com/TingYunGo/goagent/libs/tystring"
)

func (a *Action) setError(e interface{}, errType string, skipStack int) {
	if a == nil || a.stateUsed != actionUsing {
		return
	} //errorTrace 聚合,以 callstack + message
	errTime := time.Now()
	a.errors.Put(&errInfo{errTime, e, callStack(skipStack), errType})
}
func (a *Action) makeTracerID() int32 {
	return atomic.AddInt32(&a.tracerIDMaker, 1)
}

func (a *Action) unicID() string {
	if len(a.actionID) == 0 {
		a.actionID = unicID(a.time.begin, a)
	}
	return a.actionID
}
func (a *Action) getTransactionID() string {
	if _, transactionID := a.parseTrackID(); len(transactionID) > 0 {
		return transactionID
	}
	return a.unicID()
}

func (a *Action) parseTrackID() (callList, transactionID string) {
	callList, transactionID = "", ""
	if parts := strings.Split(a.trackID, ";"); len(parts) > 0 {
		for _, v := range parts {
			if tystring.SubString(v, 0, 2) == "c=" {
				callList = tystring.SubString(v, 2, len(v)-2)
			} else if tystring.SubString(v, 0, 2) == "x=" {
				transactionID = tystring.SubString(v, 2, len(v)-2)
			}
		}
	}
	return
}
