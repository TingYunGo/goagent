// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.
// Order Map

package odmap

// Iterator : the map's iterator
type Iterator struct {
	node *TreeNode
}

// IsEnd check node is last
func (it Iterator) IsEnd() bool {
	return it.node == nil
}

// Next return current right node
func (it Iterator) Next() Iterator {
	if it.node != nil {
		return Iterator{it.node.Next()}
	}
	return Iterator{nil}
}

// Pre return current left node
func (it Iterator) Pre() Iterator {
	if it.node != nil {
		return Iterator{it.node.Pre()}
	}
	return Iterator{nil}
}

// Value return node data
func (it Iterator) Value() *Paire {
	if it.node == nil {
		return nil
	}
	return it.node.Get()
}

// Map by rbtree
type Map struct {
	_Tree RbTree
	_Size uint64
}

// Size return nodes
func (m *Map) Size() uint64 {
	return m._Size
}

// Init is the map constructor
func (m *Map) Init(less func(a, b interface{}) bool) *Map {
	m._Tree.Init(less)
	m._Size = 0
	return m
}

// Clear remove all nodes
func (m *Map) Clear() {
	for m._Size > 0 {
		m.Erase(m.Begin())
	}
}

// Begin return left node
func (m *Map) Begin() Iterator {
	if m._Size > 0 {
		return Iterator{m._Tree.Begin()}
	}
	return Iterator{nil}
}

// Rbegin return right node
func (m *Map) Rbegin() Iterator {
	if m._Size > 0 {
		return Iterator{m._Tree.Rbegin()}
	}
	return Iterator{nil}
}

// End return empty node, for cmpaire
// Instead by Iterator.IsEnd
func (m *Map) End() Iterator {
	return Iterator{nil}
}

// Erase is remove a node
func (m *Map) Erase(it Iterator) {
	if it.node != nil && it.node.valid {
		m._Tree.Remove(it.node)
		m._Size--
	}
}

// Remove is remove node by key
func (m *Map) Remove(key interface{}) {
	m.Erase(m.Find(key))
}

// Set is add a k:v paire into map
func (m *Map) Set(key, value interface{}) Iterator {
	isParent, node := m._Tree.Find(key)
	if !isParent && node != nil {
		node.Get().Value = value
		return Iterator{node}
	}
	newNode := (&TreeNode{}).init(Red)
	newNode.value.first = key
	newNode.value.Value = value
	m._Tree.Insert(node, newNode)
	m._Size++
	return Iterator{newNode}
}

// Find is find node by key
func (m *Map) Find(key interface{}) Iterator {
	if isParent, node := m._Tree.Find(key); (!isParent) && (node != nil) {
		return Iterator{node}
	}
	return Iterator{nil}
}
