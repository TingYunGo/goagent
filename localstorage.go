// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package tingyun3

import (
	"sync"

	odmap "git.codemonky.net/TingYunGo/goagent/libs/map"
)

// Unit is The GoRoutine LocalStorage Unit
type Unit struct {
	lock  sync.RWMutex
	items odmap.Map
}

func (u *Unit) init() {
	u.items.Init(func(a, b interface{}) bool {
		return a.(int64) < b.(int64)
	})
}
func (u *Unit) get(gid int64) interface{} {
	u.lock.RLock()
	defer u.lock.RUnlock()
	if iterator := u.items.Find(gid); iterator != u.items.End() {
		return iterator.Value().Value
	}
	return nil
}
func (u *Unit) set(gid int64, local interface{}) {
	u.lock.Lock()
	defer u.lock.Unlock()
	u.items.Set(gid, local)
}
func (u *Unit) remove(gid int64) interface{} {
	u.lock.Lock()
	defer u.lock.Unlock()
	if iterator := u.items.Find(gid); iterator != u.items.End() {
		r := iterator.Value().Value
		u.items.Erase(iterator)
		return r
	}
	return nil
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
func routineLocalGet() interface{} {
	gid := GetGID()
	return storages[gid%localStorageNodes].get(gid)
}

// Set is Set the gorouteine local storage value
func routineLocalSet(local interface{}) {
	gid := GetGID()
	storages[gid%localStorageNodes].set(gid, local)
}

// Remove is clean the goroutine local storage
func routineLocalRemove() interface{} {
	gid := GetGID()
	return storages[gid%localStorageNodes].remove(gid)
}

// RoutineLocal : 事务线程局部存储对象
type RoutineLocal struct {
	action    *Action
	component *Component
	pointers  map[int]interface{}
}
