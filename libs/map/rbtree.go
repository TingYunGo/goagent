// Copyright 2021 冯立强 fenglq@tingyun.com.  All rights reserved.

package odmap

// Color is rbtree's colour
type Color int8

const (
	// Red colour
	Red Color = 0
	// Black colour
	Black Color = 1
)

// Paire is Node Data
type Paire struct {
	first interface{}
	Value interface{}
}

// Key method return the Node Key
func (p *Paire) Key() interface{} {
	return p.first
}

// TreeNode is node of the red/black tree
type TreeNode struct {
	Left   *TreeNode
	Right  *TreeNode
	Parent *TreeNode
	value  Paire
	color  Color
	valid  bool
}

// Get return Node Data
func (t *TreeNode) Get() *Paire {
	if t == nil {
		return nil
	}
	return &t.value
}

// isBlack is check the node color
func (t *TreeNode) isBlack() bool {
	return t.color != Red
}

// isRed is check the node color
func (t *TreeNode) isRed() bool {
	return t.color == Red
}

// set node color
func (t *TreeNode) setBlack() {
	t.color = Black
}

// set node color
func (t *TreeNode) setRed() {
	t.color = Red
}
func (t *TreeNode) init(color Color) *TreeNode {
	t.Left = nil
	t.Right = nil
	t.Parent = nil
	t.color = color
	t.valid = true
	return t
}

// copy node color
func (t *TreeNode) copyColor(other *TreeNode) {
	t.color = other.color
}

// swap node color
func (t *TreeNode) swapColor(other *TreeNode) {
	_color := other.color
	other.color = t.color
	t.color = _color
}

// Pre is left node of the current
func (t *TreeNode) Pre() *TreeNode {
	ret := t.Left
	if ret != nil {
		for ret.Right != nil {
			ret = ret.Right
		}
		return ret
	}
	top := t
	for ret = t.Parent; ret != nil; {
		if ret.Right == top {
			return ret
		}
		top = ret
		ret = top.Parent
	}
	return nil
}

// Next is right node of the current
func (t *TreeNode) Next() *TreeNode {
	top := t
	ret := t.Right
	if ret != nil {
		for ret.Left != nil {
			ret = ret.Left
		}
		return ret
	}

	for ret = top.Parent; ret != nil; {
		if ret.Left == top {
			return ret
		}
		top = ret
		ret = top.Parent
	}
	return nil
}

// RbTree is Red/Black tree
type RbTree struct {
	less func(a, b interface{}) bool
	Root *TreeNode
}

// Init is struct initializer
func (rbt *RbTree) Init(less func(a, b interface{}) bool) *RbTree {
	rbt.Root = nil
	rbt.less = less
	return rbt
}

func (rbt *RbTree) _Rotate(current *TreeNode) {
	for {
		parent := current.Parent
		if parent == nil {
			current.setBlack()
			break
		}
		if parent.isBlack() {
			break
		}
		grandparent := parent.Parent
		uncle := grandparent.Left
		if uncle == parent {
			uncle = grandparent.Right
		}
		if uncle != nil && uncle.isRed() {
			parent.setBlack()
			uncle.setBlack()
			current = grandparent
			current.setRed()
			continue
		}
		if grandparent.Left == parent && parent.Right == current {
			grandparent.Left = current
			current.Parent = grandparent
			parent.Parent = current
			parent.Right = current.Left
			if parent.Right != nil {
				parent.Right.Parent = parent
			}
			current.Left = parent
			current = parent
			continue
		}
		if grandparent.Right == parent && parent.Left == current {
			grandparent.Right = current
			current.Parent = grandparent
			parent.Parent = current
			parent.Left = current.Right
			if parent.Left != nil {
				parent.Left.Parent = parent
			}
			current.Right = parent
			current = parent
			continue
		}
		gg := grandparent.Parent
		if gg == nil {
			rbt.Root = parent
		} else if gg.Left == grandparent {
			gg.Left = parent
		} else {
			gg.Right = parent
		}
		parent.Parent = gg
		parent.setBlack()
		grandparent.setRed()
		if current == parent.Left {
			grandparent.Left = parent.Right
			if parent.Right != nil {
				parent.Right.Parent = grandparent
			}
			grandparent.Parent = parent
			parent.Right = grandparent
		} else {
			grandparent.Right = parent.Left
			if parent.Left != nil {
				parent.Left.Parent = grandparent
			}
			grandparent.Parent = parent
			parent.Left = grandparent
		}
		break
	}
}

func _swapNode(a, b *TreeNode) {
	replaceParent := func(x, a, b *TreeNode) {
		if x != nil {
			if x.Left == a {
				x.Left = b
			} else {
				x.Right = b
			}
		}
	}
	x := a.Parent
	a.Parent = b.Parent
	replaceParent(x, a, b)
	b.Parent = x
	replaceParent(a.Parent, b, a)
	x = a.Left
	a.Left = b.Left
	b.Left = x
	if x != nil {
		x.Parent = b
	}
	if a.Left != nil {
		a.Left.Parent = a
	}
	x = a.Right
	a.Right = b.Right
	b.Right = x
	if x != nil {
		x.Parent = b
	}
	if a.Right != nil {
		a.Right.Parent = a
	}
	a.swapColor(b)
}

