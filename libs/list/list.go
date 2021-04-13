// Copyright 2017-2021 冯立强 fenglq@tingyun.com.  All rights reserved.

//Package list : a list
package list

type node struct {
	pre   *node
	nxt   *node
	value interface{}
}

func (n *node) reset() {
	n.pre = nil
	n.nxt = nil
	n.value = nil
}

//Iterator : iterator of list
type Iterator struct {
	root *List
	n    *node
}

//IsEnd : iterator target is nil
func (i Iterator) IsEnd() bool {
	return i.n == nil
}

//Set : set the node value
func (i Iterator) Set(value interface{}) bool {
	if i.n == nil {
		return false
	}
	i.n.value = value
	return true
}

//Valid : check iterator valid
func (i Iterator) Valid() bool {
	if i.root == nil {
		return false
	}
	if i.n == nil {
		return true
	}
	if i.root.count == 0 {
		return false
	}
	if i.n.pre == nil {
		return i.root.first == i.n
	}
	if i.n.nxt == nil {
		return i.root.last == i.n
	}
	return true
}

//Destroy iterator
func (i Iterator) Destroy() {
	i.root = nil
	i.n = nil
}

//Remove node from list
func (i Iterator) Remove() (interface{}, bool) {
	if !i.Valid() {
		return nil, false
	}
	if i.root.count == 0 {
		return nil, false
	}
	defer i.Destroy()
	if i.n.pre == nil {
		return i.root.PopFront()
	}
	if i.n.nxt == nil {
		return i.root.PopBack()
	}
	i.n.pre.nxt = i.n.nxt
	i.n.nxt.pre = i.n.pre
	i.root.count--
	ret := i.n.value
	i.n.reset()
	return ret, true
}

//InsertFront : insert a node front of current
func (i Iterator) InsertFront(value interface{}) (Iterator, bool) {
	if i.root.count == 0 {
		return i.root.PushFront(value), true
	}
	if !i.Valid() || i.IsEnd() {
		return Iterator{i.root, nil}, false
	}
	if i.n.pre == nil {
		return i.root.PushFront(value), true
	}
	n := &node{i.n.pre, i.n, value}
	n.pre.nxt = n
	i.n.pre = n
	i.root.count++
	return Iterator{i.root, n}, true
}

//InsertBack : insert a node back of current
func (i Iterator) InsertBack(value interface{}) (Iterator, bool) {
	if i.root.count == 0 {
		return i.root.PushBack(value), true
	}
	if !i.Valid() || i.IsEnd() {
		return Iterator{i.root, nil}, false
	}
	if i.n.nxt == nil {
		return i.root.PushBack(value), true
	}
	n := &node{i.n, i.n.nxt, value}
	n.nxt.pre = n
	i.n.nxt = n
	i.root.count++
	return Iterator{i.root, n}, true
}

//Value : node value
func (i Iterator) Value() (interface{}, bool) {
	if i.n == nil {
		return nil, false
	}
	return i.n.value, true
}

//Equal : is same iterator
func (i Iterator) Equal(it Iterator) bool {
	return i.n == it.n && i.root == it.root
}

//MoveBack : move iterator to next node
func (i *Iterator) MoveBack() {
	if i.root != nil && i.n != nil {
		i.n = i.n.nxt
	}
}

//MoveFront : move iterator to pre node
func (i *Iterator) MoveFront() {
	if i.root != nil && i.n != nil {
		i.n = i.n.pre
	}
}

//Front : return pre node
func (i Iterator) Front() Iterator {
	if i.n == nil {
		return Iterator{i.root, nil}
	}
	return Iterator{i.root, i.n.pre}
}

//Back : return next node
func (i Iterator) Back() Iterator {
	if i.n == nil {
		return Iterator{i.root, nil}
	}
	return Iterator{i.root, i.n.nxt}
}

//List : a normal list
type List struct {
	first *node
	last  *node
	count int
}

//Init list
func (l *List) Init() *List {
	l.first = nil
	l.last = nil
	l.count = 0
	return l
}

//Size : node counts
func (l *List) Size() int {
	return l.count
}

//Front : pre node
func (l *List) Front() Iterator {
	return Iterator{l, l.first}
}

//Back : next node
func (l *List) Back() Iterator {
	return Iterator{l, l.last}
}

//PushBack : push node in list back
func (l *List) PushBack(value interface{}) Iterator {
	var n *node
	if l.count == 0 {
		n = &node{nil, nil, value}
		l.first = n
		l.last = n
	} else {
		n = &node{l.last, nil, value}
		l.last.nxt = n
		l.last = n
	}
	l.count++
	return Iterator{l, n}
}

//PushFront : push node in list front
func (l *List) PushFront(value interface{}) Iterator {
	var n *node
	if l.count == 0 {
		n = &node{nil, nil, value}
		l.first = n
		l.last = n
	} else {
		n = &node{nil, l.first, value}
		l.first.pre = n
		l.first = n
	}
	l.count++
	return Iterator{l, n}
}

//PopFront : pop node from front
func (l *List) PopFront() (interface{}, bool) {
	if l.count == 0 {
		return nil, false
	}
	rnode := l.first
	ret := rnode.value
	l.first = rnode.nxt
	rnode.value = nil
	l.count--
	if l.count == 0 {
		l.last = nil
	} else {
		l.first.pre = nil
		rnode.nxt = nil
	}
	return ret, true
}

//PopBack : pop nod from back
func (l *List) PopBack() (interface{}, bool) {
	if l.count == 0 {
		return nil, false
	}
	rnode := l.last
	ret := rnode.value
	l.last = rnode.pre
	rnode.value = nil
	l.count--
	if l.count == 0 {
		l.first = nil
	} else {
		l.last.nxt = nil
		rnode.pre = nil
	}
	return ret, true
}
