package list

import (
	"fmt"
	"iter"
)

// DoublyLinked represents a doubly linked list.
// The zero value is an empty list ready to use.
type DoublyLinked[T any] struct {
	len        int
	head, tail *DoublyLinkedNode[T]
}

// DoublyLinkedNode represents an element within a List.
// Its linking pointers are unexported to protect list invariants during O(1) operations.
type DoublyLinkedNode[T any] struct {
	prev, next *DoublyLinkedNode[T]
	list       *DoublyLinked[T]
	Value      T
}

// Next returns the next node in the list, or nil if n is the tail or detached.
func (n *DoublyLinkedNode[T]) Next() *DoublyLinkedNode[T] {
	if n == nil || n.list == nil {
		return nil
	}
	return n.next
}

// Prev returns the previous node in the list, or nil if n is the head or detached.
func (n *DoublyLinkedNode[T]) Prev() *DoublyLinkedNode[T] {
	if n == nil || n.list == nil {
		return nil
	}
	return n.prev
}

// New returns a new List initialized with vals.
func New[T any](vals ...T) *DoublyLinked[T] {
	l := &DoublyLinked[T]{}
	for _, v := range vals {
		l.Append(v)
	}
	return l
}

// NewWithLen returns a new List populated with initLen zero-valued elements.
func NewWithLen[T any](initLen int) *DoublyLinked[T] {
	if initLen < 0 {
		panic(fmt.Sprintf("invalid list length %d", initLen))
	}
	l := &DoublyLinked[T]{}
	var zero T
	for range initLen {
		l.Append(zero)
	}
	return l
}

// Len returns the number of elements in the list.
func (l *DoublyLinked[T]) Len() int {
	return l.len
}

// Head returns the first node in the list, or nil if the list is empty.
func (l *DoublyLinked[T]) Head() *DoublyLinkedNode[T] {
	return l.head
}

// Tail returns the last node in the list, or nil if the list is empty.
func (l *DoublyLinked[T]) Tail() *DoublyLinkedNode[T] {
	return l.tail
}

// Get returns the value of the node at i.
// Negative indices count from the end of the list, with -1 representing the tail.
func (l *DoublyLinked[T]) Get(i int) T {
	l.checkBounds(i)
	return l.nodeAt(i).Value
}

// GetNode returns the *Node[T] handle at i for later O(1) operations.
// Negative indices count from the end of the list, with -1 representing the tail.
func (l *DoublyLinked[T]) GetNode(i int) *DoublyLinkedNode[T] {
	l.checkBounds(i)
	return l.nodeAt(i)
}

// Set updates the value of the node at i to t and returns its previous value.
// Negative indices count from the end of the list, with -1 representing the tail.
func (l *DoublyLinked[T]) Set(i int, t T) T {
	l.checkBounds(i)
	n := l.nodeAt(i)
	prev := n.Value
	n.Value = t
	return prev
}

// Append inserts t at the end of the list and returns its node handle.
func (l *DoublyLinked[T]) Append(t T) *DoublyLinkedNode[T] {
	n := &DoublyLinkedNode[T]{Value: t}
	l.insertAfter(n, l.tail)
	return n
}

// Prepend inserts t at the beginning of the list and returns its node handle.
func (l *DoublyLinked[T]) Prepend(t T) *DoublyLinkedNode[T] {
	n := &DoublyLinkedNode[T]{Value: t}
	l.insertAfter(n, nil)
	return n
}

// Push inserts t at the end of the list and returns its node handle.
func (l *DoublyLinked[T]) Push(t T) *DoublyLinkedNode[T] {
	return l.Append(t)
}

// Add inserts t before index i if i >= 0, or after index i if i < 0.
// If i == l.Len(), t is appended to the end; if i == -(l.Len() + 1), t is prepended to the start.
func (l *DoublyLinked[T]) Add(t T, i int) *DoublyLinkedNode[T] {
	if l.len == 0 && (i == 0 || i == -1) {
		return l.Prepend(t)
	}
	if i == l.len || i == -1 {
		return l.Append(t)
	}
	if i == 0 || i == -(l.len+1) {
		return l.Prepend(t)
	}

	l.checkBounds(i)
	at := l.nodeAt(i)
	n := &DoublyLinkedNode[T]{Value: t}

	if i >= 0 {
		l.insertAfter(n, at.prev)
	} else {
		l.insertAfter(n, at)
	}
	return n
}

// Remove deletes the element at i, shifting subsequent elements left, and returns its value.
func (l *DoublyLinked[T]) Remove(i int) T {
	l.checkBounds(i)
	return l.RemoveNode(l.nodeAt(i))
}

// Pop removes and returns the last element in the list.
func (l *DoublyLinked[T]) Pop() T {
	return l.Remove(-1)
}

// RemoveNode detaches n from the list in O(1) time and returns its value.
func (l *DoublyLinked[T]) RemoveNode(n *DoublyLinkedNode[T]) T {
	l.checkNode(n)
	l.unlink(n)

	val := n.Value
	var zero T
	n.Value = zero
	return val
}

