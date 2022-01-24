// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"sync"
)

// Unit is The GoRoutine LocalStorage Unit
type Unit struct {
	lock  sync.RWMutex
	items map[int64]interface{}
}

func (u *Unit) init() {
	u.items = nil
}
func (u *Unit) Exec(gid int64, update func(value interface{}) interface{}) {
	u.lock.Lock()
	defer u.lock.Unlock()
	if u.items == nil {
		if v := update(nil); v != nil {
			u.items = map[int64]interface{}{gid: v}
		}
	} else {
		if v, found := u.items[gid]; !found {
			if v = update(nil); v != nil {
				u.items[gid] = v
			}
		} else {
			if n := update(v); n == nil {
				delete(u.items, gid)
				if len(u.items) == 0 {
					u.items = nil
				}
			} else if n != v {
				u.items[gid] = n
			}
		}
	}
}

const (
	// StorageNodes : The GoRoutine LocalStorage ArrayCount
	localStorageNodes = 32768
)

var storages [localStorageNodes]Unit

func init() {
	for i := 0; i < localStorageNodes; i++ {
		storages[i].init()
	}
}

// Get is Return the goroutine local storage value
func routineLocalExec(update func(value interface{}) interface{}) {
	gid := GetGID()
	storages[gid%localStorageNodes].Exec(gid, update)
}

// Remove is clean the goroutine local storage
func routineLocalRemove() interface{} {
	var r interface{} = nil
	routineLocalExec(func(local interface{}) interface{} {
		r = local
		return nil
	})
	return r
}

// Clear Routine Local Storage
func LocalClear() {
	routineLocalRemove()
}

// RoutineLocal : 事务线程局部存储对象
type RoutineLocal struct {
	action    *Action
	component *Component
	pointers  map[int]interface{}
}
