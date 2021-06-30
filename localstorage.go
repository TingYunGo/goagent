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
	u.items = make(map[int64]interface{})
}
func (u *Unit) get(gid int64) interface{} {
	u.lock.RLock()
	defer u.lock.RUnlock()
	if v, found := u.items[gid]; found {
		return v
	}
	return nil
}
func (u *Unit) set(gid int64, local interface{}) {
	u.lock.Lock()
	defer u.lock.Unlock()
	if u.items == nil {
		u.init()
	}
	u.items[gid] = local
}
func (u *Unit) remove(gid int64) interface{} {
	u.lock.Lock()
	defer u.lock.Unlock()
	if v, found := u.items[gid]; found {
		delete(u.items, gid)
		return v
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
