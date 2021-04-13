// Copyright 2016-2021 冯立强 fenglq@tingyun.com.  All rights reserved.

//Package pool 无锁消息池，多读多写, 用于goroutine 间收发消息
package pool

//接口: Put, Get, Size

import "sync/atomic"

type node struct {
	next  *node
	value interface{}
}
type lockLIST struct {
	lock int32
	head *node
	end  *node
}

func (l *lockLIST) init() {
	l.lock = 0
	l.head = nil
	l.end = nil
}

func (l *lockLIST) PushBack(n *node) bool {
	used := atomic.AddInt32(&l.lock, 1)
	defer atomic.AddInt32(&l.lock, -1)
	if used != 1 {
		return false
	}
	n.next = nil
	if l.end == nil {
		l.head = n
		l.end = n
	} else {
		l.end.next = n
		l.end = n
	}
	return true
}
func (l *lockLIST) ForEach(cb func(v interface{})) {
	used := atomic.AddInt32(&l.lock, 1)
	defer atomic.AddInt32(&l.lock, -1)
	if used != 1 || l.head == nil {
		return
	}
	for node := l.head; node != nil; node = node.next {
		cb(node.value)
	}
}
func (l *lockLIST) PopFront() *node {
	if l.head == nil {
		return nil
	}
	used := atomic.AddInt32(&l.lock, 1)
	defer atomic.AddInt32(&l.lock, -1)
	if used != 1 {
		return nil
	}
	if l.head == nil {
		return nil
	}
	ret := l.head
	l.head = ret.next
	ret.next = nil
	if l.head == nil {
		l.end = nil
	}
	return ret
}

const bucketCount = 8

type nodePool struct {
	array      [bucketCount]lockLIST
	count      int32
	indexRead  int32
	indexWrite int32
}

func (p *nodePool) init() *nodePool {
	p.count = 0
	p.indexRead = 0
	p.indexWrite = 0
	for i := 0; i < bucketCount; i++ {
		p.array[i].init()
	}
	return p
}

func (p *nodePool) Put(n *node) {
	pwrite := p.indexWrite
	for {
		for i := pwrite; i-pwrite < bucketCount; i++ {
			listID := i % bucketCount
			if p.array[listID].PushBack(n) {
				atomic.AddInt32(&p.count, 1)
				p.indexWrite = listID + 1
				return
			}
		}
	}
}

func (p *nodePool) Size() int32 {
	return p.count
}
func (p *nodePool) ForEach(cb func(v interface{})) {
	if p.count == 0 {
		return
	}
	pread := p.indexRead
	for i := pread; i-pread < bucketCount; i++ {
		readlistID := i % bucketCount
		p.array[readlistID].ForEach(cb)
	}
}

func (p *nodePool) Get() *node {
	if p.count == 0 {
		return nil
	}
	pread := p.indexRead
	for i := pread; i-pread < bucketCount; i++ {
		readlistID := i % bucketCount
		r := p.array[readlistID].PopFront()
		if r != nil {
			atomic.AddInt32(&p.count, -1)
			p.indexRead = readlistID + 1
			return r
		}
	}
	return nil
}

//Pool None lock Message pool
type Pool struct {
	pool nodePool
}

//Init : Pool init.
func (p *Pool) Init() *Pool {
	p.pool.init()
	return p
}

//New : Create a pool
func New() *Pool {
	return new(Pool).Init()
}

//Put a message in pool
func (p *Pool) Put(v interface{}) {
	p.pool.Put(&node{next: nil, value: v})
}

//Size : Count off message in pool
func (p *Pool) Size() int32 {
	return p.pool.Size()
}

//ForEach : Peek messages
func (p *Pool) ForEach(cb func(v interface{})) {
	p.pool.ForEach(cb)
}

//Get : Take out a message
func (p *Pool) Get() interface{} {
	n := p.pool.Get()
	if n == nil {
		return nil
	}
	ret := n.value
	n.value = nil
	return ret
}
