// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"strings"
	"sync/atomic"
	"time"
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

func getAppSecID(id string) string {
	sidArray := strings.Split(id, "|")
	if len(sidArray) < 2 {
		return ""
	}
	trackID := sidArray[1]
	for i := 2; i < len(sidArray); i++ {
		trackID += "|" + sidArray[i]
	}
	return trackID

}
func getTxID(id string) string {
	array := strings.Split(id, ";")
	if len(array) < 4 {
		return ""
	}
	for i := 0; i < len(array); i++ {
		paire := strings.Split(array[i], "=")
		if paire[0] == "x" {
			if len(paire) < 2 {
				return ""
			}
			xid := paire[1]
			for i := 2; i < len(paire); i++ {
				xid += "=" + paire[i]
			}
			return xid
		}
	}
	return ""
}
func getTopMetric(id string) string {
	array := strings.Split(id, ";")
	if len(array) < 4 {
		return ""
	}
	for i := 1; i < len(array); i++ {
		paire := strings.Split(array[i], "=")
		if paire[0] == "p" {
			if len(paire) < 2 {
				return ""
			}
			protocol := paire[1]
			for i := 2; i < len(paire); i++ {
				protocol += "=" + paire[i]
			}
			return "EntryTransaction/" + protocol + "/" + getAppSecID(array[0])
		}
	}
	return ""
}

func (a *Action) unicID() string {
	txID := getTxID(a.trackID)
	if txID == "" {
		return unicID(a.time.begin, a)
	}
	return txID
}
