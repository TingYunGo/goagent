// Copyright 2016-2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package pool

import (
	"sync/atomic"
)

//SerialReadPool 无锁消息池，多写,有序单读(按写入时间排序)。适用于多个消息发送者,单个消息接收者模式
type SerialReadPool struct {
	//由于atomic.AddUint64 在非8字节对齐的地址上会崩溃，所以这两个id必须放在开始位置
	writeID    uint64
	readID     uint64
	cache      map[uint64]interface{}
	p          Pool
	cacheCount int32
}
type msg struct {
	id uint64
	o  interface{}
}

//SerialNew Create a SerialReadPool
func SerialNew() *SerialReadPool {
	return (&SerialReadPool{}).Init()
}

//Init : init SerialReadPool
func (p *SerialReadPool) Init() *SerialReadPool {
	p.writeID = 0
	p.readID = 1
	p.cache = make(map[uint64]interface{})
	p.p.Init()
	p.cacheCount = 0
	return p
}

//Size : Count of message in pool
func (p *SerialReadPool) Size() int32 {
	return p.p.Size() + p.cacheCount
}
func (p *SerialReadPool) cacheGet() interface{} {
	if r, exist := p.cache[p.readID]; exist {
		delete(p.cache, p.readID)
		atomic.AddInt32(&p.cacheCount, -1)
		p.readID++
		return r
	}
	return nil
}

//Get Take out a message
func (p *SerialReadPool) Get() interface{} {
	if r := p.cacheGet(); r != nil { //cache里有可用的数据,从cache里取
		return r
	}
	for {
		u := p.p.Get()
		if u == nil { //底层的pool里没有message
			return nil
		}
		m := u.(*msg)
		r := m.o
		m.o = nil
		id := m.id
		if id == p.readID { //是最早入队的message
			p.readID++
			return r
		}
		//不是最早的那个message，扔到cache里
		atomic.AddInt32(&p.cacheCount, 1)
		p.cache[id] = r
	}
}

//Put a message in pool
func (p *SerialReadPool) Put(o interface{}) {
	id := atomic.AddUint64(&(p.writeID), 1)
	p.p.Put(&msg{id, o})
}

//ForEach Peek messages
func (p *SerialReadPool) ForEach(cb func(v interface{})) {
	p.p.ForEach(cb)
}