// MoveToFront moves n to the head of the list in O(1) time.
func (l *DoublyLinked[T]) MoveToFront(n *DoublyLinkedNode[T]) {
	l.checkNode(n)
	if l.head == n {
		return
	}
	l.unlink(n)
	l.insertAfter(n, nil)
}

// MoveToBack moves n to the tail of the list in O(1) time.
func (l *DoublyLinked[T]) MoveToBack(n *DoublyLinkedNode[T]) {
	l.checkNode(n)
	if l.tail == n {
		return
	}
	l.unlink(n)
	l.insertAfter(n, l.tail)
}

// MoveAfter moves n immediately after mark in O(1) time.
func (l *DoublyLinked[T]) MoveAfter(n, mark *DoublyLinkedNode[T]) {
	l.checkNode(n)
	l.checkNode(mark)
	if n == mark || mark.next == n {
		return
	}
	l.unlink(n)
	l.insertAfter(n, mark)
}

// MoveBefore moves n immediately before mark in O(1) time.
func (l *DoublyLinked[T]) MoveBefore(n, mark *DoublyLinkedNode[T]) {
	l.checkNode(n)
	l.checkNode(mark)
	if n == mark || n.next == mark {
		return
	}
	prev := mark.prev
	l.unlink(n)
	l.insertAfter(n, prev)
}

// Swap exchanges the positions of the nodes at indices k and j, preserving node identity.
func (l *DoublyLinked[T]) Swap(k, j int) {
	l.checkBounds(k)
	l.checkBounds(j)
	if k == j {
		return
	}
	l.SwapNodes(l.nodeAt(k), l.nodeAt(j))
}

// SwapNodes exchanges the positions of nodes a and b in O(1) time, preserving node identity.
func (l *DoublyLinked[T]) SwapNodes(a, b *DoublyLinkedNode[T]) {
	l.checkNode(a)
	l.checkNode(b)
	if a == b {
		return
	}

	if a.next == b {
		l.unlink(a)
		l.insertAfter(a, b)
		return
	}
	if b.next == a {
		l.unlink(b)
		l.insertAfter(b, a)
		return
	}

	aPrev := a.prev
	bPrev := b.prev
	l.unlink(a)
	l.unlink(b)
	l.insertAfter(a, bPrev)
	l.insertAfter(b, aPrev)
}

// ForEach iterates from head to tail, calling f(i, value) with i in [0, len-1].
func (l *DoublyLinked[T]) ForEach(f func(i int, t T)) {
	node := l.head
	for i := 0; i < l.len; i, node = i+1, node.next {
		f(i, node.Value)
	}
}

// ReverseForEach iterates from tail to head, calling f(i, value) with i in [-1, -len].
func (l *DoublyLinked[T]) ReverseForEach(f func(i int, t T)) {
	node := l.tail
	for i := -1; i >= -l.len; i, node = i-1, node.prev {
		f(i, node.Value)
	}
}

// All returns an iterator over index-value pairs from head to tail.
func (l *DoublyLinked[T]) All() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		node := l.head
		for i := 0; i < l.len; i, node = i+1, node.next {
			if !yield(i, node.Value) {
				return
			}
		}
	}
}

// Backward returns an iterator over negative-index and value pairs from tail to head.
func (l *DoublyLinked[T]) Backward() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		node := l.tail
		for i := -1; i >= -l.len; i, node = i-1, node.prev {
			if !yield(i, node.Value) {
				return
			}
		}
	}
}

func (l *DoublyLinked[T]) unlink(n *DoublyLinkedNode[T]) {
	if n.prev != nil {
		n.prev.next = n.next
	} else {
		l.head = n.next
	}

	if n.next != nil {
		n.next.prev = n.prev
	} else {
		l.tail = n.prev
	}

	n.prev = nil
	n.next = nil
	n.list = nil
	l.len--
}

func (l *DoublyLinked[T]) insertAfter(n, mark *DoublyLinkedNode[T]) {
	n.list = l
	if mark == nil {
		n.prev = nil
		n.next = l.head
		if l.head != nil {
			l.head.prev = n
		} else {
			l.tail = n
		}
		l.head = n
	} else {
		n.prev = mark
		n.next = mark.next
		if mark.next != nil {
			mark.next.prev = n
		} else {
			l.tail = n
		}
		mark.next = n
	}
	l.len++
}

func (l *DoublyLinked[T]) nodeAt(i int) *DoublyLinkedNode[T] {
	if i < 0 {
		i += l.len
	}

	if i <= l.len/2 {
		node := l.head
		for range i {
			node = node.next
		}
		return node
	}

	node := l.tail
	for range l.len - 1 - i {
		node = node.prev
	}
	return node
}

func (l *DoublyLinked[T]) checkBounds(i int) {
	if i >= l.len || i < -l.len {
		panic(fmt.Sprintf("list index %d out of bounds", i))
	}
}

func (l *DoublyLinked[T]) checkNode(n *DoublyLinkedNode[T]) {
	if n == nil || n.list != l {
		panic("list: node is nil or belongs to another list")
	}
}