func (rbt *RbTree) _removeOne(node *TreeNode) {
	child := node.Left
	if child == nil {
		child = node.Right
	}
	parent := node.Parent
	if parent == nil {
		if child == nil {
			if node.isBlack() {
				rbt.Root = nil
			}
		} else {
			child.setBlack()
			child.Parent = nil
			rbt.Root = child
		}
		return
	}
	_TempNode := (&TreeNode{}).init(Black)

	if node.isBlack() && child == nil {
		child = _TempNode //借一个子节点
	}
	if parent.Left == node {
		parent.Left = child
	} else {
		parent.Right = child
	}
	if node.isRed() {
		return
	}
	child.Parent = parent
	if node.isBlack() {
		if child.isRed() {
			child.setBlack()
		} else {
			rbt._removeCaseN(child)
		}
	}
	if child == _TempNode { //还回借来的子节点
		parent = child.Parent
		if parent.Left == child {
			parent.Left = nil
		} else {
			parent.Right = nil
		}
		child.Parent = nil
	}
}

func (rbt *RbTree) _rotateLeft(node *TreeNode) {
	parent := node.Parent
	s := node.Right
	s.Parent = parent
	node.Right = s.Left
	if node.Right != nil {
		node.Right.Parent = node
	}
	node.Parent = s
	s.Left = node
	if parent != nil {
		if parent.Left == node {
			parent.Left = s
		} else {
			parent.Right = s
		}
	} else {
		rbt.Root = s
	}
}
func (rbt *RbTree) _rotateRight(node *TreeNode) {
	parent := node.Parent
	s := node.Left
	s.Parent = parent
	node.Left = s.Right
	if node.Left != nil {
		node.Left.Parent = node
	}
	node.Parent = s
	s.Right = node
	if parent != nil {
		if parent.Left == node {
			parent.Left = s
		} else {
			parent.Right = s
		}
	} else {
		rbt.Root = s
	}
}

func (rbt *RbTree) _removeCaseN(n *TreeNode) {
	isBlack := func(p *TreeNode) bool {
		if p == nil {
			return true
		}
		return p.isBlack()
	}
	parent := n.Parent
	var s *TreeNode = nil
	for {
		if parent == nil {
			break
		}
		if parent.Left == n {
			s = parent.Right
		} else {
			s = parent.Left
		}
		if s.isRed() {
			parent.setRed()
			s.setBlack()
			if parent.Left == n {
				rbt._rotateLeft(parent)
			} else {
				rbt._rotateRight(parent)
			}
			if parent.Left == n {
				s = parent.Right
			} else {
				s = parent.Left
			}
		}
		if parent.isBlack() && s.isBlack() && isBlack(s.Left) && isBlack(s.Right) {
			s.setRed()
			n = parent
			parent = n.Parent
			continue
		}
		if parent.isRed() && s.isBlack() && isBlack(s.Left) && isBlack(s.Right) {
			s.setRed()
			parent.setBlack()
			break
		}
		if s.isBlack() {
			if n == parent.Left && isBlack(s.Right) && !isBlack(s.Left) {
				s.setRed()
				s.Left.setBlack()
				rbt._rotateRight(s)
			} else if n == parent.Right && isBlack(s.Left) && !isBlack(s.Right) {
				s.setRed()
				s.Right.setBlack()
				rbt._rotateLeft(s)
			}
			if parent.Left == n {
				s = parent.Right
			} else {
				s = parent.Left
			}
		}
		s.copyColor(parent)
		parent.setBlack()
		if n == parent.Left {
			s.Right.setBlack()
			rbt._rotateLeft(parent)
		} else {
			s.Left.setBlack()
			rbt._rotateRight(parent)
		}
		break
	}
}

// Find : get node on key
func (rbt *RbTree) Find(key interface{}) (isParent bool, r *TreeNode) {
	node := rbt.Root
	if node != nil {
		for {
			if rbt.less(key, node.value.first) {
				if node.Left == nil {
					break
				} else {
					node = node.Left
				}
			} else if rbt.less(node.value.first, key) {
				if node.Right == nil {
					break
				} else {
					node = node.Right
				}
			} else {
				return false, node
			}
		}
	}
	return true, node
}

// Insert : add node
func (rbt *RbTree) Insert(parent, node *TreeNode) {
	node.Parent = parent
	node.Left = nil
	node.Right = nil
	node.valid = true
	if parent != nil {
		node.setRed()
	} else {
		node.setBlack()
	}
	if parent == nil {
		rbt.Root = node
	} else {
		if rbt.less(node.value.first, parent.value.first) {
			parent.Left = node
		} else {
			parent.Right = node
		}
		rbt._Rotate(node)
	}
}

// Remove is remove node
func (rbt *RbTree) Remove(node *TreeNode) {
	if !node.valid {
		return
	}
	if node.Right != nil {
		next := node.Right
		for next.Left != nil {
			next = next.Left
		}
		_swapNode(node, next)
		if next.Parent == nil {
			rbt.Root = next
		}
	}
	rbt._removeOne(node)
	node.valid = false
}

// Begin return the left node
func (rbt *RbTree) Begin() *TreeNode {
	ret := rbt.Root
	if ret != nil {
		for ret.Left != nil {
			ret = ret.Left
		}
	}
	return ret
}

// Rbegin return the right node
func (rbt *RbTree) Rbegin() *TreeNode {
	ret := rbt.Root
	if ret != nil {
		for ret.Right != nil {
			ret = ret.Right
		}
	}
	return ret
}
